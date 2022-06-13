package handlers

import (
	"douyin/pkg/common"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"douyin/pkg/myjwt"
	"douyin/service"
	"github.com/gin-gonic/gin"
)

func FavorVideo(c *gin.Context) {
	userId, err := myjwt.IdentityHandler(c)
	if err != nil {
		common.SendCommonResp(c, errno.ConvertErr(err))
		return
	}
	if err := service.FavorVideo(userId, c.Query(constants.VideoIdQuery), c.Query(constants.ActionQuery)); err != nil {
		common.SendCommonResp(c, errno.ConvertErr(err))
		return
	}
	common.SendCommonResp(c, errno.Success)
}
