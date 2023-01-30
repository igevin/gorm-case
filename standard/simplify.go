package standard

import (
	"database/sql"
	"errors"
)

// Result 设计了db.Exec 和db.ExecContext返回结果的抽象
// 简化了结果返回和错误处理
type Result struct {
	err error
	res sql.Result
}

func NewResult(res sql.Result, err error) Result {
	return Result{
		err: err,
		res: res,
	}
}

func (r *Result) LastInsertId() (int64, error) {
	// LastInsertId() 和 RowsAffected() 处理逻辑类似(如下面注释掉的部分)，可以进一步提炼封装
	//if r.err != nil {
	//	return 0, r.err
	//}
	//affected, err := r.res.LastInsertId()
	//if err != nil {
	//	return 0, err
	//}
	//return affected, nil

	var f func() (int64, error)
	if r.res != nil {
		f = r.res.LastInsertId
	}
	return r.handleResult(f)
}

func (r *Result) RowsAffected() (int64, error) {
	// LastInsertId() 和 RowsAffected() 处理逻辑类似，可以进一步提炼封装
	var f func() (int64, error)
	if r.res != nil {
		f = r.res.RowsAffected
	}
	return r.handleResult(f)
}

// handleResult 封装了LastInsertId() 和 RowsAffected() 中的重复逻辑
func (r *Result) handleResult(f func() (int64, error)) (int64, error) {
	if r.err != nil {
		return 0, r.err
	}
	if f == nil {
		return 0, errors.New("handler func is nil")
	}
	affected, err := f()
	if err != nil {
		return 0, err
	}
	return affected, nil
}
