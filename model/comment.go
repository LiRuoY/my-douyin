package model

import "time"

//Comment
//User User has many comment,so can't serialize the user into the db,only to be json
type Comment struct {
	Id         int64     `json:"id,omitempty"`
	UserId     int64     `json:"-"`
	VideoId    int64     `json:"-"`
	Content    string    `json:"content,omitempty"`
	User       *User     `json:"user,omitempty" gorm:"-"`
	CreatedAt  time.Time `json:"-"`
	CreateDate string    `json:"create_date,omitempty" gorm:"-"`
	UpdatedAt  time.Time `json:"-"`
}
