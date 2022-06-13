package service

import (
	dao2 "douyin/dao"
	model2 "douyin/model"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"douyin/util"
	"errors"
)

func Feed(latestTime int64) (int64, []*model2.Video, error) {
	return NewFeedFlow(latestTime).Do()
}

type FeedFlow struct {
	nextTime   int64
	latestTime int64
	videos     []*model2.Video
	Authors    map[int64]*model2.User
}

func NewFeedFlow(latestTime int64) *FeedFlow {
	return &FeedFlow{
		latestTime: latestTime,
	}
}
func (f *FeedFlow) Do() (int64, []*model2.Video, error) {
	var err error

	f.checkParams()

	if err = f.prepareData(); err != nil {
		return 0, nil, err
	}

	f.pack()

	return f.nextTime, f.videos, nil
}
func (f *FeedFlow) checkParams() {
	if f.latestTime == 0 {
		//f.latestTime = time.Now().Unix()
		f.latestTime, _ = util.ParseStrToUnix("2022-06-02")
	}
}

func (f *FeedFlow) prepareData() error {
	latestTimeString := util.ParseUnixToStr(f.latestTime)
	if latestTimeString == "" {
		return errors.New("ParseUnixToStr err")
	}

	f.videos = make([]*model2.Video, 0, constants.DefaultVideoLimit)
	if err := dao2.NewVideoDAOImpl().QueryVideoListByTimeAndLimit(&f.videos, latestTimeString); err != nil {
		return err
	}

	if len(f.videos) == 0 {
		return errno.DataEmptyErr.WithMessage("there are not videos")
	}

	//去重
	delRe := make(map[int64]struct{}, len(f.videos)/2)
	vl := len(f.videos)
	for i := 0; i < vl; i++ {
		delRe[f.videos[i].UserID] = struct{}{}
	}

	userIds := make([]int64, len(delRe))
	for k, _ := range delRe {
		userIds = append(userIds, k)
	}

	authors, err := dao2.NewUserDaoImpl().QueryUsersByUserIds(&userIds)
	if err != nil {
		return err
	}

	//因为下面要根据video的userId来取这个author
	f.Authors = make(map[int64]*model2.User, len(authors))
	for _, author := range authors {
		f.Authors[author.Id] = author
	}

	return nil
}
func (f *FeedFlow) pack() {
	for _, video := range f.videos {
		video.Author = f.Authors[video.UserID]
	}
	//在gorm配置dns时用了解析时间，且为中国的时间（local是系统时间），所以gorm自动从数据库读到的时间字符串转为时间没问题
	f.nextTime = f.videos[len(f.videos)-1].CreatedAt.Unix()
}
