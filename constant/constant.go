package constant

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	WRONG_ID      int64  = 0
	MISTAKE       int64  = -1
	RESP_MISTAKE  int32  = 1
	RESP_SUCCESS  int32  = 0
	USER_NOT_EXIT string = "用户不存在"

	CONTEXT               = context.Background()
	REDIS   *redis.Client // Redis 缓存接口

	EXPIRE_TIME = 10 * time.Minute
)
