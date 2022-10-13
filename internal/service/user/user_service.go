package user

import (
	"context"
	"fmt"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service"
	"github.com/eininst/scheduler/internal/types"
	"github.com/eininst/scheduler/internal/util"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

const (
	STATUS_OK       = "ok"
	STATUS_DISABLEd = "disabled"
	STATUS_DEL      = "del"
)

const (
	ROLE_ADMIN  = "admin"
	ROLE_NORMAL = "normal"
)

var SORT_MAP = map[string]string{
	"createTime": "create_time",
}

type UserService interface {
	Add(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	ResetPassword(ctx context.Context, id int64, password string) error

	Enable(ctx context.Context, id int64) error
	Disable(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error

	GetById(ctx context.Context, id int64) (*model.User, error)
	Login(ctx context.Context, username string, password string) (*model.User, error)
	List(ctx context.Context, opt *types.UserOption) ([]*model.User, error)

	Count(ctx context.Context) (int64, error)
}

type userService struct {
	Store *session.Store
	DB    *gorm.DB      `inject:""`
	Rcli  *redis.Client `inject:""`
}

func NewService() UserService {
	return &userService{}
}

func (us *userService) Login(ctx context.Context,
	username string, password string) (*model.User, error) {

	var u model.User
	us.DB.WithContext(ctx).First(&u, "name = ?", username)
	if u.Id == 0 {
		return nil, service.NewServiceError("账号或密码错误")
	}
	if u.Password != util.Md5(password) {
		return nil, service.NewServiceError("账号或密码错误")
	}
	if u.Status != STATUS_OK {
		return nil, service.NewServiceError("账号不可用")
	}
	return &u, nil
}

func (us *userService) Add(ctx context.Context, user *model.User) error {
	sess := us.DB.WithContext(ctx)
	var checkUser model.User

	sess.First(&checkUser, "name = ?", user.Name)
	if checkUser.Id != 0 {
		return service.NewServiceError("账号名称重复")
	}

	user.Password = util.Md5(user.Password)
	user.CreateTime = util.FormatTime()
	user.Status = "ok"

	er := sess.Create(&user).Error
	if er != nil {
		return service.NewServiceError("新增账号失败")
	}

	return er
}

func (us *userService) Update(ctx context.Context, user *model.User) error {
	var u model.User
	sess := us.DB.WithContext(ctx)
	sess.First(&u, user.Id)
	if u.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}

	var checkUser model.User
	sess.First(&checkUser, "name = ? and id != ?", user.Name, user.Id)
	if checkUser.Id != 0 {
		return service.NewServiceError("账号名称重复")
	}

	if u.CreateTime == "" {
		u.CreateTime = util.FormatTime()
	}

	if user.Role != checkUser.Role {
		var count int64
		sess.Model(&model.User{}).Where("role = ? and status = ?", ROLE_ADMIN, STATUS_OK).Count(&count)
		if count == 1 && user.Role != u.Role && u.Role == ROLE_ADMIN {
			return service.NewServiceError("修改失败、系统需要至少保留一个管理员权限账号")
		}
	}

	er := sess.Model(&u).Updates(&model.User{
		Name:       user.Name,
		RealName:   user.RealName,
		Mail:       user.Mail,
		Role:       user.Role,
		Head:       user.Head,
		CreateTime: u.CreateTime,
	}).Error

	if er != nil {
		return service.NewServiceError("修改账号失败")
	}
	return er
}

func (us *userService) ResetPassword(ctx context.Context, id int64, password string) error {
	var u model.User
	sess := us.DB.WithContext(ctx)
	sess.First(&u, id)
	if u.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}

	er := sess.Model(&u).Updates(&model.User{
		Password: util.Md5(password),
	}).Error
	return er
}

func (us *userService) Enable(ctx context.Context, id int64) error {
	var u model.User
	sess := us.DB.WithContext(ctx)
	sess.First(&u, id)
	if u.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}
	er := sess.Model(&u).Update("status", STATUS_OK).Error

	if er != nil {
		return service.NewServiceError("启用失败")
	}
	return er
}

func (us *userService) Disable(ctx context.Context, id int64) error {
	var u model.User
	sess := us.DB.WithContext(ctx)
	sess.First(&u, id)
	if u.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}

	var count int64
	sess.Model(&model.User{}).Where("role = ? and status = ?", ROLE_ADMIN, STATUS_OK).Count(&count)

	if count == 1 && u.Role == ROLE_ADMIN {
		return service.NewServiceError("禁用失败、系统需要至少保留一个有效的管理员权限账号")
	}

	er := sess.Model(&u).Update("status", STATUS_DISABLEd).Error

	if er != nil {
		return service.NewServiceError("禁用失败")
	}
	return er
}

func (us *userService) Delete(ctx context.Context, id int64) error {
	var u model.User
	sess := us.DB.WithContext(ctx)
	sess.First(&u, id)
	if u.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}

	var count int64
	sess.Model(&model.User{}).Where("role = ? and status = ?", ROLE_ADMIN, STATUS_OK).Count(&count)
	if count == 1 && u.Role == ROLE_ADMIN {
		return service.NewServiceError("删除失败、系统需要至少保留一个管理员权限账号")
	}

	er := sess.Model(&u).Update("status", STATUS_DEL).Error

	if er != nil {
		return service.NewServiceError("删除失败")
	}
	return er
}

func (us *userService) GetById(ctx context.Context, id int64) (*model.User, error) {
	var u model.User
	us.DB.WithContext(ctx).First(&u, id)
	if u.Id == 0 {
		return nil, service.ERROR_DATA_NOT_FOUND
	}
	return &u, nil
}

func (us *userService) List(ctx context.Context, opt *types.UserOption) ([]*model.User, error) {
	var users []*model.User

	q := us.DB.WithContext(ctx)
	if opt.Name != "" {
		q = q.Where("name like ?", "%"+opt.Name+"%")
	}
	if opt.RealName != "" {
		q = q.Where("real_name like ?", "%"+opt.RealName+"%")
	}
	if opt.Status != "" {
		q = q.Where("status = ?", opt.Status)
	}
	if opt.Role != "" {
		q = q.Where("role = ?", opt.Role)
	}
	if opt.Sort != "" {
		sort, ok := SORT_MAP[opt.Sort]
		if ok {
			if opt.Dir == "ascend" {
				q = q.Order(fmt.Sprintf("%s", sort))
			} else {
				q = q.Order(fmt.Sprintf("%s desc", sort))
			}
		}
	} else {
		q = q.Order("id desc")
	}

	q.Find(&users, "status != ?", STATUS_DEL)
	return users, nil
}

func (us *userService) Count(ctx context.Context) (int64, error) {
	var count int64

	q := us.DB.WithContext(ctx).Model(&model.User{})
	er := q.Where("status != ?", STATUS_DEL).Count(&count).Error

	return count, er
}
