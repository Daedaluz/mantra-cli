package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type API struct {
	Server     string `json:"server"`
	APIKey     string `json:"api_key,omitempty"`
	Plaintext  bool   `json:"plaintext,omitempty"`
	SkipVerify bool   `json:"skip_verify,omitempty"`
}

type Context struct {
	API          string `json:"api"`
	Domain       string `json:"domain"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RegisterPath string `json:"register_path"`
	AuthPath     string `json:"auth_path"`
}

type Config struct {
	CurrentContext string              `json:"current_context"`
	APIs           map[string]*API     `json:"apis"`
	Contexts       map[string]*Context `json:"contexts"`
}

func ConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".mantra", "config.json")
}

func Load() *Config {
	path := ConfigPath()
	if path == "" {
		return &Config{APIs: map[string]*API{}, Contexts: map[string]*Context{}}
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return &Config{APIs: map[string]*API{}, Contexts: map[string]*Context{}}
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return &Config{APIs: map[string]*API{}, Contexts: map[string]*Context{}}
	}
	if cfg.APIs == nil {
		cfg.APIs = map[string]*API{}
	}
	if cfg.Contexts == nil {
		cfg.Contexts = map[string]*Context{}
	}
	return &cfg
}

func Save(cfg *Config) error {
	path := ConfigPath()
	if path == "" {
		return os.ErrNotExist
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func (c *Config) ResolveContext() (*Context, *API) {
	if c.CurrentContext == "" {
		return nil, nil
	}
	ctx, ok := c.Contexts[c.CurrentContext]
	if !ok {
		return nil, nil
	}
	api, ok := c.APIs[ctx.API]
	if !ok {
		return ctx, nil
	}
	return ctx, api
}
