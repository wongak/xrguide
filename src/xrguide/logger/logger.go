package logger

import (
	"github.com/codegangsta/martini"
	"log"
	"net/http"
	"time"
)

func Logger() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context, log *log.Logger) {
		start := time.Now()
		log.Printf("[%s] Started %s %s", start.Format(time.RFC3339), req.Method, req.URL.Path)

		rw := res.(martini.ResponseWriter)
		c.Next()

		log.Printf("[%s] Completed %v %s in %v\n", time.Now().Format(time.RFC3339), rw.Status(), http.StatusText(rw.Status()), time.Since(start))
	}
}
