package main

import (
	"net/http"
)

// In-memory session store (replace with Redis/Postgres in production)
var userSessions = make(map[string]*ContainerInfo)

// setSession stores the user-container mapping
func setSession(w http.ResponseWriter, userID string, container *ContainerInfo) {
    userSessions[userID] = container
    http.SetCookie(w, &http.Cookie{Name: "user_id", Value: userID})
}

// getSessionUserID retrieves the user ID from the session
func getSessionUserID(r *http.Request) string {
    c, err := r.Cookie("user_id")
    if err != nil {
        return ""
    }
    return c.Value
}

// getContainerForUser retrieves a user's assigned container
func getContainerForUser(userID string) *ContainerInfo {
    return userSessions[userID]
}
