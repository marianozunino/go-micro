package main

import (
	"encoding/json"
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	var response jsonResponse
	response.OK = true
	response.Message = "Broker service is up and running"
	response.Data = nil
	out, err := json.MarshalIndent(response, "", "\t")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)

}
