package service

import (
	"errors"
	"fmt"
	"github.com/ikuraoo/fastdouyin/entity"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

type FavouriteActionFlow struct {
	UId int64
	VId int64
}

type FavouriteVideo struct {
	Id     int64
	Author *entity.UserInfo
}

func FavouriteAction(userId int64, videoId int64) error {
	return NewFavouriteActionFlow(userId, videoId).Do()
}

func NewFavouriteActionFlow(userId int64, videoId int64) *FavouriteActionFlow {
	return &FavouriteActionFlow{
		UId: userId,
		VId: videoId,
	}
}

func (f *FavouriteActionFlow) Do() error {
	if err := f.checkParam(); err != nil {
		return err
	}
	err := f.action()
	if err != nil {
		return err
	}
	return nil
}

func (f *FavouriteActionFlow) checkParam() error {
	return nil
}

func (f *FavouriteActionFlow) action() error {
	isFavourite, err := entity.NewFavouriteDaoInstance().QueryByVIdAndUId(f.VId, f.UId)
	fmt.Println("isFavourite : " + strconv.FormatBool(isFavourite))
	if err != nil {
		//没有找到
		if err := entity.NewFavouriteDaoInstance().CreateFavourite(&entity.Favourite{
			UId:         f.UId,
			VId:         f.VId,
			IsFavourite: true,
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		}); err != nil {
			return err
		}

		err := entity.NewVideoDaoInstance().IncFavouriteCount(f.VId)
		if err != nil {
			return err
		}

		return nil
	}
	if isFavourite == true {
		err := entity.NewVideoDaoInstance().DecFavouriteCount(f.VId)
		if err != nil {
			return err
		}
	} else {
		err := entity.NewVideoDaoInstance().IncFavouriteCount(f.VId)
		if err != nil {
			return err
		}
	}
	err = entity.NewFavouriteDaoInstance().UpdateIsFavourite(f.VId, f.UId, !isFavourite)

	if err != nil {
		return errors.New("修改失败")
	}
	return nil
}

func FavouriteList(uid int64) ([]*entity.VideoInfo, error) {
	videopath := viper.GetString("video.videoPath")
	coverpath := viper.GetString("video.coverPath")
	favouritesOfUser, err := entity.NewFavouriteDaoInstance().QueryByUId(uid)
	if err != nil {
		return nil, err
	}
	var favouriteList []*entity.VideoInfo
	for _, fav := range *favouritesOfUser {
		if fav.IsFavourite {
			user, err := entity.NewUserDaoInstance().QueryUserById(uid)
			if err != nil {
				return nil, err
			}

			video, err := entity.NewVideoDaoInstance().QueryVideoById(fav.VId)
			if err != nil {
				return nil, err
			}

			var IsFollow bool

			follow, err := entity.NewFollowDaoInstance().QueryIsFollow(uid, video.UId)
			if err != nil {
				entity.NewFollowDaoInstance().CreateFollow(&entity.Follow{
					MyId:       uid,
					HisId:      video.Id,
					IsFollow:   false,
					CreateTime: time.Now(),
					UpdateTime: time.Now(),
				})
				IsFollow = false
			} else {
				IsFollow = follow
			}
			UserInfo := &entity.UserInfo{
				Id:            user.Id,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      IsFollow,
			}

			fmt.Println(videopath + ":" + video.PlayUrl)
			favouriteList = append(favouriteList, &entity.VideoInfo{
				video.Id,
				UserInfo,
				videopath + video.PlayUrl,
				coverpath + video.CoverUrl,
				video.CommentCount,
				video.FavouriteCount,
				true,
				video.Title,
			})
		}
	}

	return favouriteList, err
}
