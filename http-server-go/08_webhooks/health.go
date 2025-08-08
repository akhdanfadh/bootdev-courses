package main

import "net/http"

// handlerReadiness is a function to handle readiness endpoint (health check)
func handlerReadiness(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK))) // just return 200 OK
}
