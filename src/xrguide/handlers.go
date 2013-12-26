package xrguide

import (
	"bytes"
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"xrguide/content"
	"xrguide/entity/language"
)

const (
	LANGUAGE_COOKIE_NAME       = "lang"
	DEFAULT_LANGUAGE     int64 = 44
)

func setHandlers(s *Server) {
	s.Map(content.NewContent())
	s.Use(languageCookie)
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

func setLangCookie(langId int64, w http.ResponseWriter) {
	c := &http.Cookie{
		Name:  LANGUAGE_COOKIE_NAME,
		Value: strconv.FormatInt(langId, 10),
	}
	http.SetCookie(w, c)
}

// Language Cookie Middleware
func languageCookie(r *http.Request, w http.ResponseWriter, c *content.XRGuideContent, db *sql.DB) {
	langCookie, err := r.Cookie(LANGUAGE_COOKIE_NAME)
	if err != nil {
		if err == http.ErrNoCookie {
			setLangCookie(DEFAULT_LANGUAGE, w)
			return
		}
		log.Printf("Error getting cookie: %v", err)
		return
	}
	langId, err := strconv.ParseInt(langCookie.Value, 10, 64)
	if err != nil {
		log.Printf("Invalid value in cookie: %v", err)
		return
	}
	lang, err := language.LanguageById(db, langId)
	if err != nil {
		log.Printf("Error retrieving language: %v", err)
		return
	}
	if lang == nil {
		setLangCookie(DEFAULT_LANGUAGE, w)
		return
	}
	setLangCookie(lang.Id, w)
}
