package grpcgateway

import (
	"fmt"
	"os"
)

// GRPCGatewayTestContext holds test state for gRPC Gateway scenarios
type GRPCGatewayTestContext struct {
	workDir            string
	projectName        string
	projectPath        string
	cmdOutput          string
	cmdError           error
	exitCode           int
	generatedFiles     []string
	auditRequirements  map[string]string
	serviceMeshConfig  map[string]string
	securityFeatures   []string
	complianceLevel    string
	performanceMetrics map[string]float64
}

// Helper methods
func (ctx *GRPCGatewayTestContext) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (ctx *GRPCGatewayTestContext) fileExistsError(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}
	return nil
}

func (ctx *GRPCGatewayTestContext) ensureFileExists(path string) error {
	if !ctx.fileExists(path) {
		return fmt.Errorf("required file does not exist: %s", path)
	}
	return nil
}