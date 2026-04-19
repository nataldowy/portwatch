# portwatch

Lightweight CLI daemon that monitors open ports and alerts on unexpected changes.

## Installation

```bash
go install github.com/yourusername/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/portwatch.git && cd portwatch && go build -o portwatch .
```

## Usage

Start the daemon with a default scan interval of 30 seconds:

```bash
portwatch start
```

Specify a custom interval and define allowed ports:

```bash
portwatch start --interval 60 --allow 22,80,443
```

When an unexpected port is detected, portwatch will alert you in the terminal:

```
[ALERT] Unexpected port opened: 8080 (PID: 3821, Process: python3)
```

### Common Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--interval` | `30` | Scan interval in seconds |
| `--allow` | none | Comma-separated list of allowed ports |
| `--log` | stdout | Path to log file |
| `--quiet` | false | Suppress stdout output |

### Example Config File

```yaml
interval: 30
allow:
  - 22
  - 80
  - 443
log: /var/log/portwatch.log
```

Run with a config file:

```bash
portwatch start --config portwatch.yaml
```

## License

MIT © 2024 yourusername