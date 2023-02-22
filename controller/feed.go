package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/constant"
	"github.com/ikuraoo/fastdouyin/entity"
	"github.com/ikuraoo/fastdouyin/middleware"
	"github.com/ikuraoo/fastdouyin/service"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []*entity.VideoInfo `json:"video_list,omitempty"`
	NextTime  int64               `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	tokenStr := c.Query("token")
	if tokenStr == "" {
		tokenStr = c.PostForm("token")
	}

	_, claims, err := middleware.ParseToken(tokenStr)
	c.Set("my_uid", claims.UserId)
	uid, _ := c.Get("my_uid")
	var myUId int64
	if uid != nil {
		myUId = uid.(int64)

	} else {
		myUId = 0
	}
	fmt.Println("feed中uid：" + strconv.FormatInt(myUId, 10))
	videoFeed, err := service.VideoFeed(myUId, 1)
	if err != nil {
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{
				StatusCode: constant.RESP_MISTAKE,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "success",
		},
		VideoList: videoFeed,
		NextTime:  time.Now().Unix(),
	})
}
