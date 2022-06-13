package dao

import (
	model2 "douyin/model"
	"gorm.io/gorm"
	"sync"
)

type CommentDAO interface {
	CreateComment(com *model2.Comment) error
	DeleteComment(com *model2.Comment) error
	QueryCommentByUser(com *model2.Comment) error
}

type CommentDAOImpl struct {
}

var (
	commentDAO  CommentDAO
	commentOnce sync.Once
)

func NewCommentDAOImpl() CommentDAO {
	commentOnce.Do(func() {
		commentDAO = new(CommentDAOImpl)
	})
	return commentDAO
}
func (dao *CommentDAOImpl) CreateComment(com *model2.Comment) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&com).Error; err != nil {
			return err
		}
		if err := tx.Model(&model2.Video{Id: com.VideoId}).Update("comment_count", gorm.Expr("comment_count + ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})
}
func (dao *CommentDAOImpl) DeleteComment(com *model2.Comment) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model2.Video{Id: com.VideoId}).Update("comment_count", gorm.Expr("comment_count - ?", 1)).Error; err != nil {
			return err
		}

		if err := tx.Delete(&com).Error; err != nil {
			return err
		}
		return nil
	})
}

func (dao *CommentDAOImpl) QueryCommentByUser(com *model2.Comment) error {
	return DB.First(com, "user_id = ?", com.UserId).Error
}
