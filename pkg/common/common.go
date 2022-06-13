package common

import (
	"douyin/pkg/errno"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func SendCommonResp(c *gin.Context, err errno.Err) {
	c.JSON(http.StatusOK, Response{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	})
}

type UserIdentity struct {
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token,omitempty"`
}

func CheckIds(ids ...int64) error {
	for _, id := range ids {
		if id <= 0 {
			return errno.ParamErr.WithMessage(fmt.Sprintf("Id<=0 ,the id =%d  all=%v", id, ids))
		}
	}
	return nil
}
