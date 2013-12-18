package importing

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
)

func BackupDb(fileName string) error {
	info, err := os.Stat(fileName)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("Error backing up db stat: %v", err)
	}
	if os.IsNotExist(err) {
		return nil
	}
	if info.IsDir() {
		return fmt.Errorf("DB file is a directory. Cannot continue.")
	}
	orig, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("Could not open db: %v", err)
	}
	defer orig.Close()
	bak, err := os.Create(fileName + ".bak")
	if err != nil {
		return fmt.Errorf("Error creating backup file: %v", err)
	}
	defer bak.Close()
	_, err = io.Copy(bak, orig)
	return nil
}

func OpenDb(dbFileName string, rebuild *bool) (*sql.DB, error) {
	dbFileInfo, err := os.Stat(dbFileName)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("Error on stat db: %v", err)
	}
	// force rebuild on new db
	if os.IsNotExist(err) {
		*rebuild = true
	} else {
		if dbFileInfo.IsDir() {
			return nil, fmt.Errorf("Db is a directory.")
		}
	}
	db, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		return nil, fmt.Errorf("Error opening db: %v", err)
	}
	return db, nil
}

func Db(dbFileName string) (*sql.DB, error) {
	return sql.Open("sqlite3", "file:"+dbFileName+"?cache=shared")
}
