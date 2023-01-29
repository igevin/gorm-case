package basic

import (
	"context"
	"database/sql"
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

func CreateTable(db *sql.DB, ctx context.Context, createSql string) (int64, error) {
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

func InsertRow(db *sql.DB, ctx context.Context, insertSql string, data ...any) (int64, error) {
	r, err := db.ExecContext(ctx, insertSql, data...)
	res := NewResult(r, err)
	affected, err := res.RowsAffected()
	if err != nil || affected == 0 {
		return 0, err
	}
	return res.LastInsertId()
}

// QueryRow 查询结果与query语句是耦合的，这里没法解耦，只能写死，否则有过度设计之嫌
func QueryRow(db *sql.DB, ctx context.Context, data ...any) (TestModel, error) {
	query := "SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?"
	row := db.QueryRowContext(ctx, query, data...)
	if row.Err() != nil {
		return TestModel{}, row.Err()
	}
	tm := TestModel{}
	err := row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	return tm, err
}

// QueryRows 查询结果与query语句是耦合的，这里没法解耦，只能写死，否则有过度设计之嫌
func QueryRows(db *sql.DB, ctx context.Context, data ...any) ([]TestModel, error) {
	query := "SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?"
	rows, err := db.QueryContext(ctx, query, data...)
	if err != nil {
		return nil, err
	}
	res := make([]TestModel, 0)
	for rows.Next() {
		if rows.Err() != nil {
			continue
		}
		tm := TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		if err != nil {
			continue
		}
		res = append(res, tm)
	}
	return res, err
}
