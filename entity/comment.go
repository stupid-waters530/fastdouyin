package entity

import (
	"github.com/ikuraoo/fastdouyin/util"
	"sync"
	"time"
)

type Comment struct {
	Id         int64 `gorm:"column:id"`
	VId        int64 `gorm:"column:vid"`
	UId        int64 `gorm:"column:uid"`
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
	IsDeleted  bool
}

type CommentInfo struct {
	Id         int64     `json:"id,omitempty"`
	User       *UserInfo `json:"user"`
	Content    string    `json:"content,omitempty"`
	CreateDate string    `json:"create_date,omitempty"`
}

type CommentDao struct {
}

var commentDao *CommentDao //DAO(DataAccessObject)模式
var commentOnce sync.Once

func NewCommentDaoInstance() *CommentDao {
	commentOnce.Do(
		func() {
			commentDao = &CommentDao{}
		})
	return commentDao
}

func (c *CommentDao) QueryByVId(vid int64) (*[]Comment, error) {
	var comments []Comment
	err := db.Where("vid = ?", vid).Find(&comments).Error
	if err != nil {
		util.Logger.Error("find comment by vid err:" + err.Error())
		return nil, err
	}
	return &comments, nil
}

func (c *CommentDao) CreateComment(content *Comment) error {
	if err := db.Create(content).Error; err != nil {
		util.Logger.Error("insert favourite err:" + err.Error())
		return err
	}
	return nil
}
