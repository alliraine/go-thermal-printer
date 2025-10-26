# Thermal Printer Feature Guide

This guide explains every printing capability exposed by the service, and provides examples you can copy verbatim when integrating printers, building templates, or experimenting with graphics.

---

## 1. Service Overview

| Endpoint | Method | Purpose | Success Code |
|----------|--------|---------|--------------|
| `/api/v1/printer/status` | GET | Query printer status bytes (printer/offline/error/paper) | 200 |
| `/api/v1/printer/print` | POST | Print raw ESC/POS byte streams (Base64 encoded) | 201 |
| `/api/v1/printer/print-template` | POST | Render a Go template with data variables and print | 201 |
| `/api/v1/printer/print-image` | POST | Render a Base64 image as ESC/POS raster data and print | 200 |

The API runs with a worker queue so only one job hits the thermal printer at a time. Each request is wrapped in a 10 s timeout to prevent jobs from hanging forever.

---

## 2. Environment & Setup

1. **Install Go** (version ≥ the one declared in `go.mod`).
2. **Wire up your printer** and note the serial/USB device path (e.g. `/dev/ttyUSB0`).
3. **Configure runtime options** via `config.toml` (or override with `CONFIG_PATH`). Key printer settings:
   ```toml
   [printer]
   port = "/dev/ttyUSB0"
   baud_rate = 19200
   data_bits = 8
   stop_bits = 1
   parity = 0         # 0=None,1=Odd,2=Even,3=Mark,4=Space
   ```
4. **Run the service**:
   ```bash
   go run ./cmd/server
   # or, if you prefer Make tasks
   make run
   ```
5. **Verify connectivity** using the status endpoint:
   ```bash
   curl -s http://localhost:8080/api/v1/printer/status | jq
   ```

You should see four status bytes in the JSON payload. Any non-zero error fields mean the printer blocked the request.

---

## 3. Printing Raw ESC/POS Bytes

POST `/api/v1/printer/print`
```bash
curl -X POST http://localhost:8080/api/v1/printer/print \
  -H 'Content-Type: application/json' \
  -d '{
    "data": "GxsbbyB3b3JsZCEK"  # "hello world!\n" in Base64
  }'
```

**Notes**
- Payloads *must* be Base64. The service decodes and streams the bytes unchanged.
- You are responsible for wrapping formatting commands (ESC/POS opcodes). The printer will print bytes exactly as received.

---

## 4. Printing Templates

Template files live under `templates/`. Each template is a Go `text/template` file with helper functions injected at runtime. Print one with:

```bash
curl -X POST http://localhost:8080/api/v1/printer/print-template \
  -H 'Content-Type: application/json' \
  -d '{
    "templateFile": "templates/demo.tmpl",
    "variables": {
      "title": "Store 42",
      "msg": "Thank you!"
    }
  }'
```

### Template Rendering Flow
1. Template text is parsed, then executed with the `variables` map.
2. Helpers inject ESC/POS sequences.
3. The rendered buffer is written to the printer using the default code page (PC858) unless you override it yourself.

---

## 5. Template Helper Reference

Each helper emits ESC/POS sequences with sensible defaults so you don’t have to memorize opcodes. Combine them freely inside templates.

### 5.1 Text Formatting
| Helper | Description | Example |
|--------|-------------|---------|
| `bold text` | Toggle emphasized mode just for `text`. | `{{bold "Bold label"}}` |
| `underline text` | Wrap with underline on/off codes. | `{{underline "Underlined"}}` |
| `italic / italics` | ESC/POS italics (not all printers support). | `{{italic "Note"}}` |
| `fontb text` | Switch to font B while printing `text`. | `{{fontb "Small caps"}}` |
| `doubleWidth`, `doubleHeight`, `doubleSize` | Stretch text horizontally, vertically, or both. | `{{doubleSize "TOTAL"}}` |
| `invert text` | White-on-black text. | `{{invert "Reverse"}}` |

### 5.2 Alignment & Position
| Helper | Description |
|--------|-------------|
| `center text` | Centers then restores left alignment. |
| `left text` / `right text` | Forces left/right alignment. |
| `align "left|center|right"` | Returns the raw alignment command for manual composition. |

### 5.3 Spacing & Cutting
| Helper | Description |
|--------|-------------|
| `feed n` | Feed *n* full lines (0–255). |
| `feedDots n` | Feed *n* dot rows (fine control). |
| `lineSpacing dots` | Set line spacing until changed. |
| `cut` / `cut "partial"` | Trigger cutter. |
| `reset` | ESC @ plus the default code page (PC858). Always call after templates that might change fonts/code pages permanently. |

### 5.4 Word Wrapping & Rotation
- `wrap text [width]`: Soft-wraps long words, preserving whitespace, before printing. Default width: 32 chars.
  ```gotemplate
  {{wrap "Long description that should wrap" 28}}
  ```
- `rotate90 text`: Print vertical (90°) text where supported.

### 5.5 Font Options Helper
`fontOptions` accepts key/value pairs in a single call. Missing fields leave the current setting untouched.

Supported keys:
- `font`: `"A"`, `"B"`, `0`, `1`, `2`
- `width`, `height`: integers 1–8 (multipliers for double/triple size)
- `lineSpacing`: 0–255 dots
- `charSpacing`: 0–255 dots between characters
- `bold` / `emphasized`: bool (`true`/`false`)
- `underline`, `underlineLevel`: bool/int (0–2)
- `invert`, `negative`, `reverse`: bool
- `doubleStrike`: bool

Example:
```gotemplate
{{fontOptions "font" "B" "width" 2 "height" 2 "bold" true}}
{{center "Big Bold Title"}}
{{fontOptions "bold" false "width" 1 "height" 1}}
```

### 5.6 Images
Use the `image` helper when you already have Base64 data inline:
```gotemplate
{{image .logoB64 384}}
```
- `logoB64` must be a Base64 PNG/JPEG/GIF string.
- The helper resizes to `maxWidth` dots (default 384 ≈ 58 mm) and emits raster bytes plus a code-page reset.
- The API `/print-image` offers the same conversion for raw requests:
  ```bash
  curl -X POST http://localhost:8080/api/v1/printer/print-image \
    -F imageBase64=@qr.png.b64 \
    -F maxWidthDots=384
  ```

### 5.7 QR Codes
`qr` builds a QR image internally (using `skip2/go-qrcode`), converts it to raster bytes, and appends a code-page reset.

Usage:
```gotemplate
{{qr "https://example.com" "size" 8 "error" "Q" "border" 1}}
```
Options (key/value pairs):
- `size`|`module`|`scale`: dots per module (>=1). Higher numbers = larger code.
- `error`|`errorLevel`|`correction`|`ecc`: `"L"`, `"M"`, `"Q"`, `"H"` or 0–3.
- `border`|`margin`|`quietZone`: number of modules for the quiet zone (≤0 disables border).
- `maxWidth`|`width`: constrain raster width (good for very large QR codes).

Example template snippet:
```gotemplate
{{center "Scan Me"}}
{{qr .link "size" 8 "error" "H"}}
```
Payload:
```json
{
  "templateFile": "templates/link.tmpl",
  "variables": {
    "title": "Daily Promo",
    "msg": "https://example.com/promo"
  }
}
```

### 5.8 Material Icons
`icon` prints any Material Design glyph shipped with [`gio.tools/icons`](https://pkg.go.dev/gio.tools/icons).

```gotemplate
{{icon "ActionFace"}}
{{icon "navigation-menu" 72}}       {{/* width override */}}
{{icon "alert-warning" 64 3}}       {{/* width + extra feed lines */}}
```

- **Name matching is flexible** – the helper ignores case, dashes, and underscores, so `ActionFace`, `action_face`, or `action-face` all resolve to the same icon.
- **Sizing** – the first optional argument sets icon width in printer dots (defaults to 96). Icons keep their aspect ratio.
- **Feed control** – an optional second argument adds extra blank lines after the icon (0–255). Use this to separate glyphs from following text.
- **Code page reset** – the helper appends the default code page just like `image` and `qr`, so subsequent text prints correctly.
- **Behind the scenes** – icons are vector drawings rasterized through the Gio icon renderer, then converted with the same `ImageToRasterBytes` pipeline used by the `image` helper.

---

## 6. Sample Templates

### 6.1 `templates/demo.tmpl`
Demonstrates wrapping, custom fonts, rotation, spacing tweaks, QR codes, inversion, and cutting. Print via:
```bash
curl -X POST http://localhost:8080/api/v1/printer/print-template \
  -H 'Content-Type: application/json' \
  -d '{"templateFile":"templates/demo.tmpl","variables":{}}'
```

### 6.2 `templates/link.tmpl` & `templates/notify.tmpl`
Expect a title (`.title`) and message (`.msg`), automatically wrapped and centered.

### 6.3 Creating Your Own
1. Copy `templates/demo.tmpl` as a starting point.
2. Reference helpers as shown above.
3. Keep `{{reset}}` near the top when you expect to re-use the printer for other jobs—it clears prior formatting and restores PC858.
4. Always end templates with a small feed (`{{feed 2}}`) to make tearing easier, then call `{{cut}}` if your hardware supports it.

---

## 7. Troubleshooting

| Symptom | Likely Cause | Fix |
|---------|--------------|-----|
| Random characters / gibberish | Code page changed mid-print or data not wrapped correctly | Ensure templates call `{{reset}}` before text; avoid custom opcodes unless necessary. |
| Nothing prints, but API returns 201 | Printer offline/cover open OR job consumed by newlines | Check `/printer/status`; verify hardware, and ensure templates end with `{{feed}}`/`{{cut}}`. |
| Image prints as garbage | Base64 invalid or exceeds width | Validate Base64, set `maxWidth` ≤ printer width (usually 384). |
| QR code too small/large | Module size mismatch | Adjust `"size"` or `"maxWidth"` options in the `qr` helper. |

---

## 8. Extending the System

- Add new helpers in `pkg/template/template.go` and register them in `getTemplateFuncs`.
- Put shared documentation in `Docs/` so new features have clear examples.
- Run formatting/tests before committing:
  ```bash
  gofmt -w pkg/template/template.go
  go test ./...
  ```
  *(The test suite currently takes >20 s in some environments; re-run if it times out.)*

---

## 9. Support & Contribution Tips

- Prefer templates over raw byte payloads—they self-document formatting and automatically reapply the default code page after QR/images.
- Open issues/PRs with: description, hardware model, ESC/POS logs (hex) if possible.
- If you extend helpers, update this guide with examples so the feature set stays discoverable.

Happy printing!
