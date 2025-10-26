package cacher

import (
	"crypto/sha256"
	"net/http"
)

type MemoryCacheService struct {
	cacheMap map[string]*http.Response
}

func (s *MemoryCacheService) SaveResponse(req *http.Request, res *http.Response) error {
	reqHash := s.hashRequest(req)

	// TODO persist body
	res.Body = nil
	s.cacheMap[reqHash] = res

	return nil
}

func (s *MemoryCacheService) GetResponse(req *http.Request) *http.Response {
	reqHash := s.hashRequest(req)
	return s.cacheMap[reqHash]
}

func (s *MemoryCacheService) hashRequest(req *http.Request) string {
	h := sha256.New()

	h.Write([]byte(req.Method + req.URL.Path))

	return string(h.Sum(nil))
}

func NewMemoryCacheService() *MemoryCacheService {
	return &MemoryCacheService{
		cacheMap: make(map[string]*http.Response, 64),
	}
}
