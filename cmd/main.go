package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World")
}

func main() {
	// TODO: There must be a way to do this without mux?
	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	r.HandleFunc("/{path:.*}", func(w http.ResponseWriter, r *http.Request) {
		requestedPath := "public/" + mux.Vars(r)["path"]

		if requestedPath == "" || strings.HasSuffix(requestedPath, "/") {
			requestedPath += "index"
		}

		var contentType string

		if !strings.HasSuffix(requestedPath, ".html") && !strings.HasSuffix(requestedPath, ".md") {
			requestedPath += ".html"
		}

		if strings.HasSuffix(requestedPath, ".md") {
			contentType = "text/plain; charset=utf-8"
		}

		content, err := ioutil.ReadFile(requestedPath)
		if err != nil {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", contentType)
		w.Write(content)

	})

	http.HandleFunc("/api/hello/", Hello)

	http.Handle("/", r)

	log.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
