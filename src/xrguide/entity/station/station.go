package station

import (
	"database/sql"
	"fmt"
	"xrguide/db/query"
	"xrguide/entity/ware"
)

type Station struct {
	Id   string
	Name sql.NullString

	ProducedWares []*ware.Ware
}

func (s *Station) addWare(w *ware.Ware) {
	if s.ProducedWares == nil {
		s.ProducedWares = make([]*ware.Ware, 0)
	}
	s.ProducedWares = append(s.ProducedWares, w)
}

func StationsOverview(db *sql.DB, langId int64) ([]*Station, error) {
	rows, err := db.Query(query.SelectStations, langId, langId)
	if err != nil {
		return nil, fmt.Errorf("Error on query select stations: %v", err)
	}
	// a little awkward, but since we receive multiple rows
	// for each station (every production in one row),
	// we first map by station id
	stations := make(map[string]*Station, 0)
	for rows.Next() {
		var wareId string
		var wareName sql.NullString
		station := new(Station)
		err = rows.Scan(&station.Id, &station.Name, &wareId, &wareName)
		s, ok := stations[station.Id]
		if ok {
			station = s
		} else {
			stations[station.Id] = station
		}
		ware := &ware.Ware{Id: wareId, Name: wareName}
		station.addWare(ware)
	}
	// then we create a slice out of the map
	i := 0
	ret := make([]*Station, len(stations))
	for _, station := range stations {
		ret[i] = station
		i++
	}
	return ret, nil
}
