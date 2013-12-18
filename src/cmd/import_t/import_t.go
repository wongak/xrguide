package main

import (
	"database/sql"
	"encoding/xml"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"xrguide/db/schema"
	"xrguide/importing"
)

var dbFile = flag.String("db", "xrguide.db", "Database file.")
var rebuild = flag.Bool("r", false, "Whether to reinitialize db.")
var textDir = flag.String("t", ".", "Directory with text files.")
var verbose = flag.Bool("v", false, "Verbose output.")
var lang = flag.Int64("l", 0, "Language Id. If not specified all.")
var page = flag.Int64("p", 0, "Page Id. If not specified all.")

func main() {
	flag.Parse()
	err := importing.BackupDb(*dbFile)
	if err != nil {
		log.Fatal(err)
	}
	db, err := importing.OpenDb(*dbFile, rebuild)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if *rebuild {
		err = prepareDb(db)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = read(db, *textDir, *verbose, *lang, *page)
	if err != nil {
		log.Fatal(err)
	}
}

type Text struct {
	Id    int64  `xml:"id,attr"`
	Entry string `xml:",innerxml"`
}

type Page struct {
	Id      int64  `xml:"id,attr"`
	Entries []Text `xml:"t"`
}

type LangFile struct {
	XMLName xml.Name `xml:"language"`
	LangId  int64    `xml:"id,attr"`
	Pages   []Page   `xml:"page"`
}

func read(db *sql.DB, directory string, verbose bool, useLang, usePage int64) error {
	info, err := os.Stat(directory)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Could not stat text directory: %v", err)
	}
	if os.IsNotExist(err) {
		return fmt.Errorf("Text directory not found: %s", directory)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory.", directory)
	}
	dir, err := os.Open(directory)
	if err != nil {
		return fmt.Errorf("Error opening text directory: %v", err)
	}
	defer dir.Close()
	pattern := regexp.MustCompile("0001-L(\\d{3})\\.xml")
	stmt, err := db.Prepare(schema.TextInsert)
	if err != nil {
		return fmt.Errorf("Error preparing statement: %v", err)
	}
	reset, err := db.Prepare(schema.TextDeletePage)
	if err != nil {
		return fmt.Errorf("Error preparing statement: %v", err)
	}
	var langId int64
	var fileName string
	var lang LangFile
	for {
		f, err := dir.Readdir(1)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("Error on reading text directory: %v")
		}
		fileName = filepath.Join(dir.Name(), f[0].Name())
		if !pattern.MatchString(f[0].Name()) {
			log.Printf("Skipping %s", fileName)
			continue
		}
		matches := pattern.FindStringSubmatch(f[0].Name())
		langId, err = strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			log.Printf("Error on parsing filename: %s (%v)", fileName, err)
			continue
		}
		if langId == 0 {
			log.Printf("Error on filename %s. Invalid lang id.", fileName)
			continue
		}
		log.Printf("Text file %s. Language Id %d", fileName, langId)

		if useLang != 0 && useLang != langId {
			log.Printf("Skipping %s.", fileName)
			continue
		}

		file, err := os.Open(fileName)
		if err != nil {
			return fmt.Errorf("Error opening file %s.", fileName)
		}
		decoder := xml.NewDecoder(file)
		err = decoder.Decode(&lang)
		if err != nil {
			file.Close()
			return fmt.Errorf("Error decoding file %s: %v", fileName, err)
		}
		for _, page := range lang.Pages {
			if lang.LangId != langId {
				log.Printf("Language Id does not match in %s. Id %d.", fileName, lang.LangId)
				break
			}
			if usePage != 0 && usePage != page.Id {
				if verbose {
					log.Printf("Skipping page %d.", page.Id)
				}
				continue
			}
			if verbose {
				log.Printf("Lang %d Page %d.", lang.LangId, page.Id)
			}
			_, err = reset.Exec(lang.LangId, page.Id)
			if err != nil {
				return fmt.Errorf("Error on reset. Aborting: %v", err)
			}
			for _, t := range page.Entries {
				if verbose {
					log.Printf("Lang %d Page %d Text %d.", lang.LangId, page.Id, t.Id)
				}
				_, err = stmt.Exec(lang.LangId, page.Id, t.Id, t.Entry)
				if err != nil {
					return fmt.Errorf("Error on insert. Aborting: %v", err)
				}
			}
		}
		file.Close()
	}
	return nil
}

func prepareDb(db *sql.DB) error {
	var err error
	for _, sql := range schema.TextReset {
		_, err = db.Exec(sql)
		if err != nil {
			return fmt.Errorf("Error preparing db: %v", err)
		}
	}
	return nil
}
