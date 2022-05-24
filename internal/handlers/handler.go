package handlers

import (
	"URL_shortener_2/internal/services"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Handler interface {
	http.Handler
}

type handler struct {
	service services.Service
}

type request struct {
	Url string `json:"url,omitempty"`
}

func New(s services.Service) Handler {
	return &handler{service: s}

}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.processLong(w, r)
	case http.MethodGet:
		h.processShort(w, r)
	default:
		http.Error(w, "unsupported method", http.StatusMethodNotAllowed)
	}
}

//processShort retrieves shortened URL and writes to w corresponding long URL
//or error if such shortened URL is not in storage
func (h handler) processShort(w http.ResponseWriter, r *http.Request) {
	url, err := getUrlFromRequest(w, r)
	if err != nil {
		log.Printf("unable retreive url from request: %s", err)
		return
	}
	//retrieve short URL from storage
	urlStruct, err := h.service.Get(url)
	if err != nil {
		if err.Error() == "no such url is in storage" {
			respondWithError(w, "no such url", http.StatusBadRequest)
		} else {
			log.Printf("unable retrieve url from storage: %s", err)
			respondWithError(w, "server-side error", http.StatusInternalServerError)
		}
		return
	}
	respondSuccess(w, urlStruct.LongUrl)
}

//processLong retrieves original URL, saves it in storage
//and returns it's shortened version
func (h handler) processLong(w http.ResponseWriter, r *http.Request) {
	url, err := getUrlFromRequest(w, r)
	if err != nil {
		log.Printf("unable retreive url from request: %s", err)
		return
	}
	//if such URL is in storage, return it's shortened version
	url1, err := h.service.Get(url)
	if url1 != nil {
		respondSuccess(w, url1.ShortUrl)
		return
	}
	//if no such URL is in storage, save it firstly and then return it's shortened version
	shortUrl, err := h.service.Save(url)
	respondSuccess(w, shortUrl)
}

//getUrlFromRequest reads body from Request, marshall it to request struct
//and returns found URL or error
func getUrlFromRequest(w http.ResponseWriter, r *http.Request) (string, error) {
	var req request

	bodyCont, err := ioutil.ReadAll(r.Body)
	defer func() { _ = r.Body.Close() }()
	if err != nil {
		respondWithError(w, "unable to read body", http.StatusInternalServerError)
		return "", fmt.Errorf("unable to read body: %s", err)
	}
	if len(bodyCont) == 0 {
		respondWithError(w, "empty body", http.StatusBadRequest)
		return "", fmt.Errorf("empty body")
	}
	if err := json.Unmarshal(bodyCont, &req); err != nil {
		respondWithError(w, "unable recognize data in body", http.StatusBadRequest)
		return "", fmt.Errorf("unable recognize data in body: %s", err)
	}
	if req.Url == "" {
		respondWithError(w, "incorrect data in body", http.StatusBadRequest)
		return "", fmt.Errorf("incorrect data in body")
	}
	return req.Url, nil
}
