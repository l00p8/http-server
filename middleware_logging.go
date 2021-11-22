package xserver

import (
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"time"

	"github.com/l00p8/log"
)

func WithLogging(log log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			rec := httptest.NewRecorder()
			log.Debug(r.Method + " " + r.URL.String() + " " + r.Header.Get("X-Request-Id"))
			next.ServeHTTP(rec, r)

			dumpResp, _ := httputil.DumpResponse(rec.Result(), true)
			dumpReq, _ := httputil.DumpRequest(r, true)

			// we copy the captured response headers to our new response
			for k, v := range rec.Header() {
				w.Header()[k] = v
			}

			// grab the captured response body
			data := rec.Body.Bytes()
			w.WriteHeader(rec.Result().StatusCode)
			_, _ = w.Write(data)

			dur := time.Since(t1)
			log.Debug("\n" + string(dumpReq) + "\n\n" + string(dumpResp) + " " + dur.String())
		}
		return http.HandlerFunc(fn)
	}
}
