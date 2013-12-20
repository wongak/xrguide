package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func connectDb(c *Config) (*sql.DB, error) {
	if len(c.Db.Dsn) != 1 {
		return nil, fmt.Errorf("Currently there is only one DSN supported.")
	}
	var conn *sql.DB
	var err error
	for dr, dsn := range c.Db.Dsn {
		conn, err = sql.Open(dr, dsn)
		if err != nil {
			return nil, fmt.Errorf("Error connecting to sql: %v", err)
		}
	}
	return conn, nil
}
