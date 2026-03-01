# Food Tracker - Architecture

**Back to:** [[01-Projects/Food Tracker/!Hub]]

## Tech Stack

| Layer | Technology | Rationale |
|-------|------------|-----------|
| **Backend** | Go 1.21+ | Pi 4 ARM compatible, single binary, learning goal |
| **Router** | Chi | Lightweight, idiomatic, standard http interfaces |
| **Database** | SQLite (modernc.org/sqlite) | Pure Go, no CGO, cross-compiles cleanly |
| **Templating** | Go html/template | Built-in, no dependencies, server-side rendering |
| **Interactivity** | HTMX 2.0.x | HTML attributes for AJAX, no JS to write |
| **Styling** | Pico.css (classless) | CDN link, zero CSS classes needed, dark mode auto |
| **Server** | Raspberry Pi 4 | Low power, always-on home server |
| **Sync** | Syncthing (existing) | Already syncs Obsidian vault |
| **Editor** | Neovim | Developer preference |
| **AI Integration** | Claude API (optional) | Custom food entry via text/photo |

## How HTMX Works (Quick Reference)

HTMX lets any HTML element make HTTP requests and swap the response into the DOM. Instead of building a JSON API and a JavaScript frontend to consume it, the Go server returns **HTML fragments** and HTMX handles the swap.

Key attributes:
- `hx-get="/url"` — make a GET request
- `hx-post="/url"` — make a POST request
- `hx-target="#element-id"` — swap the response into this element
- `hx-swap="innerHTML"` — how to insert (innerHTML, outerHTML, beforeend, etc.)
- `hx-trigger="click"` — what triggers the request (click, change, submit, etc.)
- `hx-include="[name='field']"` — include form values in request

Example: A "Save Meal" button that POSTs checked foods and replaces the meal section with updated HTML:
```html
<button hx-post="/meals/breakfast/save" 
        hx-target="#breakfast-section"
        hx-swap="outerHTML">
  Save Breakfast
</button>
```
The Go handler processes the save, then renders and returns the updated breakfast section HTML.

## System Architecture

### Development Mode

```
┌─────────────────────────────────────────────────────────────┐
│                   ThinkPad (Development)                    │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Single Go Process (:3000)               │   │
│  │                                                      │   │
│  │   • Chi router                                      │   │
│  │   • html/template rendering                         │   │
│  │   • HTMX served via CDN (or local copy)             │   │
│  │   • Pico.css via CDN (or local copy)                │   │
│  │   • SQLite database                                 │   │
│  │   • Markdown export                                 │   │
│  └──────────────────────────────────────────────────────┘   │
│                                                             │
│  Tip: Use `air` or similar for hot-reload during dev       │
└─────────────────────────────────────────────────────────────┘
```

No proxy, no second terminal, no build step. Just `go run .` and open browser.

### Production Mode (Raspberry Pi)

```
┌─────────────────────────────────────────────────────────────┐
│                     RASPBERRY PI 4                          │
│                                                             │
│  ┌──────────────────────────────────────────────────────┐   │
│  │              Single Go Binary (:3000)                │   │
│  │                                                      │   │
│  │   • Serves full HTML pages (Go templates)           │   │
│  │   • HTMX partial responses for interactivity        │   │
│  │   • Embedded static assets (CSS, HTMX JS)           │   │
│  │   • SQLite database (SOURCE OF TRUTH)               │   │
│  │   • Markdown export to Obsidian vault               │   │
│  │   • Optional: Claude API for custom foods           │   │
│  └──────────────────────────────────────────────────────┘   │
│                    │                    │                   │
│                    ▼                    ▼                   │
│  ┌──────────────────────┐  ┌─────────────────────────────┐ │
│  │   SQLite Database    │  │   Obsidian Vault            │ │
│  │   (food-tracker.db)  │  │   (via Syncthing)           │ │
│  │                      │  │                             │ │
│  │   • foods (preset    │  │   WRITE ONLY:               │ │
│  │     + custom)        │  │   • Nutrition Logs/*.md     │ │
│  │   • daily_logs       │  │     (exported from SQLite)  │ │
│  │   • meal_entries     │  │                             │ │
│  └──────────────────────┘  └─────────────────────────────┘ │
│                            ▲                                │
│                            │ (one-way export)               │
└────────────────────────────┼────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                      SYNCTHING                              │
│         (distributes vault to all devices)                  │
└─────────────────────────────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                  CLAUDE (via Obsidian MCP)                  │
│         (reads logs for analysis when prompted)             │
└─────────────────────────────────────────────────────────────┘
```

## Data Flow

### CRITICAL: SQLite is Source of Truth

**Key architectural decision:** SQLite database is the **single source of truth**. Markdown files are **export targets only** (one-way sync).

```
SQLite Database (source of truth)
    ↓ (one-way export)
Obsidian Markdown (for reading/analysis)
```

**Why this matters:**
- ❌ NO bidirectional sync complexity
- ❌ NO fragile markdown parser with edge cases  
- ✅ Structured, queryable, reliable data
- ✅ Simple export to maintain Obsidian compatibility
- ✅ Claude can still read logs via Obsidian MCP

---

### Seeding Preset Foods (App Startup)

```
Hardcoded preset foods in Go (or simple CSV/JSON file)
  → INSERT into SQLite `foods` table on first run
  → Available in app immediately
```

**Preset foods (~10-15 common meals):**
- 3-4 breakfast options (e.g., "Coffee + Pumpkin Bread + Salmon")
- 5-8 lunch/dinner options (e.g., "Chicken Dürum", "Work Canteen Turkey + Rice")
- Standard portions pre-calculated
- Can be expanded over time

---

### Adding Custom Foods

**Phase 1: Text Description**
```
User types: "Chicken teriyaki bowl from work canteen"
  → HTMX POST to /foods/custom
  → Go backend calls Claude API directly
  → Claude estimates macros based on description
  → INSERT into SQLite
  → Return updated food list (HTMX swap)
  → Food immediately available for selection
```

**Phase 2: Photo Upload**
```
User uploads photo of nutrition label
  → HTMX POST with file to /foods/custom/photo
  → Go backend saves image temporarily
  → Call Claude API with image (vision capabilities)
  → Claude extracts: name, calories, protein, fat, carbs, serving size
  → INSERT into SQLite
  → Return updated food list
  → Food immediately available for selection
```

**Claude API Integration:**
- Simple HTTP call to `https://api.anthropic.com/v1/messages`
- Include image as base64 for vision requests
- Parse JSON response for macro data
- No Claude Code binary needed - direct API calls from Go

---

### Rendering a Page (Full Page Load)

```
Browser GET /day/2026-02-09
  → Chi router → Go handler
  → Query SQLite for day data + all available foods
  → Render full page template (base layout + day content)
  → Return complete HTML (includes HTMX + Pico.css via <head>)
```

### Saving a Meal (HTMX Partial)

```
User checks foods → clicks "Save Breakfast"
  → HTMX sends POST /meals/breakfast/save (with checked food IDs + quantities)
  → Go handler: 
      1. Saves to SQLite (meal_entries table)
      2. Calls export.GenerateDayMarkdown()
      3. Writes to Obsidian vault
      4. Renders ONLY the updated meal section + totals
  → HTMX swaps the returned HTML fragment into the page
```

### Weight Logging (HTMX Inline)

```
User types weight → blur/enter triggers hx-post="/weight"
  → Go handler: 
      1. Saves to SQLite
      2. Updates markdown export
  → Returns updated weight display HTML fragment
```

### Exporting to Obsidian (Background)

```
After any meal save or weight update:
  → Go handler calls export.GenerateDayMarkdown()
  → Reads day data from SQLite
  → Generates markdown in proper format
  → Appends/updates day entry in monthly log file
  → Writes to: /path/to/vault/02-Areas/Health & Fitness/Nutrition Logs/YYYY-MM Month.md
  → Syncthing propagates changes to all devices
  → Claude can read via Obsidian MCP when user asks for analysis
```

### Date Navigation

```
User clicks "Next Day" → standard link to /day/2026-02-10
  (full page load, no HTMX needed for navigation)
```

## Request/Response Pattern

The key insight: **every endpoint returns HTML, not JSON.**

| Action | Method | URL | Returns |
|--------|--------|-----|---------||
| View day | GET | `/day/:date` | Full HTML page |
| Save meal | POST | `/meals/:meal/save` | HTML fragment (meal section + totals) |
| Toggle day type | POST | `/day/:date/type` | HTML fragment (header) |
| Save weight | POST | `/day/:date/weight` | HTML fragment (weight display) |
| Same as yesterday | POST | `/meals/:meal/yesterday` | HTML fragment (meal section) |
| Add custom food (text) | POST | `/foods/custom` | HTML fragment (updated food list) |
| Add custom food (photo) | POST | `/foods/custom/photo` | HTML fragment (updated food list) |
| Search foods | GET | `/foods/search?q=...` | HTML fragment (search results) |

## Project Structure

```
web-health-tracker/
├── cmd/
│   └── server/
│       └── main.go             # Entry point
├── internal/
│   ├── handler/
│   │   ├── day.go              # Day view handler (full page)
│   │   ├── meals.go            # Meal save/load handlers (HTMX partials)
│   │   ├── weight.go           # Weight handler (HTMX partial)
│   │   ├── foods.go            # Food search/custom handlers
│   │   ├── claude.go           # Claude API integration
│   │   └── middleware.go       # Logging, etc.
│   ├── db/
│   │   ├── sqlite.go           # Database connection
│   │   ├── schema.go           # Table creation
│   │   ├── queries.go          # SQL queries
│   │   └── seed.go             # Preset foods seeding
│   ├── export/
│   │   └── markdown.go         # Export to Obsidian format
│   └── models/
│       └── models.go           # Structs (Food, Meal, DailyLog, etc.)
├── templates/
│   ├── base.html               # Base layout (<html>, <head>, nav)
│   ├── day.html                # Full day view page
│   ├── partials/
│   │   ├── meal_section.html   # Single meal section (reused by HTMX)
│   │   ├── day_totals.html     # Running totals footer
│   │   ├── weight_input.html   # Weight display/edit
│   │   ├── food_search.html    # Search results
│   │   └── food_checkbox.html  # Single food item
│   └── components/
│       ├── header.html         # Date nav + day type + weight
│       └── add_food_modal.html # Add food UI
├── static/
│   ├── htmx.min.js             # HTMX 2.0.x (local copy for offline)
│   ├── pico.min.css            # Pico.css (local copy for offline)
│   └── app.css                 # Minimal custom CSS overrides (if any)
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Build & Deploy Process

### Development
```bash
# Single command — that's it
go run ./cmd/server

# Optional: use `air` for hot-reload on file changes
air
```

### Production Build
```bash
# Build for Pi (cross-compile)
GOOS=linux GOARCH=arm64 go build -o dist/food-tracker ./cmd/server

# Deploy single binary + templates/static to Pi
scp -r dist/food-tracker templates/ static/ pi@raspberrypi:~/food-tracker/

# OR embed everything in the binary (no extra files to deploy)
# Uses Go's embed package to bundle templates/ and static/ into the binary
GOOS=linux GOARCH=arm64 go build -o dist/food-tracker ./cmd/server
scp dist/food-tracker pi@raspberrypi:~/food-tracker/
```

### On Raspberry Pi
```bash
# Run directly
./food-tracker

# Or via systemd (auto-start on boot)
sudo systemctl start food-tracker
```

## Why This Architecture?

### Go + html/template + HTMX (instead of Go + SvelteKit)

**Why the pivot from SvelteKit?**
- Eliminates TypeScript from the learning stack (focus on Go only)
- No Node.js dependency at all (not even for development)
- No JSON API to design — Go renders HTML directly
- No build step for frontend — templates are just HTML files
- Simpler mental model: server does everything, HTMX makes it dynamic
- Faster path to working POC

**Why Go's built-in html/template?**
- Zero dependencies (comes with Go)
- Good enough for this use case
- Template inheritance via `{{template "name" .}}` and `{{block "name" .}}`
- Learning opportunity (understanding server-side rendering)
- Can upgrade to `templ` later if templates get unwieldy

**Why HTMX?**
- Tiny (~14KB), no build step, just a `<script>` tag
- Makes standard HTML interactive via attributes
- Perfect for "forms + checkboxes + save buttons" type UI
- Natural fit with server-rendered Go templates
- Active community, stable API, not going anywhere

**Why Pico.css?**
- Classless: just write semantic HTML and it looks good
- One `<link>` tag, zero configuration
- Good form/input/button styling (important for checkbox UI)
- Auto dark mode based on system preference
- ~10KB, can swap it out later

### SQLite Driver Choice

**Using `modernc.org/sqlite` because:**
- Pure Go implementation (no CGO)
- Cross-compiles to ARM without toolchain setup
- Slightly slower than CGO-based drivers but simpler deployment
- Good enough for single-user food tracking app

### SQLite as Source of Truth (Not Markdown)

**Why SQLite → Markdown (not Markdown → SQLite)?**
- ❌ Markdown is a living document with inconsistent formatting
- ❌ Parser would break constantly as format evolves
- ❌ Bidirectional sync is complex and error-prone
- ✅ SQLite gives structured, queryable, reliable data
- ✅ Preset foods = simple INSERT statements
- ✅ Custom foods via Claude API = direct to database
- ✅ Export to markdown is trivial (one-way generation)
- ✅ Claude can still read via Obsidian MCP for analysis

---
**Tags:** #architecture #golang #htmx #chi #sqlite #pico-css
