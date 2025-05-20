package mq

import (
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	PublishEvent(exchange string, payload any, passiveQueues []string) error
}

type PublisherService struct {
	ch *amqp091.Channel
}

func NewPublisherService(ch *amqp091.Channel) *PublisherService {
	return &PublisherService{ch: ch}
}

func (p *PublisherService) PublishEvent(exchange string, payload any, passiveQueues []string) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Queue-ları əvvəlcə yaradıb bind edirik
	for _, qName := range passiveQueues {
		_, err := p.ch.QueueDeclare(
			qName, true, false, false, false, nil,
		)
		if err != nil {
			log.Printf("❌ Failed to declare queue %s: %v", qName, err)
			continue
		}

		err = p.ch.QueueBind(qName, "", exchange, false, nil)
		if err != nil {
			log.Printf("❌ Failed to bind queue %s: %v", qName, err)
		}
	}

	err = p.ch.Publish(
		exchange,
		"",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("❌ Failed to publish to %s: %v", exchange, err)
		return err
	}

	log.Printf("📤 Published event to %s: %s", exchange, string(body))
	return nil
}
