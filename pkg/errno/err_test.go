package errno

import (
	"errors"
	"fmt"
	"testing"
)

func TestErr_AppendMsg(t *testing.T) {
	t.Log(ServiceErr.AppendMsg(errors.New("abc").Error()))
}

func TestConvertErr(t *testing.T) {
	e1 := Err{}
	e2 := Err{ErrMsg: "12312"}
	e3 := errors.New("sb")
	e4 := fmt.Errorf("%#v %#v", e2, e3)
	e5 := fmt.Errorf("%#v %#v", e2.Error(), e3.Error())
	e1 = ConvertErr(e1)
	e2 = ConvertErr(e2)
	e3 = ConvertErr(e3)
	e4 = ConvertErr(e4)
	e5 = ConvertErr(e5)
	t.Logf("%#v\n%#v\n%#v\n%#v\n%#v\n", e1, e2, e3, e4, e5)
}
