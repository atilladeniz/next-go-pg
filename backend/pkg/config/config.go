package config

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// =============================================================================
// Environment Variable Registry (Parsed from .env.example)
// =============================================================================
// Validation rules are defined in .env.example using @tags in comments.
// This file parses those tags and validates accordingly.
// To add a new required variable: just add it to .env.example with @tags.

// RequiredIn defines when an environment variable is required
type RequiredIn int

const (
	RequiredNever      RequiredIn = iota // Optional in all environments
	RequiredProduction                   // Required only in production
	RequiredAlways                       // Required in all environments
)

// EnvVarType defines the type of validation to apply
type EnvVarType int

const (
	TypeString   EnvVarType = iota // Any non-empty string
	TypeURL                        // Valid URL format
	TypeHTTPSURL                   // Valid HTTPS URL (production only)
	TypePort                       // Valid port number 1-65535
	TypeEnum                       // One of allowed values
	TypeSecret                     // Secret value (non-empty, min length)
	TypeSSLMode                    // Database SSL mode
)

// EnvVar defines an environment variable with its validation rules
type EnvVar struct {
	Name        string     // Environment variable name
	Description string     // Human-readable description
	Required    RequiredIn // When this variable is required
	Type        EnvVarType // Type of validation
	AllowedVals []string   // For TypeEnum: allowed values
	MinLength   int        // For TypeSecret: minimum length
	Default     string     // Default value
}

// envRegistry is populated from .env.example at init time
var envRegistry []EnvVar

func init() {
	var err error
	envRegistry, err = parseEnvExampleFromFile()
	if err != nil {
		log.Printf("Warning: Could not parse .env.example: %v", err)
		envRegistry = []EnvVar{} // Empty registry, no validation
	}
}

// parseEnvExampleFromFile reads .env.example from filesystem
func parseEnvExampleFromFile() ([]EnvVar, error) {
	// Try multiple locations
	paths := []string{
		".env.example",
		"../.env.example",
		"../../.env.example",
	}

	// Also try relative to executable
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		paths = append(paths,
			filepath.Join(dir, ".env.example"),
			filepath.Join(dir, "..", ".env.example"),
		)
	}

	for _, path := range paths {
		if data, err := os.ReadFile(path); err == nil {
			return parseEnvExampleContent(string(data))
		}
	}

	return nil, fmt.Errorf(".env.example not found")
}

// parseEnvExampleContent parses the content of .env.example
func parseEnvExampleContent(content string) ([]EnvVar, error) {
	var vars []EnvVar
	scanner := bufio.NewScanner(strings.NewReader(content))

	var currentTags string
	var currentDesc string

	// Regex patterns
	tagPattern := regexp.MustCompile(`@(\w+)(?::(\w+))?(?:\|([^@\s]+))?`)
	envPattern := regexp.MustCompile(`^([A-Z][A-Z0-9_]*)=(.*)$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "#") {
			comment := strings.TrimPrefix(line, "#")
			comment = strings.TrimSpace(comment)

			// Check for @tags
			if strings.Contains(comment, "@") {
				currentTags = comment
			} else if comment != "" && !strings.HasPrefix(comment, "-") && !strings.HasPrefix(comment, "=") {
				// Regular description comment
				currentDesc = comment
			}
			continue
		}

		// Parse KEY=value
		if matches := envPattern.FindStringSubmatch(line); matches != nil {
			name := matches[1]
			defaultVal := matches[2]

			ev := EnvVar{
				Name:        name,
				Description: currentDesc,
				Default:     defaultVal,
				Required:    RequiredNever,
				Type:        TypeString,
			}

			// Parse @tags
			if currentTags != "" {
				tagMatches := tagPattern.FindAllStringSubmatch(currentTags, -1)
				for _, tm := range tagMatches {
					tag := tm[1]
					modifier := tm[2]
					value := tm[3]

					switch tag {
					case "required":
						if modifier == "prod" {
							ev.Required = RequiredProduction
						} else {
							ev.Required = RequiredAlways
						}
					case "type":
						switch modifier {
						case "enum":
							ev.Type = TypeEnum
							if value != "" {
								ev.AllowedVals = strings.Split(value, ",")
							}
						case "port":
							ev.Type = TypePort
						case "url":
							ev.Type = TypeURL
						case "https":
							ev.Type = TypeHTTPSURL
						case "secret":
							ev.Type = TypeSecret
						case "sslmode":
							ev.Type = TypeSSLMode
						}
					case "min":
						if n, err := strconv.Atoi(modifier); err == nil {
							ev.MinLength = n
						} else if n, err := strconv.Atoi(value); err == nil {
							ev.MinLength = n
						}
					}
				}
			}

			// Only add vars with validation rules
			if ev.Required != RequiredNever || ev.Type != TypeString {
				vars = append(vars, ev)
			}

			// Reset for next var
			currentTags = ""
			currentDesc = ""
		}
	}

	return vars, scanner.Err()
}

// =============================================================================
// Validation Error Types
// =============================================================================

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationErrors is a collection of validation errors
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	if len(e) == 0 {
		return ""
	}
	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

func (e ValidationErrors) HasErrors() bool {
	return len(e) > 0
}

// =============================================================================
// Config Structs
// =============================================================================

type Config struct {
	Port        string
	Environment string
	LogLevel    string
	FrontendURL string
	Database    DatabaseConfig
	Server      ServerConfig
	Logging     LoggingConfig
}

type LoggingConfig struct {
	AnonymizeIPs bool
	WithCaller   bool
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	Name         string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

type ServerConfig struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// =============================================================================
// Config Loading
// =============================================================================

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", ""),
			Name:         getEnv("DB_NAME", "nextgopg"),
			SSLMode:      getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns: getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
			MaxLifetime:  getEnvAsDuration("DB_MAX_LIFETIME", "5m"),
		},
		Server: ServerConfig{
			ReadTimeout:  getEnvAsDuration("SERVER_READ_TIMEOUT", "10s"),
			WriteTimeout: getEnvAsDuration("SERVER_WRITE_TIMEOUT", "10s"),
			IdleTimeout:  getEnvAsDuration("SERVER_IDLE_TIMEOUT", "60s"),
		},
		Logging: LoggingConfig{
			AnonymizeIPs: getEnvAsBool("LOG_ANONYMIZE_IPS", false),
			WithCaller:   getEnvAsBool("LOG_WITH_CALLER", false),
		},
	}
}

// =============================================================================
// Validation (Registry-Based from .env.example)
// =============================================================================

// Validate checks all registered environment variables based on their rules
func (c *Config) Validate() ValidationErrors {
	var errs ValidationErrors
	isProduction := c.Environment == "production"

	for _, ev := range envRegistry {
		value := c.getValueForEnvVar(ev.Name)

		// Check if required
		if ev.Required == RequiredAlways || (ev.Required == RequiredProduction && isProduction) {
			if value == "" {
				errs = append(errs, ValidationError{
					Field:   ev.Name,
					Message: fmt.Sprintf("is required%s", requiredInMsg(ev.Required)),
				})
				continue
			}
		}

		// Skip further validation if empty and not required
		if value == "" {
			continue
		}

		// Type-specific validation
		if err := c.validateType(ev, value, isProduction); err != nil {
			errs = append(errs, *err)
		}
	}

	return errs
}

// validateType validates a value against its type rules
func (c *Config) validateType(ev EnvVar, value string, isProduction bool) *ValidationError {
	switch ev.Type {
	case TypeEnum:
		if !contains(ev.AllowedVals, value) {
			return &ValidationError{
				Field:   ev.Name,
				Message: fmt.Sprintf("must be one of: %s (got: %s)", strings.Join(ev.AllowedVals, ", "), value),
			}
		}

	case TypePort:
		port, err := strconv.Atoi(value)
		if err != nil || port < 1 || port > 65535 {
			return &ValidationError{
				Field:   ev.Name,
				Message: fmt.Sprintf("must be a valid port number 1-65535 (got: %s)", value),
			}
		}

	case TypeURL:
		if _, err := url.ParseRequestURI(value); err != nil {
			return &ValidationError{
				Field:   ev.Name,
				Message: fmt.Sprintf("must be a valid URL (got: %s)", value),
			}
		}

	case TypeHTTPSURL:
		if _, err := url.ParseRequestURI(value); err != nil {
			return &ValidationError{
				Field:   ev.Name,
				Message: fmt.Sprintf("must be a valid URL (got: %s)", value),
			}
		}
		if isProduction && !strings.HasPrefix(value, "https://") {
			return &ValidationError{
				Field:   ev.Name,
				Message: "must use HTTPS in production",
			}
		}

	case TypeSecret:
		if ev.MinLength > 0 && len(value) < ev.MinLength {
			return &ValidationError{
				Field:   ev.Name,
				Message: fmt.Sprintf("must be at least %d characters", ev.MinLength),
			}
		}

	case TypeSSLMode:
		validModes := []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}
		if !contains(validModes, value) {
			return &ValidationError{
				Field:   ev.Name,
				Message: fmt.Sprintf("must be one of: %s (got: %s)", strings.Join(validModes, ", "), value),
			}
		}
		if isProduction && value == "disable" {
			return &ValidationError{
				Field:   ev.Name,
				Message: "should not be 'disable' in production (use 'require' or 'verify-full')",
			}
		}
	}

	return nil
}

// getValueForEnvVar gets the current value for an environment variable
func (c *Config) getValueForEnvVar(name string) string {
	switch name {
	case "ENVIRONMENT":
		return c.Environment
	case "PORT":
		return c.Port
	case "LOG_LEVEL":
		return c.LogLevel
	case "FRONTEND_URL":
		return c.FrontendURL
	case "DB_PASSWORD":
		return c.Database.Password
	case "DB_SSL_MODE":
		return c.Database.SSLMode
	default:
		// For env-only variables (secrets), read directly from environment
		return os.Getenv(name)
	}
}

// =============================================================================
// Validation Helpers
// =============================================================================

// MustValidate validates the configuration and logs fatal if there are errors
func (c *Config) MustValidate() {
	errs := c.Validate()
	if errs.HasErrors() {
		log.Fatalf("Configuration validation failed:\n  %s", strings.ReplaceAll(errs.Error(), "; ", "\n  "))
	}
}

// ValidateWithWarnings validates and logs warnings for non-critical issues
func (c *Config) ValidateWithWarnings() error {
	errs := c.Validate()

	// In development, log warnings for production requirements
	if c.Environment == "development" {
		for _, ev := range envRegistry {
			if ev.Required == RequiredProduction {
				if c.getValueForEnvVar(ev.Name) == "" {
					log.Printf("Warning: %s not set (%s)", ev.Name, ev.Description)
				}
			}
		}
	}

	// Filter out production-only errors in non-production
	if c.Environment != "production" {
		var criticalErrs ValidationErrors
		for _, err := range errs {
			if !isProductionOnlyError(err.Field) {
				criticalErrs = append(criticalErrs, err)
			}
		}
		if criticalErrs.HasErrors() {
			return criticalErrs
		}
		return nil
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

// isProductionOnlyError checks if an error field is production-only
func isProductionOnlyError(field string) bool {
	for _, ev := range envRegistry {
		if ev.Name == field && ev.Required == RequiredProduction {
			return true
		}
	}
	// Also filter HTTPS/SSL warnings in development
	return field == "FRONTEND_URL" || field == "DB_SSL_MODE"
}

// GetRegistry returns the parsed environment variable registry (for testing)
func GetRegistry() []EnvVar {
	return envRegistry
}

// =============================================================================
// Utility Functions
// =============================================================================

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Warning: Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if value == "true" || value == "1" || value == "yes" {
			return true
		}
		if value == "false" || value == "0" || value == "no" {
			return false
		}
		log.Printf("Warning: Invalid boolean value for %s: %s, using default: %v", key, value, defaultValue)
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue string) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("Warning: Invalid duration value for %s: %s, using default: %s", key, value, defaultValue)
	}
	duration, _ := time.ParseDuration(defaultValue)
	return duration
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func requiredInMsg(r RequiredIn) string {
	switch r {
	case RequiredProduction:
		return " in production"
	default:
		return ""
	}
}

// =============================================================================
// Documentation Helper
// =============================================================================

// PrintEnvDocs prints documentation for all registered environment variables
func PrintEnvDocs() {
	fmt.Println("Environment Variables (from .env.example):")
	fmt.Println("==========================================")
	fmt.Println()

	for _, ev := range envRegistry {
		required := "optional"
		switch ev.Required {
		case RequiredAlways:
			required = "required"
		case RequiredProduction:
			required = "required in production"
		}

		fmt.Printf("%-20s %s\n", ev.Name, ev.Description)
		fmt.Printf("%-20s Required: %s\n", "", required)
		if ev.Default != "" {
			fmt.Printf("%-20s Default: %s\n", "", ev.Default)
		}
		if len(ev.AllowedVals) > 0 {
			fmt.Printf("%-20s Allowed: %s\n", "", strings.Join(ev.AllowedVals, ", "))
		}
		fmt.Println()
	}
}
