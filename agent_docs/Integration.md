# Food Tracker - Integration

**Back to:** [[01-Projects/Food Tracker/!Hub]]

## Obsidian Vault Paths

### Read Access (Food Database — deferred)
```
03-Resources/Health & Fitness/Nutrition Reference.md
```
- **MVP approach:** Use hardcoded preset foods (seeded into SQLite on first run)
- **Future:** Parse this file to populate foods table dynamically
- Detect file changes via mtime for re-parsing

### Write Access (Daily Logs)
```
02-Areas/Health & Fitness/Nutrition Logs/YYYY-MM Month.md
```
- App appends/updates daily entries
- Format must match existing log structure exactly
- One file per month (append to existing)

### Example Write Path
```
02-Areas/Health & Fitness/Nutrition Logs/2026-01 January.md
02-Areas/Health & Fitness/Nutrition Logs/2026-02 February.md
```

## Syncthing Configuration

### Current Setup
- Obsidian vault synced between: ThinkPad, Phone, Raspberry Pi
- KeePass database also synced

### Food Tracker Requirements

**Pi needs:**
- Read access to entire vault (for future Nutrition Reference.md parsing)
- Write access to `Nutrition Logs/` folder

**Recommended `.stignore` on Pi:**
```
// Ignore SQLite database to prevent sync conflicts
food-tracker.db
food-tracker.db-journal
food-tracker.db-wal
food-tracker.db-shm
```

**Folder permissions:**
- `Nutrition Logs/` folder should be set to "Send & Receive" on Pi
- Rest of vault can be "Receive Only" on Pi (safer)

### Conflict Prevention

**Rule:** The app is the sole writer for current/future daily logs.

- App writes to today and forward
- Past days are read-only in app (editable in Obsidian)
- If Syncthing conflict occurs, `.sync-conflict` files will appear — manual resolution needed

## Claude Analysis Workflow

### Current (ThinkPad Only)
Claude has MCP access to Obsidian vault via `mcp-obsidian` server.

**Prompt:**
> "Check my Obsidian vault and analyze today's food intake"

**What Claude does:**
1. Reads `Nutrition Logs/2026-XX Month.md`
2. Finds today's entry
3. Calculates totals (including estimates for custom foods)
4. Compares to targets
5. Provides analysis and recommendations

### Enabling Claude on Phone

**Option A: Copy-paste from Obsidian**
- Open today's log in Obsidian on phone
- Copy the day's entry
- Paste into Claude app
- Ask for analysis

**Option B: Future — In-app Claude button (V2)**
- Add button in Food Tracker UI
- Calls Claude API directly with day's data
- Shows analysis in-app

### Analysis Prompt Template

For manual analysis, use this prompt:

```
Here's my food log for today. Please analyze:

1. Calculate estimated macros for any custom foods
2. Total calories and protein
3. Compare to my targets:
   - Rest Day: 2,040-2,140 cal | 120-140g protein
   - Workout Day: 2,400-2,600 cal | 140-210g protein
4. Flag any protein shortfalls
5. Suggest adjustments if needed

[paste day's entry here]
```

## API Endpoints Summary

**All endpoints return HTML (not JSON).** HTMX swaps the response into the DOM.

| Endpoint | Method | Returns |
|----------|--------|---------|
| `/day/:date` | GET | Full day view page |
| `/day/:date/type` | POST | Header partial (day type updated) |
| `/day/:date/weight` | POST | Weight input partial |
| `/meals/:meal/save` | POST | Meal section partial + OOB totals |
| `/meals/:meal/yesterday` | POST | Meal section partial (pre-checked) |
| `/foods/custom` | POST | Updated food list partial |
| `/foods/search?q=...` | GET | Search results partial |

## Configuration

### Environment Variables

```bash
# Path to Obsidian vault
OBSIDIAN_VAULT_PATH=/path/to/Obsidian/Omni

# Specific paths within vault (relative to vault root)
NUTRITION_LOGS_PATH=02-Areas/Health & Fitness/Nutrition Logs

# SQLite database location
DATABASE_PATH=./data/food-tracker.db

# Server config
PORT=3000
HOST=0.0.0.0  # Allow LAN access
```

### Go Config File Alternative

Can also use a `config.yaml` or embed config in binary:

```yaml
vault:
  path: /home/pi/Sync/Obsidian/Omni
  nutrition_logs: 02-Areas/Health & Fitness/Nutrition Logs

database:
  path: ./data/food-tracker.db

server:
  port: 3000
  host: 0.0.0.0
```

## Build & Cross-Compile

### Development (ThinkPad)
```bash
# Single command — that's it
go run ./cmd/server

# Optional: use `air` for hot-reload on file changes
air
```

### Production Build
```bash
# Build for Pi (cross-compile) — templates + static embedded in binary
GOOS=linux GOARCH=arm64 go build -o dist/food-tracker ./cmd/server

# Deploy single binary to Pi
scp dist/food-tracker pi@raspberrypi.local:~/food-tracker/
```

### On Raspberry Pi
```bash
# Run directly
./food-tracker

# Or via systemd (auto-start on boot)
sudo systemctl start food-tracker
```

## Network Access

### Local Network (MVP)

**Option A: Static IP**
- Assign fixed IP to Pi (e.g., 192.168.1.100)
- Access via `http://192.168.1.100:3000`

**Option B: mDNS (Avahi)**
- Install avahi-daemon on Pi
- Access via `http://raspberrypi.local:3000`
- Easier to remember, works across reboots

### External Access (Stretch Goal)

**Option A: Cloudflare Tunnel**
- Free, no port forwarding needed
- Requires Cloudflare account
- Access via custom subdomain

**Option B: Tailscale**
- VPN-based, very secure
- Access Pi from anywhere on Tailscale network
- Requires Tailscale on all devices

## systemd Service

For auto-start on Pi boot:

```ini
# /etc/systemd/system/food-tracker.service
[Unit]
Description=Food Tracker
After=network.target

[Service]
Type=simple
User=pi
WorkingDirectory=/home/pi/food-tracker
ExecStart=/home/pi/food-tracker/food-tracker
Restart=on-failure
Environment=OBSIDIAN_VAULT_PATH=/home/pi/Sync/Obsidian/Omni

[Install]
WantedBy=multi-user.target
```

Enable with:
```bash
sudo systemctl enable food-tracker
sudo systemctl start food-tracker
```

---
**Tags:** #integration #obsidian #syncthing #claude #deployment #golang
