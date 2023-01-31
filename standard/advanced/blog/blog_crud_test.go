package blog

import (
	"context"
	"database/sql"
	"github.com/igevin/gorm-case/standard"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type blogSuite struct {
	suite.Suite
	driver string
	dsn    string

	db *sql.DB
}

func TestBlog(t *testing.T) {
	//dsn := "file:blog.db?cache=shared&mode=memory"
	dsn := "file:blog.db"
	suite.Run(t, &blogSuite{
		driver: "sqlite3",
		dsn:    dsn,
	})
}

func (b *blogSuite) SetupSuite() {
	db, err := sql.Open(b.driver, b.dsn)
	b.Require().NoError(err)
	b.db = db

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	createSql := `
	CREATE TABLE IF NOT EXISTS blog(
    id          INTEGER
        primary key,
    title       TEXT not null,
    content     TEXT not null,
    author      TEXT not null,
    create_time datetime,
    update_time datetime
)
	`
	_, err = b.db.ExecContext(ctx, createSql)
	b.Require().NoError(err)
}

func (b *blogSuite) SetupTest() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	b.createBlog(ctx, 1000)
}

func (b *blogSuite) createBlog(ctx context.Context, id int64) {
	blog := createBlog(id, createWriter(1))
	s := "INSERT INTO blog (`id`, `title`, `content`, `author`, `create_time`, `update_time`) VALUES(?,?,?,?,?,?)"
	// 由于 Writer 实现了 Scanner 和 Valuer 接口，blog.Author 可以直接通过 json 序列化为文本字符串
	res, err := b.db.ExecContext(ctx, s, blog.Id, blog.Title, blog.Content, blog.Author, blog.CreateTime, blog.UpdateTime)
	result := standard.NewResult(res, err)
	affected, err := result.RowsAffected()
	b.Require().NoError(err)
	b.Assert().Equal(int64(1), affected)
}

func (b *blogSuite) TearDownTest() {
	query := "DELETE FROM `blog`"
	_, err := b.db.Exec(query)
	b.Require().NoError(err)
}

func (b *blogSuite) TestCreateBlog() {
	var id int64 = 1
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	b.createBlog(ctx, id)
}

func (b *blogSuite) TestGetBlog() {
	var id int64 = 1000
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	s := "SELECT * FROM blog WHERE id=?"
	row := b.db.QueryRowContext(ctx, s, id)
	b.Require().NoError(row.Err())

	blog := Blog{}
	// 用指针也可以
	//blog := &Blog{}
	err := row.Scan(&blog.Id, &blog.Title, &blog.Content, &blog.Author, &blog.CreateTime, &blog.UpdateTime)
	b.Require().NoError(err)
	b.Assert().Equal(id, blog.Id)
	b.Assert().Equal("Gevin", blog.Author.Username)
}

func createWriter(id int64) *Writer {
	return &Writer{
		Id:       id,
		Username: "Gevin",
	}
}

func createBlog(id int64, w *Writer) Blog {
	return Blog{
		Id:         id,
		Title:      "Demo",
		Content:    "This is a demo",
		Author:     w,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}
}
