package main

import (
	"encoding/xml"
	"flag"
)

var dbFile = flag.String("db", "xrguide.db", "Database file.")
var rebuild = flag.Bool("r", false, "Whether to reinitialize db.")
var textDir = flag.String("l", ".", "Libraries directory.")
var verbose = flag.Bool("v", false, "Verbose output.")

func main() {
	flag.Parse()
}

type Wares struct {
	XMLName xml.Name `xml:"wares"`
}

type Ware struct {
	Id        string `xml:"id,attr"`
	Name      string `xml:"name,attr"`
	Transport string `xml:"transport,attr"`
	Size      string `xml:"size,attr"`
	Volume    int    `xml:"volume,attr"`
}

type Price struct {
}
