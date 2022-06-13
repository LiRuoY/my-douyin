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

type FeedResp struct {
	common.Response
	VideoList []*model.Video `json:"video_list"`
	NextTime  int64          `json:"next_time"`
}

func FeedFlow(c *gin.Context) {
	var err error
	var latestTimeUnix int64 = 0
	latestTime := c.Query(constants.FeedLimitQuery)
	if latestTime != "" {
		latestTimeUnix, err = strconv.ParseInt(latestTime, 10, 64)
		if err != nil {
			sendFeedResp(c, errno.ConvertErr(err), nil, 0)
			return
		}
	}

	nextTime, videos, err := service.Feed(latestTimeUnix)
	if err != nil {
		sendFeedResp(c, errno.ConvertErr(err), nil, 0)
		return
	}
	sendFeedResp(c, errno.Success, videos, nextTime)
}
func sendFeedResp(c *gin.Context, err errno.Err, videoList []*model.Video, nextTime int64) {
	c.JSON(http.StatusOK, FeedResp{
		Response: common.Response{
			StatusCode: err.ErrCode,
			StatusMsg:  err.ErrMsg,
		},
		VideoList: videoList,
		NextTime:  nextTime,
	})
}
