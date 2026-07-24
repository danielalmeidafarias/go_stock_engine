package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestRequestTimeoutMiddlewareReturnsGatewayTimeout(t *testing.T) {
	r := gin.New()
	r.Use(requestTimeoutMiddleware(true, time.Nanosecond))
	r.GET("/", func(c *gin.Context) {
		<-c.Request.Context().Done()
		respondTimeout(c)
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("status = %d, want %d", w.Code, http.StatusGatewayTimeout)
	}
}
