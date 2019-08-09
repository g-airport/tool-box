package http_client

import (
	"log"
	"net/http"
)

// ----------------------------

func LoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

// authMiddleware verifies the auth token on the request matches the
// one defined in the environment
func AuthMiddleware(next http.Handler) http.Handler {
	// authToken is the token that must be used on all requests
	authToken := getEnv("AUTH_TOKEN", "")

	// Return the Handlerfunc that asserts the auth token
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authToken != "" {
			if r.Header.Get("X-Auth-Token") == authToken {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
		return
	})
}

