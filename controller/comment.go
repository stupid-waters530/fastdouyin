package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/ikuraoo/fastdouyin/entity"
	"github.com/ikuraoo/fastdouyin/service"
	"net/http"
	"strconv"
)

type CommentListResponse struct {
	StatusCode  int32                 `json:"status_code"`
	StatusMsg   string                `json:"status_msg,omitempty"`
	CommentList []*entity.CommentInfo `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment *entity.CommentInfo `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	uid, _ := c.Get("my_uid")
	userId := uid.(int64)

	vid := c.Query("video_id")
	VideoId, err := strconv.ParseInt(vid, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "parse int error"})
	}
	actionType := c.Query("action_type")

	if actionType == "1" {
		text := c.Query("comment_text")
		err := service.CommentAction(userId, VideoId, text)
		if err != nil {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "service comment action error"})
		}
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "Comment action successfully!"})
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	vid := c.Query("video_id")
	VideoId, err := strconv.ParseInt(vid, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "parse int error"})
	}

	comments, err := service.CommentList(VideoId)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "get commentList error"})
	}

	c.JSON(http.StatusOK, CommentListResponse{
		StatusCode:  0,
		StatusMsg:   "Video loads successfully",
		CommentList: comments,
	})
}
