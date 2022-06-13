package util

import (
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	unix, _ := ParseStrToUnix("2022-05-01 12:00:00")
	unix2, _ := ParseStrToUnix("2022-06-04 01:19:40")
	t.Log(unix)
	t.Log(unix2)
	str := ParseUnixToStr(1654171967)
	t.Log(str)
	str1 := ParseUnixToStr(1654737619)
	//
	t.Log(str1)
	t.Log(CommentDateFormat(time.Now()))

}
