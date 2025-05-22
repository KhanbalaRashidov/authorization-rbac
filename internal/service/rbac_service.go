package service

import (
	"log"
	"ms-authz/internal/domain/model"
	"ms-authz/internal/domain/repository"
	"ms-authz/internal/infrastructure/mq"
	"strings"
	"sync"
)

type RBACService struct {
	uow       repository.UnitOfWork
	cache     sync.Map // map[roleName][]permissionCode
	publisher mq.Publisher
}

func NewRBACService(uow repository.UnitOfWork, publisher mq.Publisher) *RBACService {
	s := &RBACService{
		uow:       uow,
		publisher: publisher,
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
		if strings.EqualFold(p, permission) {
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
//func (s *RBACService) PublishCacheReload() {
//	s.PublishCacheEvent("RBAC_CACHE_RELOAD", map[string]any{})
//}

func (s *RBACService) PublishCacheEvent(event string, payload map[string]any) {
	message := map[string]any{
		"event": event,
	}
	for k, v := range payload {
		message[k] = v
	}

	_ = s.publisher.PublishEvent("rbac.update.fanout", message, []string{})
}
