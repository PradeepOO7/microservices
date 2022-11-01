package main

import (
	"broker/events"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type MailPayload struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type RequestPayload struct {
	Action string       `json:"action"`
	Auth   Authenticate `json:"auth,omitempty"`
	Log    LogEntry     `json:"log,omitempty"`
	Mail   MailPayload  `json:"mail,omitempty"`
}	

type Authenticate struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogEntry struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {

	res := jsonResponse{
		Error:   false,
		Message: "Broker service is working",
	}

	app.writeJson(w, res, http.StatusOK)

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch requestPayload.Action {
	case "auth":
		app.Authenticate(w, requestPayload.Auth)
	case "log":
		//app.logItem(w, requestPayload.Log)
		app.logEventViaRabbitMQ(w,requestPayload.Log)
	case "mail":
		app.SendMail(w,requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New(" Action not allowed"))
	}

}

func (app *Config) Authenticate(w http.ResponseWriter, authPayload Authenticate) {
	jsonData, _ := json.MarshalIndent(authPayload, "", "\t")

	request, err := http.NewRequest("POST", "http://authentication-service:8080/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("Invalid Credentials"))
		return
	} else if resp.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("Invalid Credentials"))
		return
	}

	var payload jsonResponse

	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if payload.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	result := jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    payload.Data,
	}

	app.writeJson(w, result, 202)

}

func (app *Config) logItem(w http.ResponseWriter, logPayload LogEntry) {
	jsonData, _ := json.MarshalIndent(logPayload, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-service:8080/log", bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("Something went wrong"))
		return
	}

	var payload jsonResponse

	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if payload.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	result := jsonResponse{
		Error:   false,
		Message: "logged",
		Data:    payload.Data,
	}

	app.writeJson(w, result, http.StatusAccepted)

}

func (app *Config) SendMail(w http.ResponseWriter, msg MailPayload) {
	jsonData, _ := json.MarshalIndent(msg, "", "\t")

	request, err := http.NewRequest("POST", "http://mailer-service:8080/mail", bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("Something went wrong"))
		return
	}

	var payload jsonResponse

	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if payload.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	result := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Sent mail to %s \n",msg.To),
		Data:    payload.Data,
	}

	app.writeJson(w, result, http.StatusAccepted)

}

func (app *Config) logEventViaRabbitMQ(w http.ResponseWriter,l LogEntry){
	err:=app.pushToQueue(l)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	result := jsonResponse{
		Error:   false,
		Message: "logged via RabbitMQ",
		Data:    l,
	}

	app.writeJson(w, result, http.StatusAccepted)

}

func (app *Config)pushToQueue(l LogEntry)error{
	emitter,err:= events.NewEventEmitter(app.Rabbit)
	if err!=nil{
		return err
	}

	jsonData, _ := json.MarshalIndent(l, "", "\t")
	err = emitter.Push(string(jsonData),"log.INFO")
	if err!=nil{
		return err
	}
	return nil
}