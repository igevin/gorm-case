package practice

import (
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

type createSuite struct {
	suite.Suite
	dsn string

	db *gorm.DB
}

func TestCreate(t *testing.T) {
	//dsn := "file:test.db?cache=shared&mode=memory"
	dsn := "file:create.db"
	suite.Run(t, &createSuite{dsn: dsn})
}

func (c *createSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(c.dsn), &gorm.Config{})
	c.Require().NoError(err)
	c.db = db

	c.Require().NoError(c.db.AutoMigrate(&User{}))
}

func (c *createSuite) TearDownTest() {
	c.db.Unscoped().Where("1=1").Delete(&User{})
}

func (c *createSuite) TestCreate() {
	// 常规创建
	var id uint = 1000
	u := newUser(id)
	res := c.db.Create(u)
	c.Require().NoError(res.Error)
	c.Assert().Equal(id, u.ID)
	c.Assert().Equal(int64(1), res.RowsAffected)

	// 当 id=0 或 id 未传值时，会自动创建id，id=lastId+1
	u = newUser(0)
	res = c.db.Create(u)
	c.Require().NoError(res.Error)
	c.Assert().Equal(id+1, u.ID)
}

func (c *createSuite) TestCreateSelectively() {
	var id uint = 1000
	u := newUser(id)
	// 只保存选中的字段
	res := c.db.Select("Name").Create(u)
	c.Require().NoError(res.Error)
	c.Assert().NotEqual(id, u.ID)
	c.T().Log(u)

	user := &User{}
	res = c.db.First(user, u.ID)
	c.Require().NoError(res.Error)
	c.Assert().NotEqual(u.Age, user.Age)
	c.Assert().Equal(0, user.Age)
	c.Assert().True(user.Birthday.IsZero())
}

func (c *createSuite) TestCreateNegligibly() {
	var id uint = 1000
	u := newUser(id)
	// 保存时忽略指定字段
	res := c.db.Omit("id", "name").Create(u)
	c.Assert().NoError(res.Error)
	c.Assert().NotEqual(id, u.ID)

	user := &User{}
	res = c.db.First(user, u.ID)
	c.Require().NoError(res.Error)
	c.Assert().Equal(u.Birthday.UTC(), user.Birthday.UTC())
	c.Assert().Equal(u.Age, user.Age)
	c.Assert().Equal("", user.Name)
}

func (c *createSuite) TestCreateBatch() {

}

func (c *createSuite) TestCreateBatchV2() {

}

func (c *createSuite) TestCreateByMap() {

}

func (c *createSuite) TestCreateBySqlClause() {

}

func (c *createSuite) TestCreateHandlingCustomType() {

}

func (c *createSuite) TestCreateHandlingCustomTypeWithClause() {

}

func newUser(id uint) *User {
	return &User{
		Model:    gorm.Model{ID: id},
		Name:     "Tom",
		Age:      18,
		Birthday: time.Now(),
	}
}
