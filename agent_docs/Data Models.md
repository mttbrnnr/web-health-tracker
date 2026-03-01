# Food Tracker - Data Models

**Back to:** [[01-Projects/Food Tracker/!Hub]]

## SQLite Schema

```sql
-- Cached from Nutrition Reference.md
CREATE TABLE foods (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT,  -- 'proteins', 'breads', 'dairy', etc.
    calories_per_serving REAL,
    protein_per_serving REAL,
    fat_per_serving REAL,
    carbs_per_serving REAL,
    serving_description TEXT,  -- "half pack (62.5g)"
    is_favorite BOOLEAN DEFAULT FALSE,
    source_line TEXT  -- for debugging parse issues
);

-- Custom foods added via app (no macros)
CREATE TABLE custom_foods (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,  -- "large portion", "with extra cheese"
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Daily log entries
CREATE TABLE daily_logs (
    id INTEGER PRIMARY KEY,
    date DATE NOT NULL UNIQUE,
    day_type TEXT CHECK(day_type IN ('rest', 'workout')) DEFAULT 'rest',
    weight_kg REAL,
    notes TEXT,
    synced_to_obsidian BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Individual meal entries
CREATE TABLE meal_entries (
    id INTEGER PRIMARY KEY,
    daily_log_id INTEGER REFERENCES daily_logs(id),
    meal_type TEXT CHECK(meal_type IN ('breakfast', 'lunch', 'dinner', 'snacks', 'shake')),
    food_id INTEGER REFERENCES foods(id),  -- NULL if custom
    custom_food_id INTEGER REFERENCES custom_foods(id),  -- NULL if standard
    quantity REAL DEFAULT 1,  -- multiplier for serving
    saved_at DATETIME,  -- NULL until "Save Meal" pressed
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Offline sync queue (V1.5)
CREATE TABLE sync_queue (
    id INTEGER PRIMARY KEY,
    action TEXT,  -- 'meal_save', 'weight_log', etc.
    payload TEXT,  -- JSON
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    synced_at DATETIME
);
```

## Table Relationships

```
daily_logs (1) ←──── (many) meal_entries
                              │
                              ├── food_id → foods (standard items)
                              │
                              └── custom_food_id → custom_foods (free text)
```

## Markdown Output Format

Must match existing format in `Nutrition Logs/` exactly:

```markdown
### Friday, January 31
**Type:** Rest Day  
**Weight:** 83.1kg

**Meals:**
- **Breakfast:** Coffee + pumpkin bread + cream cheese + smoked salmon
  - 404 kcal | 24g protein
- **Lunch:** Chicken dürüm (large)
  - 690 kcal | 50g protein
- **Snacks:** Club Mate + [Custom: handful of nuts, medium portion]
  - 100 kcal | 0g protein (+ custom item)
- **Dinner:** [pending]
- **Shake:** Double scoop
  - 222 kcal | 48g protein

**Total:** 1,416 kcal | 72g protein

**Analysis:** [left blank for Claude]
```

## Format Rules

### Standard Foods
- Listed by name from `foods` table
- Macros calculated from `calories_per_serving` × `quantity`
- Multiple items joined with ` + `

### Custom Foods
- Wrapped in `[Custom: name, description]`
- No macros displayed (Claude estimates during analysis)
- Shows `(+ custom item)` in meal total line

### Pending Meals
- Unsaved meals show as `[pending]`
- No macro line for pending meals

### Day Types
- `Rest Day` → targets 2,040-2,140 cal | 120-140g protein
- `Workout Day` → targets 2,400-2,600 cal | 140-210g protein

## Parsing Nutrition Reference.md

The parser must extract from sections like:

```markdown
### Smoked Salmon - Stremellachs (Norwegian) 🔥
**Per 100g:**
- Calories: 249 kcal | Protein: 24g ⚡ | Fat: 17g | Carbs: 0g

**Servings:**
- Half pack (62.5g): 156 kcal, 15g protein
- Full pack (125g): 311 kcal, 30g protein
```

**Parser should extract:**
- Name: "Smoked Salmon - Stremellachs (Norwegian)"
- Category: (from section header, e.g., "Proteins")
- Multiple serving options as separate `foods` rows
- Each serving: name variant, calories, protein, fat, carbs

**Edge cases to handle:**
- Emoji markers (🔥, ⚡, ✅, ⚠️)
- Missing macros (some items only have cal/protein)
- Combo items in "Quick Reference Combos" section
- Items without per-serving breakdowns

---
**Tags:** #data-models #sqlite #markdown #schema
