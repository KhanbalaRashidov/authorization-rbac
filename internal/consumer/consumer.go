package consumer

import (
"encoding/json"
"log"
"ms-authz/internal/service"
"github.com/rabbitmq/amqp091-go"
)

func StartConsumers(ch *amqp091.Channel, authSvc *service.AuthService, rbacSvc *service.RBACService) {
	declareExchangeAndQueue(ch, "auth.tokens.fanout", "auth.tokens.queue", func(d amqp091.Delivery) {
		var e struct {
			Event string `json:"event"`
			JTI   string `json:"jti"`
			Exp   int64  `json:"exp"`
		}
		if err := json.Unmarshal(d.Body, &e); err == nil && e.Event == "TOKEN_BLACKLISTED" {
			log.Println("✅ Received TOKEN_BLACKLISTED")
			authSvc.HandleBlacklistEvent(e.JTI, e.Exp)
		}
	})

	declareExchangeAndQueue(ch, "rbac.update.fanout", "rbac.update.queue", func(d amqp091.Delivery) {
		var e struct {
			Event string `json:"event"`
		}
		if err := json.Unmarshal(d.Body, &e); err == nil && e.Event == "RBAC_CACHE_RELOAD" {
			log.Println("✅ Received RBAC_CACHE_RELOAD")
			rbacSvc.ReloadCache()
		}
	})
}

func declareExchangeAndQueue(ch *amqp091.Channel, exchangeName, queueName string, handle func(amqp091.Delivery)) {
	// Declare durable fanout exchange
	must(ch.ExchangeDeclare(exchangeName, "fanout", true, false, false, false, nil))

	// Declare durable queue
	queue, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
	must(err)

	// Bind queue to exchange
	must(ch.QueueBind(queue.Name, "", exchangeName, false, nil))

	// Start consumer
	msgs, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	must(err)

	go func() {
		for d := range msgs {
			handle(d)
		}
	}()
}

func must(err error) {
	if err != nil {
		log.Fatalf("❌ %v", err)
	}
}
