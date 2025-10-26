package cacher

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMemCacheServicePersistsResponses(t *testing.T) {
	service := NewMemoryCacheService()

	req := httptest.NewRequest("GET", "localhost:8080/qqq", nil)
	res := &http.Response{
		StatusCode: 404,
		Status:     "404 Not found",
	}

	storedRes := service.GetResponse(req)
	if storedRes != nil {
		// Case 1: Should be nil before the save
		t.Errorf("Expected response to be nil first")
	}

	service.SaveResponse(req, res)
	storedRes = service.GetResponse(req)

	if storedRes == nil {
		// Case 2: Should not be nil after the save
		t.Errorf("Expected response but got nil")

	} else if storedRes.StatusCode != 404 {
		// Case 3: Should be the same as saved
		t.Errorf("Expected code %d; Got %d", 404, storedRes.StatusCode)
	}
}
