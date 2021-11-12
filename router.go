package http_server

import (
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
)

type Muxer interface {
	Mux() chi.Router
}

type Router interface {
	Healthers(healthers ...Healther)

	Get(prefix string, fn http.HandlerFunc)

	Post(prefix string, fn http.HandlerFunc)

	Put(prefix string, fn http.HandlerFunc)

	Delete(prefix string, fn http.HandlerFunc)

	Patch(prefix string, fn http.HandlerFunc)

	Head(prefix string, fn http.HandlerFunc)

	Muxer
}

type router struct {
	mux    chi.Router
	Config Config
	mw     OutputMiddleware
}

func (r *router) Healthers(healthers ...Healther) {
	r.mux.Get("/_health", healthHandler(healthers...))
}

func (r *router) Get(prefix string, fn http.HandlerFunc) {
	r.mux.Get(prefix, fn)
}

func (r *router) Post(prefix string, fn http.HandlerFunc) {
	r.mux.Post(prefix, fn)
}

func (r *router) Put(prefix string, fn http.HandlerFunc) {
	r.mux.Put(prefix, fn)
}

func (r *router) Patch(prefix string, fn http.HandlerFunc) {
	r.mux.Patch(prefix, fn)
}

func (r *router) Head(prefix string, fn http.HandlerFunc) {
	r.mux.Head(prefix, fn)
}

func (r *router) Delete(prefix string, fn http.HandlerFunc) {
	r.mux.Delete(prefix, fn)
}

func (r *router) Mux() chi.Router {
	return r.mux
}

func NewRouter(cfg Config) Router {
	r := &router{mux: chi.NewRouter(), Config: cfg}
	lmt := tollbooth.NewLimiter(float64(cfg.RateLimit), nil)

	timeout := r.Config.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	//r.Mux.Use(WithLogging(cfg.Logger.Logger.Desugar()))
	r.mux.Use(chiMiddleware.Logger)
	r.mux.Use(chiMiddleware.RequestID)
	r.mux.Use(chiMiddleware.StripSlashes)
	r.mux.Use(chiMiddleware.Recoverer)
	r.mux.Use(chiMiddleware.Timeout(timeout))
	r.mux.Use(prometheusMiddleware)
	r.mux.Use(rateLimitter(lmt))

	return r
}
