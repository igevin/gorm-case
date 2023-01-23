package basic

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func TestUserCrud(t *testing.T) {
	db, err := prepare()
	assert.NoError(t, err)
	if err != nil {
		return
	}
	err = UserCrud(db)
	assert.NoError(t, err)
}

func TestUserV2Crud(t *testing.T) {
	db, err := prepare()
	assert.NoError(t, err)
	if err != nil {
		return
	}
	err = UserV2Crud(db)
	assert.NoError(t, err)
}

func TestUserV3Crud(t *testing.T) {
	db, err := prepare()
	assert.NoError(t, err)
	if err != nil {
		return
	}
	err = UserV3Crud(db)
	assert.NoError(t, err)
}

func prepare() (db *gorm.DB, err error) {
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&User{}); err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&UserV2{}); err != nil {
		return nil, err
	}
	if err = db.AutoMigrate(&UserV3{}); err != nil {
		return nil, err
	}
	return
}
