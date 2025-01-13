package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	DefaultProvider string            `mapstructure:"default_provider"`
	DefaultRegion   string            `mapstructure:"default_region"`
	Providers       map[string]Provider `mapstructure:"providers"`
}

type Provider struct {
	Region      string            `mapstructure:"region"`
	Credentials map[string]string `mapstructure:"credentials"`
}

func LoadConfig(cfgFile string) (*Config, error) {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".cloud-cli")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create default config
			return createDefaultConfig()
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

func createDefaultConfig() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	config := &Config{
		DefaultProvider: "aws",
		DefaultRegion:   "us-west-2",
		Providers: map[string]Provider{
			"aws": {
				Region: "us-west-2",
				Credentials: map[string]string{
					"access_key_id":     "",
					"secret_access_key": "",
				},
			},
			"gcp": {
				Region: "us-central1",
				Credentials: map[string]string{
					"credentials_file": "",
				},
			},
			"azure": {
				Region: "eastus",
				Credentials: map[string]string{
					"subscription_id": "",
					"tenant_id":      "",
					"client_id":      "",
					"client_secret": "",
				},
			},
		},
	}

	configFile := filepath.Join(home, ".cloud-cli.yaml")
	if err := saveConfig(configFile, config); err != nil {
		return nil, err
	}

	return config, nil
}

func saveConfig(file string, config *Config) error {
	viper.SetConfigFile(file)
	
	for provider, cfg := range config.Providers {
		viper.Set(fmt.Sprintf("providers.%s.region", provider), cfg.Region)
		for k, v := range cfg.Credentials {
			viper.Set(fmt.Sprintf("providers.%s.credentials.%s", provider, k), v)
		}
	}
	
	viper.Set("default_provider", config.DefaultProvider)
	viper.Set("default_region", config.DefaultRegion)

	return viper.WriteConfig()
}
