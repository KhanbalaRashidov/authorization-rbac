package mq

import (
	"encoding/json"
	"log"
	"ms-authz/internal/service"

	"github.com/rabbitmq/amqp091-go"
)

func StartConsumers(ch *amqp091.Channel, authSvc *service.AuthService, rbacSvc *service.RBACService) {
	// Token Blacklist Consumer
	tokenQ, _ := ch.QueueDeclare("", false, true, true, false, nil)
	_ = ch.QueueBind(tokenQ.Name, "", "auth.tokens.fanout", false, nil)
	msgs1, _ := ch.Consume(tokenQ.Name, "", true, false, false, false, nil)

	go func() {
		for d := range msgs1 {
			var event struct {
				Event string `json:"event"`
				JTI   string `json:"jti"`
				Exp   int64  `json:"exp"`
			}
			if err := json.Unmarshal(d.Body, &event); err == nil && event.Event == "TOKEN_BLACKLISTED" {
				log.Println("[MQ] Received TOKEN_BLACKLISTED event")
				authSvc.HandleBlacklistEvent(event.JTI, event.Exp)
			}
		}
	}()

	// RBAC Update Consumer
	rbacQ, _ := ch.QueueDeclare("", false, true, true, false, nil)
	_ = ch.QueueBind(rbacQ.Name, "", "rbac.update.fanout", false, nil)
	msgs2, _ := ch.Consume(rbacQ.Name, "", true, false, false, false, nil)

	go func() {
		for d := range msgs2 {
			var event struct {
				Event string `json:"event"`
			}
			if err := json.Unmarshal(d.Body, &event); err == nil && event.Event == "RBAC_CACHE_RELOAD" {
				log.Println("[MQ] Received RBAC_CACHE_RELOAD event")
				rbacSvc.ReloadCache()
			}
		}
	}()
}
