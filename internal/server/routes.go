package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)

	r.Get("/api/v1/posts", s.HandlerGetPosts)
	r.Post("/api/v1/posts", s.HandlerCreatePost)

	return r
}

func internalError(err error, w http.ResponseWriter) {
	if err != nil {
		log.Fatalf("Could not get posts from DB, Err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (s *Server) HandlerGetPosts(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]interface{})

	posts, err := s.db.GetAllPosts()
	if err != nil {
		log.Fatalf("Could not get posts from DB, Err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp["posts"] = posts
	resp["total_posts"] = len(posts)

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Could not marshal response, Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) HandlerCreatePost(w http.ResponseWriter, r *http.Request) {
	type PostParams struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Slug    string `json:"slug"`
	}

	decoder := json.NewDecoder(r.Body)

	params := PostParams{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Fatalf("Could not decode params. Err %v", err)
	}

	errs := make(map[string][]string)

	if params.Title == "" {
		errs["title"] = append(errs["title"], "Title is required")
	}

	if params.Content == "" {
		errs["content"] = append(errs["content"], "Content is required")
	}

	if params.Slug == "" {
		errs["slug"] = append(errs["slug"], "Slug is required")
	}

	if len(errs) > 0 {
		fmt.Println(errs)
		jsonResp, err := json.Marshal(errs)
		if err != nil {
			log.Fatalf("Could not marshal errors, Err: %v", err)
		}

		_, _ = w.Write(jsonResp)
		return
	}

	result, err := s.db.CreatePost(params.Title, params.Content, params.Slug)
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf("Could not insert: %w", err)))
		return
	}

	insertedId, _ := result.LastInsertId()
	params.ID = int(insertedId)

	jsonResp, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("Could not marshal response, Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World 2"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
