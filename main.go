package main

import (
	"log"
	"net/http"
)

func main() {
    http.HandleFunc("/login", handleLogin)
    http.HandleFunc("/callback", handleCallback)
    http.HandleFunc("/lab/", proxyToLab)

    log.Println("GoHub running on :8080")
    http.ListenAndServe(":8080", nil)
}
