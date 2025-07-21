# Installation Guide

## Quick Install (Recommended)

### Using Go Install
```bash
go install github.com/francknouama/go-starter@latest
```

### Using Homebrew
**Currently unavailable** due to publishing issues. Use Go install or binary download instead.

## Binary Downloads

Download the latest release for your platform from [GitHub Releases](https://github.com/francknouama/go-starter/releases/latest).

### Linux/macOS
```bash
# Linux AMD64
curl -L https://github.com/francknouama/go-starter/releases/download/v1.3.1/go-starter_1.3.1_Linux_x86_64.tar.gz | tar -xz
sudo mv go-starter /usr/local/bin/

# macOS Apple Silicon
curl -L https://github.com/francknouama/go-starter/releases/download/v1.3.1/go-starter_1.3.1_Darwin_arm64.tar.gz | tar -xz
sudo mv go-starter /usr/local/bin/

# macOS Intel
curl -L https://github.com/francknouama/go-starter/releases/download/v1.3.1/go-starter_1.3.1_Darwin_x86_64.tar.gz | tar -xz
sudo mv go-starter /usr/local/bin/
```

### Windows
```powershell
# PowerShell
Invoke-WebRequest -Uri "https://github.com/francknouama/go-starter/releases/download/v1.3.1/go-starter_1.3.1_Windows_x86_64.zip" -OutFile "go-starter.zip"
Expand-Archive go-starter.zip -DestinationPath .
# Add to PATH or move to desired location
```

## Package Managers

### Linux Packages
Available from [GitHub Releases](https://github.com/francknouama/go-starter/releases/latest):

```bash
# Debian/Ubuntu
wget https://github.com/francknouama/go-starter/releases/download/v1.1.0/go-starter_1.1.0_linux_amd64.deb
sudo dpkg -i go-starter_1.1.0_linux_amd64.deb

# RHEL/CentOS/Fedora  
wget https://github.com/francknouama/go-starter/releases/download/v1.1.0/go-starter_1.1.0_linux_amd64.rpm
sudo rpm -i go-starter_1.1.0_linux_amd64.rpm

# Alpine Linux
wget https://github.com/francknouama/go-starter/releases/download/v1.1.0/go-starter_1.1.0_linux_amd64.apk
sudo apk add --allow-untrusted go-starter_1.1.0_linux_amd64.apk
```

## From Source

```bash
git clone https://github.com/francknouama/go-starter.git
cd go-starter
make install
```

## Verify Installation

```bash
go-starter version
go-starter list
```

## Next Steps

- 🚀 **[Quick Start Guide](GETTING_STARTED.md)** - Your first project in 30 seconds
- 📖 **[Project Types Guide](BLUEPRINTS.md)** - Choose the right template
- ⚙️ **[Configuration Guide](CONFIGURATION.md)** - Customize your setup