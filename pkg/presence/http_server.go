package presence

import (
	"context"
	"embed"
	"log"
	"net"
	"net/http"
)

//go:embed static
var staticFiles embed.FS

func RunHTTPServer(ctx context.Context) error {
	cfg := ctx.Value(ConfigKey).(Config)
	log.Printf("Starting HTTP server on %s", cfg.HTTPListen)

	m := http.NewServeMux()
	m.Handle("/static/", http.FileServer(http.FS(staticFiles)))

	m.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dhcpWorker := r.Context().Value(DHCPWorkerKey).(*DHCPWorker)
		if !dhcpWorker.ConnectionOK() {
			errText := ""
			if dhcpWorker.LastError != nil {
				errText = dhcpWorker.LastError.Error()
			}
			ExecuteTemplateWithLayout(r.Context(), w, "message.html", map[string]any{
				"Title":            "Presence unavailable",
				"Details":          "Connection with the DHCP server has failed. Please try again in a few minutes.",
				"TechnicalDetails": errText,
			})
			return
		}
		ExecuteTemplateWithLayout(r.Context(), w, "index.html", map[string]any{
			"UsersInside": []string{"user1", "user2", "user3"},
		})
	}))

	m.Handle("/enroll", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ExecuteTemplateWithLayout(r.Context(), w, "enroll.html", map[string]any{
			"MAC":        "00:00:00:00:00:00",
			"Username":   "user1",
			"DeviceName": "user1's phone",
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
