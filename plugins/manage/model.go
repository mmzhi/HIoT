package manage

import (
	"fmt"
	"strings"
	"time"
)

// Response 应答结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// fail 错误的应答
func fail(code int, message string) *Response {
	if message == "" {
		message = "Error"
	}
	if code == 0 {
		code = -1
	}
	return &Response{
		Code:    code,
		Message: message,
	}
}

// success 成功的应答
func success(data interface{}) *Response {
	return &Response{
		Code:    0,
		Message: "OK",
		Data:    data,
	}
}

// Page 应答结构体
type Page struct {
	Total   int `json:"total"`
	Size    int `json:"size"`
	Current int `json:"current"`
	Pages   int `json:"pages"`
}

// Datetime 日期时间控件
type Datetime struct {
	time.Time
}

func (dt Datetime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", dt.Time.Format(time.RFC3339))
	return []byte(stamp), nil
}

func (dt Datetime) UnmarshalJSON(data []byte) (err error) {
	s := strings.Trim(string(data), "\"")
	if s == "null" {
		dt.Time = time.Time{}
		return
	}
	t, err := time.Parse(time.RFC3339, s)
	dt.Time = t
	return
}

// PDatetime time转换成Datetime指针
func PDatetime(t *time.Time) *Datetime {
	if t == nil {
		return nil
	}
	dt := Datetime{*t}
	return &dt
}
