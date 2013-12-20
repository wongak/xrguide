package main

import (
	"database/sql"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"xrguide/db/schema"
	"xrguide/importing"
	"xrguide/text"
)

var dbType = flag.String("dbt", "sqlite3", "Database backend type.")
var dsn = flag.String("dsn", "xrguide.db", "DSN")
var rebuild = flag.Bool("r", false, "Whether to reinitialize db.")
var libDir = flag.String("l", ".", "Libraries directory.")
var verbose = flag.Bool("v", false, "Verbose output.")

func main() {
	flag.Parse()

	// check for library
	waresFileName := filepath.Join(*libDir, "wares.xml")
	waresFileInfo, err := os.Stat(waresFileName)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error on stat wares.xml: %v", err)
	}
	if os.IsNotExist(err) {
		log.Fatal("wares.xml not found in ", waresFileName)
	}
	if waresFileInfo.IsDir() {
		log.Fatal("wares.xml is a directory.")
	}

	var database *importing.ImportDb
	if *dbType == "sqlite3" {
		err := importing.BackupDb(*dsn)
		if err != nil {
			log.Fatal(err)
		}
	}
	database, err = importing.Connect(*dbType, *dsn)
	if err != nil {
		log.Fatal(err)
	}
	db, err := database.OpenDb(*dsn)
	if err != nil {
		log.Fatal(err)
	}
	if *rebuild || database.RequireRebuild() {
		err = prepareDb(database)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer db.Close()
	err = read(database, waresFileName, *verbose)
	if err != nil {
		log.Fatal(err)
	}
}

type Wares struct {
	XMLName xml.Name `xml:"wares"`

	Wares []Ware `xml:"ware"`
}

type Default struct {
	XMLName xml.Name `xml:"wares"`
	Ware
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
	Max     int `xml:"max,attr"`
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

func read(database *importing.ImportDb, waresFileName string, verbose bool) error {
	waresFile, err := os.Open(waresFileName)
	if err != nil {
		return fmt.Errorf("Error opening wares file: %v", err)
	}
	defer waresFile.Close()
	dec := xml.NewDecoder(waresFile)

	db, err := database.Db()
	if err != nil {
		return fmt.Errorf("Error connecting to db: %v", err)
	}
	insertWare, err := db.Prepare(schema.WareInsertWare)
	if err != nil {
		return fmt.Errorf("Error preparing SQL: %v", err)
	}

	var wares Wares
	err = dec.Decode(&wares)
	if err != nil {
		return fmt.Errorf("Error decoding XML: %v", err)
	}
	var imported, skipped int
	for _, ware := range wares.Wares {
		var namePageId, nameTextId, descPageId, descTextId sql.NullInt64
		var rawName, specialist sql.NullString
		namePageId.Int64, nameTextId.Int64, err = text.ParseTextRef(ware.Name)
		if err != nil {
			rawName.String = ware.Name
			rawName.Valid = true
		} else {
			namePageId.Valid, nameTextId.Valid = true, true
		}
		descPageId.Int64, descTextId.Int64, err = text.ParseTextRef(ware.Description)
		if err == nil {
			descPageId.Valid = true
			descTextId.Valid = true
		}
		if ware.Specialist != "" {
			specialist.String = ware.Specialist
			specialist.Valid = true
		}
		_, err := insertWare.Exec(
			ware.Id,
			namePageId,
			nameTextId,
			descPageId,
			descTextId,
			rawName,
			ware.Transport,
			specialist,
			ware.Size,
			ware.Volume,
			ware.Price.Min,
			ware.Price.Average,
			ware.Price.Max,
		)
		if err != nil {
			log.Printf("DB error: %v", err)
			skipped++
		}
		imported++
	}
	return nil
}

func prepareDb(database *importing.ImportDb) error {
	var err error
	db, err := database.Db()
	if err != nil {
		return err
	}
	for _, sql := range schema.WaresReset {
		_, err = db.Exec(*sql)
		if err != nil {
			return fmt.Errorf("Error preparing db: %v", err)
		}
	}
	for _, sql := range schema.WaresCreateIndexes {
		_, err = db.Exec(sql)
		if err != nil {
			return fmt.Errorf("Error creating index: %v", err)
		}
	}
	return nil
}
