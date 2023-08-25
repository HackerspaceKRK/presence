package presence

import (
	"context"
	"embed"
	"html/template"
	"log"
	"net"
	"net/http"
)

//go:embed static
var staticFiles embed.FS

//go:embed templates
var templateFiles embed.FS

func RunHTTPServer(ctx context.Context) error {
	cfg := ctx.Value(ConfigKey).(Config)
	log.Printf("Starting HTTP server on %s", cfg.HTTPListen)

	m := http.NewServeMux()
	m.Handle("/static", http.FileServer(http.FS(staticFiles)))
	tpls := template.Must(template.ParseFS(templateFiles, "templates/*.html"))

	m.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tpls.ExecuteTemplate(w, "layout.html", map[string]any{
			"Config": cfg,
		})
	}))

	httpServer := &http.Server{
		Addr:    cfg.HTTPListen,
		Handler: m,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	return httpServer.ListenAndServe()
}
