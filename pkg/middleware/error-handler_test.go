package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestErrorHandlerMiddlewareSingleResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(NewErrorHandlerMiddleware().Add())
	router.GET("/test", func(c *gin.Context) {
		c.Error(errors.New("first error"))
		c.Error(errors.New("second error"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	body := w.Body.Bytes()
	if !json.Valid(body) {
		t.Fatalf("expected valid JSON body, got %q", string(body))
	}

	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("failed to unmarshal response body: %v", err)
	}

	if payload["error"] != "first error" {
		t.Fatalf("expected first error message, got %q", payload["error"])
	}
}
