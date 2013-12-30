package xrguide

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"xrguide/content"
	"xrguide/entity/language"
	"xrguide/entity/ware"
)

var defaultOverviewOrder = func() string {
	return " ORDER BY name_text.text ASC"
}

type WaresHandler struct {
	s *Server
}

func (w *WaresHandler) registerRoutes() {
	w.s.Get("/wares", w.GetWares)
}

func (w *WaresHandler) GetWares(
	db *sql.DB,
	tmpl *template.Template,
	c *content.XRGuideContent,
	r *http.Request,
	resp http.ResponseWriter,
	l *log.Logger,
) {
	var err error
	lEntry, ok := c.Data["lang"]
	if !ok {
		err = fmt.Errorf("Language not set in content.")
		content.HandleError(err, l, tmpl, resp)
		return
	}
	lang, ok := lEntry.(*language.Language)
	if !ok {
		err = fmt.Errorf("Error on cast language.")
		content.HandleError(err, l, tmpl, resp)
		return
	}
	c.Data["wares"], err = ware.WaresOverview(db, lang.Id, defaultOverviewOrder)
	if err != nil {
		content.HandleError(err, l, tmpl, resp)
		return
	}
	err = tmpl.ExecuteTemplate(resp, "wares.tmpl.html", c)
	if err != nil {
		content.HandleError(err, l, tmpl, resp)
		return
	}
}

func NewWaresHandler(s *Server) *WaresHandler {
	h := &WaresHandler{s}
	h.registerRoutes()
	return h
}
