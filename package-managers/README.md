# Package Manager Integration Guide

This directory contains configuration and setup files for distributing go-starter across multiple package managers and platforms.

## 📦 Supported Package Managers

### 1. Go Modules (Primary)
- **Status**: ✅ Ready - v2.1 with Progressive Disclosure
- **Installation**: `go install github.com/francknouama/go-starter/cmd/go-starter@latest`
- **Update**: Automatic with Go toolchain
- **Platform**: All Go-supported platforms
- **Features**: 20 blueprints, progressive disclosure system, simplified logger architecture

### 2. Homebrew (macOS/Linux)
- **Status**: 🔄 In Progress
- **Tap**: `francknouama/go-starter`
- **Installation**: `brew install francknouama/go-starter/go-starter`
- **Update**: `brew upgrade go-starter`

### 3. APT (Debian/Ubuntu)
- **Status**: 🕐 Planned
- **Repository**: `ppa:francknouama/go-starter`
- **Installation**: `sudo apt install go-starter`

### 4. YUM/DNF (RedHat/Fedora)
- **Status**: 🕐 Planned
- **Repository**: Custom RPM repository
- **Installation**: `sudo dnf install go-starter`

### 5. Chocolatey (Windows)
- **Status**: 🕐 Planned
- **Package**: `choco install go-starter`
- **Update**: `choco upgrade go-starter`

### 6. Scoop (Windows)
- **Status**: 🕐 Planned
- **Bucket**: `scoop bucket add go-starter https://github.com/francknouama/scoop-go-starter`
- **Installation**: `scoop install go-starter`

### 7. Snap (Universal Linux)
- **Status**: 🕐 Planned
- **Installation**: `sudo snap install go-starter`
- **Update**: Automatic

### 8. Flatpak (Universal Linux)
- **Status**: 🕐 Planned
- **Installation**: `flatpak install flathub org.francknouama.GoStarter`

## 🚀 Quick Installation

### Recommended (Go Users)
```bash
go install github.com/francknouama/go-starter/cmd/go-starter@latest
```

### Script Installation (All Platforms)
```bash
curl -sSL https://raw.githubusercontent.com/francknouama/go-starter/main/scripts/install.sh | bash
```

### Platform-Specific

#### macOS
```bash
# Homebrew (when available)
brew tap francknouama/go-starter
brew install go-starter

# Manual download
curl -L -o go-starter https://github.com/francknouama/go-starter/releases/latest/download/go-starter-darwin-amd64
chmod +x go-starter
sudo mv go-starter /usr/local/bin/
```

#### Linux
```bash
# Ubuntu/Debian (when available)
sudo add-apt-repository ppa:francknouama/go-starter
sudo apt update
sudo apt install go-starter

# Manual download
curl -L -o go-starter https://github.com/francknouama/go-starter/releases/latest/download/go-starter-linux-amd64
chmod +x go-starter
sudo mv go-starter /usr/local/bin/
```

#### Windows
```powershell
# Chocolatey (when available)
choco install go-starter

# Scoop (when available)
scoop bucket add go-starter https://github.com/francknouama/scoop-go-starter
scoop install go-starter

# Manual download
Invoke-WebRequest -Uri "https://github.com/francknouama/go-starter/releases/latest/download/go-starter-windows-amd64.exe" -OutFile "go-starter.exe"
# Move to PATH location
```

## 📋 Package Manager Specifications

### Homebrew Formula
- **Location**: `homebrew-formula/go-starter.rb`
- **Multi-platform**: ✅ (macOS Intel/ARM, Linux Intel/ARM)
- **Auto-update**: Via formula updates
- **Testing**: Built-in formula tests

### APT Package
- **Control File**: `debian/control`
- **Build Dependencies**: Go 1.21+
- **Package Dependencies**: None (static binary)
- **Architecture**: amd64, arm64

### RPM Package
- **Spec File**: `rpm/go-starter.spec`
- **Build Requirements**: golang >= 1.21
- **Runtime Dependencies**: None
- **Architecture**: x86_64, aarch64

### Chocolatey Package
- **Nuspec**: `chocolatey/go-starter.nuspec`
- **Install Script**: PowerShell-based
- **Dependencies**: None
- **Checksum**: SHA256 verification

### Snap Package
- **Snapcraft**: `snap/snapcraft.yaml`
- **Confinement**: strict
- **Grade**: stable
- **Architecture**: amd64, arm64

## 🔧 Maintenance

### Version Updates
All package managers are automatically updated through CI/CD when a new GitHub release is created:

1. **GitHub Release** triggers package builds
2. **Homebrew Formula** is updated with new checksums
3. **APT/YUM repositories** are rebuilt
4. **Chocolatey/Scoop** packages are submitted
5. **Snap/Flatpak** are automatically updated

### Manual Update Process
If automatic updates fail:

1. Update version in package configuration
2. Regenerate checksums for new binaries
3. Test installation on target platform
4. Submit to package manager repository
5. Verify availability after publication

### Quality Assurance
Each package manager integration includes:

- ✅ Installation testing
- ✅ Version verification
- ✅ Help command validation
- ✅ Basic functionality test
- ✅ Uninstallation verification

## 📊 Distribution Statistics

### Platform Coverage
- **Go Modules**: 100% (all Go platforms)
- **Homebrew**: 80% (macOS/Linux primary)
- **APT**: 60% (Debian-based Linux)
- **YUM/DNF**: 40% (RedHat-based Linux)
- **Windows**: 70% (via multiple managers)

### Installation Methods
1. **Go Install**: Recommended for developers
2. **Package Managers**: Recommended for system installation
3. **Direct Download**: For CI/CD and automation
4. **Script Install**: Quick setup for new users

## 🎯 Roadmap

### Phase 1: Core Distribution (Current)
- ✅ Go Modules
- 🔄 Homebrew
- ✅ Direct downloads
- ✅ Installation script

### Phase 2: Linux Package Managers
- 🕐 APT (Debian/Ubuntu)
- 🕐 YUM/DNF (RedHat/Fedora)
- 🕐 AUR (Arch Linux)

### Phase 3: Windows Package Managers
- 🕐 Chocolatey
- 🕐 Scoop
- 🕐 WinGet

### Phase 4: Universal Linux
- 🕐 Snap
- 🕐 Flatpak
- 🕐 AppImage

### Phase 5: Enterprise & Cloud
- 🕐 Docker images
- 🕐 Kubernetes Helm charts
- 🕐 Cloud provider marketplaces

## 🤝 Contributing

### Adding New Package Manager
1. Create configuration in appropriate subdirectory
2. Add to CI/CD automation (`.github/workflows/`)
3. Update this README with installation instructions
4. Test on target platform
5. Submit pull request

### Testing Package Installation
```bash
# Test script for all platforms
./scripts/test-package-installation.sh

# Test specific package manager
./scripts/test-homebrew.sh
./scripts/test-apt.sh
```

## 📞 Support

### Package Manager Issues
- **Homebrew**: Report to [tap repository](https://github.com/francknouama/homebrew-go-starter)
- **APT/YUM**: Check repository status
- **Windows**: Verify package manager is updated
- **General**: [Main repository issues](https://github.com/francknouama/go-starter/issues)

### Verification
```bash
# Verify installation
go-starter version

# Check installation path
which go-starter

# Validate functionality
go-starter --help
```

## 🔄 Version History & Compatibility

### Current Release: v2.1+
- **Progressive Disclosure**: Complete implementation with smart help filtering
- **Simplified Logger Architecture**: 60-90% code reduction while maintaining functionality  
- **20 Blueprints**: From simple CLI (8 files) to enterprise architectures
- **Enhanced Security**: Progressive disclosure security validation
- **Phase 2 Complete**: Enterprise production features with observability

### Version Compatibility Matrix

| Version | Features | Blueprint Count | Progressive Disclosure | Logger Architecture |
|---------|----------|-----------------|------------------------|-------------------|
| **v2.1+** | Full feature set | 20 blueprints | ✅ Complete system | ✅ Simplified (60-90% reduction) |
| **v2.0** | Enhanced blueprints | 12 blueprints | 🔄 Basic implementation | ❌ Complex implementation |
| **v1.4+** | Core functionality | 8 blueprints | ❌ Not available | ❌ Legacy logger system |

### Upgrade Recommendations
- **From v1.x**: Major upgrade recommended for progressive disclosure and simplified architecture
- **From v2.0**: Minor upgrade recommended for simplified logger architecture
- **New Users**: Install latest v2.1+ for complete feature set

---

**Last updated**: 2025-08-24  
**Version compatibility**: v2.1+ (Recommended) | v1.4+ (Minimum)  
**Go version**: 1.21+ required