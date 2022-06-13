package service

import (
	dao2 "douyin/dao"
	model2 "douyin/model"
	"douyin/pkg/common"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"douyin/util"
)

func CommentVideo(userId, videoId, commentId *int64, actionType *int32, commentText *string) (*model2.Comment, error) {
	return NewCommentVideoFlow(userId, videoId, commentId, actionType, commentText).Do()
}

type CommentVideoFlow struct {
	actionType    int32
	commentAuthor *model2.User
	comment       *model2.Comment
}

func NewCommentVideoFlow(userId, videoId, commentId *int64, actionType *int32, commentText *string) *CommentVideoFlow {
	return &CommentVideoFlow{
		comment: &model2.Comment{
			Id:      *commentId,
			UserId:  *userId,
			VideoId: *videoId,
			Content: *commentText,
		},
		actionType: *actionType,
	}
}

func (f *CommentVideoFlow) Do() (*model2.Comment, error) {
	if err := f.checkParams(); err != nil {
		return nil, err
	}

	if err := f.commentAndFindAuthor(); err != nil {
		return nil, err
	}

	f.pack()

	return f.comment, nil
}

func (f *CommentVideoFlow) checkParams() error {
	if err := common.CheckIds(f.comment.VideoId, f.comment.UserId); err != nil {
		return err
	}

	switch f.actionType {
	case constants.ForWardAction:
		if f.comment.Content == "" {
			return errno.ParamErr.WithMessage("can't publish empty commentText")
		}
	case constants.ReverseAction:
		//删除评论不应该有content
		f.comment.Content = ""
		if f.comment.Id <= 0 {
			return errno.ParamErr.WithMessage("commentId must >0 ")
		}
	default:
		return errno.ParamErr.WithMessage("has not the action")
	}
	return nil
}

func (f *CommentVideoFlow) commentAndFindAuthor() error {
	var err error
	switch f.actionType {
	case constants.ForWardAction:
		if err = dao2.NewCommentDAOImpl().CreateComment(f.comment); err != nil {
			return err
		}
	case constants.ReverseAction:
		//因为要返回要删除的评论，所以传comment
		if err = dao2.NewCommentDAOImpl().QueryCommentByUser(f.comment); err != nil {
			return err
		}
		if err = dao2.NewCommentDAOImpl().DeleteComment(f.comment); err != nil {
			return err
		}
	}

	f.commentAuthor = &model2.User{Id: f.comment.UserId}
	if err = dao2.NewUserDaoImpl().QueryUserByID(f.commentAuthor); err != nil {
		return err
	}
	return nil
}

func (f *CommentVideoFlow) pack() {
	f.comment.User = f.commentAuthor
	f.comment.CreateDate = util.CommentDateFormat(f.comment.CreatedAt)
}
