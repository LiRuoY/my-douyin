package handlers

import (
	"douyin/pkg/common"
	"douyin/pkg/errno"
	"douyin/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserLogin(c *gin.Context) {
	var err error
	var login UserReq
	if err = c.ShouldBind(&login); err != nil {
		fmt.Println(err)
		SendUserResp(c, errno.ConvertErr(err), common.UserIdentity{})
		return
	}
	identity, err := service.UserLogin(&login.UserName, &login.Password)
	if err != nil {
		SendUserResp(c, errno.ConvertErr(err), common.UserIdentity{})
		return
	}
	SendUserResp(c, errno.Success, *identity)
}

func UserRegister(c *gin.Context) {
	var err error
	var register UserReq

	if err = c.ShouldBind(&register); err != nil {
		SendUserResp(c, errno.ConvertErr(err), common.UserIdentity{})
		return
	}

	identity, err := service.UserRegister(&register.UserName, &register.Password)
	if err != nil {
		SendUserResp(c, errno.ConvertErr(err), common.UserIdentity{})
		return
	}

	SendUserResp(c, errno.Success, *identity)
}

type UserReq struct {
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type UserResp struct {
	common.Response
	common.UserIdentity
}

func SendUserResp(c *gin.Context, err errno.Err, identity common.UserIdentity) {
	c.JSON(http.StatusOK, UserResp{
		Response: common.Response{
			StatusCode: err.ErrCode,
			StatusMsg:  err.ErrMsg,
		},
		UserIdentity: identity,
	})
}
