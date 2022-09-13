package types

type Page[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
	Extra any   `json:"extra"`
}

type Task struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	Spec       string `json:"spec"`
	Url        string `json:"url"`
	Method     string `json:"method"`
	Body       string `json:"body"`
	Leader     string `json:"leader"`
	Timeout    int64  `json:"timeout"`
	MaxRetries int    `json:"maxRetries"`
	Mail       string `json:"mail"`
}

type TaskExcuteInfo struct {
	TaskId    string `json:"taskId"`
	TaskObj   string `json:"taskObj"`
	TaskName  string `json:"taskName"`
	TaskUrl   string `json:"taskUrl"`
	Code      int    `json:"code"`
	Response  string `json:"response"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}
