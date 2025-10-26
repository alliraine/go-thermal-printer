<div align="center">

# go-thermal-printer

Lightweight, concurrent-friendly REST API for printing to ESC/POS thermal printers (currently verified only on EPSON TM-T88II; other ESC/POS models may work but are untested). Supports raw ESC/POS bytes, template rendering with helpers (bold, underline, italics, font B), status querying, and containerized deployment.

![License](https://img.shields.io/badge/license-MIT-green)
![Go Version](https://img.shields.io/github/go-mod/go-version/jonasclaes/go-thermal-printer)
![Status](https://img.shields.io/badge/api-v1-blue)

</div>

## ✨ Features

- RESTful API (Gin) with versioned routes under `/api/v1`
- Print raw ESC/POS byte payloads (Base64 friendly for JSON transport)
- Print using text templates with variable substitution
- Built-in template functions: `bold`, `underline`, `italic` / `italics`, `fontb`
- Safe sequential hardware access via internal print/status worker & channels
- Query printer status (printer, offline, error, paper) via dedicated endpoint
- Configurable via TOML + `CONFIG_PATH` environment variable override
- Docker & docker-compose ready
- Graceful context-based timeouts for print/status ops

## 🖨 Supported Printer

| Model | Status | Notes |
|-------|--------|-------|
| EPSON TM-T88II | ✅ Verified | Development & runtime testing performed on this model |
| Other ESC/POS printers | ? (Untested) | May work if they implement standard ESC/POS command set & status bytes |

If you successfully run another model, consider opening an issue so it can be listed.

## 🔢 Versioning & Releases

Semantic Versioning (MAJOR.MINOR.PATCH):

- Backwards-compatible additions increase MINOR
- Bug fixes increase PATCH
- Breaking API / contract changes increase MAJOR

"Latest" references the most recent published (non pre-release) GitHub Release tag (e.g. `v1.2.3`). Pre-releases (like `v1.3.0-rc1`) are not treated as "latest" for documentation examples.

Suggested Docker pull pattern (after first release is published):
```bash
# Exact version
docker pull ghcr.io/jonasclaes/go-thermal-printer:v1.2.3
# Track latest stable
docker pull ghcr.io/jonasclaes/go-thermal-printer:latest
```

Until a first stable release is cut, treat `main` (or the feature branch you build from) as potentially unstable.

## 🧱 Architecture Overview

```
┌────────────┐   HTTP (Gin)    ┌────────────────┐        ┌────────────────────┐
│  Client    │ ──────────────▶ │  Controllers   │ ─────▶ │  PrinterService    │
└────────────┘                 └────────────────┘        │ (timeout wrapper)  │
                                                 ▲                       └─────────┬──────────┘
                                                 │                                 │ channels
                                                 │                       ┌─────────▼──────────┐
                                                 │                       │   PrintService     │
                                                 │                       │  (worker goroutine)│
                                                 │                       ├─────────┬──────────┤
                                                 │                       │ printQ   │ statusQ │
                                                 │                       └────┬─────┴────┬────┘
                                                 │                            │          │
                                                 │                            ▼          ▼
                                                 │                        ESC/POS    ESC/POS
                                                 │                        serial hw  status cmds
                                                 │
                                        Templates (text/template + helper funcs)
```

## 📦 Endpoints

### Authentication: API Key

All API endpoints require an API key for access. Include the `X-Api-Key` header in your requests:

```http
X-Api-Key: <your-api-key-here>
```

The API key is configured in your TOML file under the `[server]` section as `api_key`.

Example `config.toml`:

```toml
[server]
host = "127.0.0.1"
port = 8080
api_key = "your-secret-key"
```

If the header is missing or invalid, requests will be rejected with an authentication error.


Base URL: `http://<host>:<port>` (default `http://127.0.0.1:8080`)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Liveness probe |
| GET | `/api/v1/printer/status` | Returns raw status bytes (printer/offline/error/paper) |
| POST | `/api/v1/printer/print` | Print raw ESC/POS payload (JSON) |
| POST | `/api/v1/printer/print-template` | Render & print template file with variables |

### Request / Response Examples

Health:
```http
GET /health
200 OK
```

Printer status:
```http
GET /api/v1/printer/status
X-Api-Key: <your-api-key-here>
200 OK
{
   "printerStatus": 0,
   "offlineStatus": 0,
   "errorStatus": 0,
   "continuousPaperStatus": 0
}
```

Raw print (data is Base64 in example):
```http
POST /api/v1/printer/print
X-Api-Key: <your-api-key-here>
Content-Type: application/json

{
   "data": "G3QT1Qo="
}
```

Template print:
```http
POST /api/v1/printer/print-template
X-Api-Key: <your-api-key-here>
Content-Type: application/json

{
   "templateFile": "templates/receipt.tmpl",
   "variables": {
      "storeName": "Coffee & More",
      "date": "2025-08-07",
      "time": "14:30:25",
      "orderNumber": "67890",
      "items": {"Coffee": 3.50, "Sandwich": 8.75},
      "subtotal": 12.25,
      "tax": 0.98,
      "total": 13.23,
      "customerName": "Jane Smith"
   }
}
```

## 🧪 Template System

Templates are standard Go `text/template` files. Example (`templates/receipt.tmpl`):

```gotemplate
{{ bold .storeName }}\n
Order: {{ .orderNumber }}\n
{{ underline "Items" }}\n
{{ range $name, $price := .items -}}
{{$name}}  {{$price}}\n
{{ end -}}

Subtotal: {{ .subtotal }}\n
Tax: {{ .tax }}\n
TOTAL: {{ bold (printf "%.2f" .total) }}\n
Thank you {{ italic .customerName }}!\n
```

Custom helpers wrap ESC/POS commands, producing styled output directly.

## ⚙️ Configuration

Environment variable:
- `CONFIG_PATH` (default: `config.toml` if present, else falls back to embedded defaults)

TOML structure:
```toml
[server]
host = "127.0.0.1"
port = 8080

[printer]
port = "/dev/ttyUSB0"   # e.g. Linux /dev/ttyUSB0, macOS /dev/tty.usbserial*, Windows COM3
baud_rate = 19200
data_bits = 8
stop_bits = 1            # 1 or 2
parity = 0               # 0=None,1=Odd,2=Even,3=Mark,4=Space
usb_mode = false         # Set true to stream bytes to a USB printer node (status polling disabled)
```

**USB mode** streams ESC/POS bytes directly to a raw USB printer device (for example `/dev/usb/lp0`).
When enabled (`usb_mode = true`) the service bypasses the serial driver entirely and uses the
USB transport in `pkg/escpos/usb.go`. Printer status queries require a bidirectional serial
connection, so they return an error while USB mode is active.

### Selecting the Serial Port

List ports (Linux):
```bash
ls /dev/ttyUSB* /dev/ttyACM* 2>/dev/null
```
Grant permissions (temporary):
```bash
sudo chmod a+rw /dev/ttyUSB0
```
Persistent approach: add your user to `dialout` (Linux):
```bash
sudo usermod -aG dialout $USER
```

## 🚀 Quick Start (Local)

```bash
git clone https://github.com/jonasclaes/go-thermal-printer.git
cd go-thermal-printer
cp config.example.toml config.toml
# adjust printer.port etc.
go run ./cmd/go-thermal-printer
```

Test health:
```bash
curl http://localhost:8080/health
```

Print using `requests.http` (REST Client / curl). Ensure Base64 content decodes to valid ESC/POS.

## 🐳 Docker

Build image:
```bash
docker build -t go-thermal-printer .
```

Run (Linux example):
```bash
docker run --rm \
   --device /dev/ttyUSB0 \
   -p 8080:8080 \
   -v $(pwd)/config.docker.toml:/app/config.toml \
   -v $(pwd)/templates:/app/templates:ro \
   -e CONFIG_PATH=/app/config.toml \
   go-thermal-printer
```

docker-compose:
```bash
docker compose up --build
```

Note: container must access the serial device (`--device`). On macOS with USB adapters inside Docker Desktop, device passthrough may vary.

## 🔐 Security Considerations

- Expose only on trusted networks; unauthenticated print endpoints can be abused
- Consider adding an API key / auth proxy in front (future enhancement)
- Validate template variables if coming from untrusted clients

## 🧵 Concurrency Model

All serial I/O is funneled through a single worker goroutine (`PrintService.worker`) using channels:
- `printQueue` (buffered) for print jobs
- `statusQueue` for status requests
This ensures commands never interleave on the serial line.

## ⏱ Timeouts

Each public operation (print/status) wraps requests with a 10s context timeout in `PrinterService`. Adjust there if needed.

## 🧰 Development

Run with live reload (suggested tool [air] or [fresh]) – not included by default. Minimal flow:
```bash
go build ./...
go test ./...
go run ./cmd/go-thermal-printer
```

## 🐛 Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| Permission denied opening port | User lacks group / device perms | Add user to `dialout`, adjust udev rules |
| Garbled characters | Wrong baud / code page | Match printer settings; ensure `baud_rate` correct |
| Nothing prints | Wrong port or cable | Verify port exists; try different USB adapter |
| Status always zero | Printer not replying / flow control | Confirm printer supports status commands & cable supports bi-directional comm |
| Template styles not applied | Printer resets unexpectedly | Ensure initialization not overridden mid-print |

Enable verbose serial debugging by wrapping `serial.Open` with logging (not yet built-in – PRs welcome).

## 🛣 Roadmap (Ideas)

- [ ] Authentication / API key middleware
- [ ] Structured logging (zap / zerolog)
- [ ] Metrics (Prometheus endpoint)
- [ ] Graceful shutdown & port close on SIGTERM
- [ ] Support for images / QR codes
- [ ] Hot reload of templates
- [ ] Pluggable transport (network printers / USB raw)

## 🤝 Contributing

Fork, create a feature branch, open a PR. Please include:
- Description & rationale
- Tests if adding logic
- Updated docs / examples

## 📄 License

MIT – see `LICENSE`.

## 🙋 FAQ

Q: How do I generate raw ESC/POS bytes?
A: Use a library or record printer output; encode the bytes (Base64) for JSON. You can also craft templates using helper functions.

Q: Can I run multiple printers?
A: Not yet; would require managing multiple `PrintService` instances bound to different serial ports (future enhancement).

Q: Does it support Windows?
A: Should, provided the serial library can open `COMx` ports. Adjust config accordingly.

---

Questions or feature requests? Open an issue.

</div>
