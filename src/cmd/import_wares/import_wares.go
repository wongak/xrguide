package main

import (
	"flag"
)

var dbFile = flag.String("db", "xrguide.db", "Database file.")
var rebuild = flag.Bool("r", false, "Whether to reinitialize db.")
var textDir = flag.String("l", ".", "Libraries directory.")
var verbose = flag.Bool("v", false, "Verbose output.")

func main() {
	flag.Parse()
}
