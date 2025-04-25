package common

import (
	"fmt"
	"os"
	"strconv"

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

func GetEnvInt(key string) int {
	val := GetEnv(key)

	intVal, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("Environment variable %s is not an integer: %v", key, err))
	}

	return intVal
}

func MockEnv(key, v string) {
	if mock == nil {
		mock = make(map[string]string)
	}
	mock[key] = v
}
