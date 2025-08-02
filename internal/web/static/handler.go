package static

import (
	"io/fs"
	"log"
	"net/http"
	"os"

	"pagos-cesar/web"
)

// noListFileSystem is a custom file system that prevents directory listing.
type noListFileSystem struct {
	fs http.FileSystem
}

// Open opens the file with the given name, but returns os.ErrNotExist for directories.
func (nlfs noListFileSystem) Open(name string) (http.File, error) {
	f, err := nlfs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, os.ErrNotExist
	}
	return f, nil
}

type StaticHandler struct {
	logger *log.Logger
}

func NewStaticHandler(logger *log.Logger) *StaticHandler {
	return &StaticHandler{
		logger: logger,
	}
}

func (sh *StaticHandler) RegisterRoutes(mux *http.ServeMux) {

	staticFS, err := fs.Sub(web.WebFS, "static")
	if err != nil {
		sh.logger.Fatalf("failed to create static sub file system: %v", err)
	}

	fileServer := http.FileServer(noListFileSystem{http.FS(staticFS)})

	h := http.StripPrefix("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sh.logger.Printf("Serving static file: %s", r.URL.Path)
		fileServer.ServeHTTP(w, r)
	}))

	mux.Handle("/static/", h)
}
