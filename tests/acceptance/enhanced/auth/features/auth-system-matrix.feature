Feature: Authentication System Matrix Testing
  As a go-starter user building secure applications
  I want to ensure that all authentication systems work correctly with various configurations
  So that user authentication and authorization are secure and reliable

  Background:
    Given I have the go-starter CLI available
    And all templates are properly initialized
    And I am testing authentication system combinations

  @critical @auth @jwt @p0
  Scenario Outline: JWT authentication system combinations
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | jwt         |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And JWT authentication should be properly configured
    And JWT middleware should be implemented
    And token generation should work correctly
    And token validation should be secure
    And JWT signing should use secure algorithms
    And token expiration should be configurable
    And refresh token support should be available
    And JWT claims should be properly structured
    And middleware should handle invalid tokens gracefully
    And authentication routes should be secured

    Examples:
      | framework | database | orm  | logger  |
      | gin       | postgres | gorm | zap     |
      | echo      | mysql    | sqlx | slog    |
      | fiber     | sqlite   | gorm | zerolog |
      | chi       | postgres | sqlx | logrus  |
      | gin       | mysql    | gorm | slog    |

  @critical @auth @oauth2 @p0
  Scenario Outline: OAuth2 authentication system combinations
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | oauth2      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And OAuth2 authentication should be properly configured
    And OAuth2 flows should be implemented
    And provider configuration should be available
    And authorization code flow should work
    And access token management should be secure
    And scope validation should be implemented
    And state parameter should prevent CSRF attacks
    And token refresh should be supported
    And user profile retrieval should work
    And provider-specific handlers should be available

    Examples:
      | framework | database | orm  | logger  |
      | gin       | postgres | gorm | zap     |
      | echo      | mysql    | sqlx | slog    |
      | fiber     | sqlite   | gorm | zerolog |
      | chi       | postgres | sqlx | logrus  |

  @critical @auth @session @p0
  Scenario Outline: Session-based authentication combinations
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | session     |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And session-based authentication should be properly configured
    And session middleware should be implemented
    And session storage should be configured
    And session security should be enforced
    And session cookies should be secure
    And session expiration should be managed
    And CSRF protection should be enabled
    And session regeneration should work
    And concurrent session handling should be supported
    And session cleanup should be automated

    Examples:
      | framework | database | orm  | logger  |
      | gin       | postgres | gorm | zap     |
      | echo      | mysql    | sqlx | slog    |
      | fiber     | sqlite   | gorm | zerolog |
      | chi       | postgres | sqlx | logrus  |

  @integration @auth @user-management
  Scenario Outline: User management system integration
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And user registration should be implemented
    And password hashing should be secure
    And user authentication should work
    And password reset functionality should be available
    And email verification should be supported
    And user profile management should be implemented
    And account lockout should prevent brute force attacks
    And user roles and permissions should be supported
    And audit logging should track authentication events

    Examples:
      | framework | database | orm  | auth    | logger |
      | gin       | postgres | gorm | jwt     | zap    |
      | echo      | mysql    | sqlx | oauth2  | slog   |
      | fiber     | sqlite   | gorm | session | zerolog|

  @security @auth @password-security
  Scenario Outline: Password security and hashing
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And password hashing should use bcrypt or stronger
    And password complexity requirements should be enforced
    And password history should be maintained
    And brute force protection should be implemented
    And account lockout should be configurable
    And password reset should be secure
    And timing attack prevention should be implemented
    And secure password storage should be used
    And password rotation should be supported

    Examples:
      | framework | database | auth    | logger |
      | gin       | postgres | jwt     | zap    |
      | echo      | mysql    | oauth2  | slog   |
      | fiber     | sqlite   | session | zerolog|

  @security @auth @token-security
  Scenario Outline: Token security and management
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | auth         | jwt         |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And JWT tokens should use strong signing algorithms
    And token expiration should be properly configured
    And token blacklisting should be supported
    And token rotation should be implemented
    And secure token storage should be used
    And token transmission should be secure
    And token validation should be thorough
    And token claims should be validated
    And token revocation should be supported

    Examples:
      | framework | database | logger |
      | gin       | postgres | zap    |
      | echo      | mysql    | slog   |
      | fiber     | sqlite   | zerolog|

  @integration @auth @middleware
  Scenario Outline: Authentication middleware integration
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And authentication middleware should be properly configured
    And middleware should integrate with <framework>
    And protected routes should require authentication
    And public routes should be accessible
    And middleware should handle authentication errors gracefully
    And CORS should be properly configured for auth
    And preflight requests should be handled
    And authentication headers should be validated
    And middleware ordering should be correct

    Examples:
      | framework | database | auth    | logger |
      | gin       | postgres | jwt     | zap    |
      | echo      | mysql    | oauth2  | slog   |
      | fiber     | sqlite   | session | zerolog|
      | chi       | postgres | jwt     | logrus |

  @integration @auth @authorization
  Scenario Outline: Role-based authorization systems
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And role-based access control should be implemented
    And permission system should be flexible
    And role hierarchy should be supported
    And resource-based permissions should work
    And authorization middleware should be available
    And admin panel should be role-protected
    And API endpoints should enforce permissions
    And role assignment should be secure
    And permission inheritance should work correctly

    Examples:
      | framework | database | orm  | auth    | logger |
      | gin       | postgres | gorm | jwt     | zap    |
      | echo      | mysql    | sqlx | oauth2  | slog   |
      | fiber     | sqlite   | gorm | session | zerolog|

  @performance @auth @optimization
  Scenario Outline: Authentication performance optimization
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And authentication should have minimal performance overhead
    And token validation should be fast
    And database queries should be optimized
    And caching should be used for user data
    And session lookups should be efficient
    And authentication middleware should be lightweight
    And concurrent authentication should be supported
    And performance metrics should be available

    Examples:
      | framework | database | auth    | logger |
      | gin       | postgres | jwt     | zap    |
      | echo      | mysql    | oauth2  | slog   |
      | fiber     | sqlite   | session | zerolog|

  @security @auth @csrf-protection
  Scenario Outline: CSRF protection implementation
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And CSRF protection should be enabled
    And CSRF tokens should be properly generated
    And CSRF validation should be enforced
    And SameSite cookie attributes should be configured
    And Origin header validation should be implemented
    And CSRF exemptions should be minimal
    And Double Submit Cookie pattern should be used if applicable
    And CSRF errors should be handled gracefully

    Examples:
      | framework | database | auth    | logger |
      | gin       | postgres | session | zap    |
      | echo      | mysql    | session | slog   |
      | fiber     | sqlite   | session | zerolog|

  @integration @auth @api-keys
  Scenario Outline: API key authentication system
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | orm          | <orm>       |
      | auth         | api-key     |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And API key authentication should be implemented
    And API key generation should be secure
    And API key validation should be efficient
    And API key scoping should be supported
    And API key rotation should be available
    And Rate limiting should be enforced per key
    And API key revocation should work
    And Usage tracking should be implemented
    And API key management endpoints should be secured

    Examples:
      | framework | database | orm  | logger |
      | gin       | postgres | gorm | zap    |
      | echo      | mysql    | sqlx | slog   |
      | fiber     | sqlite   | gorm | zerolog|

  @monitoring @auth @audit-logging
  Scenario Outline: Authentication audit and monitoring
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And authentication events should be logged
    And failed login attempts should be tracked
    And successful logins should be recorded
    And logout events should be logged
    And suspicious activity should be detected
    And Login patterns should be analyzed
    And Security alerts should be generated
    And Audit trails should be tamper-proof
    And Compliance reporting should be available

    Examples:
      | framework | database | auth    | logger |
      | gin       | postgres | jwt     | zap    |
      | echo      | mysql    | oauth2  | slog   |
      | fiber     | sqlite   | session | zerolog|

  @integration @auth @multi-factor
  Scenario Outline: Multi-factor authentication support
    When I generate a web API project with configuration:
      | type         | web-api     |
      | framework    | <framework> |
      | database     | <database>  |
      | auth         | <auth>      |
      | logger       | <logger>    |
      | go_version   | 1.23       |
    Then the project should compile successfully
    And MFA foundation should be available
    And TOTP support should be implemented
    And SMS verification should be configurable
    And Email verification should be available
    And Backup codes should be generated
    And MFA enforcement should be configurable
    And Recovery options should be secure
    And MFA bypass should be logged
    And Device trust should be manageable

    Examples:
      | framework | database | auth | logger |
      | gin       | postgres | jwt  | zap    |
      | echo      | mysql    | jwt  | slog   |