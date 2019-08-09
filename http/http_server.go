package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/g-airport/go-infra/graceful"
)

// this used to test graceful close

func main() {
	// Declare the router
	m := mux.NewRouter()

	// Bind the API endpoints to router
	m.HandleFunc("/v1/graceful", GracefulHandler()).Methods(http.MethodGet)

	// Listen and Serve

	srv := &http.Server{
		Handler: m,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func GracefulHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		graceful.AddOne()
		go func() {
			defer graceful.Done()
			time.Sleep(2 * time.Second)
			log.Println("this is test graceful success")
		}()
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("<p>Hi</p>"))
	}
}
