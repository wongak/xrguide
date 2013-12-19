package main

import (
	"database/sql"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"
	"xrguide/db/schema"
	"xrguide/importing"
)

const MAX_WORKERS = 50

var dbType = flag.String("dbt", "sqlite3", "Database type.")
var dsn = flag.String("dsn", "xrguide.db", "Database DSN.")
var rebuild = flag.Bool("r", false, "Whether to reinitialize db.")
var textDir = flag.String("t", ".", "Directory with text files.")
var verbose = flag.Bool("v", false, "Verbose output.")
var lang = flag.Int64("l", 0, "Language Id. If not specified all.")
var page = flag.Int64("p", 0, "Page Id. If not specified all.")
var workers = flag.Int("w", 10, "Number of pages to process concurrently.")

func main() {
	flag.Parse()

	if *workers <= 0 {
		*workers = 10
	}
	if *workers > MAX_WORKERS {
		*workers = MAX_WORKERS
	}

	var database *importing.ImportDb

	if *dbType == "sqlite3" {
		err := importing.BackupDb(*dsn)
		if err != nil {
			log.Fatal(err)
		}
	}
	database, err := importing.Connect(*dbType, *dsn)
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
	err = read(database, *textDir, *verbose, *lang, *page)
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

type workload struct {
	Lang LangFile
	Page Page
}

func read(database *importing.ImportDb, directory string, verbose bool, useLang, usePage int64) error {
	var working sync.WaitGroup

	work := make(chan *workload, *workers)

	for i := 0; i < *workers; i++ {
		db, err := database.Db()
		if err != nil {
			log.Panicf("Error opening db: %v", err)
		}
		go func(db *sql.DB) {
			for {
				select {
				case w := <-work:
					if w == nil {
						db.Close()
						return
					}
					working.Add(1)
					stmt, err := db.Prepare(schema.TextInsert)
					if err != nil {
						log.Panicf("Error preparing statement: %v", err)
					}
					reset, err := db.Prepare(schema.TextDeletePage)
					if err != nil {
						log.Panicf("Error preparing statement: %v", err)
					}
					_, err = reset.Exec(w.Lang.LangId, w.Page.Id)
					for err != nil && err == sqlite3.ErrLocked {
						time.Sleep(time.Second)
						_, err = reset.Exec(w.Lang.LangId, w.Page.Id)
					}
					if err != nil {
						log.Panicf("Error on reset. Aborting: %v", err)
					}
					for _, t := range w.Page.Entries {
						if verbose {
							log.Printf("Lang %d Page %d Text %d.", w.Lang.LangId, w.Page.Id, t.Id)
						}
						_, err = stmt.Exec(w.Lang.LangId, w.Page.Id, t.Id, t.Entry)
						if err != nil {
							log.Panicf("Error on insert. Aborting: %v", err)
						}
					}

					log.Printf("Finished Lang %d Page %d.", w.Lang.LangId, w.Page.Id)
					working.Done()
				}
			}
		}(db)
	}

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
			wl := &workload{lang, page}
			work <- wl
		}
		file.Close()
	}
	working.Wait()
	close(work)
	return nil
}

func prepareDb(database *importing.ImportDb) error {
	var err error
	err = database.SetIgnoreForeignKeys(true)
	if err != nil {
		return err
	}
	db, err := database.Db()
	if err != nil {
		return err
	}
	for _, sql := range schema.TextReset {
		_, err = db.Exec(*sql)
		if err != nil {
			return fmt.Errorf("Error preparing db: %v", err)
		}
	}
	err = database.SetIgnoreForeignKeys(false)
	if err != nil {
		return err
	}
	return nil
}
