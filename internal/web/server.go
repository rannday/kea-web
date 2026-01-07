package web

import (
	"net/http"
	"time"

	"github.com/rannday/kea-web/internal/web/handlers"
)

type Server struct {
  httpServer *http.Server
}

func NewServer(addr string) *http.Server {
  handlers.SetBundledAssets(handlers.BundledCSS, handlers.BundledJS)
  
  s := &Server{}
  mux := routes(s)

  s.httpServer = &http.Server{
    Addr:         addr,
    Handler:      mux,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  60 * time.Second,
  }

  return s.httpServer
}
