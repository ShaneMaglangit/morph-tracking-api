package router

import (
	"encoding/json"
	"net/http"
)

// MorphHandler returns a list of the most recently morphed Axies
func (deps *Deps) MorphHandler(w http.ResponseWriter, r *http.Request) {
	page := getIntParams(r.URL.Query().Get("page"), 0)
	asc := getBoolParams(r.URL.Query().Get("asc"), false)
	byId := getBoolParams(r.URL.Query().Get("byId"), false)
	axies := deps.db.SelectAxies(page, asc, byId)
	//Return results
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(axies)
}
