package entity

import (
	"database/sql"
	"fmt"
	"xrguide/db/query"
)

type Ware struct {
	Id           string
	Name         sql.NullString  `json:",omitempty"`
	Description  *sql.NullString `json:",omitempty"`
	NameRaw      sql.NullString  `json:",omitempty"`
	Transport    string          `json:",omitempty"`
	Specialist   sql.NullString  `json:",omitempty"`
	Size         string          `json:",omitempty"`
	Volume       int             `json:",omitempty"`
	PriceMin     int             `json:",omitempty"`
	PriceAverage int             `json:",omitempty"`
	PriceMax     int             `json:",omitempty"`
	Container    string          `json:",omitempty"`
	Icon         string          `json:",omitempty"`

	Productions map[string]*Production `json:",omitempty"`
	UsedIn      []*Ware                `json:",omitEmpty"`
}

type Production struct {
	Method string
	Time   int
	Amount int
	Text   sql.NullString

	Wares []*ProductionWare
	Stat  *ProductionStats
}

type ProductionStats struct {
	Method           string
	TotalAverage     int
	ProductionValue  int
	AverageYield     int
	AverageYieldPerH float32
}

type ProductionWare struct {
	Primary bool
	Ware    *Ware
	Amount  int

	AverageCost int
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
	ware.Productions = make(map[string]*Production)
	q = query.WaresSelectProductions
	rows, err := db.Query(q, languageId, wareId)
	if err != nil {
		return nil, fmt.Errorf("Error querying ware production: %v", err)
	}
	productionWares, err := db.Prepare(query.WaresSelectProductionWares)
	if err != nil {
		return nil, fmt.Errorf("Error preparing querying production wares: %v", err)
	}
	for rows.Next() {
		prod := new(Production)
		err = rows.Scan(&prod.Method, &prod.Time, &prod.Amount, &prod.Text)
		if err != nil {
			return nil, fmt.Errorf("Error scanning ware production: %v", err)
		}
		// select production wares
		wares, err := productionWares.Query(languageId, wareId, prod.Method)
		if err != nil {
			return nil, fmt.Errorf("Error querying production wares: %v", err)
		}
		prod.Wares = make([]*ProductionWare, 0)
		for wares.Next() {
			prodWare := new(ProductionWare)
			prodWare.Ware = new(Ware)
			err = wares.Scan(&prodWare.Primary, &prodWare.Ware.Id, &prodWare.Ware.Name, &prodWare.Amount, &prodWare.AverageCost)
			if err != nil {
				return nil, fmt.Errorf("Error scanning production ware: %v", err)
			}
			prod.Wares = append(prod.Wares, prodWare)
		}
		ware.Productions[prod.Method] = prod
	}
	// efficiency
	productionEfficiency, err := db.Query(query.WaresSelectProductionEfficiency, wareId)
	if err != nil {
		return nil, fmt.Errorf("Error querying efficiency: %v", err)
	}
	for productionEfficiency.Next() {
		stat := new(ProductionStats)
		err = productionEfficiency.Scan(&stat.Method, &stat.TotalAverage, &stat.ProductionValue, &stat.AverageYield, &stat.AverageYieldPerH)
		if err != nil {
			return nil, fmt.Errorf("Error scanning production stats: %v", err)
		}
		if _, ok := ware.Productions[stat.Method]; !ok {
			continue
		}
		ware.Productions[stat.Method].Stat = stat
	}
	// used in
	usedIn, err := db.Query(query.WaresSelectUsedIn, languageId, wareId)
	if err != nil {
		return nil, fmt.Errorf("Error querying used in: %v", err)
	}
	ware.UsedIn = make([]*Ware, 0)
	for usedIn.Next() {
		used := new(Ware)
		err = usedIn.Scan(&used.Id, &used.Name, &used.NameRaw)
		if err != nil {
			return nil, fmt.Errorf("Error scanning used in: %v", err)
		}
		ware.UsedIn = append(ware.UsedIn, used)
	}
	return ware, nil
}
