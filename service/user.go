package service

import (
	dao2 "douyin/dao"
	model2 "douyin/model"
	"douyin/pkg/common"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"douyin/pkg/myjwt"
	"douyin/util"
	"fmt"
)

const (
	login    = 1
	register = 2
)

func UserLogin(userName, password *string) (*common.UserIdentity, error) {
	return NewUserFlow(userName, password, login).Do()
}
func UserRegister(userName, password *string) (*common.UserIdentity, error) {
	return NewUserFlow(userName, password, register).Do()
}
func UserInfo(userID int64) (*model2.User, error) {
	if userID <= 0 {
		return nil, errno.ParamErr.WithMessage("userId<=0")
	}

	userInfo := &model2.User{Id: userID}
	err := dao2.NewUserDaoImpl().QueryUserByID(userInfo)
	if err != nil {
		return nil, errno.ConvertErr(err).AppendMsg(fmt.Sprintf("can't find the user %d", userID))
	}

	return userInfo, nil
}

type UserFlow struct {
	identity *common.UserIdentity
	opt      int

	name     string
	password string
	id       int64
}

func NewUserFlow(name, password *string, opt int) *UserFlow {
	return &UserFlow{
		name:     *name,
		password: *password,
		opt:      opt,
	}
}

func (f *UserFlow) Do() (*common.UserIdentity, error) {
	var err error
	if err = f.checkParams(); err != nil {
		return nil, err
	}

	if err = f.prepareDate(); err != nil {
		return nil, err
	}

	if err = f.pack(); err != nil {
		return nil, err
	}

	return f.identity, nil
}

func (f *UserFlow) checkParams() error {
	nl, pl := len(f.name), len(f.password)
	if nl <= 0 || nl > constants.UserMaxLength {
		return errno.NameErr.AppendMsg(fmt.Sprintf("not %d", nl))
	}
	if pl <= 0 || pl > constants.UserMaxLength {
		return errno.PasswordErr.AppendMsg(fmt.Sprintf("max length is %d,not %d", constants.UserMaxLength, pl))
	}
	return nil
}

func (f *UserFlow) prepareDate() error {
	switch f.opt {
	case login:
		return f.userLogin()
	case register:
		return f.userRegister()
	}
	return fmt.Errorf("user service has not the option %d", f.opt)
}

func (f *UserFlow) userLogin() error {
	login, exist := dao2.NewUserLoginDaoImpl().QueryUserExistByName(f.name)
	if !exist {
		return errno.UserNotExistErr
	}
	if !util.Equal(f.password, login.Password) {
		return errno.PasswordErr
	}
	f.id = login.UserID
	return nil
}

func (f *UserFlow) userRegister() error {
	if exist := dao2.NewUserLoginDaoImpl().IsUserExistByName(f.name); exist {
		return errno.UserAlreadyExistErr
	}

	user := model2.User{
		Name: f.name,
		UserLogin: &model2.UserLogin{
			Name:     f.name,
			Password: util.Encrypted(f.password),
		},
	}
	if err := dao2.NewUserDaoImpl().CreateUser(&user); err != nil {
		return errno.ConvertErr(err)
	}

	f.id = user.Id
	return nil
}

func (f *UserFlow) pack() error {
	token, err := myjwt.ReleaseToken(f.id)
	if err != nil {
		return errno.ConvertErr(err)
	}

	f.identity = &common.UserIdentity{
		UserId: f.id,
		Token:  token,
	}
	return nil
}
