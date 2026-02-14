#!/bin/bash
set -euo pipefail

. ./env

cmd=${1:-}
if [[ -z "${cmd}" ]]; then
  echo "Usage: $0 <admin|client|domainAdmin> [subcommand [args...]]" >&2
  exit 1
fi
shift

# Build command-specific args
args=()
case "${cmd}" in
  admin)
    args+=("--api-key=${API_KEY}")
    ;;
  client)
    args+=("--client-id=${CLIENT_ID}" "--client-secret=${CLIENT_SECRET}" "-d" "${DOMAIN}")
    ;;
  domainAdmin)
    args+=("--client-id=${CLIENT_ID}" "--client-secret=${CLIENT_SECRET}" "-d" "${DOMAIN}")
    ;;
esac

# Execute with command-scoped flags following the command
./mantra_cli "${cmd}" "${args[@]}" "$@"
# Example:
# ./mantra_cli domainAdmin --client-id=... --client-secret=... -d authn.se createUser -n "Phone"
