package validation

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
)

// TestFeatures runs the Cross-System Validation BDD tests using godog
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			ctx := NewCrossSystemValidationTestContext()
			
			s.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
				// Setup before each scenario
				return ctx, nil
			})
			
			s.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				// Cleanup after each scenario
				return ctx, nil
			})
			
			ctx.RegisterSteps(s)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},  
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run cross-system validation feature tests")
	}
}