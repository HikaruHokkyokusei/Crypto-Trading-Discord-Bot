package utils

import "os"

func GetEnv(key string, defaultValue string) string {
	if value, isPresent := os.LookupEnv(key); isPresent {
		return value
	} else {
		return defaultValue
	}
}
