package xrguide

import (
	"database/sql"
	"html/template"
	"net/http"
)

type WaresHandler struct {
	s *Server
}

func (w *WaresHandler) registerRoutes() {
	w.s.Get("/wares")
}

func (w *WaresHandler) GetWares(db *sql.DB, tmpl *template.Template, r *http.Request, resp http.ResponseWriter) {

}

func NewWaresHandler(s *Server) *WaresHandler {
	h := &WaresHandler{s}
	h.registerRoutes()
	return h
}
