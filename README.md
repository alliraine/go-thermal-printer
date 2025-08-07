# go-thermal-printer
Thermal printer web API written in Go, with the EPSON TM-T88II in mind.

## Configuration

The application uses TOML configuration files to manage settings. You can specify the location of your configuration file using the `CONFIG_PATH` environment variable.

### Environment Variable
- `CONFIG_PATH`: Path to the TOML configuration file (default: `config.toml`)

### Configuration File Structure

```toml
[server]
host = "127.0.0.1"    # Server host address
port = 8080           # Server port

[printer]
port = "/dev/ttyUSB0" # Serial port (Windows: "COM1", Linux: "/dev/ttyUSB0")
baud_rate = 19200     # Serial communication baud rate
data_bits = 8         # Number of data bits
stop_bits = 1         # Number of stop bits (1 or 2)
parity = 0            # Parity: 0=None, 1=Odd, 2=Even, 3=Mark, 4=Space
```

### Usage Examples

Using default config file:
```bash
./go-thermal-printer
```

Using custom config file:
```bash
CONFIG_PATH=/path/to/your/config.toml ./go-thermal-printer
```

### Getting Started

1. Copy the example configuration file:
   ```bash
   cp config.example.toml config.toml
   ```

2. Edit `config.toml` with your printer's serial port settings

3. Run the application:
   ```bash
   go run cmd/go-thermal-printer/main.go
   ```
