package xserver

import (
	"net/http"

	"github.com/go-chi/chi"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

func NewRouterWithTracing(router Router) Router {
	return &routerWithTracing{router}
}

type routerWithTracing struct {
	router Router
}

func (r *routerWithTracing) Healthers(healthers ...Healther) {
	r.router.Healthers(healthers...)
}

func (r *routerWithTracing) Get(prefix string, fn http.HandlerFunc) {
	r.router.Get(prefix, otelhttp.NewHandler(fn, "GET "+prefix, otelhttp.WithTracerProvider(otel.GetTracerProvider())).ServeHTTP)
}

func (r *routerWithTracing) Post(prefix string, fn http.HandlerFunc) {
	r.router.Post(prefix, otelhttp.NewHandler(fn, "POST "+prefix, otelhttp.WithTracerProvider(otel.GetTracerProvider())).ServeHTTP)
}

func (r *routerWithTracing) Put(prefix string, fn http.HandlerFunc) {
	r.router.Put(prefix, otelhttp.NewHandler(fn, "PUT "+prefix, otelhttp.WithTracerProvider(otel.GetTracerProvider())).ServeHTTP)
}

func (r *routerWithTracing) Patch(prefix string, fn http.HandlerFunc) {
	r.router.Patch(prefix, otelhttp.NewHandler(fn, "PATCH "+prefix, otelhttp.WithTracerProvider(otel.GetTracerProvider())).ServeHTTP)
}

func (r *routerWithTracing) Head(prefix string, fn http.HandlerFunc) {
	r.router.Head(prefix, otelhttp.NewHandler(fn, "HEAD "+prefix, otelhttp.WithTracerProvider(otel.GetTracerProvider())).ServeHTTP)
}

func (r *routerWithTracing) Delete(prefix string, fn http.HandlerFunc) {
	r.router.Delete(prefix, otelhttp.NewHandler(fn, "DELETE "+prefix, otelhttp.WithTracerProvider(otel.GetTracerProvider())).ServeHTTP)
}

func (r *routerWithTracing) Mux() chi.Router {
	return r.router.Mux()
}
