package youzan_gin_middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"log"
	"net"
	"net/http"
	"testing"
	"time"
)

func init() {
	testService := createTestServer()
	testListener, _ := net.Listen("tcp", ":9334")
	testSv := &http.Server{Handler: testService}
	go func() { log.Fatal(testSv.Serve(testListener)) }()
}

func TestMiddleWare_youzan(t *testing.T) {
	codeSecond := sendRequest(http.MethodGet, "http://localhost:9334/youzan_webhook", map[string]string{
		"Client-Id": "test1",
	})
	assert.Equal(t, 200, codeSecond)

	codeSecond = sendRequest(http.MethodGet, "http://localhost:9334/youzan_webhook", map[string]string{
		"Client-Id": "test",
		"Event-Sign": "",
	})
	assert.Equal(t, 200, codeSecond)

	codeSecond = sendRequest(http.MethodGet, "http://localhost:9334/youzan_webhook", map[string]string{
		"Client-Id": "test",
		"Event-Sign": "05a671c66aefea124cc08b76ea6d30bb",
	})
	assert.Equal(t, 200, codeSecond)
}

func createLoginServerMock() *gin.Engine {
	r := gin.New()
	r.GET("/auth/me", func(context *gin.Context) {
		value := context.Request.Header["Authorization"]
		if value[0] == "test" {
			context.JSON(http.StatusOK, gin.H{})
			context.Done()
			return
		}
		context.AbortWithStatus(401)
	})
	return r
}

func createTestServer() *gin.Engine {
	r := gin.New()
	r.Use(YouZanAPIAuth("test", "test"))
	r.GET("/youzan_webhook", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})
	return r
}

func sendRequest(method string, url string, headers map[string]string) int {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	req, _ := http.NewRequest(method, url, nil)
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	res, _ := httpClient.Do(req)
	defer res.Body.Close()
	return res.StatusCode
}
