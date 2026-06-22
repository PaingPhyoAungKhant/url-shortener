package config_test

import (
	"testing"
	"time"

	"github.com/PaingPhyoAungKhant/url-shortener/internal/config"
)

func TestLoad_Defaults(t *testing.T) {
	cfg := config.Load()

	cases := []struct {
		name string
		got  any
		want any
	}{
		{"App.Env", cfg.App.Env, "development"},
		{"App.LogLevel", cfg.App.LogLevel, "info"},
		{"App.LogFormat", cfg.App.LogFormat, "json"},
		{"Server.Port", cfg.Server.Port, 8081},
		{"Server.ReadTimeout", cfg.Server.ReadTimeout, 15 * time.Second},
		{"Server.WriteTimeout", cfg.Server.WriteTimeout, 15 * time.Second},
		{"Server.IdleTimeout", cfg.Server.IdleTimeout, 60 * time.Second},
		{"Server.ShutdownTimeout", cfg.Server.ShutdownTimeout, 10 * time.Second},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Errorf("Load().%s = %v; want %v", tc.name, tc.got, tc.want)
			}
		})
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("SERVER_PORT", "9090")
	t.Setenv("SERVER_READ_TIMEOUT", "30s")

	cfg := config.Load()

	if cfg.App.Env != "production" {
		t.Errorf("App.Env = %q; want %q", cfg.App.Env, "production")
	}
	if cfg.Server.Port != 9090 {
		t.Errorf("Server.Port = %d; want 9090", cfg.Server.Port)
	}
	if cfg.Server.ReadTimeout != 30*time.Second {
		t.Errorf("Server.ReadTimeout = %v; want 30s", cfg.Server.ReadTimeout)
	}
}

func TestValidate(t *testing.T) {
	validBase := func() *config.Config {
		return &config.Config{
			App: config.AppConfig{
				Env:       "development",
				LogLevel:  "info",
				LogFormat: "json",
			},
			Server: config.ServerConfig{
				Port:            8080,
				ReadTimeout:     15 * time.Second,
				WriteTimeout:    15 * time.Second,
				IdleTimeout:     60 * time.Second,
				ShutdownTimeout: 10 * time.Second,
			},
		}
	}

	cases := []struct {
		name    string
		mutate  func(*config.Config)
		wantErr bool
	}{
		{
			name:    "valid config passes",
			mutate:  func(_ *config.Config) {},
			wantErr: false,
		},
		// App.Env
		{
			name:    "App.Env production is valid",
			mutate:  func(c *config.Config) { c.App.Env = "production" },
			wantErr: false,
		},
		{
			name:    "App.Env test is valid",
			mutate:  func(c *config.Config) { c.App.Env = "test" },
			wantErr: false,
		},
		{
			name:    "App.Env unknown value is rejected",
			mutate:  func(c *config.Config) { c.App.Env = "staging" },
			wantErr: true,
		},
		// Server.Port
		{
			name:    "Server.Port 1 is valid",
			mutate:  func(c *config.Config) { c.Server.Port = 1 },
			wantErr: false,
		},
		{
			name:    "Server.Port 65535 is valid",
			mutate:  func(c *config.Config) { c.Server.Port = 65535 },
			wantErr: false,
		},
		{
			name:    "Server.Port 0 is rejected",
			mutate:  func(c *config.Config) { c.Server.Port = 0 },
			wantErr: true,
		},
		{
			name:    "Server.Port 65536 is rejected",
			mutate:  func(c *config.Config) { c.Server.Port = 65536 },
			wantErr: true,
		},
		// App.LogLevel
		{
			name:    "App.LogLevel debug is valid",
			mutate:  func(c *config.Config) { c.App.LogLevel = "debug" },
			wantErr: false,
		},
		{
			name:    "App.LogLevel unknown value is rejected",
			mutate:  func(c *config.Config) { c.App.LogLevel = "trace" },
			wantErr: true,
		},
		// App.LogFormat
		{
			name:    "App.LogFormat text is valid",
			mutate:  func(c *config.Config) { c.App.LogFormat = "text" },
			wantErr: false,
		},
		{
			name:    "App.LogFormat unknown value is rejected",
			mutate:  func(c *config.Config) { c.App.LogFormat = "xml" },
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := validBase()
			tc.mutate(cfg)
			err := cfg.Validate()
			if tc.wantErr && err == nil {
				t.Error("Validate() = nil; want an error")
			}
			if !tc.wantErr && err != nil {
				t.Errorf("Validate() = %v; want nil", err)
			}
		})
	}
}
