package web

import (
	"net/http"

	"github.com/rannday/kea-web/internal/web/handlers"
	"github.com/rannday/kea-web/internal/web/handlers/pages"
)

func routes(_ any) http.Handler {
  mux := http.NewServeMux()

  mux.HandleFunc("/", pages.HandleIndex)

  mux.HandleFunc("/sw.js", handlers.ServiceWorker())
  mux.HandleFunc("/robots.txt", handlers.RobotsTxt())
  mux.HandleFunc("/site.webmanifest", handlers.Manifest())

  mux.Handle("/css/", handlers.Stylesheets("static"))
  mux.Handle("/js/", handlers.Javascripts("static"))

  // Other static files (favicon, icons, images, etc.)
  mux.Handle("/static/", http.StripPrefix("/static/", handlers.StaticFileHandler("static")))

  return mux
}
