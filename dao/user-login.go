package dao

import (
	"douyin/model"
	"sync"
)

type UserLoginDAO interface {
	IsUserExistByName(name string) bool
	QueryUserExistByName(name string) (*model.UserLogin, bool)
}
type UserLoginDAOImpl struct {
}

var (
	userLoginDAO  UserLoginDAO
	userLoginOnce sync.Once
)

func NewUserLoginDaoImpl() UserLoginDAO {
	userLoginOnce.Do(func() {
		userLoginDAO = new(UserLoginDAOImpl)
	})
	return userLoginDAO
}

func (dao *UserLoginDAOImpl) IsUserExistByName(name string) bool {
	var c int64
	DB.Model(&model.UserLogin{}).Where("name=?", name).Count(&c)
	if c == 0 {
		return false
	}
	return true
}

func (dao *UserLoginDAOImpl) QueryUserExistByName(name string) (*model.UserLogin, bool) {
	u := model.UserLogin{}
	err := DB.First(&u, "name=?", name).Error
	if err != nil {
		return nil, false
	}

	return &u, true
}
