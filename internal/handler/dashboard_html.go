package handler

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

type DashboardHandler struct {
	logger    *log.Logger
	templates *template.Template
}

func NewDashboardHandler(logger *log.Logger, templates fs.FS) *DashboardHandler {

	t := template.Must(template.ParseFS(templates, "*.html"))

	return &DashboardHandler{
		logger:    logger,
		templates: t,
	}
}

func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	data := struct {
		Title string
	}{
		Title: "Admin Dashboard",
	}

	if err := h.templates.ExecuteTemplate(w, "dashboard.html", data); err != nil {
		http.Error(w, "Failed to render dashboard: ", http.StatusInternalServerError)
		h.logger.Fatal("Template error:", err)
	}

}
