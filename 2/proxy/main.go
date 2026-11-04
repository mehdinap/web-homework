// A minimal reverse proxy that authenticates every request with a JWT (HS256)
// before forwarding it to the legacy service on the internal network.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	listenAddr  = env("LISTEN_ADDR", ":8080")
	upstreamURL = env("UPSTREAM_URL", "http://legacy-httpbin:8081")
	secret      = []byte(env("JWT_SECRET", "dev-only-secret-change-me"))
)

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func main() {
	target, err := url.Parse(upstreamURL)
	if err != nil {
		log.Fatalf("invalid UPSTREAM_URL: %v", err)
	}

	// Everything that is not the login endpoint goes through the JWT check and
	// then straight to the legacy service.
	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/login", login)
	mux.Handle("/", auth(httputil.NewSingleHostReverseProxy(target)))

	log.Printf("listening on %s, forwarding to %s", listenAddr, target)
	log.Fatal(http.ListenAndServe(listenAddr, mux))
}

// login hands out a token for the demo user so the proxy can be tried by hand.
func login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if body.Username != env("LOGIN_USER", "student") || body.Password != env("LOGIN_PASSWORD", "hunter2") {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   body.Username,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	}).SignedString(secret)
	if err != nil {
		http.Error(w, "could not sign token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"access_token": token})
}

// auth rejects requests that do not carry a valid bearer token.
func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || raw == "" {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		// Pin the algorithm to HS256 so a token claiming "alg: none" or an
		// asymmetric algorithm can never be accepted.
		token, err := jwt.Parse(raw, func(t *jwt.Token) (any, error) {
			return secret, nil
		}, jwt.WithValidMethods([]string{"HS256"}))
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Keep the token out of the internal hop.
		r.Header.Del("Authorization")
		next.ServeHTTP(w, r)
	})
}
