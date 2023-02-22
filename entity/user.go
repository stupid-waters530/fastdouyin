package entity

import (
	"errors"
	"fmt"
	"github.com/ikuraoo/fastdouyin/constant"
	"github.com/ikuraoo/fastdouyin/util"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

type User struct {
	Id            int64
	Name          string
	Password      string
	FollowCount   int64
	FollowerCount int64
	CreateTime    time.Time
	UpdateTime    time.Time
	IsDeleted     bool
}

type UserDao struct {
}

type UserInfo struct {
	Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

var userDao *UserDao
var userOnce sync.Once

func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

func (*UserDao) QueryById(id int64) (*User, error) {
	var user User
	err := db.Where("id = ?", id).Find(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (*UserDao) QueryByName(name string) (*User, error) {
	var user User
	err := db.Where("name = ?", name).Find(&user).Error
	if err != nil {
		return nil, errors.New("查询出错")
	}
	if user.Id == constant.WRONG_ID {
		return nil, errors.New(constant.USER_NOT_EXIT)
	}
	return &user, nil
}

func (*UserDao) CreateUser(user *User) error {
	err := db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (*UserDao) IncUserFollow(uid int64) error {
	err := db.Model(User{}).Where("id = ?", uid).UpdateColumn("follower_count", gorm.Expr("follow_count + ?", 1)).Error
	if err != nil {
		util.Logger.Error("inc user follow count error")
		return err
	}
	return nil
}

func (*UserDao) DecUserFollow(uid int64) error {
	err := db.Model(User{}).Where("id = ?", uid).UpdateColumn("follower_count", gorm.Expr("follow_count - ?", 1)).Error
	if err != nil {
		util.Logger.Error("dec user follow count error")
		return err
	}
	return nil
}

func (*UserDao) IncUserFollower(uid int64) error {
	err := db.Model(User{}).Where("id = ?", uid).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1)).Error
	if err != nil {
		util.Logger.Error("inc user follower count error")
		return err
	}
	return nil
}

func (*UserDao) DecUserFollower(uid int64) error {
	err := db.Model(User{}).Where("id = ?", uid).UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1)).Error
	if err != nil {
		util.Logger.Error("dec user follower count error")
		return err
	}
	return nil
}

func (u *UserDao) QueryUserById(id int64) (*User, error) {
	var user User
	err := db.Model(User{}).Where("id = ?", id).Find(&user).Error
	fmt.Println("QueryUserById:" + strconv.FormatInt(id, 10))
	if err != nil {
		util.Logger.Error("user queried by id is not exist!")
		return nil, err
	}
	return &user, err
}
