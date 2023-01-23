package basic

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
}

func UserCrud(db *gorm.DB) error {
	u := User{Name: "Tom"}
	db.Create(&u)
	db.First(&u)
	fmt.Println("User: ", u)
	// 只要有 `gorm.DeletedAt` 类型的DeleteAt字段，就会触发软删除
	db.Delete(&u)
	return nil
}

type UserV2 struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

func UserV2Crud(db *gorm.DB) error {
	u := UserV2{Name: "Tom"}
	db.Create(&u)
	db.First(&u)
	fmt.Println("UserV2: ", u)
	// 没有DeleteAt字段，不会软删除
	db.Delete(&u)
	return nil
}

type UserV3 struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	Name      string
}

func UserV3Crud(db *gorm.DB) error {
	u := UserV3{Name: "Tom"}
	db.Create(&u)
	db.First(&u)
	fmt.Println("UserV3: ", u)
	// 有DeleteAt字段，但不是"gorm.DeleteAt"类型，也不会软删除
	db.Delete(&u)
	return nil
}
