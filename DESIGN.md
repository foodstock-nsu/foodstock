# Design System Document

## 1. Overview & Creative North Star: "The Culinary Sanctuary"

This design system is built upon the "Culinary Sanctuary" philosophy. In a world of cluttered, loud food delivery interfaces, we lean into high-end editorial minimalism. We treat food not as a commodity, but as an experience.

The "North Star" for this system is **Tactile Ethereality**. We move away from the "template" look by utilizing intentional asymmetry—placing hero imagery slightly off-center—and using high-contrast typography scales that feel like a premium lifestyle magazine. By prioritizing breathing room and soft tonal transitions over rigid grids and borders, we create a digital environment that feels calm, premium, and inherently trustworthy.

---

## 2. Colors & Surface Philosophy

The palette is anchored in sophisticated neutrals, allowing the vibrant Emerald Primary to act as a beacon for action.

### The "No-Line" Rule

**Explicit Instruction:** Traditional 1px solid borders are strictly prohibited for sectioning or containment. Boundaries must be defined through background color shifts.

- To separate a header from a body, transition from `surface` to `surface-container-low`.
- To define a card, use `surface-container-lowest` sitting on a `surface` background.

### Surface Hierarchy & Nesting

Treat the UI as a series of physical layers, like stacked sheets of fine vellum.

- **Base Layer:** `surface` (#f8f9fa)
- **Secondary Sections:** `surface-container-low` (#f3f4f5)
- **Floating/Interactive Cards:** `surface-container-lowest` (#ffffff)
- **Overlays/Modals:** `surface-bright` (#f8f9fa)

### The "Glass & Gradient" Rule

To elevate the PWA experience, use **Glassmorphism** for floating navigation bars or sticky headers. Use the `surface` token at 80% opacity with a `20px` backdrop-blur.

- **Signature Textures:** For main CTAs (e.g., "Checkout"), use a subtle linear gradient: `primary` (#006c49) to `primary_container` (#10b981). This provides a "glow" that flat colors cannot replicate.

---

## 3. Typography

The system utilizes a dual-sans approach to create editorial depth.

- **Display & Headlines (Manrope):** Chosen for its geometric precision and modern warmth. Use `display-lg` for hero marketing moments and `headline-md` for restaurant names.
- **Body & Utility (Inter):** Chosen for its exceptional legibility at small sizes on mobile screens.

**Hierarchy Intent:**

- **The Power Gap:** Create high contrast between `headline-lg` and `body-md`. Large headlines should dominate the space, followed by significant whitespace, then compact, readable body text.
- **Letter Spacing:** Headlines should have a `-0.02em` tracking for a "tighter," more professional editorial feel.

---

## 4. Elevation & Depth

We convey hierarchy through **Tonal Layering** rather than structural lines.

### The Layering Principle

Depth is achieved by "stacking" surface tiers. Place a `surface-container-lowest` card on a `surface-container-low` section to create a soft, natural lift.

### Ambient Shadows

Shadows must be "unseen but felt."

- **Token:** `Shadow-Soft`
- **Values:** `0px 12px 32px`
- **Color:** Use `on-surface` at **4% opacity**. It should mimic natural, ambient light. Forbid dark grey or heavy drop shadows.

### The "Ghost Border" Fallback

If accessibility requires a container edge (e.g., a search input), use a **Ghost Border**: `outline-variant` (#bbcabf) at **20% opacity**. Never use 100% opaque borders.

---

## 5. Components

### Buttons

- **Primary:** Rounded `full`. Gradient fill (`primary` to `primary_container`). White text (`on-primary`).
- **Secondary:** Rounded `full`. Background: `secondary_container`. Text: `on-secondary_container`. No border.
- **Tertiary:** No background. Text: `primary`. Weight: Semi-bold.

### Cards & Lists

- **Constraint:** Forbid divider lines.
- **The Alternative:** Use vertical whitespace (32px) or a subtle background shift to `surface-container-high` to separate list items.
- **Rounding:** All food cards must use `lg` (2rem) corner radius to evoke a soft, organic feel.

### Input Fields

- **Style:** Minimalist. Background: `surface-container-lowest`.
- **Interaction:** On focus, the "Ghost Border" increases to 40% opacity, and a subtle `Shadow-Soft` is applied.

### Bespoke Component: The "Ingredient Float"

For food detail pages, use overlapping imagery where the dish (PNG) sits 20px outside its `surface-container` container, breaking the bounding box to create a 3D, high-end feel.

---

## 6. Do's and Don'ts

### Do

- **Do** use 24px-32px padding for all containers to maintain the "Apple-like" aesthetic.
- **Do** use `primary_fixed_dim` for "Ready" or "Delivered" status tags.
- **Do** prioritize imagery. All food photos should have a consistent 2:3 or 1:1 aspect ratio with `md` rounded corners.

### Don't

- **Don't** use black (#000000). Use `on-background` (#191c1d) for text to maintain softness.
- **Don't** use standard "Select" dropdowns. Use a bottom-sheet (Glassmorphism) for a native PWA feel.
- **Don't** crowd the interface. If an element feels "stuck" to another, double the whitespace.
