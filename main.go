package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	// 读取botkeep自动注入的PORT环境变量
	port := os.Getenv("PORT")
	if port == "" {
		port = "57516"
	}
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	// 健康检查路由 /health，用来保活防休眠
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// 示例首页
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("IPTV Control Center Go Version | Running OK"))
	})

	fmt.Printf("✅ Server listen on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("❌ Start failed: %v\n", err)
		os.Exit(1)
	}
}
