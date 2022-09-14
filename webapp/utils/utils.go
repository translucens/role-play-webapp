package utils

import (
	"os"
)

func GetEnvProjectID() string {
	return getEnv("GOOGLE_CLOUD_PROJECT", "")
}

func GetEnvDBInstanceID() string {
	return getEnv("DB_INSTANCE_ID", "scstore")
}

func GetEnvDBName() string {
	return getEnv("DB_NAME", "scstore")
}

func getEnv(key, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}
