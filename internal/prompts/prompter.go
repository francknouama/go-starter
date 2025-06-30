package prompts

import (
	"os"

	"github.com/francknouama/go-starter/internal/prompts/bubbletea"
	"github.com/francknouama/go-starter/internal/prompts/interfaces"
	"github.com/francknouama/go-starter/internal/prompts/survey"
)

// PrompterFactory creates the appropriate prompter based on configuration
type PrompterFactory struct {
	useEnhancedUI bool
}

// NewPrompterFactory creates a new prompter factory
func NewPrompterFactory(useEnhancedUI bool) *PrompterFactory {
	return &PrompterFactory{
		useEnhancedUI: useEnhancedUI,
	}
}

// CreatePrompter creates a prompter instance based on the configuration
func (f *PrompterFactory) CreatePrompter() interfaces.Prompter {
	// Check if we should use enhanced UI (Bubble Tea)
	if f.useEnhancedUI && isTerminal() {
		// Try to create Bubble Tea prompter
		// If it fails, fall back to Survey
		if p := createBubbleTeaPrompter(); p != nil {
			return p
		}
	}
	
	// Default to Survey prompter
	return createSurveyPrompter()
}

// Helper function to check if we're in a terminal
func isTerminal() bool {
	stat, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

// createBubbleTeaPrompter creates a new Bubble Tea prompter
func createBubbleTeaPrompter() interfaces.Prompter {
	return bubbletea.NewPrompter()
}

// createSurveyPrompter creates a new Survey prompter
func createSurveyPrompter() interfaces.Prompter {
	return survey.NewPrompter()
}

// NewDefault creates a default prompter using Bubble Tea with Survey fallback
func NewDefault() interfaces.Prompter {
	factory := NewPrompterFactory(true)
	return factory.CreatePrompter()
}

// NewSurveyFallback creates a Survey prompter for fallback scenarios
func NewSurveyFallback() interfaces.Prompter {
	return survey.NewPrompter()
}