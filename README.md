# WireGuard Exporter Go

A Prometheus exporter for WireGuard VPN, written in **Go**. This is a complete rewrite inspired by the [original Rust version](https://github.com/kbknapp/wireguard_exporter), designed for better performance, maintainability, and extensibility.

---

## 🚀 Features

- ✅ **Complete metrics coverage** matching the Rust implementation
- ✅ **Real-time interface statistics** (total interfaces, peers per interface, bytes transferred)
- ✅ **Per-peer metrics** (bytes, handshake timing, endpoint information)
- ✅ **Peer aliases** - assign friendly names to peer public keys
- ✅ **Scrape performance metrics** - monitor collection duration and success
- ✅ **Configurable collection interval**
- ✅ **Verbose logging mode** for debugging
- ✅ **Lightweight** - minimal resource usage
- ✅ **No external dependencies** except Prometheus client and WireGuard control library

---

## 📊 Exported Metrics

| Metric Name | Type | Labels | Description |
|-------------|------|--------|-------------|
| `wireguard_interfaces_total` | Gauge | - | Total number of WireGuard interfaces |
| `wireguard_scrape_success` | Gauge | - | Whether the last scrape was successful (1=success, 0=failure) |
| `wireguard_scrape_duration_milliseconds` | Gauge | - | Duration of the scrape in milliseconds |
| `wireguard_peers_total` | Gauge | `interface` | Total number of peers per interface |
| `wireguard_bytes_total` | Gauge | `interface`, `direction` | Total bytes per interface (rx/tx) |
| `wireguard_peer_bytes_total` | Gauge | `interface`, `peer`, `alias`, `direction` | Total bytes per peer (rx/tx) |
| `wireguard_duration_since_latest_handshake` | Gauge | `interface`, `peer`, `alias` | Seconds since last handshake |
| `wireguard_peer_endpoint` | Gauge | `interface`, `peer`, `alias`, `endpoint_ip` | Peer endpoint information (static value of 1) |

---

## 🔧 Installation

### Prerequisites

- Go 1.22 or higher
- WireGuard installed on your system
- Root/sudo access (required to read WireGuard statistics)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/abolfazlnorozian/wireguard_exporter_go.git
cd wireguard_exporter_go

# Build the binary
go build -o wireguard_exporter ./cmd/wireguard_exporter

# Run with sudo (required for WireGuard access)
sudo ./wireguard_exporter
```

---

## 🎯 Usage

### Command Line Options

```bash
./wireguard_exporter [OPTIONS]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-listen-address` | `:9586` | Address to listen on for metrics server |
| `-metrics-path` | `/metrics` | Path under which to expose metrics |
| `-interval` | `15s` | Interval between metric collections |
| `-alias` | - | Comma-separated list of `publicKey:alias` entries |
| `-verbose` | `false` | Enable verbose logging |
| `-version` | - | Show version and exit |

### Examples

**Basic usage:**
```bash
sudo ./wireguard_exporter
```

**Custom port and collection interval:**
```bash
sudo ./wireguard_exporter -listen-address :9587 -interval 10s
```

**With peer aliases:**
```bash
sudo ./wireguard_exporter -alias "abc123...:server1,def456...:server2" -verbose
```

**Access metrics:**
```bash
curl http://localhost:9586/metrics
```

---

## 🏷️ Using Peer Aliases

Instead of seeing long public keys in your metrics, you can assign friendly names:

```bash
sudo ./wireguard_exporter \
  -alias "q2JWEKWfLPU5UjG2Sq31xx2GsSjdhKNtdT/X/tFVyjs=:kevin,2ELWFmGnqhtRpu4r2PUKc0cw+ELtuMPLd6l0KsoCUBQ=:jane"
```

This will add an `alias` label to peer metrics:

```
wireguard_peer_bytes_total{interface="wg0",peer="q2JW...",alias="kevin",direction="rx"} 123456
```

---

## 📈 Prometheus Configuration

Add this to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'wireguard'
    static_configs:
      - targets: ['localhost:9586']
    scrape_interval: 30s
```

---

## 🐳 Docker Deployment

Create a `Dockerfile`:

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o wireguard_exporter ./cmd/wireguard_exporter

FROM alpine:latest
RUN apk --no-cache add ca-certificates wireguard-tools
COPY --from=builder /app/wireguard_exporter /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/wireguard_exporter"]
```

Build and run:

```bash
docker build -t wireguard_exporter .
docker run --net=host --cap-add=NET_ADMIN wireguard_exporter
```

---

## 🖥️ systemd Service

Create `/etc/systemd/system/wireguard-exporter.service`:

```ini
[Unit]
Description=WireGuard Prometheus Exporter
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/wireguard_exporter -listen-address :9586 -interval 15s
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl daemon-reload
sudo systemctl enable wireguard-exporter
sudo systemctl start wireguard-exporter
sudo systemctl status wireguard-exporter
```

---

## 🏗️ Architecture

```
wireguard_exporter/
├── cmd/
│   └── wireguard_exporter/
│       └── main.go              # Entry point
├── internal/
│   └── metrics/
│       └── collector.go         # Metrics collection logic
├── pkg/
│   └── config/
│       └── config.go            # Configuration and flag parsing
├── go.mod
└── README.md
```

---

## 🔍 Troubleshooting

**Error: "cannot open wgctrl client"**
- Make sure you're running with `sudo` or as root
- Verify WireGuard is installed: `which wg`

**Error: "cannot list wireguard devices"**
- Check if WireGuard interfaces exist: `sudo wg show`
- Ensure WireGuard kernel module is loaded: `lsmod | grep wireguard`

**No metrics appearing:**
- Check if the exporter is running: `curl http://localhost:9586/metrics`
- Enable verbose mode: `-verbose`

---

## 🧪 Development

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Format code
go fmt ./...

# Run linter
golangci-lint run

# Build
go build -o wireguard_exporter ./cmd/wireguard_exporter
```

---

## 📝 License

MIT License — free to use, modify, and distribute.

---

## 🤝 Contributing

Contributions, bug reports, and feature ideas are welcome!  
Please open an issue or submit a pull request.

---

## 🙏 Acknowledgments

This project is inspired by the excellent [wireguard_exporter](https://github.com/kbknapp/wireguard_exporter) by Kevin K. written in Rust.

---

## 📚 References

- [WireGuard Official Site](https://www.wireguard.com/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Original Rust Implementation](https://github.com/kbknapp/wireguard_exporter)

