# ascii Package

This package provides ASCII art rendering functionality for the go-starter CLI tool with consistent styling and configuration options.

## Overview

The ascii package is responsible for displaying stylized text banners and logos in the terminal, enhancing the user experience with visual branding. It uses lipgloss for consistent styling across all CLI commands and supports comprehensive configuration options.

## Key Components

### Types

- **`BannerStyle`** - Enumeration for banner display styles (Full, Minimal, None)
- **`BannerConfig`** - Configuration structure for banner display behavior

### Core Functions

- **`Banner() string`** - Returns the full ASCII art banner with default config
- **`BannerWithConfig(config *BannerConfig) string`** - Returns banner with custom configuration
- **`Logo() string`** - Returns smaller ASCII logo with default config
- **`LogoWithConfig(config *BannerConfig) string`** - Returns logo with custom configuration
- **`PrintLogo()`** - Prints the logo to stdout
- **`PrintWelcome()`** - Shows welcome message with full banner
- **`RenderBanner(text string) string`** - Renders custom text as a styled banner

### Configuration Functions

- **`DefaultConfig() *BannerConfig`** - Returns sensible default configuration
- **`ConfigFromEnv() *BannerConfig`** - Creates configuration from environment variables
- **`GetBannerConfig(quiet, noBanner bool, style string) *BannerConfig`** - Creates config from CLI flags

## Banner Styles

### StyleFull
Full ASCII art banner with "GOLANG STARTER" text in multiple colors.

### StyleMinimal  
Compact logo version for space-constrained displays.

### StyleNone
No banner display (respects quiet mode and user preferences).

## Configuration Options

### Environment Variables

```bash
# Disable all banners
export GO_STARTER_BANNER=false

# Set banner style
export GO_STARTER_BANNER_STYLE=minimal  # full, minimal, none

# Disable colors (respects NO_COLOR standard)
export NO_COLOR=1
```

### CLI Flags

```bash
# Quiet mode (suppresses banners)
go-starter new --quiet

# Disable banner specifically
go-starter new --no-banner

# Set banner style
go-starter new --banner-style=minimal
```

### Programmatic Configuration

```go
// Custom banner configuration
config := &ascii.BannerConfig{
    Enabled:    true,
    Style:      ascii.StyleMinimal,
    Colors:     true,
    ShowOnHelp: true,
    Quiet:      false,
}

banner := ascii.BannerWithConfig(config)
```

## Features

- **Consistent Styling**: Uses lipgloss for unified visual experience across all commands
- **Environment Awareness**: Respects NO_COLOR standard and terminal capabilities
- **Configurable Display**: Multiple banner styles and display options
- **CLI Integration**: Seamless integration with command-line flags
- **Smart Defaults**: Sensible defaults with easy customization
- **Color Management**: Automatic color detection and fallback support

## Usage Examples

### Basic Usage

```go
import "github.com/francknouama/go-starter/internal/ascii"

// Print welcome banner with default config
ascii.PrintWelcome()

// Print logo only
ascii.PrintLogo()

// Get banner as string
banner := ascii.Banner()
fmt.Print(banner)
```

### Advanced Configuration

```go
// Create custom configuration
config := ascii.ConfigFromEnv()
config.Style = ascii.StyleMinimal
config.Colors = false

// Use custom configuration
ascii.PrintWelcomeWithConfig(config)

// Render custom banner
customBanner := ascii.RenderBannerWithConfig("My Project", config)
fmt.Print(customBanner)
```

### CLI Integration

```go
// In command implementation
bannerConfig := ascii.GetBannerConfig(quiet, noBanner, bannerStyle)

if bannerConfig.Enabled {
    ascii.PrintWelcomeWithConfig(bannerConfig)
}
```

## Command Integration

### Root Command (`go-starter --help`)
- Shows full banner in help text
- Respects environment configuration
- Uses consistent lipgloss styling

### Version Command (`go-starter version`)
- Shows logo with version information
- Maintains visual consistency
- Supports color configuration

### New Command (`go-starter new`)
- Full welcome banner in interactive mode
- Minimal banner in direct mode
- Configurable via CLI flags

### List Command (`go-starter list`)
- Minimal header banner
- Consistent with other commands
- Environment-aware display

## Dependencies

- **github.com/charmbracelet/lipgloss** - Terminal styling and colors
- **github.com/muesli/termenv** - Terminal environment detection