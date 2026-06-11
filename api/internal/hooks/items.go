package hooks

import (
	"log"

	"github.com/pocketbase/pocketbase/core"
)

// RegisterItemHooks attaches lifecycle hooks to the "items" collection.
// These hooks fire regardless of whether the record is created via the API,
// the Admin UI, or programmatically in Go code.
func RegisterItemHooks(app core.App) {
	// OnRecordCreate fires around the creation of a new record.
	// Calling e.Next() proceeds with the actual save — logic before it
	// runs pre-save, logic after it runs post-save.
	app.OnRecordCreate("items").BindFunc(func(e *core.RecordEvent) error {
		name := e.Record.GetString("name")
		log.Printf("[hook] Creating item: %q", name)

		// --- Pre-save logic ---
		// Example: normalize the name before saving
		// e.Record.Set("name", strings.TrimSpace(name))

		// Proceed with the actual database save.
		if err := e.Next(); err != nil {
			return err
		}

		// --- Post-save logic ---
		log.Printf("[hook] Item created successfully: id=%s name=%q", e.Record.Id, name)

		// Example: you could trigger a notification, update a counter,
		// write to an audit log, call an external API, etc.

		return nil
	})
}
