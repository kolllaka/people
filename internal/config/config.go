package config

import (
	"log"
	"os"
	"strings"

	"github.com/KoLLlaka/sobes/internal/model"
	"github.com/joho/godotenv"
)

const (
	HOST             = "HOST"
	PORT             = "PORT"
	FROMKAFKABROKERS = "FROMKAFKABROKERS"
	FROMKAFKATOPIC   = "FROMKAFKATOPIC"
	FROMKAFKAGROUPID = "FROMKAFKAGROUPID"
	TOKAFKABROKERS   = "TOKAFKABROKERS"
	TOKAFKATOPIC     = "TOKAFKATOPIC"
	TOKAFKAGROUPID   = "TOKAFKAGROUPID"
	USER             = "USER"
	PASSWORD         = "PASSWORD"
	DBNAME           = "DBNAME"
	SSLMODE          = "SSLMODE"
)

// loads values from .env into the system
func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// loads values from .env to Config
func NewConfig() model.Config {
	fromKafkaBrokers := strings.Split(getEnv(FROMKAFKABROKERS, ""), ",")
	toKafkaBrokers := strings.Split(getEnv(TOKAFKABROKERS, ""), ",")

	return model.Config{
		Host: getEnv(HOST, "localhost"),
		Port: getEnv(PORT, "8080"),
		FromKafka: model.Kafka{
			Brokers: fromKafkaBrokers,
			Topic:   getEnv(FROMKAFKATOPIC, ""),
			GroupID: getEnv(FROMKAFKAGROUPID, ""),
		},
		ToKafka: model.Kafka{
			Brokers: toKafkaBrokers,
			Topic:   getEnv(TOKAFKATOPIC, ""),
			GroupID: getEnv(TOKAFKAGROUPID, ""),
		},
		DBconf: model.DBconf{
			User:     getEnv(USER, ""),
			Password: getEnv(PASSWORD, ""),
			Dbname:   getEnv(DBNAME, ""),
			Sslmode:  getEnv(SSLMODE, ""),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
