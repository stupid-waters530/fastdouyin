package entity

import (
	"github.com/ikuraoo/fastdouyin/constant"
	"github.com/ikuraoo/fastdouyin/util"
	"strconv"
	"sync"
	"time"
)

type Follow struct {
	ID         int64     `gorm:"primarykey"`
	MyId       int64     `gorm:"column:my_uid"`
	HisId      int64     `gorm:"column:his_uid"`
	IsFollow   bool      `gorm:"column:is_follow"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

type FollowDao struct {
}

var followDao *FollowDao
var followOnce sync.Once

func NewFollowDaoInstance() *FollowDao {
	followOnce.Do(
		func() {
			followDao = &FollowDao{}
		})
	return followDao
}

func (c *FollowDao) QueryByUId(myUId int64) (*[]Follow, error) {
	var followList []Follow
	err := db.Where("my_uid = ?", myUId).Find(&followList).Error
	if err != nil {
		util.Logger.Error("find followList by UId err:" + err.Error())
		return nil, err
	}
	return &followList, nil
}

func (*FollowDao) QueryIsFollow(myUId int64, hisUId int64) (bool, error) {
	var follow Follow
	err := db.Where("my_uid = ?", strconv.Itoa(int(myUId))).Where("his_uid = ?", strconv.Itoa(int(hisUId))).Find(&follow).Error
	if err != nil {
		return false, err
	}
	if follow.ID == constant.WRONG_ID {
		return false, nil
	}
	return true, nil
}

func (c *FollowDao) UpdateFollow(myUId int64, hisUId int64, isFollow bool) error {
	err := db.Model(Follow{}).Where("my_uid = ?", myUId).Where("his_uid = ?", hisUId).Update("is_follow", isFollow).Error
	if err != nil {
		util.Logger.Error("update follow err:" + err.Error())
		return err
	}
	return nil
}

func (c *FollowDao) CreateFollow(follow *Follow) error {
	if err := db.Create(follow).Error; err != nil {
		util.Logger.Error("insert relation err:" + err.Error())
		return err
	}
	return nil
}
