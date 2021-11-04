package manage

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/ruixiaoedu/hiot/model"
	"net/http"
	"strings"
)

// 该文件负责授权处理

// BasicAuth 授权方法
func (e *Engine) BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, username, password := decodeAuthorization(c.Request.Header.Get("Authorization"))
		if !auth || username != e.config.Username || password != e.config.Password {
			c.Header("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			c.AbortWithStatusJSON(http.StatusUnauthorized, failWithError(model.ErrPermissionDenied))
			return
		}
		c.Set(gin.AuthUserKey, username)
	}
}

func decodeAuthorization(encodeString string) (auth bool, username string, password string) {
	if !strings.HasPrefix(encodeString, "Basic ") {
		return false, "", ""
	}
	encodeString = strings.TrimPrefix(encodeString, "Basic ")
	base, err := base64.StdEncoding.DecodeString(encodeString)
	if err != nil {
		return false, "", ""
	}
	splitBase := strings.SplitN(string(base), ":", 2)
	if len(splitBase) != 2 {
		return false, "", ""
	}

	return true, splitBase[0], splitBase[1]
}
