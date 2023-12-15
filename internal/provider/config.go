package provider

import "os"

type DefaultUser struct {
	ID        string
	UserID    string
	Email     string
	FirstName string
	LastName  string
}

type Config struct {
	Port            string
	PublicBaseURL   string
	ProviderName    string
	PackageVersion  string
	JWTSecret       string
	TokenQueryParam string
	SessionParam    string
	DefaultUser     DefaultUser
}

func LoadConfig() Config {
	return Config{
		Port:            env("PORT", "8080"),
		PublicBaseURL:   os.Getenv("PUBLIC_BASE_URL"),
		ProviderName:    env("PROVIDER_NAME", "Tata Consulting"),
		PackageVersion:  env("PACKAGE_VERSION", "v0.1.0"),
		JWTSecret:       env("JWT_SECRET", "change-me-for-real-deployments"),
		TokenQueryParam: env("TOKEN_QUERY_PARAM", "token"),
		SessionParam:    env("SESSION_QUERY_PARAM", "session_cookie"),
		DefaultUser: DefaultUser{
			ID:        env("DEFAULT_USER_ID", "0f6f66c7-68f4-4b11-8d9a-d5f27f95ad8e"),
			UserID:    env("DEFAULT_USER_HANDLE", "hamza-mohd"),
			Email:     env("DEFAULT_USER_EMAIL", "hamza.mohd@tata-consulting.example"),
			FirstName: env("DEFAULT_USER_FIRST_NAME", "Mohd"),
			LastName:  env("DEFAULT_USER_LAST_NAME", "Hamza Shaikh"),
		},
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
