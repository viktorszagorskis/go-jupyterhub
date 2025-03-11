package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func proxyToLab(w http.ResponseWriter, r *http.Request) {
    userID := getSessionUserID(r)
    container := getContainerForUser(userID)

    target, _ := url.Parse(container.URL)
    proxy := httputil.NewSingleHostReverseProxy(target)
    proxy.ServeHTTP(w, r)
}
