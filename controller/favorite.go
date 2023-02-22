package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/entity"
	"github.com/ikuraoo/fastdouyin/service"
	"net/http"
	"strconv"
)

type PublishListResponse struct {
	StatusCode    int32  `json:"status_code"`
	StatusMsg     string `json:"status_msg,omitempty"`
	FavouriteList []*entity.VideoInfo
}

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	uid, _ := c.Get("my_uid")
	vid := c.Query("video_id")
	//fmt.Println(uid)
	//fmt.Println(vid)
	//userId, err := strconv.ParseInt(uid, 10, 64)
	userId := uid.(int64)
	videoId, err := strconv.ParseInt(vid, 10, 64)

	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "id convert error"},
		})
		return
	}
	err = service.FavouriteAction(userId, videoId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "favourite change error"})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	uid := c.Query("my_uid")
	userId, err := strconv.ParseInt(uid, 10, 64)

	favouriteList, err := service.FavouriteList(userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Video loads Failed",
		})
	}
	c.JSON(http.StatusOK, PublishListResponse{
		StatusCode:    0,
		StatusMsg:     "publishList successfully",
		FavouriteList: favouriteList,
	})
}
