package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"net/http"
	"time"

	"github.com/ashtanko/octo-server/store"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	ctxKeyRequestID         = iota
	checkViolationErrorCode = "23514"
)

type Server struct {
	// RootRouter is the starting point for all HTTP requests to the server.
	RootRouter      *mux.Router
	Server          *http.Server
	Logger          *logrus.Logger
	Store           store.Store
	exitSignal      chan struct{}
	didFinishListen chan struct{}
	goroutineCount  int32
}

func NewServer(store store.Store) (*Server, error) {

	rootRouter := mux.NewRouter()
	logger := logrus.New()

	s := &Server{
		exitSignal: make(chan struct{}, 1),
		RootRouter: rootRouter,
		Logger:     logger,
		Store:      store,
	}

	s.configureRouter()
	s.InitAccount()
	s.InitTransaction()

	return s, nil
}

func (s *Server) Start(port int) error {
	portStr := fmt.Sprintf(":%d", port)
	s.Logger.Infof("Starting server on port %d", port)

	var handler http.Handler = s.RootRouter

	s.Server = &http.Server{
		Addr:         portStr,
		Handler:      handler,
		ReadTimeout:  time.Duration(10) * time.Second,
		WriteTimeout: time.Duration(10) * time.Second,
		IdleTimeout:  time.Duration(10) * time.Second,
	}

	s.didFinishListen = make(chan struct{})

	go func() {
		var err error
		err = s.Server.ListenAndServe()

		if err != nil && err != http.ErrServerClosed {
			s.Logger.Fatal("Error starting server", err)
			time.Sleep(time.Second)
		}

		close(s.didFinishListen)
	}()

	return nil
}

func (s *Server) StopHTTPServer() {
	if s.Server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		didShutdown := false
		for s.didFinishListen != nil && !didShutdown {
			if err := s.Server.Shutdown(ctx); err != nil {
				s.Logger.Warn("Unable to shutdown server", err)
			}
			select {
			case <-s.didFinishListen:
				didShutdown = true
			}
		}
		err := s.Server.Close()
		if err != nil {
			s.Logger.Fatal("Error close server")
		}
		s.Server = nil
	}
}

func (s *Server) Shutdown() {
	s.Logger.Info("Stopping Server")
	s.StopHTTPServer()

	if s.Store != nil {
		s.Store.Close()
	}

	_, timeoutCancel := context.WithTimeout(context.Background(), time.Second*15)
	defer timeoutCancel()

	s.Logger.Info("Server stopped")

	// should just write the "server stopped" record, the rest are already flushed.
	_, timeoutCancel2 := context.WithTimeout(context.Background(), time.Second*5)
	defer timeoutCancel2()
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.RootRouter.ServeHTTP(w, r)
}

func (s *Server) configureRouter() {
	s.RootRouter.Use(s.logRequest)
	s.RootRouter.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
}

func (s *Server) error(w http.ResponseWriter, code int, err error) {
	s.respond(w, code, map[string]string{"error": err.Error()})
}

func (s *Server) respond(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			s.Logger.Error(err)
		}
	}
}

func (s *Server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.Logger.WithFields(logrus.Fields{
			"Addr":       r.RemoteAddr,
			"Request id": r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	code int
}
