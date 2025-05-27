package config

import (
	"fmt"
	"os"
	"strconv"
)

func Get(key string, def string) string {
	if key, ok := os.LookupEnv(key); ok {
		return key
	}
	return def
}

func GetInt(key string, def int) int {
	value, err := strconv.Atoi(Get(key, fmt.Sprintf("%d", def)))
	if err != nil {
		return def
	}
	return value
}

func GetBool(key string, def bool) bool {
	value, err := strconv.ParseBool(Get(key, fmt.Sprintf("%t", def)))
	if err != nil {
		return def
	}
	return value
}
