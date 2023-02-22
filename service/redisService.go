package service

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ikuraoo/fastdouyin/constant"
	"github.com/ikuraoo/fastdouyin/entity"
	"math/rand"
	"strconv"
	"time"
)

func AddUserInfoByUserIDToRedis(user *entity.User) error {
	// 定义 key
	userKey := fmt.Sprintf("user:%d", user.Id)

	// 使用 pipeline
	_, err := constant.REDIS.TxPipelined(constant.CONTEXT, func(pipe redis.Pipeliner) error {
		pipe.HSet(constant.CONTEXT, userKey, "id", user.Id)
		pipe.HSet(constant.CONTEXT, userKey, "name", user.Name)
		pipe.HSet(constant.CONTEXT, userKey, "password", user.Password)
		pipe.HSet(constant.CONTEXT, userKey, "follow_count", user.FollowerCount)
		pipe.HSet(constant.CONTEXT, userKey, "follower_count", user.FollowerCount)
		pipe.HSet(constant.CONTEXT, userKey, "create_time", user.CreateTime.UnixMilli())
		// 设置过期时间
		pipe.Expire(constant.CONTEXT, userKey, constant.EXPIRE_TIME+time.Duration(rand.Float64()*constant.EXPIRE_TIME.Seconds())*time.Second)
		return nil
	})
	return err
}

func GetUserInfoByUserIDFromRedis(userID int64) (*entity.User, error) {
	// 定义 key
	userKey := fmt.Sprintf("user:%d", userID)

	var user entity.User

	if result := constant.REDIS.Exists(constant.CONTEXT, userKey).Val(); result <= 0 {
		return nil, errors.New("not found in cache")
	}
	// 使用 pipeline
	commands, err := constant.REDIS.TxPipelined(constant.CONTEXT, func(pipe redis.Pipeliner) error {
		pipe.HGetAll(constant.CONTEXT, userKey)
		pipe.HGet(constant.CONTEXT, userKey, "create_time").Val()
		// 设置过期时间
		pipe.Expire(constant.CONTEXT, userKey, constant.EXPIRE_TIME+time.Duration(rand.Float64()*constant.EXPIRE_TIME.Seconds())*time.Second)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err = commands[0].(*redis.StringStringMapCmd).Scan(&user); err != nil {
		fmt.Println(err)
		return nil, err
	}

	timeUnixMilliStr := commands[1].(*redis.StringCmd).Val()
	timeUnixMilli, _ := strconv.ParseInt(timeUnixMilliStr, 10, 64)
	user.UpdateTime = time.UnixMilli(timeUnixMilli)
	return &user, nil
}
