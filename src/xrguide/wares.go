package xrguide

import (
	"database/sql"
	"encoding/json"
	"github.com/codegangsta/martini"
	"html/template"
	"log"
	"net/http"
	"xrguide/content"
	"xrguide/entity"
)

var defaultOverviewOrder = func() string {
	return " ORDER BY name_text.text ASC"
}

type WaresHandler struct {
	s *Server
}

func NewWaresHandler(s *Server) *WaresHandler {
	h := &WaresHandler{s}
	h.registerRoutes()
	return h
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
	wares, err := entity.WaresOverview(db, lang.Id, defaultOverviewOrder)
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
		err = encoder.Encode(wares)
		if err != nil {
			content.HandleHttpError(err, http.StatusInternalServerError, l, resp)
		}
		return
	}
	c.Data["wares"] = &wares
	err = tmpl.ExecuteTemplate(resp, "wares.tmpl.html", c)
	if err != nil {
		content.HandleError(err, l, tmpl, resp)
		return
	}
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
	ware, err := entity.GetWare(db, lang.Id, wareId)
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
		err = encoder.Encode(ware)
		if err != nil {
			content.HandleHttpError(err, http.StatusInternalServerError, l, resp)
		}
		return
	}
	c.Data["ware"] = ware
	err = tmpl.ExecuteTemplate(resp, "ware.tmpl.html", c)
	if err != nil {
		content.HandleError(err, l, tmpl, resp)
		return
	}
}
