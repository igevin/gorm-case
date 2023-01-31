package blog

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type Writer struct {
	Id       int64
	Username string
}

// Value 实现了 driver.Valuer 接口
// 可以把数据值转换为数据库的列数据
func (w *Writer) Value() (driver.Value, error) {
	return json.Marshal(w)
}

// Scan 实现了 sql.Scanner 接口
// 调用 Rows.Scan() 接口时，可以把数据库的列数据 转换成结构体的属性值
func (w *Writer) Scan(src any) error {
	var bs []byte
	switch val := src.(type) {
	case []byte:
		bs = val
	case *[]byte:
		bs = *val
	case string:
		bs = []byte(val)
	case sql.RawBytes:
		bs = val
	case *sql.RawBytes:
		bs = *val
	default:
		return fmt.Errorf("ekit：JsonColumn.Scan 不支持 src 类型 %v", src)
	}
	return json.Unmarshal(bs, w)
}

type Blog struct {
	Id      int64
	Title   string
	Content string
	// 这里用指针，才能正常利用 Scanner 和 Valuer 接口，让自定义类型与数据库如期交互数据
	Author     *Writer
	CreateTime time.Time
	UpdateTime time.Time
}
