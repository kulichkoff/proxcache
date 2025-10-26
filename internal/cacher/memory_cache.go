package cacher

import (
	"bytes"
	"crypto/sha256"
	"io"
	"net/http"
)

type MemoryCacheService struct {
	cacheMap map[string]*http.Response
	bodyMap  map[string][]byte
}

func (s *MemoryCacheService) SaveResponse(req *http.Request, res *http.Response) error {
	reqHash := s.hashRequest(req)

	if res.Body != nil {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			// TODO wrap with own errors
			return err
		}
		s.bodyMap[reqHash] = bodyBytes
		res.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	s.cacheMap[reqHash] = res

	return nil
}

func (s *MemoryCacheService) GetResponse(req *http.Request) *http.Response {
	reqHash := s.hashRequest(req)
	res := s.cacheMap[reqHash]
	bodyBytes, hasBody := s.bodyMap[reqHash]
	if hasBody {
		res.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}
	return res
}

func (s *MemoryCacheService) hashRequest(req *http.Request) string {
	h := sha256.New()

	h.Write([]byte(req.Method + req.URL.Path))

	return string(h.Sum(nil))
}

func NewMemoryCacheService() *MemoryCacheService {
	return &MemoryCacheService{
		cacheMap: make(map[string]*http.Response, 64),
		bodyMap:  make(map[string][]byte, 64),
	}
}
