// Proxy server MUST add the headers to the response that indicate whether
// the response is from the cache or the server:
//
// # If the response is from the cache
// X-Cache: HIT
//
// # If the response is from the origin server
// X-Cache: MISS

package server

import (
	"net/http"
	"proxcache/internal/cacher"
)

type CacheService interface {
	SaveResponse(req *http.Request, res *http.Response) error
	GetResponse(req *http.Request) *http.Response
}

type ProxyServer struct {
	originURL    string
	router       *http.ServeMux
	cacheService CacheService
}

func (s *ProxyServer) Serve(addr string) error {
	s.router.HandleFunc("/", s.handleHTTP)

	return http.ListenAndServe(addr, s.router)
}

func (s *ProxyServer) handleHTTP(w http.ResponseWriter, req *http.Request) {
	res := s.cacheService.GetResponse(req)
	if res != nil {
		res.Header.Set("X-Cache", "HIT")
		s.writeResponse(w, res)
		return
	}

	originReq, _ := http.NewRequest(req.Method, s.originURL+req.URL.Path, nil)
	client := http.Client{}
	res, _ = client.Do(originReq)
	s.cacheService.SaveResponse(req, res)

	res.Header.Set("X-Cache", "MISS")
	s.writeResponse(w, res)
}

func (s *ProxyServer) writeResponse(w http.ResponseWriter, res *http.Response) {
	for key, values := range res.Header {
		for _, value := range values {
			w.Header()
			w.Header().Add(key, value)
		}
	}

	// TODO write body

	w.WriteHeader(res.StatusCode)
}

func NewProxyServer(origin string) *ProxyServer {
	cacheService := cacher.NewMemoryCacheService()
	mux := http.NewServeMux()

	return &ProxyServer{
		originURL:    origin,
		router:       mux,
		cacheService: cacheService,
	}
}
