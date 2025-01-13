package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DB      *sqlx.DB
	//go:embed frontend/dist
	frontend embed.FS
)

var buildMode = "prod"

func main() {
	port := flag.String("port", "8080", "Port on which to launch the app")
	flag.Parse()
	mux := http.NewServeMux()

	if buildMode == "dev" {
		frontendURL, err := url.Parse("http://localhost:5173")
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(frontendURL)
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		})
	} else {
		dist, _ := fs.Sub(frontend, "frontend/dist")
		mux.Handle("/", http.FileServer(http.FS(dist)))
	}

	server := http.Server{
		Addr:         fmt.Sprintf(":%s", *port),
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
		IdleTimeout:  time.Minute,
		Handler:      mux,
	}

	log.Println("Server listening on port 8080")
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
