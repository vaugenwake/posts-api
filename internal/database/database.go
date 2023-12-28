package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/joho/godotenv/autoload"
)

type Post struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type Service interface {
	Health() map[string]string
	CreatePost(title string, content string, slug string) (sql.Result, error)
	GetAllPosts() ([]Post, error)
}

type service struct {
	db *sql.DB
}

var (
	dbname   = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	// Opening a driver typically will not attempt to connect to the database.
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	s := &service{db: db}
	return s
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (s *service) CreatePost(title string, content string, slug string) (sql.Result, error) {
	insertQuery := "INSERT INTO `posts` (title, content, slug, published_at) VALUES (?, ?, ?, NOW())"
	return s.db.ExecContext(context.Background(), insertQuery, title, content, slug)
}

func (s *service) GetAllPosts() ([]Post, error) {
	selectQuert := "SELECT id, title, slug FROM `posts` WHERE published_at < NOW() ORDER BY id DESC"
	result, err := s.db.Query(selectQuert)
	if err != nil {
		return nil, err
	}

	defer result.Close()

	var posts []Post

	for result.Next() {
		var post Post
		if err := result.Scan(&post.ID, &post.Title, &post.Title); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
