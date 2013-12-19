package main

import (
	"encoding/xml"
	"flag"
)

var dbType = flag.String("dbt", "sqlite3", "Database backend type.")
var dsn = flag.String("dsn", "xrguide.db", "DSN")
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
	Id          string `xml:"id,attr"`
	Name        string `xml:"name,attr"`
	Description string `xml:"description,attr"`
	Transport   string `xml:"transport,attr"`
	Specialist  string `xml:"specialist,attr"`
	Size        string `xml:"size,attr"`
	Volume      int    `xml:"volume,attr"`
	Tags        string `xml:"tags,attr"`

	Price      Price        `xml:"price"`
	Production []Production `xml:"production"`
}

type Price struct {
	Min     int `xml:"min,attr"`
	Average int `xml:"average,attr"`
	Max     int `xml:"average,attr"`
}

type Production struct {
	Time   int    `xml:"time,attr"`
	Amount int    `xml:"amount,attr"`
	Method string `xml:"method,attr"`
	Name   string `xml:"name,attr"`

	Effects   []ProductionEffect `xml:"effects>effect"`
	Primary   []ProductionWare   `xml:"primary>ware"`
	Secondary []ProductionWare   `xml:"secondary>ware"`

	Container   Container   `xml:"container"`
	Icon        Icon        `xml:"icon"`
	Restriction Restriction `xml:"resriction"`
}

type ProductionEffect struct {
	Type    string  `xml:"type,attr"`
	Product float32 `xml:"product,attr"`
}

type ProductionWare struct {
	Ware   string `xml:"ware,attr"`
	Amount int    `xml:"amount,attr"`
}

type Container struct {
	Ref string `xml:"ref,attr"`
}

type Icon struct {
	Active string `xml:"active,attr"`
}

type Restriction struct {
	Sell float32 `xml:"sell,attr"`
}
