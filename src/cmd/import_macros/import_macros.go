package main

import (
	"database/sql"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"xrguide/db/schema"
	"xrguide/importing"
	"xrguide/text"
)

var dbType = flag.String("dbt", "sqlite3", "Database backend type.")
var dsn = flag.String("dsn", "xrguide.db", "DSN")
var rebuild = flag.Bool("r", false, "Whether to reinitialize db.")
var macrosDirName = flag.String("s", ".", "Structures directory.")
var verbose = flag.Bool("v", false, "Verbose output.")

func main() {
	flag.Parse()

	// check for library
	macrosDir, err := os.Stat(*macrosDirName)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal("Error on stat macros directory: %v", err)
	}
	if os.IsNotExist(err) {
		log.Fatal("Macros directory not found in ", *macrosDirName)
	}
	if !macrosDir.IsDir() {
		log.Fatal("%s is not a directory.", *macrosDirName)
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
	err = read(database, *macrosDirName, *verbose)
	if err != nil {
		log.Fatal(err)
	}
}

type Macros struct {
	XMLName xml.Name `xml:"macros"`
	Macros  []Macro  `xml:"macro"`
}

type Macro struct {
	Id          string        `xml:"name,attr"`
	Class       string        `xml:"class,attr"`
	Ident       MacroIdent    `xml:"properties>identification"`
	Production  *Production   `xml:"properties>production"`
	Connections []*Connection `xml:"connections>connection"`

	FileName string
}

type MacroIdent struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"description,attr"`
	Unique      bool   `xml:"unique,attr"`
}

type Production struct {
	Wares string `xml:"wares,attr"`
}

type Connection struct {
	Build ConnectionBuild `xml:"build"`
	Macro ConnectionMacro `xml:"macro"`
}

type ConnectionBuild struct {
	Mode     *string `xml:"mode,attr"`
	Group    *string `xml:"group,attr"`
	Sequence *string `xml:"sequence,attr"`
	Stage    *int    `xml:"stage,attr"`
}

type ConnectionMacro struct {
	Ref        *string                `xml:"ref,attr"`
	Components []*ConnectionComponent `xml:"component"`
}

type ConnectionComponent struct {
	Ref string `xml:"ref,attr"`
}

func read(database *importing.ImportDb, macrosDirName string, verbose bool) error {
	db, err := database.Db()
	if err != nil {
		return err
	}
	macroFilePattern := regexp.MustCompile(".*?_macro\\.xml")
	var files []string
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if !macroFilePattern.MatchString(info.Name()) {
			return nil
		}
		files = append(files, path)
		return nil
	}
	err = filepath.Walk(macrosDirName, walkFunc)
	if err != nil {
		return fmt.Errorf("Error finding macro XMLs: %v", err)
	}
	insertStmt, err := db.Prepare(schema.MacrosInsertMacro)
	if err != nil {
		return fmt.Errorf("Error preparing insert stmt: %v", err)
	}
	insertProd, err := db.Prepare(schema.MacrosInsertProduction)
	if err != nil {
		return fmt.Errorf("Error preparing insert stmt: %v", err)
	}
	insertConn, err := db.Prepare(schema.MacrosInsertConnection)
	if err != nil {
		return fmt.Errorf("Error preparing insert stmt: %v", err)
	}
	for _, m := range files {
		mFile, err := os.Open(m)
		if err != nil {
			return fmt.Errorf("Error opening %s: %v", m, err)
		}
		macros := new(Macros)
		decoder := xml.NewDecoder(mFile)
		err = decoder.Decode(macros)
		if err != nil {
			mFile.Close()
			return fmt.Errorf("Error decoding: %v", err)
		}
		mFile.Close()
		fmt.Printf("%s\n", m)
		for _, macro := range macros.Macros {
			var namePageId, nameTextId, descPageId, descTextId sql.NullInt64
			namePageId.Int64, nameTextId.Int64, err = text.ParseTextRef(macro.Ident.Name)
			if err == nil {
				namePageId.Valid, nameTextId.Valid = true, true
			}
			descPageId.Int64, descTextId.Int64, err = text.ParseTextRef(macro.Ident.Description)
			if err == nil {
				descPageId.Valid, descTextId.Valid = true, true
			}
			_, err = insertStmt.Exec(macro.Id, macro.Class, namePageId, nameTextId, descPageId, descTextId, macro.Ident.Unique, filepath.Base(m))
			if err != nil {
				return fmt.Errorf("Error on insert: %v", err)
			}
			if macro.Production != nil {
				wares := strings.Split(macro.Production.Wares, " ")
				for _, ware := range wares {
					_, err = insertProd.Exec(macro.Id, ware)
					if err != nil {
						return fmt.Errorf("Error on insert: %v", err)
					}
				}
			}
			if macro.Connections != nil {
				for _, conn := range macro.Connections {
					// skip components for now
					if conn.Macro.Ref == nil {
						continue
					}
					_, err = insertConn.Exec(
						macro.Id,
						conn.Macro.Ref,
						conn.Build.Mode,
						conn.Build.Group,
						conn.Build.Sequence,
						conn.Build.Stage,
					)
					if err != nil {
						return fmt.Errorf("Error on insert connection: %v", err)
					}
				}
			}
		}
	}
	return nil
}

func prepareDb(database *importing.ImportDb) error {
	var err error
	db, err := database.Db()
	if err != nil {
		return err
	}
	for _, sql := range schema.MacrosReset {
		_, err = db.Exec(*sql)
		if err != nil {
			return fmt.Errorf("Error preparing db: %v", err)
		}
	}
	for _, sql := range schema.MacrosCreateIndexes {
		_, err = db.Exec(sql)
		if err != nil {
			return fmt.Errorf("Error creating index: %v", err)
		}
	}
	return nil
}
