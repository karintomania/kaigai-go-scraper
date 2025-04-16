package common

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var mock map[string]string

func GetEnv(key string) string {
	// use mock value if exists
	if val, ok := mock[key]; ok {
		return val
	}

	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Environment variable %s not set", key))
	}
	return val
}

func MockEnv(key, v string) {
	if mock == nil {
		mock = make(map[string]string)
	}
	mock[key] = v
}
