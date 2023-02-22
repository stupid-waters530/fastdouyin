package service

import (
	"fmt"
	"github.com/ikuraoo/fastdouyin/entity"
	"strconv"
	"time"
)

type RelationActionFlow struct {
	MyUId      int64
	HisUId     int64
	ActionType string
}

func RelationAction(myUId int64, hisUId int64, actionType string) error {
	return NewRelationActionFlow(myUId, hisUId, actionType).Do()
}

func NewRelationActionFlow(myUId int64, hisUId int64, actionType string) *RelationActionFlow {
	return &RelationActionFlow{
		MyUId:      myUId,
		HisUId:     hisUId,
		ActionType: actionType,
	}
}

func (c *RelationActionFlow) Do() error {
	if err := c.checkParam(); err != nil {
		return err
	}
	err := c.action()
	if err != nil {
		return err
	}
	return nil
}

func (c *RelationActionFlow) checkParam() error {
	return nil
}

func (c *RelationActionFlow) action() error {
	var isFollow bool
	if c.ActionType == "1" {
		isFollow = true
	} else {
		isFollow = false
	}
	follow, err := entity.NewFollowDaoInstance().QueryIsFollow(c.MyUId, c.HisUId)
	fmt.Println("follow action: " + strconv.FormatBool(follow))
	if follow == true && isFollow == false {
		err := entity.NewUserDaoInstance().DecUserFollow(c.MyUId)
		if err != nil {
			return err
		}
		err = entity.NewUserDaoInstance().DecUserFollower(c.HisUId)
		if err != nil {
			return err
		}
	} else if follow == false && isFollow == true {
		err := entity.NewUserDaoInstance().IncUserFollow(c.MyUId)
		if err != nil {
			return err
		}
		err = entity.NewUserDaoInstance().IncUserFollower(c.HisUId)
		if err != nil {
			return err
		}
		err = entity.NewFollowDaoInstance().CreateFollow(&entity.Follow{
			MyId:       c.MyUId,
			HisId:      c.HisUId,
			IsFollow:   isFollow,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		})
		if err != nil {
			return err
		}
	}

	if err != nil {
		fmt.Println("没有查询到")
		err := entity.NewFollowDaoInstance().CreateFollow(&entity.Follow{
			MyId:       c.MyUId,
			HisId:      c.HisUId,
			IsFollow:   isFollow,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		})
		if err != nil {
			return err
		}
		return nil
	}
	err = entity.NewFollowDaoInstance().UpdateFollow(c.MyUId, c.HisUId, isFollow)
	if err != nil {
		return err
	}
	return nil
}

func RelationFollowList(uid int64) ([]*entity.UserInfo, error) {
	var userList []*entity.UserInfo
	followList, err := entity.NewFollowDaoInstance().QueryByUId(uid)
	if err != nil {
		return nil, err
	}
	for _, follow := range *followList {
		user, err := entity.NewUserDaoInstance().QueryUserById(follow.HisId)
		if err != nil {
			return nil, err
		}
		isFollow, _ := entity.NewFollowDaoInstance().QueryIsFollow(uid, follow.HisId)
		userList = append(userList, &entity.UserInfo{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      isFollow,
		})
	}
	return userList, nil
}

func RelationFollowerList(uid int64) ([]*entity.UserInfo, error) {
	var protoUserList []*entity.UserInfo
	followerList, err := entity.NewFollowDaoInstance().QueryByUId(uid)
	if err != nil {
		return nil, err
	}
	for _, follow := range *followerList {
		user, err := entity.NewUserDaoInstance().QueryUserById(follow.MyId)
		if err != nil {
			return nil, err
		}
		protoUserList = append(protoUserList, &entity.UserInfo{
			Id:            user.Id,
			Name:          user.Name,
			FollowCount:   user.FollowCount,
			FollowerCount: user.FollowerCount,
			IsFollow:      true,
		})
	}
	return protoUserList, nil
}
