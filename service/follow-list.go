package service

import (
	"douyin/dao"
	"douyin/model"
	"douyin/pkg/common"
	"douyin/pkg/errno"
	"fmt"
)

//FollowList
//isFollower表示是要followerList还是followList
func FollowList(userId int64, isFollower bool) ([]*model.User, error) {
	return NewFollowListFlow(userId, isFollower).Do()
}

type FollowListFlow struct {
	userId     int64
	isFollower bool
	userList   []*model.User
}

func NewFollowListFlow(userId int64, isFollower bool) *FollowListFlow {
	return &FollowListFlow{
		userId:     userId,
		isFollower: isFollower,
	}
}

func (f *FollowListFlow) Do() ([]*model.User, error) {
	if err := f.checkParams(); err != nil {
		return nil, err
	}
	if err := f.prepareData(); err != nil {
		return nil, err
	}

	if err := f.pack(); err != nil {
		return nil, err
	}

	return f.userList, nil
}

func (f *FollowListFlow) checkParams() error {
	return common.CheckIds(f.userId)
}

func (f *FollowListFlow) prepareData() error {
	var err error
	switch f.isFollower {
	case true:
		if f.userList, err = dao.NewUserDaoImpl().QueryFollowerList(f.userId); err != nil {
			return err
		}
		if len(f.userList) == 0 {
			return errno.DataEmptyErr.WithMessage(fmt.Sprintf("can't find the follower for user %d", f.userId))
		}
	case false:
		if f.userList, err = dao.NewUserDaoImpl().QueryFollowList(f.userId); err != nil {
			return err
		}
		if len(f.userList) == 0 {
			return errno.DataEmptyErr.WithMessage(fmt.Sprintf("can't find the follow for user %d", f.userId))
		}
	}
	return nil
}
func (f *FollowListFlow) followListPack() error {
	for _, followedUser := range f.userList {
		followedUser.IsFollow = true
	}
	return nil
}

func (f *FollowListFlow) followerListPack() error {
	//查询user是否也关注了他的follower
	ids := make([]int64, len(f.userList))
	for _, follower := range f.userList {
		ids = append(ids, follower.Id)
	}
	those, err := dao.NewUserDaoImpl().QueryIsFollowThose(f.userId, &ids)
	if err != nil {
		return err
	}

	//如果没找到user有没有对其follower关注，就是没有，直接返回nil就行，也不用在user的显示中把IsFollow改为true
	if len(those) == 0 {
		return nil
	}

	for _, user := range f.userList {
		if _, ok := those[user.Id]; ok {
			user.IsFollow = true
		}
	}
	return nil
}

func (f *FollowListFlow) pack() error {
	switch f.isFollower {
	case true:
		if err := f.followerListPack(); err != nil {
			return err
		}

	case false:
		if err := f.followListPack(); err != nil {
			return err
		}
	}
	return nil
}
