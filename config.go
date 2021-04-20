package main

import "fmt"

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
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

type Config struct {
	Port int
	Env  string
}

func (c Config) IsProd() bool {
	return c.Env == "prod"
}

func DefaultConfig() Config {
	return Config{
		Port: 3000,
		Env:  "dev",
	}
}

//# models/users.go
//const userPwPepper = "+&_|U;_?=r]}~7NZTVf>|^eG>QwL{!^eYkX=TN.4C\".3D$fXo`"
//const hmacSecretKey = "secret-hmac-key"
//
//# models/services.go
//db, err := gorm.Open(postgres.Open(connectionInfo),
//&gorm.Config{
//Logger: logger.Default.LogMode(logger.Info),
//})
