package main

import (
	"logger/data"
	"net/http"
)

type Request struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var requestPayload Request
	err := app.readJson(w, r, &requestPayload)

	if err != nil {
		app.errorJSON(w, err)
		return
	}
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
	err = event.Insert(event)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: "logged",
		Data:    requestPayload,
	}

	app.writeJson(w, response, 202)

}
