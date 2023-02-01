package generic

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

func JsonColumnScanError(src any) error {
	return fmt.Errorf("JsonColumn.Scan 不支持 src 类型 %v", src)
}

type JsonColumn[T any] struct {
	Val   T
	Valid bool
}

func NewJsonColumn[T any](obj T) JsonColumn[T] {
	return JsonColumn[T]{
		Val:   obj,
		Valid: true,
	}
}

// Scan 将 src 转化为对象
// src 的类型必须是 []byte, *[]byte, string, sql.RawBytes, *sql.RawBytes 之一
func (j *JsonColumn[T]) Scan(src any) error {
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
		//return fmt.Errorf("JsonColumn.Scan 不支持 src 类型 %v", src)
		return JsonColumnScanError(src)
	}

	if err := json.Unmarshal(bs, &j.Val); err != nil {
		return err
	}
	j.Valid = true
	return nil
}

// Value 返回一个 json string，类型是 []byte
func (j *JsonColumn[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	return json.Marshal(j.Val)
}
