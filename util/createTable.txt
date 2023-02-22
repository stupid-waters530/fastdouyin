package util

import (
	"fmt"
	"github.com/ikuraoo/fastdouyin/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateTable() {
	//host := viper.GetString("datasource.host")
	//port := viper.GetString("datasource.port")
	//database := viper.GetString("datasource.database")
	//username := viper.GetString("datasource.username")
	//password := viper.GetString("datasource.password")
	//charset := viper.GetString("datasource.charset")
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
	//	username,
	//	password,
	//	host,
	//	port,
	//	database,
	//	charset)

	dsn := "root:root@tcp(127.0.0.1:3306)/dousheng?charset=utf8mb4&parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}

	err = db.AutoMigrate(&entity.Video{})
	if err != nil {
		fmt.Println("建表失败")
	}
}
