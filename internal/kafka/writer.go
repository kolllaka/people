package kafka

import (
	"context"
	"strings"

	"github.com/KoLLlaka/sobes/internal/logger"
	"github.com/KoLLlaka/sobes/internal/model"
	kafkago "github.com/segmentio/kafka-go"
)

type Writer struct {
	writer *kafkago.Writer
	logger logger.MyLogger
}

// create a new writer to Kafka.
func NewKafkaWriter(conf model.Kafka, logger logger.MyLogger) *Writer {
	addr := strings.Join(conf.Brokers, ", ")
	writer := &kafkago.Writer{
		Addr:  kafkago.TCP(addr),
		Topic: conf.Topic,
	}

	return &Writer{
		writer: writer,
		logger: logger,
	}
}

// function to send message to Kafka
// message recieved from messagesChan channel
func (k *Writer) WriteMessages(ctx context.Context, messagesChan chan kafkago.Message, messageCommitChan chan kafkago.Message) error {
	for {
		select {
		case <-ctx.Done():
			k.logger.LogWarningLvl(
				"ctx.Done()",
				"kafka",
				"WriteMessages",
				ctx.Err(),
				nil,
			)

			return ctx.Err()
		case msg := <-messagesChan:
			err := k.writer.WriteMessages(ctx, kafkago.Message{
				Value: msg.Value,
			})
			if err != nil {
				k.logger.LogErrorLvl(
					"failed to write message",
					"kafka",
					"WriteMessages",
					err,
					msg.Value,
				)

				return err
			}

			select {
			case <-ctx.Done():
				k.logger.LogWarningLvl(
					"ctx.Done()",
					"kafka",
					"WriteMessages",
					ctx.Err(),
					nil,
				)
			case messageCommitChan <- msg:
				k.logger.LogTraceLvl(
					"message from WriteMessages send to messageCommitChan",
					"kafka",
					"WriteMessages",
					nil,
					string(msg.Value),
				)
			}
		}
	}
}
