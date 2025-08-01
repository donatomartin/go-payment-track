package handler

import (
	"html/template"
	"log"
	"net/http"
)

type DashboardHandler struct {
	logger    *log.Logger
	templates *template.Template
}

func NewDashboardHandler(logger *log.Logger) *DashboardHandler {

	templates := template.Must(template.ParseGlob("web/templates/*.html"))

	return &DashboardHandler{
		logger:    logger,
		templates: templates,
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
