package handler

import (
	"github.com/gohive/core/logger"
	"github.com/segmentio/kafka-go"
)

type MessageHandler interface {
	Handle(message kafka.Message) error
}

type OrderHandler struct{}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) Handle(message kafka.Message) error {
	logger.Infof("processing order message: %s", string(message.Value))
	// TODO: implement order processing logic
	return nil
}

type NotificationHandler struct{}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

func (h *NotificationHandler) Handle(message kafka.Message) error {
	logger.Infof("processing notification message: %s", string(message.Value))
	// TODO: implement notification processing logic
	return nil
}
