package service

import (
	dao2 "douyin/dao"
	model2 "douyin/model"
	"douyin/pkg/common"
	"douyin/pkg/errno"
	"douyin/util"
)

func CommentListOfVideo(videoId int64) ([]*model2.Comment, error) {
	return NewCommentListOfVideoFlow(videoId).Do()
}

type CommentListOfVideoFlow struct {
	videoId        int64
	commentList    []*model2.Comment
	commentAuthors map[int64]*model2.User
}

func NewCommentListOfVideoFlow(video int64) *CommentListOfVideoFlow {
	return &CommentListOfVideoFlow{
		videoId: video,
	}
}
func (f *CommentListOfVideoFlow) Do() ([]*model2.Comment, error) {
	if err := common.CheckIds(f.videoId); err != nil {
		return nil, err
	}
	if err := f.prepareData(); err != nil {
		return nil, err
	}

	f.pack()

	return f.commentList, nil
}

func (f *CommentListOfVideoFlow) prepareData() error {
	if err := dao2.NewVideoDAOImpl().QueryCommentListByVideoId(f.videoId, &f.commentList); err != nil {
		return err
	}
	if len(f.commentList) == 0 {
		return errno.DataEmptyErr.WithMessage("the commentList of video is null")
	}

	//去重
	delRe := make(map[int64]struct{}, len(f.commentList)/2)
	cl := len(f.commentList)
	for i := 0; i < cl; i++ {
		delRe[f.commentList[i].UserId] = struct{}{}
	}

	userIds := make([]int64, len(delRe))
	for id, _ := range delRe {
		userIds = append(userIds, id)
	}

	users, err := dao2.NewUserDaoImpl().QueryUsersByUserIds(&userIds)
	if err != nil {
		return err
	}

	f.commentAuthors = make(map[int64]*model2.User, len(users))
	for _, user := range users {
		f.commentAuthors[user.Id] = user
	}
	return nil
}

func (f *CommentListOfVideoFlow) pack() {
	for _, comment := range f.commentList {
		comment.User = f.commentAuthors[comment.UserId]
		comment.CreateDate = util.CommentDateFormat(comment.CreatedAt)
	}
}
