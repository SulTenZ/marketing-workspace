package app

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func lookupMap(values map[string]string) Lookup {
	return func(key string) (string, bool) { value, ok := values[key]; return value, ok }
}

func devEnv() map[string]string {
	return map[string]string{"APP_ENV": "development", "APP_ORIGIN": "http://127.0.0.1:5173"}
}

func TestConfigValidEnvironments(t *testing.T) {
	for _, env := range []string{"development", "preview", "production"} {
		t.Run(env, func(t *testing.T) {
			values := devEnv()
			values["APP_ENV"], values["VERCEL_ENV"] = env, env
			if env != "development" {
				values["APP_ORIGIN"] = "https://web.example.invalid"
			}
			cfg, err := LoadConfig(lookupMap(values))
			if err != nil {
				t.Fatal(err)
			}
			if cfg.Environment != env || cfg.AppOrigin != values["APP_ORIGIN"] ||
				cfg.RequestTimeout != 10*time.Second || cfg.HTTPAddr != "127.0.0.1:8080" {
				t.Fatalf("configuration not preserved: %+v", cfg)
			}
		})
	}
}

func TestConfigFailsClosed(t *testing.T) {
	cases := []struct{ field, value string }{
		{"APP_ENV", ""}, {"APP_ENV", "staging"}, {"APP_ENV", "production"},
		{"APP_ORIGIN", ""}, {"APP_ORIGIN", "*"}, {"APP_ORIGIN", "null"},
		{"APP_ORIGIN", "http://localhost:5173/"}, {"APP_ORIGIN", "http://user:private@localhost"},
		{"APP_ORIGIN", "http://localhost?token=private"}, {"APP_ORIGIN", "http://localhost#private"},
		{"APP_ORIGIN", "http://localhost#"}, {"APP_ORIGIN", "http://localhost:65536"},
		{"APP_ORIGIN", "http://localhost:80"}, {"APP_ORIGIN", "http://localhost:"},
		{"APP_ORIGIN", "http://LOCALHOST:5173"}, {"APP_ORIGIN", "http://localhost:5173\n"},
		{"APP_ORIGIN", "https://web.example.invalid"}, {"APP_ORIGIN", "http://127.1:5173"},
		{"HTTP_ADDR", "0.0.0.0:8080"}, {"HTTP_ADDR", "private.invalid:8080"},
		{"HTTP_ADDR", "127.0.0.1:65536"}, {"HTTP_ADDR", ""},
		{"REQUEST_TIMEOUT", ""}, {"REQUEST_TIMEOUT", "999h"}, {"REQUEST_TIMEOUT", "0s"},
		{"REQUEST_TIMEOUT", "1ms"}, {"REQUEST_TIMEOUT", "invalid-private-value"},
		{"LOG_LEVEL", "private-value"}, {"LOG_LEVEL", ""},
		{"VERCEL_ENV", "preview"}, {"VERCEL_ENV", ""},
		{"PUSH_ENABLED", "true"}, {"CRON_ENABLED", "true"}, {"PUSH_ENABLED", ""},
	}
	for _, test := range cases {
		t.Run(test.field+"/"+test.value, func(t *testing.T) {
			values := devEnv()
			values[test.field] = test.value
			cfg, err := LoadConfig(lookupMap(values))
			var invalid *ConfigError
			if !errors.As(err, &invalid) || cfg != (Config{}) {
				t.Fatalf("invalid config accepted: %+v, %v", cfg, err)
			}
			if strings.Contains(err.Error(), "private") {
				t.Fatalf("configuration value leaked: %v", err)
			}
		})
	}
}

func TestNonDevelopmentRequiresHTTPSAndCanonicalHost(t *testing.T) {
	for _, origin := range []string{
		"http://web.example.invalid", "https://localhost", "https://[::1]", "https://*",
		"https://foo..invalid", "https://bad_name.invalid", "https://web.example.invalid.",
		"https://web.example.invalid:443", "https://[::1%25private]",
	} {
		t.Run(origin, func(t *testing.T) {
			values := map[string]string{"APP_ENV": "preview", "APP_ORIGIN": origin}
			if _, err := LoadConfig(lookupMap(values)); err == nil {
				t.Fatal("unsafe or non-canonical origin accepted")
			}
		})
	}
}

func TestConfigOverridesAndNoDatabaseLookup(t *testing.T) {
	values := devEnv()
	values["HTTP_ADDR"], values["REQUEST_TIMEOUT"], values["LOG_LEVEL"] = "[::1]:0", "3s", "warn"
	values["PUSH_ENABLED"], values["CRON_ENABLED"] = "false", "false"
	lookup := func(key string) (string, bool) {
		if strings.Contains(key, "DATABASE") || strings.Contains(key, "SECRET") {
			t.Fatalf("M1-02 must not inspect database/secrets: %s", key)
		}
		value, ok := values[key]
		return value, ok
	}
	cfg, err := LoadConfig(lookup)
	if err != nil || cfg.HTTPAddr != "[::1]:0" || cfg.RequestTimeout != 3*time.Second {
		t.Fatalf("valid overrides rejected: %+v %v", cfg, err)
	}
	values["APP_ORIGIN"] = "http://[::1]:5173"
	if _, err := LoadConfig(lookupMap(values)); err != nil {
		t.Fatal(err)
	}
}
