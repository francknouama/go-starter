Feature: Expanded Configuration Matrix Testing
  As a go-starter maintainer
  I want comprehensive matrix testing across all configuration dimensions
  So that we ensure all combinations work correctly and identify edge cases

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized
    And matrix testing mode is enabled

  Scenario: Complete framework and database matrix coverage
    Given I test all framework and database combinations
    When I generate projects with the full matrix:
      | Framework | Database   | Driver    | ORM     | Expected Result |
      | gin       | postgresql | postgres  | gorm    | Success         |
      | gin       | postgresql | postgres  | sqlx    | Success         |
      | gin       | postgresql | postgres  | sqlc    | Success         |
      | gin       | postgresql | postgres  | (none)  | Success         |
      | gin       | mysql      | mysql     | gorm    | Success         |
      | gin       | mysql      | mysql     | sqlx    | Success         |
      | gin       | mysql      | mysql     | (none)  | Success         |
      | gin       | sqlite     | sqlite3   | gorm    | Success         |
      | gin       | sqlite     | sqlite3   | (none)  | Success         |
      | gin       | mongodb    | mongo     | (none)  | Success         |
      | echo      | postgresql | postgres  | gorm    | Success         |
      | echo      | mysql      | mysql     | gorm    | Success         |
      | echo      | sqlite     | sqlite3   | gorm    | Success         |
      | fiber     | postgresql | postgres  | gorm    | Success         |
      | fiber     | mysql      | mysql     | gorm    | Success         |
      | fiber     | sqlite     | sqlite3   | gorm    | Success         |
      | chi       | postgresql | postgres  | gorm    | Success         |
      | chi       | mysql      | mysql     | sqlx    | Success         |
      | chi       | mongodb    | mongo     | (none)  | Success         |
    Then all valid combinations should generate successfully
    And invalid ORM selections should be prevented
    And database-specific features should be correctly configured

  Scenario: Authentication and authorization matrix
    Given I test authentication across all frameworks
    When I generate projects with auth configurations:
      | Framework | AuthType  | Session | Database   | TokenStorage | Expected |
      | gin       | jwt       | false   | postgresql | redis        | Success  |
      | gin       | jwt       | false   | mysql      | memory       | Success  |
      | gin       | oauth2    | true    | postgresql | database     | Success  |
      | gin       | api-key   | false   | sqlite     | memory       | Success  |
      | gin       | session   | true    | postgresql | redis        | Success  |
      | echo      | jwt       | false   | mysql      | redis        | Success  |
      | echo      | oauth2    | true    | postgresql | database     | Success  |
      | fiber     | jwt       | false   | postgresql | memory       | Success  |
      | fiber     | api-key   | false   | mysql      | database     | Success  |
      | chi       | jwt       | false   | postgresql | redis        | Success  |
      | chi       | session   | true    | mysql      | database     | Success  |
    Then authentication middleware should be properly configured
    And session management should work when enabled
    And token storage should match configuration
    And auth endpoints should be correctly implemented

  Scenario: Logger implementation matrix across architectures
    Given I test all logger types with different architectures
    When I generate projects with logger configurations:
      | Architecture | Logger  | LogLevel | Output  | Structured | Performance |
      | standard     | slog    | info     | json    | true       | Good        |
      | standard     | zap     | debug    | console | true       | Excellent   |
      | standard     | logrus  | warn     | json    | true       | Good        |
      | standard     | zerolog | error    | json    | true       | Excellent   |
      | clean        | slog    | info     | json    | true       | Good        |
      | clean        | zap     | debug    | console | true       | Excellent   |
      | ddd          | slog    | info     | json    | true       | Good        |
      | ddd          | zerolog | debug    | json    | true       | Excellent   |
      | hexagonal    | zap     | info     | json    | true       | Excellent   |
      | hexagonal    | logrus  | debug    | text    | false      | Good        |
    Then logger initialization should match architecture patterns
    And log levels should be properly configured
    And structured logging should work when enabled
    And performance characteristics should match expectations

  Scenario: Deployment target matrix validation
    Given I test deployment configurations
    When I generate projects with deployment targets:
      | Blueprint    | Docker | K8s | Lambda | Vercel | Railway | Terraform | Success |
      | web-api      | true   | true| false  | false  | true    | true      | Yes     |
      | cli          | true   | false| false  | false  | false   | false     | Yes     |
      | lambda       | false  | false| true   | false  | false   | true      | Yes     |
      | microservice | true   | true | false  | false  | true    | true      | Yes     |
      | monolith     | true   | true | false  | true   | true    | false     | Yes     |
    Then deployment configurations should be generated correctly
    And Dockerfiles should be optimized for each blueprint
    And Kubernetes manifests should follow best practices
    And cloud-specific configurations should be valid

  Scenario: Testing framework and coverage matrix
    Given I test different testing configurations
    When I generate projects with testing setups:
      | Framework | TestLib  | Mocking   | Coverage | Integration | E2E     |
      | gin       | testify  | mockery   | true     | true        | true    |
      | gin       | ginkgo   | gomock    | true     | true        | false   |
      | echo      | stdlib   | manual    | false    | true        | false   |
      | echo      | testify  | mockery   | true     | true        | true    |
      | fiber     | testify  | gomock    | true     | false       | false   |
      | chi       | ginkgo   | mockery   | true     | true        | true    |
    Then test files should be generated appropriately
    And mocking frameworks should be integrated
    And coverage tools should be configured
    And test commands should work correctly

  Scenario: Middleware combination matrix
    Given I test middleware stacks across frameworks
    When I configure middleware combinations:
      | Framework | CORS | RateLimit | Compression | Metrics | Tracing | Auth | Cache |
      | gin       | yes  | yes       | yes         | yes     | yes     | jwt  | redis |
      | gin       | yes  | no        | yes         | yes     | no      | none | none  |
      | echo      | yes  | yes       | no          | yes     | yes     | jwt  | memory|
      | fiber     | yes  | yes       | yes         | no      | no      | api  | redis |
      | chi       | no   | yes       | yes         | yes     | yes     | jwt  | none  |
    Then middleware should be registered in correct order
    And middleware configuration should be consistent
    And performance impact should be acceptable
    And middleware conflicts should be prevented

  Scenario: Database migration tool matrix
    Given I test migration tools with databases
    When I configure migration setups:
      | Database   | MigrationTool | Direction | Versioning | Rollback | Success |
      | postgresql | migrate       | up/down   | timestamp  | yes      | Yes     |
      | postgresql | goose         | up/down   | sequential | yes      | Yes     |
      | postgresql | atlas         | up/down   | hash       | yes      | Yes     |
      | mysql      | migrate       | up/down   | timestamp  | yes      | Yes     |
      | mysql      | goose         | up/down   | sequential | yes      | Yes     |
      | sqlite     | migrate       | up/down   | timestamp  | yes      | Yes     |
      | mongodb    | migrate-mongo | up/down   | timestamp  | yes      | Yes     |
    Then migration files should be properly structured
    And migration commands should work correctly
    And rollback functionality should be tested
    And schema versioning should be tracked

  Scenario: Caching strategy matrix
    Given I test caching implementations
    When I configure caching strategies:
      | Framework | CacheType | Backend   | TTL    | Invalidation | Serialization |
      | gin       | local     | memory    | 5min   | manual       | json          |
      | gin       | distributed | redis   | 1hour  | auto         | msgpack       |
      | echo      | local     | memory    | 10min  | ttl-only     | gob           |
      | echo      | distributed | redis   | 30min  | event-based  | json          |
      | fiber     | hybrid    | redis+mem | 15min  | manual       | json          |
    Then cache implementations should work correctly
    And cache invalidation should function properly
    And serialization should be efficient
    And cache metrics should be available

  Scenario: API versioning matrix
    Given I test API versioning strategies
    When I implement versioning approaches:
      | Framework | Strategy  | Location | Migration | Deprecation | Documentation |
      | gin       | url-path  | /v1/     | manual    | headers     | openapi       |
      | gin       | header    | Accept   | auto      | sunset      | swagger       |
      | echo      | url-path  | /api/v1/ | manual    | headers     | openapi       |
      | echo      | subdomain | v1.api   | auto      | redirect    | swagger       |
      | fiber     | query     | ?v=1     | manual    | headers     | openapi       |
    Then versioning should work as configured
    And version migration should be supported
    And deprecation notices should be included
    And API documentation should reflect versions

  Scenario: Error handling and recovery matrix
    Given I test error handling strategies
    When I configure error handling:
      | Architecture | ErrorType    | Recovery  | Logging | Client Response | Monitoring |
      | standard     | panic        | recover   | error   | 500            | sentry     |
      | standard     | validation   | none      | warn    | 400            | metrics    |
      | clean        | domain       | retry     | error   | 422            | sentry     |
      | ddd          | application  | none      | warn    | 400            | logs       |
      | hexagonal    | adapter      | circuit   | error   | 503            | all        |
    Then error handling should follow architecture patterns
    And recovery mechanisms should work correctly
    And error responses should be consistent
    And monitoring integration should function

  Scenario: Configuration management matrix
    Given I test configuration approaches
    When I setup configuration management:
      | Type    | Source      | Format | Validation | Hot-reload | Secrets     |
      | simple  | env         | flat   | basic      | no         | env         |
      | standard| file+env    | yaml   | schema     | no         | env         |
      | advanced| consul      | json   | schema     | yes        | vault       |
      | cloud   | aws-param   | json   | strict     | yes        | aws-secrets |
      | k8s     | configmap   | yaml   | schema     | yes        | k8s-secrets |
    Then configuration loading should work correctly
    And validation should catch errors
    And hot-reload should function when enabled
    And secrets should be properly managed

  Scenario: Matrix test execution optimization
    Given I have a large configuration matrix
    When I execute matrix tests with optimization
    Then execution should be efficient:
      | Strategy              | Time Reduction | Resource Usage | Coverage Impact |
      | Parallel execution    | 60-70%        | +200% CPU      | None            |
      | Smart sampling        | 80-85%        | Normal         | -5%             |
      | Incremental testing   | 70-75%        | Normal         | None            |
      | Priority-based        | 75-80%        | Normal         | -2%             |
    And critical paths should always be tested
    And optimization should not compromise quality
    And results should be deterministic