package main

import (
	categoryRoutes "catalog-service/internal/category/routes"
	"catalog-service/internal/database"
	catalogMiddleware "catalog-service/internal/middleware"
	redisClient "catalog-service/internal/redis"

	"os"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"catalog-service/internal/logger"
	requestMiddleware "catalog-service/internal/middleware"

	fiberCors "github.com/gofiber/fiber/v2/middleware/cors"
	fiberRecover "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/joho/godotenv"
)

func main() {

	// Load ENV

	err := godotenv.Load()

	if err != nil {
		logger.Log.Fatal(".env file not loaded")
	}

	logger.InitLogger()

	defer logger.Log.Sync()

	logger.Log.Info("starting catalog service")

	// Connect PostgreSQL

	database.ConnectPostgres()

	logger.Log.Info("postgres connected")

	// Connect Redis

	redisClient.ConnectRedis()

	redisClient.InitFiberStorage()

	logger.Log.Info("redis connected")

	// Fiber App

	app := fiber.New(fiber.Config{
		AppName: "Ahmed Catalog Service",
	})

	// =========================
	// Global Middleware
	// =========================

	// Panic Recovery

	app.Use(fiberRecover.New())

	// Request ID

	app.Use(requestMiddleware.RequestContextMiddleware())

	app.Use(requestMiddleware.RequestLoggerMiddleware())

	// CORS

	app.Use(fiberCors.New(fiberCors.Config{
		AllowOrigins: "http://localhost:3000, https://ae.ahmedalmaghribi.com, https://ksa.ahmedalmaghribi.com, https://qa.ahmedalmaghribi.com, https://kw.ahmedalmaghribi.com, https://bh.ahmedalmaghribi.com, https://om.ahmedalmaghribi.com",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// Tenant Middleware

	app.Use(catalogMiddleware.TenantMiddleware())

	// =========================
	// Health Check
	// =========================

	api := app.Group("/api/v1")

	app.Get("/health", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "catalog-service",
		})
	})

	// =========================
	// Routes
	// =========================

	categoryRoutes.SetupCategoryRoutes(api)

	// =========================
	// Start Server
	// =========================

	port := os.Getenv("APP_PORT")

	if port == "" {
		port = "8081"
	}

	logger.Log.Info("🚀 Catalog Service Running On Port:",
		zap.String("port", port),
	)

	logger.Log.Fatal("server stopped",
		zap.Error(app.Listen(":"+port)),
	)
}
