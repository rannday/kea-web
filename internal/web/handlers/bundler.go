package handlers

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets_dist/**
var assetsDistFS embed.FS

func BundledCSSFS() http.FileSystem {
  sub, err := fs.Sub(assetsDistFS, "assets_dist/css")
  if err != nil {
    panic(err)
  }
  return http.FS(sub)
}

func BundledJSFS() http.FileSystem {
  sub, err := fs.Sub(assetsDistFS, "assets_dist/js")
  if err != nil {
    panic(err)
  }
  return http.FS(sub)
}
