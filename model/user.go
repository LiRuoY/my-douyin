package model

//User
//has a UserLogin,except create and delete,should omit it,but only create and update will associate
//has many Videos
//favor_videos for has Videos table,so make other table to store
//follow_relations for has User table,so make other table to store
//has many Comments
type User struct {
	Id            int64      `json:"id,omitempty"`
	Name          string     `json:"name,omitempty"`
	FollowCount   int64      `json:"follow_count"`
	FollowerCount int64      `json:"follower_count"`
	IsFollow      bool       `json:"is_follow"`
	UserLogin     *UserLogin `json:"-"`
	Videos        []*Video   `json:"-"`
	FavorVideos   []*Video   `json:"-" gorm:"many2many:favor_videos"`
	Follows       []*User    `json:"-" gorm:"many2many:follows"`
	Comments      []*Comment `json:"-" `
}
