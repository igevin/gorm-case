package main

import (
	"fmt"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func crud(db *gorm.DB) {
	// Create
	create(db)

	var p, p1 Product

	// Read
	fmt.Println("Reading...")
	read(db, &p, &p1)

	// Update
	fmt.Println("Updating...")
	update(db, &p, &p1)

	// Delete
	fmt.Println("Deleting...")
	del(db, &p)

	// 永久删除
	deletePermanently(db, &p)

}

func create(db *gorm.DB) {
	fmt.Println("Creating...")
	db.Create(&Product{Code: "G001", Price: 100})
}

func read(db *gorm.DB, p *Product, p1 *Product) {
	// p 可以为零值，也可以有属性值
	// 若 p 有属性值(如p2)，则属性值为查询条件
	db.First(p)
	fmt.Println("p: ", p)
	p2 := Product{Model: gorm.Model{ID: p.ID}}
	db.First(&p2)
	fmt.Println("p2: ", p2)

	// 指定ID
	var p0 Product
	db.First(&p0, p.ID)
	fmt.Println("p0: ", p0)

	// 传入查询条件
	db.First(p1, "Code=?", "G001")
	fmt.Println("p1: ", p1)

	// 总结：三种查询方式
	//db.First(p)
	//db.First(p, 1)
	//db.First(p, "code=?", "001")
}

func update(db *gorm.DB, p, p1 *Product) {
	// 更新单个对象，用这种方式：
	//db.Model(p).Update()
	// &p 是二级指针，可以用
	db.Model(&p).Update("Price", 200)
	db.First(&p)
	fmt.Println("p: ", *p)
	fmt.Println("p1, before update: ", *p1)
	// p1 是一级指针，可以用
	db.Model(p1).Updates(Product{Code: "G001-1", Price: 300})
	db.First(p1)
	fmt.Println("p1, after update: ", *p1)

	// 不用指针会报错，如：
	//p2 := Product{
	//	Model: p1.Model,
	//	Code:  p1.Code,
	//	Price: p1.Price,
	//}
	//db.Model(p2).Updates(Product{Code: "G001-2", Price: 300})
	//db.First(p2)
	//fmt.Println("p2: ", p2)
}

func del(db *gorm.DB, p *Product) {
	// 删除：
	//db.Delete(p)

	// If your model includes a gorm.DeletedAt field (which is included in gorm.Model),
	// it will get soft delete ability automatically!
	db.Delete(&p)
	fmt.Println("Deleted")

	// 查询被软删除的记录
	db.Unscoped().First(&p)
	fmt.Println("soft deleted data: ", p)
}

func deletePermanently(db *gorm.DB, p *Product) {
	fmt.Println("永久删除...")
	db.Unscoped().Delete(p)
	fmt.Printf("永久删除成功：%v\n", *p)

	// 不需要先删除，再永久删除
	db.Create(&Product{
		Model: gorm.Model{ID: p.ID},
		Code:  p.Code,
		Price: p.Price,
	})
	var p1 Product
	db.First(&p1, p.ID)
	fmt.Println("new created p1: ", p1)
	db.Unscoped().Delete(p1)
	fmt.Println("直接永久删除 p1 成功")
}
