package basic

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3" // 一定不要忘记导入驱动
)

func CreateDb() (*sql.DB, error) {
	driverName := "sqlite3"
	// dsn: dataSourceName，这里是创建了一个基于sqlite3的内存数据库
	dsn := "file:test.db?cache=shared&mode=memory"
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
	result := NewResult(res, err)
	return result.RowsAffected()
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
	//if r.err != nil {
	//	return 0, nil
	//}
	//lastId, err := r.res.LastInsertId()
	//if err != nil {
	//	return 0, err
	//}
	//return lastId, nil
	return r.handleResult(r.res.LastInsertId)
}

func (r *Result) RowsAffected() (int64, error) {
	//if r.err != nil {
	//	return 0, nil
	//}
	//affected, err := r.res.RowsAffected()
	//if err != nil {
	//	return 0, err
	//}
	//return affected, nil
	return r.handleResult(r.res.RowsAffected)
}

// handleResult 进一步重构，提炼重复代码
func (r *Result) handleResult(f func() (int64, error)) (int64, error) {
	if r.err != nil {
		return 0, nil
	}
	affected, err := f()
	if err != nil {
		return 0, err
	}
	return affected, nil
}
