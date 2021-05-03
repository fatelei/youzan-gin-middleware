package youzan_gin_middleware

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
)

func YouZanAPIAuth(clientID string, clientSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentClientID := c.Request.Header.Get("Client-Id")
		eventSign := c.Request.Header.Get("Event-Sign")
		if currentClientID != clientID {
			c.AbortWithStatusJSON(200, gin.H{"code": 0, "msg": "success"})
			log.Printf("invalid client id: %s", currentClientID)
			return
		}

		sign := md5.New()
		rawBody, _ := ioutil.ReadAll(c.Request.Body)
		_, _ = io.WriteString(sign, clientID)
		_, _ = io.WriteString(sign, string(rawBody))
		_, _ = io.WriteString(sign, clientSecret)
		currentSign := sign.Sum(nil)
		if fmt.Sprintf("%x", currentSign) != eventSign {
			c.AbortWithStatusJSON(200, gin.H{"code": 0, "msg": "success"})
			log.Printf("invalid event sign: %s", eventSign)
			return
		}
		c.Next()
	}
}
