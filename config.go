package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

var (
	defaultPort = "5000"
)

type Config struct {
	SendGridAPIKey string `toml:"sendgrid_api_key"`
	Port           string `toml:"port"`
}

func ReadConfig(data string) (Config, error) {
	var cfg Config
	meta, err := toml.Decode(data, &cfg)
	if err != nil {
		return Config{}, err
	}

	// build set of present config keys
	present := make(map[string]bool)
	for _, keys := range meta.Keys() {
		key := keys[0]
		present[key] = true
	}

	required := []string{}

	// gather missing keys
	missing := []string{}
	for _, key := range required {
		if _, ok := present[key]; !ok {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		msg := strings.Join(missing, ", ")
		return Config{}, fmt.Errorf("missing config values: %s", msg)
	}

	// handle defaults
	if cfg.Port == "" {
		cfg.Port = defaultPort
	}

	return cfg, nil
}

func ReadConfigFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	return ReadConfig(string(data))
}
