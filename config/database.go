package config

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() {
	user := GetEnv("DB_USER", "root")
	password := GetEnv("DB_PASSWORD", "")
	host := GetEnv("DB_HOST", "127.0.0.1")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/mygorm?charset=utf8mb4&parseTime=True&loc=Local", user, password, host)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Db = db
}
