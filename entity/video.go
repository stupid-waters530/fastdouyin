package entity

import (
	"errors"
	"fmt"
	"github.com/ikuraoo/fastdouyin/util"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Video struct {
	Id             int64     `gorm:"column:id"`
	UId            int64     `gorm:"column:uid"`
	PlayUrl        string    `gorm:"column:play_url"`
	CoverUrl       string    `gorm:"column:cover_url"`
	CommentCount   int64     `gorm:"column:comment_count"`
	FavouriteCount int64     `gorm:"column:favourite_count"`
	Title          string    `gorm:"column:title"`
	CreateTime     time.Time `gorm:"column:create_time"`
	UpdateTime     time.Time `gorm:"column:update_time"`
	IsDeleted      bool      `gorm:"column:is_deleted"`
}

type VideoInfo struct {
	Id             int64     `json:"id,omitempty"`
	Author         *UserInfo `json:"author"`
	PlayUrl        string    `json:"play_url" json:"play_url,omitempty"`
	CoverUrl       string    `json:"cover_url,omitempty"`
	FavouriteCount int64     `json:"favorite_count,omitempty"`
	CommentCount   int64     `json:"comment_count,omitempty"`
	IsFavorite     bool      `json:"is_favorite,omitempty"`
	Title          string    `json:"title,omitempty"`
}

type VideoDao struct {
}

var videoDao *VideoDao
var videoOnce sync.Once

func NewVideoDaoInstance() *VideoDao {
	videoOnce.Do(
		func() {
			videoDao = &VideoDao{}
		})
	return videoDao
}

func (*VideoDao) QueryVideos(maxNum int64) (*[]Video, error) {
	var videos []Video
	err := db.Limit(int(maxNum)).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return &videos, nil
}

func (*VideoDao) QueryVideoListByUId(uid int64) (*[]Video, error) {
	var videoList []Video
	err := db.Where("uid = ?", uid).Find(&videoList).Error
	if err != nil {
		util.Logger.Error("find videos by uid err:" + err.Error())
		return nil, err
	}
	return &videoList, nil
}

func (*VideoDao) QueryVideoById(vid int64) (*Video, error) {
	var videoList Video
	err := db.Where("id = ?", vid).Find(&videoList).Error
	if err != nil {
		util.Logger.Error("find video by vid err:" + err.Error())
		return nil, err
	}
	return &videoList, nil
}

func (*VideoDao) CreateVideo(video *Video) error {
	err := db.Create(video).Error
	fmt.Println("存视频到数据库")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (v *VideoDao) IncCommentCount(vid int64) error {
	err := db.Model(Video{}).Where("id = ?", vid).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1)).Error
	if err != nil {
		util.Logger.Error("inc video comment count error")
		return err
	}
	return nil
}

func (v *VideoDao) IncFavouriteCount(vid int64) error {
	err := db.Model(Video{}).Where("id = ?", vid).UpdateColumn("favourite_count", gorm.Expr("favourite_count + ?", 1)).Error
	if err != nil {
		util.Logger.Error("inc video favourite count error")
		return err
	}
	return nil
}

func (v *VideoDao) DecFavouriteCount(vid int64) error {
	err := db.Model(Video{}).Where("id = ?", vid).UpdateColumn("favourite_count", gorm.Expr("favourite_count - ?", 1)).Error
	if err != nil {
		util.Logger.Error("dec video favourite count error")
		return err
	}
	return nil
}

func (*VideoDao) QueryVideoCountByUserId(userId int64, count *int64) error {
	if count == nil {
		return errors.New("QueryVideoCountByUserId count 空指针")
	}
	return db.Model(&Video{}).Where("uid=?", userId).Count(count).Error
}
