package demo

import (
	"github.com/igevin/gorm-case/standard/advanced/generic"
	"time"
)

type Writer struct {
	Id       int64
	Username string
}

type Blog struct {
	Id      int64
	Title   string
	Content string
	// 这里用指针，才能正常利用 Scanner 和 Valuer 接口，让自定义类型与数据库如期交互数据
	Author     *generic.JsonColumn[Writer]
	Tags       *generic.JsonColumn[[]string]
	CreateTime time.Time
	UpdateTime time.Time
}
