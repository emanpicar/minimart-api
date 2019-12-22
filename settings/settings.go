package settings

import (
	"os"
)

func getEnv(envName, envDefault string) string {
	if envValue := os.Getenv(envName); envValue != "" {
		return envValue
	}

	return envDefault
}

func GetLogLevel() string {
	return getEnv("LOG_LEVEL", "info")
}

func GetDBHost() string {
	return getEnv("DB_HOST", "localhost")
}

func GetDBPort() string {
	return getEnv("DB_PORT", "5432")
}

func GetDBUser() string {
	return getEnv("DB_USER", "secretdbuser")
}

func GetDBPass() string {
	return getEnv("DB_PASS", "secretdbpass")
}

func GetServerHost() string {
	return getEnv("SERVER_HOST", "")
}

func GetServerPort() string {
	return getEnv("SERVER_PORT", "9988")
}

func GetServerPublicKey() string {
	return getEnv("SERVER_PUBLIC_KEY", "./certs/cert.pem")
}

func GetServerPrivateKey() string {
	return getEnv("SERVER_PRIVATE_KEY", "./certs/key.pem")
}

func GetTokenSecret() string {
	return getEnv("TOKEN_SECRET", "notSoSecret")
}
