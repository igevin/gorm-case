package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	// 打开db session会话
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移schema
	_ = db.AutoMigrate(&Product{})

	// Create
	fmt.Println("Creating...")
	db.Create(&Product{Code: "G001", Price: 100})

	// Read
	fmt.Println("Reading...")
	var p, p1 Product
	// get by id
	db.First(&p)
	// 指定ID
	//db.First(&p, 1)
	fmt.Println(p)
	db.First(&p1, "Code=?", "G001")
	fmt.Println(p1)
	// get by model，可以指定ID
	p2 := Product{Model: gorm.Model{ID: p.ID}}
	db.First(&p2)
	fmt.Println(p2)

	// Update
	fmt.Println("Updating...")
	db.Model(&p).Update("Price", 200)
	db.First(&p)
	fmt.Println(p)
	fmt.Println(p1)
	db.Model(&p1).Updates(Product{Code: "G001-1", Price: 300})
	db.First(&p1)
	fmt.Println(p1)

	// Delete
	fmt.Println("Deleting...")
	// If your model includes a gorm.DeletedAt field (which is included in gorm.Model),
	// it will get soft delete ability automatically!
	db.Delete(&p)
	fmt.Println("Deleted")

	// 查询被软删除的记录
	db.Unscoped().First(&p)
	fmt.Println("soft deleted: ", p)

	// 永久删除
	db.Unscoped().Delete(&p)

	// 不需要先删除，再永久删除
	db.Create(&p)
	db.Unscoped().Delete(&p)
}
