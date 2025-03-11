package main

import (
	"context"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

var (
    keycloakIssuer = "http://localhost:8081/realms/myrealm"
    clientID       = "gohub"
    clientSecret   = "9olf4EW3dNzDq3LHSd73zyb72gsx6iHu"
    redirectURL    = "http://localhost:8080/callback"

    provider      *oidc.Provider
    oauth2Config oauth2.Config
)

func init() {
    var err error
    provider, err = oidc.NewProvider(context.Background(), keycloakIssuer)
    if err != nil {
        log.Fatalf("Failed to connect to Keycloak: %v", err)
    }

    oauth2Config = oauth2.Config{
        ClientID:     clientID,
        ClientSecret: clientSecret,
        RedirectURL:  redirectURL,
        Endpoint:     provider.Endpoint(),
        Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
    }
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, oauth2Config.AuthCodeURL("state"), http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
    ctx := context.Background()
    code := r.URL.Query().Get("code")

    token, err := oauth2Config.Exchange(ctx, code)
    if err != nil {
        http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
        return
    }

    idToken, err := provider.Verifier(&oidc.Config{ClientID: clientID}).Verify(ctx, token.Extra("id_token").(string))
    if err != nil {
        http.Error(w, "Invalid ID Token", http.StatusUnauthorized)
        return
    }

    var claims struct {
        Email string `json:"email"`
    }
    if err := idToken.Claims(&claims); err != nil {
        http.Error(w, "Failed to parse claims", http.StatusInternalServerError)
        return
    }

    container := getOrSpawnContainer(claims.Email) // 🔹 Get the running JupyterLab container

    setSession(w, claims.Email, container) // 🔹 Store session info

    // 🔹 Redirect user to their specific JupyterLab container
    http.Redirect(w, r, container.URL, http.StatusFound)
}
