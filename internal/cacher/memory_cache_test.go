package cacher

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_MemCacheService(t *testing.T) {
	t.Run("Should have nil response in first", func(t *testing.T) {
		service := NewMemoryCacheService()
		req := httptest.NewRequest("GET", "localhost:8080/qqq", nil)
		storedRes := service.GetResponse(req)
		if storedRes != nil {
			t.Errorf("Expected response to be nil first")
		}
	})

	t.Run("Should persist response", func(t *testing.T) {
		service := NewMemoryCacheService()
		req := httptest.NewRequest("GET", "localhost:8080/qqq", nil)
		res := &http.Response{
			StatusCode: 404,
			Status:     "404 Not found",
		}

		service.SaveResponse(req, res)
		storedRes := service.GetResponse(req)

		if storedRes.StatusCode != 404 {
			t.Errorf("Expected status code: %d; Got: %d", 404, storedRes.StatusCode)
		}
	})

	t.Run("Should persist response body", func(t *testing.T) {
		testStr := "JSON CONtent here 22!!"
		service := NewMemoryCacheService()
		req := httptest.NewRequest("GET", "localhost:8080/qqq", nil)
		res := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(testStr))),
		}

		service.SaveResponse(req, res)
		storedRes := service.GetResponse(req)

		bodyBytes, _ := io.ReadAll(storedRes.Body)
		bodyStr := string(bodyBytes)
		if bodyStr != testStr {
			t.Errorf("Expected content %s; Got %s", testStr, bodyStr)
		}
	})

	// TODO test original response body still accessible after save
}
