package impl

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
