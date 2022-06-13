package errno

import (
	"douyin/pkg/constants"
	"errors"
	"fmt"
)

const (
	SuccessCode             = 0
	ServiceErrCode          = 10001
	ParamErrCode            = 10002
	UserAlreadyExistErrCode = 10003
	UserNotExistErrCode     = 10004
	IdentityKeyClassErrCode = 10005
	AssertTypeErrCode       = 10006
	DataEmptyErrCode        = 10007
)

type Err struct {
	ErrCode int32
	ErrMsg  string
}

func (e Err) Error() string {
	return fmt.Sprintf("err_code=%d,err_message=%s", e.ErrCode, e.ErrMsg)
}

func NewErr(code int, msg string) Err {
	return Err{
		ErrCode: int32(code),
		ErrMsg:  msg,
	}
}
func (e Err) AppendMsg(msg string) Err {
	e.ErrMsg += "," + msg
	return e
}
func (e Err) WithMessage(msg string) Err {
	e.ErrMsg = msg
	return e
}

var (
	Success             = NewErr(SuccessCode, "Success")
	ServiceErr          = NewErr(ServiceErrCode, "Service can't run successful")
	ParamErr            = NewErr(ParamErrCode, "request params is wrong")
	NameErr             = NewErr(ParamErrCode, fmt.Sprintf("Wrong name has been given,max length is %d", constants.UserMaxLength))
	PasswordErr         = NewErr(ParamErrCode, "Wrong password has been given")
	UserNotExistErr     = NewErr(UserNotExistErrCode, "User does not exists")
	UserAlreadyExistErr = NewErr(UserAlreadyExistErrCode, "User already exists")
	IdTypeErr           = NewErr(AssertTypeErrCode, "given userID type is not int64")
	DataEmptyErr        = NewErr(DataEmptyErrCode, "it is empty")
)

// ConvertErr convert error to Err
//only if the err is a Err,return it immediately,else regard it as ServiceErr
func ConvertErr(err error) Err {
	Er := Err{}
	if errors.As(err, &Er) {
		return Er
	}
	s := ServiceErr
	s.ErrMsg = err.Error()
	return s
}
