package main

import (
	"log"
	"net/http"
)

func (app *app) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("internal server error: ", err)
	writeJSON(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *app) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	log.Printf("forbidden request: IP: %s Method: %s Path:%s", r.RemoteAddr, r.Method, r.URL.Path)
	writeJSON(w, http.StatusForbidden, "forbidden")
}

func (app *app) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("Bad request: ", err)
	writeJSON(w, http.StatusBadRequest, err.Error())
}

func (app *app) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Println("error not found: ", err)
	writeJSON(w, http.StatusNotFound, "not found")
}
