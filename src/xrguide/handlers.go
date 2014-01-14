package xrguide

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini"
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

func handlePage(
	templateName string,
	t *template.Template,
	l *log.Logger,
	c *content.XRGuideContent,
) (respCode int, body string) {
	var buf bytes.Buffer
	respCode = http.StatusOK
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
	c *content.XRGuideContent,
) (int, string) {
	return handlePage("index.tmpl.html", t, l, c)
}

// About page
func about(
	t *template.Template,
	l *log.Logger,
	c *content.XRGuideContent,
) (respCode int, body string) {
	return handlePage("about.tmpl.html", t, l, c)
}

func setLangCookie(langId int64, w http.ResponseWriter) {
	c := &http.Cookie{
		Name:  LANGUAGE_COOKIE_NAME,
		Value: strconv.FormatInt(langId, 10),
	}
	http.SetCookie(w, c)
}

// Language Cookie Middleware
func languageCookie(r *http.Request, w http.ResponseWriter, c *content.XRGuideContent, db *sql.DB, ctx martini.Context) {
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
		setLangCookie(DEFAULT_LANGUAGE, w)
		return
	}
	lang, err := language.LanguageById(db, langId)
	if err != nil {
		log.Printf("Error retrieving language: %v", err)
		setLangCookie(DEFAULT_LANGUAGE, w)
		return
	}
	if lang == nil {
		lang, err = language.LanguageById(db, DEFAULT_LANGUAGE)
		if err != nil {
			log.Printf("Error retrieving default language: %v", err)
			return
		}
		setLangCookie(lang.Id, w)
	}
	c.Data["lang"] = lang
	ctx.Map(c)

	ctx.Next()
}

// Return the set language from the content object
// or return an error
func contentLanguage(c *content.XRGuideContent) (lang *language.Language, err error) {
	lEntry, ok := c.Data["lang"]
	if !ok {
		err = fmt.Errorf("Language not set in content.")
		return
	}
	lang, ok = lEntry.(*language.Language)
	if !ok {
		err = fmt.Errorf("Error on cast language.")
		return
	}
	return
}

func isJsonRequest(r *http.Request) bool {
	return r.Header.Get("Accept") == "application/json"
}
