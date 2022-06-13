package handlers

import (
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

func FavorVideoList(c *gin.Context) {
	id := c.Query(constants.UserIdQuery)
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		SendVideoListResp(c, errno.ConvertErr(err), nil)
		return
	}

	favorVideoList, err := service.FavorVideoList(userId)
	if err != nil {
		SendVideoListResp(c, errno.ConvertErr(err), nil)
		return
	}
	SendVideoListResp(c, errno.Success, favorVideoList)
}
