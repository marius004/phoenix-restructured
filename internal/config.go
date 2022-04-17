package internal

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbName     string
	DbPassword string

	CookieLifetime int
	JwtSecret      string

	ServerHost string
	ServerPort string
}

func NewConfig() *Config {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalln("Error loading .env file", err)
	}

	var (
		dbHost     = os.Getenv("DB_HOST")
		dbPort     = os.Getenv("DB_PORT")
		dbUser     = os.Getenv("DB_USER")
		dbName     = os.Getenv("DB_NAME")
		dbPassword = os.Getenv("DB_PASSWORD")

		jwtSecret = os.Getenv("JWT_SECRET")

		serverHost = os.Getenv("SERVER_HOST")
		serverPort = os.Getenv("SERVER_PORT")
	)

	cookieLifetime, err := strconv.Atoi(os.Getenv("COOKIE_LIFETIME"))
	if err != nil {
		panic(err)
	}

	return &Config{
		DbHost:     dbHost,
		DbPort:     dbPort,
		DbUser:     dbUser,
		DbName:     dbName,
		DbPassword: dbPassword,

		CookieLifetime: cookieLifetime,
		JwtSecret:      jwtSecret,

		ServerHost: serverHost,
		ServerPort: serverPort,
	}
}
