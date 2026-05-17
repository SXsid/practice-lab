package main

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

// we we want to update the hash/json vlaue se hset hmmm.....
func (app *App) SetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	defer r.Body.Close()
	res, err := app.redis.Set(r.Context(), id, string(data), time.Minute*5).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func (app *App) GetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := app.redis.Get(r.Context(), id).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// just bytes so no marshal but in hget we have map so we need to marshal
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(res))
}

func (app *App) HSetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	defer r.Body.Close()
	_, err = app.redis.HSet(r.Context(), id, req).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("ok"))
}

func (app *App) HGetUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	res, err := app.redis.HGetAll(r.Context(), id).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "applicaton/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}
