package xserver

import (
	"encoding/json"
	"net/http"
)

type Healther interface {
	Health() error
}

func healthHandler(healthers ...Healther) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var errs []string
		for _, h := range healthers {
			if err := h.Health(); err != nil {
				errs = append(errs, err.Error())
			}
		}
		if len(errs) > 0 {
			d, err := json.Marshal(errs)
			if err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
			w.WriteHeader(http.StatusServiceUnavailable)
			_, err = w.Write(d)
		}
	}
}
