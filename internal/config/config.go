package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Port        string `mapstructure:"port"`
	DatabaseURL string `mapstructure:"database_url"`
	AppEnv      string `mapstructure:"app_env"`
	Auth        AuthConfig
	Database    DatabaseConfig
}

type AuthConfig struct {
	JWTSecret              []byte
	CustomerPortalOrigin   string        `mapstructure:"customer_portal_origin"`
	StaffPortalOrigin      string        `mapstructure:"staff_portal_origin"`
	CustomerAccessTokenTTL time.Duration `mapstructure:"access_token_ttl_customer"`
	StaffAccessTokenTTL    time.Duration `mapstructure:"access_token_ttl_staff"`
	RefreshTokenTTL        time.Duration `mapstructure:"refresh_token_ttl"`
	CookieDomain           string        `mapstructure:"cookie_domain"`
	CookieSecure           bool          `mapstructure:"cookie_secure"`
	CookieSameSite         string        `mapstructure:"cookie_same_site"`
	JWTClockSkew           time.Duration `mapstructure:"jwt_clock_skew"`
	RateLimitWindow        time.Duration `mapstructure:"rate_limit_window"`
	CustomerLoginLimit     int           `mapstructure:"rate_limit_customer_login"`
	StaffLoginLimit        int           `mapstructure:"rate_limit_staff_login"`
	CustomerRegisterLimit  int           `mapstructure:"rate_limit_customer_register"`
}

type DatabaseConfig struct {
	MaxConns        int32         `mapstructure:"max_conns"`
	MinConns        int32         `mapstructure:"min_conns"`
	MaxConnLifetime time.Duration `mapstructure:"max_conn_lifetime"`
	MaxConnIdleTime time.Duration `mapstructure:"max_conn_idle_time"`
}

func Load() (Config, error) {
	v := viper.New()

	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")

	for _, key := range configKeys() {
		v.RegisterAlias(key, envKey(key))
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return Config{}, fmt.Errorf("read config: %w", err)
		}
	}

	setDefaults(v)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	for _, key := range configKeys() {
		if err := v.BindEnv(key); err != nil {
			return Config{}, fmt.Errorf("bind env %s: %w", key, err)
		}
	}

	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	if cfg.Port == "" {
		return Config{}, errors.New("config: PORT is required")
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("config: DATABASE_URL is required")
	}
	if cfg.AppEnv == "" {
		return Config{}, errors.New("config: APP_ENV is required")
	}
	if cfg.Auth.CustomerPortalOrigin == "" {
		return Config{}, errors.New("config: AUTH_CUSTOMER_PORTAL_ORIGIN is required")
	}
	if cfg.Auth.StaffPortalOrigin == "" {
		return Config{}, errors.New("config: AUTH_STAFF_PORTAL_ORIGIN is required")
	}

	jwtSecret, err := base64.RawURLEncoding.DecodeString(v.GetString("auth.jwt_secret"))
	if err != nil {
		return Config{}, fmt.Errorf("config: AUTH_JWT_SECRET must be a base64url without padding: %w", err)
	}
	if len(jwtSecret) < 32 {
		return Config{}, errors.New("config: AUTH_JWT_SECRET must decoded to at least 32 bytes")
	}
	cfg.Auth.JWTSecret = jwtSecret

	return cfg, nil
}

func configKeys() []string {
	return []string{
		"port",
		"database_url",
		"app_env",
		"auth.jwt_secret",
		"auth.customer_portal_origin",
		"auth.staff_portal_origin",
		"auth.access_token_ttl_customer",
		"auth.access_token_ttl_staff",
		"auth.refresh_token_ttl",
		"auth.cookie_domain",
		"auth.cookie_secure",
		"auth.cookie_same_site",
		"auth.jwt_clock_skew",
		"auth.rate_limit_window",
		"auth.rate_limit_customer_login",
		"auth.rate_limit_staff_login",
		"auth.rate_limit_customer_register",
		"database.max_conns",
		"database.min_conns",
		"database.max_conn_lifetime",
		"database.max_conn_idle_time",
	}
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("auth.access_token_ttl_customer", 15*time.Minute)
	v.SetDefault("auth.access_token_ttl_staff", 5*time.Minute)
	v.SetDefault("auth.refresh_token_ttl", 30*24*time.Hour)
	v.SetDefault("auth.cookie_domain", "")
	v.SetDefault("auth.cookie_secure", true)
	v.SetDefault("auth.cookie_same_site", "strict")
	v.SetDefault("auth.jwt_clock_skew", 10*time.Second)
	v.SetDefault("auth.rate_limit_window", time.Minute)
	v.SetDefault("auth.rate_limit_customer_login", 10)
	v.SetDefault("auth.rate_limit_staff_login", 5)
	v.SetDefault("auth.rate_limit_customer_register", 5)
	v.SetDefault("database.max_conns", 10)
	v.SetDefault("database.min_conns", 1)
	v.SetDefault("database.max_conn_lifetime", time.Hour)
	v.SetDefault("database.max_conn_idle_time", 15*time.Minute)
}

func envKey(key string) string {
	return strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
}
