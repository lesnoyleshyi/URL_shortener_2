package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Url   string `json:"url,omitempty"`
	Error string `json:"error,omitempty"`
}

func respondWithError(w http.ResponseWriter, respMsg string, status int) {
	//if caller forget to pass respMsg, answer shouldn't be "{}"
	if respMsg == "" {
		respMsg = "unknown error"
	}
	responseStruct := response{Error: respMsg}

	w.WriteHeader(status)
	response, err := json.Marshal(responseStruct)
	if err != nil {
		http.Error(w, "server-side problems", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, "server-side problems", http.StatusInternalServerError)
	}
}

func respondSuccess(w http.ResponseWriter, url string) {
	if url == "" {
		respondWithError(w, "server-side error", http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(response{Url: url})
	if err != nil {
		log.Printf("unable marshall response: %s", err)
		respondWithError(w, "server-side error", http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		log.Printf("unable write response: %s", err)
		respondWithError(w, "server-side problems", http.StatusInternalServerError)
	}
}
