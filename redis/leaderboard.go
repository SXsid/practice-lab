package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *App) LeaderBoradRouter() {
	app.Router.Get("/post/{id}/view", app.ViewPostById)
	app.Router.Post("/leaderBoard/score", app.UpdateScoreinLeaderBorad)
	app.Router.Get("/leaderBoard/{room_id}", app.GetTopPerformerFromLeaderBoard)
	app.Router.Get("/leaderBoard/{user_id}/rank", app.GetUserRankFromLeaderBoard)
}

// can easily repalced by atomic updae in postgress
func (app *App) ViewPostById(w http.ResponseWriter, r *http.Request) {
	count, err := app.redis.Incr(r.Context(), fmt.Sprintf("post:%s", chi.URLParam(r, "id"))).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%d", count)
}

func (app *App) UpdateScoreinLeaderBorad(w http.ResponseWriter, r *http.Request) {
	var data UpdatScoreRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := app.redis.ZIncrBy(r.Context(), data.RoomID, float64(data.Score), data.UserID).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	leaderBoard, err := app.GetLeaderBoard(r.Context(), data.RoomID)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(&leaderBoard)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (app *App) GetLeaderBoard(ctx context.Context, key string) ([]LeaderBoardEntry, error) {
	data := make([]LeaderBoardEntry, 0)

	val, err := app.redis.ZRevRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	for _, v := range val {
		data = append(data, LeaderBoardEntry{
			UserID: v.Member.(string),
			Score:  v.Score,
		})
	}
	return data, nil
}

func (app *App) GetTopPerformerFromLeaderBoard(w http.ResponseWriter, r *http.Request) {
	leaderBoard, err := app.GetLeaderBoard(r.Context(), chi.URLParam(r, "room_id"))
	data := GetLeaderBoardResp{
		Data: leaderBoard,
	}
	res, err := json.Marshal(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func (app *App) GetUserRankFromLeaderBoard(w http.ResponseWriter, r *http.Request) {
}
