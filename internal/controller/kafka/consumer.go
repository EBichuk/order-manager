package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"order-manager/internal/models"
	"order-manager/pkg/errorx"
	"strings"

	"github.com/segmentio/kafka-go"
)

type service interface {
	GetOrderByUID(string) (*models.Order, error)
	SaveOrder(*models.Order) error
}

type Consumer struct {
	reader *kafka.Reader
	s      service
	log    *slog.Logger
}

// конфиг
func NewConsumer(s service, log *slog.Logger, topic string, brokers string) *Consumer {
	brokersList := strings.Split(brokers, ",")
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokersList,
		Topic:   topic,
		GroupID: "0",
	})

	return &Consumer{
		reader: r,
		s:      s,
		log:    log,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	c.log.Info("Starting kafka consumer")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			m, err := c.reader.FetchMessage(ctx)
			if err != nil {
				break
			}

			c.log.Info("Got message from kafka", slog.Int("Partition", m.Partition), slog.Int("Offset", int(m.Offset)))

			new_order := models.Order{}
			err = json.Unmarshal(m.Value, &new_order)
			if err != nil {
				c.log.Error("Failed to unmarshal message")
				continue
			}

			err = c.s.SaveOrder(&new_order)
			if err != nil {
				c.log.Warn("Not saved order", slog.String("Error", err.Error()))
				if errors.Is(err, errorx.ErrOrderValidation) {
					err = c.reader.CommitMessages(ctx, m)
					if err != nil {
						c.log.Error("Failed to commit message", slog.String("Error", err.Error()))
					}
				}
				continue
			}

			err = c.reader.CommitMessages(ctx, m)
			if err != nil {
				c.log.Error("Failed to commit message", slog.String("Error", err.Error()))
			} else {
				c.log.Info("Commited messsage")
			}
		}
	}
}

func (c *Consumer) Stop() error {
	c.log.Info("Consumer is stopping")
	if err := c.reader.Close(); err != nil {
		c.log.Error("Failed to stop reader", slog.String("Error", err.Error()))
		return err
	}
	return nil
}
