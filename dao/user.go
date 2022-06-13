package dao

import (
	model2 "douyin/model"
	"douyin/pkg/constants"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
)

type UserDAO interface {
	CreateUser(u *model2.User) error
	QueryUserByID(user *model2.User) error
	QueryUsersByUserIds(userIds *[]int64) ([]*model2.User, error)
	QueryCommentAuthor(u *model2.User, commentId int64) error

	QueryUserWithVideosByUserId(userId int64) (*model2.User, error)
	QueryFavorVideoList(userId int64) ([]*model2.Video, error)
	QueryFollowList(userId int64) ([]*model2.User, error)
	QueryFollowerList(userId int64) ([]*model2.User, error)

	QueryIsFavor(userId, videoId int64) error
	QueryIsFavorThose(userId int64, those *[]int64) (map[int64]struct{}, error)
	FavorVideo(userId, videoId int64) error
	CancelFavorVideo(userId, videoId int64) error

	QueryIfFollow(userId, toUserId int64) bool
	QueryIsFollowThose(userId int64, those *[]int64) (map[int64]struct{}, error)
	FollowUser(userId, toUserId int64) error
	CancelFollowUser(userId, toUserId int64) error
}
type UserDAOImpl struct {
}

var (
	userDAO  UserDAO
	userOnce sync.Once
)

func NewUserDaoImpl() UserDAO {
	userOnce.Do(func() {
		userDAO = new(UserDAOImpl)
	})
	return userDAO
}

func (dao *UserDAOImpl) CreateUser(u *model2.User) error {
	return DB.Create(u).Error
}

func (dao *UserDAOImpl) QueryUserByID(user *model2.User) error {
	err := DB.First(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *UserDAOImpl) QueryUserWithVideosByUserId(userId int64) (*model2.User, error) {
	user := model2.User{Id: userId}
	err := DB.Preload("Videos", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc")
	}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (dao *UserDAOImpl) QueryUsersByUserIds(userIds *[]int64) ([]*model2.User, error) {
	var users []*model2.User
	err := DB.Find(&users, *userIds).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (dao *UserDAOImpl) QueryIsFavor(userId, videoId int64) error {
	res := make(map[string]interface{}, 1)
	return DB.Table(constants.FavorVideoM2MTable).Take(&res, "video_id = ? and user_id = ?", videoId, userId).Error
}

func (dao *UserDAOImpl) QueryIsFavorThose(userId int64, those *[]int64) (map[int64]struct{}, error) {
	//先拿到user在those中的favorite video id
	var res []map[string]interface{}
	err := DB.Table(constants.FavorVideoM2MTable).Select("video_id").Where("user_id = ? and video_id in ?", userId, *those).Find(&res).Error
	if err != nil {
		return nil, err
	}
	//把id都放在这
	ids := make(map[int64]struct{}, len(res))
	for _, m := range res {
		for _, v := range m {
			i, ok := v.(int64)
			if !ok {
				return nil, errors.New(fmt.Sprintf("video_id %T is not int64", v))
			}
			ids[i] = struct{}{}
		}
	}
	return ids, nil
}

func (dao *UserDAOImpl) FavorVideo(userId, videoId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model2.Video{Id: videoId}).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err != nil {
			return err
		}

		if err := tx.Model(&model2.User{Id: userId}).Association("FavorVideos").Append(&model2.Video{Id: videoId}); err != nil {
			return err
		}
		return nil
	})

}

func (dao *UserDAOImpl) CancelFavorVideo(userId, videoId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model2.Video{Id: videoId}).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error
		if err != nil {
			return err
		}

		err = tx.Model(&model2.User{Id: userId}).Association("FavorVideos").Delete(&model2.Video{Id: videoId})
		if err != nil {
			return err
		}
		return nil
	})
}
func (dao *UserDAOImpl) QueryFavorVideoList(userId int64) ([]*model2.Video, error) {
	var favorVideos []*model2.Video
	if err := DB.Model(&model2.User{Id: userId}).Association("FavorVideos").Find(&favorVideos); err != nil {
		return nil, err
	}
	return favorVideos, nil
}

func (dao *UserDAOImpl) QueryCommentAuthor(u *model2.User, commentId int64) error {
	if err := DB.Debug().Preload("Comments", "id = ?", commentId).First(u).Error; err != nil {
		return err
	}
	return nil
}

func (dao *UserDAOImpl) QueryIfFollow(userId, toUserId int64) bool {
	res := make(map[string]interface{}, 1)
	err := DB.Table(constants.FollowsM2MTable).Where("user_id = ? and follow_id = ?", userId, toUserId).Take(&res).Error
	if err != nil {
		return false
	}
	return true
}

func (dao *UserDAOImpl) FollowUser(userId, toUserId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		user := model2.User{Id: userId}
		toUser := model2.User{Id: toUserId}
		if err := tx.Model(&user).Update("follow_count", gorm.Expr("follow_count + ?", 1)).Error; err != nil {
			return err
		}
		if err := tx.Model(&toUser).Update("follower_count", gorm.Expr("follower_count + ?", 1)).Error; err != nil {
			return err
		}
		if err := tx.Model(&user).Association("Follows").Append(&toUser); err != nil {
			return err
		}
		return nil
	})
}

func (dao *UserDAOImpl) CancelFollowUser(userId, toUserId int64) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		user := model2.User{Id: userId}
		toUser := model2.User{Id: toUserId}
		if err := tx.Model(&user).Update("follow_count", gorm.Expr("follow_count - ?", 1)).Error; err != nil {
			return err
		}
		if err := tx.Model(&toUser).Update("follower_count", gorm.Expr("follower_count - ?", 1)).Error; err != nil {
			return err
		}
		if err := tx.Model(&user).Association("Follows").Delete(&toUser); err != nil {
			return err
		}
		return nil
	})
}

func (dao *UserDAOImpl) QueryFollowList(userId int64) ([]*model2.User, error) {
	var userList []*model2.User
	if err := DB.Model(&model2.User{Id: userId}).Debug().Association("Follows").Find(&userList); err != nil {
		return nil, err
	}
	return userList, nil
}

func (dao *UserDAOImpl) QueryFollowerList(userId int64) ([]*model2.User, error) {
	//子表查询
	//SELECT * FROM `users` WHERE id IN (SELECT user_id FROM `follows` WHERE follow_id = 1);
	var userList []*model2.User
	if err := DB.Model(&model2.User{}).Where("id in (?)",
		DB.Table(constants.FollowsM2MTable).Select("user_id").Where("follow_id = ?", userId)).
		Find(&userList).Error; err != nil {
		return nil, err
	}
	return userList, nil
}
func (dao *UserDAOImpl) QueryIsFollowThose(userId int64, those *[]int64) (map[int64]struct{}, error) {
	//先拿到user在those中的关注的user id
	var res []map[string]interface{}
	err := DB.Table(constants.FollowsM2MTable).Select("follow_id").Where("user_id = ? and follow_id in ?", userId, *those).Find(&res).Error
	if err != nil {
		return nil, err
	}
	//把id都放在这
	ids := make(map[int64]struct{}, len(res))
	for _, m := range res {
		for _, v := range m {
			i, ok := v.(int64)
			if !ok {
				return nil, errors.New(fmt.Sprintf("user_id %T is not int64", v))
			}
			ids[i] = struct{}{}
		}
	}
	return ids, nil
}
