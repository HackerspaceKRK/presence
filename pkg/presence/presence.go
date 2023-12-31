package presence

import (
	"context"
	"log"
)

func Run() {
	ctx := context.Background()
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx = context.WithValue(ctx, ConfigKey, cfg)
	ctx = WithDHCPWorker(ctx)
	if err := RunHTTPServer(ctx); err != nil {
		log.Fatalf("Failed to run HTTP server: %v", err)
	}

}
