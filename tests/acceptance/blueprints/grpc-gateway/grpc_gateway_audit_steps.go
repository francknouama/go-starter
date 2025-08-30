package grpcgateway

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/cucumber/godog"
)

// RegisterProductionAuditSteps registers step definitions for production audit scenarios
func (ctx *GRPCGatewayTestContext) RegisterProductionAuditSteps(s *godog.ScenarioContext) {
	// Background steps
	s.Step(`^I have the project audit requirements$`, ctx.iHaveProjectAuditRequirements)
	
	// Health Check Protocol steps
	s.Step(`^I generate a grpc-gateway project with health checks enabled$`, ctx.iGenerateGRPCGatewayWithHealthChecks)
	s.Step(`^the service should implement gRPC health checking protocol$`, ctx.serviceShouldImplementHealthCheckingProtocol)
	s.Step(`^the health service should be registered with the gRPC server$`, ctx.healthServiceShouldBeRegistered)
	s.Step(`^the health endpoints should support both gRPC and REST protocols$`, ctx.healthEndpointsShouldSupportBothProtocols)
	s.Step(`^the health checks should verify database connectivity$`, ctx.healthChecksShouldVerifyDatabase)
	s.Step(`^the health checks should support Kubernetes (liveness|readiness) probes$`, ctx.healthChecksShouldSupportKubernetesProbes)
	
	// Distributed Tracing steps
	s.Step(`^I generate a grpc-gateway with OpenTelemetry tracing$`, ctx.iGenerateGRPCGatewayWithTracing)
	s.Step(`^the service should include gRPC tracing interceptors$`, ctx.serviceShouldIncludeGRPCTracingInterceptors)
	s.Step(`^the gateway should include HTTP tracing middleware$`, ctx.gatewayShouldIncludeHTTPTracingMiddleware)
	s.Step(`^the service should generate trace spans for all requests$`, ctx.serviceShouldGenerateTraceSpans)
	s.Step(`^the tracing should include request correlation IDs$`, ctx.tracingShouldIncludeCorrelationIDs)
	s.Step(`^the trace context should propagate between gRPC and REST layers$`, ctx.traceContextShouldPropagate)
	
	// Metrics Collection steps
	s.Step(`^I generate a grpc-gateway with Prometheus metrics$`, ctx.iGenerateGRPCGatewayWithMetrics)
	s.Step(`^the service should include gRPC metrics interceptors$`, ctx.serviceShouldIncludeGRPCMetricsInterceptors)
	s.Step(`^the gateway should include HTTP metrics middleware$`, ctx.gatewayShouldIncludeHTTPMetricsMiddleware)
	s.Step(`^the metrics should distinguish between gRPC and REST requests$`, ctx.metricsShouldDistinguishProtocols)
	s.Step(`^the service should export standard gRPC metrics$`, ctx.serviceShouldExportStandardGRPCMetrics)
	s.Step(`^the service should export custom business metrics$`, ctx.serviceShouldExportCustomMetrics)
	
	// Error Mapping steps
	s.Step(`^I generate a grpc-gateway with enhanced error mapping$`, ctx.iGenerateGRPCGatewayWithErrorMapping)
	s.Step(`^the service should map domain errors to appropriate gRPC codes$`, ctx.serviceShouldMapDomainErrors)
	s.Step(`^the gateway should translate gRPC errors to HTTP status codes$`, ctx.gatewayShouldTranslateErrors)
	s.Step(`^the error responses should be consistent across protocols$`, ctx.errorResponsesShouldBeConsistent)
	s.Step(`^the errors should include proper error context$`, ctx.errorsShouldIncludeContext)
	
	// Authentication & Authorization steps
	s.Step(`^I generate a grpc-gateway with RBAC authentication$`, ctx.iGenerateGRPCGatewayWithRBAC)
	s.Step(`^the service should include gRPC authentication interceptors$`, ctx.serviceShouldIncludeAuthInterceptors)
	s.Step(`^the gateway should include HTTP authentication middleware$`, ctx.gatewayShouldIncludeAuthMiddleware)
	s.Step(`^the service should support role-based access control$`, ctx.serviceShouldSupportRBAC)
	s.Step(`^the authentication should validate JWT tokens properly$`, ctx.authShouldValidateJWTTokens)
	s.Step(`^the service should support granular authorization rules$`, ctx.serviceShouldSupportGranularAuth)
	
	// Rate Limiting steps
	s.Step(`^I generate a grpc-gateway with intelligent rate limiting$`, ctx.iGenerateGRPCGatewayWithRateLimiting)
	s.Step(`^the service should include gRPC rate limiting interceptors$`, ctx.serviceShouldIncludeGRPCRateLimitingInterceptors)
	s.Step(`^the gateway should include HTTP rate limiting middleware$`, ctx.gatewayShouldIncludeHTTPRateLimitingMiddleware)
	s.Step(`^the rate limits should be configurable per method$`, ctx.rateLimitsShouldBeConfigurablePerMethod)
	s.Step(`^the rate limits should support different algorithms$`, ctx.rateLimitsShouldSupportAlgorithms)
	s.Step(`^the service should return appropriate rate limit headers$`, ctx.serviceShouldReturnRateLimitHeaders)
	
	// Configuration Management steps
	s.Step(`^I generate a grpc-gateway with advanced configuration$`, ctx.iGenerateGRPCGatewayWithAdvancedConfig)
	s.Step(`^the service should support environment-based configuration$`, ctx.serviceShouldSupportEnvironmentConfig)
	s.Step(`^the configuration should be validated at startup$`, ctx.configShouldBeValidatedAtStartup)
	s.Step(`^the service should support configuration hot-reloading$`, ctx.serviceShouldSupportConfigHotReloading)
	s.Step(`^the configuration should include secure secret management$`, ctx.configShouldIncludeSecretManagement)
	
	// Performance Optimization steps
	s.Step(`^I generate a grpc-gateway with performance optimizations$`, ctx.iGenerateGRPCGatewayWithPerformanceOptimizations)
	s.Step(`^the service should include connection pooling$`, ctx.serviceShouldIncludeConnectionPooling)
	s.Step(`^the gateway should support response compression$`, ctx.gatewayShouldSupportResponseCompression)
	s.Step(`^the service should optimize for HTTP/2 usage$`, ctx.serviceShouldOptimizeForHTTP2)
	s.Step(`^the service should include request buffering$`, ctx.serviceShouldIncludeRequestBuffering)
	s.Step(`^the performance should be measurable with benchmarks$`, ctx.performanceShouldBeMeasurable)
	
	// Container and Kubernetes steps
	s.Step(`^I generate a grpc-gateway with container support$`, ctx.iGenerateGRPCGatewayWithContainerSupport)
	s.Step(`^the project should include optimized Dockerfile$`, ctx.projectShouldIncludeOptimizedDockerfile)
	s.Step(`^the container should support multi-stage builds$`, ctx.containerShouldSupportMultiStageBuilds)
	s.Step(`^the service should include proper health check endpoints$`, ctx.serviceShouldIncludeHealthCheckEndpoints)
	s.Step(`^the container should follow security best practices$`, ctx.containerShouldFollowSecurityBestPractices)
	
	// Security Hardening steps
	s.Step(`^I generate a grpc-gateway with security hardening$`, ctx.iGenerateGRPCGatewayWithSecurityHardening)
	s.Step(`^the service should implement defense-in-depth$`, ctx.serviceShouldImplementDefenseInDepth)
	s.Step(`^all communications should be encrypted with TLS 1.3$`, ctx.allCommunicationsShouldBeEncryptedTLS13)
	s.Step(`^the service should validate all input rigorously$`, ctx.serviceShouldValidateAllInputRigorously)
	s.Step(`^the authentication should follow OWASP guidelines$`, ctx.authShouldFollowOWASPGuidelines)
	s.Step(`^the service should implement proper CORS policies$`, ctx.serviceShouldImplementCORSPolicies)
}

// Implementation methods for production audit steps

func (ctx *GRPCGatewayTestContext) iHaveProjectAuditRequirements() error {
	ctx.auditRequirements = map[string]string{
		"health_checks":        "gRPC health protocol required",
		"distributed_tracing":  "OpenTelemetry integration required",
		"metrics_collection":   "Prometheus metrics required",
		"error_mapping":        "Consistent error handling required",
		"authentication":       "JWT + RBAC required",
		"rate_limiting":        "Intelligent rate limiting required",
		"configuration":        "Hot-reload configuration required",
		"performance":          "Optimized for production load",
		"container":            "Multi-stage Docker builds",
		"security":             "Defense-in-depth security",
	}
	return nil
}

func (ctx *GRPCGatewayTestContext) iGenerateGRPCGatewayWithHealthChecks() error {
	return ctx.generateGRPCGatewayProject(map[string]string{
		"health-checks": "enabled",
		"observability": "enabled",
	})
}

func (ctx *GRPCGatewayTestContext) serviceShouldImplementHealthCheckingProtocol() error {
	// Check for gRPC health service implementation
	healthServiceFile := filepath.Join(ctx.projectPath, "internal/services/health.go")
	if !ctx.fileExists(healthServiceFile) {
		return fmt.Errorf("health service file not found: %s", healthServiceFile)
	}
	
	// Verify health service implements grpc_health_v1.HealthServer
	content, err := os.ReadFile(healthServiceFile)
	if err != nil {
		return fmt.Errorf("failed to read health service file: %v", err)
	}
	
	if !strings.Contains(string(content), "grpc_health_v1.HealthServer") {
		return fmt.Errorf("health service does not implement gRPC health protocol")
	}
	
	return nil
}

func (ctx *GRPCGatewayTestContext) healthServiceShouldBeRegistered() error {
	// Check if health service is registered in gRPC server
	grpcServerFile := filepath.Join(ctx.projectPath, "internal/server/grpc.go")
	content, err := os.ReadFile(grpcServerFile)
	if err != nil {
		return fmt.Errorf("failed to read gRPC server file: %v", err)
	}
	
	if !strings.Contains(string(content), "RegisterHealthServer") {
		return fmt.Errorf("health service is not registered with gRPC server")
	}
	
	return nil
}

func (ctx *GRPCGatewayTestContext) healthEndpointsShouldSupportBothProtocols() error {
	// Check for gRPC health proto definition
	healthProto := filepath.Join(ctx.projectPath, "proto/health/v1/health.proto")
	if !ctx.fileExists(healthProto) {
		return fmt.Errorf("health proto file not found")
	}
	
	// Check for REST health endpoints in gateway
	gatewayFile := filepath.Join(ctx.projectPath, "internal/server/gateway.go")
	content, err := os.ReadFile(gatewayFile)
	if err != nil {
		return fmt.Errorf("failed to read gateway file: %v", err)
	}
	
	if !strings.Contains(string(content), "/health") {
		return fmt.Errorf("REST health endpoints not found in gateway")
	}
	
	return nil
}

func (ctx *GRPCGatewayTestContext) healthChecksShouldVerifyDatabase() error {
	// Check if health service includes database connectivity check
	healthServiceFile := filepath.Join(ctx.projectPath, "internal/services/health.go")
	content, err := os.ReadFile(healthServiceFile)
	if err != nil {
		return fmt.Errorf("failed to read health service file: %v", err)
	}
	
	if !strings.Contains(string(content), "database") && !strings.Contains(string(content), "db") {
		return fmt.Errorf("health service does not include database connectivity check")
	}
	
	return nil
}

func (ctx *GRPCGatewayTestContext) healthChecksShouldSupportKubernetesProbes(probeType string) error {
	// Check for appropriate health check endpoints for Kubernetes probes
	expectedEndpoints := map[string]string{
		"liveness":  "/health/live",
		"readiness": "/health/ready",
	}
	
	endpoint, exists := expectedEndpoints[probeType]
	if !exists {
		return fmt.Errorf("unknown probe type: %s", probeType)
	}
	
	// Verify the endpoint is implemented
	gatewayFile := filepath.Join(ctx.projectPath, "internal/server/gateway.go")
	content, err := os.ReadFile(gatewayFile)
	if err != nil {
		return fmt.Errorf("failed to read gateway file: %v", err)
	}
	
	if !strings.Contains(string(content), endpoint) {
		return fmt.Errorf("Kubernetes %s probe endpoint not found: %s", probeType, endpoint)
	}
	
	return nil
}

func (ctx *GRPCGatewayTestContext) iGenerateGRPCGatewayWithTracing() error {
	return ctx.generateGRPCGatewayProject(map[string]string{
		"observability": "enabled",
		"tracing":       "opentelemetry",
	})
}

func (ctx *GRPCGatewayTestContext) serviceShouldIncludeGRPCTracingInterceptors() error {
	// Check for tracing interceptors in gRPC server
	interceptorsFile := filepath.Join(ctx.projectPath, "internal/interceptors/enhanced.go")
	if !ctx.fileExists(interceptorsFile) {
		return fmt.Errorf("interceptors file not found")
	}
	
	content, err := os.ReadFile(interceptorsFile)
	if err != nil {
		return fmt.Errorf("failed to read interceptors file: %v", err)
	}
	
	if !strings.Contains(string(content), "tracing") && !strings.Contains(string(content), "opentelemetry") {
		return fmt.Errorf("gRPC tracing interceptors not found")
	}
	
	return nil
}

func (ctx *GRPCGatewayTestContext) gatewayShouldIncludeHTTPTracingMiddleware() error {
	// Check for tracing middleware in HTTP gateway
	middlewareDir := filepath.Join(ctx.projectPath, "internal/middleware")
	files, err := os.ReadDir(middlewareDir)
	if err != nil {
		return fmt.Errorf("failed to read middleware directory: %v", err)
	}
	
	tracingFound := false
	for _, file := range files {
		if strings.Contains(file.Name(), "tracing") || strings.Contains(file.Name(), "request_id") {
			tracingFound = true
			break
		}
	}
	
	if !tracingFound {
		return fmt.Errorf("HTTP tracing middleware not found")
	}
	
	return nil
}

func (ctx *GRPCGatewayTestContext) serviceShouldGenerateTraceSpans() error {
	// Check if service generates trace spans
	return ctx.verifyImplementationPattern("trace spans", []string{
		"internal/interceptors/enhanced.go",
		"internal/middleware/",
	}, []string{"span", "trace", "opentelemetry"})
}

func (ctx *GRPCGatewayTestContext) tracingShouldIncludeCorrelationIDs() error {
	// Check for correlation ID implementation
	return ctx.verifyImplementationPattern("correlation IDs", []string{
		"internal/middleware/request_id.go",
	}, []string{"correlation", "request_id", "trace_id"})
}

func (ctx *GRPCGatewayTestContext) traceContextShouldPropagate() error {
	// Check for trace context propagation between layers
	return ctx.verifyImplementationPattern("trace context propagation", []string{
		"internal/interceptors/enhanced.go",
		"internal/middleware/",
	}, []string{"propagate", "context", "metadata"})
}

// Helper methods

func (ctx *GRPCGatewayTestContext) generateGRPCGatewayProject(options map[string]string) error {
	ctx.projectName = "test-grpc-gateway-audit"
	ctx.projectPath = filepath.Join(ctx.workDir, ctx.projectName)
	
	args := []string{
		"new", ctx.projectName,
		"--type=grpc-gateway",
		"--logger=slog",
		"--output=" + ctx.workDir,
		"--module=github.com/test/" + ctx.projectName,
		"--no-git",
		"--quiet",
	}
	
	// Add optional features
	for key, value := range options {
		if value == "enabled" {
			args = append(args, "--"+key)
		} else {
			args = append(args, "--"+key+"="+value)
		}
	}
	
	cmd := exec.Command("go-starter", args...)
	output, err := cmd.CombinedOutput()
	ctx.cmdError = err
	ctx.cmdOutput = string(output)
	
	if ctx.cmdError != nil {
		return fmt.Errorf("failed to generate gRPC Gateway project: %v\nOutput: %s", ctx.cmdError, ctx.cmdOutput)
	}
	
	return nil
}

// fileExists helper method removed - using shared method from main test file

func (ctx *GRPCGatewayTestContext) verifyImplementationPattern(feature string, files []string, patterns []string) error {
	for _, file := range files {
		path := filepath.Join(ctx.projectPath, file)
		
		// Handle directory patterns
		if strings.HasSuffix(file, "/") {
			dirFiles, err := filepath.Glob(path + "*.go")
			if err != nil {
				continue
			}
			
			for _, dirFile := range dirFiles {
				if ctx.checkPatternsInFile(dirFile, patterns) {
					return nil
				}
			}
		} else {
			if ctx.fileExists(path) && ctx.checkPatternsInFile(path, patterns) {
				return nil
			}
		}
	}
	
	return fmt.Errorf("%s implementation not found", feature)
}

func (ctx *GRPCGatewayTestContext) checkPatternsInFile(filePath string, patterns []string) bool {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	
	contentStr := strings.ToLower(string(content))
	
	for _, pattern := range patterns {
		if strings.Contains(contentStr, strings.ToLower(pattern)) {
			return true
		}
	}
	
	return false
}

// Placeholder implementations for remaining steps (to be implemented as needed)
func (ctx *GRPCGatewayTestContext) iGenerateGRPCGatewayWithMetrics() error {
	return ctx.generateGRPCGatewayProject(map[string]string{"metrics": "prometheus"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldIncludeGRPCMetricsInterceptors() error {
	return ctx.verifyImplementationPattern("gRPC metrics interceptors", 
		[]string{"internal/interceptors/enhanced.go"}, 
		[]string{"metrics", "prometheus"})
}

func (ctx *GRPCGatewayTestContext) gatewayShouldIncludeHTTPMetricsMiddleware() error {
	return ctx.verifyImplementationPattern("HTTP metrics middleware", 
		[]string{"internal/middleware/"}, 
		[]string{"metrics", "prometheus"})
}

func (ctx *GRPCGatewayTestContext) metricsShouldDistinguishProtocols() error {
	return ctx.verifyImplementationPattern("protocol-specific metrics", 
		[]string{"internal/metrics/collector.go"}, 
		[]string{"grpc", "http", "protocol"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldExportStandardGRPCMetrics() error {
	return ctx.verifyImplementationPattern("standard gRPC metrics", 
		[]string{"internal/metrics/collector.go"}, 
		[]string{"grpc_requests", "grpc_duration"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldExportCustomMetrics() error {
	return ctx.verifyImplementationPattern("custom business metrics", 
		[]string{"internal/metrics/collector.go"}, 
		[]string{"custom", "business"})
}

func (ctx *GRPCGatewayTestContext) iGenerateGRPCGatewayWithErrorMapping() error {
	return ctx.generateGRPCGatewayProject(map[string]string{"error-mapping": "enhanced"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldMapDomainErrors() error {
	return ctx.verifyImplementationPattern("domain error mapping", 
		[]string{"internal/errors/"}, 
		[]string{"domain", "mapping", "grpc.status"})
}

func (ctx *GRPCGatewayTestContext) gatewayShouldTranslateErrors() error {
	return ctx.verifyImplementationPattern("error translation", 
		[]string{"internal/middleware/error_handler.go"}, 
		[]string{"translate", "http.status", "grpc.codes"})
}

func (ctx *GRPCGatewayTestContext) errorResponsesShouldBeConsistent() error {
	return ctx.verifyImplementationPattern("consistent error responses", 
		[]string{"internal/errors/"}, 
		[]string{"consistent", "format"})
}

func (ctx *GRPCGatewayTestContext) errorsShouldIncludeContext() error {
	return ctx.verifyImplementationPattern("error context", 
		[]string{"internal/errors/"}, 
		[]string{"context", "details"})
}

// Add similar placeholder implementations for all other step methods
// This provides a comprehensive foundation that can be extended

// Additional utility methods
func (ctx *GRPCGatewayTestContext) iGenerateGRPCGatewayWithRBAC() error {
	return ctx.generateGRPCGatewayProject(map[string]string{"auth-type": "jwt", "rbac": "enabled"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldIncludeAuthInterceptors() error {
	return ctx.verifyImplementationPattern("auth interceptors", 
		[]string{"internal/interceptors/enhanced.go"}, 
		[]string{"auth", "jwt", "interceptor"})
}

func (ctx *GRPCGatewayTestContext) gatewayShouldIncludeAuthMiddleware() error {
	return ctx.verifyImplementationPattern("auth middleware", 
		[]string{"internal/middleware/auth.go"}, 
		[]string{"jwt", "auth", "middleware"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldSupportRBAC() error {
	return ctx.verifyImplementationPattern("RBAC support", 
		[]string{"internal/middleware/auth.go", "internal/services/auth.go"}, 
		[]string{"rbac", "role", "permission"})
}

func (ctx *GRPCGatewayTestContext) authShouldValidateJWTTokens() error {
	return ctx.verifyImplementationPattern("JWT validation", 
		[]string{"internal/services/auth.go"}, 
		[]string{"jwt", "validate", "token"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldSupportGranularAuth() error {
	return ctx.verifyImplementationPattern("granular authorization", 
		[]string{"internal/middleware/auth.go"}, 
		[]string{"granular", "method", "resource"})
}

// Continue with remaining placeholder methods...
func (ctx *GRPCGatewayTestContext) iGenerateGRPCGatewayWithRateLimiting() error {
	return ctx.generateGRPCGatewayProject(map[string]string{"rate-limiting": "intelligent"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldIncludeGRPCRateLimitingInterceptors() error {
	return ctx.verifyImplementationPattern("gRPC rate limiting", 
		[]string{"internal/interceptors/enhanced.go"}, 
		[]string{"rate", "limit", "grpc"})
}

func (ctx *GRPCGatewayTestContext) gatewayShouldIncludeHTTPRateLimitingMiddleware() error {
	return ctx.verifyImplementationPattern("HTTP rate limiting", 
		[]string{"internal/middleware/rate_limiter.go"}, 
		[]string{"rate", "limit", "http"})
}

func (ctx *GRPCGatewayTestContext) rateLimitsShouldBeConfigurablePerMethod() error {
	return ctx.verifyImplementationPattern("per-method rate limits", 
		[]string{"internal/middleware/rate_limiter.go"}, 
		[]string{"method", "configurable"})
}

func (ctx *GRPCGatewayTestContext) rateLimitsShouldSupportAlgorithms() error {
	return ctx.verifyImplementationPattern("rate limiting algorithms", 
		[]string{"internal/middleware/rate_limiter.go"}, 
		[]string{"token_bucket", "sliding_window", "algorithm"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldReturnRateLimitHeaders() error {
	return ctx.verifyImplementationPattern("rate limit headers", 
		[]string{"internal/middleware/rate_limiter.go"}, 
		[]string{"header", "X-RateLimit", "Retry-After"})
}

// Additional methods following the same pattern...
// Each method should implement appropriate verification logic

func (ctx *GRPCGatewayTestContext) iGenerateGRPCGatewayWithAdvancedConfig() error {
	return ctx.generateGRPCGatewayProject(map[string]string{"config": "advanced"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldSupportEnvironmentConfig() error {
	return ctx.verifyImplementationPattern("environment config", 
		[]string{"internal/config/config.go"}, 
		[]string{"environment", "env", "config"})
}

func (ctx *GRPCGatewayTestContext) configShouldBeValidatedAtStartup() error {
	return ctx.verifyImplementationPattern("config validation", 
		[]string{"internal/config/config.go"}, 
		[]string{"validate", "startup"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldSupportConfigHotReloading() error {
	return ctx.verifyImplementationPattern("config hot reloading", 
		[]string{"internal/config/config.go"}, 
		[]string{"reload", "watch", "hot"})
}

func (ctx *GRPCGatewayTestContext) configShouldIncludeSecretManagement() error {
	return ctx.verifyImplementationPattern("secret management", 
		[]string{"internal/config/config.go"}, 
		[]string{"secret", "vault", "encrypted"})
}

func (ctx *GRPCGatewayTestContext) iGenerateGRPCGatewayWithPerformanceOptimizations() error {
	return ctx.generateGRPCGatewayProject(map[string]string{"performance": "optimized"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldIncludeConnectionPooling() error {
	return ctx.verifyImplementationPattern("connection pooling", 
		[]string{"internal/server/grpc.go", "internal/server/gateway.go"}, 
		[]string{"pool", "connection"})
}

func (ctx *GRPCGatewayTestContext) gatewayShouldSupportResponseCompression() error {
	return ctx.verifyImplementationPattern("response compression", 
		[]string{"internal/server/gateway.go"}, 
		[]string{"compression", "gzip"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldOptimizeForHTTP2() error {
	return ctx.verifyImplementationPattern("HTTP/2 optimization", 
		[]string{"internal/server/gateway.go"}, 
		[]string{"http2", "h2"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldIncludeRequestBuffering() error {
	return ctx.verifyImplementationPattern("request buffering", 
		[]string{"internal/server/gateway.go"}, 
		[]string{"buffer", "batching"})
}

func (ctx *GRPCGatewayTestContext) performanceShouldBeMeasurable() error {
	return ctx.verifyImplementationPattern("performance benchmarks", 
		[]string{"tests/"}, 
		[]string{"benchmark", "performance", "load"})
}

func (ctx *GRPCGatewayTestContext) iGenerateGRPCGatewayWithContainerSupport() error {
	return ctx.generateGRPCGatewayProject(map[string]string{"container": "optimized"})
}

func (ctx *GRPCGatewayTestContext) projectShouldIncludeOptimizedDockerfile() error {
	dockerFile := filepath.Join(ctx.projectPath, "Dockerfile")
	return ctx.verifyFileExists("optimized Dockerfile", dockerFile)
}

func (ctx *GRPCGatewayTestContext) containerShouldSupportMultiStageBuilds() error {
	_ = filepath.Join(ctx.projectPath, "Dockerfile")
	return ctx.verifyImplementationPattern("multi-stage builds", 
		[]string{"Dockerfile"}, 
		[]string{"FROM", "AS", "stage"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldIncludeHealthCheckEndpoints() error {
	return ctx.verifyImplementationPattern("health check endpoints", 
		[]string{"internal/services/health.go"}, 
		[]string{"health", "liveness", "readiness"})
}

func (ctx *GRPCGatewayTestContext) containerShouldFollowSecurityBestPractices() error {
	return ctx.verifyImplementationPattern("container security", 
		[]string{"Dockerfile"}, 
		[]string{"USER", "nonroot", "readonly"})
}

func (ctx *GRPCGatewayTestContext) iGenerateGRPCGatewayWithSecurityHardening() error {
	return ctx.generateGRPCGatewayProject(map[string]string{"security": "hardened"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldImplementDefenseInDepth() error {
	return ctx.verifyImplementationPattern("defense in depth", 
		[]string{"internal/middleware/"}, 
		[]string{"security", "defense", "hardening"})
}

func (ctx *GRPCGatewayTestContext) allCommunicationsShouldBeEncryptedTLS13() error {
	return ctx.verifyImplementationPattern("TLS 1.3 encryption", 
		[]string{"internal/tls/config.go"}, 
		[]string{"tls", "1.3", "encryption"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldValidateAllInputRigorously() error {
	return ctx.verifyImplementationPattern("input validation", 
		[]string{"internal/middleware/"}, 
		[]string{"validate", "sanitize", "input"})
}

func (ctx *GRPCGatewayTestContext) authShouldFollowOWASPGuidelines() error {
	return ctx.verifyImplementationPattern("OWASP compliance", 
		[]string{"internal/middleware/auth.go"}, 
		[]string{"owasp", "secure", "guidelines"})
}

func (ctx *GRPCGatewayTestContext) serviceShouldImplementCORSPolicies() error {
	return ctx.verifyImplementationPattern("CORS policies", 
		[]string{"internal/middleware/security.go"}, 
		[]string{"cors", "origin", "policy"})
}

// Helper method for file existence verification
func (ctx *GRPCGatewayTestContext) verifyFileExists(description string, filePath string) error {
	if !ctx.fileExists(filePath) {
		return fmt.Errorf("%s not found: %s", description, filePath)
	}
	return nil
}