package pages

import (
	"net/http"

	"github.com/rannday/kea-web/internal/web/handlers"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/" {
    http.NotFound(w, r)
    return
  }

  handlers.RenderTemplate(w, "index", handlers.PageData{
    Title: "Kea Web",
  })
}
