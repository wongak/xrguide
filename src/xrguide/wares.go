package xrguide

import (
	"database/sql"
	"github.com/codegangsta/martini"
	"html/template"
	"log"
	"net/http"
	"xrguide/content"
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
	w.s.Get("/ware/:id", w.GetWare)
}

func (w *WaresHandler) GetWares(
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

func (w *WaresHandler) GetWare(
	db *sql.DB,
	tmpl *template.Template,
	c *content.XRGuideContent,
	r *http.Request,
	resp http.ResponseWriter,
	l *log.Logger,
	params martini.Params,
) {
	lang, err := contentLanguage(c)
	if err != nil {
		content.HandleError(err, l, tmpl, resp)
		return
	}
	wareId, ok := params["id"]
	if !ok {
		resp.WriteHeader(http.StatusNotFound)
		return
	}
	c.Data["ware"], err = ware.GetWare(db, lang.Id, wareId)
	if err != nil {
		content.HandleError(err, l, tmpl, resp)
		return
	}
	err = tmpl.ExecuteTemplate(resp, "ware.tmpl.html", c)
	if err != nil {
		content.HandleError(err, l, tmpl, resp)
		return
	}
}
