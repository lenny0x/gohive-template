package consumer

import (
	"context"

	"github.com/gohive/core/logger"
	"github.com/gohive/demo-worker-kafka/handler"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	readers  map[string]*kafka.Reader
	handlers map[string]handler.MessageHandler
}

func NewConsumer() *Consumer {
	return &Consumer{
		readers:  make(map[string]*kafka.Reader),
		handlers: make(map[string]handler.MessageHandler),
	}
}

func (c *Consumer) RegisterHandler(topic string, reader *kafka.Reader, h handler.MessageHandler) {
	c.readers[topic] = reader
	c.handlers[topic] = h
}

func (c *Consumer) Start(ctx context.Context) error {
	for topic, reader := range c.readers {
		go c.consumeTopic(ctx, topic, reader)
	}

	<-ctx.Done()
	return ctx.Err()
}

func (c *Consumer) consumeTopic(ctx context.Context, topic string, reader *kafka.Reader) {
	h, ok := c.handlers[topic]
	if !ok {
		logger.Warnf("no handler for topic: %s", topic)
		return
	}

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			logger.Errorf("failed to read message: %v", err)
			continue
		}

		if err := h.Handle(msg); err != nil {
			logger.Errorf("failed to handle message: %v", err)
			continue
		}
	}
}

func (c *Consumer) Close() error {
	for _, reader := range c.readers {
		if err := reader.Close(); err != nil {
			logger.Errorf("failed to close reader: %v", err)
		}
	}
	return nil
}
