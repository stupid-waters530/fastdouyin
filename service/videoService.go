package service

import (
	"fmt"
	"github.com/ikuraoo/fastdouyin/constant"
	"github.com/ikuraoo/fastdouyin/entity"
	"github.com/spf13/viper"
	"log"
	"strconv"
	"time"
)

func VideoFeed(myUId int64, LatestTime int64) ([]*entity.VideoInfo, error) {
	MaxNum := viper.GetString("video.maxNum")
	maxNum, err := strconv.ParseInt(MaxNum, 10, 64)
	if err != nil {
		return nil, err
	}

	videos, err := GetVideos(maxNum)
	videosWithUsers, err := CombinationVideosAndUsers(myUId, videos)

	//fmt.Println("--------------------------------------------")
	//fmt.Println(videosWithUsers[0].PlayUrl)
	return videosWithUsers, err
}

func GetVideos(maxNum int64) (*[]entity.Video, error) {
	videos, err := entity.NewVideoDaoInstance().QueryVideos(maxNum)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func CombinationVideosAndUsers(myUId int64, videos *[]entity.Video) ([]*entity.VideoInfo, error) {
	videoPath := viper.GetString("video.videoPath")
	coverPath := viper.GetString("video.coverPath")
	var video entity.Video
	var author *entity.User
	var err error
	var videosWithUsers = make([]*entity.VideoInfo, 0)
	var isFollow bool

	for _, video = range *videos {
		author, err = entity.NewUserDaoInstance().QueryById(video.UId)
		if err != nil {
			continue
		}
		if myUId != constant.WRONG_ID {
			isFollow, err = entity.NewFollowDaoInstance().QueryIsFollow(myUId, author.Id)
			if err != nil {
				fmt.Println(err)
			}
		}
		isFavourite, _ := entity.NewFavouriteDaoInstance().QueryByVIdAndUId(video.Id, myUId)
		fmt.Println("feed中isFavourite：" + strconv.FormatBool(isFavourite))
		VideoInfo := entity.VideoInfo{
			Id: video.Id,
			Author: &entity.UserInfo{
				Id:            author.Id,
				Name:          author.Name,
				FollowCount:   author.FollowCount,
				FollowerCount: author.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:        videoPath + video.PlayUrl,
			CoverUrl:       coverPath + video.CoverUrl,
			CommentCount:   video.CommentCount,
			FavouriteCount: video.FavouriteCount,
			IsFavorite:     isFavourite,
			Title:          video.Title,
		}
		videosWithUsers = append(videosWithUsers, &VideoInfo)
	}
	return videosWithUsers, nil
}

func VideoPublish(uid int64, title, filename, covername string) error {

	video := &entity.Video{
		//Id:             0,
		UId:            uid,
		PlayUrl:        filename,
		CoverUrl:       covername,
		CommentCount:   0,
		FavouriteCount: 0,
		Title:          title,
		CreateTime:     time.Now(),
		UpdateTime:     time.Now(),
		IsDeleted:      false,
	}
	err := entity.NewVideoDaoInstance().CreateVideo(video)

	if err != nil {
		return err
	}
	return nil
}

func PublishList(uid int64) ([]*entity.VideoInfo, error) {
	videos, err := entity.NewVideoDaoInstance().QueryVideoListByUId(uid)
	if err != nil {
		return nil, err
	}

	videosWithUsers, err := CombinationVideosAndUsers(uid, videos)
	if err != nil {
		return nil, err
	}
	return videosWithUsers, nil
}

func NewFileCoverName(userId int64) string {
	var count int64
	count = 0
	err := entity.NewVideoDaoInstance().QueryVideoCountByUserId(userId, &count)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%d-%d.jpg", userId, count)
}
