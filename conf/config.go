package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type Config struct {
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	Pepper   string         `json:"pepper"`
	HMACKey  string         `json:"hmac_key"`
	Database PostgresConfig `json:"database"`
	Mailgun  MailgunConfig  `json:"mailgun"`
}

type MailgunConfig struct {
	APIKey       string `json:"api_key"`
	PublicAPIKey string `json:"public_api_key"`
	Domain       string `json:"domain"`
}

func (c PostgresConfig) PsqlConnInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Name)
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)
}

func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "adam",
		Password: "",
		Name:     "picapp",
	}
}

func DefaultConfig() Config {
	return Config{
		Port:     3000,
		Env:      "dev",
		Pepper:   "+&_|U;_?=r]}~7NZTVf>|^eG>QwL{!^eYkX=TN.4C\\\".3D$fXo`",
		HMACKey:  "secret-hmac-key",
		Database: DefaultPostgresConfig(),
	}
}

func LoadConfig(confRequired bool) Config {
	f, err := os.Open(".config")
	if err != nil {
		if confRequired {
			panic(err)
		}
		fmt.Println("Using the default config...")
		return DefaultConfig()
	}
	var c Config
	dec := json.NewDecoder(f)
	err = dec.Decode(&c)
	if err != nil {
		panic(err)
	}
	fmt.Println(".config loaded successfully...")
	return c
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}
