package model

import "time"

//Video
//CreatedAt 和 UpdatedAt 要放入数据库
type Video struct {
	Id            int64      `json:"id,omitempty"`
	UserID        int64      `json:"-"`
	Author        *User      `json:"author" gorm:"-"` //User has many Video,but need json
	Title         string     `json:"title,omitempty"`
	PlayUrl       string     `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string     `json:"cover_url,omitempty"`
	FavoriteCount int64      `json:"favorite_count"`
	IsFavorite    bool       `json:"is_favorite" gorm:"-"`
	CommentCount  int64      `json:"comment_count"`
	Users         []*User    `json:"-" gorm:"many2many:favor_videos"` //in the same table with User's FavorVideo
	Comments      []*Comment `json:"-"`                               //has many Comments
	CreatedAt     time.Time  `json:"-"`
	UpdatedAt     time.Time  `json:"-"`
}
