package content

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type XRGuideContent struct {
	Data map[string]interface{}
}

func NewContent() *XRGuideContent {
	m := make(map[string]interface{})
	c := &XRGuideContent{
		Data: m,
	}
	c.Data["title"] = "X Rebirth Guide"
	return c
}

func HandleError(err error, log *log.Logger, t *template.Template, w http.ResponseWriter) {
	c := NewContent()
	c.Data["timestamp"] = time.Now().UTC().UnixNano()
	log.Printf("Error [%d]", c.Data["timestamp"])
	log.Printf("Errorinfo: %v", err)
	tmplErr := t.ExecuteTemplate(w, "error.tmpl.html", c)
	if tmplErr != nil {
		log.Printf("Template error on error page :/ (%v)", tmplErr)
	}
}

func HandleHttpError(err error, status int, log *log.Logger, w http.ResponseWriter) {
	ts := time.Now().UTC().UnixNano()
	log.Printf("Error [%d]", ts)
	log.Printf("Errorinfo: %v", err)
	w.WriteHeader(status)
	w.Write([]byte(strconv.FormatInt(ts, 10)))
}
