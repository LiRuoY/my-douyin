package dao

import (
	model2 "douyin/model"
	"douyin/pkg/constants"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
)

type VideoDAO interface {
	CreateVideo(video *model2.Video) error
	QueryVideoListByTimeAndLimit(videoFeed *[]*model2.Video, latestTime string) error
	QueryCommentListByVideoId(videoId int64, commentList *[]*model2.Comment) error
}

type VideoDAOImpl struct {
}

var (
	videoDAO  VideoDAO
	videoOnce sync.Once
)

func NewVideoDAOImpl() VideoDAO {
	videoOnce.Do(func() {
		videoDAO = new(VideoDAOImpl)
	})
	return videoDAO
}

func (dao *VideoDAOImpl) CreateVideo(video *model2.Video) error {
	return DB.Create(video).Error
}
func (dao *VideoDAOImpl) QueryVideoListByTimeAndLimit(videoFeed *[]*model2.Video, latestTime string) error {
	if err := DB.Where("created_at >= ?", latestTime).Order("created_at desc").
		Limit(constants.DefaultVideoLimit).Find(videoFeed).Error; err != nil {
		return fmt.Errorf("err=%v, latestTime=%s", err, latestTime)
	}
	return nil
}
func (dao *VideoDAOImpl) QueryCommentListByVideoId(videoId int64, commentList *[]*model2.Comment) error {
	err := DB.Model(&model2.Video{Id: videoId}).Order("created_at desc").Association("Comments").Find(commentList)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}
	return nil
}
