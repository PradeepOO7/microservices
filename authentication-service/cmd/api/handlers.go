package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	err := app.readJson(w, r, &requestPayload)
	log.Printf("json data is %v",requestPayload)
	if err != nil {
		log.Panicf("error is %s\n",err)
		app.errorJSON(w, errors.New("Error unmarshalling"))
		return
	}

	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		log.Printf("GetByEmail %s\n",err)
		app.errorJSON(w, errors.New("InvalidCredentials"))
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		log.Printf("PasswordMatches %s\n",err)
		app.errorJSON(w, errors.New("Invalid Password"))
		return
	}

	err=app.logRequest("Authentication",fmt.Sprintf("%s logged in \n",user.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Success",
		Data:    user,
	}

	app.writeJson(w, payload, http.StatusAccepted)

}

func (app *Config)logRequest(name,data string)error{
	var logPayload struct{
		Name string `json:"name"`
		Data string `json:"data"`
	}
	jsonData, _ := json.MarshalIndent(logPayload, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-service:8080/log", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		
		return errors.New("Something went wrong")
	}

	var payload jsonResponse

	err = json.NewDecoder(resp.Body).Decode(&payload)
	if err != nil {
		return errors.New("Something went wrong")
	}

	if payload.Error {
		return errors.New("Something went wrong")
	}

	return nil
}