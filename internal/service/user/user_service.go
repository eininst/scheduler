package user

import (
	"context"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service"
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

type UserService interface {
	Add(ctx context.Context, user *model.User) error
	Update(ctx context.Context, user *model.User) error
	Enable(ctx context.Context, id int64) error
	Disable(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error

	GetById(ctx context.Context, id int64) (*model.User, error)
	Login(ctx context.Context, username string, password string) (*model.User, error)
	List(ctx context.Context) ([]*model.User, error)
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
	return &u, nil
}

func (us *userService) Add(ctx context.Context, user *model.User) error {
	sess := us.DB.WithContext(ctx)
	var checkUser model.User

	sess.First(&checkUser, "name = ?", user.Name)
	if checkUser.Id != 0 {
		return service.NewServiceError("账号名称重复")
	}

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
	sess.First(&checkUser, "name = ?", user.Name)
	if checkUser.Id != 0 {
		return service.NewServiceError("账号名称重复")
	}

	er := sess.Model(&u).Updates(&model.User{
		Name:     user.Name,
		RealName: user.RealName,
		Mail:     user.Mail,
		Role:     user.Role,
		Head:     user.Head,
		Password: user.Password,
	}).Error

	if er != nil {
		return service.NewServiceError("修改账号失败")
	}
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

func (us *userService) List(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	us.DB.WithContext(ctx).Find(&users, "status != ?", STATUS_DEL)
	return users, nil
}
