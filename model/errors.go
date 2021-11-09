package model

import "strconv"

type ErrorCode string

type Error struct {
	Code    int
	Message string
}

func (err *Error) Error() string {
	return strconv.Itoa(err.Code) + " - " + err.Message
}

var (
	ErrUnknown          = &Error{-1, "unknown error"}         // 未知错误
	ErrDataNotExist     = &Error{400001, "data not exist"}    // 数据不存在
	ErrInvalidFormat    = &Error{400002, "invalid format"}    // 格式错误
	ErrPermissionDenied = &Error{400003, "permission denied"} // 权限不足
	ErrDuplicateData    = &Error{400004, "duplicate data"}    // 重复数据
	ErrOverLengthData   = &Error{400005, "over length data"}  // 数据不存在
	ErrDatabase         = &Error{500001, "repository error"}  // 未知数据库异常
)
