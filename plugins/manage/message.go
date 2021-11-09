package manage

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/ruixiaoedu/hiot/model"
	"net/http"
	"regexp"
)

// MessageController 消息控制器
type MessageController manage

// NewMessageController 新建MessageController
func NewMessageController(m manage) MessageController {
	return MessageController(m)
}

// Routes 创建路由
func (ctr MessageController) Routes(route *gin.RouterGroup) {
	route = route.Group("/message")
	route.POST("/publish", ctr.publish)
}

// MessagePublishRequest 消息发布
type MessagePublishRequest struct {
	Topic   *string `json:"topic" binding:"required"` // 要接收消息的设备的自定义Topic，格式必须为usr/{productId}/{deviceId}/{topics...}
	Qos     *byte   `json:"qos"`                      // 指定消息的发送方式
	Payload *string `json:"payload"`                  // 要发送的消息主体，base64编码
}

// publish 发布信息
func (ctr *MessageController) publish(c *gin.Context) {
	var (
		req MessagePublishRequest
		err error
	)
	if err = c.ShouldBindJSON(&req); err != nil {
		if i, ok := err.(validator.ValidationErrors); ok {
			fmt.Println("Error" + i.Error())
		}
		c.JSON(http.StatusBadRequest, failWithError(model.ErrInvalidFormat))
		return
	}

	var payload []byte
	if req.Payload != nil && *req.Payload != "" {
		if payload, err = base64.StdEncoding.DecodeString(*req.Payload); err != nil {
			c.JSON(http.StatusBadRequest, failWithError(model.ErrInvalidFormat))
			return
		}
	}

	if len(payload) > 4*1024*1024 {
		c.JSON(http.StatusBadRequest, failWithError(model.ErrOverLengthData))
		return
	}

	usrRegexp := regexp.MustCompile(`^usr/([\d\w-_]{1,32})/([\d\w-_]{1,32})/`)
	params := usrRegexp.FindStringSubmatch(*req.Topic)
	if len(params) != 3 {
		c.JSON(http.StatusBadRequest, failWithError(model.ErrInvalidFormat))
		return
	}
	productId := params[1]
	deviceId := params[2]

	if _, err = ctr.engine.DB().Device().Get(productId, deviceId); err != nil {
		c.JSON(http.StatusBadRequest, failWithError(err))
		return
	}

	// 这里不检查topic可用性，由下一级检查
	err = ctr.engine.Core().Publish(*req.Topic, *req.Qos, payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, failWithError(err))
		return
	}

	c.JSON(http.StatusOK, success(nil))
}
