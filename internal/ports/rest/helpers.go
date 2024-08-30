package rest

import (
	"encoding/json"
	"github.com/Kugeki/kode_test_task/internal/ports/rest/dto"
	"log/slog"
	"net/http"
)

func respond(w http.ResponseWriter, code int, data interface{}, log *slog.Logger) {
	var jsonData []byte
	var err error

	if data == nil {
		log.Warn("respond helper: data is nil")
	}

	jsonData, err = json.Marshal(data)

	if err != nil {
		log.Error("respond helper: json marshal error", slog.Any("error", err))

		jsonData = []byte("{\"error\": \"json marshal error\"}")
		code = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(code)
	_, err = w.Write(jsonData)
	if err != nil {
		log.Error("respond helper: response write json data error", slog.Any("error", err))
	}
}

func respondError(w http.ResponseWriter, code int, err error, log *slog.Logger) {
	respond(w, code, &dto.HTTPError{Error: err.Error()}, log)
}

func respondUnauthorized(w http.ResponseWriter, err error, log *slog.Logger) {
	w.Header().Set("WWW-Authenticate", `Bearer realm="api"`)
	respondError(w, http.StatusUnauthorized, err, log)
}
