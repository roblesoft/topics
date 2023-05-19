package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/roblesoft/topics/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	var (
		testSuite = mock.NewTestSuit("3000")
		server    = &Server{Port: testSuite.Port, Service: testSuite.Service}
		router    = gin.Default()
	)

	router.GET("/api/v1/healthcheck", server.HealthCheck)
	server.Router = router

	t.Run("http status ok", func(t *testing.T) {
		var (
			w      = httptest.NewRecorder()
			req, _ = http.NewRequest("GET", "/api/v1/healthcheck", nil)
		)
		server.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
