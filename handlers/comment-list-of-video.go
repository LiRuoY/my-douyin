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

type CommentListOfVideoResp struct {
	common.Response
	CommentList []*model.Comment `json:"comment_list"`
}

func CommentListOfVideo(c *gin.Context) {
	vid := c.Query(constants.VideoIdQuery)
	videoId, err := strconv.ParseInt(vid, 10, 64)
	if err != nil {
		SendCommentListOfVideoResp(c, errno.ConvertErr(err), nil)
		return
	}
	commentOfVideoList, err := service.CommentListOfVideo(videoId)
	if err != nil {
		SendCommentListOfVideoResp(c, errno.ConvertErr(err), nil)
		return
	}
	SendCommentListOfVideoResp(c, errno.Success, commentOfVideoList)

}
func SendCommentListOfVideoResp(c *gin.Context, err errno.Err, commentList []*model.Comment) {
	c.JSON(http.StatusOK, CommentListOfVideoResp{
		Response: common.Response{
			StatusCode: err.ErrCode,
			StatusMsg:  err.ErrMsg,
		},
		CommentList: commentList,
	})
}
