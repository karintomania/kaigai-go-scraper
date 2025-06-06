package httpserver

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
)

type Server struct {
	dbConn     *sql.DB
	pr         *db.PageRepository
	tr         *db.TweetRepository
	push       pushFunc
	httpServer *http.Server
}

func NewServer() *Server {
	dbConn := db.GetDbConnection(common.GetEnv("db_path"))

	return &Server{
		dbConn: dbConn,
		pr:     db.NewPageRepository(dbConn),
		tr:     db.NewTweetRepository(dbConn),
		push: func() (string, error) {
			options := []string{"push"}
			return common.RunGitCommand(options)
		},
		httpServer: &http.Server{},
	}
}

func NewTestServer(
	dbConn *sql.DB,
	dateStr string,
	push pushFunc,
) *Server {
	return &Server{
		dbConn:     dbConn,
		pr:         db.NewPageRepository(dbConn),
		push:       push,
		httpServer: &http.Server{},
	}
}

func (s *Server) Start() {
	defer s.dbConn.Close()

	gph := NewGetPageHandler(s.pr)

	http.HandleFunc("/", gph.getPages)

	ph := NewPublishHandler(s.push, s.pr, s.tr)

	http.HandleFunc("/publish", ph.handle)

	port := fmt.Sprintf(":%s", common.GetEnv("server_port"))
	s.httpServer.Addr = port

	slog.Info("starting serever", "port", port)

	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("error on starting server", "error", err)
	}
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		return err
	}

	slog.Info("Server gracefully stopped")
	return s.dbConn.Close()
}
