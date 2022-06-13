package service

import (
	"douyin/dao"
	model2 "douyin/model"
	"douyin/pkg/common"
	"douyin/pkg/errno"
)

func FavorVideoList(userId int64) ([]*model2.Video, error) {
	return NewFavorVideoListFlow(userId).Do()
}

type FavorVideoListFlow struct {
	userId         int64
	favorVideoList []*model2.Video
	authors        map[int64]*model2.User
}

func NewFavorVideoListFlow(userId int64) *FavorVideoListFlow {
	return &FavorVideoListFlow{
		userId: userId,
	}
}

func (f *FavorVideoListFlow) Do() ([]*model2.Video, error) {
	if err := f.checkParams(); err != nil {
		return nil, err
	}
	if err := f.getFavorVideoList(); err != nil {
		return nil, err
	}

	if len(f.favorVideoList) == 0 {
		return nil, errno.DataEmptyErr.WithMessage("user has not favorite video")
	}

	if err := f.getAuthors(); err != nil {
		return nil, err
	}

	f.pack()

	return f.favorVideoList, nil
}

func (f *FavorVideoListFlow) checkParams() error {
	return common.CheckIds(f.userId)
}
func (f *FavorVideoListFlow) getFavorVideoList() error {
	var err error
	f.favorVideoList, err = dao.NewUserDaoImpl().QueryFavorVideoList(f.userId)
	if err != nil {
		return err
	}
	return nil
}

func (f *FavorVideoListFlow) getAuthors() error {
	var err error
	var author *model2.User
	f.authors = make(map[int64]*model2.User, len(f.favorVideoList)/2)
	for _, video := range f.favorVideoList {
		author = &model2.User{Id: video.UserID}
		err = dao.NewUserDaoImpl().QueryUserByID(author)
		if err != nil {
			return err
		}
		f.authors[author.Id] = author
	}
	return nil
}

func (f *FavorVideoListFlow) pack() {
	for _, video := range f.favorVideoList {
		video.Author = f.authors[video.UserID]
		if video.FavoriteCount != 0 {
			video.IsFavorite = true
		}
	}
}
