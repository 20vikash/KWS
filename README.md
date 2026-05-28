# KWS — Self-Hosted Cloud Platform

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8.svg)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-✓-2496ED.svg)](https://docker.com)

KWS is a self-hosted cloud platform that gives each user a private, VPN-protected LXC container with a browser-based VS Code environment, managed PostgreSQL databases, and custom domain publishing — all behind WireGuard.

---

## Architecture

```
                    ┌─────────────────┐
                    │    Internet     │
                    └────────┬────────┘
                             │ :80/:443 (TLS)
                    ┌────────▼────────┐
                    │   nginx_proxy   │
                    └────────┬────────┘
                             │ proxy_pass :8080
                    ┌────────▼────────┐
                    │   kws_gateway   │   Go app (chi router)
                    │   (host net)    │   ─ sessions / auth / MQ
                    └──┬───┬───┬───┬──┘
                       │   │   │   │
          ┌────────────┼───┼───┼───┼──────────────────┐
          │            │   │   │   │                  │
   ┌──────▼──┐  ┌──────▼──┐┌──▼───┐  ┌──────────────▼──┐
   │Postgres │  │  Redis   ││RabbitMQ│ │Postgres  Adminer│
   │  (main) │  │ (cache)  ││  (MQ) │ │(pg service)     │
   └─────────┘  └──────────┘└───────┘ └─────────────────┘
                                       ┌──────────────┐
                                       │   dnsmasq    │
                                       └──────────────┘

              ┌─────────────────────────────────┐
              │         LXD Host (via socket)    │
              │  ├── lxdbr0 bridge (172.30.0/24) │
              │  ├── wg0 WireGuard (10.0.0/24)   │
              │  └── Per-user LXC containers     │
              │      (Ubuntu 22.04 + code-server)│
              └─────────────────────────────────┘
```

**Stack**: Go · PostgreSQL · Redis · RabbitMQ · LXD (LXC) · WireGuard · Nginx · Docker Compose

**Per-user resources**: Isolated LXC container with nesting, bridged networking, SSH, and code-server (VS Code in browser). Users can register up to 3 WireGuard devices, provision PostgreSQL databases, and publish domains with automatic HTTPS.

---

## Prerequisites

| Requirement | Notes |
|---|---|
| **Ubuntu 22.04+** | Server or VM with snapd |
| **Docker Engine** | Installed by Ansible if you use it |
| **LXD** (snap) | Installed by Ansible if you use it |
| **WireGuard** | Installed by ansible if you use it |
| **Go 1.25+** | Needed to build kws_install binary |
| **Domain name** | Pointed to your server's public IP |
| **SSL certificate** | Wildcard recommended (e.g. LetsEncrypt) |

---

## Quick Start

### 1. Clone

```bash
git clone https://github.com/20vikash/kws.git
cd kws
```

### 2. Run the installer

```bash
cd install
go run .       # or use the pre-built binary: ./kws_install
```

The installer asks for:
- **Domain name** and **public IP** (If you don't have a public IP, give private IP to set ip up locally)
- **Wildcard SSL certificate paths** (defaults to LetsEncrypt paths derived from your domain)
- **Gmail SMTP credentials** (for email verification)
- **Service passwords** (auto-generated if you press Enter)
- **Instance memory limit**
- **WireGuard keypair** (auto-generated — public key shown for client configs)
- **Ansible target** (SSH details — press Enter at the IP prompt to skip)

It generates:
- `kws_config.yaml` — platform configuration
- `.env` — Docker Compose + Go app secrets
- `nginx/conf.d/main.conf` + `00-default.conf` — TLS termination
- `dnsmasq.conf` — internal DNS for LXC containers
- `ansible/inventory` — Ansible target details (if you provided SSH info)

### 3. Provision the server (Ansible)

If you provided SSH details during install, the installer offers to run the playbook automatically. Otherwise:

```bash
cp ansible/inventory.example ansible/inventory   # edit with your server IP
ansible-playbook -i ansible/inventory ansible/playbook.yaml
```

**What it does**: apt update/upgrade → installs Docker, LXD (snap), WireGuard, `golang-migrate` → configures LXD bridge + storage pool → sets up iptables (IP forwarding, NAT, WG↔LXD forwarding) → writes WireGuard config.

### 4. Run database migrations

```bash
make migrate_up
```

### 5. Start the platform

```bash
make up
```

This launches all Docker Compose services, attaches them to the LXD bridge, and tails logs.

---

## Makefile Reference

| Command | Description |
|---|---|
| `make up` | Start all services, attach to LXD bridge, tail logs |
| `make down` | Stop and remove all containers |
| `make stop` | Stop containers (preserves state) |
| `make start` | Restart stopped containers + reattach bridge |
| `make dv` | Delete main DB volumes (⚠️ destroys data) |
| `make dvs` | Delete PG service volume (⚠️ destroys data) |
| `make migrate_up` | Run all pending database migrations |
| `make migrate_down-N` | Rollback N migrations |
| `make migrate_down-all` | Rollback all migrations |
| `make create_migration` | Create a new migration pair |

---

## Configuration

### kws_config.yaml

Generated by `kws_install`. Key sections:

```yaml
domain: example.com
public_ip: 203.0.113.10

network:
  lxd_bridge_subnet: 172.30.0.0/24
  lxd_bridge_gateway: 172.30.0.1
  lxc_ip_start: 11              # .2-.10 reserved for Docker services
  dns_ip: 172.30.0.4

wireguard:
  interface_name: wg0
  address: 10.0.0.1/24
  listen_port: 51820
  keepalive_sec: 25

services:
  bridge_attach:                # Docker containers attached to LXD bridge
    - container: postgres.kws.services
      bridge: lxdbr0
      ip_cidr: 172.30.0.2/24
    - container: adminer.kws.services
      bridge: lxdbr0
      ip_cidr: 172.30.0.3/24
    - container: dnsmasq_kws
      bridge: lxdbr0
      ip_cidr: 172.30.0.4/24

limits:
  max_wg_devices_per_user: 3
  max_service_db_users: 5
  max_service_db_databases: 10
  user_domain_limit: 3
```

All values have sensible defaults. Edit `kws_config.yaml` and restart the gateway to apply changes.

### Environment variables (.env)

Secrets and connection strings. Generated by `kws_install`. Notable vars:

- `DB_*` — main Postgres connection
- `REDIS_*` — Redis cache
- `MQ_*` — RabbitMQ
- `GMAIL_*` / `SMTP_*` — outbound email
- `WG_PRIVATE_KEY` / `WG_PUBLIC_KEY` — WireGuard server keypair
- `PG_SERVICE_*` — user-facing Postgres service
- `SERVICES_SUBNET` / `SERVICES_GATEWAY` — Docker compose network
- `ATTACH_SERVICES` — bridge attach triples for `make up`

---

## IP Allocation

```
lxdbr0 (172.30.0.0/24)
├── .1          Gateway
├── .2          postgres.kws.services (Docker)
├── .3          adminer.kws.services  (Docker)
├── .4          dnsmasq_kws           (Docker)
├── .5 - .10    Reserved for future services
└── .11 - .254  Per-user LXC instances (auto-allocated)

wg0 (10.0.0.0/24)
├── .1          wg0 interface
└── .2 - .254   Per-user WireGuard peers
```

LXC instance IPs start at `.11` (configurable via `lxc_ip_start`). WireGuard peer IPs start at `.2`.

---

## Directory Structure

```
kws/
├── ansible/               # Server provisioning playbook
│   ├── playbook.yaml
│   └── roles/
│       ├── bootstrap/     # apt, migrate CLI
│       ├── docker/        # Docker CE + compose
│       ├── lxd/           # LXD snap + preseed
│       ├── wireguard/     # WG tools + key setup
│       └── iptables/      # Forwarding, NAT, UFW
├── install/               # kws_install binary (Go)
├── nginx/
│   ├── nginx.conf         # Global nginx config
│   └── conf.d/            # Server blocks (generated)
├── src/
│   ├── cmd/               # HTTP handlers + main
│   ├── consts/            # Config accessors + services
│   ├── internal/          # Core packages
│   │   ├── database/      # Postgres + Redis + migrations
│   │   ├── store/         # Data access layer
│   │   ├── docker/        # Docker API + nginx reload
│   │   ├── lxd/           # LXC lifecycle
│   │   ├── wg/            # WireGuard operations
│   │   ├── nginx/         # Dynamic nginx config gen
│   │   ├── mq/            # RabbitMQ pool
│   │   ├── gmail/         # SMTP sender
│   │   ├── kwsconfig/     # kws_config.yaml loader
│   │   └── env.go         # .env loader
│   ├── lxd/               # LXD helpers
│   ├── models/            # Data models
│   └── web/               # HTML templates + JS + CSS
├── util/
│   └── attach_to_bridge   # Attach Docker container to LXD bridge
├── compose.yaml           # Docker Compose
├── Makefile
├── dnsmasq.conf           # Internal DNS (generated)
└── kws_config.yaml        # Platform config (generated)
```

---

## License

MIT © 2025 Vikash

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run `go vet ./...` and `go build ./...` from `src/` and `install/`
5. Submit a pull request

Pre-commit hooks (whitespace, YAML, JSON validation) run on every PR via GitHub Actions.
