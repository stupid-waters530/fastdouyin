package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/constant"
	"github.com/ikuraoo/fastdouyin/entity"
	"github.com/ikuraoo/fastdouyin/middleware"
	"github.com/ikuraoo/fastdouyin/service"
	"github.com/ikuraoo/fastdouyin/util"
	"net/http"
	"path/filepath"
	"strconv"
)

type VideoListResponse struct {
	Response
	VideoList []*entity.VideoInfo `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {

	userId, _ := c.Get("my_uid")
	uid, _ := userId.(int64)
	title := c.PostForm("title")
	data, err := c.FormFile("data")

	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: int32(constant.MISTAKE), StatusMsg: err.Error()})
	}

	filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", userId, filename)
	saveFile := filepath.Join("./public/videos/", finalName)
	fmt.Println("获得数据uid:" + strconv.FormatInt(uid, 10))
	fmt.Println("获得数据filename:" + filename)

	if err = c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: int32(constant.MISTAKE), StatusMsg: err.Error()})
	}

	coverName := service.NewFileCoverName(uid)
	saveCoverFile := filepath.Join("./public/covers/", coverName)
	err = util.SnapShotFromVideo(saveFile, saveCoverFile, 1)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: int32(constant.MISTAKE), StatusMsg: err.Error()})
	}

	if err = service.VideoPublish(uid, title, finalName, coverName); err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: int32(constant.MISTAKE), StatusMsg: err.Error()})
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	var uid int64
	tokenStr := c.Query("token")
	if tokenStr == "" {
		tokenStr = c.PostForm("token")
	}
	token, claims, err := middleware.ParseToken(tokenStr)
	if err != nil || !token.Valid {
		uid = 0
	}
	uid = claims.UserId
	fmt.Println(uid)
	var hisUId int64
	hisuid := c.Query("user_id")
	hisUId, _ = strconv.ParseInt(hisuid, 10, 64)
	//fmt.Println("PublishList获取数据uid：" + strconv.FormatInt(uid, 10))
	videoFeed, err := service.PublishList(hisUId)
	fmt.Printf("得到%d个视频\n", len(videoFeed))
	//fmt.Print(videoFeed[0].PlayUrl + "\n")
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: int32(constant.MISTAKE), StatusMsg: err.Error()})
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoFeed,
	})
}
