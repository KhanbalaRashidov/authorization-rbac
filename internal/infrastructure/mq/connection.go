package mq

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type MQ struct {
	Connection *amqp091.Connection
	Channel    *amqp091.Channel
}

// NewMQ RabbitMQ bağlantısını qurur və exchange-ləri yaradır
func NewMQ(rabbitURL string) (*MQ, error) {
	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Exchange-ləri fanout növündə yaradırıq
	exchanges := []string{"auth.tokens.fanout", "rbac.update.fanout"}
	for _, ex := range exchanges {
		err := ch.ExchangeDeclare(
			ex,       // name
			"fanout", // type
			true,     // durable
			false,    // auto-deleted
			false,    // internal
			false,    // no-wait
			nil,      // arguments
		)
		if err != nil {
			_ = ch.Close()
			_ = conn.Close()
			return nil, fmt.Errorf("failed to declare exchange '%s': %w", ex, err)
		}
	}

	log.Println("[MQ] Connected and exchanges declared")
	return &MQ{
		Connection: conn,
		Channel:    ch,
	}, nil
}

// Close RabbitMQ connection və channel-i bağlayır
func (mq *MQ) Close() {
	if mq.Channel != nil {
		if err := mq.Channel.Close(); err != nil {
			log.Println("Error closing MQ channel:", err)
		}
	}
	if mq.Connection != nil {
		if err := mq.Connection.Close(); err != nil {
			log.Println("Error closing MQ connection:", err)
		}
	}
}
