package xrguide

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"xrguide/content"
)

func setHandlers(s *Server) {
	s.Get("/", index)
	s.Get("/about", about)
}

func handlePage(templateName string, t *template.Template, l *log.Logger) (respCode int, body string) {
	var buf bytes.Buffer
	respCode = http.StatusOK
	c := content.NewContent()
	err := t.ExecuteTemplate(&buf, templateName, c)
	if err != nil {
		log.Printf("Template error: %v", err)
		respCode = http.StatusInternalServerError
		return
	}
	body = buf.String()
	return
}

// The Index Page
func index(
	t *template.Template,
	l *log.Logger,
) (int, string) {
	return handlePage("index.tmpl.html", t, l)
}

// About page
func about(
	t *template.Template,
	l *log.Logger,
) (respCode int, body string) {
	return handlePage("about.tmpl.html", t, l)
}
