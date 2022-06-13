package constants

import "time"

const (

	//other
	DefaultVideoLimit = 30

	ForWardAction = 1
	ReverseAction = 2

	TimeLayout        = "2006-01-02 15:04:05"
	CommentTimeLayout = "1-2"
	ChinaTimeZone     = "Asia/Shanghai"

	//db
	FavorVideoM2MTable = "favor_videos"
	FollowsM2MTable    = "follows"

	TitleMaxLength = 32

	//gin
	UserMaxLength = 32

	//data
	VideoDataForm  = "data"
	VideoTitleForm = "title"

	TokenQuery     = "token"
	FeedLimitQuery = "latest_time"
	UserIdQuery    = "user_id"
	VideoIdQuery   = "video_id"
	ActionQuery    = "action_type"
	NameQuery      = "username"
	PasswordQuery  = "password"
	FollowQuery    = "to_user_id"
)

var TimeZone *time.Location

func init() {
	timeZone, err := time.LoadLocation(ChinaTimeZone)
	if err != nil {
		panic(err)
	}
	TimeZone = timeZone
}

var VideoFormat = [...]string{
	".mp4",
}

var KeyId = map[string]int64{}
