package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/laetificat/slogger/pkg/slogger"
)

// LogMiddleWare is the middleware for the http routers to log the requests.
type LogMiddleWare struct {
	next http.Handler
}

/*
NewLogMiddleWare returns a new LogMiddleWare struct.
*/
func NewLogMiddleWare(next http.Handler) *LogMiddleWare {
	return &LogMiddleWare{next: next}
}

/*
ServeHTTP wraps the ServeHTTP and adds a logger.
*/
func (m *LogMiddleWare) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	defer slogger.Info(fmt.Sprintf("Request: %s %s, time taken: %s", r.Method, r.URL.Path, time.Since(t)))

	m.next.ServeHTTP(w, r)
}
