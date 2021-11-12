package http_server

import (
	hh "github.com/l00p8/utils"
	"net/http"
)

type Handler func(r *http.Request) hh.Response

type Middleware func(handler Handler) Handler

type OutputMiddleware func(handler Handler) http.HandlerFunc
