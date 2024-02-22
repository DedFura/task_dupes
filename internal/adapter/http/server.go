package httpserver

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"task/internal/usecase"
	"time"

	"github.com/sirupsen/logrus"
)

type ServerInterface interface {
	Start(port string) error
}

type Server struct {
	service usecase.ConnectionChecker
	mux     *http.ServeMux
	log     *logrus.Logger
}

func NewServer(service usecase.ConnectionChecker, logger *logrus.Logger) ServerInterface {
	server := &Server{
		service: service,
		mux:     http.NewServeMux(),
		log:     logger,
	}
	server.routes()

	return server
}

func (s *Server) routes() {
	s.mux.HandleFunc("/", s.handleCheck)
}

func (s *Server) Start(port string) error {
	return http.ListenAndServe(port, s.mux)
}

func (s *Server) handleCheck(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	segments := strings.Split(r.URL.Path, "/")
	if len(segments) < 3 {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user1, err := strconv.ParseInt(segments[len(segments)-2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user1 ID", http.StatusBadRequest)
		return
	}

	user2, err := strconv.ParseInt(segments[len(segments)-1], 10, 64)
	if err != nil {
		http.Error(w, "Invalid user2 ID", http.StatusBadRequest)
		return
	}

	dupes, err := s.service.CheckDupes(user1, user2)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf(`{"dupes": %t}`, dupes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))

	elapsed := time.Since(start)
	s.log.Infof("HandleCheck took %s", elapsed)
}
