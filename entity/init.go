package entity

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ikuraoo/fastdouyin/constant"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var db *gorm.DB

func Init() error {
	var err error
	err = InitDB()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	err = RedisInit()
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return err
}

func InitDB() error {
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)
	//dsn := "root:root@tcp(127.0.0.1:3306)/dousheng?charset=utf8mb4&parseTime=true"
	fmt.Println(dsn)
	var err error
	db, err = gorm.Open(
		mysql.Open(dsn),
		&gorm.Config{},
	)
	return err
}

func RedisInit() (err error) {
	constant.REDIS = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host") + ":" + viper.GetString("redis.port"),
		Password: viper.GetString("redis.password"), // no password set
		DB:       viper.GetInt("redis.db"),          // use default DB
		PoolSize: viper.GetInt("redis.poolsize"),    // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = constant.REDIS.Ping(ctx).Result()
	return err
}
