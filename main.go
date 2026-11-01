package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("❌ ERROR: PORT environment variable is not set!")
		os.Exit(1)
	}
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("IPTV Control Center Go Version | Running OK"))
	})

	fmt.Printf("✅ Server listen on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("❌ Listen failed: %v\n", err)
		os.Exit(1)
	}
}
