package manage

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (e *Engine) BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, username, password := authorization(c.Request.Header.Get("Authorization"))
		if !auth || username != e.config.Username || password != e.config.Password {
			c.Header("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(gin.AuthUserKey, username)
	}
}

func authorization(encodeString string) (auth bool, username string, password string) {
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
