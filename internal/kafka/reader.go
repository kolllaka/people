package kafka

import (
	"context"

	"github.com/KoLLlaka/sobes/internal/logger"
	"github.com/KoLLlaka/sobes/internal/model"
	"github.com/pkg/errors"
	kafkago "github.com/segmentio/kafka-go"
)

type Reader struct {
	Reader *kafkago.Reader
	logger logger.MyLogger
}

// create a new channel of kafkago.Message with size size
func MakeChan(size int) chan kafkago.Message {
	return make(chan kafkago.Message, size)
}

// create a new reader from Kafka.
func NewKafkaReader(conf model.Kafka, logger logger.MyLogger) *Reader {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: conf.Brokers,
		Topic:   conf.Topic,
		GroupID: conf.GroupID,
	})

	return &Reader{
		Reader: reader,
		logger: logger,
	}
}

// function to receives message from Kafka
// message set to messagesChan channel
func (k *Reader) FetchMessage(ctx context.Context, messagesChan chan<- kafkago.Message) error {
	for {
		msg, err := k.Reader.FetchMessage(ctx)
		if err != nil {
			k.logger.LogErrorLvl(
				"failed to fetch message",
				"kafka",
				"FetchMessage",
				err,
				nil,
			)

			return err
		}

		select {
		case <-ctx.Done():
			k.logger.LogWarningLvl(
				"ctx.Done()",
				"kafka",
				"FetchMessage",
				ctx.Err(),
				nil,
			)

			return ctx.Err()
		case messagesChan <- msg:
			k.logger.LogTraceLvl(
				"message fetched and sent to a channel",
				"kafka",
				"FetchMessage",
				nil,
				string(msg.Value),
			)
		}
	}
}

// function to commit Kafka message
// commit receives from messageCommitChan channel
func (k *Reader) CommitMessages(ctx context.Context, messageCommitChan chan kafkago.Message) error {
	for {
		select {
		case <-ctx.Done():
			k.logger.LogWarningLvl(
				"ctx.Done()",
				"kafka",
				"CommitMessages",
				ctx.Err(),
				nil,
			)

			continue
		case msg := <-messageCommitChan:
			err := k.Reader.CommitMessages(ctx, msg)
			if err != nil {
				k.logger.LogErrorLvl(
					"Reader.CommitMessages",
					"kafka",
					"CommitMessages",
					err,
					nil,
				)

				return errors.Wrap(err, "Reader.CommitMessages")
			}
			k.logger.LogTraceLvl(
				"commited an msg",
				"kafka",
				"CommitMessages",
				nil,
				string(msg.Value),
			)
		}
	}
}
