package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	DBHost        string
	DBPort        int
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	redisPort, _ := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))

	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        dbPort,
		DBUser:        getEnv("DB_USER", "library"),
		DBPassword:    getEnv("DB_PASSWORD", "library123"),
		DBName:        getEnv("DB_NAME", "library_db"),
		JWTSecret:     getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     redisPort,
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
	}, nil
}

// func (c *Config) DBConnString() string {
// 	return fmt.Sprintf(
// 		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
// 	)
// }

func (c *Config) DBConnString() string {
	// Если в системе есть DATABASE_URL (как в Railway), используем её целиком
	if envURL := os.Getenv("DATABASE_URL"); envURL != "" {
		return envURL
	}

	// Иначе (на локалке) собираем по кусочкам
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName,
	)
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
