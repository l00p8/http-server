package http_server

import (
	"encoding/json"
	"net/http"
)

func (fac *HttpMiddlewareFactory) JSON(handler Handler) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		resp := handler(r)
		if resp == nil {
			return
		}
		respBody := resp.Response()
		if respBody == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.StatusCode())
			return
		}
		data, err := json.Marshal(respBody)
		if err != nil {
			_, err = w.Write([]byte("{}"))
			if err != nil {
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for k, val := range resp.Headers() {
			w.Header().Set(k, val)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode())
		_, err = w.Write(data)
		if err != nil {
			return
		}
	}
	return fn
}
