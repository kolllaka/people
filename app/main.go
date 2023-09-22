package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/KoLLlaka/sobes/internal/config"
	"github.com/KoLLlaka/sobes/internal/db"
	"github.com/KoLLlaka/sobes/internal/handler"
	"github.com/KoLLlaka/sobes/internal/kafka"
	"github.com/KoLLlaka/sobes/internal/logger"
	"github.com/KoLLlaka/sobes/internal/model"
	"github.com/KoLLlaka/sobes/internal/services"

	externalapi "github.com/KoLLlaka/sobes/internal/externalAPI"
	"github.com/sirupsen/logrus"
)

func main() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	logrusLogger := logrus.New()
	logger := logger.NewMyLooger(logrusLogger)
	conf := config.NewConfig()
	logger.LogTraceLvl("config", "main", "main", nil, conf)

	confFromKafka := conf.FromKafka
	reader := kafka.NewKafkaReader(confFromKafka, logger)

	confToKafka := conf.ToKafka
	writer := kafka.NewKafkaWriter(confToKafka, logger)

	ctx := context.Background()
	messages := kafka.MakeChan(1e3)
	errormessages := kafka.MakeChan(1e3)
	messageCommitChan := kafka.MakeChan(1e3)

	// ? DB
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", conf.User, conf.Password, conf.Dbname, conf.Sslmode)
	postgresDB, err := db.NewDB(dsn)
	if err != nil {
		logger.LogFatalLvl(err)
	}
	store := db.NewStore(postgresDB, logger)

	// ? enrich
	enrich := externalapi.NewEnrich(logger)
	// ? service
	peopleService := services.NewPeopleService(&store, &enrich, &logger)

	// ? server
	server := handler.NewServer(logger, peopleService)
	server.StartServer(conf.Host, conf.Port)

	go func() {
		reader.FetchMessage(ctx, messages)
	}()

	go func() {
		writer.WriteMessages(ctx, errormessages, messageCommitChan)
	}()

	go func() {
		reader.CommitMessages(ctx, messageCommitChan)
	}()

	go func() {
		for {
			messageFromKafka := model.MessageFromKafka{}
			msg := <-messages

			err := json.Unmarshal(msg.Value, &messageFromKafka)
			if err != nil {
				logger.LogTraceLvl("wrong data from Kafka", "main", "main", err, string(msg.Value))
				messageToKafka := model.MessageToKafka{
					MessageFromKafka: msg.Value,
					Error:            err,
				}
				msg.Value, _ = json.Marshal(messageToKafka)
				errormessages <- msg

				continue
			}

			messageCommitChan <- msg

			messageToDB := model.MessageToDB{
				Name:       messageFromKafka.Name,
				Surname:    messageFromKafka.Surname,
				Patronymic: messageFromKafka.Patronymic,
			}

			peopleService.EnrichPeople(&messageToDB)
			peopleService.AddPeople(messageToDB)
		}
	}()

	<-sc
	logger.LogTraceLvl("Stopping...", "main", "main", nil, nil)
	logger.LogTraceLvl("postgres DB closing...", "main", "main", nil, nil)
	postgresDB.Close()
	logger.LogTraceLvl("postgres DB close", "main", "main", nil, nil)
}
