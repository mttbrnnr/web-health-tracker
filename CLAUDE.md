# web-health-tracker

## What
Personal health tracking web app. Quick meal logging via checkboxes, weight tracking, and day type toggling. In a further step I want to add gym workout tracking. Runs on a Raspberry Pi 4 on the local network. Exports daily logs to an Obsidian vault in an existing markdown format.

## Tech Stack
- Go 1.21+
- Chi router
- SQLite via modernc.org/sqlite (pure Go, no CGO — required for ARM64 cross-compile)
- html/template for server-side rendering
- HTMX 2.0.x for dynamic interactions
- Pico.css (classless variant) for styling
- Target: Raspberry Pi 4 (linux/arm64)

## Project Structure
web-health-tracker/
├── cmd/server/        # Entry point (main.go)
├── internal/
│   ├── handler/       # HTTP handlers
│   ├── db/            # SQLite connection, schema, queries, seed data
│   ├── export/        # Markdown export to Obsidian
│   └── models/        # Shared structs
├── templates/
│   ├── partials/      # HTMX-swappable fragments
│   └── components/    # Reusable template pieces
├── static/            # htmx.min.js, pico.min.css
└── agent_docs/        # Detailed specs (read when relevant)

## How to Run
go run ./cmd/server
# → http://localhost:3000

## How to Test
go test ./...

## Key Conventions
- All endpoints return HTML, never JSON
- HTMX swaps HTML fragments — handlers return partials, not full pages
- SQLite is the source of truth; markdown export is one-way (SQLite → Obsidian)
- Pico.css is the classless variant — write semantic HTML only, no CSS classes
- Use HTMX OOB swaps to update day totals when saving a meal
- Pure Go only — no CGO anywhere, must cross-compile cleanly to linux/arm64

## Environment Variables
OBSIDIAN_VAULT_PATH   # Absolute path to Obsidian vault root
DATABASE_PATH         # Path to SQLite file (default: ./data/food-tracker.db)
PORT                  # Server port (default: 3000)

## Detailed Specs (read when relevant)
- Architecture & data flow:          agent_docs/architecture.md
- Database schema & markdown format: agent_docs/data-models.md
- UI mockups & HTMX interaction map: agent_docs/ui-design.md
- Obsidian paths & environment:      agent_docs/integration.md
