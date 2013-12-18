package importing

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
)

type ImportDb struct {
	requireRebuild bool
	dsn            string
	db             *sql.DB

	OpenDb func(dsn string) (*sql.DB, error)
	Db     func() (*sql.DB, error)
}

func (i ImportDb) RequireRebuild() bool {
	return i.requireRebuild
}

func Connect(dbType, dsn string) (*ImportDb, error) {
	var db ImportDb
	switch dbType {
	case "sqlite3":
		db.OpenDb = func(dsn string) (*sql.DB, error) {
			db.dsn = dsn
			dbFileInfo, err := os.Stat(dsn)
			if err != nil && !os.IsNotExist(err) {
				return nil, fmt.Errorf("Error on stat db: %v", err)
			}
			// force rebuild on new db
			if os.IsNotExist(err) {
				db.requireRebuild = true
			} else {
				if dbFileInfo.IsDir() {
					return nil, fmt.Errorf("Db is a directory.")
				}
			}
			sqlite, err := sql.Open("sqlite3", dsn)
			if err != nil {
				return nil, fmt.Errorf("Error opening db: %v", err)
			}
			return sqlite, err
		}
		db.Db = func() (*sql.DB, error) {
			return sql.Open("sqlite3", "file:"+db.dsn+"?cache=shared")
		}
	case "mysql":
		db.OpenDb = func(dsn string) (*sql.DB, error) {
			var err error
			db.dsn = dsn
			db.db, err = sql.Open("mysql", dsn)
			if err != nil {
				return nil, fmt.Errorf("Error connecting to MySQL: %v", err)
			}
			return db.db, nil
		}
		db.Db = func() (*sql.DB, error) {
			return db.db, nil
		}
	default:
		return nil, fmt.Errorf("Invalid db type: %s", dbType)
	}
	return &db, nil
}

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
