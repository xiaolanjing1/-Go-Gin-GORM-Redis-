package config

import (
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() {
	dsn := "root:your@tcp(127.0.0.1:3306)/mygorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Db = db
}
