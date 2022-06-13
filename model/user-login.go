package model

type UserLogin struct {
	ID       int64
	UserID   int64
	Name     string
	Password string
	//LoginStatus need friend relations
}
