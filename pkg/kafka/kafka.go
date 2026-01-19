package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	Brokers []string
	GroupID string
}

type Reader struct {
	*kafka.Reader
}

func NewReader(cfg Config, topic string) *Reader {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		GroupID:        cfg.GroupID,
		Topic:          topic,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	})
	return &Reader{Reader: r}
}

type Writer struct {
	*kafka.Writer
}

func NewWriter(brokers []string, topic string) *Writer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		MaxAttempts:  5,
	}
	return &Writer{Writer: w}
}

func (w *Writer) SendMessage(ctx context.Context, key, value []byte) error {
	return w.WriteMessages(ctx, kafka.Message{
		Key:   key,
		Value: value,
	})
}
