package app

import (
	"log/slog"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Lookup membatasi pembacaan konfigurasi ke nama yang dipakai fondasi HTTP.
type Lookup func(string) (string, bool)

type Config struct {
	Environment    string
	AppOrigin      string
	HTTPAddr       string
	RequestTimeout time.Duration
	LogLevel       slog.Level
}

// ConfigError hanya menyimpan nama field; nilai environment tidak masuk error/log.
type ConfigError struct{ Field string }

func (e *ConfigError) Error() string { return "Konfigurasi tidak valid: " + e.Field }

func LoadConfig(lookup Lookup) (Config, error) {
	cfg := Config{
		HTTPAddr:       "127.0.0.1:8080",
		RequestTimeout: 10 * time.Second,
		LogLevel:       slog.LevelInfo,
	}
	cfg.Environment, _ = lookup("APP_ENV")
	cfg.AppOrigin, _ = lookup("APP_ORIGIN")
	if marker, ok := lookup("VERCEL_ENV"); ok && (marker == "" || marker != cfg.Environment) {
		return Config{}, &ConfigError{Field: "VERCEL_ENV"}
	}
	if addr, ok := lookup("HTTP_ADDR"); ok {
		cfg.HTTPAddr = addr
	}
	if timeout, ok := lookup("REQUEST_TIMEOUT"); ok {
		value, err := time.ParseDuration(timeout)
		if err != nil {
			return Config{}, &ConfigError{Field: "REQUEST_TIMEOUT"}
		}
		cfg.RequestTimeout = value
	}
	if level, ok := lookup("LOG_LEVEL"); ok {
		switch level {
		case "debug":
			cfg.LogLevel = slog.LevelDebug
		case "info":
			cfg.LogLevel = slog.LevelInfo
		case "warn":
			cfg.LogLevel = slog.LevelWarn
		case "error":
			cfg.LogLevel = slog.LevelError
		default:
			return Config{}, &ConfigError{Field: "LOG_LEVEL"}
		}
	}
	// Cron/push belum diimplementasikan. Jangan menerima konfigurasi yang seolah
	// mengaktifkannya, pada environment mana pun, sebelum milestone terkait.
	for _, key := range []string{"PUSH_ENABLED", "CRON_ENABLED"} {
		if value, ok := lookup(key); ok && value != "false" {
			return Config{}, &ConfigError{Field: key}
		}
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	switch c.Environment {
	case "development", "preview", "production":
	default:
		return &ConfigError{Field: "APP_ENV"}
	}
	u, err := url.Parse(c.AppOrigin)
	if err != nil || u == nil || u.Hostname() == "" || u.User != nil ||
		u.Opaque != "" || u.Path != "" || u.RawQuery != "" || u.ForceQuery ||
		u.Fragment != "" || u.RawFragment != "" ||
		(u.Scheme != "http" && u.Scheme != "https") ||
		strings.ContainsAny(c.AppOrigin, " \t\r\n\\") {
		return &ConfigError{Field: "APP_ORIGIN"}
	}
	if port := u.Port(); port != "" {
		n, err := strconv.Atoi(port)
		if err != nil || n < 1 || n > 65535 {
			return &ConfigError{Field: "APP_ORIGIN"}
		}
	}
	loopback := strings.EqualFold(u.Hostname(), "localhost")
	if ip := net.ParseIP(u.Hostname()); ip != nil {
		if ip.String() != u.Hostname() {
			return &ConfigError{Field: "APP_ORIGIN"}
		}
		loopback = ip.IsLoopback()
	} else if !validHostname(u.Hostname()) {
		return &ConfigError{Field: "APP_ORIGIN"}
	}
	if (c.Environment == "development" && !loopback) ||
		(c.Environment != "development" && (u.Scheme != "https" || loopback)) {
		return &ConfigError{Field: "APP_ORIGIN"}
	}
	// Simpan origin dalam bentuk canonical yang persis dikirim browser. Menolak
	// default port eksplisit/trailing slash menghindari konfigurasi yang tak cocok.
	if u.String() != c.AppOrigin || u.Host != strings.ToLower(u.Host) ||
		(u.Scheme == "http" && u.Port() == "80") ||
		(u.Scheme == "https" && u.Port() == "443") || strings.HasSuffix(u.Host, ":") {
		return &ConfigError{Field: "APP_ORIGIN"}
	}
	host, port, err := net.SplitHostPort(c.HTTPAddr)
	ip := net.ParseIP(host)
	n, portErr := strconv.Atoi(port)
	if err != nil || ip == nil || !ip.IsLoopback() || portErr != nil || n < 0 || n > 65535 {
		return &ConfigError{Field: "HTTP_ADDR"}
	}
	if c.RequestTimeout < time.Second || c.RequestTimeout > 30*time.Second {
		return &ConfigError{Field: "REQUEST_TIMEOUT"}
	}
	if c.LogLevel != slog.LevelDebug && c.LogLevel != slog.LevelInfo &&
		c.LogLevel != slog.LevelWarn && c.LogLevel != slog.LevelError {
		return &ConfigError{Field: "LOG_LEVEL"}
	}
	return nil
}

func configField(err error) string {
	if field, ok := err.(*ConfigError); ok {
		return field.Field
	}
	return "CONFIG"
}

func validHostname(host string) bool {
	if len(host) > 253 {
		return false
	}
	for _, label := range strings.Split(host, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, char := range label {
			if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') || char == '-') {
				return false
			}
		}
	}
	return true
}
