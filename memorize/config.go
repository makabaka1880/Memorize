package memorize

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type SM2Config struct {
	EFInitial      float64 `toml:"ef_initial"`
	EFMin          float64 `toml:"ef_min"`
	IntervalFirst  int     `toml:"interval_first"`
	IntervalSecond int     `toml:"interval_second"`
	EFCoefA        float64 `toml:"ef_coef_a"`
	EFCoefB        float64 `toml:"ef_coef_b"`
	EFCoefC        float64 `toml:"ef_coef_c"`
}

type Config struct {
	SM2 SM2Config `toml:"sm2"`
}

var cfg *Config

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if _, err := toml.Decode(string(data), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func GetConfig() *Config {
	return cfg
}

func SetConfig(c *Config) {
	cfg = c
}
