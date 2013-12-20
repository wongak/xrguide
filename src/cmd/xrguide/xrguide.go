// X Rebirth Guide
//
// Server
package main

import (
	"flag"
	"log"
	"xrguide"
)

var addr = flag.String("a", ":8080", "Listen address")
var cfgFile = flag.String("c", "cfg/xrguide.cfg.json", "Path to config file.")

func main() {
	flag.Parse()
	cfg, err := loadConfig(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	db, err := connectDb(cfg)
	if err != nil {
		log.Fatal(err)
	}
	srv, err := xrguide.NewServer(db, cfg.Html.HtmlSrcDir)
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}
	log.Fatal(srv.ListenAndServe(*addr))
}
