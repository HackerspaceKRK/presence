package presence

import (
	
	"embed"
	"html/template"
	"log"
	
	"net/http"
)


//go:embed templates
var templateFiles embed.FS

var tpls = template.Must(template.ParseFS(templateFiles, "templates/*.html"))

func ExecuteTemplateWithLayout(w http.ResponseWriter, name string, data map[string]any) {
	commonParams := map[string]any{
		"Config": cfg,
	}

	mergedParams := map[string]any{}
	for k, v := range commonParams {
		mergedParams[k] = v
	}
	for k, v := range data {
		mergedParams[k] = v
	}
	buf := &bytes.Buffer{}
	if err := tpls.ExecuteTemplate(buf, name, mergedParams); err != nil {
		log.Printf("Failed to execute template %s: %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	commonParams["Content"] = template.HTML(buf.String())
	if err := tpls.ExecuteTemplate(w, "layout.html", commonParams); err != nil {
		log.Printf("Failed to execute template layout.html: %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
