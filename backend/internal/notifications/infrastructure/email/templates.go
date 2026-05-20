package email

import (
	"embed"
	"html/template"
)

//go:embed templates/*.html
var templateFS embed.FS

// emailTemplates holds every email body loaded once at startup. Each
// template is keyed by its filename (e.g. "magic_link.html").
var emailTemplates = template.Must(
	template.New("email").ParseFS(templateFS, "templates/*.html"),
)
