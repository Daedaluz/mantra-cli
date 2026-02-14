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
