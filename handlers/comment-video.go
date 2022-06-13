package handlers

import (
	"douyin/model"
	"douyin/pkg/common"
	"douyin/pkg/errno"
	"douyin/pkg/myjwt"
	"douyin/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CommentReq struct {
	UserId      int64  `form:"user_id,omitempty"`
	VideoId     int64  `form:"video_id,omitempty"`
	ActionType  int32  `form:"action_type,omitempty"`
	CommentText string `form:"comment_text,omitempty"`
	CommentId   int64  `form:"comment_id,omitempty"`
}
type CommentResp struct {
	common.Response
	Comment *model.Comment `json:"comment,omitempty"`
}

func CommentVideo(c *gin.Context) {
	com := CommentReq{}
	//get方法才能用Query
	if err := c.ShouldBindQuery(&com); err != nil {
		SendCommentResp(c, errno.ConvertErr(err), nil)
		return
	}
	userId, err := myjwt.IdentityHandler(c)
	if err != nil {
		SendCommentResp(c, errno.ConvertErr(err), nil)
		return
	}
	comment, err := service.CommentVideo(&userId, &com.VideoId, &com.CommentId, &com.ActionType, &com.CommentText)
	if err != nil {
		SendCommentResp(c, errno.ConvertErr(err), nil)
		return
	}

	SendCommentResp(c, errno.Success, comment)
}

func SendCommentResp(c *gin.Context, err errno.Err, comment *model.Comment) {
	c.JSON(http.StatusOK, CommentResp{
		Response: common.Response{
			StatusCode: err.ErrCode,
			StatusMsg:  err.ErrMsg,
		},
		Comment: comment,
	})
}
