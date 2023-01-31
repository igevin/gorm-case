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
	size := 100
	users := createUsers(size)
	// 传入slice，则进行批量创建
	// 这里传入&users 或 user 都可以
	res := c.db.Create(&users)
	c.Assert().NoError(res.Error)
	for i := 0; i < size; i++ {
		c.Assert().Equal(uint(i+1), users[i].ID)
	}
}

func (c *createSuite) TestCreatePerBatch() {
	size, batchSize := 100, 5
	users := createUsers(size)
	// 可以控制每个批次的大小
	res := c.db.CreateInBatches(&users, batchSize)
	c.Assert().NoError(res.Error)
	for i := 0; i < size; i++ {
		c.Assert().Equal(uint(i+1), users[i].ID)
	}
}

func (c *createSuite) TestCreatePerBatchV2() {
	size, batchSize := 100, 5
	users := createUsers(size)
	// 可以在创建db的时候，就配置好每个批次的大小：
	// db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{
	//  CreateBatchSize: 1000,
	//})
	// 也可以配置到session 中去：
	db := c.db.Session(&gorm.Session{CreateBatchSize: batchSize})
	res := db.Create(&users)
	c.Assert().NoError(res.Error)
	for i := 0; i < size; i++ {
		c.Assert().Equal(uint(i+1), users[i].ID)
	}
}

func (c *createSuite) TestCreateByMap() {
	var id uint = 1000
	// 创建单个记录
	row := map[string]any{
		"id": id,
		// key不区分大小写，都会对应到column
		"name":     "Tom",
		"Birthday": time.Now(),
	}
	res := c.db.Model(&User{}).Create(row)
	c.Assert().NoError(res.Error)
	u := &User{Model: gorm.Model{ID: id}}
	res = c.db.First(u)
	c.Require().NoError(res.Error)
	c.Assert().Equal(id, u.ID)
	c.Assert().Zero(u.Age)

	// 批量创建
	rows := []User{
		*newUser(2001),
		*newUser(2002),
	}
	res = c.db.Model(&User{}).Create(rows)
	c.Require().NoError(res.Error)

	rows2 := []*User{
		newUser(3001),
		newUser(3002),
	}
	res = c.db.Model(&User{}).Create(rows2)
	c.Require().NoError(res.Error)

	// 也可以传入 rows3 的指针
	rows3 := []*User{
		newUser(4001),
		newUser(4002),
	}
	res = c.db.Model(&User{}).Create(&rows3)
	c.Require().NoError(res.Error)
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

func createUsers(size int) []User {
	res := make([]User, 0, size)
	for i := 0; i < size; i++ {
		res = append(res, *newUser(0))
	}
	return res
}
