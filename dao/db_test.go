package dao

import (
	"fmt"
	"testing"
)

func TestQueryUserExistByName(t *testing.T) {
	Init()
	login, ok := NewUserLoginDaoImpl().QueryUserExistByName("lym")
	if !ok {
		t.Fatal()
	}
	t.Logf("%#v\n", login)

}
func TestUserDAO_QueryFollowList(t *testing.T) {
	Init()
	if us, err := NewUserDaoImpl().QueryFollowList(1); err != nil {
		t.Fatal(err)
	} else {
		for _, u := range us {
			fmt.Println(u)
		}
		fmt.Println(len(us))
	}

}
func TestUserDAO_QueryFollowerList(t *testing.T) {
	Init()
	if us, err := NewUserDaoImpl().QueryFollowerList(1); err != nil {
		t.Fatal(err)
	} else {
		for _, u := range us {
			fmt.Println(u)
		}
		fmt.Println(len(us))
	}

}

func TestQueryIsFollowThose(t *testing.T) {
	Init()
	if us, err := NewUserDaoImpl().QueryIsFollowThose(1, &[]int64{2, 3, 4}); err != nil {
		t.Fatal(err)
	} else {
		for id, _ := range us {
			fmt.Println(id)
		}
		fmt.Println("userLen=", len(us))
	}

}
