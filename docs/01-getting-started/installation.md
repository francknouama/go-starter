# Installation Guide

Install go-starter on any platform in minutes. Choose the method that works best for you.

## 🚀 Quick Install (Recommended)

### Option 1: Go Install (Easiest)
```bash
go install github.com/francknouama/go-starter@latest
```
✅ **Pros**: Automatic updates, always latest version  
⚠️ **Requires**: Go 1.21+ installed

### Option 2: Homebrew (macOS/Linux)
```bash
brew install francknouama/tap/go-starter
```
✅ **Pros**: Easy updates with `brew upgrade`  
⚠️ **Note**: Currently in review process

### Option 3: Binary Download (Universal)
Download from [GitHub Releases](https://github.com/francknouama/go-starter/releases/latest)

## 📦 Platform-Specific Installation

### 🍎 macOS

#### Homebrew (Recommended)
```bash
# Add tap (one-time setup)
brew tap francknouama/tap

# Install go-starter
brew install go-starter

# Verify installation
go-starter version
```

#### Binary Download
```bash
# Apple Silicon (M1/M2)
curl -L https://github.com/francknouama/go-starter/releases/latest/download/go-starter-darwin-arm64.tar.gz | tar -xz
sudo mv go-starter /usr/local/bin/

# Intel Macs
curl -L https://github.com/francknouama/go-starter/releases/latest/download/go-starter-darwin-amd64.tar.gz | tar -xz
sudo mv go-starter /usr/local/bin/
```

### 🐧 Linux

#### Package Managers
```bash
# Ubuntu/Debian
wget https://github.com/francknouama/go-starter/releases/latest/download/go-starter_linux_amd64.deb
sudo dpkg -i go-starter_linux_amd64.deb

# RHEL/CentOS/Fedora
wget https://github.com/francknouama/go-starter/releases/latest/download/go-starter_linux_amd64.rpm
sudo rpm -i go-starter_linux_amd64.rpm

# Alpine Linux
wget https://github.com/francknouama/go-starter/releases/latest/download/go-starter_linux_amd64.apk
sudo apk add --allow-untrusted go-starter_linux_amd64.apk
```

#### Binary Download
```bash
# AMD64 (Most common)
curl -L https://github.com/francknouama/go-starter/releases/latest/download/go-starter-linux-amd64.tar.gz | tar -xz
sudo mv go-starter /usr/local/bin/

# ARM64 (Raspberry Pi 4, etc.)
curl -L https://github.com/francknouama/go-starter/releases/latest/download/go-starter-linux-arm64.tar.gz | tar -xz
sudo mv go-starter /usr/local/bin/
```

### 🪟 Windows

#### PowerShell (Administrator)
```powershell
# Download and extract
Invoke-WebRequest -Uri "https://github.com/francknouama/go-starter/releases/latest/download/go-starter-windows-amd64.zip" -OutFile "go-starter.zip"
Expand-Archive go-starter.zip -DestinationPath "C:\Program Files\go-starter"

# Add to PATH
$env:PATH += ";C:\Program Files\go-starter"
[Environment]::SetEnvironmentVariable("PATH", $env:PATH, [EnvironmentVariableTarget]::Machine)
```

#### Manual Installation
1. Download `go-starter-windows-amd64.zip` from [releases](https://github.com/francknouama/go-starter/releases/latest)
2. Extract to `C:\Program Files\go-starter\`
3. Add `C:\Program Files\go-starter\` to your PATH environment variable
4. Restart command prompt

#### Chocolatey (Coming Soon)
```bash
choco install go-starter
```

## 🛠️ Development Installation

### From Source
```bash
# Clone repository
git clone https://github.com/francknouama/go-starter.git
cd go-starter

# Install dependencies and build
make install

# Verify installation
go-starter version
```

### Docker (For Isolated Usage)
```bash
# Pull image
docker pull ghcr.io/francknouama/go-starter:latest

# Create alias for easy usage
echo 'alias go-starter="docker run --rm -v $(pwd):/workspace -w /workspace ghcr.io/francknouama/go-starter:latest"' >> ~/.bashrc
source ~/.bashrc

# Use normally
go-starter new my-project --type=web-api
```

## ✅ Verify Installation

After installation, verify everything works:

```bash
# Check version
go-starter version
# Expected: go-starter v2.x.x

# List available blueprints
go-starter list
# Expected: List of project types

# Test generation (dry run)
go-starter new test-project --type=cli --dry-run
# Expected: File list without creating files

# Check help system
go-starter --help
# Expected: Command overview and usage
```

## 🔧 Configuration (Optional)

### Global Configuration
```bash
# Create config directory
mkdir -p ~/.config/go-starter

# Create basic config file
cat > ~/.config/go-starter/config.yaml << EOF
default:
  author: "Your Name"
  email: "your.email@example.com"
  license: "MIT"
  go_version: "1.21"
  logger: "slog"
EOF
```

### Team Configuration
For team standardization, see [Configuration Guide](../02-user-guides/configuration.md).

## 🆘 Troubleshooting

### Common Issues

#### "command not found: go-starter"
```bash
# Check if binary is in PATH
which go-starter

# If using Go install, check GOPATH
echo $GOPATH
ls $GOPATH/bin/

# Add GOPATH/bin to PATH if missing
export PATH=$PATH:$GOPATH/bin
```

#### Permission Denied (Linux/macOS)
```bash
# Make binary executable
chmod +x go-starter

# Move to system directory with proper permissions
sudo mv go-starter /usr/local/bin/
```

#### "cannot verify signature" (macOS)
```bash
# Allow unsigned binary (one-time)
sudo spctl --add --label "go-starter" --type exec /usr/local/bin/go-starter
```

#### Windows Execution Policy
```powershell
# Allow script execution (if needed)
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

### Getting Help

- **Quick Issues**: [Troubleshooting Guide](../02-user-guides/troubleshooting.md)
- **Common Questions**: [FAQ](../02-user-guides/faq.md)
- **Report Issues**: [GitHub Issues](https://github.com/francknouama/go-starter/issues)

## 🎯 Next Steps

### New Users
✅ Installation complete! Continue with:
- **[Quick Start](quick-start.md)** - Generate your first project in 5 minutes
- **[Getting Started](getting-started.md)** - Complete tutorial with examples

### Experienced Users
Jump directly to:
- **[Blueprint Selection](../02-user-guides/blueprint-selection.md)** - Choose your project type
- **[CLI Reference](../03-reference/cli-commands.md)** - Full command documentation

### Teams
Set up standardization:
- **[Configuration Guide](../02-user-guides/configuration.md)** - Team settings and standards

---

**🎉 Ready to start building?** You're all set to create amazing Go projects!