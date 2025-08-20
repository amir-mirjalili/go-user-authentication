package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/amir-mirjalili/go-user-authentication/docs"
	"github.com/amir-mirjalili/go-user-authentication/internal/db"
	"github.com/amir-mirjalili/go-user-authentication/internal/handlers"
	"github.com/amir-mirjalili/go-user-authentication/internal/middlewares"
	"github.com/amir-mirjalili/go-user-authentication/internal/repository"
	"github.com/amir-mirjalili/go-user-authentication/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title OTP Authentication Service API
// @version 1.0
// @description OTP-based auth API
// @contact.name API Support
// @contact.email support@example.com
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	godotenv.Load()

	// Connect to database
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer func(conn *sql.DB) {
		err := db.Close(conn)
		if err != nil {
			log.Fatalf("DB connection failed: %v", err)
		}
	}(database.DB)

	log.Println("✅ Connected to", database.Dialect)

	otpRepository := repository.NewOtpRepository(database.DB)
	userRepository := repository.NewUserRepository(database.DB)

	otpService := services.NewOTPService(otpRepository)
	userService := services.NewUserService(userRepository)
	authService := services.NewAuthService(otpService, userService, "")

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Gin setup
	r := gin.Default()
	if err := r.SetTrustedProxies(nil); err != nil {
		log.Fatalf("Failed to set trusted proxies: %v", err)
	}

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		auth.POST("/send-otp", authHandler.SendOTP)
		auth.POST("/verify-otp", authHandler.VerifyOTP)

		users := api.Group("/users")
		users.Use(middlewares.AuthMiddleware(os.Getenv("JWT_SECRET")))
		{
			users.GET("/me", userHandler.GetCurrentUser)
			users.GET("/list", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"timestamp": log.Default().Writer(),
		})
	})

	// Run server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
