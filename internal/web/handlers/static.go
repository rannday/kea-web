package handlers

import (
	"embed"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rannday/kea-web/internal/utils"
)

// Embedded defaults
//go:embed static_defaults/**
var staticDefaultsFS embed.FS

func setStaticCacheHeaders(w http.ResponseWriter, filePath, ext string) {
  // “Always revalidate” files (unhashed / user likely to edit)
  if ext == ".html" || ext == ".webmanifest" || ext == ".ico" || ext == ".txt" || ext == ".xml" {
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    return
  }

  // If filename looks fingerprinted, allow immutable caching.
  // Examples: bundle.<hash>.css, app.<hash>.js, logo.<hash>.png
  // Also allow query-based cache busting via version param if you ever do that.
  base := filepath.Base(filePath)
  if looksFingerprinted(base) {
    w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
    return
  }

  // Default: cache but revalidate (good “override directory” behavior)
  w.Header().Set("Cache-Control", "public, max-age=0, must-revalidate")
}

func looksFingerprinted(name string) bool {
  // very simple heuristic: contains a dot-separated token with >= 8 hex chars
  // e.g. "bundle.1a2b3c4d5e6f.css" or "app.9f3c2b1a.js"
  parts := strings.Split(name, ".")
  if len(parts) < 3 {
    return false
  }
  for _, p := range parts[1 : len(parts)-1] { // middle parts
    if len(p) >= 8 && isHex(p) {
      return true
    }
  }
  return false
}

func isHex(s string) bool {
  for _, r := range s {
    if !(r >= '0' && r <= '9') && !(r >= 'a' && r <= 'f') && !(r >= 'A' && r <= 'F') {
      return false
    }
  }
  return true
}

// StaticFileHandler serves files from the /static directory
func StaticFileHandler(staticDir string) http.Handler {
  fs := http.FileServer(http.Dir(staticDir))

  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // r.URL.Path is already stripped by the router. Example: "/icons/x.png"
    filePath := strings.TrimPrefix(r.URL.Path, "/") // just remove leading slash

    fullPath := filepath.Join(staticDir, filePath)
    utils.Debug("Attempting to serve: %s", fullPath)

    if _, err := http.Dir(staticDir).Open(filePath); err != nil {
      utils.Debug("File not found: %s", fullPath)
      http.NotFound(w, r)
      return
    }

    ext := strings.ToLower(filepath.Ext(filePath))

    // Cache policy: see #3 below (fixed)
    setStaticCacheHeaders(w, filePath, ext)

    // Optional: content-type override map if you want it
    // (http.FileServer usually gets this right; keep yours if you prefer)
    mimeTypes := map[string]string{
      ".css":         "text/css",
      ".js":          "application/javascript",
      ".webmanifest": "application/manifest+json",
      ".png":         "image/png",
      ".jpg":         "image/jpeg",
      ".jpeg":        "image/jpeg",
      ".gif":         "image/gif",
      ".svg":         "image/svg+xml",
      ".webp":        "image/webp",
      ".ico":         "image/x-icon",
    }
    if contentType, ok := mimeTypes[ext]; ok {
      w.Header().Set("Content-Type", contentType)
    }

    fs.ServeHTTP(w, r)
  })
}

// RobotsTxt serves /robots.txt.
// Disk override: ./static/robots.txt
// Fallback: embedded default
func RobotsTxt() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/robots.txt" {
      http.NotFound(w, r)
      return
    }

    // Try disk override first
    diskPath := filepath.Join("static", "robots.txt")
    if data, err := os.ReadFile(diskPath); err == nil {
      w.Header().Set("Content-Type", "text/plain; charset=utf-8")
      w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
      _, _ = w.Write(data)
      return
    }

    // Embedded fallback
    data, err := staticDefaultsFS.ReadFile("static_defaults/robots.txt")
    if err != nil {
      utils.Error("Embedded robots.txt missing: %v", err)
      http.NotFound(w, r)
      return
    }

    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    _, _ = w.Write(data)
  }
}

// ServiceWorker serves /sw.js.
// Disk override: ./static/sw.js
// Fallback: embedded default
func ServiceWorker() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/sw.js" {
      http.NotFound(w, r)
      return
    }

    // Try disk override first
    diskPath := filepath.Join("static", "sw.js")
    if data, err := os.ReadFile(diskPath); err == nil {
      w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
      w.Header().Set("Service-Worker-Allowed", "/")
      w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
      _, _ = w.Write(data)
      return
    }

    // Embedded fallback
    data, err := staticDefaultsFS.ReadFile("static_defaults/sw.js")
    if err != nil {
      utils.Error("Embedded sw.js missing: %v", err)
      http.NotFound(w, r)
      return
    }

    w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
    w.Header().Set("Service-Worker-Allowed", "/")
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    _, _ = w.Write(data)
  }
}

// Manifest serves /site.webmanifest.
// Disk override: ./static/site.webmanifest
// Fallback: embedded default
func Manifest() http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/site.webmanifest" {
      http.NotFound(w, r)
      return
    }

    diskPath := filepath.Join("static", "site.webmanifest")
    if data, err := os.ReadFile(diskPath); err == nil {
      w.Header().Set("Content-Type", "application/manifest+json; charset=utf-8")
      w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
      _, _ = w.Write(data)
      return
    }

    data, err := staticDefaultsFS.ReadFile("static_defaults/site.webmanifest")
    if err != nil {
      utils.Error("Embedded site.webmanifest missing: %v", err)
      http.NotFound(w, r)
      return
    }

    w.Header().Set("Content-Type", "application/manifest+json; charset=utf-8")
    w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
    _, _ = w.Write(data)
  }
}

// Stylesheets serves /css/*
// Disk override: <staticDir>/css/<file>
// Fallback: embedded bundled CSS (assets_dist/css/*)
func Stylesheets(staticDir string) http.Handler {
  diskRoot := filepath.Join(staticDir, "css")

  // Embedded FS contains files at root like: "bundle.<hash>.css"
  embedded := http.FileServer(BundledCSSFS())

  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Expect URL like /css/<file>
    name := filepath.Base(r.URL.Path)

    // Basic sanity: only serve css
    if !strings.HasSuffix(strings.ToLower(name), ".css") {
      http.NotFound(w, r)
      return
    }

    // Disk override first
    diskPath := filepath.Join(diskRoot, name)
    if data, err := os.ReadFile(diskPath); err == nil {
      w.Header().Set("Content-Type", "text/css; charset=utf-8")
      w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
      w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat)) // cheap-ish; optional
      _, _ = w.Write(data)
      return
    }

    // Embedded fallback: strip "/css/" so embedded FileServer sees "/bundle.<hash>.css"
    w.Header().Set("Content-Type", "text/css; charset=utf-8")
    w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
    http.StripPrefix("/css/", embedded).ServeHTTP(w, r)
  })
}

// Javascripts serves /js/*
// Disk override: <staticDir>/js/<file>
// Fallback: embedded bundled JS (assets_dist/js/*)
func Javascripts(staticDir string) http.Handler {
  diskRoot := filepath.Join(staticDir, "js")

  embedded := http.FileServer(BundledJSFS())

  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Expect URL like /js/<file>
    name := filepath.Base(r.URL.Path)

    // Basic sanity: only serve js
    if !strings.HasSuffix(strings.ToLower(name), ".js") {
      http.NotFound(w, r)
      return
    }

    // Disk override first
    diskPath := filepath.Join(diskRoot, name)
    if data, err := os.ReadFile(diskPath); err == nil {
      w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
      w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
      w.Header().Set("Last-Modified", time.Now().UTC().Format(http.TimeFormat)) // optional
      _, _ = w.Write(data)
      return
    }

    // Embedded fallback
    w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
    w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
    http.StripPrefix("/js/", embedded).ServeHTTP(w, r)
  })
}
