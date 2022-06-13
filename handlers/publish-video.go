package handlers

import (
	"douyin/config"
	"douyin/pkg/common"
	"douyin/pkg/constants"
	"douyin/pkg/errno"
	"douyin/pkg/myjwt"
	"douyin/service"
	"douyin/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"path/filepath"
)

func PublicVideo(c *gin.Context) {
	file, err := getFile(c)
	if err != nil {
		common.SendCommonResp(c, errno.ConvertErr(err))
		return
	}

	title := c.PostForm(constants.VideoTitleForm)

	UserId, err := myjwt.IdentityHandler(c)
	if err != nil {
		common.SendCommonResp(c, errno.ConvertErr(err))
		return
	}
	userFolder, err := getUserFolder(UserId)
	if err != nil {
		common.SendCommonResp(c, errno.ConvertErr(err))
		return
	}

	videoPath, err := downloadFileToLocal(c, file, userFolder)
	if err != nil {
		common.SendCommonResp(c, errno.ConvertErr(err))
		return
	}

	//coverURL暂时不要
	if err = service.PublishVideo(UserId, title, videoPath, ""); err != nil {
		common.SendCommonResp(c, errno.ConvertErr(err))
		return
	}

	common.SendCommonResp(c, errno.Success)
}

func downloadFileToLocal(c *gin.Context, file *multipart.FileHeader, userFolder string) (string, error) {
	//videoPath  static/id/filename
	videoPath := userFolder + util.GetVideoFileName(file.Filename)
	videoLocalStorePath := config.CStatic.LocalStorePath + videoPath

	if err := c.SaveUploadedFile(file, videoLocalStorePath); err != nil {
		return "", err
	}

	return videoPath, nil
}
func getUserFolder(userId int64) (userFolder string, err error) {

	if userFolder, err = util.NewFileByUserIDIfNotExist(userId); err != nil {
		return
	}
	return userFolder, nil
}

func getFile(c *gin.Context) (*multipart.FileHeader, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}
	files := form.File[constants.VideoDataForm]
	if files == nil || len(files) == 0 {
		return nil, errno.ParamErr.WithMessage("data(video) is nil")
	}

	if len(files) != 1 {
		return nil, errno.ParamErr.WithMessage("only can post one video")
	}

	file := files[0]
	if err = checkFileFormat(filepath.Ext(file.Filename)); err != nil {
		return nil, err
	}
	return file, nil
}
func checkFileFormat(ext string) error {
	for _, format := range constants.VideoFormat {
		if ext == format {
			return nil
		}
	}
	return errno.ParamErr.WithMessage(fmt.Sprintf("can not post the video type,%s", ext))
}
