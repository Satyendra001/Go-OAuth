package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// PgConfig set configuration
type DBConfig struct {
	Host     string
	Port     int64
	User     string
	Password string
	Database string
	Scheme   string
}

var DBEnvs = getDBConfig()

func getDBConfig() DBConfig {
	godotenv.Load()

	return DBConfig{
		// Make sure user is created, DB exists and user is assigned all the privileges
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Host:     getEnv("DB_HOST", "127.0.0.1"),
		Database: getEnv("DB_NAME", "infiviz_test"),
		Port:     getEnvAsInt("PORT", 5432),
		Scheme:   getEnv("SCHEMA", "public"),
	}
}

type DBUser struct {
	Id        int
	Email     string
	Username  string
	FirstName string
	LastName  string
}

func getEnv(key string, fallback string) string {

	val, ok := os.LookupEnv(key)

	if ok {
		return val
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	val, ok := os.LookupEnv(key)

	if ok {
		intData, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return fallback
		}
		return intData
	}

	return fallback
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func GenerateToken() (string, error) {
	// Generate a 32-byte random token
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	tokenStr := hex.EncodeToString(bytes)
	return tokenStr[0:40], nil
}
