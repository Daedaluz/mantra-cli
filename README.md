# mantra-cli

A command-line interface for the Mantra authentication platform. Mantra provides passwordless authentication using WebAuthn/FIDO2 security keys through a challenge-response model.

## Overview

mantra-cli lets you manage domains, users, and clients on a Mantra server, and initiate authentication or registration flows that are completed by scanning a QR code with a FIDO2 authenticator.

## Authentication Flow

```
CLI sends Sign/CreateUser request via gRPC
        │
        ▼
Server returns a challenge (id + secret)
        │
        ▼
CLI signs JWT tokens using the challenge secret
and embeds them in URLs displayed as QR codes
        │
        ▼
User scans the QR code and completes the
WebAuthn challenge on their authenticator
        │
        ▼
CLI polls the server until the challenge is
completed, rejected, or expired
        │
        ▼
Server returns the signed assertion
(authenticator data, signature, user handle)
```

## Installation

Download a prebuilt binary from the [Releases](../../releases) page, or build from source:

```sh
go install github.com/daedaluz/mantra-cli@latest
```

## Usage

```
mantra-cli [command] [flags]
```

### Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-s, --server` | Server hostname | `mantra-api.inits.se:443` (or `$SERVER`) |
| `--plaintext` | Use unencrypted gRPC | `false` |
| `--skip-verify` | Skip TLS certificate verification | `false` |

Flag values are resolved in this order: **CLI flags > environment variables > active context > hardcoded defaults**.

### Context Management

mantra-cli supports kubectl-style context management. Connection parameters and client credentials are stored in `~/.mantra/config.json` so you don't have to pass them on every invocation.

The config file has two concepts:

- **APIs** — Define Mantra server connections (address, optional API key, TLS settings).
- **Contexts** — Reference an API and add domain/client credentials plus URL paths.

#### Setting up a context

```sh
# Add an API server
mantra-cli api add prod --server mantra-api.inits.se:443 --api-key <key>

# Or a local dev server using h2c (plaintext gRPC)
mantra-cli api add local --server localhost:8080 --plaintext

# Add a context that references the API
mantra-cli context add myctx --api prod \
  --domain example.com \
  --client-id <id> --client-secret <secret>

# Switch to a context (the first context is auto-selected)
mantra-cli context use myctx
```

Once a context is active, commands pick up server, plaintext, skip-verify, domain, client-id, client-secret, api-key, auth-path, and register-path automatically. You can still override any value with flags or environment variables.

#### `api` commands

| Command | Description |
|---------|-------------|
| `api list` | List all configured APIs |
| `api add <name>` | Add an API (`--server`, `--api-key`, `--plaintext`, `--skip-verify`) |

#### `context` commands

| Command | Description |
|---------|-------------|
| `context list` | List all contexts (`*` marks the active one) |
| `context add <name>` | Add a context (`--api` required, `--domain`, `--client-id`, `--client-secret`, `--register-path`, `--auth-path`) |
| `context use <name>` | Switch the active context |

### Commands

#### `admin` — Platform Administration

Requires `--api-key`. Manage domains at the platform level.

```sh
# Create a domain
mantra-cli admin createDomain example.com "Example" "My domain"

# List all domains
mantra-cli admin listDomains

# Delete a domain
mantra-cli admin deleteDomain example.com
```

#### `domainAdmin` — Domain Administration

Requires `--domain`, `--client-id`, and `--client-secret`. Manage users and clients within a domain.

```sh
# Create a user (displays a QR code for registration)
mantra-cli domainAdmin createUser --domain example.com \
  --client-id <id> --client-secret <secret> \
  -u user123 -n "Alice"

# Authenticate a user
mantra-cli domainAdmin authenticate --domain example.com \
  --client-id <id> --client-secret <secret> \
  -u user123
```

With an active context, the above simplifies to:

```sh
mantra-cli domainAdmin createUser -u user123 -n "Alice"
mantra-cli domainAdmin authenticate -u user123
```

#### `client` — Client API

Requires `--domain`, `--client-id`, and `--client-secret`. Initiate authentication challenges from an application.

```sh
# Request a signature challenge
mantra-cli client sign --domain example.com \
  --client-id <id> --client-secret <secret> \
  -u user123 -m "Approve transaction"
```

## Concepts

**Domains** — Top-level organizational units. Each domain has its own set of users, clients, and configuration.

**Users** — End users who register FIDO2 security keys and authenticate within a domain.

**Clients** — Applications that integrate with Mantra. Admin clients can manage users and keys; regular clients initiate authentication flows.

**Keys** — WebAuthn/FIDO2 credentials registered by users. Each user can have multiple keys across different authenticators.

## Building

```sh
go build -o mantra-cli .
```
