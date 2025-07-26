class GoStarter < Formula
  desc "Comprehensive Go project generator with modern best practices"
  homepage "https://github.com/francknouama/go-starter"
  version "1.4.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/francknouama/go-starter/releases/download/v1.4.0/go-starter-v1.4.0-darwin-amd64.tar.gz"
      sha256 "placeholder_sha256_for_macos_intel" # Will be updated by release automation

      def install
        bin.install "go-starter-v1.4.0-darwin-amd64" => "go-starter"
      end
    end

    if Hardware::CPU.arm?
      url "https://github.com/francknouama/go-starter/releases/download/v1.4.0/go-starter-v1.4.0-darwin-arm64.tar.gz"
      sha256 "placeholder_sha256_for_macos_arm64" # Will be updated by release automation

      def install
        bin.install "go-starter-v1.4.0-darwin-arm64" => "go-starter"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/francknouama/go-starter/releases/download/v1.4.0/go-starter-v1.4.0-linux-amd64.tar.gz"
      sha256 "placeholder_sha256_for_linux_intel" # Will be updated by release automation

      def install
        bin.install "go-starter-v1.4.0-linux-amd64" => "go-starter"
      end
    end

    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/francknouama/go-starter/releases/download/v1.4.0/go-starter-v1.4.0-linux-arm64.tar.gz"
      sha256 "placeholder_sha256_for_linux_arm64" # Will be updated by release automation

      def install
        bin.install "go-starter-v1.4.0-linux-arm64" => "go-starter"
      end
    end
  end

  test do
    # Test that the binary runs and shows version
    assert_match "go-starter version", shell_output("#{bin}/go-starter version")
    
    # Test that help command works
    assert_match "Comprehensive Go project generator", shell_output("#{bin}/go-starter --help")
    
    # Test list command
    assert_match "Available project types", shell_output("#{bin}/go-starter list")
    
    # Test dry-run generation (doesn't create files)
    assert_match "Files to be generated", shell_output("#{bin}/go-starter new test-project --type=cli --complexity=simple --dry-run --module=github.com/test/project")
  end

  def caveats
    <<~EOS
      go-starter has been installed! 🚀

      Quick start:
        go-starter new my-project --type=cli

      Interactive mode:
        go-starter new

      Advanced mode with all options:
        go-starter new --advanced

      For help and documentation:
        go-starter --help
        https://github.com/francknouama/go-starter
    EOS
  end
end