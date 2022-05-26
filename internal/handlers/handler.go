package handlers

import (
	"URL_shortener_2/internal/repository"
	"URL_shortener_2/internal/services"
	"encoding/json"
	"errors"
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
		respondWithError(w, "unsupported method", http.StatusMethodNotAllowed)
	}
}

//processShort retrieves shortened URL and writes to w corresponding long URL
//or error if such shortened URL is not in storage
func (h handler) processShort(w http.ResponseWriter, r *http.Request) {
	shortUrl, err := getUrlFromRequest(w, r)
	if err != nil {
		log.Printf("unable retreive shortUrl from request: %s", err)
		return
	}
	//retrieve URL from storage
	longUrl, err := h.service.Get(shortUrl)
	if err != nil {
		if errors.Is(err, repository.ErrNoSuchUrl) {
			respondWithError(w, "no such shortUrl", http.StatusBadRequest)
		} else {
			log.Printf("unable retrieve shortUrl from storage: %s", err)
			respondWithError(w, "server-side error", http.StatusInternalServerError)
		}
		return
	}
	respondSuccess(w, longUrl)
}

//processLong retrieves original URL, saves it in storage
//and returns it's shortened version
func (h handler) processLong(w http.ResponseWriter, r *http.Request) {
	longUrl, err := getUrlFromRequest(w, r)
	if err != nil {
		log.Printf("unable retreive longUrl from request: %s", err)
		return
	}
	//if such URL is in storage, return it's shortened version
	shortUrl, err := h.service.Get(longUrl)
	if err == nil {
		respondSuccess(w, shortUrl)
		return
	}
	//if no such URL is in storage, save it firstly and then return it's shortened version
	if errors.Is(err, repository.ErrNoSuchUrl) {
		shortUrl, err = h.service.Save(longUrl)
		w.WriteHeader(http.StatusCreated)
		respondSuccess(w, shortUrl)
		return
	} else {
		//if some other error occurred, we mustn't save this longUrl
		log.Printf("unable retrieve longUrl from storage: %s", err)
		respondWithError(w, "server-side error", http.StatusInternalServerError)
	}
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
