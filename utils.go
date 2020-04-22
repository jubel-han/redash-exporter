package main

import (
	"strconv"
	"fmt"
	"os"
)

func getEnv(key string, defaultVal string) string {
	if envVal, ok := os.LookupEnv(key); ok {
		return envVal
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if envVal, ok := os.LookupEnv(key); ok {
		if v, err := strconv.Atoi(envVal); err == nil {
			return v
		}
	}
	return defaultVal
}

func logPanic(err error) {
	if err != nil {
		fmt.Print(err)
		panic(err)
	}
}
