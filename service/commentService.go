package service

import (
	"fmt"
	"github.com/ikuraoo/fastdouyin/entity"
	"strconv"
	"time"
)

type CommentActionFlow struct {
	UId     int64
	VId     int64
	Content string
}

func CommentAction(userId int64, videoId int64, content string) error {
	return NewCommentActionFlow(userId, videoId, content).Do()
}

func NewCommentActionFlow(userId int64, videoId int64, content string) *CommentActionFlow {
	return &CommentActionFlow{
		UId:     userId,
		VId:     videoId,
		Content: content,
	}
}

func (c *CommentActionFlow) Do() error {
	if err := c.checkParam(); err != nil {
		return err
	}
	err := c.action()
	if err != nil {
		return err
	}
	return nil
}

func (c *CommentActionFlow) checkParam() error {
	return nil
}

func (c *CommentActionFlow) action() error {
	err := entity.NewCommentDaoInstance().CreateComment(&entity.Comment{
		VId:        c.VId,
		UId:        c.UId,
		Content:    c.Content,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		IsDeleted:  false,
	})
	if err != nil {
		return err
	}
	//增加评论数
	err = entity.NewVideoDaoInstance().IncCommentCount(c.VId)
	if err != nil {
		return err
	}
	return nil
}

func CommentList(vid int64) ([]*entity.CommentInfo, error) {

	comments, err := entity.NewCommentDaoInstance().QueryByVId(vid)
	if err != nil {
		return nil, err
	}
	var commentList []*entity.CommentInfo
	for _, comment := range *comments {
		fmt.Println(comment.UId)
		fmt.Println("-------------------------------------------------")
		user, err := entity.NewUserDaoInstance().QueryUserById(comment.UId)
		fmt.Println(strconv.FormatInt(user.Id, 10))
		if err != nil {
			return nil, err
		}
		UserInfo := &entity.UserInfo{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      false,
		}
		month := comment.CreateTime.Format("01")
		date := comment.CreateTime.Format("02")

		commentList = append(
			commentList,
			&entity.CommentInfo{
				Id:         comment.Id,
				User:       UserInfo,
				Content:    comment.Content,
				CreateDate: month + "-" + date,
			},
		)
	}
	return commentList, nil
}
