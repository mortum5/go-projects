package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/mortum5/go-projects/dgraph-go/model"
)

type Repository interface {
	Get(string) (model.Url, bool)
	Set(model.Url)
}

type Server struct {
	srv     *http.Server
	r       Repository
	errChan chan error
}

func New(r Repository) *Server {
	return &Server{
		r:       r,
		errChan: make(chan error, 1),
	}
}

func (s *Server) Run() {
	router := http.NewServeMux()
	router.HandleFunc("GET /urls/", s.getHandler())
	router.HandleFunc("POST /urls/", s.postHandler())

	srv := &http.Server{
		Addr:              ":9090",
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      5 * time.Second,
		MaxHeaderBytes:    1 << 20,
		Handler:           router,
	}

	s.srv = srv

	go func() {
		log.Println("Server started on", srv.Addr)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.errChan <- fmt.Errorf("server: %v", err)
		}
	}()

}

func (s *Server) getHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimPrefix(r.URL.Path, "/urls/")

		url, ok := s.r.Get(key)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		json, _ := json.Marshal(url)

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}

func (s *Server) postHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var params map[string]string = make(map[string]string)
		err := json.NewDecoder(r.Body).Decode(&params)
		if err != nil {
			log.Printf("post handler: %v", err)
		}
		defer r.Body.Close()

		url := model.New(params["url"], params["slug"], time.Hour)
		s.r.Set(url)
		log.Println("create new url", "url", params["url"], "slug", params["slug"])

		w.Write([]byte("ok"))
	}
}

func (s *Server) Stop() error {

	if err := s.srv.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("server: stops %v", err)
	}
	log.Println("Server stops")
	return nil
}

func (s *Server) Error() chan error {
	return s.errChan
}
