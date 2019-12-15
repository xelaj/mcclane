package api

import (
	"encoding/json"
	"github.com/xelaj/mcclane/internal/db/pg"
	"github.com/xelaj/mcclane/internal/model"
	"net/http"
)

func AddHotLocation(db *pg.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := new(model.HotLocation)
		err := json.NewDecoder(r.Body).Decode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = db.AddHotLocation(r.Context(), *res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func AddNews(db *pg.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := new(model.News)
		err := json.NewDecoder(r.Body).Decode(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = db.AddNews(r.Context(), *res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
