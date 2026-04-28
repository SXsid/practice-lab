package app

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DSN      string
	RedisURL string
	LogLevel string
}

func getEnv(key string, required bool) string {
	res := os.Getenv(key)
	if res == "" && required {
		panic(fmt.Sprintf("requred key not found key: %s", key))
	}
	return res
}

func getEnvInt(key string, required bool) int {
	value := getEnv(key, required)
	res, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return res
}

func NewConfig() *Config {
	return &Config{
		DSN:      getEnv("DSN", true),
		RedisURL: getEnv("RedisURL", true),
		LogLevel: getEnv("LogLevel", false),
	}
}
