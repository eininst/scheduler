package types

type Page[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
	Extra any   `json:"extra"`
}

type PageReq struct {
	PageSize int `json:"pageSize" form:"pageSize"`
	Current  int `json:"current" form:"current"`
}
type TaskOption struct {
	*PageReq
	UserId int64  `json:"userId" form:"userId"`
	Name   string `json:"name" form:"name"`
	Group  string `json:"group" form:"group"`
	Status string `json:"status" form:"status"`
	Sort   string `json:"sort" form:"sort"`
	Dir    string `json:"dir" form:"dir"`
}

type TaskDTO struct {
	Id           int64  `gorm:"primary_key" json:"id"`
	UserId       int64  `json:"userId"`
	UserName     string `json:"userName"`
	UserRealName string `json:"userRealName"`
	UserHead     string `json:"userHead"`
	UserMail     string `json:"userMail"`
	Name         string `json:"name"`
	Group        string `json:"group"`
	Cron         string `json:"cron"`
	Url          string `json:"url"`
	Method       string `json:"method"`
	ContentType  string `json:"contentType"`
	Body         string `json:"body"`
	Timeout      int64  `json:"timeout"`
	MaxRetries   int    `json:"maxRetries"`
	Desc         string `json:"desc"`
	Status       string `json:"status"`
	CreateTime   string `json:"createTime"`
}

type TaskChangeUser struct {
	UserId  int64   `json:"userId" form:"userId"`
	TaskIds []int64 `json:"taskIds" form:"taskIds"`
}

type TaskBatch struct {
	TaskIds []int64 `json:"taskIds" form:"taskIds"`
}

type TaskExcuteOption struct {
	*PageReq
	UserId    int64  `json:"userId" form:"userId"`
	TaskId    int64  `json:"taskId" form:"taskId"`
	TaskName  string `json:"taskName" form:"taskName"`
	TaskGroup string `json:"taskGroup" form:"taskGroup"`
	Status    string `json:"status" form:"status"`
	Duration  int64  `json:"duration" form:"duration"`
	StartTime string `json:"start_time" form:"start_time"`
	EndTime   string `json:"end_time" form:"end_time"`
	Sort      string `json:"sort" form:"sort"`
	Dir       string `json:"dir" form:"dir"`
}

type TaskExcuteDTO struct {
	Id           int64  `gorm:"primary_key" json:"id"`
	UserId       int64  `json:"userId" `
	UserName     string `json:"userName" `
	UserRealName string `json:"userRealName" `
	UserHead     string `json:"userHead" `
	TaskId       int64  `json:"taskId" `
	TaskName     string `json:"taskName"`
	TaskGroup    string `json:"taskGroup"`
	TaskUrl      string `json:"taskUrl"`
	TaskObj      string `json:"taskObj"`
	Code         int    `json:"code"`
	Response     string `json:"response"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
	Duration     int64  `json:"duration"`
	CreateTime   string `json:"createTime"`
}

type UserOption struct {
	Name     string `json:"name" form:"name"`
	RealName string `json:"realName" form:"realName"`
	Role     string `json:"role" form:"role"`
	Status   string `json:"status" form:"status"`
	Sort     string `json:"sort" form:"sort"`
	Dir      string `json:"dir" form:"dir"`
}

type DashboardChart struct {
	Date  string `json:"date"`
	Code  int64  `json:"code"`
	Count int64  `json:"count"`
}

type Dashboard struct {
	TaskCount    int64 `json:"taskCount"`
	TaskRunCount int64 `json:"runCount"`

	SchedulerCount int64             `json:"schedulerCount"`
	StartTime      string            `json:"startTime"`
	Chart          []*DashboardChart `json:"chart"`
}
