package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI       string
	DatabaseName   string
	ServerPort     string
	TestMode       bool
	JWTSecret      string
	JWTExpiration  time.Duration
	PasswordPepper string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	testMode, _ := strconv.ParseBool(os.Getenv("TEST_MODE"))

	defaultExpiration := int64(3600)

	jwtExpiration, _ := strconv.ParseInt(os.Getenv("JWT_EXPIRATION"), 10, 64)

	if jwtExpiration == 0 {
		jwtExpiration = defaultExpiration
	}

	return &Config{
		MongoURI:       getEnv("MONGO_URI", "mongodb://localhost:27017"),
		DatabaseName:   getEnv("DATABASE_NAME", "todo_db"),
		ServerPort:     getEnv("SERVER_PORT", ":8080"),
		TestMode:       testMode,
		JWTSecret:      getEnv("JWT_SECRET", "very-secret-key"),
		JWTExpiration:  time.Duration(jwtExpiration) * time.Second,
		PasswordPepper: getEnv("PASSWORD_PEPPER", "pepper"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
