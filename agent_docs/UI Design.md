# Food Tracker - UI Design

**Back to:** [[01-Projects/Food Tracker/!Hub]]

## Design Principles

- **Mobile-first** — designed for phone browser, works on desktop
- **Large tap targets** — easy to use with thumbs
- **Minimal taps** — common foods are checkboxes, not searches
- **Server-calculated totals** — totals update when meals are saved (not on every checkbox toggle)
- **Per-meal saving** — save button per section, not whole day
- **Classless styling** — Pico.css handles it, semantic HTML only

## Main Screen (Day View)

```
┌─────────────────────────────────────────────┐
│  ◀ Jan 30    Friday, Jan 31    Feb 1 ▶     │
│           [Rest Day ▼]  [83.1 kg ✏️]        │
├─────────────────────────────────────────────┤
│                                             │
│  BREAKFAST                      404 | 24g   │
│  ┌─────────────────────────────────────┐   │
│  │ ☑ Coffee (flat white)          72/2 │   │
│  │ ☑ Pumpkin bread + cream cheese 176/8│   │
│  │ ☑ Smoked salmon (half)        156/15│   │
│  │ ☐ Smoked trout (2/3 can)      144/18│   │
│  │ ☐ Greek yogurt + granola      388/20│   │
│  │ [+ Add food...]                     │   │
│  └─────────────────────────────────────┘   │
│  [💾 Save Breakfast]   [↺ Same as Yesterday]│
│                                             │
│  LUNCH                            ___/___   │
│  ┌─────────────────────────────────────┐   │
│  │ ☐ Chicken Dürüm (large)       690/50│   │
│  │ ☐ Work canteen turkey         565/43│   │
│  │ ☐ Weißwurst + pretzel         500/18│   │
│  │ [+ Add food...]                     │   │
│  └─────────────────────────────────────┘   │
│  [💾 Save Lunch]                            │
│                                             │
│  DINNER                           ___/___   │
│  ┌─────────────────────────────────────┐   │
│  │ ☐ Gnocchi + burrata           630/26│   │
│  │ ☐ Pizza (3 squares)           695/37│   │
│  │ ☐ Pho Bo                      500/25│   │
│  │ [+ Add food...]                     │   │
│  └─────────────────────────────────────┘   │
│  [💾 Save Dinner]                           │
│                                             │
│  SNACKS                           ___/___   │
│  ┌─────────────────────────────────────┐   │
│  │ ☐ Club Mate                   100/0 │   │
│  │ ☐ Beer (500ml)                210/2 │   │
│  │ ☐ Grapes (150g)               100/1 │   │
│  │ [+ Add food...]                     │   │
│  └─────────────────────────────────────┘   │
│  [💾 Save Snacks]                           │
│                                             │
│  SHAKE                            ___/___   │
│  ┌─────────────────────────────────────┐   │
│  │ ☐ Single scoop (24g) ⚠️       111/24│   │
│  │ ☐ Double scoop (48g) ✅       222/48│   │
│  │ ☐ Double + milk              318/53│   │
│  └─────────────────────────────────────┘   │
│  [💾 Save Shake]                            │
│                                             │
├─────────────────────────────────────────────┤
│  TODAY'S TOTAL                              │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  1,416 kcal    │    72g protein             │
│  Target: 2,040-2,140 │ 120-140g             │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  ⚠️ 48g protein remaining                   │
└─────────────────────────────────────────────┘
```

## Add Food Modal

Triggered by `[+ Add food...]` button in any meal section.

```
┌─────────────────────────────────────────────┐
│  Add to Breakfast                     [X]   │
├─────────────────────────────────────────────┤
│  [🔍 Search foods...                    ]   │
│                                             │
│  FAVORITES                                  │
│  • Coffee (flat white)                      │
│  • Pumpkin bread + cream cheese             │
│  • Smoked salmon (half pack)                │
│                                             │
│  RECENT                                     │
│  • Club Mate                                │
│  • Greek yogurt (150g)                      │
│                                             │
│  ALL FOODS                                  │
│  ▸ Proteins                                 │
│  ▸ Breads & Grains                          │
│  ▸ Dairy                                    │
│  ▸ Beverages                                │
│  ▸ Restaurant/Takeout                       │
│  ▸ Snacks                                   │
│                                             │
│  ─────────────────────────────────────────  │
│  [+ Add Custom Food]                        │
│                                             │
└─────────────────────────────────────────────┘
```

## Add Custom Food Modal

For foods not in the database — no macros, just name + description for Claude to estimate.

```
┌─────────────────────────────────────────────┐
│  Add Custom Food                      [X]   │
├─────────────────────────────────────────────┤
│  Name:                                      │
│  [Homemade pasta bake                   ]   │
│                                             │
│  Description/Portion:                       │
│  [Large bowl, lots of cheese, some veg  ]   │
│                                             │
│  [Add to Lunch]                             │
└─────────────────────────────────────────────┘
```

## Template & HTMX Breakdown

### Templates (Go html/template files)

**Full Page Templates:**
- `base.html` — Base layout: `<html>`, `<head>` (HTMX + Pico.css), `<body>` wrapper
- `day.html` — Day view: extends base, includes all sections for a given date

**Partial Templates (returned by HTMX endpoints):**
- `partials/meal_section.html` — Single meal with checkboxes, subtotals, save button. Returned after save to swap updated state.
- `partials/day_totals.html` — Running totals footer. Returned via OOB swap alongside meal saves.
- `partials/weight_input.html` — Weight display/edit field. Returned after weight save.
- `partials/food_search.html` — Search results list. Returned on search input.
- `partials/food_checkbox.html` — Single food checkbox row (for adding new items to a meal).

**Shared Components (included via `{{template}}`):**
- `components/header.html` — Date nav links + day type dropdown + weight input
- `components/add_food_modal.html` — Add food dialog UI

### HTMX Interactions Map

| User Action | HTMX Trigger | Endpoint | Target | Swap |
|-------------|-------------|----------|--------|------|
| Click "Save Breakfast" | `hx-post` on button | `POST /meals/breakfast/save` | `#breakfast-section` | `outerHTML` |
| Change day type dropdown | `hx-post` + `hx-trigger="change"` | `POST /day/:date/type` | `#day-type` | `outerHTML` |
| Enter weight + blur | `hx-post` + `hx-trigger="blur"` | `POST /day/:date/weight` | `#weight-input` | `outerHTML` |
| Click "Same as Yesterday" | `hx-post` on button | `POST /meals/breakfast/yesterday` | `#breakfast-section` | `outerHTML` |
| Type in food search | `hx-get` + `hx-trigger="keyup changed delay:300ms"` | `GET /foods/search?q=...` | `#search-results` | `innerHTML` |
| Submit custom food | `hx-post` on form | `POST /foods/custom` | `#meal-food-list` | `beforeend` |

### OOB (Out-of-Band) Swaps

When saving a meal, the Go handler returns TWO things in the response:
1. The updated meal section (primary target)
2. The updated day totals (via `hx-swap-oob="true"` attribute on the totals element)

This lets one POST update multiple parts of the page.

## Interaction States

### Checkbox Behavior
- Checkboxes are standard HTML `<input type="checkbox">` inside a `<form>`
- No HTMX on individual checkbox toggle (that would be too many requests)
- State is local until "Save" is clicked
- After save: server returns the section with checked state from DB

### Save Button States
- **Default** — "Save Breakfast"
- **Saving** — HTMX adds `htmx-request` class automatically (can style with CSS)
- **Saved** — Server returns section with subtle "✓ Saved" indicator
- **Error** — Server returns section with error message

### Loading Indicators
- HTMX has built-in support via `htmx-indicator` class
- Pico.css has a loading spinner on `<button aria-busy="true">`
- Combine: set `aria-busy="true"` during request via HTMX

## Mobile Optimization

- **Touch targets:** minimum 44×44px — Pico.css handles this well by default
- **Date navigation:** plain links (fast, no JS needed)
- **Sticky header:** CSS `position: sticky` on the date/weight header
- **Sticky footer:** CSS `position: sticky` on totals
- **No swipe gestures for V1** — keep it simple, links work fine

## Color Scheme

Pico.css handles light/dark mode automatically based on system preference. Custom accents via CSS variables if needed:
- **Protein warning:** amber/orange
- **Protein danger:** red (significantly under)
- **Protein good:** green (on track)
- **Save success:** brief green indicator
- **Custom food:** subtle different background

---
**Tags:** #ui-design #mockups #mobile #htmx #templates
