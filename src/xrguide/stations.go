package xrguide

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"xrguide/content"
	"xrguide/entity"
)

type StationsHandler struct {
	s *Server
}

func NewStationsHandler(s *Server) *StationsHandler {
	h := &StationsHandler{s}
	h.registerRoutes()
	return h
}

func (s *StationsHandler) registerRoutes() {
	s.s.Get("/stations", s.GetStations)
}

func (s *StationsHandler) GetStations(
	db *sql.DB,
	tmpl *template.Template,
	c *content.XRGuideContent,
	r *http.Request,
	resp http.ResponseWriter,
	l *log.Logger,
) {
	lang, err := contentLanguage(c)
	if err != nil {
		content.HandleError(err, l, tmpl, resp)
		return
	}
	stations, err := entity.StationsOverview(db, lang.Id)
	if err != nil {
		if isJsonRequest(r) {
			content.HandleHttpError(err, http.StatusInternalServerError, l, resp)
		} else {
			content.HandleError(err, l, tmpl, resp)
		}
		return
	}
	if isJsonRequest(r) {
		resp.Header().Add("Content-Type", "application/json")
		encoder := json.NewEncoder(resp)
		err = encoder.Encode(stations)
		if err != nil {
			content.HandleHttpError(err, http.StatusInternalServerError, l, resp)
		}
		return
	}
	err = tmpl.ExecuteTemplate(resp, "stations.tmpl.html", c)
	if err != nil {
		content.HandleError(err, l, tmpl, resp)
		return
	}
}
