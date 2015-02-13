package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"path"
)

//go:generate go-bindata data/ assets/...

var flagPort int
var flagStanzaBaseDir string

func init() {
	flag.IntVar(&flagPort, "port", 8080, "port to listen on")
	flag.StringVar(&flagStanzaBaseDir, "stanza-base-dir", ".", "stanza base directory")
}

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	assetsHandler := http.FileServer(http.Dir(flagStanzaBaseDir))

	sp, err := NewStanzaProvider(flagStanzaBaseDir)
	if err != nil {
		log.Fatal(err)
	}
	if err := sp.Build(); err != nil {
		log.Fatal(err)
	}

	assetsDir := "assets"
	log.Printf("generating assets under %s", path.Join(flagStanzaBaseDir, assetsDir))
	if err := RestoreAssets(flagStanzaBaseDir, assetsDir); err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		assetsHandler.ServeHTTP(w, req)
	})

	addr := fmt.Sprintf(":%d", flagPort)
	log.Println("listening on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
