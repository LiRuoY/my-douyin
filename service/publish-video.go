package service

import (
	"douyin/dao"
	"douyin/model"
	"douyin/pkg/common"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"douyin/util"
	"fmt"
)

func PublishVideo(userID int64, title, playPath, coverPath string) error {
	return NewPublishVideoFlow(userID, title, playPath, coverPath).Do()
}

type PublishVideoFlow struct {
	userID                   int64
	title, playUrl, coverUrl string
}

func NewPublishVideoFlow(userID int64, title, playPath, coverPath string) *PublishVideoFlow {
	return &PublishVideoFlow{
		userID:   userID,
		title:    title,
		playUrl:  playPath,
		coverUrl: coverPath,
	}
}

func (f *PublishVideoFlow) Do() error {
	var err error

	if err = f.checkParams(); err != nil {
		return err
	}

	if err = f.prepareDate(); err != nil {
		return err
	}

	if err = f.publish(); err != nil {
		return err
	}

	return nil
}
func (f *PublishVideoFlow) checkParams() error {
	tl := len(f.title)
	if tl <= 0 || tl > constants.TitleMaxLength {
		return errno.ParamErr.WithMessage(fmt.Sprintf("title max length is %d, can't be %d", constants.TitleMaxLength, tl))
	}

	return common.CheckIds(f.userID)
}
func (f *PublishVideoFlow) prepareDate() error {
	f.playUrl = util.GetVideoPlayURL(f.playUrl)
	f.coverUrl = util.GetVideoCoverURL()
	return nil
}
func (f *PublishVideoFlow) publish() error {
	video := model.Video{
		UserID:   f.userID,
		Title:    f.title,
		PlayUrl:  f.playUrl,
		CoverUrl: f.coverUrl,
	}
	return dao.NewVideoDAOImpl().CreateVideo(&video)
}
