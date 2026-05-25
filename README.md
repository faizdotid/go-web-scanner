# Web Scanner

Fast concurrent web vulnerability scanner powered by configurable exploit modules.

## Features

- Configurable exploit modules via `config.json`
- Concurrent worker pool with adjustable goroutines
- Context-aware HTTP requests with timeouts
- Pre-compiled validation regex per exploit
- Buffered, thread-safe result writing
- TLS certificate verification disabled for testing

## Installation

```bash
go build -o go-web-scanner main.go
```

## Usage

### List available exploits

```bash
./go-web-scanner
# or
./go-web-scanner -exploit 0
```

### Run a scan

```bash
./go-web-scanner -list targets.txt -exploit 1 -workers 50 -output results
```

### Options

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `./files/config.json` | Path to configuration file |
| `-list` | `""` | File containing target URLs |
| `-exploit` | `0` | Exploit index to run (1-based) |
| `-workers` | `20` | Number of concurrent workers |
| `-output` | `results` | Output directory for findings |

## Configuration

Edit `files/config.json` to add or modify exploit modules:

```json
{
    "exploits": [
        {
            "type": "wordpress",
            "description": "WordPress Configuration Backup",
            "file_path": "./files/wordpress.txt",
            "validation_criteria": "DB_USER|DB_PASSWORD|table_prefix",
            "save_as": "wordpress_backup.txt",
            "response": "body"
        }
    ],
    "configuration": {
        "timeout": 7,
        "request_headers": ["Mozilla/5.0 ..."]
    }
}
```

### Response types

- `body` — Match against response body
- `header` — Match against `Content-Type` header

## Disclaimer

This tool is intended for authorized security testing and research only. Always obtain proper permission before testing systems you do not own.
