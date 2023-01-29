package basic

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3" // 一定不要忘记导入驱动
)

func CreateDb() (*sql.DB, error) {
	driverName := "sqlite3"
	// dsn: dataSourceName，这里是创建了一个基于sqlite3的内存数据库
	dsn := "file:standard.db?cache=shared&mode=memory"
	// 创建一个普通的sqlite3 数据库
	//dsn := "file:standard.db"
	return sql.Open(driverName, dsn)
}

func CreateTable(db *sql.DB, ctx context.Context) (int64, error) {
	createSql := `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`
	res, err := db.ExecContext(ctx, createSql)
	// 这部分代码可以再封装提炼出去
	//if err != nil {
	//	return 0, nil
	//}
	//affected, err := res.RowsAffected()
	//if err != nil {
	//	return 0, err
	//}
	//return affected, nil
	result := NewResult(res, err)
	return result.RowsAffected()
}

func InsertRow(db *sql.DB, ctx context.Context, data ...any) (int64, error) {
	s := "INSERT INTO test_model(`id`, `first_name`, `age`, `last_name`) VALUES(?, ?, ?, ?)"
	r, err := db.ExecContext(ctx, s, data...)
	//if err != nil {
	//	return 0, err
	//}
	res := NewResult(r, err)
	affected, err := res.RowsAffected()
	if err != nil || affected == 0 {
		return 0, err
	}
	return res.LastInsertId()
}

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
