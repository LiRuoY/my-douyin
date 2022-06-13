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

func FollowUser(userId int64, toUserId, actionType string) error {
	return NewFollowUserFlow(userId, toUserId, actionType).Do()
}

type FollowUserFlow struct {
	tuid, act        string
	userId, toUserId int64
	actionType       int32
}

func NewFollowUserFlow(userId int64, toUserId, actionType string) *FollowUserFlow {
	return &FollowUserFlow{
		userId: userId,
		tuid:   toUserId,
		act:    actionType,
	}
}

func (f *FollowUserFlow) Do() error {
	if err := f.checkAndConvertParams(); err != nil {
		return err
	}

	if err := f.doAction(); err != nil {
		return err
	}
	return nil
}

func (f *FollowUserFlow) checkAndConvertParams() error {
	toUserId, err := strconv.ParseInt(f.tuid, 10, 64)
	if err != nil {
		return errno.ParamErr.WithMessage(err.Error())
	}
	actionType, err := strconv.ParseInt(f.act, 10, 32)
	if err != nil {
		return errno.ParamErr.WithMessage(err.Error())
	}

	if err = common.CheckIds(f.userId, toUserId); err != nil {
		return err
	}
	if toUserId == f.userId {
		return errors.New("can follow self")
	}

	switch actionType {
	case constants.ForWardAction:
	case constants.ReverseAction:
	default:
		return errno.ParamErr.WithMessage(fmt.Sprintf("has not the action %v", actionType))
	}
	f.toUserId = toUserId
	f.actionType = int32(actionType)
	return nil
}

func (f *FollowUserFlow) doAction() error {
	isFollow := dao.NewUserDaoImpl().QueryIfFollow(f.userId, f.toUserId)
	switch f.actionType {
	case constants.ForWardAction:
		if isFollow {
			return errors.New("the user has follow the user")
		}
		if err := dao.NewUserDaoImpl().FollowUser(f.userId, f.toUserId); err != nil {
			return err
		}
	case constants.ReverseAction:
		if !isFollow {
			return errors.New("the user has not follow the user")
		}
		if err := dao.NewUserDaoImpl().CancelFollowUser(f.userId, f.toUserId); err != nil {
			return err
		}
	}
	return nil
}
