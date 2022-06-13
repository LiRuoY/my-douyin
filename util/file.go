package util

import (
	"douyin/config"
	"douyin/pkg/errno"
	"errors"
	"fmt"
	"os"
	"time"
)

//if not exists the userFolder,to new one,return userFolder(id+/) and err
//otherwise,return userFolder and  nil
func NewFileByUserIDIfNotExist(userID int64) (string, error) {
	userFolder := fmt.Sprintf("%d/", userID)
	//userLocalFolderPath := constants.StaticLocalFolderPath + userFolder
	userLocalFolderPath := config.CStatic.LocalStorePath + userFolder
	if _, err := os.Stat(userLocalFolderPath); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(userLocalFolderPath, os.ModeDir)
		if err != nil {
			return userFolder, errno.ConvertErr(err)
		}
	}

	return userFolder, nil
}
func GetVideoFileName(fileName string) string {
	y, m, d := time.Now().Date()
	return fmt.Sprintf("%04d%02d%02d%s", y, m, d, fileName)
}

func GetVideoPlayURL(videoStorePath string) string {
	//return constants.StaticUrl + videoStorePath
	return config.CStatic.GetStaticUrl() + videoStorePath
}
func GetVideoCoverURL() string {
	return config.CStatic.GetStaticUrl() + "onlycover.jpg"
}
