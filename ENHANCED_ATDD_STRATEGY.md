# Enhanced ATDD Strategy for Template Quality

## Problem Statement

Current ATDD tests focus on functional correctness but miss template generation quality issues:
- Unused imports (`"fmt"`, `"os"`, `models` package in raw SQL)  
- Unused variables (`dsn`, `errorHandler`)
- Configuration mismatches (postgres vs postgresql)
- Framework inconsistencies (Gin imports in Fiber projects)
- Missing dependencies (bcrypt, golang.org/x/crypto)
- Conditional import inconsistencies (models imported when ORM not used)

## Enhanced Coverage Strategy

### 1. Static Code Quality Validation

```go
func TestGeneratedCodeQuality(t *testing.T) {
    configs := []TestConfig{
        {Framework: "gin", Database: "postgresql", ORM: "", Logger: "slog"},
        {Framework: "fiber", Database: "mysql", ORM: "gorm", Logger: "zap"},
        {Framework: "echo", Database: "sqlite", ORM: "", Logger: "logrus"},
        // ... all permutations
    }
    
    for _, config := range configs {
        t.Run(config.Name(), func(t *testing.T) {
            // Generate project
            projectPath := generateProject(t, config)
            
            // 1. Compilation check
            assertCompilationSuccess(t, projectPath)
            
            // 2. Linting validation
            assertNoLintingErrors(t, projectPath)
            
            // 3. Import analysis
            assertNoUnusedImports(t, projectPath)
            
            // 4. Variable analysis  
            assertNoUnusedVariables(t, projectPath)
            
            // 5. Configuration consistency
            assertConfigurationConsistency(t, projectPath, config)
        })
    }
}
```

### 2. Template Logic Validation

```go
func TestTemplateConditionalLogic(t *testing.T) {
    testCases := []struct {
        name     string
        config   TestConfig
        validate func(t *testing.T, projectPath string, config TestConfig)
    }{
        {
            name: "PostgreSQL selection generates PostgreSQL code",
            config: TestConfig{Database: "postgresql"},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                // Check database connection file
                connectionFile := filepath.Join(projectPath, "internal/database/connection.go")
                content := readFile(t, connectionFile)
                
                assert.Contains(t, content, `_ "github.com/lib/pq"`)
                assert.Contains(t, content, `driverName = "postgres"`)
                assert.NotContains(t, content, `_ "github.com/mattn/go-sqlite3"`)
            },
        },
        {
            name: "Fiber framework generates only Fiber imports",
            config: TestConfig{Framework: "fiber"},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                // Check all Go files for framework imports
                goFiles := findGoFiles(projectPath)
                for _, file := range goFiles {
                    content := readFile(t, file)
                    assert.NotContains(t, content, `"github.com/gin-gonic/gin"`)
                    if strings.Contains(content, "fiber") {
                        assert.Contains(t, content, `"github.com/gofiber/fiber/v2"`)
                    }
                }
            },
        },
        {
            name: "Raw SQL database connection does not import models package",
            config: TestConfig{DatabaseORM: ""},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                // Check database connection file
                connectionFile := filepath.Join(projectPath, "internal/database/connection.go")
                content := readFile(t, connectionFile)
                
                // Should not import models when using raw SQL
                assert.NotContains(t, content, `"internal/models"`)
                assert.NotContains(t, content, `/internal/models"`)
                
                // But should still have database/sql import
                assert.Contains(t, content, `"database/sql"`)
            },
        },
        {
            name: "GORM database connection imports models package",
            config: TestConfig{DatabaseORM: "gorm"},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                // Check database connection file
                connectionFile := filepath.Join(projectPath, "internal/database/connection.go")
                content := readFile(t, connectionFile)
                
                // Should import models when using GORM
                assert.Contains(t, content, `/internal/models"`)
                
                // Should use models in AutoMigrate
                assert.Contains(t, content, "models.User{}")
                assert.Contains(t, content, "AutoMigrate")
                
                // Should have GORM imports instead of database/sql
                assert.Contains(t, content, `"gorm.io/gorm"`)
                assert.NotContains(t, content, `"database/sql"`)
            },
        },
        {
            name: "Authentication enabled includes required middleware imports",
            config: TestConfig{AuthType: "jwt", Framework: "gin"},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                // Check middleware files have correct imports
                authMiddleware := filepath.Join(projectPath, "internal/middleware/auth.go")
                content := readFile(t, authMiddleware)
                
                // Should have JWT library when auth is JWT
                assert.Contains(t, content, `"github.com/golang-jwt/jwt/v5"`)
                
                // Should have framework-specific imports
                assert.Contains(t, content, `"github.com/gin-gonic/gin"`)
                assert.NotContains(t, content, `"github.com/gofiber/fiber/v2"`)
                
                // Should not have unused crypto imports
                if !strings.Contains(content, "bcrypt") {
                    assert.NotContains(t, content, `"golang.org/x/crypto/bcrypt"`)
                }
            },
        },
        {
            name: "No authentication means no auth-related imports",
            config: TestConfig{AuthType: ""},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                // Auth middleware file should not exist
                authMiddleware := filepath.Join(projectPath, "internal/middleware/auth.go")
                assert.NoFileExists(t, authMiddleware)
                
                // Main server file should not import auth packages
                mainFile := filepath.Join(projectPath, "cmd/server/main.go")
                content := readFile(t, mainFile)
                assert.NotContains(t, content, `"github.com/golang-jwt/jwt/v5"`)
                assert.NotContains(t, content, `"golang.org/x/crypto"`)
                
                // go.mod should not have auth dependencies
                goMod := readFile(t, filepath.Join(projectPath, "go.mod"))
                assert.NotContains(t, goMod, "golang.org/x/crypto")
                assert.NotContains(t, goMod, "github.com/golang-jwt/jwt")
            },
        },
        {
            name: "Logger selection generates consistent imports across files",
            config: TestConfig{Logger: "zap"},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                goFiles := findGoFiles(projectPath)
                
                for _, file := range goFiles {
                    content := readFile(t, file)
                    
                    // If file uses logging, should use selected logger
                    if strings.Contains(content, "logger") && strings.Contains(content, "import") {
                        if strings.Contains(content, "zap") {
                            assert.Contains(t, content, `"go.uber.org/zap"`)
                            assert.NotContains(t, content, `"github.com/sirupsen/logrus"`)
                            assert.NotContains(t, content, `"github.com/rs/zerolog"`)
                        }
                    }
                }
            },
        },
        {
            name: "Database migrations match selected database driver",
            config: TestConfig{DatabaseDriver: "postgresql", DatabaseORM: ""},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                // Check migration files use PostgreSQL syntax
                migrationFiles := findSQLFiles(projectPath, "migrations")
                
                for _, file := range migrationFiles {
                    content := readFile(t, file)
                    
                    // Should use PostgreSQL-specific syntax
                    if strings.Contains(content, "CREATE TABLE") {
                        assert.Contains(t, content, "SERIAL PRIMARY KEY")
                        assert.Contains(t, content, "VARCHAR")
                        assert.NotContains(t, content, "AUTO_INCREMENT") // MySQL
                        assert.NotContains(t, content, "INTEGER PRIMARY KEY AUTOINCREMENT") // SQLite
                    }
                }
            },
        },
        {
            name: "Docker compose services match database selection",
            config: TestConfig{DatabaseDriver: "postgresql"},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                dockerCompose := filepath.Join(projectPath, "docker-compose.yml")
                if fileExists(dockerCompose) {
                    content := readFile(t, dockerCompose)
                    
                    // Should have PostgreSQL service
                    assert.Contains(t, content, "postgres:")
                    assert.Contains(t, content, "image: postgres:")
                    
                    // Should not have other database services
                    assert.NotContains(t, content, "image: mysql:")
                    assert.NotContains(t, content, "image: redis:")
                }
            },
        },
        {
            name: "Configuration files match selected options",
            config: TestConfig{DatabaseDriver: "postgresql", Logger: "zap", Framework: "gin"},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                configFiles := findYAMLFiles(projectPath, "configs")
                
                for _, file := range configFiles {
                    content := readFile(t, file)
                    
                    // Database configuration should match selection
                    if strings.Contains(content, "database") {
                        assert.Contains(t, content, "postgres")
                        assert.NotContains(t, content, "mysql")
                        assert.NotContains(t, content, "sqlite")
                    }
                    
                    // Logger configuration should match selection
                    if strings.Contains(content, "logger") || strings.Contains(content, "log") {
                        assert.Contains(t, content, "zap")
                        assert.NotContains(t, content, "logrus")
                    }
                }
            },
        },
        {
            name: "Error handling imports are present when needed",
            config: TestConfig{Framework: "gin"},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                // Check error handling files
                errorFiles := findGoFiles(projectPath)
                
                for _, file := range errorFiles {
                    if strings.Contains(file, "error") || strings.Contains(file, "handler") {
                        content := readFile(t, file)
                        
                        // Should have fmt for error formatting
                        if strings.Contains(content, "fmt.Errorf") {
                            assert.Contains(t, content, `"fmt"`)
                        }
                        
                        // Should not have unused error packages
                        if !strings.Contains(content, "errors.") {
                            assert.NotContains(t, content, `"errors"`)
                        }
                    }
                }
            },
        },
        {
            name: "Repository layer matches ORM selection",
            config: TestConfig{DatabaseORM: "gorm"},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                repoFiles := findGoFiles(filepath.Join(projectPath, "internal/repository"))
                
                for _, file := range repoFiles {
                    content := readFile(t, file)
                    
                    // Should use GORM patterns
                    assert.Contains(t, content, `"gorm.io/gorm"`)
                    assert.Contains(t, content, ".Create(")
                    assert.Contains(t, content, ".Find(")
                    
                    // Should not use raw SQL patterns
                    assert.NotContains(t, content, `"database/sql"`)
                    assert.NotContains(t, content, ".Exec(")
                    assert.NotContains(t, content, ".Query(")
                }
            },
        },
        {
            name: "Test files match project configuration",
            config: TestConfig{DatabaseDriver: "postgresql", Framework: "gin"},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                testFiles := findGoFiles(filepath.Join(projectPath, "tests"))
                
                for _, file := range testFiles {
                    content := readFile(t, file)
                    
                    // Integration tests should use selected database
                    if strings.Contains(content, "testcontainers") {
                        assert.Contains(t, content, "postgres")
                        assert.NotContains(t, content, "mysql")
                    }
                    
                    // HTTP tests should use selected framework
                    if strings.Contains(content, "http") {
                        if strings.Contains(content, "gin") {
                            assert.Contains(t, content, `"github.com/gin-gonic/gin"`)
                        }
                    }
                }
            },
        },
        {
            name: "Security middleware matches framework",
            config: TestConfig{Framework: "fiber"},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                securityFile := filepath.Join(projectPath, "internal/middleware/security_headers.go")
                if fileExists(securityFile) {
                    content := readFile(t, securityFile)
                    
                    // Should use Fiber-specific context
                    assert.Contains(t, content, `"github.com/gofiber/fiber/v2"`)
                    assert.Contains(t, content, "*fiber.Ctx")
                    
                    // Should not use other framework contexts
                    assert.NotContains(t, content, "*gin.Context")
                    assert.NotContains(t, content, "echo.Context")
                }
            },
        },
        {
            name: "Environment files match configuration",
            config: TestConfig{DatabaseDriver: "postgresql", AuthType: "jwt"},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                envFiles := [...]string{".env.example", ".env.dev", ".env.prod"}
                
                for _, envFile := range envFiles {
                    envPath := filepath.Join(projectPath, envFile)
                    if fileExists(envPath) {
                        content := readFile(t, envPath)
                        
                        // Should have database-specific variables
                        assert.Contains(t, content, "DB_HOST")
                        assert.Contains(t, content, "POSTGRES_")
                        
                        // Should have auth-specific variables when auth enabled
                        assert.Contains(t, content, "JWT_SECRET")
                        
                        // Should not have unused database variables
                        assert.NotContains(t, content, "MYSQL_")
                        assert.NotContains(t, content, "REDIS_")
                    }
                }
            },
        },
        {
            name: "Makefile commands match project type",
            config: TestConfig{DatabaseDriver: "postgresql"},
            validate: func(t *testing.T, projectPath string, config TestConfig) {
                makefile := filepath.Join(projectPath, "Makefile")
                if fileExists(makefile) {
                    content := readFile(t, makefile)
                    
                    // Should have database-specific commands
                    assert.Contains(t, content, "migrate")
                    assert.Contains(t, content, "postgres")
                    
                    // Docker commands should reference correct services
                    if strings.Contains(content, "docker-compose") {
                        assert.Contains(t, content, "postgres")
                        assert.NotContains(t, content, "mysql")
                    }
                }
            },
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            projectPath := generateProject(t, tc.config)
            tc.validate(t, projectPath, tc.config)
        })
    }
}
```

### 3. Configuration Matrix Testing

```go
func TestConfigurationMatrix(t *testing.T) {
    frameworks := []string{"gin", "echo", "fiber", "chi"}
    databases := []string{"postgresql", "mysql", "sqlite"}
    orms := []string{"", "gorm"}
    loggers := []string{"slog", "zap", "logrus", "zerolog"}
    
    // Test critical combinations
    criticalCombinations := []TestConfig{
        {Framework: "gin", Database: "postgresql", ORM: "", Logger: "slog"},
        {Framework: "fiber", Database: "mysql", ORM: "gorm", Logger: "zap"},
        {Framework: "echo", Database: "sqlite", ORM: "", Logger: "logrus"},
        // Add more strategic combinations
    }
    
    for _, config := range criticalCombinations {
        t.Run(config.Name(), func(t *testing.T) {
            projectPath := generateProject(t, config)
            
            // Validate each aspect matches configuration
            validateFrameworkConsistency(t, projectPath, config.Framework)
            validateDatabaseConsistency(t, projectPath, config.Database, config.ORM)
            validateLoggerConsistency(t, projectPath, config.Logger)
            
            // Ensure no cross-contamination
            assertNoIncorrectImports(t, projectPath, config)
        })
    }
}
```

### 4. Static Analysis Integration

```go
func assertNoLintingErrors(t *testing.T, projectPath string) {
    // Run golangci-lint
    cmd := exec.Command("golangci-lint", "run", "--timeout=5m", "./...")
    cmd.Dir = projectPath
    
    output, err := cmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Linting failed: %s\nOutput: %s", err, string(output))
    }
    
    // Should have no warnings/errors
    assert.Empty(t, strings.TrimSpace(string(output)), "Generated code should pass linting")
}

func assertNoUnusedImports(t *testing.T, projectPath string) {
    // Use go/ast to parse all Go files and check for unused imports
    err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
        if !strings.HasSuffix(path, ".go") || strings.Contains(path, "vendor/") {
            return nil
        }
        
        fset := token.NewFileSet()
        node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
        if err != nil {
            return err
        }
        
        // Check for unused imports using go/types
        unusedImports := findUnusedImports(node, fset)
        assert.Empty(t, unusedImports, "File %s has unused imports: %v", path, unusedImports)
        
        return nil
    })
    assert.NoError(t, err)
}
```

### 5. Dependency Validation

```go
func TestDependencyConsistency(t *testing.T) {
    testCases := []struct {
        name     string
        config   TestConfig
        validate func(t *testing.T, projectPath string)
    }{
        {
            name: "Auth enabled includes bcrypt dependency",
            config: TestConfig{AuthType: "jwt"},
            validate: func(t *testing.T, projectPath string) {
                goMod := readFile(t, filepath.Join(projectPath, "go.mod"))
                assert.Contains(t, goMod, "golang.org/x/crypto")
            },
        },
        {
            name: "PostgreSQL includes correct driver",
            config: TestConfig{Database: "postgresql"},
            validate: func(t *testing.T, projectPath string) {
                goMod := readFile(t, filepath.Join(projectPath, "go.mod"))
                assert.Contains(t, goMod, "github.com/lib/pq")
            },
        },
        {
            name: "Framework dependencies are correctly versioned",
            config: TestConfig{Framework: "gin"},
            validate: func(t *testing.T, projectPath string) {
                goMod := readFile(t, filepath.Join(projectPath, "go.mod"))
                
                // Should have Gin with correct version
                assert.Contains(t, goMod, "github.com/gin-gonic/gin v1.9")
                assert.Contains(t, goMod, "github.com/gin-contrib/cors v1.4")
                
                // Should not have other framework dependencies
                assert.NotContains(t, goMod, "github.com/labstack/echo")
                assert.NotContains(t, goMod, "github.com/gofiber/fiber")
            },
        },
        {
            name: "Logger dependencies match selection",
            config: TestConfig{Logger: "zerolog"},
            validate: func(t *testing.T, projectPath string) {
                goMod := readFile(t, filepath.Join(projectPath, "go.mod"))
                
                // Should have zerolog dependency
                assert.Contains(t, goMod, "github.com/rs/zerolog")
                
                // Should not have other logger dependencies
                assert.NotContains(t, goMod, "go.uber.org/zap")
                assert.NotContains(t, goMod, "github.com/sirupsen/logrus")
            },
        },
        {
            name: "Test dependencies are appropriate for project type",
            config: TestConfig{DatabaseDriver: "postgresql"},
            validate: func(t *testing.T, projectPath string) {
                goMod := readFile(t, filepath.Join(projectPath, "go.mod"))
                
                // Should have testify for assertions
                assert.Contains(t, goMod, "github.com/stretchr/testify")
                
                // Should have testcontainers for database testing
                assert.Contains(t, goMod, "github.com/testcontainers/testcontainers-go")
                assert.Contains(t, goMod, "testcontainers-go/modules/postgres")
            },
        },
        {
            name: "No conflicting dependencies in go.mod",
            config: TestConfig{Framework: "gin", Logger: "zap"},
            validate: func(t *testing.T, projectPath string) {
                goMod := readFile(t, filepath.Join(projectPath, "go.mod"))
                
                // Should not have multiple HTTP routers
                ginCount := strings.Count(goMod, "gin-gonic")
                echoCount := strings.Count(goMod, "labstack/echo")
                fiberCount := strings.Count(goMod, "gofiber/fiber")
                
                httpRouters := 0
                if ginCount > 0 { httpRouters++ }
                if echoCount > 0 { httpRouters++ }  
                if fiberCount > 0 { httpRouters++ }
                
                assert.LessOrEqual(t, httpRouters, 1, "Should not have multiple HTTP framework dependencies")
                
                // Should not have multiple logger libraries (excluding stdlib)
                zapCount := strings.Count(goMod, "go.uber.org/zap")
                logrusCount := strings.Count(goMod, "sirupsen/logrus")
                zerologCount := strings.Count(goMod, "rs/zerolog")
                
                loggers := 0
                if zapCount > 0 { loggers++ }
                if logrusCount > 0 { loggers++ }
                if zerologCount > 0 { loggers++ }
                
                assert.LessOrEqual(t, loggers, 1, "Should not have multiple logger dependencies")
            },
        },
        {
            name: "Raw SQL projects don't include GORM dependencies",
            config: TestConfig{DatabaseDriver: "postgresql", DatabaseORM: ""},
            validate: func(t *testing.T, projectPath string) {
                goMod := readFile(t, filepath.Join(projectPath, "go.mod"))
                
                // Should have raw SQL driver
                assert.Contains(t, goMod, "github.com/lib/pq", "Should have PostgreSQL driver for raw SQL")
                
                // Should NOT have GORM dependencies
                assert.NotContains(t, goMod, "gorm.io/gorm", "Raw SQL project should not have GORM dependency")
                assert.NotContains(t, goMod, "gorm.io/driver/postgres", "Raw SQL project should not have GORM PostgreSQL driver")
                
                // Connection file should use database/sql
                connectionFile := findDatabaseConnectionFile(projectPath)
                if connectionFile != "" {
                    content := readFile(t, connectionFile)
                    assert.Contains(t, content, `"database/sql"`, "Should use database/sql for raw SQL")
                    assert.NotContains(t, content, `"gorm.io/gorm"`, "Should not import GORM")
                }
            },
        },
        {
            name: "GORM projects don't include raw SQL drivers", 
            config: TestConfig{DatabaseDriver: "postgresql", DatabaseORM: "gorm"},
            validate: func(t *testing.T, projectPath string) {
                goMod := readFile(t, filepath.Join(projectPath, "go.mod"))
                
                // Should have GORM dependencies
                assert.Contains(t, goMod, "gorm.io/gorm", "GORM project should have GORM dependency")
                assert.Contains(t, goMod, "gorm.io/driver/postgres", "GORM project should have GORM PostgreSQL driver")
                
                // Should NOT have raw SQL drivers (GORM includes its own)
                assert.NotContains(t, goMod, "github.com/lib/pq", "GORM project should not have raw SQL driver dependency")
                
                // Connection file should use GORM
                connectionFile := findDatabaseConnectionFile(projectPath)
                if connectionFile != "" {
                    content := readFile(t, connectionFile)
                    assert.Contains(t, content, `"gorm.io/gorm"`, "Should use GORM imports")
                    assert.NotContains(t, content, `"database/sql"`, "Should not import database/sql when using GORM")
                }
            },
        },
        {
            name: "MySQL raw SQL uses correct driver",
            config: TestConfig{DatabaseDriver: "mysql", DatabaseORM: ""},
            validate: func(t *testing.T, projectPath string) {
                goMod := readFile(t, filepath.Join(projectPath, "go.mod"))
                
                // Should have MySQL raw SQL driver
                assert.Contains(t, goMod, "github.com/go-sql-driver/mysql", "Should have MySQL driver for raw SQL")
                
                // Should not have PostgreSQL or SQLite drivers
                assert.NotContains(t, goMod, "github.com/lib/pq", "Should not have PostgreSQL driver")
                assert.NotContains(t, goMod, "github.com/mattn/go-sqlite3", "Should not have SQLite driver")
                
                // Should not have GORM dependencies
                assert.NotContains(t, goMod, "gorm.io/gorm", "Raw SQL project should not have GORM")
                assert.NotContains(t, goMod, "gorm.io/driver/mysql", "Raw SQL should not have GORM MySQL driver")
            },
        },
        {
            name: "Integration tests import consistency with ORM selection",
            config: TestConfig{DatabaseDriver: "postgresql", DatabaseORM: ""},
            validate: func(t *testing.T, projectPath string) {
                integrationTestFile := filepath.Join(projectPath, "tests/integration/api_test.go")
                if fileExists(integrationTestFile) {
                    content := readFile(t, integrationTestFile)
                    
                    // Raw SQL project should import database/sql in tests
                    assert.Contains(t, content, `"database/sql"`, "Integration tests should import database/sql for raw SQL projects")
                    
                    // Should NOT import GORM in tests
                    assert.NotContains(t, content, `"gorm.io/gorm"`, "Integration tests should not import GORM for raw SQL projects")
                    
                    // Should use *sql.DB type casting
                    assert.Contains(t, content, "(*sql.DB)", "Integration tests should cast to *sql.DB for raw SQL projects")
                    assert.NotContains(t, content, "(*gorm.DB)", "Integration tests should not cast to *gorm.DB for raw SQL projects")
                }
            },
        },
        {
            name: "Integration tests import GORM when ORM selected",
            config: TestConfig{DatabaseDriver: "postgresql", DatabaseORM: "gorm"},
            validate: func(t *testing.T, projectPath string) {
                integrationTestFile := filepath.Join(projectPath, "tests/integration/api_test.go")
                if fileExists(integrationTestFile) {
                    content := readFile(t, integrationTestFile)
                    
                    // GORM project should import GORM in tests
                    assert.Contains(t, content, `"gorm.io/gorm"`, "Integration tests should import GORM for GORM projects")
                    
                    // Should NOT import database/sql in tests (GORM handles it)
                    assert.NotContains(t, content, `"database/sql"`, "Integration tests should not import database/sql for GORM projects")
                    
                    // Should use *gorm.DB type casting
                    assert.Contains(t, content, "(*gorm.DB)", "Integration tests should cast to *gorm.DB for GORM projects")
                    assert.NotContains(t, content, "(*sql.DB)", "Integration tests should not cast to *sql.DB for GORM projects")
                }
            },
        },
        {
            name: "Test framework imports match project framework",
            config: TestConfig{Framework: "fiber"},
            validate: func(t *testing.T, projectPath string) {
                testFiles := findGoFiles(filepath.Join(projectPath, "tests"))
                
                for _, file := range testFiles {
                    content := readFile(t, file)
                    
                    // If file imports any framework, should be Fiber
                    if strings.Contains(content, "gin-gonic") || strings.Contains(content, "labstack/echo") || strings.Contains(content, "gofiber/fiber") {
                        assert.Contains(t, content, `"github.com/gofiber/fiber/v2"`, "Test files should use Fiber when Fiber is selected")
                        assert.NotContains(t, content, `"github.com/gin-gonic/gin"`, "Test files should not import Gin when Fiber is selected")
                        assert.NotContains(t, content, `"github.com/labstack/echo"`, "Test files should not import Echo when Fiber is selected")
                    }
                }
            },
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            projectPath := generateProject(t, tc.config)
            tc.validate(t, projectPath)
        })
    }
}
```

### 6. Architecture Pattern Validation

```go
func TestArchitecturePatternConsistency(t *testing.T) {
    testCases := []struct {
        name         string  
        architecture string
        validate     func(t *testing.T, projectPath string)
    }{
        {
            name: "Clean Architecture follows dependency rule",
            architecture: "clean",
            validate: func(t *testing.T, projectPath string) {
                // Domain layer should not import infrastructure
                domainFiles := findGoFiles(filepath.Join(projectPath, "internal/domain"))
                for _, file := range domainFiles {
                    content := readFile(t, file)
                    assert.NotContains(t, content, "/internal/infrastructure")
                    assert.NotContains(t, content, "/internal/adapters")
                    assert.NotContains(t, content, `"gorm.io/gorm"`)
                    assert.NotContains(t, content, `"github.com/gin-gonic/gin"`)
                }
                
                // Application layer should not import infrastructure
                appFiles := findGoFiles(filepath.Join(projectPath, "internal/application"))
                for _, file := range appFiles {
                    content := readFile(t, file)
                    assert.NotContains(t, content, "/internal/infrastructure")
                    assert.NotContains(t, content, `"gorm.io/gorm"`)
                }
            },
        },
        {
            name: "DDD patterns have proper aggregate boundaries",
            architecture: "ddd",
            validate: func(t *testing.T, projectPath string) {
                // Domain services should not directly depend on repositories
                domainServices := findGoFiles(filepath.Join(projectPath, "internal/domain/services"))
                for _, file := range domainServices {
                    content := readFile(t, file)
                    assert.NotContains(t, content, "/internal/infrastructure/persistence")
                    assert.NotContains(t, content, `"gorm.io/gorm"`)
                    assert.NotContains(t, content, `"database/sql"`)
                }
                
                // Entities should not import external packages
                entityFiles := findGoFiles(filepath.Join(projectPath, "internal/domain/entities"))
                for _, file := range entityFiles {
                    content := readFile(t, file)
                    assert.NotContains(t, content, `"gorm.io/gorm"`)
                    assert.NotContains(t, content, "/internal/infrastructure")
                    assert.NotContains(t, content, "/internal/application")
                }
            },
        },
        {
            name: "Hexagonal architecture has proper port/adapter separation",
            architecture: "hexagonal", 
            validate: func(t *testing.T, projectPath string) {
                // Domain should only use interfaces (ports)
                domainFiles := findGoFiles(filepath.Join(projectPath, "internal/domain"))
                for _, file := range domainFiles {
                    content := readFile(t, file)
                    assert.NotContains(t, content, "/internal/adapters")
                    assert.NotContains(t, content, `"gorm.io/gorm"`)
                    assert.NotContains(t, content, `"github.com/gin-gonic/gin"`)
                }
                
                // Primary adapters should not know about secondary adapters
                primaryAdapters := findGoFiles(filepath.Join(projectPath, "internal/adapters/primary"))
                for _, file := range primaryAdapters {
                    content := readFile(t, file)
                    assert.NotContains(t, content, "/internal/adapters/secondary")
                    assert.NotContains(t, content, `"gorm.io/gorm"`)
                    assert.NotContains(t, content, `"database/sql"`)
                }
            },
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            config := TestConfig{Architecture: tc.architecture}
            projectPath := generateProject(t, config)
            tc.validate(t, projectPath)
        })
    }
}
```

### 7. Cross-Blueprint Consistency Testing

```go
func TestCrossBlueprintConsistency(t *testing.T) {
    blueprints := []string{"web-api-standard", "web-api-clean", "web-api-ddd", "web-api-hexagonal"}
    commonConfig := TestConfig{
        Framework: "gin",
        Database: "postgresql", 
        Logger: "slog",
        AuthType: "jwt",
    }
    
    // Generate all blueprint types with same config
    projectPaths := make(map[string]string)
    for _, blueprint := range blueprints {
        config := commonConfig
        config.Blueprint = blueprint
        projectPaths[blueprint] = generateProject(t, config)
    }
    
    t.Run("Common files have consistent patterns", func(t *testing.T) {
        commonFiles := []string{
            "go.mod",
            "cmd/server/main.go", 
            "internal/config/config.go",
            "docker-compose.yml",
        }
        
        for _, file := range commonFiles {
            blueprintContents := make(map[string]string)
            
            // Read file from each blueprint
            for blueprint, projectPath := range projectPaths {
                filePath := filepath.Join(projectPath, file)
                if fileExists(filePath) {
                    blueprintContents[blueprint] = readFile(t, filePath)
                }
            }
            
            // Validate consistency
            for blueprint, content := range blueprintContents {
                // All should use same framework
                assert.Contains(t, content, "gin", "Blueprint %s file %s should use gin", blueprint, file)
                
                // All should use same database
                if strings.Contains(file, "go.mod") || strings.Contains(file, "docker") {
                    assert.Contains(t, content, "postgres", "Blueprint %s file %s should use postgres", blueprint, file)
                }
                
                // All should have auth dependencies when auth enabled
                if strings.Contains(file, "go.mod") {
                    assert.Contains(t, content, "golang.org/x/crypto", "Blueprint %s should have crypto dependency", blueprint)
                }
            }
        }
    })
    
    t.Run("Logger implementation is consistent", func(t *testing.T) {
        for blueprint, projectPath := range projectPaths {
            loggerFiles := findGoFiles(filepath.Join(projectPath, "internal/logger"))
            
            for _, file := range loggerFiles {
                content := readFile(t, file)
                
                // Should use slog consistently
                if strings.Contains(content, "logger") {
                    assert.Contains(t, content, "log/slog", "Blueprint %s should use slog in %s", blueprint, file)
                    assert.NotContains(t, content, "go.uber.org/zap", "Blueprint %s should not use zap in %s", blueprint, file)
                }
            }
        }
    })
    
    t.Run("Database connection patterns are consistent", func(t *testing.T) {
        for blueprint, projectPath := range projectPaths {
            // Find database connection files (may be in different locations per architecture)
            connectionFiles := []string{
                "internal/database/connection.go",
                "internal/infrastructure/persistence/database.go", 
                "internal/adapters/secondary/persistence/database.go",
            }
            
            for _, connFile := range connectionFiles {
                filePath := filepath.Join(projectPath, connFile)
                if fileExists(filePath) {
                    content := readFile(t, filePath)
                    
                    // Should use PostgreSQL driver
                    assert.Contains(t, content, `"github.com/lib/pq"`, "Blueprint %s should use postgres driver", blueprint)
                    assert.NotContains(t, content, `"github.com/mattn/go-sqlite3"`, "Blueprint %s should not use sqlite driver", blueprint)
                    
                    // Should have consistent connection patterns
                    assert.Contains(t, content, "postgres", "Blueprint %s should reference postgres", blueprint)
                }
            }
        }
    })
}
```

### 8. Performance and Resource Testing

```go
func TestGeneratedCodePerformance(t *testing.T) {
    testCases := []struct {
        name   string
        config TestConfig
        check  func(t *testing.T, projectPath string)
    }{
        {
            name: "Database connections use proper pooling",
            config: TestConfig{DatabaseDriver: "postgresql"},
            check: func(t *testing.T, projectPath string) {
                connectionFiles := findGoFiles(projectPath)
                
                for _, file := range connectionFiles {
                    if strings.Contains(file, "connection") || strings.Contains(file, "database") {
                        content := readFile(t, file)
                        
                        if strings.Contains(content, "sql.Open") || strings.Contains(content, "gorm.Open") {
                            // Should configure connection pool
                            assert.Contains(t, content, "SetMaxIdleConns", "File %s should configure idle connections", file)
                            assert.Contains(t, content, "SetMaxOpenConns", "File %s should configure max connections", file)
                            assert.Contains(t, content, "SetConnMaxLifetime", "File %s should configure connection lifetime", file)
                        }
                    }
                }
            },
        },
        {
            name: "HTTP server has proper timeouts",
            config: TestConfig{Framework: "gin"},
            check: func(t *testing.T, projectPath string) {
                mainFile := filepath.Join(projectPath, "cmd/server/main.go")
                content := readFile(t, mainFile)
                
                // Should configure server timeouts
                if strings.Contains(content, "http.Server") {
                    assert.Contains(t, content, "ReadTimeout", "Server should have read timeout")
                    assert.Contains(t, content, "WriteTimeout", "Server should have write timeout")
                    assert.Contains(t, content, "IdleTimeout", "Server should have idle timeout")
                }
            },
        },
        {
            name: "No resource leaks in generated code",
            config: TestConfig{DatabaseDriver: "postgresql"},
            check: func(t *testing.T, projectPath string) {
                goFiles := findGoFiles(projectPath)
                
                for _, file := range goFiles {
                    content := readFile(t, file)
                    
                    // Check for proper resource cleanup
                    if strings.Contains(content, ".Query(") || strings.Contains(content, ".QueryRow(") {
                        assert.Contains(t, content, "defer", "File %s should have defer cleanup for queries", file)
                    }
                    
                    // Check for proper connection closing
                    if strings.Contains(content, "http.Client") {
                        assert.Contains(t, content, "defer", "File %s should cleanup HTTP resources", file)
                    }
                }
            },
        },
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            projectPath := generateProject(t, tc.config)
            tc.check(t, projectPath)
        })
    }
}
```

## Matrix-Based Test Case Generation Strategy

### Configuration Matrix Dimensions

To systematically identify ALL possible BDD test cases, we need to create a **multi-dimensional matrix** based on user prompt selections:

#### **Primary Selection Dimensions**
```
Blueprint Type:    [web-api-standard, web-api-clean, web-api-ddd, web-api-hexagonal, cli, library, lambda, microservice, monolith]
Framework:         [gin, echo, fiber, chi, stdlib]  
Database Driver:   [postgresql, mysql, sqlite, redis, mongodb, ""] (none)
Database ORM:      [gorm, sqlx, sqlc, ent, ""] (raw SQL)
Logger:           [slog, zap, logrus, zerolog]
Authentication:   [jwt, oauth2, session, api-key, ""] (none)
Architecture:     [standard, clean, ddd, hexagonal, event-driven]
Go Version:       [1.21, 1.22, 1.23]
```

#### **Validation Dimensions**
```
Import Categories:     [framework, database, logger, auth, testing, stdlib, models, middleware]
Dependency Types:      [go.mod entries, versions, conflicts, unused]
File Categories:       [.go, .yaml, .sql, .md, .tmpl, Dockerfile, Makefile]
Code Patterns:        [interfaces, structs, middleware, handlers, services, repositories]
Configuration Types:  [env vars, config files, docker-compose, kubernetes]
```

### **Matrix-Generated Test Categories**

#### **1. Import Consistency Matrix**
```
Test Case = [Blueprint] × [Framework] × [Database] × [ORM] × [Logger] × [Auth] × [File Type]

Examples:
- web-api-standard × gin × postgresql × "" × slog × jwt × connection.go
  → Should import: gin, lib/pq, log/slog, jwt, crypto
  → Should NOT import: gorm, echo, fiber, zap, logrus

- web-api-clean × fiber × mysql × gorm × zap × "" × repository.go  
  → Should import: fiber, gorm, mysql driver, zap
  → Should NOT import: gin, lib/pq, jwt, slog

- cli × cobra × "" × "" × zerolog × "" × main.go
  → Should import: cobra, zerolog
  → Should NOT import: gin, database drivers, jwt
```

#### **2. Dependency Management Matrix**
```
Test Case = [Framework] × [Database] × [ORM] × [Logger] × [Auth] = go.mod Validation

Examples:
gin × postgresql × gorm × zap × jwt:
  ✅ Should have: gin, gorm, gorm/postgres, zap, jwt, crypto
  ❌ Should NOT have: echo, fiber, lib/pq, slog, logrus

echo × mysql × "" × slog × "":
  ✅ Should have: echo, go-sql-driver/mysql  
  ❌ Should NOT have: gin, gorm, zap, jwt, crypto
```

#### **3. Architecture Pattern Matrix**
```
Test Case = [Architecture] × [Blueprint] × [Framework] = Code Structure Validation

Examples:
clean × web-api × gin:
  ✅ Domain layer should NOT import: gin, gorm, infrastructure
  ✅ Application layer should NOT import: gin, gorm
  ✅ Infrastructure layer CAN import: gin, gorm

ddd × web-api × fiber:
  ✅ Entities should NOT import: fiber, gorm, application
  ✅ Domain services should NOT import: fiber, persistence
  ✅ Application services CAN import: domain, NOT infrastructure
```

#### **4. Template Logic Matrix**  
```
Test Case = [Selection Value] × [Template Condition] = Generated Content

Examples:
DatabaseDriver="postgresql" × {{if eq .DatabaseDriver "postgresql"}}:
  ✅ Should generate PostgreSQL-specific code
  ❌ Should NOT generate MySQL/SQLite code

AuthType="" × {{if ne .AuthType ""}}:
  ✅ Should NOT generate auth middleware
  ❌ Should NOT have JWT imports
```

### **Systematic Test Case Generation Algorithm**

```go
type TestMatrix struct {
    Blueprint    []string
    Framework    []string
    Database     []string
    ORM          []string
    Logger       []string
    Auth         []string
    Architecture []string
}

func GenerateAllTestCases(matrix TestMatrix) []TestCase {
    var testCases []TestCase
    
    // Generate all combinations
    for _, blueprint := range matrix.Blueprint {
        for _, framework := range matrix.Framework {
            for _, database := range matrix.Database {
                for _, orm := range matrix.ORM {
                    for _, logger := range matrix.Logger {
                        for _, auth := range matrix.Auth {
                            for _, arch := range matrix.Architecture {
                                config := TestConfig{
                                    Blueprint: blueprint,
                                    Framework: framework,
                                    Database: database,
                                    ORM: orm,
                                    Logger: logger,
                                    Auth: auth,
                                    Architecture: arch,
                                }
                                
                                // Generate validation rules for this combination
                                testCases = append(testCases, 
                                    generateImportTests(config),
                                    generateDependencyTests(config),
                                    generateArchitectureTests(config),
                                    generateTemplateLogicTests(config),
                                    generateConfigurationTests(config),
                                )
                            }
                        }
                    }
                }
            }
        }
    }
    
    return testCases
}

func generateImportTests(config TestConfig) []TestCase {
    expected := calculateExpectedImports(config)
    forbidden := calculateForbiddenImports(config)
    
    return []TestCase{
        {
            Name: fmt.Sprintf("Import consistency for %s", config.String()),
            Validate: func(t *testing.T, projectPath string) {
                for _, file := range findRelevantFiles(projectPath, config) {
                    content := readFile(t, file)
                    
                    for _, imp := range expected[filepath.Base(file)] {
                        assert.Contains(t, content, imp, 
                            "File %s should import %s for config %s", file, imp, config)
                    }
                    
                    for _, imp := range forbidden[filepath.Base(file)] {
                        assert.NotContains(t, content, imp,
                            "File %s should NOT import %s for config %s", file, imp, config)
                    }
                }
            },
        },
    }
}
```

### **Configuration-Specific Expected Outcomes**

#### **Import Expectation Rules**
```go
var ImportRules = map[string]ImportRule{
    "framework": {
        "gin":    []string{`"github.com/gin-gonic/gin"`},
        "echo":   []string{`"github.com/labstack/echo/v4"`},
        "fiber":  []string{`"github.com/gofiber/fiber/v2"`},
        "chi":    []string{`"github.com/go-chi/chi/v5"`},
    },
    "database": {
        "postgresql+raw": []string{`"github.com/lib/pq"`, `"database/sql"`},
        "postgresql+gorm": []string{`"gorm.io/gorm"`, `"gorm.io/driver/postgres"`},
        "mysql+raw": []string{`"github.com/go-sql-driver/mysql"`, `"database/sql"`},
        "mysql+gorm": []string{`"gorm.io/gorm"`, `"gorm.io/driver/mysql"`},
    },
    "logger": {
        "slog":    []string{`"log/slog"`},
        "zap":     []string{`"go.uber.org/zap"`},
        "logrus":  []string{`"github.com/sirupsen/logrus"`},
        "zerolog": []string{`"github.com/rs/zerolog"`},
    },
    "auth": {
        "jwt":     []string{`"github.com/golang-jwt/jwt/v5"`, `"golang.org/x/crypto"`},
        "oauth2":  []string{`"golang.org/x/oauth2"`},
        "session": []string{`"github.com/gorilla/sessions"`},
    },
}

var ForbiddenImports = map[string][]string{
    "framework": {
        "gin":   []string{`"github.com/labstack/echo"`, `"github.com/gofiber/fiber"`},
        "echo":  []string{`"github.com/gin-gonic/gin"`, `"github.com/gofiber/fiber"`},
        "fiber": []string{`"github.com/gin-gonic/gin"`, `"github.com/labstack/echo"`},
    },
    "database": {
        "postgresql+raw":  []string{`"gorm.io/gorm"`, `"gorm.io/driver/postgres"`},
        "postgresql+gorm": []string{`"github.com/lib/pq"`, `"database/sql"`},
        "mysql+raw":       []string{`"gorm.io/gorm"`, `"gorm.io/driver/mysql"`},
        "mysql+gorm":      []string{`"github.com/go-sql-driver/mysql"`, `"database/sql"`},
    },
}
```

#### **File-Specific Rules**
```go
var FileTypeRules = map[string]FileRule{
    "connection.go": {
        RequiresDatabase: true,
        ImportRules: map[string][]string{
            "postgresql+raw": []string{`"database/sql"`, `"github.com/lib/pq"`},
            "postgresql+gorm": []string{`"gorm.io/gorm"`, `"gorm.io/driver/postgres"`},
        },
        ForbiddenImports: map[string][]string{
            "postgresql+raw": []string{`"gorm.io/gorm"`},
            "postgresql+gorm": []string{`"database/sql"`, `"github.com/lib/pq"`},
        },
    },
    "api_test.go": {
        RequiresDatabase: true,
        ImportRules: map[string][]string{
            "postgresql+raw": []string{`"database/sql"`},
            "postgresql+gorm": []string{`"gorm.io/gorm"`},
        },
        TypeCastingRules: map[string][]string{
            "postgresql+raw": []string{"(*sql.DB)"},
            "postgresql+gorm": []string{"(*gorm.DB)"},
        },
    },
    "auth.go": {
        RequiresAuth: true,
        ImportRules: map[string][]string{
            "jwt+gin": []string{`"github.com/gin-gonic/gin"`, `"github.com/golang-jwt/jwt/v5"`},
            "jwt+fiber": []string{`"github.com/gofiber/fiber/v2"`, `"github.com/golang-jwt/jwt/v5"`},
        },
        ContextTypeRules: map[string][]string{
            "gin": []string{"*gin.Context"},
            "fiber": []string{"*fiber.Ctx"},
            "echo": []string{"echo.Context"},
        },
    },
}
```

### **Matrix Calculation & Test Case Generation**

#### **Total Matrix Size Analysis**
```
Selection Dimensions:
- Blueprint Type: 9 options (web-api-standard, web-api-clean, web-api-ddd, web-api-hexagonal, cli, library, lambda, microservice, monolith)
- Framework: 5 options (gin, echo, fiber, chi, stdlib)
- Database Driver: 6 options (postgresql, mysql, sqlite, redis, mongodb, "")
- Database ORM: 5 options (gorm, sqlx, sqlc, ent, "")
- Logger: 4 options (slog, zap, logrus, zerolog)
- Authentication: 6 options (jwt, oauth2, session, api-key, basic, "")
- Architecture: 5 options (standard, clean, ddd, hexagonal, event-driven)

Total Theoretical Combinations: 9 × 5 × 6 × 5 × 4 × 6 × 5 = 40,500 combinations
```

#### **Combination Validation & Filtering**
```go
type CombinationValidator struct {
    ValidCombinations map[string][]string
    InvalidCombinations map[string][]string
    RequiredDependencies map[string][]string
}

var CombinationRules = CombinationValidator{
    // Valid combinations only
    ValidCombinations: map[string][]string{
        "blueprint": {
            "web-api-*": []string{"gin", "echo", "fiber", "chi"}, // Web APIs need web frameworks
            "cli": []string{"cobra", ""}, // CLI doesn't need web framework
            "library": []string{""}, // Libraries don't need frameworks
            "lambda": []string{"gin", "echo", "fiber", ""}, // Lambda can have web framework for API Gateway
        },
        "database+orm": {
            "postgresql": []string{"gorm", "sqlx", "sqlc", ""}, // PostgreSQL supports all ORMs
            "mysql": []string{"gorm", "sqlx", ""}, // MySQL doesn't support sqlc well
            "sqlite": []string{"gorm", "sqlx", ""}, 
            "redis": []string{""}, // Redis doesn't use traditional ORMs
            "mongodb": []string{""}, // MongoDB uses its own driver
            "": []string{""}, // No database = no ORM
        },
        "architecture+blueprint": {
            "clean": []string{"web-api-*"}, // Clean architecture only for web APIs
            "ddd": []string{"web-api-*"}, // DDD only for complex web APIs
            "hexagonal": []string{"web-api-*"}, // Hexagonal only for web APIs
            "standard": []string{"*"}, // Standard works for all
            "event-driven": []string{"web-api-*", "microservice"}, // Event-driven for distributed systems
        },
    },
    
    InvalidCombinations: map[string][]string{
        "cli+database": []string{"postgresql", "mysql", "sqlite"}, // CLI rarely needs databases
        "library+framework": []string{"gin", "echo", "fiber"}, // Libraries don't use web frameworks
        "lambda+database": []string{"postgresql", "mysql"}, // Lambda typically uses serverless databases
    },
}

func (cv *CombinationValidator) IsValidCombination(config TestConfig) bool {
    // Check blueprint + framework compatibility
    if !cv.isValidBlueprintFramework(config.Blueprint, config.Framework) {
        return false
    }
    
    // Check database + ORM compatibility  
    if !cv.isValidDatabaseORM(config.Database, config.ORM) {
        return false
    }
    
    // Check architecture + blueprint compatibility
    if !cv.isValidArchitectureBlueprint(config.Architecture, config.Blueprint) {
        return false
    }
    
    return true
}
```

#### **Priority Matrix for Test Case Generation**
```go
type PriorityLevel int
const (
    Critical PriorityLevel = iota // Must work, highest user impact
    High                         // Important combinations, common usage
    Medium                       // Less common but still valid
    Low                          // Edge cases, rare usage
)

var CombinationPriority = map[string]PriorityLevel{
    // Critical combinations - most common user selections
    "web-api-standard+gin+postgresql+raw+slog+jwt": Critical,
    "web-api-standard+gin+postgresql+gorm+slog+jwt": Critical,
    "web-api-standard+echo+mysql+gorm+zap+jwt": Critical,
    "web-api-standard+fiber+postgresql+raw+zerolog+none": Critical,
    "cli+none+none+none+slog+none": Critical,
    "library+none+none+none+slog+none": Critical,
    
    // High priority - common architectural patterns
    "web-api-clean+gin+postgresql+gorm+slog+jwt": High,
    "web-api-ddd+echo+postgresql+sqlx+zap+oauth2": High,
    "web-api-hexagonal+fiber+mysql+gorm+logrus+session": High,
    "microservice+gin+redis+none+zap+jwt": High,
    
    // Medium priority - valid but less common
    "web-api-standard+chi+sqlite+sqlx+zerolog+api-key": Medium,
    "lambda+gin+none+none+slog+jwt": Medium,
    "monolith+echo+postgresql+gorm+zap+session": Medium,
    
    // Low priority - edge cases
    "web-api-standard+stdlib+mongodb+none+logrus+basic": Low,
}
```

#### **Systematic Test Case Generation Engine**
```go
type TestCaseGenerator struct {
    Matrix TestMatrix
    Rules ValidationRules
    Priority map[string]PriorityLevel
}

func (tcg *TestCaseGenerator) GenerateAllTestCases() ([]FeatureFile, error) {
    var featureFiles []FeatureFile
    
    // Generate all valid combinations
    validCombinations := tcg.generateValidCombinations()
    
    // Group by priority
    prioritizedCombinations := tcg.prioritizeCombinations(validCombinations)
    
    // Generate test cases for each combination
    for priority, combinations := range prioritizedCombinations {
        for _, combination := range combinations {
            testCases := tcg.generateTestCasesForCombination(combination)
            
            featureFile := FeatureFile{
                Name: fmt.Sprintf("%s_%s.feature", combination.Category(), priority),
                Scenarios: tcg.convertToGherkinScenarios(testCases),
                Priority: priority,
            }
            
            featureFiles = append(featureFiles, featureFile)
        }
    }
    
    return featureFiles, nil
}

func (tcg *TestCaseGenerator) generateTestCasesForCombination(config TestConfig) []TestCase {
    var testCases []TestCase
    
    // 1. Import validation test cases
    testCases = append(testCases, tcg.generateImportTests(config)...)
    
    // 2. Dependency management test cases  
    testCases = append(testCases, tcg.generateDependencyTests(config)...)
    
    // 3. File generation test cases
    testCases = append(testCases, tcg.generateFileGenerationTests(config)...)
    
    // 4. Architecture compliance test cases
    testCases = append(testCases, tcg.generateArchitectureTests(config)...)
    
    // 5. Configuration consistency test cases
    testCases = append(testCases, tcg.generateConfigurationTests(config)...)
    
    // 6. Template logic test cases
    testCases = append(testCases, tcg.generateTemplateLogicTests(config)...)
    
    return testCases
}
```

#### **Comprehensive Validation Rules Engine**
```go
type ValidationRules struct {
    ImportRules        map[string]ImportRule
    DependencyRules    map[string]DependencyRule
    FileGenerationRules map[string]FileRule
    ArchitectureRules  map[string]ArchitectureRule
    ConfigurationRules map[string]ConfigRule
}

// Import validation rules - exactly what should be imported/forbidden for each combination
var ImportValidationRules = map[string]ImportRule{
    "web-api-standard+gin+postgresql+raw+slog+jwt": {
        ExpectedImports: map[string][]string{
            "connection.go": []string{`"database/sql"`, `"github.com/lib/pq"`, `"log/slog"`},
            "main.go": []string{`"github.com/gin-gonic/gin"`, `"log/slog"`},
            "auth.go": []string{`"github.com/golang-jwt/jwt/v5"`, `"golang.org/x/crypto"`, `"github.com/gin-gonic/gin"`},
            "api_test.go": []string{`"database/sql"`, `"github.com/gin-gonic/gin"`, `"github.com/stretchr/testify"`},
        },
        ForbiddenImports: map[string][]string{
            "connection.go": []string{`"gorm.io/gorm"`, `"github.com/labstack/echo"`, `"go.uber.org/zap"`},
            "main.go": []string{`"gorm.io/gorm"`, `"github.com/labstack/echo"`, `"github.com/gofiber/fiber"`},
            "auth.go": []string{`"github.com/labstack/echo"`, `"github.com/gofiber/fiber"`},
            "api_test.go": []string{`"gorm.io/gorm"`, `"github.com/labstack/echo"`},
        },
        TypeCastingRules: map[string][]string{
            "api_test.go": []string{"(*sql.DB)"}, // Should cast to *sql.DB, not *gorm.DB
        },
    },
    
    "web-api-clean+fiber+mysql+gorm+zap+oauth2": {
        ExpectedImports: map[string][]string{
            "database.go": []string{`"gorm.io/gorm"`, `"gorm.io/driver/mysql"`, `"go.uber.org/zap"`},
            "main.go": []string{`"github.com/gofiber/fiber/v2"`, `"go.uber.org/zap"`},
            "auth.go": []string{`"golang.org/x/oauth2"`, `"github.com/gofiber/fiber/v2"`},
            "api_test.go": []string{`"gorm.io/gorm"`, `"github.com/gofiber/fiber/v2"`},
        },
        ForbiddenImports: map[string][]string{
            "database.go": []string{`"database/sql"`, `"github.com/lib/pq"`, `"log/slog"`},
            "main.go": []string{`"github.com/gin-gonic/gin"`, `"github.com/labstack/echo"`},
            "auth.go": []string{`"github.com/golang-jwt/jwt/v5"`, `"github.com/gin-gonic/gin"`},
            "api_test.go": []string{`"database/sql"`, `"github.com/gin-gonic/gin"`},
        },
        TypeCastingRules: map[string][]string{
            "api_test.go": []string{"(*gorm.DB)"}, // Should cast to *gorm.DB, not *sql.DB
        },
        ContextTypeRules: map[string][]string{
            "auth.go": []string{"*fiber.Ctx"}, // Should use Fiber context, not Gin
        },
    },
    
    "cli+cobra+none+none+zerolog+none": {
        ExpectedImports: map[string][]string{
            "main.go": []string{`"github.com/spf13/cobra"`, `"github.com/rs/zerolog"`},
            "cmd/root.go": []string{`"github.com/spf13/cobra"`, `"github.com/rs/zerolog"`},
        },
        ForbiddenImports: map[string][]string{
            "main.go": []string{`"github.com/gin-gonic/gin"`, `"database/sql"`, `"gorm.io/gorm"`, `"github.com/golang-jwt/jwt"`},
            "cmd/root.go": []string{`"github.com/gin-gonic/gin"`, `"database/sql"`},
        },
    },
}

// Architecture compliance rules - what each architecture pattern should/shouldn't allow
var ArchitectureValidationRules = map[string]ArchitectureRule{
    "clean": {
        LayerRules: map[string]LayerRule{
            "domain": {
                CanImport: []string{"internal/domain", "standard library"},
                CannotImport: []string{"internal/infrastructure", "internal/adapters", "github.com/gin-gonic", "gorm.io/gorm"},
                RequiredPatterns: []string{"interfaces", "entities", "value objects"},
            },
            "application": {
                CanImport: []string{"internal/domain", "internal/application", "standard library"},
                CannotImport: []string{"internal/infrastructure", "github.com/gin-gonic", "gorm.io/gorm"},
                RequiredPatterns: []string{"use cases", "application services"},
            },
            "infrastructure": {
                CanImport: []string{"*"}, // Infrastructure can import anything
                CannotImport: []string{}, // No restrictions
                RequiredPatterns: []string{"repositories", "external services"},
            },
        },
    },
    
    "ddd": {
        LayerRules: map[string]LayerRule{
            "domain/entities": {
                CanImport: []string{"internal/domain", "standard library"},
                CannotImport: []string{"internal/application", "internal/infrastructure", "gorm.io/gorm", "github.com/gin-gonic"},
                RequiredPatterns: []string{"domain entities", "domain events"},
            },
            "domain/services": {
                CanImport: []string{"internal/domain", "standard library"},
                CannotImport: []string{"internal/infrastructure", "gorm.io/gorm"},
                RequiredPatterns: []string{"domain services", "interfaces"},
            },
            "application": {
                CanImport: []string{"internal/domain", "internal/application", "standard library"},
                CannotImport: []string{"internal/infrastructure", "gorm.io/gorm"},
                RequiredPatterns: []string{"application services", "command handlers"},
            },
        },
    },
    
    "hexagonal": {
        LayerRules: map[string]LayerRule{
            "domain": {
                CanImport: []string{"internal/domain", "standard library"},
                CannotImport: []string{"internal/adapters", "gorm.io/gorm", "github.com/gin-gonic"},
                RequiredPatterns: []string{"ports", "domain logic"},
            },
            "adapters/primary": {
                CanImport: []string{"internal/domain", "internal/adapters/primary", "web frameworks"},
                CannotImport: []string{"internal/adapters/secondary", "gorm.io/gorm", "database/sql"},
                RequiredPatterns: []string{"HTTP handlers", "controllers"},
            },
            "adapters/secondary": {
                CanImport: []string{"internal/domain", "internal/adapters/secondary", "database libraries"},
                CannotImport: []string{"internal/adapters/primary", "github.com/gin-gonic"},
                RequiredPatterns: []string{"repository implementations", "external services"},
            },
        },
    },
}
```

#### **BDD Feature File Generation**
```go
func (tcg *TestCaseGenerator) convertToGherkinScenarios(testCases []TestCase) []GherkinScenario {
    var scenarios []GherkinScenario
    
    for _, testCase := range testCases {
        scenario := GherkinScenario{
            Name: testCase.Name,
            Given: tcg.generateGivenSteps(testCase.Config),
            When: tcg.generateWhenSteps(testCase.Config),
            Then: tcg.generateThenSteps(testCase.ValidationRules),
        }
        scenarios = append(scenarios, scenario)
    }
    
    return scenarios
}

func (tcg *TestCaseGenerator) generateGivenSteps(config TestConfig) []string {
    var steps []string
    
    steps = append(steps, "Given the go-starter CLI tool is available")
    steps = append(steps, "And I am in a clean working directory")
    
    if config.Blueprint != "" {
        steps = append(steps, fmt.Sprintf(`And I want to create a "%s" project`, config.Blueprint))
    }
    
    if config.Framework != "" {
        steps = append(steps, fmt.Sprintf(`And I select "%s" as my web framework`, config.Framework))
    }
    
    if config.Database != "" {
        steps = append(steps, fmt.Sprintf(`And I select "%s" as my database driver`, config.Database))
    }
    
    if config.ORM != "" {
        steps = append(steps, fmt.Sprintf(`And I select "%s" as my ORM`, config.ORM))
    } else if config.Database != "" {
        steps = append(steps, "And I choose raw SQL without an ORM")
    }
    
    if config.Logger != "" {
        steps = append(steps, fmt.Sprintf(`And I select "%s" as my logger`, config.Logger))
    }
    
    if config.Auth != "" {
        steps = append(steps, fmt.Sprintf(`And I select "%s" authentication`, config.Auth))
    } else {
        steps = append(steps, "And I choose not to include authentication")
    }
    
    if config.Architecture != "standard" {
        steps = append(steps, fmt.Sprintf(`And I select "%s" architecture pattern`, config.Architecture))
    }
    
    return steps
}

func (tcg *TestCaseGenerator) generateThenSteps(rules ValidationRules) []string {
    var steps []string
    
    // Import validation steps
    for file, expectedImports := range rules.ExpectedImports {
        for _, imp := range expectedImports {
            steps = append(steps, fmt.Sprintf(`Then "%s" should import "%s"`, file, imp))
        }
    }
    
    for file, forbiddenImports := range rules.ForbiddenImports {
        for _, imp := range forbiddenImports {
            steps = append(steps, fmt.Sprintf(`And "%s" should not import "%s"`, file, imp))
        }
    }
    
    // Dependency validation steps
    for _, dep := range rules.ExpectedDependencies {
        steps = append(steps, fmt.Sprintf(`And go.mod should contain "%s"`, dep))
    }
    
    for _, dep := range rules.ForbiddenDependencies {
        steps = append(steps, fmt.Sprintf(`And go.mod should not contain "%s"`, dep))
    }
    
    // Architecture validation steps
    for layer, rule := range rules.ArchitectureRules {
        for _, forbidden := range rule.CannotImport {
            steps = append(steps, fmt.Sprintf(`And "%s" layer should not import "%s"`, layer, forbidden))
        }
    }
    
    // Compilation and quality steps
    steps = append(steps, "And the generated project should compile without errors")
    steps = append(steps, "And golangci-lint should pass without warnings")
    steps = append(steps, "And there should be no unused imports")
    steps = append(steps, "And there should be no unused variables")
    
    return steps
}
```

## Recommended Feature File Structure

Given this matrix-based approach, we should organize tests into **systematic feature files** based on the generated combinations:

### Core Quality Features
```
tests/acceptance/quality/
├── static_analysis.feature          # Import/variable analysis, linting
├── compilation_validation.feature   # Build success, go vet, basic quality
├── dependency_management.feature    # go.mod consistency, conflict detection
└── template_logic.feature          # Conditional generation, fallback logic
```

### Configuration Consistency Features
```  
tests/acceptance/configuration/
├── database_consistency.feature     # Database driver, ORM, migration consistency
├── framework_consistency.feature    # Framework imports, middleware, context types
├── logger_consistency.feature       # Logger selection, import consistency
├── authentication_consistency.feature # Auth imports, middleware, dependencies
└── environment_consistency.feature  # Config files, env vars, Docker settings
```

### Architecture Pattern Features
```
tests/acceptance/architecture/
├── clean_architecture.feature       # Clean arch dependency rules, layer isolation
├── ddd_patterns.feature             # Domain boundaries, aggregate rules
├── hexagonal_architecture.feature   # Port/adapter separation, dependency direction
└── cross_blueprint_consistency.feature # Shared patterns across architectures
```

### Blueprint-Specific Features
```
tests/acceptance/blueprints/
├── web_api_standard.feature         # Standard architecture specific tests
├── web_api_clean.feature           # Clean architecture specific tests  
├── web_api_ddd.feature             # DDD specific tests
├── web_api_hexagonal.feature       # Hexagonal specific tests
├── cli_blueprints.feature          # CLI-specific quality tests
└── lambda_blueprints.feature       # Lambda-specific quality tests
```

### Performance & Resource Features
```
tests/acceptance/performance/
├── resource_management.feature      # Connection pooling, timeouts, cleanup
├── security_standards.feature       # Security headers, auth patterns, CORS
├── deployment_readiness.feature     # Docker, Makefile, deployment configs
└── integration_testing.feature      # Generated tests work correctly
```

### Generated Test Quality Features
```
tests/acceptance/test_quality/
├── test_import_consistency.feature   # Test file imports match project config
├── test_compilation.feature          # Generated tests compile and run
├── test_database_consistency.feature # Test DB setup matches project ORM
└── test_framework_consistency.feature # Test framework usage matches selection
```

### Matrix Testing Features
```
tests/acceptance/matrix/
├── framework_database_matrix.feature    # All framework × database combinations
├── logger_architecture_matrix.feature   # All logger × architecture combinations  
├── auth_framework_matrix.feature        # All auth × framework combinations
└── critical_combinations.feature        # High-priority configuration sets
```

## Sample Feature File Structure

### Example: `dependency_management.feature`
```gherkin
Feature: Dependency Management Consistency
  As a developer using go-starter
  I want generated projects to have correct dependencies
  So that my project compiles without conflicts

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Raw SQL PostgreSQL project has correct dependencies
    Given I want to create a web API with raw SQL
    When I generate a project with PostgreSQL and no ORM
    Then the go.mod should contain "github.com/lib/pq"
    And the go.mod should not contain "gorm.io/gorm"
    And the go.mod should not contain "gorm.io/driver/postgres"
    And the database connection should use "database/sql"

  Scenario: GORM PostgreSQL project has correct dependencies  
    Given I want to create a web API with GORM
    When I generate a project with PostgreSQL and GORM ORM
    Then the go.mod should contain "gorm.io/gorm"
    And the go.mod should contain "gorm.io/driver/postgres"
    And the go.mod should not contain "github.com/lib/pq"
    And the database connection should use GORM imports

  Scenario: No conflicting database drivers in go.mod
    Given I select PostgreSQL as my database
    When I generate any web API project
    Then the go.mod should not contain MySQL drivers
    And the go.mod should not contain SQLite drivers
    And the go.mod should not contain multiple PostgreSQL drivers

  Scenario: Framework dependencies are exclusive
    Given I select Gin as my framework
    When I generate any web API project  
    Then the go.mod should contain Gin dependencies
    And the go.mod should not contain Echo dependencies
    And the go.mod should not contain Fiber dependencies
    And the go.mod should not contain Chi dependencies
```

### Example: `static_analysis.feature`
```gherkin
Feature: Static Code Analysis Quality
  As a developer using go-starter
  I want generated code to pass static analysis
  So that my project follows Go best practices

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Generated code has no unused imports
    Given I generate a project with any valid configuration
    When I analyze all Go files for unused imports
    Then there should be no unused import statements
    And golangci-lint should pass without import warnings

  Scenario: Generated code has no unused variables
    Given I generate a project with any valid configuration  
    When I analyze all Go files for unused variables
    Then there should be no unused variable declarations
    And go vet should pass without variable warnings

  Scenario: Conditional imports match actual usage
    Given I generate a project with raw SQL database
    When I analyze the database connection file
    Then models package should not be imported
    And only database/sql should be imported
    And no GORM packages should be imported

  Scenario: Generated code compiles cleanly
    Given I generate a project with any valid configuration
    When I run "go build ./..."
    Then the build should succeed without errors
    And there should be no compilation warnings
```

### Example: `framework_consistency.feature`
```gherkin
Feature: Framework Consistency Across Files
  As a developer selecting a web framework
  I want all generated code to use my selected framework consistently
  So that there are no framework conflicts or mixing

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario Outline: Framework imports are consistent across all files
    Given I select "<framework>" as my web framework
    When I generate a web API project
    Then all middleware files should import "<expected_import>"
    And no files should import "<forbidden_gin>"
    And no files should import "<forbidden_echo>"
    And no files should import "<forbidden_fiber>"
    And context types should be "<context_type>"

    Examples:
      | framework | expected_import              | forbidden_gin        | forbidden_echo        | forbidden_fiber        | context_type  |
      | gin       | github.com/gin-gonic/gin     | N/A                 | labstack/echo         | gofiber/fiber         | *gin.Context  |
      | echo      | github.com/labstack/echo/v4  | gin-gonic/gin       | N/A                   | gofiber/fiber         | echo.Context  |
      | fiber     | github.com/gofiber/fiber/v2  | gin-gonic/gin       | labstack/echo         | N/A                   | *fiber.Ctx    |
      | chi       | github.com/go-chi/chi/v5     | gin-gonic/gin       | labstack/echo         | gofiber/fiber         | http.ResponseWriter |

  Scenario: Authentication middleware matches selected framework
    Given I select Fiber as my framework and JWT authentication
    When I generate a web API project
    Then auth middleware should use "*fiber.Ctx" context type
    And auth middleware should import "github.com/gofiber/fiber/v2"
    And auth middleware should not import Gin or Echo packages
```

### Example: `test_import_consistency.feature`
```gherkin
Feature: Generated Test Import Consistency
  As a developer using go-starter
  I want generated test files to have correct imports
  So that my tests compile and run without errors

  Background:
    Given the go-starter CLI tool is available
    And I am in a clean working directory

  Scenario: Raw SQL integration tests don't import GORM
    Given I select raw SQL database with PostgreSQL
    When I generate a web API project
    Then integration test files should import "database/sql"
    And integration test files should not import "gorm.io/gorm"
    And database setup should use "*sql.DB" type casting
    And repository creation should use "*sql.DB" parameter

  Scenario: GORM integration tests don't import database/sql
    Given I select GORM ORM with PostgreSQL
    When I generate a web API project  
    Then integration test files should import "gorm.io/gorm"
    And integration test files should not import "database/sql"
    And database setup should use "*gorm.DB" type casting
    And repository creation should use "*gorm.DB" parameter

  Scenario: Framework-specific test setup
    Given I select Fiber as my web framework
    When I generate a web API project
    Then integration test files should import "github.com/gofiber/fiber/v2"
    And integration test files should not import Gin packages
    And test router setup should use Fiber API
    And request creation should use Fiber patterns

  Scenario: Test dependency imports match project selection
    Given I select Echo framework with Zap logger
    When I generate a web API project
    Then integration test files should import "github.com/labstack/echo/v4"
    And integration test files should import Zap logger packages
    And integration test files should not import other framework packages
    And integration test files should not import other logger packages

  Scenario: Authentication test imports are conditional
    Given I select a project without authentication
    When I generate a web API project
    Then integration test files should not import JWT packages
    And integration test files should not import "golang.org/x/crypto"
    And test setup should not include auth service creation
    And test routes should not include auth endpoints

  Scenario: Authentication test imports are present when needed
    Given I select JWT authentication
    When I generate a web API project
    Then integration test files should import JWT packages
    And integration test files should import "golang.org/x/crypto"  
    And test setup should include auth service creation
    And test routes should include auth endpoints
    And test cases should include authentication scenarios
```

## Zero-Downtime Migration Strategy

### Migration Philosophy

Instead of replacing existing tests immediately, we implement a **parallel development approach** that maintains continuous ATDD coverage while building superior template quality validation.

### Migration Timeline with Continuous Coverage

#### **Phase 1: Foundation + Parallel Development (Week 1-2)**

**Maintain Existing Coverage:**
```bash
# ✅ KEEP existing tests running in CI
tests/acceptance/blueprints/           # Continue legacy functional tests
├── web-api/features/*.feature         # Maintain baseline coverage
├── monolith/monolith_atdd_test.go    # Keep functional validation
└── */features/*.feature              # Current ATDD coverage
```

**Build Enhanced Tests in Parallel:**
```bash
# 🆕 BUILD alongside existing (not replacement)
tests/acceptance/enhanced/            # New enhanced validation
├── quality/
│   ├── static_analysis_test.go      # Import/variable analysis
│   ├── dependency_management_test.go # go.mod consistency
│   └── compilation_validation_test.go # Build + lint validation
└── infrastructure/
    └── test_helpers.go              # Supporting functions
```

**CI Configuration (Dual Coverage):**
```yaml
jobs:
  legacy-atdd:           # ✅ MAINTAIN existing coverage
    runs-on: ubuntu-latest
    steps:
      - name: Run existing BDD tests
        run: go test ./tests/acceptance/blueprints/... -v
  
  enhanced-atdd-pilot:   # 🆕 ADD pilot enhanced tests
    runs-on: ubuntu-latest
    steps:
      - name: Run enhanced quality tests
        run: go test ./tests/acceptance/enhanced/quality/... -v
```

**Week 1-2 Deliverables:**
- [x] Legacy tests continue providing 100% baseline coverage
- [ ] Create `static_analysis_test.go` with unused import detection
- [ ] Create `dependency_management_test.go` with go.mod validation
- [ ] Validate enhanced tests catch issues legacy tests miss
- [ ] Enhanced tests run as **additional** CI jobs (not replacements)

#### **Phase 2: Expansion + Matrix Testing (Week 3-4)**

**Enhanced Test Expansion:**
```bash
tests/acceptance/enhanced/
├── quality/                         # ✅ Completed in Phase 1
├── configuration/                   # 🆕 Matrix-based testing
│   ├── framework_consistency_test.go # Cross-contamination prevention
│   ├── database_consistency_test.go  # Driver/ORM combinations
│   └── logger_consistency_test.go    # Logger import validation
└── matrix/                          # 🆕 Combination testing
    └── critical_combinations_test.go # High-priority configs
```

**CI Configuration (Expanded):**
```yaml
jobs:
  legacy-atdd:           # ✅ STILL running existing
    runs-on: ubuntu-latest
    steps:
      - name: Run existing BDD tests  
        run: go test ./tests/acceptance/blueprints/... -v
  
  enhanced-atdd:         # 🆕 EXPANDED enhanced coverage
    runs-on: ubuntu-latest
    strategy:
      matrix:
        category: [quality, configuration, matrix]
    steps:
      - name: Run enhanced tests
        run: go test ./tests/acceptance/enhanced/${{ matrix.category }}/... -v
```

**Week 3-4 Deliverables:**
- [x] Legacy tests still providing safety net coverage
- [ ] Enhanced tests covering 50%+ of enhanced strategy requirements
- [ ] Matrix-based testing for critical configuration combinations
- [ ] Framework cross-contamination prevention validation
- [ ] Performance monitoring: enhanced tests should run faster

#### **Phase 3: Architecture + Comprehensive Coverage (Week 5-6)**

**Complete Enhanced Implementation:**
```bash
tests/acceptance/enhanced/
├── quality/              # ✅ Core quality validation
├── configuration/        # ✅ Matrix-based config testing
├── architecture/         # 🆕 Dependency rule enforcement
│   ├── clean_architecture_test.go
│   ├── ddd_patterns_test.go
│   └── hexagonal_architecture_test.go
├── matrix/              # ✅ Critical combinations
├── performance/         # 🆕 Resource management
│   └── resource_management_test.go
└── test_quality/        # 🆕 Generated test validation
    └── test_import_consistency_test.go
```

**CI Configuration (Validation Phase):**
```yaml
jobs:
  legacy-atdd:           # ✅ STILL running (safety net)
    runs-on: ubuntu-latest
    steps:
      - name: Run existing BDD tests
        run: go test ./tests/acceptance/blueprints/... -v
  
  enhanced-atdd-full:    # 🆕 COMPREHENSIVE enhanced coverage
    runs-on: ubuntu-latest
    strategy:
      matrix:
        category: [quality, configuration, architecture, matrix, performance]
    steps:
      - name: Run enhanced ATDD tests
        run: go test ./tests/acceptance/enhanced/${{ matrix.category }}/... -v -parallel 4
```

**Week 5-6 Deliverables:**
- [x] Legacy tests continue as safety net
- [ ] Enhanced tests achieve 90%+ coverage of enhanced strategy
- [ ] Architecture dependency rule validation
- [ ] Performance improvement: enhanced tests run 70%+ faster than legacy
- [ ] Quality validation: enhanced tests catch template issues legacy tests miss

#### **Phase 4: Validation + Cutover Preparation (Week 7)**

**Coverage Validation:**
```bash
# Validate enhanced tests provide superior coverage
./scripts/validate-enhanced-coverage.sh

# Compare issue detection capabilities
./scripts/compare-test-quality.sh
```

**Gradual Feature Flag Rollout:**
```yaml
jobs:
  atdd-tests:
    runs-on: ubuntu-latest
    steps:
      - name: Choose ATDD approach based on feature flag
        run: |
          if [[ "${{ github.event_name }}" == "pull_request" ]]; then
            # Pull requests use enhanced tests (faster feedback)
            echo "ATDD_MODE=enhanced" >> $GITHUB_ENV
          else
            # Main branch uses both (extra safety during transition)
            echo "ATDD_MODE=both" >> $GITHUB_ENV
          fi
```

**Week 7 Deliverables:**
- [ ] Comprehensive comparison of legacy vs enhanced test effectiveness
- [ ] Enhanced tests proven to catch 90%+ more template quality issues
- [ ] Performance validation: enhanced tests complete in <15 minutes vs 45+ minutes legacy
- [ ] Team approval for legacy test retirement
- [ ] Rollback plan documented and tested

#### **Phase 5: Legacy Retirement + Production Deployment (Week 8)**

**Production CI Configuration:**
```yaml
jobs:
  enhanced-atdd-production:  # ✅ FULL enhanced coverage only
    runs-on: ubuntu-latest
    strategy:
      matrix:
        category: [quality, configuration, architecture, matrix, performance, security]
    steps:
      - name: Run production enhanced ATDD
        run: go test ./tests/acceptance/enhanced/${{ matrix.category }}/... -v -parallel 4
```

**Legacy Test Retirement:**
```bash
# Archive legacy tests (don't delete - keep for reference)
mkdir -p archive/legacy-tests/
mv tests/acceptance/blueprints/ archive/legacy-tests/
mv tests/acceptance/features/ archive/legacy-tests/

# Update documentation
git add archive/legacy-tests/
git commit -m "archive: retire legacy ATDD tests in favor of enhanced strategy

- Legacy tests archived to archive/legacy-tests/
- Enhanced ATDD strategy now provides superior template quality validation
- Performance improvement: 78% faster test execution  
- Quality improvement: 90% more template issues detected"
```

**Week 8 Deliverables:**
- [ ] Legacy tests retired and archived
- [ ] Enhanced tests running in production CI/CD
- [ ] Performance metrics: <15 minute ATDD execution vs previous 45+ minutes
- [ ] Quality metrics: Template quality issues caught pre-production
- [ ] Documentation updated for new ATDD maintenance

### Risk Mitigation & Rollback Strategy

#### **Continuous Validation Checkpoints**

**Before Each Phase:**
```bash
# Validate enhanced tests don't introduce regressions
./scripts/validate-migration-safety.sh

# Compare coverage metrics
./scripts/coverage-comparison.sh legacy enhanced

# Performance benchmarking
./scripts/benchmark-test-performance.sh
```

#### **Rollback Plan**

**If Enhanced Tests Have Issues:**
```bash
# Quick rollback to legacy tests
git revert <enhanced-test-commits>
cp .github/workflows/ci-legacy-backup.yml .github/workflows/ci.yml
git commit -m "rollback: restore legacy ATDD tests due to enhanced test issues"
```

#### **Success Criteria for Each Phase**

**Phase 1 Success Criteria:**
- [ ] Enhanced tests run successfully alongside legacy tests
- [ ] Enhanced tests catch at least 1 template quality issue legacy tests miss
- [ ] CI execution time increases by <10%

**Phase 2 Success Criteria:**
- [ ] Enhanced tests cover matrix-based configuration testing
- [ ] Framework cross-contamination prevention working
- [ ] Enhanced tests execution time <20 minutes

**Phase 3 Success Criteria:**
- [ ] Enhanced tests achieve comprehensive architecture validation
- [ ] Enhanced tests catch 70%+ more issues than legacy tests
- [ ] Enhanced tests run 50%+ faster than legacy tests

**Phase 4 Success Criteria:**
- [ ] Enhanced tests proven superior in coverage and performance
- [ ] Team confidence in enhanced approach
- [ ] Rollback plan tested and ready

**Phase 5 Success Criteria:**  
- [ ] Legacy tests successfully retired
- [ ] Enhanced tests provide complete ATDD coverage
- [ ] CI performance improvement achieved
- [ ] Template quality improvement measurable

### Benefits of Zero-Downtime Approach

#### **✅ Continuous Coverage**
- **0% ATDD coverage loss** during migration
- Legacy tests provide safety net while enhanced tests mature  
- CI continues catching regressions throughout transition

#### **⚡ Performance Improvement Timeline**
```
Week 1-2: Legacy (45min) + Enhanced pilot (5min) = 50min total
Week 3-4: Legacy (45min) + Enhanced expanded (8min) = 53min total
Week 5-6: Legacy (45min) + Enhanced full (12min) = 57min total  
Week 7-8: Enhanced only (12min) = 78% improvement ⚡
```

#### **🔍 Quality Improvement Timeline**
```
Week 1-2: Legacy detection + Enhanced pilot detection (10% more issues)
Week 3-4: Legacy + Enhanced catching 30% more template issues
Week 5-6: Legacy + Enhanced catching 60% more template issues
Week 7-8: Enhanced only catching 90% more issues than legacy alone
```

This zero-downtime migration ensures we never lose ATDD coverage while progressively building superior template quality validation that aligns with the enhanced strategy requirements.

## Parallel Execution Benefits

### Why Multiple Feature Files Enable Parallel Testing

The structured approach with **multiple focused feature files** provides significant performance and organizational benefits:

#### **🚀 Parallel Execution Advantages**

1. **Independent Test Execution**
   ```bash
   # Run different feature categories in parallel
   go test ./tests/acceptance/quality/... &          # Static analysis tests
   go test ./tests/acceptance/configuration/... &    # Configuration tests  
   go test ./tests/acceptance/architecture/... &     # Architecture tests
   go test ./tests/acceptance/performance/... &      # Performance tests
   wait  # Wait for all parallel jobs to complete
   ```

2. **Resource Utilization**
   - **CPU**: Multiple test suites can run on different cores
   - **I/O**: File system operations distributed across feature files
   - **Memory**: Smaller test scopes reduce memory pressure per process

3. **Faster Feedback Loops**
   ```bash
   # Sequential (current): ~45 minutes for full ATDD suite
   # Parallel (proposed): ~12 minutes with 4-way parallelism
   ```

#### **⚡ CI/CD Pipeline Optimization**

```yaml
# GitHub Actions Example
name: Enhanced ATDD Testing
on: [push, pull_request]

jobs:
  quality-tests:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        feature: [static_analysis, compilation_validation, dependency_management, template_logic]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
      - name: Run Quality Feature Tests
        run: go test ./tests/acceptance/quality/${{ matrix.feature }}_test.go -v

  configuration-tests:
    runs-on: ubuntu-latest  
    strategy:
      matrix:
        feature: [database_consistency, framework_consistency, logger_consistency, authentication_consistency]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
      - name: Run Configuration Feature Tests
        run: go test ./tests/acceptance/configuration/${{ matrix.feature }}_test.go -v

  architecture-tests:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        feature: [clean_architecture, ddd_patterns, hexagonal_architecture, cross_blueprint_consistency]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
      - name: Run Architecture Feature Tests
        run: go test ./tests/acceptance/architecture/${{ matrix.feature }}_test.go -v

  matrix-tests:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        combination: [gin-postgres, fiber-mysql, echo-sqlite, chi-postgres]
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v4
      - name: Run Matrix Combination Tests
        run: go test ./tests/acceptance/matrix/ -run TestCombination${{ matrix.combination }} -v
```

#### **📊 Performance Comparison**

| Approach | Total Time | Parallelism | Resource Usage | Maintainability |
|----------|------------|-------------|----------------|-----------------|
| **Single Suite** | ~45 min | None | High per-process | Poor - monolithic |
| **Category-Based** | ~12 min | 4x parallel | Distributed | Good - focused |
| **Feature-Based** | ~8 min | 8x parallel | Optimized | Excellent - granular |

#### **🔧 Local Development Benefits**

```bash
# Run only relevant tests during development
npm run test:quality          # Quick static analysis feedback
npm run test:database         # Only database-related tests  
npm run test:architecture     # Only architecture validation

# Run critical path tests first
npm run test:critical         # Essential quality gates
npm run test:full            # Complete test suite (CI only)
```

#### **📈 Scalability Advantages**

1. **Team Parallel Development**
   - **Frontend team**: Focus on `static_analysis.feature`
   - **Backend team**: Focus on `database_consistency.feature`  
   - **DevOps team**: Focus on `deployment_readiness.feature`
   - **Architecture team**: Focus on `architecture/*.feature`

2. **Incremental Test Development**
   ```bash
   # Add new feature files without disrupting existing ones
   tests/acceptance/new_category/
   ├── new_feature.feature
   └── new_feature_test.go
   ```

3. **Selective Test Execution**
   ```bash
   # Run only changed feature areas
   git diff --name-only | grep "blueprints/web-api" && npm run test:web-api
   git diff --name-only | grep "logger" && npm run test:logger
   ```

#### **🎯 Test Isolation Benefits**

1. **Failure Isolation**: One failing feature doesn't block others
2. **Clear Ownership**: Each feature file has focused responsibility  
3. **Easier Debugging**: Smaller test scopes make issues easier to locate
4. **Independent Versioning**: Feature tests can evolve independently

#### **⚙️ Makefile Integration**

```make
# Parallel test execution targets
.PHONY: test-quality test-config test-arch test-matrix test-all-parallel

test-quality:
	@echo "Running quality tests..."
	go test ./tests/acceptance/quality/... -v -parallel 4

test-config:
	@echo "Running configuration tests..."  
	go test ./tests/acceptance/configuration/... -v -parallel 4

test-arch:
	@echo "Running architecture tests..."
	go test ./tests/acceptance/architecture/... -v -parallel 4

test-matrix:
	@echo "Running matrix tests..."
	go test ./tests/acceptance/matrix/... -v -parallel 2

test-all-parallel: 
	@echo "Running all ATDD tests in parallel..."
	$(MAKE) test-quality & \
	$(MAKE) test-config & \
	$(MAKE) test-arch & \
	$(MAKE) test-matrix & \
	wait
	@echo "All parallel tests completed!"

# Quick feedback for developers
test-critical:
	go test ./tests/acceptance/quality/static_analysis_test.go -v
	go test ./tests/acceptance/quality/compilation_validation_test.go -v  
	go test ./tests/acceptance/configuration/dependency_management_test.go -v
```

## Quality Metrics

### Template Quality Score
- **Static Analysis**: No linting errors (20%)
- **Import Cleanliness**: No unused imports (20%) 
- **Variable Usage**: No unused variables (15%)
- **Configuration Consistency**: User selections applied (25%)
- **Compilation Success**: Clean builds (20%)

### Coverage Targets
- **Blueprint Combinations**: 90% of critical combinations tested
- **Template Conditions**: 100% of conditional logic paths tested
- **Framework Combinations**: All framework × database × logger combinations
- **Error Detection**: 95% of template issues caught before release

## Benefits

1. **Prevent Template Regressions**: Catch issues before they reach users
2. **Improve Code Quality**: Generated code follows Go best practices
3. **Configuration Accuracy**: User selections are correctly applied
4. **Faster Development**: Automated detection vs manual discovery
5. **User Confidence**: Generated projects work correctly out-of-the-box

This enhanced ATDD strategy would have caught all the template issues we've been fixing manually.