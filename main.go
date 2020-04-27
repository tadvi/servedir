package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	user := flag.String("username", "admin", "username for login")
	pass := flag.String("password", "admin8989", "password")

	dir := flag.String("dir", "download/", "directory to serve")
	port := flag.Int("port", 8989, "port to listen on")
	setNoCache := flag.Bool("no-cache", false, "set no-cache header on requests")

	flag.Parse()

	log.Println("Serving files from", *dir, "on port", *port)

	fserver := http.FileServer(http.Dir(*dir))

	http.HandleFunc("/", BasicAuth(user, pass, func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.URL)
		if *setNoCache {
			w.Header().Set("Cache-Control", "no-cache")
		}

		fserver.ServeHTTP(w, r)
	}))

	portStr := fmt.Sprintf(":%d", *port)
	err := http.ListenAndServe(portStr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func BasicAuth(user, pass *string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		u, p, ok := req.BasicAuth()
		if !ok || len(strings.TrimSpace(u)) < 1 || len(strings.TrimSpace(p)) < 1 {
			unauthorised(w)
			return
		}

		// This is a dummy check for credentials.
		if u != *user || p != *pass {
			unauthorised(w)
			return
		}

		// If required, Context could be updated to include authentication
		// related data so that it could be used in consequent steps.
		handler(w, req)
	}
}

func unauthorised(rw http.ResponseWriter) {
	rw.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
	rw.WriteHeader(http.StatusUnauthorized)
}
