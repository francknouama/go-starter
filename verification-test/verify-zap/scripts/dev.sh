#!/bin/bash

# Development script for verify-zap
set -e

echo "🚀 Starting verify-zap development environment..."

# Check if .env file exists
if [ ! -f .env ]; then
    echo "📋 Creating .env file from .env.example..."
    cp .env.example .env
    echo "✅ Please edit .env file with your configuration"
fi

# Install dependencies
echo "📦 Installing dependencies..."
go mod download
go mod tidy

# Install development tools if not present
if ! command -v air &> /dev/null; then
    echo "🔧 Installing air for hot reload..."
    go install github.com/cosmtrek/air@latest
fi

# Start the application with hot reload
echo "🔥 Starting application with hot reload..."
echo "📍 Application will be available at http://localhost:8080"
echo "🩺 Health check: http://localhost:8080/health"
echo ""
echo "Press Ctrl+C to stop..."

# Check if air config exists, if not create a basic one
if [ ! -f .air.toml ]; then
    cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/server"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html", "yaml", "yml"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
EOF
fi

# Start with air for hot reload
air