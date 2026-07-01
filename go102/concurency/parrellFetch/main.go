package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

var ErrTimeout = errors.New("times up")

type (
	UserService struct{}
	Profile     struct {
		Name string `json:"name"`
	}
)

func (u *UserService) fetchProfile(ctx context.Context, id string, ch chan<- Response) {
	select {
	case <-time.After(time.Millisecond * 300):
		ch <- Response{UserReponse: UserReponse{Profile: Profile{Name: "sid"}}, service: "profile"}
	case <-ctx.Done():
		ch <- Response{err: ErrTimeout}
	}
}

type (
	NotificationService struct{}
	Notification        struct{}
)

func (n *NotificationService) fetchNotification(ctx context.Context, id string, ch chan<- Response) {
	select {
	case <-time.After(time.Millisecond * 300):
		ch <- Response{service: "notification"}
	case <-ctx.Done():
		ch <- Response{err: ErrTimeout}
	}
}

type (
	PostsServices struct{}
	Post          struct{}
)

func (p *PostsServices) fetchPosts(ctx context.Context, id string, ch chan<- Response) {
	select {
	case <-time.After(time.Millisecond * 300):
		ch <- Response{service: "post"}
	case <-ctx.Done():
		ch <- Response{err: ErrTimeout}
	}
}

type UserHandler struct {
	UserService         *UserService
	NotificationService *NotificationService
	PostsServices       *PostsServices
}

type Response struct {
	UserReponse
	err     error
	service string
}

type UserReponse struct {
	Profile       Profile        `json:"profile"`
	Posts         []Post         `json:"posts"`
	Notifications []Notification `json:"notification"`
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	ch := make(chan Response, 3)
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*300)
	defer cancel()
	go func(h *UserHandler) {
		var wg sync.WaitGroup

		wg.Go(func() {
			h.NotificationService.fetchNotification(ctx, userID, ch)
		})
		wg.Go(func() {
			h.PostsServices.fetchPosts(ctx, userID, ch)
		})
		wg.Go(func() {
			h.UserService.fetchProfile(ctx, userID, ch)
		})
		wg.Wait()
		close(ch)
	}(h)

	result := make([]Response, 0, 3)
	for range 3 {
		result = append(result, <-ch)
	}
	var res UserReponse
	for _, data := range result {
		if data.err != nil && !errors.Is(data.err, ErrTimeout) {
			http.Error(w, data.err.Error(), http.StatusInternalServerError)
			return
		}
		switch data.service {
		case "post":
			res.Posts = data.UserReponse.Posts
		case "profile":
			res.Profile = data.UserReponse.Profile
		case "notification":
			res.Notifications = data.UserReponse.Notifications
		default:
			http.Error(w, "ivalid meethod", http.StatusInternalServerError)
			return

		}
	}
	byte, _ := json.Marshal(&res)
	w.Write(byte)
}

func main() {
	u := &UserHandler{}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /user/{id}", u.getUser)
	req := httptest.NewRequest(http.MethodGet, "/user/123", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	fmt.Println(w.Body.String())
	fmt.Println(w.Code)
}
