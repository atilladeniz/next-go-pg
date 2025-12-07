package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseEnvExampleContent tests the parsing of .env.example content
func TestParseEnvExampleContent(t *testing.T) {
	content := `
# @required @type:enum|dev,prod
# Environment
ENV=dev

# @required:prod @type:secret @min:16
# API Key
API_KEY=test-key

# @type:port
# Server port
PORT=8080

# No validation tags - should not be in registry
OPTIONAL_VAR=value
`

	vars, err := parseEnvExampleContent(content)
	require.NoError(t, err)

	// Should only have 3 vars (OPTIONAL_VAR has no validation)
	assert.Len(t, vars, 3)

	// Check ENV
	env := findVar(vars, "ENV")
	require.NotNil(t, env)
	assert.Equal(t, RequiredAlways, env.Required)
	assert.Equal(t, TypeEnum, env.Type)
	assert.Equal(t, []string{"dev", "prod"}, env.AllowedVals)
	assert.Equal(t, "dev", env.Default)

	// Check API_KEY
	apiKey := findVar(vars, "API_KEY")
	require.NotNil(t, apiKey)
	assert.Equal(t, RequiredProduction, apiKey.Required)
	assert.Equal(t, TypeSecret, apiKey.Type)
	assert.Equal(t, 16, apiKey.MinLength)

	// Check PORT
	port := findVar(vars, "PORT")
	require.NotNil(t, port)
	assert.Equal(t, RequiredNever, port.Required)
	assert.Equal(t, TypePort, port.Type)
}

// TestEnvExampleParsing tests that .env.example is properly parsed from filesystem
func TestEnvExampleParsing(t *testing.T) {
	registry := GetRegistry()

	// Should have parsed vars from .env.example file
	assert.NotEmpty(t, registry, "Registry should not be empty")

	// Check for expected variables
	expectedVars := []string{"ENVIRONMENT", "PORT", "LOG_LEVEL", "FRONTEND_URL", "GO_JWT_SECRET", "WEBHOOK_SECRET", "DB_PASSWORD", "DB_SSL_MODE"}
	for _, name := range expectedVars {
		ev := findVar(registry, name)
		assert.NotNil(t, ev, "Expected %s in registry", name)
	}
}

// TestConfigLoad tests configuration loading from environment
func TestConfigLoad(t *testing.T) {
	os.Setenv("PORT", "9090")
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("FRONTEND_URL", "https://example.com")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSL_MODE", "require")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("ENVIRONMENT")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("FRONTEND_URL")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSL_MODE")
	}()

	cfg := Load()

	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "staging", cfg.Environment)
	assert.Equal(t, "debug", cfg.LogLevel)
	assert.Equal(t, "https://example.com", cfg.FrontendURL)
	assert.Equal(t, "db.example.com", cfg.Database.Host)
	assert.Equal(t, "5433", cfg.Database.Port)
	assert.Equal(t, "testuser", cfg.Database.User)
	assert.Equal(t, "testpass", cfg.Database.Password)
	assert.Equal(t, "testdb", cfg.Database.Name)
	assert.Equal(t, "require", cfg.Database.SSLMode)
}

// TestConfigDefaults tests default configuration values
func TestConfigDefaults(t *testing.T) {
	envVars := []string{
		"PORT", "ENVIRONMENT", "LOG_LEVEL", "FRONTEND_URL",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL_MODE",
	}
	for _, v := range envVars {
		os.Unsetenv(v)
	}

	cfg := Load()

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, "http://localhost:3000", cfg.FrontendURL)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "5432", cfg.Database.Port)
	assert.Equal(t, "postgres", cfg.Database.User)
	assert.Equal(t, "", cfg.Database.Password)
	assert.Equal(t, "nextgopg", cfg.Database.Name)
	assert.Equal(t, "disable", cfg.Database.SSLMode)
}

// TestValidateEnvironment tests environment validation
func TestValidateEnvironment(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		hasError    bool
	}{
		{"development", "development", false},
		{"staging", "staging", false},
		{"production", "production", false},
		{"invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Port:        "8080",
				Environment: tt.environment,
				LogLevel:    "info",
				FrontendURL: "http://localhost:3000",
				Database:    DatabaseConfig{SSLMode: "disable"},
			}

			errs := cfg.Validate()
			hasEnvError := hasErrorForField(errs, "ENVIRONMENT")
			assert.Equal(t, tt.hasError, hasEnvError)
		})
	}
}

// TestValidateLogLevel tests log level validation
func TestValidateLogLevel(t *testing.T) {
	validLevels := []string{"trace", "debug", "info", "warn", "error", "fatal"}
	invalidLevels := []string{"invalid", "verbose"}

	for _, level := range validLevels {
		t.Run("valid_"+level, func(t *testing.T) {
			cfg := &Config{
				Port:        "8080",
				Environment: "development",
				LogLevel:    level,
				FrontendURL: "http://localhost:3000",
				Database:    DatabaseConfig{SSLMode: "disable"},
			}

			errs := cfg.Validate()
			assert.False(t, hasErrorForField(errs, "LOG_LEVEL"), "Expected no error for valid level: %s", level)
		})
	}

	for _, level := range invalidLevels {
		t.Run("invalid_"+level, func(t *testing.T) {
			cfg := &Config{
				Port:        "8080",
				Environment: "development",
				LogLevel:    level,
				FrontendURL: "http://localhost:3000",
				Database:    DatabaseConfig{SSLMode: "disable"},
			}

			errs := cfg.Validate()
			assert.True(t, hasErrorForField(errs, "LOG_LEVEL"), "Expected error for invalid level: %s", level)
		})
	}
}

// TestValidatePort tests port validation
func TestValidatePort(t *testing.T) {
	tests := []struct {
		name     string
		port     string
		hasError bool
	}{
		{"valid port", "8080", false},
		{"port 1", "1", false},
		{"port 65535", "65535", false},
		{"port 0", "0", true},
		{"port 65536", "65536", true},
		{"negative port", "-1", true},
		{"non-numeric", "abc", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Port:        tt.port,
				Environment: "development",
				LogLevel:    "info",
				FrontendURL: "http://localhost:3000",
				Database:    DatabaseConfig{SSLMode: "disable"},
			}

			errs := cfg.Validate()
			assert.Equal(t, tt.hasError, hasErrorForField(errs, "PORT"))
		})
	}
}

// TestValidateFrontendURL tests frontend URL validation
func TestValidateFrontendURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		hasError bool
	}{
		{"valid http", "http://localhost:3000", false},
		{"valid https", "https://example.com", false},
		{"with path", "https://example.com/app", false},
		{"invalid url", "not-a-url", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Port:        "8080",
				Environment: "development",
				LogLevel:    "info",
				FrontendURL: tt.url,
				Database:    DatabaseConfig{SSLMode: "disable"},
			}

			errs := cfg.Validate()
			assert.Equal(t, tt.hasError, hasErrorForField(errs, "FRONTEND_URL"))
		})
	}
}

// TestValidateProduction tests production-specific validations
func TestValidateProduction(t *testing.T) {
	t.Run("production requires DB_PASSWORD", func(t *testing.T) {
		cfg := &Config{
			Port:        "8080",
			Environment: "production",
			LogLevel:    "info",
			FrontendURL: "https://example.com",
			Database: DatabaseConfig{
				Password: "",
				SSLMode:  "require",
			},
		}
		os.Setenv("GO_JWT_SECRET", "this-is-a-very-long-secret-key-for-jwt")
		os.Setenv("WEBHOOK_SECRET", "webhook-secret-16ch")
		defer func() {
			os.Unsetenv("GO_JWT_SECRET")
			os.Unsetenv("WEBHOOK_SECRET")
		}()

		errs := cfg.Validate()
		assert.True(t, hasErrorForField(errs, "DB_PASSWORD"))
	})

	t.Run("production requires GO_JWT_SECRET", func(t *testing.T) {
		cfg := &Config{
			Port:        "8080",
			Environment: "production",
			LogLevel:    "info",
			FrontendURL: "https://example.com",
			Database: DatabaseConfig{
				Password: "password",
				SSLMode:  "require",
			},
		}
		os.Unsetenv("GO_JWT_SECRET")
		os.Setenv("WEBHOOK_SECRET", "webhook-secret-16ch")
		defer os.Unsetenv("WEBHOOK_SECRET")

		errs := cfg.Validate()
		assert.True(t, hasErrorForField(errs, "GO_JWT_SECRET"))
	})

	t.Run("production requires WEBHOOK_SECRET", func(t *testing.T) {
		cfg := &Config{
			Port:        "8080",
			Environment: "production",
			LogLevel:    "info",
			FrontendURL: "https://example.com",
			Database: DatabaseConfig{
				Password: "password",
				SSLMode:  "require",
			},
		}
		os.Setenv("GO_JWT_SECRET", "this-is-a-very-long-secret-key-for-jwt")
		os.Unsetenv("WEBHOOK_SECRET")
		defer os.Unsetenv("GO_JWT_SECRET")

		errs := cfg.Validate()
		assert.True(t, hasErrorForField(errs, "WEBHOOK_SECRET"))
	})

	t.Run("production warns about disabled SSL", func(t *testing.T) {
		cfg := &Config{
			Port:        "8080",
			Environment: "production",
			LogLevel:    "info",
			FrontendURL: "https://example.com",
			Database: DatabaseConfig{
				Password: "password",
				SSLMode:  "disable",
			},
		}
		os.Setenv("GO_JWT_SECRET", "this-is-a-very-long-secret-key-for-jwt")
		os.Setenv("WEBHOOK_SECRET", "webhook-secret-16ch")
		defer func() {
			os.Unsetenv("GO_JWT_SECRET")
			os.Unsetenv("WEBHOOK_SECRET")
		}()

		errs := cfg.Validate()
		assert.True(t, hasErrorForField(errs, "DB_SSL_MODE"))
	})

	t.Run("production warns about non-HTTPS frontend", func(t *testing.T) {
		cfg := &Config{
			Port:        "8080",
			Environment: "production",
			LogLevel:    "info",
			FrontendURL: "http://example.com",
			Database: DatabaseConfig{
				Password: "password",
				SSLMode:  "require",
			},
		}
		os.Setenv("GO_JWT_SECRET", "this-is-a-very-long-secret-key-for-jwt")
		os.Setenv("WEBHOOK_SECRET", "webhook-secret-16ch")
		defer func() {
			os.Unsetenv("GO_JWT_SECRET")
			os.Unsetenv("WEBHOOK_SECRET")
		}()

		errs := cfg.Validate()
		assert.True(t, hasErrorForField(errs, "FRONTEND_URL"))
	})

	t.Run("valid production config", func(t *testing.T) {
		cfg := &Config{
			Port:        "8080",
			Environment: "production",
			LogLevel:    "info",
			FrontendURL: "https://example.com",
			Database: DatabaseConfig{
				Password: "password",
				SSLMode:  "require",
			},
		}
		os.Setenv("GO_JWT_SECRET", "this-is-a-very-long-secret-key-for-jwt")
		os.Setenv("WEBHOOK_SECRET", "webhook-secret-16ch")
		defer func() {
			os.Unsetenv("GO_JWT_SECRET")
			os.Unsetenv("WEBHOOK_SECRET")
		}()

		errs := cfg.Validate()
		assert.False(t, errs.HasErrors(), "Expected no errors, got: %v", errs)
	})
}

// TestSecretMinLength tests that secrets have minimum length requirements
func TestSecretMinLength(t *testing.T) {
	cfg := &Config{
		Port:        "8080",
		Environment: "production",
		LogLevel:    "info",
		FrontendURL: "https://example.com",
		Database: DatabaseConfig{
			Password: "password",
			SSLMode:  "require",
		},
	}

	t.Run("JWT secret too short", func(t *testing.T) {
		os.Setenv("GO_JWT_SECRET", "short")
		os.Setenv("WEBHOOK_SECRET", "webhook-secret-16ch")
		defer func() {
			os.Unsetenv("GO_JWT_SECRET")
			os.Unsetenv("WEBHOOK_SECRET")
		}()

		errs := cfg.Validate()
		hasLengthError := false
		for _, err := range errs {
			if err.Field == "GO_JWT_SECRET" && err.Message == "must be at least 32 characters" {
				hasLengthError = true
				break
			}
		}
		assert.True(t, hasLengthError)
	})

	t.Run("Webhook secret too short", func(t *testing.T) {
		os.Setenv("GO_JWT_SECRET", "this-is-a-very-long-secret-key-for-jwt")
		os.Setenv("WEBHOOK_SECRET", "short")
		defer func() {
			os.Unsetenv("GO_JWT_SECRET")
			os.Unsetenv("WEBHOOK_SECRET")
		}()

		errs := cfg.Validate()
		hasLengthError := false
		for _, err := range errs {
			if err.Field == "WEBHOOK_SECRET" && err.Message == "must be at least 16 characters" {
				hasLengthError = true
				break
			}
		}
		assert.True(t, hasLengthError)
	})
}

// TestValidationErrors tests the ValidationErrors type
func TestValidationErrors(t *testing.T) {
	t.Run("empty errors", func(t *testing.T) {
		var errs ValidationErrors
		assert.False(t, errs.HasErrors())
		assert.Empty(t, errs.Error())
	})

	t.Run("single error", func(t *testing.T) {
		errs := ValidationErrors{
			{Field: "PORT", Message: "is required"},
		}
		assert.True(t, errs.HasErrors())
		assert.Equal(t, "PORT: is required", errs.Error())
	})

	t.Run("multiple errors", func(t *testing.T) {
		errs := ValidationErrors{
			{Field: "PORT", Message: "is required"},
			{Field: "DATABASE", Message: "connection failed"},
		}
		assert.True(t, errs.HasErrors())
		assert.Contains(t, errs.Error(), "PORT: is required")
		assert.Contains(t, errs.Error(), "DATABASE: connection failed")
	})
}

// TestValidateWithWarnings tests development mode warnings
func TestValidateWithWarnings(t *testing.T) {
	t.Run("development mode - no critical errors", func(t *testing.T) {
		cfg := &Config{
			Port:        "8080",
			Environment: "development",
			LogLevel:    "info",
			FrontendURL: "http://localhost:3000",
			Database:    DatabaseConfig{SSLMode: "disable"},
		}
		os.Unsetenv("GO_JWT_SECRET")
		os.Unsetenv("WEBHOOK_SECRET")

		err := cfg.ValidateWithWarnings()
		assert.NoError(t, err)
	})

	t.Run("development mode - critical error", func(t *testing.T) {
		cfg := &Config{
			Port:        "invalid",
			Environment: "development",
			LogLevel:    "info",
			FrontendURL: "http://localhost:3000",
			Database:    DatabaseConfig{SSLMode: "disable"},
		}

		err := cfg.ValidateWithWarnings()
		assert.Error(t, err)
	})
}

// TestGetDatabaseURL tests database URL generation
func TestGetDatabaseURL(t *testing.T) {
	cfg := &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "secret",
			Name:     "testdb",
			SSLMode:  "disable",
		},
	}

	url := cfg.GetDatabaseURL()
	assert.Contains(t, url, "host=localhost")
	assert.Contains(t, url, "port=5432")
	assert.Contains(t, url, "user=postgres")
	assert.Contains(t, url, "password=secret")
	assert.Contains(t, url, "dbname=testdb")
	assert.Contains(t, url, "sslmode=disable")
}

// TestSSLModeValidation tests SSL mode validation
func TestSSLModeValidation(t *testing.T) {
	validModes := []string{"disable", "allow", "prefer", "require", "verify-ca", "verify-full"}

	for _, mode := range validModes {
		t.Run("valid_"+mode, func(t *testing.T) {
			cfg := &Config{
				Port:        "8080",
				Environment: "development",
				LogLevel:    "info",
				FrontendURL: "http://localhost:3000",
				Database:    DatabaseConfig{SSLMode: mode},
			}

			errs := cfg.Validate()
			// Check for SSL mode format errors (not production warnings)
			hasSSLFormatError := false
			for _, err := range errs {
				if err.Field == "DB_SSL_MODE" && !contains([]string{"should not be 'disable' in production (use 'require' or 'verify-full')"}, err.Message) {
					hasSSLFormatError = true
					break
				}
			}
			assert.False(t, hasSSLFormatError, "Expected no format error for valid SSL mode: %s", mode)
		})
	}

	t.Run("invalid SSL mode", func(t *testing.T) {
		cfg := &Config{
			Port:        "8080",
			Environment: "development",
			LogLevel:    "info",
			FrontendURL: "http://localhost:3000",
			Database:    DatabaseConfig{SSLMode: "invalid-mode"},
		}

		errs := cfg.Validate()
		assert.True(t, hasErrorForField(errs, "DB_SSL_MODE"))
	})
}

// TestGetEnvHelpers tests helper functions
func TestGetEnvHelpers(t *testing.T) {
	t.Run("getEnv with value", func(t *testing.T) {
		os.Setenv("TEST_VAR", "test-value")
		defer os.Unsetenv("TEST_VAR")

		result := getEnv("TEST_VAR", "default")
		assert.Equal(t, "test-value", result)
	})

	t.Run("getEnv without value", func(t *testing.T) {
		os.Unsetenv("TEST_VAR_NOT_SET")
		result := getEnv("TEST_VAR_NOT_SET", "default")
		assert.Equal(t, "default", result)
	})

	t.Run("getEnvAsInt with valid int", func(t *testing.T) {
		os.Setenv("TEST_INT", "42")
		defer os.Unsetenv("TEST_INT")

		result := getEnvAsInt("TEST_INT", 0)
		assert.Equal(t, 42, result)
	})

	t.Run("getEnvAsInt with invalid int", func(t *testing.T) {
		os.Setenv("TEST_INT_INVALID", "not-a-number")
		defer os.Unsetenv("TEST_INT_INVALID")

		result := getEnvAsInt("TEST_INT_INVALID", 100)
		assert.Equal(t, 100, result)
	})

	t.Run("getEnvAsBool with true", func(t *testing.T) {
		for _, val := range []string{"true", "1", "yes"} {
			os.Setenv("TEST_BOOL", val)
			result := getEnvAsBool("TEST_BOOL", false)
			assert.True(t, result, "expected true for value: %s", val)
		}
		os.Unsetenv("TEST_BOOL")
	})

	t.Run("getEnvAsBool with false", func(t *testing.T) {
		for _, val := range []string{"false", "0", "no"} {
			os.Setenv("TEST_BOOL", val)
			result := getEnvAsBool("TEST_BOOL", true)
			assert.False(t, result, "expected false for value: %s", val)
		}
		os.Unsetenv("TEST_BOOL")
	})

	t.Run("getEnvAsDuration with valid duration", func(t *testing.T) {
		os.Setenv("TEST_DURATION", "30s")
		defer os.Unsetenv("TEST_DURATION")

		result := getEnvAsDuration("TEST_DURATION", "10s")
		assert.Equal(t, "30s", result.String())
	})

	t.Run("getEnvAsDuration with invalid duration", func(t *testing.T) {
		os.Setenv("TEST_DURATION_INVALID", "invalid")
		defer os.Unsetenv("TEST_DURATION_INVALID")

		result := getEnvAsDuration("TEST_DURATION_INVALID", "10s")
		assert.Equal(t, "10s", result.String())
	})
}

// TestContainsHelper tests the contains helper function
func TestContainsHelper(t *testing.T) {
	slice := []string{"a", "b", "c"}

	assert.True(t, contains(slice, "a"))
	assert.True(t, contains(slice, "b"))
	assert.True(t, contains(slice, "c"))
	assert.False(t, contains(slice, "d"))
	assert.False(t, contains(slice, ""))
	assert.False(t, contains(nil, "a"))
	assert.False(t, contains([]string{}, "a"))
}

// TestRequiredInMsg tests the requiredInMsg helper
func TestRequiredInMsg(t *testing.T) {
	assert.Equal(t, "", requiredInMsg(RequiredNever))
	assert.Equal(t, " in production", requiredInMsg(RequiredProduction))
	assert.Equal(t, "", requiredInMsg(RequiredAlways))
}

// Helper functions for tests

func findVar(vars []EnvVar, name string) *EnvVar {
	for i := range vars {
		if vars[i].Name == name {
			return &vars[i]
		}
	}
	return nil
}

func hasErrorForField(errs ValidationErrors, field string) bool {
	for _, err := range errs {
		if err.Field == field {
			return true
		}
	}
	return false
}
