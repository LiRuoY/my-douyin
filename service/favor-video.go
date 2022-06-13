package service

import (
	"douyin/dao"
	"douyin/pkg/common"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"errors"
	"fmt"
	"strconv"
)

func FavorVideo(userId int64, videoId, action string) error {
	return NewFavorVideoFlow(userId, videoId, action).Do()
}

type FavorVideoFlow struct {
	vid, act        string
	userId, videoId int64
	action          int32
}

func NewFavorVideoFlow(uid int64, vid, act string) *FavorVideoFlow {
	return &FavorVideoFlow{
		userId: uid,
		vid:    vid,
		act:    act,
	}
}

func (f *FavorVideoFlow) Do() error {
	if err := f.convertParams(); err != nil {
		return err
	}

	if err := f.checkParams(); err != nil {
		return err
	}

	if err := f.doAction(); err != nil {
		return err
	}
	return nil
}
func (f *FavorVideoFlow) convertParams() error {
	videoId, err := strconv.ParseInt(f.vid, 10, 64)
	if err != nil {
		return errno.ParamErr.WithMessage(err.Error() + " videoId")
	}
	action, err := strconv.ParseInt(f.act, 10, 32)
	if err != nil {
		return errno.ParamErr.WithMessage(err.Error() + " action")
	}

	f.videoId = videoId
	f.action = int32(action)
	return nil
}

func (f *FavorVideoFlow) checkParams() error {
	if err := common.CheckIds(f.userId, f.videoId); err != nil {
		return err
	}

	switch f.action {
	case constants.ForWardAction:
	case constants.ReverseAction:
	default:
		return errno.ParamErr.WithMessage(fmt.Sprintf("has not the action,%d", f.action))
	}

	return nil
}

func (f *FavorVideoFlow) checkIfFavor() error {
	return dao.NewUserDaoImpl().QueryIsFavor(f.userId, f.videoId)
}

func (f *FavorVideoFlow) doAction() error {
	//err != nil，没找到，
	err := f.checkIfFavor()

	switch f.action {
	case constants.ForWardAction:
		if err == nil {
			//return dao.NewUserDaoImpl().CancelFavorVideo(f.userId, f.videoId)
			return errors.New("the user already has given a favor to the video")
		}
		return dao.NewUserDaoImpl().FavorVideo(f.userId, f.videoId)
	case constants.ReverseAction:
		if err != nil {
			//return dao.NewUserDaoImpl().FavorVideo(f.userId, f.videoId)
			return errors.New("the user has not  given a favor to the video")
		}
		return dao.NewUserDaoImpl().CancelFavorVideo(f.userId, f.videoId)
	}
	return nil
}
