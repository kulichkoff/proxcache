package main

import (
	"fmt"
	"proxcache/pkg/server"
)

func main() {
	// TODO parse command options by flags

	// TODO run server
	proxy := server.NewProxyServer("https://dummyjson.com")
	fmt.Println("Cache Proxy Server started on http://localhost:8080")
	proxy.Serve(":8080")
}
