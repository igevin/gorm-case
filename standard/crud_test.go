package standard

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateDb(t *testing.T) {
	_ = createDbForTest(t)
}

func TestCreateTable(t *testing.T) {
	db := createDbForTest(t)
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	createTableForTest(t, db, ctx)
}

func createDbForTest(t *testing.T) *sql.DB {
	db, err := CreateDb()
	require.NoError(t, err)
	return db
}

func createTableForTest(t *testing.T, db *sql.DB, ctx context.Context) {
	createSql := `
CREATE TABLE IF NOT EXISTS test_model(
   id INTEGER PRIMARY KEY,
   first_name TEXT NOT NULL,
   age INTEGER,
   last_name TEXT NOT NULL
)
`
	affected, err := CreateTable(db, ctx, createSql)
	require.NoError(t, err)
	assert.Equal(t, int64(0), affected)
}

func TestCrud(t *testing.T) {
	db := createDbForTest(t)
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	createTableForTest(t, db, ctx)

	testInsertRow(t, db, ctx)
	testQueryRow(t, db, ctx)
	testQueryRows(t, db, ctx)
	testInsertRowTimeout(t, db, ctx)
}

func testInsertRow(t *testing.T, db *sql.DB, ctx context.Context) {
	s := "INSERT INTO test_model(`id`, `first_name`, `age`, `last_name`) VALUES(?, ?, ?, ?)"
	id, err := InsertRow(db, ctx, s, 1, "Tom", 18, "Jerry")
	require.NoError(t, err)
	assert.Equal(t, int64(1), id)
}

func testInsertRowTimeout(t *testing.T, db *sql.DB, ctx context.Context) {
	time.Sleep(time.Second * 2)
	s := "INSERT INTO test_model(`id`, `first_name`, `age`, `last_name`) VALUES(?, ?, ?, ?)"
	_, err := InsertRow(db, ctx, s, 2, "Tom", 18, "Jerry")
	assert.Equal(t, context.DeadlineExceeded, err)
}

func testQueryRow(t *testing.T, db *sql.DB, ctx context.Context) {
	row, err := QueryRow(db, ctx, 1)
	require.NoError(t, err)
	t.Log(row, row.LastName.String)
}

func testQueryRows(t *testing.T, db *sql.DB, ctx context.Context) {
	rows, err := QueryRows(db, ctx, 1)
	require.NoError(t, err)
	t.Log(rows)
}
