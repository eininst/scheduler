package model

type SchedulerUser struct {
	Id         int64  `gorm:"primary_key" json:"id"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	RealName   string `json:"realName"`
	Role       string `json:"role"`
	Head       string `json:"head"`
	Mail       string `json:"mail"`
	Status     string `json:"status"`
	CreateTime string `json:"createTime"`
}

type SchedulerTask struct {
	Id          int64  `gorm:"primary_key" json:"id"`
	UserId      int64  `json:"userId"`
	Name        string `json:"name"`
	Group       string `json:"group"`
	Spec        string `json:"spec"`
	Url         string `json:"url"`
	Method      string `json:"method"`
	ContentType string `json:"contentType"`
	Body        string `json:"body"`
	Timeout     int64  `json:"timeout"`
	MaxRetries  int    `json:"maxRetries"`
	Desc        string `json:"desc"`
	Status      string `json:"status"`
	CreateTime  string `json:"createTime"`
}

type SchedulerTaskExcute struct {
	Id         int64  `gorm:"primary_key" json:"id"`
	UserId     int64  `json:"userId"`
	TaskId     int64  `json:"taskId"`
	TaskName   string `json:"taskName"`
	TaskUrl    string `json:"taskUrl"`
	TaskObj    string `json:"taskObj"`
	Code       int    `json:"code"`
	Response   string `json:"response"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Duration   int64  `json:"duration"`
	CreateTime string `json:"createTime"`
}
