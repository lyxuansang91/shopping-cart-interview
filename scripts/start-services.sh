#!/usr/bin/env bash
set -euo pipefail

# Load common environment variables
if [[ -f .env.local ]]; then
  while IFS= read -r line || [[ -n "$line" ]]; do
    if [[ $line =~ ^[A-Za-z_][A-Za-z0-9_]*= ]]; then
      export "$line"
    fi
  done < .env.local
fi

# Store PIDs of all processes
declare -a service_pids

# Ensure cleanup on Ctrl-C
cleanup() {
  echo "⏹ Stopping services..."

  # Kill all Air processes first
  pkill -f "air -c .air.toml" || true

  # Kill all service processes
  for pid in "${service_pids[@]}"; do
    if kill -0 "$pid" 2>/dev/null; then
      kill "$pid" 2>/dev/null || true
    fi
  done

  # Wait for all processes to terminate
  wait
  exit
}
trap cleanup SIGINT SIGTERM

# List of services
services=(cart)

# Spin up each service
for s in "${services[@]}"; do
  (
    cd "services/$s" || exit
    # Load service-specific environment variables
    if [[ -f ".env.local" ]]; then
      while IFS= read -r line || [[ -n "$line" ]]; do
        if [[ $line =~ ^[A-Za-z_][A-Za-z0-9_]*= ]]; then
          export "$line"
        fi
      done < .env.local
    fi

    # Start the service with Air
    air -c .air.toml &
    pid=$!
    service_pids+=($pid)
    printf "🚀 %-15s PID %d\n" "$s-service" "$pid"
  )
done

# Wait for all services
wait
