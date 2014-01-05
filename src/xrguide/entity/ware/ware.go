package ware

import (
	"database/sql"
	"fmt"
	"xrguide/db/query"
)

type Ware struct {
	Id           string
	Name         sql.NullString
	Description  sql.NullString
	NameRaw      sql.NullString
	Transport    string
	Specialist   sql.NullString
	Size         string
	Volume       int
	PriceMin     int
	PriceAverage int
	PriceMax     int
	Container    string
	Icon         string
}

func WaresOverview(db *sql.DB, languageId int64, order func() string) ([]*Ware, error) {
	q := query.WaresSelectWaresOverview
	q += order()
	rows, err := db.Query(q, languageId)
	if err != nil {
		return nil, fmt.Errorf("Error querying wares overview: %v", err)
	}
	wares := make([]*Ware, 0)
	for rows.Next() {
		ware := new(Ware)
		err = rows.Scan(&ware.Id, &ware.Name, &ware.NameRaw, &ware.Transport, &ware.PriceAverage, &ware.Icon)
		if err != nil {
			return nil, fmt.Errorf("Error scanning wares overview: %v", err)
		}
		wares = append(wares, ware)
	}
	return wares, nil
}

func GetWare(db *sql.DB, languageId int64, wareId string) (*Ware, error) {
	q := query.WaresSelectWare
	row := db.QueryRow(q, languageId, languageId, wareId)
	ware := new(Ware)
	err := row.Scan(&ware.Id, &ware.Name, &ware.Description, &ware.NameRaw, &ware.Transport, &ware.Specialist, &ware.Size, &ware.Volume, &ware.PriceMin, &ware.PriceAverage, &ware.PriceMax, &ware.Container, &ware.Icon)
	if err != nil {
		return nil, fmt.Errorf("Error scanning ware: %v", err)
	}
	return ware, nil
}
