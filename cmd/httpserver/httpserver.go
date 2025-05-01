package httpserver

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/karintomania/kaigai-go-scraper/db"
)

type Server struct {
	dbConn  *sql.DB
	pr      *db.PageRepository
	dateStr string
	publishFun func() error
}

func NewServer() *Server {
	dbConn := db.GetDbConnection(common.GetEnv("db_path"))

	return &Server{
		dbConn:  dbConn,
		pr:      db.NewPageRepository(dbConn),
		dateStr: time.Now().Format("2006-01-02"),
		// TODO: add proper publish function
		publishFun: func() error {
			return nil
		},
	}
}

func NewTestServer(
	dbConn *sql.DB,
	dateStr string,
	publishFun func() error,
) *Server {
	return &Server{
		dbConn:  dbConn,
		pr:      db.NewPageRepository(dbConn),
		dateStr: time.Now().Format("2006-01-02"),
		publishFun: publishFun,
	}
}

func (s *Server) Start() {
	defer s.dbConn.Close()

	gph := NewGetPageHandler(s.pr, s.dateStr)

	http.HandleFunc("/", gph.getPages)

	ph := NewPublishHandler(s.publishFun)

	http.HandleFunc("/publish", ph.handle)

	port := fmt.Sprintf(":%s", common.GetEnv("server_port"))

	slog.Info("starting serever", "port", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		slog.Error("error on starting server", "error", err)
	}
}
