package util

import (
	"douyin/pkg/constants"
	"time"
)

func CommentDateFormat(date time.Time) string {
	return date.Format(constants.CommentTimeLayout)
}

func ParseUnixToStr(latestTime int64) string {
	t := time.Unix(latestTime, 0)
	latestTimeString := t.Format(constants.TimeLayout)
	return latestTimeString
}

//ParseStrToUnix
//这个函数有点多余，但当个教训吧，ParseInLocation和Parse是不一样的
//还有方便测试时间，拿到时间戳，对比
func ParseStrToUnix(latestTime string) (int64, error) {
	parse, err := time.ParseInLocation(constants.TimeLayout, latestTime, constants.TimeZone)
	if err != nil {
		return 0, err
	}
	return parse.Unix(), nil
}
