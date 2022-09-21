package model

import (
	"fmt"
	"github.com/eininst/scheduler/configs"
)

var DefaultTablePrefix = "cron_"

type User struct {
	Id         int64  `gorm:"primary_key" json:"id"`
	Name       string `json:"name" gorm:"type:varchar(100);unique;default:''"`
	Password   string `json:"password" gorm:"type:varchar(100);default:''"`
	RealName   string `json:"realName" gorm:"type:varchar(32);default:''"`
	Role       string `json:"role" gorm:"type:varchar(32);default:''"`
	Head       string `json:"head" gorm:"type:varchar(200);default:''"`
	Mail       string `json:"mail" gorm:"type:varchar(128);default:''"`
	Status     string `json:"status" gorm:"type:varchar(32);default:''"`
	CreateTime string `json:"createTime" gorm:"type:varchar(32);default:''"`
}

func (this *User) TableName() string {
	tablePrefix := configs.Get("tablePrefix").String()
	if tablePrefix == "" {
		tablePrefix = DefaultTablePrefix
	}
	return fmt.Sprintf("%s%s", tablePrefix, "user")
}

type Task struct {
	Id          int64  `gorm:"primary_key" json:"id"`
	UserId      int64  `json:"userId" gorm:"type:bigint(20);index;default:0"`
	Name        string `json:"name" gorm:"type:varchar(100);unique;default:''"`
	Group       string `json:"group" gorm:"type:varchar(100);index;default:''"`
	Spec        string `json:"spec" gorm:"type:varchar(64);default:''"`
	Url         string `json:"url" gorm:"type:varchar(256);default:''"`
	Method      string `json:"method" gorm:"type:varchar(10);default:''"`
	ContentType string `json:"contentType" gorm:"type:varchar(64);default:''"`
	Body        string `json:"body" gorm:"type:varchar(1024);default:''"`
	Timeout     int64  `json:"timeout" gorm:"type:int;default:0"`
	MaxRetries  int    `json:"maxRetries" gorm:"type:int;default:0"`
	Desc        string `json:"desc" gorm:"type:varchar(2048);default:''"`
	Status      string `json:"status" gorm:"type:varchar(32);index;default:''"`
	CreateTime  string `json:"createTime" gorm:"type:varchar(32);index;default:''"`
}

func (this *Task) TableName() string {
	tablePrefix := configs.Get("tablePrefix").String()
	if tablePrefix == "" {
		tablePrefix = DefaultTablePrefix
	}
	return fmt.Sprintf("%s%s", tablePrefix, "task")
}

type TaskExcute struct {
	Id         int64  `gorm:"primary_key" json:"id"`
	UserId     int64  `json:"userId" gorm:"type:bigint(20);index;default:0"`
	TaskId     int64  `json:"taskId" gorm:"type:bigint(20);index;default:0"`
	TaskName   string `json:"taskName" gorm:"type:varchar(100);index;default:''"`
	TaskGroup  string `json:"taskGroup" gorm:"type:varchar(100);index;default:''"`
	TaskUrl    string `json:"taskUrl" gorm:"type:longtext;"`
	TaskObj    string `json:"taskObj" gorm:"type:varchar(2048);default:''"`
	Code       int    `json:"code" gorm:"type:int;index;default:0"`
	Response   string `json:"response" gorm:"type:longtext;"`
	StartTime  string `json:"start_time" gorm:"type:varchar(32);index;default:''"`
	EndTime    string `json:"end_time" gorm:"type:varchar(32);index;default:''"`
	Duration   int64  `json:"duration" gorm:"type:int;index;default:0"`
	CreateTime string `json:"createTime" gorm:"type:varchar(32);index;default:''"`
}

func (this *TaskExcute) TableName() string {
	tablePrefix := configs.Get("tablePrefix").String()
	if tablePrefix == "" {
		tablePrefix = DefaultTablePrefix
	}
	return fmt.Sprintf("%s%s", tablePrefix, "task_excute")
}
