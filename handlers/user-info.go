package handlers

import (
	"douyin/model"
	"douyin/pkg/common"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UserInfo(c *gin.Context) {
	id := c.Query(constants.UserIdQuery)
	userID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		SendUserInfoResp(c, errno.ConvertErr(err), model.User{})
	}
	userInfo, err := service.UserInfo(userID)
	if err != nil {
		SendUserInfoResp(c, errno.ConvertErr(err), model.User{})
	}
	SendUserInfoResp(c, errno.Success, *userInfo)
}

type UserInfoResp struct {
	common.Response
	User model.User `json:"user"`
}

func SendUserInfoResp(c *gin.Context, err errno.Err, user model.User) {
	c.JSON(http.StatusOK, UserInfoResp{
		Response: common.Response{
			StatusCode: err.ErrCode,
			StatusMsg:  err.ErrMsg,
		},
		User: user,
	})
}
