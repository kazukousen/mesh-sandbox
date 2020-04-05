package main

import (
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/echo", echoEnv())
	http.ListenAndServe(":8080", mux)
}

func echoEnv() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := os.Getenv("MSG")
		w.Write([]byte(msg))
	})
}
