package main

import (
	"fmt"
	"io"
	"net/http"

	"cloud.google.com/go/datastore"
)

type API struct {
	datastoreClient *datastore.Client
	config          Config
}

func NewAPI(d *datastore.Client, c Config) *API {
	return &API{
		datastoreClient: d,
		config:          c,
	}
}

// count?Lookup=1
func (a *API) counts(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}
	q := datastore.NewQuery(kindDatastore)
	q = q.Namespace("parzello")

	if lookUp := r.Form.Get("Lookup"); len(lookUp) > 0 {
		q = q.Filter("Lookup =", lookUp)
	}

	if info := r.Form.Get("Info"); len(info) > 0 {
		q = q.Filter("Info =", info)
	}

	n, err := a.datastoreClient.Count(r.Context(), q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, err.Error())
		return
	}
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	w.Header().Set("Expires", "0")
	w.Header().Set("content-type", "application/json")
	fmt.Fprintf(w, "%d", n)
}
