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

type VideoListResp struct {
	common.Response
	VideoList []*model.Video `json:"video_list"`
}

func VideoList(c *gin.Context) {
	id := c.Query(constants.UserIdQuery)
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		SendVideoListResp(c, errno.ConvertErr(err), nil)
		return
	}

	videoList, err := service.VideoList(userId)
	if err != nil {
		SendVideoListResp(c, errno.ConvertErr(err), nil)
		return
	}
	SendVideoListResp(c, errno.Success, videoList)

}

func SendVideoListResp(c *gin.Context, err errno.Err, videoList []*model.Video) {
	c.JSON(http.StatusOK, VideoListResp{
		Response: common.Response{
			StatusCode: err.ErrCode,
			StatusMsg:  err.ErrMsg,
		},
		VideoList: videoList,
	})
}
