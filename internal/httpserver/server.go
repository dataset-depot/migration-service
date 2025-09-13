package httpserver

import (
	"context"
	"net/http"
	"time"
)

type Opts struct {
	Addr					string
	ReadTimeout		time.Duration
	WriteTimeout	time.Duration
}

type Server struct {
	http 	*http.Server
}

func New(o Opts, mux http.Handler) *Server {
	return &Server{
		http: &http.Server{
			Addr:					o.Addr,
			Handler:			mux,
			ReadTimeout:	o.ReadTimeout,
			WriteTimeout: o.WriteTimeout,
		},
	}
}

func (s *Server) Start() error { return s.http.ListenAndServe() }
func (s *Server) Shutdown(ctx context.Context) error { return s.http.Shutdown(ctx) }

