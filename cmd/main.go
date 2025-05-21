package main

// @title AuthZ API
// @version 1.0
// @description Role-Permission və JWT yoxlama mikroservisi
// @host localhost:8000
// @BasePath /

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
	_ "ms-authz/docs"
	"ms-authz/internal/domain/model"
	"ms-authz/internal/handler"
	"ms-authz/internal/infrastructure/cache"
	"ms-authz/internal/infrastructure/db"
	"ms-authz/internal/infrastructure/mq"
	//"ms-authz/internal/consumer"
	"ms-authz/internal/service"
	"ms-authz/pkg/jwtutil"
	"os"
)

func main() {
	// Load configs (or use .env in real project)
	dbDSN := os.Getenv("DB_DSN")           // e.g. "host=localhost user=postgres dbname=authz sslmode=disable password=secret"
	rabbitURL := os.Getenv("RABBITMQ_URL") // e.g. "amqp://guest:guest@localhost:5672/"

	dbConn, err := gorm.Open(postgres.Open(dbDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ failed to connect to database:", err)
	}

	// ✅ Auto migrate models
	if err := dbConn.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.RolePermission{},
	); err != nil {
		log.Fatal("❌ migration failed:", err)
	}

	uow := db.NewUnitOfWork(dbConn)
	tokenRepo := cache.NewTokenRepository()
	cache.StartTokenCleanupService(tokenRepo, time.Minute*5)
	keyProvider := jwtutil.NewFileKeyProvider(os.Getenv("PUBLIC_KEY_DIR"))

	mqConn, err := mq.NewMQ(rabbitURL)
	if err != nil {
		log.Fatal("❌ failed to connect to RabbitMQ:", err)
	}
	defer mqConn.Close()

	publisher := mq.NewPublisherService(mqConn.Channel)

	authService := service.NewAuthService(tokenRepo, keyProvider)
	rbacService := service.NewRBACService(uow, publisher)

	// Start RabbitMQ consumers (fanout listeners)
	//consumer.StartConsumers(mqConn.Channel, authService, rbacService)

	app := fiber.New()

	authorizeHandler := handler.NewAuthorizeHandler(authService, rbacService,publisher)
	authorizeHandler.RegisterRoutes(app)

	rbacAdminHandler := handler.NewRBACAdminHandler(uow, rbacService)
	rbacAdminHandler.RegisterRoutes(app)

	app.Use(logger.New())
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	log.Println("✅ ms-authz service started on port", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("❌ fiber failed:", err)
	}
}
