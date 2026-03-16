package memorize

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type SM2Config struct {
	EFInitial       float64 `toml:"ef_initial"`
	EFMin           float64 `toml:"ef_min"`
	EFCoefA         float64 `toml:"ef_coef_a"`
	EFCoefB         float64 `toml:"ef_coef_b"`
	EFCoefC         float64 `toml:"ef_coef_c"`
	IntervalFirst   int     `toml:"interval_first"`
	IntervalSecond  int     `toml:"interval_second"`
	GraduateMinRep  int     `toml:"graduate_min_rep"`
	GraduateMaxIntv int     `toml:"graduate_max_intv"`
	GraduateMinEF   float64 `toml:"graduate_min_ef"`
}

type Config struct {
	SM2 SM2Config `toml:"sm2"`
}

var cfg *Config

func DefaultConfig() *Config {
	return &Config{
		SM2: SM2Config{
			EFInitial:       2.5,
			EFMin:           1.3,
			EFCoefA:         0.1,
			EFCoefB:         0.08,
			EFCoefC:         0.02,
			IntervalFirst:   1,
			IntervalSecond:  6,
			GraduateMinRep:  5,
			GraduateMaxIntv: 180,
			GraduateMinEF:   2.8,
		},
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	config := DefaultConfig() // start from defaults so partial TOML files still work
	if _, err := toml.Decode(string(data), config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	return config, nil
}

func GetConfig() *Config {
	return cfg
}

func SetConfig(c *Config) {
	cfg = c
}
