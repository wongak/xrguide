// The main X Rebirth Guide server
package xrguide

import (
	"database/sql"
	"fmt"
	"github.com/codegangsta/martini"
	"html/template"
	"net/http"
	"os"
	"path"
)

type Server struct {
	*martini.Martini
	martini.Router

	htmlSrcDir string
}

func NewServer(db *sql.DB, htmlSrcDir string) (*Server, error) {
	// martini initialization
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Action(r.Handle)

	// default mapping
	m.Map(db)

	s := &Server{
		Martini: m,
		Router:  r,

		htmlSrcDir: htmlSrcDir,
	}
	err := s.initDefaults()
	if err != nil {
		return nil, fmt.Errorf("Error initializing server: %v", err)
	}
	setHandlers(s)
	return s, nil
}

func (s *Server) initDefaults() error {
	// template
	pattern := path.Join(s.htmlSrcDir, "*", "*.tmpl.html")
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("Error on parsing templates: %v", err)
	}
	s.Map(tmpl)

	// serving static directory
	staticPath := path.Join(s.htmlSrcDir, "static")
	info, err := os.Stat(staticPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("Static directory %s does not exist.", staticPath)
		}
		return fmt.Errorf("Error on stat static directory %s: %v", staticPath, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("Static path %s is not a directory.", staticPath)
	}
	s.Use(martini.Static(staticPath))
	return nil
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s)
}
