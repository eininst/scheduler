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
}
