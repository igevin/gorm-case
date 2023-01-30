package suite

import (
	"context"
	"database/sql"
	"github.com/igevin/gorm-case/standard"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"strings"
	"testing"
	"time"
)

type sqlTestSuite struct {
	suite.Suite

	driver string
	dsn    string

	db *sql.DB
}

func TestSQLite(t *testing.T) {
	suite.Run(t, &sqlTestSuite{
		driver: "sqlite3",
		dsn:    "file:test.db?cache=shared&mode=memory",
		//dsn: "file:test.db",
	})
}

func (s *sqlTestSuite) SetupSuite() {
	db, err := sql.Open(s.driver, s.dsn)
	require.NoError(s.T(), err)
	s.db = db

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	createSql := `
CREATE TABLE IF NOT EXISTS test_model(
   id INTEGER PRIMARY KEY,
   first_name TEXT NOT NULL,
   age INTEGER,
   last_name TEXT NOT NULL
)
`
	_, err = s.db.ExecContext(ctx, createSql)
	require.NoError(s.T(), err)
}

// 每个测试执行一次
func (s *sqlTestSuite) SetupTest() {
	// 测试用例包含"query"或"delete"时，先给数据库新增一条数据
	name := strings.ToLower(s.T().Name())
	if !strings.Contains(name, "query") &&
		!strings.Contains(name, "delete") {
		return
	}
	row := s.defaultRow()
	id, err := insertRow(context.Background(), s.db, row.Id, row.FirstName, row.Age, row.LastName)
	require.NoError(s.T(), err)
	require.Equal(s.T(), int64(1000), id)
}

// 每个测试执行一次
func (s *sqlTestSuite) TearDownTest() {
	query := "DELETE FROM `test_model`"
	_, err := s.db.Exec(query)
	require.NoError(s.T(), err)
}

// 全部测试总计执行一次
func (s *sqlTestSuite) TearDownSuite() {
	query := "DELETE FROM `test_model`"
	_, err := s.db.Exec(query)
	require.NoError(s.T(), err)
}

func insertRow(ctx context.Context, db *sql.DB, data ...any) (int64, error) {
	query := "INSERT INTO test_model(`id`, `first_name`, `age`, `last_name`) VALUES(?, ?, ?, ?)"
	return standard.InsertRow(db, ctx, query, data...)
}

func (s *sqlTestSuite) TestInsert() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	id, err := insertRow(ctx, s.db, 1, "Tom", 18, "Jerry")
	require.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), id)
}

func (s *sqlTestSuite) TestInsertTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	time.Sleep(time.Millisecond * 2)
	_, err := insertRow(ctx, s.db, 1, "Tom", 18, "Jerry")
	assert.Equal(s.T(), context.DeadlineExceeded, err)
}

func (s *sqlTestSuite) TestQueryRow() {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	expected := s.defaultRow()

	row, err := standard.QueryRow(s.db, ctx, 1000)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), expected, row)
}

func (s *sqlTestSuite) TestQueryRows() {
	row := s.defaultRow()
	res, err := standard.QueryRows(s.db, context.Background(), row.Id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), 1, len(res))
}

func (s *sqlTestSuite) TestDeleteRow() {
	row := s.defaultRow()

	affected, err := standard.DeleteRows(s.db, context.Background(), row.Id)
	require.NoError(s.T(), err)
	require.Equal(s.T(), int64(1), affected)
}

func (s *sqlTestSuite) defaultRow() standard.TestModel {
	return standard.TestModel{
		Id:        1000,
		FirstName: "Tom",
		Age:       18,
		LastName: &sql.NullString{
			String: "Jerry",
			Valid:  true,
		},
	}
}
