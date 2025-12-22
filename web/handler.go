package web

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

//go:generate go run ./cmd/buildui

//go:embed dist
var embedded embed.FS

var (
	distFS fs.FS
	files  http.Handler
)

func init() {
	if info, err := os.Stat("web/dist"); err == nil && info.IsDir() {
		distFS = os.DirFS("web/dist")
	} else {
		sub, err := fs.Sub(embedded, "dist")
		if err != nil {
			panic(err)
		}
		distFS = sub
	}
	files = http.FileServer(http.FS(distFS))
}

// Handler serves the compiled UI bundle with client-side routing support.
func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if serveIfExists(w, r) {
			return
		}

		f, err := distFS.Open("index.html")
		if err != nil {
			http.Error(w, "ui bundle missing", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		info, err := f.Stat()
		if err != nil {
			http.Error(w, "ui bundle missing", http.StatusInternalServerError)
			return
		}
		serveContent(w, r, f, info)
	})
}

func serveIfExists(w http.ResponseWriter, r *http.Request) bool {
	cleaned := path.Clean(r.URL.Path)
	if cleaned == "/" || cleaned == "." {
		files.ServeHTTP(w, r)
		return true
	}

	target := strings.TrimPrefix(cleaned, "/")
	if target == "" {
		target = "index.html"
	}

	f, err := distFS.Open(target)
	if err != nil {
		return false
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return false
	}
	serveContent(w, r, f, info)
	return true
}

func serveContent(w http.ResponseWriter, r *http.Request, f fs.File, info fs.FileInfo) {
	if seeker, ok := f.(io.ReadSeeker); ok {
		_, _ = seeker.Seek(0, io.SeekStart)
		http.ServeContent(w, r, info.Name(), info.ModTime(), seeker)
		return
	}

	data, err := io.ReadAll(f)
	if err != nil {
		http.Error(w, "ui bundle missing", http.StatusInternalServerError)
		return
	}
	http.ServeContent(w, r, info.Name(), info.ModTime(), bytes.NewReader(data))
}
