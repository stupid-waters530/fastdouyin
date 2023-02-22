package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TokenParse() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenStr := context.Query("token")
		if tokenStr == "" {
			tokenStr = context.PostForm("token")
		}
		fmt.Println("----------------------------------------------------------------------------------------")
		//fmt.Println("tokenStr : " + tokenStr)
		if tokenStr == "" {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "没有权限",
			})
		}

		token, claims, err := ParseToken(tokenStr)
		if err != nil || !token.Valid {
			context.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "没有权限",
			})
			context.Abort()
			return
		}

		context.Set("my_uid", claims.UserId)

	}
}
