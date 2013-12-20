package content

import (
	"html/template"
	"log"
	"net/http"
	"time"
)

type XRModDirContent struct {
	Data map[string]interface{}
}

func NewContent() *XRModDirContent {
	m := make(map[string]interface{})
	c := &XRModDirContent{
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
