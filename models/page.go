package models

type Page struct {
	PageNum int `json:"pageNum"`
	PageSize int `json:"pageSize"`
	TotalPage int `json:"totalPage"`
	TotalSize int `json:"totalSize"`
	Data interface{} `json:"data"`
}
