package service

import (
	"douyin/dao"
	model2 "douyin/model"
	"douyin/pkg/common"
	"douyin/pkg/errno"
)

func VideoList(userID int64) ([]*model2.Video, error) {
	return NewVideoListFlow(userID).Do()
}

type VideoListFlow struct {
	userID int64
	user   *model2.User
	videos []*model2.Video
}

func NewVideoListFlow(userID int64) *VideoListFlow {
	return &VideoListFlow{
		userID: userID,
	}
}
func (f *VideoListFlow) Do() ([]*model2.Video, error) {
	var err error

	if err = f.checkParams(); err != nil {
		return nil, err
	}

	if f.user, err = dao.NewUserDaoImpl().QueryUserWithVideosByUserId(f.userID); err != nil {
		return nil, err
	}

	if err = f.pack(); err != nil {
		return nil, err
	}

	if len(f.videos) == 0 {
		return nil, errno.DataEmptyErr.WithMessage("user has not videos")
	}

	return f.videos, nil
}
func (f *VideoListFlow) checkParams() error {
	return common.CheckIds(f.userID)
}

func (f *VideoListFlow) pack() error {
	videoIds := make([]int64, len(f.videos))
	f.videos = f.user.Videos
	for _, video := range f.videos {
		video.Author = f.user
		videoIds = append(videoIds, video.Id)
	}

	those, err := dao.NewUserDaoImpl().QueryIsFavorThose(f.userID, &videoIds)
	if err != nil {
		return err
	}
	for _, video := range f.videos {
		if _, ok := those[video.Id]; ok {
			video.IsFavorite = true
		}
	}
	return nil
}
