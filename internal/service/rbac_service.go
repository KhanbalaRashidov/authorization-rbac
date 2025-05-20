package service

import (
	"encoding/json"
	"log"
	"ms-authz/internal/domain/model"
	"ms-authz/internal/domain/repository"
	"sync"

	"github.com/rabbitmq/amqp091-go"
)

type RBACService struct {
	uow             repository.UnitOfWork
	cache           sync.Map // map[roleName][]permissionCode
	rabbitMQChannel *amqp091.Channel
}

func NewRBACService(uow repository.UnitOfWork, rabbitCh *amqp091.Channel) *RBACService {
	s := &RBACService{
		uow:             uow,
		rabbitMQChannel: rabbitCh,
	}
	s.LoadCache()
	return s
}

// Sistemdəki bütün rolları və permission-ları yaddaşa yükləyir
func (s *RBACService) LoadCache() {
	roles, _ := s.uow.RoleRepo().GetAll()
	for _, role := range roles {
		s.loadRole(role)
	}
}

// RBAC cache-də permission yoxlama
func (s *RBACService) HasPermission(roleName string, permission string) bool {
	val, ok := s.cache.Load(roleName)
	if !ok {
		return false
	}
	perms := val.([]string)
	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}

// Cache-də konkret bir rolu yüklə
func (s *RBACService) loadRole(role model.Role) {
	perms, _ := s.uow.RolePermissionRepo().GetPermissionsByRoleID(role.ID)
	var names []string
	for _, p := range perms {
		names = append(names, p.Name)
	}
	s.cache.Store(role.Name, names)
}

// CRUD sonrası və ya MQ ilə çağırıla bilər
func (s *RBACService) ReloadCache() {
	log.Println("Reloading RBAC cache...")
	s.LoadCache()
}

// MQ ilə digər instansiyalara xəbər göndərir
func (s *RBACService) PublishCacheReload() {
	event := struct {
		Event string `json:"event"`
	}{
		Event: "RBAC_CACHE_RELOAD",
	}

	body, _ := json.Marshal(event)
	err := s.rabbitMQChannel.Publish(
		"rbac.update.fanout", // Exchange
		"",                   // Routing key
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Println("❌ Failed to publish RBAC_CACHE_RELOAD event:", err)
	} else {
		log.Println("📤 RBAC_CACHE_RELOAD event published")
	}
}
