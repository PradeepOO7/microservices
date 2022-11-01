package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (app *Config) readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	
	if err != nil {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

func (app *Config) writeJson(w http.ResponseWriter, data any, status int, headers ...http.Header) error {
	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusBad := http.StatusBadRequest
	if len(status) > 0 {
		statusBad = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()
	payload.Data = nil
	return app.writeJson(w, payload, statusBad)
}
