package language

import (
	"database/sql"
	"time"
	"xrguide/db/query"
)

type Language struct {
	Id   int64
	Name string
}

var (
	langIdCache   = make(map[int64]*Language)
	langNameCache = make(map[string]*Language)
	setCache      = make(chan *Language)
	getIdCache    = make(chan map[int64]*Language)
	getNameCache  = make(chan map[string]*Language)
)

func init() {
	go func() {
		select {
		case l := <-setCache:
			langIdCache[l.Id] = l
			langNameCache[l.Name] = l

		case getIdCache <- langIdCache:

		case getNameCache <- langNameCache:
		}
	}()
}

func LanguageById(db *sql.DB, id int64) (*Language, error) {
	var lang *Language
	to := time.After(50 * time.Millisecond)
	select {
	case cache := <-getIdCache:
		lang, ok := cache[id]
		if ok {
			return lang, nil
		}
	case <-to:
	}
	lang = new(Language)
	q := query.SelectLanguage + " WHERE id = ?"
	row := db.QueryRow(q, id)
	err := row.Scan(&lang.Id, &lang.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	to = time.After(50 * time.Millisecond)
	select {
	case setCache <- lang:
	case <-to:
	}
	return lang, nil
}

func LanguageByName(db *sql.DB, name string) (*Language, error) {
	var lang *Language
	to := time.After(50 * time.Millisecond)
	select {
	case cache := <-getNameCache:
		lang, ok := cache[name]
		if ok {
			return lang, nil
		}
	case <-to:
	}
	lang = new(Language)
	q := query.SelectLanguage + " WHERE name = ?"
	row := db.QueryRow(q, name)
	err := row.Scan(&lang.Id, &lang.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	to = time.After(50 * time.Millisecond)
	select {
	case setCache <- lang:
	case <-to:
	}
	return lang, nil

}
