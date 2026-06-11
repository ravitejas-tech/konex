package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"github.com/ravitejas/konex/api/internal/handlers"
	"github.com/ravitejas/konex/api/internal/hooks"
)

func main() {
	// Load .env file before PocketBase initialization.
	// This allows you to set custom app-level environment variables.
	// PocketBase's own config (--dir, --http) is handled via CLI flags.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Optional: override PocketBase's data directory via env var.
	// This injects the --dir flag so PocketBase uses the specified path.
	if dataDir := os.Getenv("PB_DATA_DIR"); dataDir != "" {
		found := false
		for _, arg := range os.Args {
			if arg == "--dir" {
				found = true
				break
			}
		}
		if !found {
			os.Args = append(os.Args, "--dir", dataDir)
		}
	}

	app := pocketbase.New()

	// Register custom API routes via the OnServe hook.
	// This fires after PocketBase's internal initialization is complete
	// and the router is ready to accept route bindings.
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		handlers.RegisterRoutes(se)
		return se.Next()
	})

	// Register record lifecycle hooks.
	// These fire on every record create/update/delete regardless of
	// whether the operation comes from the API, Admin UI, or Go code.
	hooks.RegisterItemHooks(app)

	// Start the PocketBase server.
	// This blocks and serves:
	//   - Admin UI at /_/
	//   - Built-in REST API at /api/
	//   - Custom routes at /api/custom/
	//
	// Use CLI flags to configure:
	//   ./api serve --http=0.0.0.0:8090 --dir=./pb_data
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
