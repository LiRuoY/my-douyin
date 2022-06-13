package handlers

import (
	"douyin/model"
	"douyin/pkg/common"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"douyin/pkg/myjwt"
	service2 "douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserListResp struct {
	common.Response
	UserList []*model.User `json:"user_list"`
}

func FollowUser(c *gin.Context) {
	userId, err := myjwt.IdentityHandler(c)
	if err != nil {
		common.SendCommonResp(c, errno.ConvertErr(err))
		return
	}
	if err := service2.FollowUser(userId, c.Query(constants.FollowQuery), c.Query(constants.ActionQuery)); err != nil {
		common.SendCommonResp(c, errno.ConvertErr(err))
		return
	}
	common.SendCommonResp(c, errno.Success)
}
func FollowList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query(constants.UserIdQuery), 10, 64)
	if err != nil {
		SendUserListResp(c, errno.ConvertErr(err), nil)
		return
	}
	followList, err := service2.FollowList(userId, false)
	if err != nil {
		SendUserListResp(c, errno.ConvertErr(err), nil)
		return
	}

	SendUserListResp(c, errno.Success, followList)
}
func FollowerList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query(constants.UserIdQuery), 10, 64)
	if err != nil {
		SendUserListResp(c, errno.ConvertErr(err), nil)
		return
	}

	followList, err := service2.FollowList(userId, true)
	if err != nil {
		SendUserListResp(c, errno.ConvertErr(err), nil)
		return
	}

	SendUserListResp(c, errno.Success, followList)
}
func SendUserListResp(c *gin.Context, err errno.Err, userList []*model.User) {
	c.JSON(http.StatusOK, UserListResp{
		Response: common.Response{
			StatusCode: err.ErrCode,
			StatusMsg:  err.ErrMsg,
		},
		UserList: userList,
	})
}
