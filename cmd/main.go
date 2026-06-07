// // package main

// // import (
// // 	"context"
// // 	"fmt"
// // 	"log"
// // 	"mobile-library/config"
// // 	"mobile-library/internal/handlers"
// // 	"mobile-library/internal/middleware"
// // 	"mobile-library/internal/repositories"
// // 	"mobile-library/internal/services"
// // 	"net/http"
// // 	"os"
// // 	"os/signal"
// // 	"syscall"
// // 	"time"

// // 	"github.com/gin-gonic/gin"
// // 	"github.com/redis/go-redis/v9"
// // 	"go.uber.org/zap"

// // 	// Инициализация документации (папка docs создастся после swag init)
// // 	_ "mobile-library/docs"

// // 	// Пакеты для Swagger
// // 	swaggerFiles "github.com/swaggo/files"
// // 	ginSwagger "github.com/swaggo/gin-swagger"
// // )

// // // @title           Mobile Library API
// // // @version         1.0
// // // @description     API server for Shogun's Library.
// // // @termsOfService  http://swagger.io/terms/

// // // @contact.name   Komron Nazarov
// // // @contact.url    https://github.com/shogun

// // // @host melodious-friendship-production-e718.up.railway.app
// // // @BasePath  /

// // // @securityDefinitions.apikey ApiKeyAuth
// // // @in header
// // // @name Authorization

// // func main() {
// // 	// 1. Config
// // 	cfg, err := config.Load()
// // 	if err != nil {
// // 		log.Fatalf("Failed to load config: %v", err)
// // 	}

// // 	// 2. Logger
// // 	logger, _ := zap.NewProduction()
// // 	defer logger.Sync()

// // 	// 3. DB
// // 	// db, err := repositories.NewDB(cfg.DBConnString())
// // 	// if err != nil {
// // 	// 	logger.Fatal("Failed to connect DB", zap.Error(err))
// // 	// }
// // 	// defer db.Close()
// // 	db, err := repositories.NewDB(cfg.DBConnString())
// // 	if err != nil {
// // 		// Выводим реальную ошибку из err
// // 		logger.Fatal("Failed to connect DB", zap.Error(err))
// // 	}

// // 	// 4. Run migrations
// // 	if err := runMigrations(db); err != nil {
// // 		logger.Fatal("Migration failed", zap.Error(err))
// // 	}

// // 	// 5. Redis
// // 	rdb := redis.NewClient(&redis.Options{
// // 		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
// // 		Password: cfg.RedisPassword,
// // 		DB:       cfg.RedisDB,
// // 	})

// // 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// // 	defer cancel()
// // 	if err := rdb.Ping(ctx).Err(); err != nil {
// // 		logger.Warn("Redis not available, continuing without cache", zap.Error(err))
// // 	}

// // 	// 6. Repositories
// // 	userRepo := repositories.NewUserRepository(db)
// // 	bookRepo := repositories.NewBookRepository(db)
// // 	borrowRepo := repositories.NewBorrowRepository(db)
// // 	favoriteRepo := repositories.NewFavoriteRepository(db)
// // 	notifRepo := repositories.NewNotificationRepository(db)

// // 	// 7. Services
// // 	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
// // 	userService := services.NewUserService(userRepo)
// // 	bookService := services.NewBookService(bookRepo)
// // 	borrowService := services.NewBorrowService(borrowRepo, bookRepo, notifRepo)
// // 	notifService := services.NewNotificationService(notifRepo)

// // 	// 8. Handlers
// // 	authHandler := handlers.NewAuthHandler(authService)
// // 	userHandler := handlers.NewUserHandler(userService)
// // 	bookHandler := handlers.NewBookHandler(bookService)
// // 	borrowHandler := handlers.NewBorrowHandler(borrowService)
// // 	favHandler := handlers.NewFavoriteHandler(favoriteRepo)
// // 	notifHandler := handlers.NewNotificationHandler(notifService)

// // 	// 9. Gin Setup
// // 	router := gin.New()
// // 	router.Use(gin.Recovery())
// // 	router.Use(middleware.LoggingMiddleware(logger))
// // 	router.Use(middleware.ErrorHandler())
// // 	router.Use(CORSMiddleware())

// // 	// --- Routes ---

// // 	// Swagger UI (Доступен по /swagger/index.html)
// // 	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// // 	// Public Routes
// // 	authGroup := router.Group("/auth")
// // 	{
// // 		authGroup.POST("/register", authHandler.Register)
// // 		authGroup.POST("/login", authHandler.Login)
// // 	}

// // 	publicBooks := router.Group("/books")
// // 	{
// // 		publicBooks.GET("", bookHandler.GetAll)
// // 		publicBooks.GET("/:id", bookHandler.GetByID)
// // 		publicBooks.GET("/search", bookHandler.Search)
// // 	}

// // 	// Protected Routes (JWT)
// // 	authM := middleware.AuthMiddleware(cfg.JWTSecret)

// // 	userGroup := router.Group("/user", authM)
// // 	{
// // 		userGroup.GET("/profile", userHandler.GetProfile)
// // 		userGroup.PUT("/profile", userHandler.UpdateProfile)
// // 		userGroup.PUT("/change-password", userHandler.ChangePassword)
// // 		userGroup.GET("/history", userHandler.GetHistory)
// // 	}

// // 	borrowGroup := router.Group("/borrow", authM)
// // 	{
// // 		borrowGroup.POST("/:id", borrowHandler.BorrowBook)
// // 		borrowGroup.POST("/:id/return", borrowHandler.ReturnBook)
// // 	}

// // 	favGroup := router.Group("/favorites", authM)
// // 	{
// // 		favGroup.POST("/:book_id", favHandler.Add)
// // 		favGroup.DELETE("/:book_id", favHandler.Remove)
// // 		favGroup.GET("", favHandler.List)
// // 	}

// // 	notifGroup := router.Group("/notifications", authM)
// // 	{
// // 		notifGroup.GET("", notifHandler.List)
// // 		notifGroup.PUT("/:id/read", notifHandler.MarkRead)
// // 	}

// // 	router.GET("/health", func(c *gin.Context) {
// // 		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
// // 	})

// // 	// 10. Graceful Shutdown
// // 	serverPort := cfg.ServerPort
// // 	if envPort := os.Getenv("PORT"); envPort != "" {
// // 		serverPort = envPort
// // 	}

// // 	srv := &http.Server{
// // 		Addr:    ":" + serverPort,
// // 		Handler: router,
// // 	}

// // 	go func() {
// // 		logger.Info("Server starting", zap.String("port", serverPort))
// // 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// // 			logger.Fatal("Server failed", zap.Error(err))
// // 		}
// // 	}()

// // 	quit := make(chan os.Signal, 1)
// // 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// // 	<-quit

// // 	logger.Info("Shutting down server...")

// // 	ctxShut, cancelShut := context.WithTimeout(context.Background(), 5*time.Second)
// // 	defer cancelShut()

// // 	if err := srv.Shutdown(ctxShut); err != nil {
// // 		logger.Error("Server forced shutdown", zap.Error(err))
// // 	}

// // 	if rdb != nil {
// // 		rdb.Close()
// // 		logger.Info("Redis connection closed")
// // 	}

// // 	db.Close()
// // 	logger.Info("Database connection closed")
// // 	logger.Info("Server stopped. Bye!")
// // }

// // // CORSMiddleware ...
// // func CORSMiddleware() gin.HandlerFunc {
// // 	return func(c *gin.Context) {
// // 		c.Header("Access-Control-Allow-Origin", "*")
// // 		c.Header("Access-Control-Allow-Credentials", "true")
// // 		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
// // 		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

// // 		if c.Request.Method == "OPTIONS" {
// // 			c.AbortWithStatus(204)
// // 			return
// // 		}
// // 		c.Next()
// // 	}
// // }

// // // runMigrations ...
// // func runMigrations(db *repositories.DB) error {
// // 	query := `
// //         CREATE TABLE IF NOT EXISTS users (
// //             id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL,
// //             email VARCHAR(255) UNIQUE NOT NULL, phone VARCHAR(20) NOT NULL,
// //             password VARCHAR(255) NOT NULL, created_at TIMESTAMP DEFAULT NOW()
// //         );
// //         CREATE TABLE IF NOT EXISTS books (
// //             id SERIAL PRIMARY KEY, title VARCHAR(255) NOT NULL,
// //             author VARCHAR(255) NOT NULL, description TEXT,
// //             category VARCHAR(100), year INTEGER,
// //             available_copies INTEGER DEFAULT 0,
// //             image_url VARCHAR(500), created_at TIMESTAMP DEFAULT NOW()
// //         );
// //         CREATE TABLE IF NOT EXISTS borrows (
// //             id SERIAL PRIMARY KEY,
// //             user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// //             book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
// //             borrow_date TIMESTAMP NOT NULL, due_date TIMESTAMP NOT NULL,
// //             return_date TIMESTAMP, status VARCHAR(20) DEFAULT 'active',
// //             created_at TIMESTAMP DEFAULT NOW()
// //         );
// //         CREATE TABLE IF NOT EXISTS favorites (
// //             id SERIAL PRIMARY KEY,
// //             user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// //             book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
// //             created_at TIMESTAMP DEFAULT NOW(),
// //             UNIQUE(user_id, book_id)
// //         );
// //         CREATE TABLE IF NOT EXISTS notifications (
// //             id SERIAL PRIMARY KEY,
// //             user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// //             message TEXT NOT NULL, is_read BOOLEAN DEFAULT FALSE,
// //             created_at TIMESTAMP DEFAULT NOW()
// //         );
// //     `
// // 	return db.Exec(query)
// // }

// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"mobile-library/config"
// 	"mobile-library/internal/handlers"
// 	"mobile-library/internal/middleware"
// 	"mobile-library/internal/repositories"
// 	"mobile-library/internal/services"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/redis/go-redis/v9"
// 	"go.uber.org/zap"

// 	// Инициализация документации
// 	_ "mobile-library/docs"

// 	// Пакеты для Swagger
// 	swaggerFiles "github.com/swaggo/files"
// 	ginSwagger "github.com/swaggo/gin-swagger"
// )

// // @title           Mobile Library API
// // @version         1.0
// // @description     API server for Shogun's Library.
// // @contact.name    Komron Nazarov
// // @host            melodious-friendship-production-e718.up.railway.app
// // @BasePath        /
// // @securityDefinitions.apikey ApiKeyAuth
// // @in header
// // @name Authorization
// // @description Введите токен в формате: Bearer <your_token>

// func main() {
// 	// 1. Config
// 	cfg, err := config.Load()
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}

// 	// 2. Logger
// 	logger, _ := zap.NewProduction()
// 	defer logger.Sync()

// 	// 3. DB - Подключение через исправленный DBConnString
// 	db, err := repositories.NewDB(cfg.DBConnString())
// 	if err != nil {
// 		logger.Fatal("Failed to connect DB", zap.Error(err))
// 	}
// 	defer db.Close()

// 	// 4. Run migrations
// 	if err := runMigrations(db); err != nil {
// 		logger.Fatal("Migration failed", zap.Error(err))
// 	}

// 	// 5. Redis
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
// 		Password: cfg.RedisPassword,
// 		DB:       cfg.RedisDB,
// 	})

// 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()
// 	if err := rdb.Ping(ctx).Err(); err != nil {
// 		logger.Warn("Redis not available, continuing without cache", zap.Error(err))
// 	}

// 	// 6. Repositories
// 	userRepo := repositories.NewUserRepository(db)
// 	bookRepo := repositories.NewBookRepository(db)
// 	borrowRepo := repositories.NewBorrowRepository(db)
// 	favoriteRepo := repositories.NewFavoriteRepository(db)
// 	notifRepo := repositories.NewNotificationRepository(db)

// 	// 7. Services
// 	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
// 	userService := services.NewUserService(userRepo)
// 	bookService := services.NewBookService(bookRepo)
// 	borrowService := services.NewBorrowService(borrowRepo, bookRepo, notifRepo)
// 	notifService := services.NewNotificationService(notifRepo)

// 	// 8. Handlers
// 	authHandler := handlers.NewAuthHandler(authService)
// 	userHandler := handlers.NewUserHandler(userService)
// 	bookHandler := handlers.NewBookHandler(bookService)
// 	borrowHandler := handlers.NewBorrowHandler(borrowService)
// 	favHandler := handlers.NewFavoriteHandler(favoriteRepo)
// 	notifHandler := handlers.NewNotificationHandler(notifService)

// 	// 9. Gin Setup
// 	router := gin.New()
// 	router.Use(gin.Recovery())
// 	router.Use(middleware.LoggingMiddleware(logger))
// 	router.Use(middleware.ErrorHandler())
// 	router.Use(CORSMiddleware())

// 	// --- Routes ---

// 	// Swagger UI
// 	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// 	// PUBLIC: Регистрация и логин доступны всем
// 	authGroup := router.Group("/auth")
// 	{
// 		authGroup.POST("/register", authHandler.Register)
// 		authGroup.POST("/login", authHandler.Login)
// 	}

// 	router.GET("/health", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
// 	})

// 	// PROTECTED: Все остальные роуты требуют JWT токен
// 	authM := middleware.AuthMiddleware(cfg.JWTSecret)

// 	// Книги теперь защищены (фронтенд будет доволен)
// 	bookGroup := router.Group("/books", authM)
// 	{
// 		bookGroup.GET("", bookHandler.GetAll)
// 		bookGroup.GET("/:id", bookHandler.GetByID)
// 		bookGroup.GET("/search", bookHandler.Search)
// 		bookGroup.POST("", bookHandler.Create)
// 	}

// 	userGroup := router.Group("/user", authM)
// 	{
// 		userGroup.GET("/profile", userHandler.GetProfile)
// 		userGroup.PUT("/profile", userHandler.UpdateProfile)
// 		userGroup.PUT("/change-password", userHandler.ChangePassword)
// 		userGroup.GET("/history", userHandler.GetHistory)
// 	}

// 	borrowGroup := router.Group("/borrow", authM)
// 	{
// 		borrowGroup.POST("/:id", borrowHandler.BorrowBook)
// 		borrowGroup.POST("/:id/return", borrowHandler.ReturnBook)
// 	}

// 	favGroup := router.Group("/favorites", authM)
// 	{
// 		favGroup.POST("/:book_id", favHandler.Add)
// 		favGroup.DELETE("/:book_id", favHandler.Remove)
// 		favGroup.GET("", favHandler.List)
// 	}

// 	notifGroup := router.Group("/notifications", authM)
// 	{
// 		notifGroup.GET("", notifHandler.List)
// 		notifGroup.PUT("/:id/read", notifHandler.MarkRead)
// 	}


// 	// ADMIN ROUTES: Для веб-сайта админ-панели
//     adminM := middleware.AdminOnly()
//     adminGroup := router.Group("/admin", authM, adminM)
//     {
//         // Статистика для дашборда (те самые карточки и графики)
//         adminGroup.GET("/stats", adminHandler.GetStats)

//         // Управление пользователями
//         adminGroup.GET("/members", adminHandler.GetAllMembers)
//         adminGroup.POST("/members/:id/accept", adminHandler.AcceptAdminRequest)

//         // Расширенное управление книгами
//         adminGroup.DELETE("/books/:id", bookHandler.Delete)
//         adminGroup.PUT("/books/:id", bookHandler.Update)
//     }


// 	// 10. Start Server
// 	serverPort := cfg.ServerPort
// 	if envPort := os.Getenv("PORT"); envPort != "" {
// 		serverPort = envPort
// 	}

// 	srv := &http.Server{
// 		Addr:    ":" + serverPort,
// 		Handler: router,
// 	}

// 	go func() {
// 		logger.Info("Server starting", zap.String("port", serverPort))
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			logger.Fatal("Server failed", zap.Error(err))
// 		}
// 	}()

// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit

// 	logger.Info("Shutting down server...")

// 	ctxShut, cancelShut := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancelShut()

// 	if err := srv.Shutdown(ctxShut); err != nil {
// 		logger.Error("Server forced shutdown", zap.Error(err))
// 	}

// 	db.Close()
// 	logger.Info("Server stopped. Bye!")
// }

// // // CORSMiddleware для работы с фронтендом
// // func CORSMiddleware() gin.HandlerFunc {
// // 	return func(c *gin.Context) {
// // 		c.Header("Access-Control-Allow-Origin", "*")
// // 		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
// // 		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

// // 		if c.Request.Method == "OPTIONS" {
// // 			c.AbortWithStatus(204)
// // 			return
// // 		}
// // 		c.Next()
// // 	}
// // }

// // CORSMiddleware для работы с фронтендом
// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		// Разрешаем запросы с любого источника
// 		c.Header("Access-Control-Allow-Origin", "*")

// 		// Разрешаем стандартные и кастомные заголовки, которые важны для JWT
// 		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")

// 		// Разрешаем все основные методы
// 		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")

// 		// Позволяем браузеру кэшировать результаты preflight-запроса на 12 часов
// 		c.Header("Access-Control-Max-Age", "43200")

// 		// Если это preflight-запрос (OPTIONS), просто отвечаем 204 и прерываем цепочку
// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}

// 		c.Next()
// 	}
// }

// // // runMigrations создает таблицы при первом запуске
// // func runMigrations(db *repositories.DB) error {
// // 	query := `
// //         CREATE TABLE IF NOT EXISTS users (
// //             id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL,
// //             email VARCHAR(255) UNIQUE NOT NULL, phone VARCHAR(20) NOT NULL,
// //             password VARCHAR(255) NOT NULL, created_at TIMESTAMP DEFAULT NOW()
// //         );
// //         CREATE TABLE IF NOT EXISTS books (
// //             id SERIAL PRIMARY KEY, title VARCHAR(255) NOT NULL,
// //             author VARCHAR(255) NOT NULL, description TEXT,
// //             category VARCHAR(100), year INTEGER,
// //             available_copies INTEGER DEFAULT 0,
// //             image_url VARCHAR(500), created_at TIMESTAMP DEFAULT NOW()
// //         );
// //         CREATE TABLE IF NOT EXISTS borrows (
// //             id SERIAL PRIMARY KEY,
// //             user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// //             book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
// //             borrow_date TIMESTAMP NOT NULL, due_date TIMESTAMP NOT NULL,
// //             return_date TIMESTAMP, status VARCHAR(20) DEFAULT 'active',
// //             created_at TIMESTAMP DEFAULT NOW()
// //         );
// //     `
// // 	return db.Exec(query)
// // }

// // runMigrations создает таблицы при первом запуске
// // func runMigrations(db *repositories.DB) error {
// // 	query := `
// //         CREATE TABLE IF NOT EXISTS users (
// //             id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL,
// //             email VARCHAR(255) UNIQUE NOT NULL, phone VARCHAR(20) NOT NULL,
// //             password VARCHAR(255) NOT NULL, created_at TIMESTAMP DEFAULT NOW()
// //         );
// //         CREATE TABLE IF NOT EXISTS books (
// //             id SERIAL PRIMARY KEY, title VARCHAR(255) NOT NULL,
// //             author VARCHAR(255) NOT NULL, description TEXT,
// //             category VARCHAR(100), year INTEGER,
// //             available_copies INTEGER DEFAULT 0,
// //             image_url VARCHAR(500), created_at TIMESTAMP DEFAULT NOW()
// //         );
// //         CREATE TABLE IF NOT EXISTS borrows (
// //             id SERIAL PRIMARY KEY,
// //             user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// //             book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
// //             borrow_date TIMESTAMP NOT NULL, due_date TIMESTAMP NOT NULL,
// //             return_date TIMESTAMP, status VARCHAR(20) DEFAULT 'active',
// //             created_at TIMESTAMP DEFAULT NOW()
// //         );
// //         CREATE TABLE IF NOT EXISTS favorites (
// //             id SERIAL PRIMARY KEY,
// //             user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// //             book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
// //             created_at TIMESTAMP DEFAULT NOW(),
// //             UNIQUE(user_id, book_id)
// //         );
// //         CREATE TABLE IF NOT EXISTS notifications (
// //             id SERIAL PRIMARY KEY,
// //             user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// //             message TEXT NOT NULL, is_read BOOLEAN DEFAULT FALSE,
// //             created_at TIMESTAMP DEFAULT NOW()
// //         );
// //     `
// // 	return db.Exec(query)
// // }


// func runMigrations(db *repositories.DB) error {
//     query := `
//         -- Обновляем таблицу пользователей (добавляем роль и позицию)
//         ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'user';
//         ALTER TABLE users ADD COLUMN IF NOT EXISTS job_position VARCHAR(50) DEFAULT 'Student';
//         ALTER TABLE users ADD COLUMN IF NOT EXISTS is_pending_admin BOOLEAN DEFAULT FALSE;
//         ALTER TABLE users ADD COLUMN IF NOT EXISTS date_of_birth VARCHAR(20);

//         -- Обновляем таблицу книг (добавляем статус и тех. детали для админки)
//         ALTER TABLE books ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'Available';
//         ALTER TABLE books ADD COLUMN IF NOT EXISTS page_count INTEGER DEFAULT 0;
//         ALTER TABLE books ADD COLUMN IF NOT EXISTS language VARCHAR(50) DEFAULT 'English';

//         -- Остальные таблицы (если их нет)
//         CREATE TABLE IF NOT EXISTS users (
//             id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL,
//             email VARCHAR(255) UNIQUE NOT NULL, phone VARCHAR(20) NOT NULL,
//             password VARCHAR(255) NOT NULL, created_at TIMESTAMP DEFAULT NOW()
//         );
//         CREATE TABLE IF NOT EXISTS books (
//             id SERIAL PRIMARY KEY, title VARCHAR(255) NOT NULL,
//             author VARCHAR(255) NOT NULL, description TEXT,
//             category VARCHAR(100), year INTEGER,
//             available_copies INTEGER DEFAULT 0,
//             image_url VARCHAR(500), created_at TIMESTAMP DEFAULT NOW()
//         );
//         CREATE TABLE IF NOT EXISTS borrows (
//             id SERIAL PRIMARY KEY,
//             user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
//             book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
//             borrow_date TIMESTAMP NOT NULL, due_date TIMESTAMP NOT NULL,
//             return_date TIMESTAMP, status VARCHAR(20) DEFAULT 'active',
//             created_at TIMESTAMP DEFAULT NOW()
//         );
//         CREATE TABLE IF NOT EXISTS favorites (
//             id SERIAL PRIMARY KEY,
//             user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
//             book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
//             created_at TIMESTAMP DEFAULT NOW(),
//             UNIQUE(user_id, book_id)
//         );
//         CREATE TABLE IF NOT EXISTS notifications (
//             id SERIAL PRIMARY KEY,
//             user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
//             message TEXT NOT NULL, is_read BOOLEAN DEFAULT FALSE,
//             created_at TIMESTAMP DEFAULT NOW()
//         );
//     `
//     return db.Exec(query)
// }



































































// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"mobile-library/config"
// 	"mobile-library/internal/handlers"
// 	"mobile-library/internal/middleware"
// 	"mobile-library/internal/repositories"
// 	"mobile-library/internal/services"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/redis/go-redis/v9"
// 	"go.uber.org/zap"

// 	// Инициализация документации
// 	_ "mobile-library/docs"

// 	// Пакеты для Swagger
// 	swaggerFiles "github.com/swaggo/files"
// 	ginSwagger "github.com/swaggo/gin-swagger"
// )

// // @title           Mobile Library API
// // @version         1.0
// // @description     API server for Shogun's Library (Mobile & Web Admin).
// // @contact.name     Komron Nazarov
// // @host            melodious-friendship-production-e718.up.railway.app
// // @BasePath        /
// // @securityDefinitions.apikey ApiKeyAuth
// // @in header
// // @name Authorization
// // @description Введите токен в формате: Bearer <your_token>

// func main() {
// 	// 1. Config
// 	cfg, err := config.Load()
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}

// 	// 2. Logger
// 	logger, _ := zap.NewProduction()
// 	defer logger.Sync()

// 	// 3. DB
// 	db, err := repositories.NewDB(cfg.DBConnString())
// 	if err != nil {
// 		logger.Fatal("Failed to connect DB", zap.Error(err))
// 	}
// 	defer db.Close()

// 	// 4. Run migrations (Обновлено для админки)
// 	if err := runMigrations(db); err != nil {
// 		logger.Fatal("Migration failed", zap.Error(err))
// 	}

// 	// 5. Redis
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
// 		Password: cfg.RedisPassword,
// 		DB:       cfg.RedisDB,
// 	})

// 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()
// 	if err := rdb.Ping(ctx).Err(); err != nil {
// 		logger.Warn("Redis not available, continuing without cache", zap.Error(err))
// 	}

// 	// 6. Repositories
// 	userRepo := repositories.NewUserRepository(db)
// 	bookRepo := repositories.NewBookRepository(db)
// 	borrowRepo := repositories.NewBorrowRepository(db)
// 	favoriteRepo := repositories.NewFavoriteRepository(db)
// 	notifRepo := repositories.NewNotificationRepository(db)

// 	// 7. Services
// 	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
// 	userService := services.NewUserService(userRepo)
// 	bookService := services.NewBookService(bookRepo)
// 	borrowService := services.NewBorrowService(borrowRepo, bookRepo, notifRepo)
// 	notifService := services.NewNotificationService(notifRepo)

// 	// // 8. Handlers
// 	// authHandler := handlers.NewAuthHandler(authService)
// 	// userHandler := handlers.NewUserHandler(userService)
// 	// bookHandler := handlers.NewBookHandler(bookService)
// 	// borrowHandler := handlers.NewBorrowHandler(borrowService)
// 	// favHandler := handlers.NewFavoriteHandler(favoriteRepo)
// 	// notifHandler := handlers.NewNotificationHandler(notifService)
// 	// adminHandler := handlers.NewAdminHandler(db) // <--- НОВЫЙ ХЕНДЛЕР ДЛЯ ВЕБА
// 	// // 8. Handlers
//     // authHandler := handlers.NewAuthHandler(authService)
//     // userHandler := handlers.NewUserHandler(userService)
	
    
//     // // ВАЖНО: Убедись, что все эти переменные объявлены!
//     // bookHandler := handlers.NewBookHandler(bookService) 
//     // borrowHandler := handlers.NewBorrowHandler(borrowService)
//     // favHandler := handlers.NewFavoriteHandler(favoriteRepo)
//     // notifHandler := handlers.NewNotificationHandler(notifService)
//     // adminHandler := handlers.NewAdminHandler(db, notifRepo)

// 	// // 9. Gin Setup
// 	// router := gin.New()
// 	// router.Use(gin.Recovery())
// 	// router.Use(middleware.LoggingMiddleware(logger))
// 	// router.Use(middleware.ErrorHandler())
// 	// router.Use(CORSMiddleware())






// 	// 8. Handlers
//     authHandler := handlers.NewAuthHandler(authService)
//     userHandler := handlers.NewUserHandler(userService)
//     bookHandler := handlers.NewBookHandler(bookService)
//     borrowHandler := handlers.NewBorrowHandler(borrowService)
//     favHandler := handlers.NewFavoriteHandler(favoriteRepo)
//     notifHandler := handlers.NewNotificationHandler(notifService)
//     adminHandler := handlers.NewAdminHandler(db, notifRepo)

//     // 9. Gin Setup
//     router := gin.New()
//     router.Use(gin.Recovery())
//     router.Use(middleware.LoggingMiddleware(logger))
//     router.Use(middleware.ErrorHandler())
//     router.Use(CORSMiddleware())

// 	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

//     // Обязательно определяем middleware перед использованием в группах
//     authM := middleware.AuthMiddleware(cfg.JWTSecret)
//     adminM := middleware.AdminOnly()

// 	// // --- Routes ---

// 	// // Swagger UI
// 	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// 	// // PUBLIC
// 	// authGroup := router.Group("/auth")
// 	// {
// 	// 	authGroup.POST("/register", authHandler.Register)
// 	// 	authGroup.POST("/login", authHandler.Login)
// 	// 	authGroup.POST("/members/:id/accept", adminHandler.AcceptAdminRequest) // <--- ДОБАВИЛИ СЮДА
// 	// }

// 	// router.GET("/health", func(c *gin.Context) {
// 	// 	c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
// 	// })

// 	// // PROTECTED
// 	// authM := middleware.AuthMiddleware(cfg.JWTSecret)

// 	// // Книги (Mobile/Web)
// 	// bookGroup := router.Group("/books", authM)
// 	// {
// 	// 	bookGroup.GET("", bookHandler.GetAll)
// 	// 	bookGroup.GET("/:id", bookHandler.GetByID)
// 	// 	bookGroup.GET("/search", bookHandler.Search)
// 	// 	bookGroup.POST("", bookHandler.Create) 
// 	// }

// 	// userGroup := router.Group("/user", authM)
// 	// {
// 	// 	userGroup.GET("/profile", userHandler.GetProfile)
// 	// 	userGroup.PUT("/profile", userHandler.UpdateProfile)
// 	// 	userGroup.PUT("/change-password", userHandler.ChangePassword)
// 	// 	userGroup.GET("/history", userHandler.GetHistory)
// 	// }

// 	// borrowGroup := router.Group("/borrow", authM)
// 	// {
// 	// 	borrowGroup.POST("/:id", borrowHandler.BorrowBook)
// 	// 	borrowGroup.POST("/:id/return", borrowHandler.ReturnBook)
// 	// }

// 	// favGroup := router.Group("/favorites", authM)
// 	// {
// 	// 	favGroup.POST("/:book_id", favHandler.Add)
// 	// 	favGroup.DELETE("/:book_id", favHandler.Remove)
// 	// 	favGroup.GET("", favHandler.List)
// 	// }

// 	// notifGroup := router.Group("/notifications", authM)
// 	// {
// 	// 	notifGroup.GET("", notifHandler.List)
// 	// 	notifGroup.PUT("/:id/read", notifHandler.MarkRead)
// 	// }

// 	// // --- ADMIN ROUTES (WEB DASHBOARD) ---
// 	// adminM := middleware.AdminOnly()
// 	// adminGroup := router.Group("/admin", authM, adminM)
// 	// {
// 	// 	// Статистика для графиков на видео
// 	// 	adminGroup.GET("/stats", adminHandler.GetStats)

// 	// 	// Управление пользователями
// 	// 	adminGroup.GET("/members", adminHandler.GetAllMembers)
// 	// 	adminGroup.POST("/members/:id/accept", adminHandler.AcceptAdminRequest)

// 	// 	// Удаление/Правка книг через веб
// 	// 	adminGroup.DELETE("/books/:id", bookHandler.Delete)
// 	// 	adminGroup.PUT("/books/:id", bookHandler.Update)
// 	// }

// // 	// --- Routes ---

// //     // Swagger UI
// //     router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// //     // PUBLIC
// //     authGroup := router.Group("/auth")
// //     {
// //         authGroup.POST("/register", authHandler.Register)
// //         authGroup.POST("/login", authHandler.Login)
// //         // Убрали AcceptAdminRequest отсюда, так как он должен быть защищен
// //     }

// //     router.GET("/health", func(c *gin.Context) {
// //         c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
// //     })

// //     // PROTECTED MIDDLEWARE
// //     authM := middleware.AuthMiddleware(cfg.JWTSecret)
// //     adminM := middleware.AdminOnly()

// //     // Книги (Mobile/Web)
// //     bookGroup := router.Group("/books")
// //     bookGroup.Use(authM)
// //     {
// //         bookGroup.GET("", bookHandler.GetAll)
// //         bookGroup.GET("/:id", bookHandler.GetByID)
// //         bookGroup.GET("/search", bookHandler.Search)
// //         bookGroup.POST("", bookHandler.Create) 
// //     }

// //     userGroup := router.Group("/user")
// //     userGroup.Use(authM)
// //     {
// //         userGroup.GET("/profile", userHandler.GetProfile)
// //         userGroup.PUT("/profile", userHandler.UpdateProfile)
// //         userGroup.PUT("/change-password", userHandler.ChangePassword)
// //         userGroup.GET("/history", userHandler.GetHistory)
// //     }

// //     borrowGroup := router.Group("/borrow")
// //     borrowGroup.Use(authM)
// //     {
// //         borrowGroup.POST("/:id", borrowHandler.BorrowBook)
// //         borrowGroup.POST("/:id/return", borrowHandler.ReturnBook)
// //     }

// //     favGroup := router.Group("/favorites")
// //     favGroup.Use(authM)
// //     {
// //         favGroup.POST("/:book_id", favHandler.Add)
// //         favGroup.DELETE("/:book_id", favHandler.Remove)
// //         favGroup.GET("", favHandler.List)
// //     }

// //     notifGroup := router.Group("/notifications")
// //     notifGroup.Use(authM)
// //     {
// //         notifGroup.GET("", notifHandler.List)
// //         notifGroup.PUT("/:id/read", notifHandler.MarkRead)
// //     }

// //     // --- ADMIN ROUTES (WEB DASHBOARD) ---
// //     // Теперь они надежно защищены и имеют явную группу
// //     // adminGroup := router.Group("/admin")
// //     // adminGroup.Use(authM, adminM)
// //     // {
// //     //     adminGroup.GET("/stats", adminHandler.GetStats)
// //     //     adminGroup.GET("/members", adminHandler.GetAllMembers)
// //     //     adminGroup.POST("/members/:id/accept", adminHandler.AcceptAdminRequest)
// //     //     adminGroup.DELETE("/books/:id", bookHandler.Delete)
// //     //     adminGroup.PUT("/books/:id", bookHandler.Update)
// //     // }

// // // --- ADMIN ROUTES (WEB DASHBOARD) ---
// //     // Обязательно используем middleware для безопасности!
// // //     adminGroup := router.Group("/admin")
// // //     adminGroup.Use(authM, adminM) 
// // //     {
// // //         // Stats & Dashboard
// // //         adminGroup.GET("/stats", adminHandler.GetStats)
// // //         adminGroup.GET("/chart", adminHandler.GetChartData)     // Нужно реализовать в admin_handler
// // //         adminGroup.GET("/overdue", adminHandler.GetOverdueMembers) // Нужно реализовать в admin_handler

// // //         // Books Management
// // //         booksGroup := adminGroup.Group("/books")
// // //         {
// // //             booksGroup.GET("", bookHandler.GetAll)
// // //             booksGroup.POST("", bookHandler.Create)
// // //             booksGroup.PUT("/:id", bookHandler.Update)
// // //             booksGroup.DELETE("/:id", bookHandler.Delete)
// // // 			adminGroup.GET("/members/:user_id/books", adminHandler.GetUserBookshelf)
// // // 			// Заявки на получение
// // // adminGroup.GET("/receive-requests", adminHandler.GetReceiveRequests) // Метод Get на чтение списка
// // // adminGroup.POST("/receive-requests/:id/accept", adminHandler.AcceptReceiveRequest)
// // // adminGroup.DELETE("/receive-requests/:id", adminHandler.DeleteReceiveRequest) // Просто отклонить
// // // // Заявки на возврат
// // // adminGroup.GET("/return-requests", adminHandler.GetReturnRequests) // Метод GET (нужно добавить)
// // // adminGroup.POST("/return-requests/:id/accept", adminHandler.AcceptReturnRequest)
// // // adminGroup.DELETE("/return-requests/:id", adminHandler.DeleteReturnRequest) // Отклонить
// // //         }

// // // --- ADMIN ROUTES (WEB DASHBOARD) ---
// // adminGroup := router.Group("/admin")
// // adminGroup.Use(authM, adminM) 
// // {
// //     // 1. Stats & Dashboard
// //     adminGroup.GET("/stats", adminHandler.GetStats)
// //     adminGroup.GET("/chart", adminHandler.GetChartData)
// //     adminGroup.GET("/overdue", adminHandler.GetOverdueMembers)

// //     // 2. Books Management
// //     booksGroup := adminGroup.Group("/books")
// //     {
// //         booksGroup.GET("", bookHandler.GetAll)
// //         booksGroup.POST("", bookHandler.Create)
// //         booksGroup.PUT("/:id", bookHandler.Update)
// //         booksGroup.DELETE("/:id", bookHandler.Delete)
// //     }

// //     // 3. Filters Management
// //     filtersGroup := adminGroup.Group("/filters")
// //     {
// //         filtersGroup.GET("", bookHandler.GetFilters)
// //         filtersGroup.POST("", bookHandler.AddFilter)
// //         filtersGroup.PUT("/:id", bookHandler.UpdateFilter)
// //         filtersGroup.DELETE("/:id", bookHandler.DeleteFilter)
// //     }

// //     // 4. Members & User Bookshelf
// //     adminGroup.GET("/members", adminHandler.GetAllMembers)
// //     adminGroup.POST("/members/:id/accept", adminHandler.AcceptAdminRequest)
// //     adminGroup.DELETE("/members/:id", adminHandler.DeleteMember)
// //     adminGroup.GET("/members/:user_id/books", adminHandler.GetUserBookshelf)

// //     // 5. Receive Book Requests (Заявки на получение)
// //     adminGroup.GET("/receive-requests", adminHandler.GetReceiveRequests)
// //     adminGroup.POST("/receive-requests/:id/accept", adminHandler.AcceptReceiveRequest)
// //     adminGroup.DELETE("/receive-requests/:id", adminHandler.DeleteReceiveRequest)

// //     // 6. Return Book Requests (Заявки на возврат)
// //     adminGroup.GET("/return-requests", adminHandler.GetReturnRequests)
// //     adminGroup.POST("/return-requests/:id/accept", adminHandler.AcceptReturnRequest)
// //     adminGroup.DELETE("/return-requests/:id", adminHandler.DeleteReturnRequest)
    
// //     // 7. Notifications (Уведомления)
// //     adminGroup.GET("/notifications", adminHandler.GetNotifications)
// //     adminGroup.POST("/notifications", adminHandler.AddNotification)
// // }

// //         // Фильтры (CRUD)
// //     filtersGroup := adminGroup.Group("/filters")
// //     {
// //         filtersGroup.GET("", bookHandler.GetFilters)
// //         filtersGroup.POST("", bookHandler.AddFilter)
// //         filtersGroup.PUT("/:id", bookHandler.UpdateFilter)
// //         filtersGroup.DELETE("/:id", bookHandler.DeleteFilter)
// //     }
// //     }
// // // --- ADMIN ROUTES (WEB DASHBOARD) ---
// // 	adminGroup := router.Group("/admin", authM, adminM)
// // 	{
// // 		// 1. Статистика
// // 		adminGroup.GET("/stats", adminHandler.GetStats)
// // 		adminGroup.GET("/chart", adminHandler.GetChartData)
// // 		adminGroup.GET("/overdue", adminHandler.GetOverdueMembers)

// // 		// 2. Управление книгами
// // 		adminGroup.GET("/books", bookHandler.GetAll)
// // 		adminGroup.POST("/books", bookHandler.Create)
// // 		adminGroup.PUT("/books/:id", bookHandler.Update)
// // 		adminGroup.DELETE("/books/:id", bookHandler.Delete)

// // 		// 3. Управление фильтрами
// // 		adminGroup.GET("/filters", bookHandler.GetFilters)
// // 		adminGroup.POST("/filters", bookHandler.AddFilter)
// // 		adminGroup.PUT("/filters/:id", bookHandler.UpdateFilter)
// // 		adminGroup.DELETE("/filters/:id", bookHandler.DeleteFilter)

// // 		// 4. Управление пользователями
// // 		adminGroup.GET("/members", adminHandler.GetAllMembers)
// // 		adminGroup.POST("/members/:id/accept", adminHandler.AcceptAdminRequest)
// // 		adminGroup.DELETE("/members/:id", adminHandler.DeleteMember)
// // 		adminGroup.GET("/members/:user_id/books", adminHandler.GetUserBookshelf)

// // 		// 5. Заявки
// // 		adminGroup.GET("/receive-requests", adminHandler.GetReceiveRequests)
// // 		adminGroup.POST("/receive-requests/:id/accept", adminHandler.AcceptReceiveRequest)
// // 		adminGroup.DELETE("/receive-requests/:id", adminHandler.DeleteReceiveRequest)

// // 		adminGroup.GET("/return-requests", adminHandler.GetReturnRequests)
// // 		adminGroup.POST("/return-requests/:id/accept", adminHandler.AcceptReturnRequest)
// // 		adminGroup.DELETE("/return-requests/:id", adminHandler.DeleteReturnRequest)

// // 		// 6. Уведомления
// // 		adminGroup.GET("/notifications", adminHandler.GetNotifications)
// // 		adminGroup.POST("/notifications", adminHandler.AddNotification)
// // 	}

// // 	// 10. Start Server


// // --- ADMIN ROUTES (WEB DASHBOARD) ---
//     adminGroup := router.Group("/admin", authM, adminM)
//     {
//         adminGroup.GET("/stats", adminHandler.GetStats)
//         adminGroup.GET("/chart", adminHandler.GetChartData)
//         adminGroup.GET("/overdue", adminHandler.GetOverdueMembers)
//         adminGroup.GET("/books", bookHandler.GetAll)
//         adminGroup.POST("/books", bookHandler.Create)
//         adminGroup.PUT("/books/:id", bookHandler.Update)
//         adminGroup.DELETE("/books/:id", bookHandler.Delete)
//         adminGroup.GET("/filters", bookHandler.GetFilters)
//         adminGroup.POST("/filters", bookHandler.AddFilter)
//         adminGroup.PUT("/filters/:id", bookHandler.UpdateFilter)
//         adminGroup.DELETE("/filters/:id", bookHandler.DeleteFilter)
//         adminGroup.GET("/members", adminHandler.GetAllMembers)
//         adminGroup.POST("/members/:id/accept", adminHandler.AcceptAdminRequest)
//         adminGroup.DELETE("/members/:id", adminHandler.DeleteMember)
//         adminGroup.GET("/members/:user_id/books", adminHandler.GetUserBookshelf)
//         adminGroup.GET("/receive-requests", adminHandler.GetReceiveRequests)
//         adminGroup.POST("/receive-requests/:id/accept", adminHandler.AcceptReceiveRequest)
//         adminGroup.DELETE("/receive-requests/:id", adminHandler.DeleteReceiveRequest)
//         adminGroup.GET("/return-requests", adminHandler.GetReturnRequests)
//         adminGroup.POST("/return-requests/:id/accept", adminHandler.AcceptReturnRequest)
//         adminGroup.DELETE("/return-requests/:id", adminHandler.DeleteReturnRequest)
//         adminGroup.GET("/notifications", adminHandler.GetNotifications)
//         adminGroup.POST("/notifications", adminHandler.AddNotification)
//     }

// 	// Смена пароля (переиспользуем готовый хендлер юзера, так как админ — это тоже юзер с ролью admin)
//     adminGroup.POST("/profile/change-password", userHandler.ChangePassword)
// 	adminGroup.GET("/profile/admins", adminHandler.GetAdminManagement) // Управление админами
//     // Заглушки, чтобы Go не ругался на неиспользуемые переменные
//     // (Используй эти переменные в своих будущих публичных роутах)
//     _, _, _, _, _ = authHandler, userHandler, borrowHandler, favHandler, notifHandler

// 	// 10. Start Server
// 	serverPort := cfg.ServerPort
// 	if envPort := os.Getenv("PORT"); envPort != "" {
// 		serverPort = envPort
// 	}

// 	srv := &http.Server{
// 		Addr:    ":" + serverPort,
// 		Handler: router,
// 	}

// 	go func() {
// 		logger.Info("Server starting", zap.String("port", serverPort))
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			logger.Fatal("Server failed", zap.Error(err))
// 		}
// 	}()

// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit

// 	logger.Info("Shutting down server...")
// 	ctxShut, cancelShut := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancelShut()

// 	if err := srv.Shutdown(ctxShut); err != nil {
// 		logger.Error("Server forced shutdown", zap.Error(err))
// 	}

// 	db.Close()
// 	logger.Info("Server stopped. Bye!")
// }

// // CORSMiddleware
// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Header("Access-Control-Allow-Origin", "*")
// 		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
// 		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")
// 		c.Header("Access-Control-Max-Age", "43200")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	}
// }

// // runMigrations
// func runMigrations(db *repositories.DB) error {
// 	query := `
// 		-- Добавляем колонки для админки и юзеров
// 		ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'user';
// 		ALTER TABLE users ADD COLUMN IF NOT EXISTS job_position VARCHAR(50) DEFAULT 'Student';
// 		ALTER TABLE users ADD COLUMN IF NOT EXISTS is_pending_admin BOOLEAN DEFAULT FALSE;
// 		ALTER TABLE users ADD COLUMN IF NOT EXISTS date_of_birth VARCHAR(20);

// 		ALTER TABLE books ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'Available';
// 		ALTER TABLE books ADD COLUMN IF NOT EXISTS page_count INTEGER DEFAULT 0;
// 		ALTER TABLE books ADD COLUMN IF NOT EXISTS language VARCHAR(50) DEFAULT 'English';

// 		-- Основные таблицы
// 		CREATE TABLE IF NOT EXISTS users (
// 			id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL,
// 			email VARCHAR(255) UNIQUE NOT NULL, phone VARCHAR(20) NOT NULL,
// 			password VARCHAR(255) NOT NULL, created_at TIMESTAMP DEFAULT NOW(),
// 			role VARCHAR(20) DEFAULT 'user', job_position VARCHAR(50) DEFAULT 'Student'
// 		);
// 		CREATE TABLE IF NOT EXISTS books (
// 			id SERIAL PRIMARY KEY, title VARCHAR(255) NOT NULL,
// 			author VARCHAR(255) NOT NULL, description TEXT,
// 			category VARCHAR(100), year INTEGER,
// 			available_copies INTEGER DEFAULT 0,
// 			image_url VARCHAR(500), created_at TIMESTAMP DEFAULT NOW(),
// 			status VARCHAR(20) DEFAULT 'Available'
// 		);
// 		CREATE TABLE IF NOT EXISTS borrows (
// 			id SERIAL PRIMARY KEY,
// 			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// 			book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
// 			borrow_date TIMESTAMP NOT NULL, due_date TIMESTAMP NOT NULL,
// 			return_date TIMESTAMP, status VARCHAR(20) DEFAULT 'active',
// 			created_at TIMESTAMP DEFAULT NOW()
// 		);
// 		CREATE TABLE IF NOT EXISTS favorites (
// 			id SERIAL PRIMARY KEY,
// 			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// 			book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
// 			created_at TIMESTAMP DEFAULT NOW(),
// 			UNIQUE(user_id, book_id)
// 		);
// 		CREATE TABLE IF NOT EXISTS notifications (
// 			id SERIAL PRIMARY KEY,
// 			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// 			message TEXT NOT NULL, is_read BOOLEAN DEFAULT FALSE,
// 			created_at TIMESTAMP DEFAULT NOW()
// 		);

// 		CREATE TABLE IF NOT EXISTS categories (
//     		id SERIAL PRIMARY KEY,
//     		name VARCHAR(100) UNIQUE NOT NULL
// 		);
// 	`
// 	return db.Exec(query)
// }


















// package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"mobile-library/config"
// 	"mobile-library/internal/handlers"
// 	"mobile-library/internal/middleware"
// 	"mobile-library/internal/repositories"
// 	"mobile-library/internal/services"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/redis/go-redis/v9"
// 	"go.uber.org/zap"

// 	// Init swagger docs
// 	_ "mobile-library/docs"

// 	swaggerFiles "github.com/swaggo/files"
// 	ginSwagger "github.com/swaggo/gin-swagger"
// )

// // @title           Mobile Library API
// // @version         1.0
// // @description     API server for Shogun's Library (Mobile & Web Admin).
// // @contact.name    Komron Nazarov
// // @host            melodious-friendship-production-e718.up.railway.app
// // @BasePath        /
// // @securityDefinitions.apikey ApiKeyAuth
// // @in header
// // @name Authorization
// // @description Введите токен в формате: Bearer <your_token>

// func main() {
// 	// 1. Config
// 	cfg, err := config.Load()
// 	if err != nil {
// 		log.Fatalf("Failed to load config: %v", err)
// 	}

// 	// 2. Logger
// 	logger, _ := zap.NewProduction()
// 	defer logger.Sync()

// 	// 3. DB
// 	db, err := repositories.NewDB(cfg.DBConnString())
// 	if err != nil {
// 		logger.Fatal("Failed to connect DB", zap.Error(err))
// 	}
// 	defer db.Close()

// 	// 4. Run migrations
// 	if err := runMigrations(db); err != nil {
// 		logger.Fatal("Migration failed", zap.Error(err))
// 	}

// 	// 5. Redis
// 	rdb := redis.NewClient(&redis.Options{
// 		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
// 		Password: cfg.RedisPassword,
// 		DB:       cfg.RedisDB,
// 	})

// 	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
// 	defer cancel()
// 	if err := rdb.Ping(ctx).Err(); err != nil {
// 		logger.Warn("Redis not available, continuing without cache", zap.Error(err))
// 	}

// 	// 6. Repositories
// 	userRepo := repositories.NewUserRepository(db)
// 	bookRepo := repositories.NewBookRepository(db)
// 	borrowRepo := repositories.NewBorrowRepository(db)
// 	favoriteRepo := repositories.NewFavoriteRepository(db)
// 	notifRepo := repositories.NewNotificationRepository(db)

// 	// 7. Services
// 	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
// 	userService := services.NewUserService(userRepo)
// 	bookService := services.NewBookService(bookRepo)
// 	borrowService := services.NewBorrowService(borrowRepo, bookRepo, notifRepo)
// 	notifService := services.NewNotificationService(notifRepo)

// 	// 8. Handlers
// 	authHandler := handlers.NewAuthHandler(authService)
// 	userHandler := handlers.NewUserHandler(userService)
// 	bookHandler := handlers.NewBookHandler(bookService)
// 	borrowHandler := handlers.NewBorrowHandler(borrowService)
// 	favHandler := handlers.NewFavoriteHandler(favoriteRepo)
// 	notifHandler := handlers.NewNotificationHandler(notifService)
// 	adminHandler := handlers.NewAdminHandler(db, notifRepo)

// 	// 9. Gin Setup
// 	router := gin.New()
// 	router.Use(gin.Recovery())
// 	router.Use(middleware.LoggingMiddleware(logger))
// 	router.Use(middleware.ErrorHandler())
// 	router.Use(CORSMiddleware())

// 	// Swagger UI
// 	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// 	// Health check
// 	router.GET("/health", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
// 	})

// 	// =============================================================================
// 	// PUBLIC ROUTES
// 	// =============================================================================
// 	authGroup := router.Group("/auth")
// 	{
// 		authGroup.POST("/register", authHandler.Register)
// 		authGroup.POST("/login", authHandler.Login)
// 	}

// 	// Middleware Definition
// 	authM := middleware.AuthMiddleware(cfg.JWTSecret)
// 	adminM := middleware.AdminOnly()

// 	// =============================================================================
// 	// PROTECTED MOBILE APPLICATION ROUTES
// 	// =============================================================================
// 	bookGroup := router.Group("/books")
// 	bookGroup.Use(authM)
// 	{
// 		bookGroup.GET("", bookHandler.GetAll)
// 		bookGroup.GET("/:id", bookHandler.GetByID)
// 		bookGroup.GET("/search", bookHandler.Search)
// 		bookGroup.POST("", bookHandler.Create)
// 	}

// 	userGroup := router.Group("/user")
// 	userGroup.Use(authM)
// 	{
// 		userGroup.GET("/profile", userHandler.GetProfile)
// 		userGroup.PUT("/profile", userHandler.UpdateProfile)
// 		userGroup.PUT("/change-password", userHandler.ChangePassword)
// 		userGroup.GET("/history", userHandler.GetHistory)
// 	}

// 	borrowGroup := router.Group("/borrow")
// 	borrowGroup.Use(authM)
// 	{
// 		borrowGroup.POST("/:id", borrowHandler.BorrowBook)
// 		borrowGroup.POST("/:id/return", borrowHandler.ReturnBook)
// 	}

// 	favGroup := router.Group("/favorites")
// 	favGroup.Use(authM)
// 	{
// 		favGroup.POST("/:book_id", favHandler.Add)
// 		favGroup.DELETE("/:book_id", favHandler.Remove)
// 		favGroup.GET("", favHandler.List)
// 	}

// 	notifGroup := router.Group("/notifications")
// 	notifGroup.Use(authM)
// 	{
// 		notifGroup.GET("", notifHandler.List)
// 		notifGroup.PUT("/:id/read", notifHandler.MarkRead)
// 	}

// 	// =============================================================================
// 	// ADMIN ROUTES (WEB DASHBOARD)
// 	// =============================================================================
// 	adminGroup := router.Group("/admin")
// 	adminGroup.Use(authM, adminM)
// 	{
// 		// 1. Stats & Dashboard
// 		adminGroup.GET("/stats", adminHandler.GetStats)
// 		adminGroup.GET("/chart", adminHandler.GetChartData)
// 		adminGroup.GET("/overdue", adminHandler.GetOverdueMembers)

// 		// 2. Books Management
// 		booksGroup := adminGroup.Group("/books")
// 		{
// 			booksGroup.GET("", bookHandler.GetAll)
// 			booksGroup.POST("", bookHandler.Create)
// 			booksGroup.PUT("/:id", bookHandler.Update)
// 			booksGroup.DELETE("/:id", bookHandler.Delete)
// 		}

// 		// 3. Filters Management (Временные заглушки через adminHandler, чтобы компилировалось)
// 		// filtersGroup := adminGroup.Group("/filters")
// 		// {
// 		// 	filtersGroup.GET("", adminHandler.GetChartData) // Заглушка, пока нет GetFilters
// 		// 	filtersGroup.POST("", adminHandler.AcceptReturnRequest)
// 		// 	filtersGroup.PUT("/:id", adminHandler.AcceptReturnRequest)
// 		// 	filtersGroup.DELETE("/:id", adminHandler.DeleteReturnRequest)
// 		// }

// // 3. Filters Management
// 		filtersGroup := adminGroup.Group("/filters")
// 		{
// 			filtersGroup.GET("", adminHandler.GetFilters)
// 			filtersGroup.POST("", adminHandler.AddFilter)
// 			filtersGroup.PUT("/:id", adminHandler.UpdateFilter)
// 			filtersGroup.DELETE("/:id", adminHandler.DeleteFilter)
// 		}

// 		// 4. Members & User Bookshelf
// 		adminGroup.GET("/members", adminHandler.GetAllMembers)
// 		adminGroup.POST("/members/:id/accept", adminHandler.AcceptAdminRequest)
// 		adminGroup.DELETE("/members/:id", adminHandler.DeleteMember)
// 		adminGroup.GET("/members/:user_id/books", adminHandler.GetUserBookshelf)

// 		// 5. Receive Book Requests
// 		adminGroup.GET("/receive-requests", adminHandler.GetReceiveRequests)
// 		adminGroup.POST("/receive-requests/:id/accept", adminHandler.AcceptReceiveRequest)
// 		adminGroup.DELETE("/receive-requests/:id", adminHandler.DeleteReceiveRequest)

// 		// 6. Return Book Requests
// 		adminGroup.GET("/return-requests", adminHandler.GetReturnRequests)
// 		adminGroup.POST("/return-requests/:id/accept", adminHandler.AcceptReturnRequest)
// 		adminGroup.DELETE("/return-requests/:id", adminHandler.DeleteReturnRequest)

// 		// 7. Notifications
// 		adminGroup.GET("/notifications", adminHandler.GetNotifications)
// 		adminGroup.POST("/notifications", adminHandler.AddNotification)

// 		// 8. Profile & Admin Management
// 		adminGroup.POST("/profile/change-password", userHandler.ChangePassword)
// 		adminGroup.GET("/profile/admins", adminHandler.GetAdminManagement)
// 	}

// 	// 10. Start Server
// 	serverPort := cfg.ServerPort
// 	if envPort := os.Getenv("PORT"); envPort != "" {
// 		serverPort = envPort
// 	}

// 	srv := &http.Server{
// 		Addr:    ":" + serverPort,
// 		Handler: router,
// 	}

// 	go func() {
// 		logger.Info("Server starting", zap.String("port", serverPort))
// 		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			logger.Fatal("Server failed", zap.Error(err))
// 		}
// 	}()

// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit

// 	logger.Info("Shutting down server...")
// 	ctxShut, cancelShut := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancelShut()

// 	if err := srv.Shutdown(ctxShut); err != nil {
// 		logger.Error("Server forced shutdown", zap.Error(err))
// 	}

// 	db.Close()
// 	logger.Info("Server stopped. Bye!")
// }

// // CORSMiddleware
// func CORSMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Header("Access-Control-Allow-Origin", "*")
// 		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
// 		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")
// 		c.Header("Access-Control-Max-Age", "43200")

// 		if c.Request.Method == "OPTIONS" {
// 			c.AbortWithStatus(204)
// 			return
// 		}
// 		c.Next()
// 	}
// }

// // // runMigrations
// // func runMigrations(db *repositories.DB) error {
// // 	query := `
// //         ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'user';
// //         ALTER TABLE users ADD COLUMN IF NOT EXISTS job_position VARCHAR(50) DEFAULT 'Student';
// //         ALTER TABLE users ADD COLUMN IF NOT EXISTS is_pending_admin BOOLEAN DEFAULT FALSE;
// //         ALTER TABLE users ADD COLUMN IF NOT EXISTS date_of_birth VARCHAR(20);

// //         ALTER TABLE books ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'Available';
// //         ALTER TABLE books ADD COLUMN IF NOT EXISTS page_count INTEGER DEFAULT 0;
// //         ALTER TABLE books ADD COLUMN IF NOT EXISTS language VARCHAR(50) DEFAULT 'English';

// //         CREATE TABLE IF NOT EXISTS users (
// //             id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL,
// //             email VARCHAR(255) UNIQUE NOT NULL, phone VARCHAR(20) NOT NULL,
// //             password VARCHAR(255) NOT NULL, created_at TIMESTAMP DEFAULT NOW(),
// //             role VARCHAR(20) DEFAULT 'user', job_position VARCHAR(50) DEFAULT 'Student'
// //         );
// //         CREATE TABLE IF NOT EXISTS books (
// //             id SERIAL PRIMARY KEY, title VARCHAR(255) NOT NULL,
// //             author VARCHAR(255) NOT NULL, description TEXT,
// //             category VARCHAR(100), year INTEGER,
// //             available_copies INTEGER DEFAULT 0,
// //             image_url VARCHAR(500), created_at TIMESTAMP DEFAULT NOW(),
// //             status VARCHAR(20) DEFAULT 'Available'
// //         );
// //         CREATE TABLE IF NOT EXISTS borrows (
// //             id SERIAL PRIMARY KEY,
// //             user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// //             book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
// //             borrow_date TIMESTAMP NOT NULL, due_date TIMESTAMP NOT NULL,
// //             return_date TIMESTAMP, status VARCHAR(20) DEFAULT 'active',
// //             created_at TIMESTAMP DEFAULT NOW()
// //         );
// //         CREATE TABLE IF NOT EXISTS favorites (
// //             id SERIAL PRIMARY KEY,
// //             user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// //             book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
// //             created_at TIMESTAMP DEFAULT NOW(),
// //             UNIQUE(user_id, book_id)
// //         );
// //         CREATE TABLE IF NOT EXISTS notifications (
// //             id SERIAL PRIMARY KEY,
// //             user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// //             message TEXT NOT NULL, is_read BOOLEAN DEFAULT FALSE,
// //             created_at TIMESTAMP DEFAULT NOW()
// //         );
// //         CREATE TABLE IF NOT EXISTS categories (
// //             id SERIAL PRIMARY KEY,
// //             name VARCHAR(100) UNIQUE NOT NULL
// //         );
// //     `
// // 	return db.Exec(query)
// // }
// // runMigrations
// func runMigrations(db *repositories.DB) error {
// 	query := `
// 		ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'user';
// 		ALTER TABLE users ADD COLUMN IF NOT EXISTS job_position VARCHAR(50) DEFAULT 'Student';
// 		ALTER TABLE users ADD COLUMN IF NOT EXISTS is_pending_admin BOOLEAN DEFAULT FALSE;
// 		ALTER TABLE users ADD COLUMN IF NOT EXISTS date_of_birth VARCHAR(20);

// 		ALTER TABLE books ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'Available';
// 		ALTER TABLE books ADD COLUMN IF NOT EXISTS page_count INTEGER DEFAULT 0;
// 		ALTER TABLE books ADD COLUMN IF NOT EXISTS language VARCHAR(50) DEFAULT 'English';
// 		ALTER TABLE books ADD COLUMN IF NOT EXISTS bg_image_url VARCHAR(500) DEFAULT '';

// 		-- Добавляем колонки для гибких уведомлений, как просит фронтенд
// 		ALTER TABLE notifications ADD COLUMN IF NOT EXISTS type VARCHAR(20) DEFAULT 'news';
// 		ALTER TABLE notifications ADD COLUMN IF NOT EXISTS title VARCHAR(255) DEFAULT '';
// 		ALTER TABLE notifications ADD COLUMN IF NOT EXISTS notification_image_url VARCHAR(500) DEFAULT '';

// 		CREATE TABLE IF NOT EXISTS users (
// 			id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL,
// 			email VARCHAR(255) UNIQUE NOT NULL, phone VARCHAR(20) NOT NULL,
// 			password VARCHAR(255) NOT NULL, created_at TIMESTAMP DEFAULT NOW(),
// 			role VARCHAR(20) DEFAULT 'user', job_position VARCHAR(50) DEFAULT 'Student'
// 		);
// 		CREATE TABLE IF NOT EXISTS books (
// 			id SERIAL PRIMARY KEY, title VARCHAR(255) NOT NULL,
// 			author VARCHAR(255) NOT NULL, description TEXT,
// 			category VARCHAR(100), year INTEGER,
// 			available_copies INTEGER DEFAULT 0,
// 			image_url VARCHAR(500), created_at TIMESTAMP DEFAULT NOW(),
// 			status VARCHAR(20) DEFAULT 'Available',
// 			bg_image_url VARCHAR(500) DEFAULT ''
// 		);
// 		CREATE TABLE IF NOT EXISTS borrows (
// 			id SERIAL PRIMARY KEY,
// 			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// 			book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
// 			borrow_date TIMESTAMP NOT NULL, due_date TIMESTAMP NOT NULL,
// 			return_date TIMESTAMP, status VARCHAR(20) DEFAULT 'active',
// 			created_at TIMESTAMP DEFAULT NOW()
// 		);
// 		CREATE TABLE IF NOT EXISTS favorites (
// 			id SERIAL PRIMARY KEY,
// 			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// 			book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
// 			created_at TIMESTAMP DEFAULT NOW(),
// 			UNIQUE(user_id, book_id)
// 		);
// 		CREATE TABLE IF NOT EXISTS notifications (
// 			id SERIAL PRIMARY KEY,
// 			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// 			message TEXT NOT NULL, is_read BOOLEAN DEFAULT FALSE,
// 			created_at TIMESTAMP DEFAULT NOW(),
// 			type VARCHAR(20) DEFAULT 'news',
// 			title VARCHAR(255) DEFAULT '',
// 			notification_image_url VARCHAR(500) DEFAULT ''
// 		);
// 		CREATE TABLE IF NOT EXISTS categories (
// 			id SERIAL PRIMARY KEY,
// 			name VARCHAR(100) UNIQUE NOT NULL
// 		);

// 		-- НОВАЯ ТАБЛИЦА ДЛЯ ОТЗЫВОВ И РЕЙТИНГА КНИГ
// 		CREATE TABLE IF NOT EXISTS reviews (
// 			id SERIAL PRIMARY KEY,
// 			book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
// 			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
// 			rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
// 			review TEXT,
// 			review_category VARCHAR(100) DEFAULT '',
// 			created_at TIMESTAMP DEFAULT NOW()
// 		);
// 	`
// 	return db.Exec(query)
// }








package main

import (
	"context"
	"fmt"
	"log"
	"mobile-library/config"
	"mobile-library/internal/handlers"
	"mobile-library/internal/middleware"
	"mobile-library/internal/repositories"
	"mobile-library/internal/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	// Init swagger docs
	_ "mobile-library/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Mobile Library API
// @version         1.0
// @description     API server for Shogun's Library (Mobile & Web Admin).
// @contact.name    Komron Nazarov
// @host            melodious-friendship-production-e718.up.railway.app
// @BasePath        /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Введите токен в формате: Bearer <your_token>

func main() {
	// 1. Config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 3. DB
	db, err := repositories.NewDB(cfg.DBConnString())
	if err != nil {
		logger.Fatal("Failed to connect DB", zap.Error(err))
	}
	defer db.Close()

	// 4. Run migrations
	if err := runMigrations(db); err != nil {
		logger.Fatal("Migration failed", zap.Error(err))
	}

	// 5. Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Warn("Redis not available, continuing without cache", zap.Error(err))
	}

	// 6. Repositories
	userRepo := repositories.NewUserRepository(db)
	bookRepo := repositories.NewBookRepository(db)
	borrowRepo := repositories.NewBorrowRepository(db)
	favoriteRepo := repositories.NewFavoriteRepository(db)
	notifRepo := repositories.NewNotificationRepository(db)
	reviewRepo := repositories.NewReviewRepository(db) // Инициализируем репозиторий отзывов

	// 7. Services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)
	userService := services.NewUserService(userRepo)
	bookService := services.NewBookService(bookRepo)
	borrowService := services.NewBorrowService(borrowRepo, bookRepo, notifRepo)
	notifService := services.NewNotificationService(notifRepo)

	// 8. Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	bookHandler := handlers.NewBookHandler(bookService)
	borrowHandler := handlers.NewBorrowHandler(borrowService)
	favHandler := handlers.NewFavoriteHandler(favoriteRepo)
	notifHandler := handlers.NewNotificationHandler(notifService)
	adminHandler := handlers.NewAdminHandler(db, notifRepo)
	reviewHandler := handlers.NewReviewHandler(reviewRepo) // Инициализируем хендлер отзывов

	// 9. Gin Setup
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.ErrorHandler())
	router.Use(CORSMiddleware())

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now().Unix()})
	})

	// =============================================================================
	// PUBLIC ROUTES
	// =============================================================================
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// Middleware Definition
	authM := middleware.AuthMiddleware(cfg.JWTSecret)
	adminM := middleware.AdminOnly()

	// =============================================================================
	// PROTECTED MOBILE APPLICATION ROUTES
	// =============================================================================
	bookGroup := router.Group("/books")
	bookGroup.Use(authM)
	{
		bookGroup.GET("", bookHandler.GetAll)
		bookGroup.GET("/:id", bookHandler.GetByID)
		bookGroup.GET("/search", bookHandler.Search)
		bookGroup.POST("", bookHandler.Create)

		// Добавляем ручки для работы с отзывами к книгам
		bookGroup.POST("/:id/reviews", reviewHandler.CreateReview) // Оставить отзыв
		bookGroup.GET("/:id/reviews", reviewHandler.GetBookReviews) // Посмотреть отзывы
	}

	userGroup := router.Group("/user")
	userGroup.Use(authM)
	{
		userGroup.GET("/profile", userHandler.GetProfile)
		userGroup.PUT("/profile", userHandler.UpdateProfile)
		userGroup.PUT("/change-password", userHandler.ChangePassword)
		userGroup.GET("/history", userHandler.GetHistory)
	}

	borrowGroup := router.Group("/borrow")
	borrowGroup.Use(authM)
	{
		borrowGroup.POST("/:id", borrowHandler.BorrowBook)
		borrowGroup.POST("/:id/return", borrowHandler.ReturnBook)
	}

	favGroup := router.Group("/favorites")
	favGroup.Use(authM)
	{
		favGroup.POST("/:book_id", favHandler.Add)
		favGroup.DELETE("/:book_id", favHandler.Remove)
		favGroup.GET("", favHandler.List)
	}

	notifGroup := router.Group("/notifications")
	notifGroup.Use(authM)
	{
		notifGroup.GET("", notifHandler.List)
		notifGroup.PUT("/:id/read", notifHandler.MarkRead)
	}

	// =============================================================================
	// ADMIN ROUTES (WEB DASHBOARD)
	// =============================================================================
	adminGroup := router.Group("/admin")
	adminGroup.Use(authM, adminM)
	{
		// 1. Stats & Dashboard
		adminGroup.GET("/stats", adminHandler.GetStats)
		adminGroup.GET("/chart", adminHandler.GetChartData)
		adminGroup.GET("/overdue", adminHandler.GetOverdueMembers)

		// 2. Books Management
		booksGroup := adminGroup.Group("/books")
		{
			booksGroup.GET("", bookHandler.GetAll)
			booksGroup.POST("", bookHandler.Create)
			booksGroup.PUT("/:id", bookHandler.Update)
			booksGroup.DELETE("/:id", bookHandler.Delete)
		}

		// 3. Filters Management
		filtersGroup := adminGroup.Group("/filters")
		{
			filtersGroup.GET("", adminHandler.GetFilters)
			filtersGroup.POST("", adminHandler.AddFilter)
			filtersGroup.PUT("/:id", adminHandler.UpdateFilter)
			filtersGroup.DELETE("/:id", adminHandler.DeleteFilter)
		}

		// 4. Members & User Bookshelf
		adminGroup.GET("/members", adminHandler.GetAllMembers)
		adminGroup.POST("/members/:id/accept", adminHandler.AcceptAdminRequest)
		adminGroup.DELETE("/members/:id", adminHandler.DeleteMember)
		adminGroup.GET("/members/:user_id/books", adminHandler.GetUserBookshelf)

		// 5. Receive Book Requests
		adminGroup.GET("/receive-requests", adminHandler.GetReceiveRequests)
		adminGroup.POST("/receive-requests/:id/accept", adminHandler.AcceptReceiveRequest)
		adminGroup.DELETE("/receive-requests/:id", adminHandler.DeleteReceiveRequest)

		// 6. Return Book Requests
		adminGroup.GET("/return-requests", adminHandler.GetReturnRequests)
		adminGroup.POST("/return-requests/:id/accept", adminHandler.AcceptReturnRequest)
		adminGroup.DELETE("/return-requests/:id", adminHandler.DeleteReturnRequest)

		// 7. Notifications
		adminGroup.GET("/notifications", adminHandler.GetNotifications)
		adminGroup.POST("/notifications", adminHandler.AddNotification)

		// 8. Profile & Admin Management
		adminGroup.POST("/profile/change-password", userHandler.ChangePassword)
		adminGroup.GET("/profile/admins", adminHandler.GetAdminManagement)
	}

	// 10. Start Server
	serverPort := cfg.ServerPort
	if envPort := os.Getenv("PORT"); envPort != "" {
		serverPort = envPort
	}

	srv := &http.Server{
		Addr:    ":" + serverPort,
		Handler: router,
	}

	go func() {
		logger.Info("Server starting", zap.String("port", serverPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")
	ctxShut, cancelShut := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShut()

	if err := srv.Shutdown(ctxShut); err != nil {
		logger.Error("Server forced shutdown", zap.Error(err))
	}

	db.Close()
	logger.Info("Server stopped. Bye!")
}

// CORSMiddleware
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT, DELETE")
		c.Header("Access-Control-Max-Age", "43200")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

// runMigrations
func runMigrations(db *repositories.DB) error {
	query := `
		ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(20) DEFAULT 'user';
		ALTER TABLE users ADD COLUMN IF NOT EXISTS job_position VARCHAR(50) DEFAULT 'Student';
		ALTER TABLE users ADD COLUMN IF NOT EXISTS is_pending_admin BOOLEAN DEFAULT FALSE;
		ALTER TABLE users ADD COLUMN IF NOT EXISTS date_of_birth VARCHAR(20);

		ALTER TABLE books ADD COLUMN IF NOT EXISTS status VARCHAR(20) DEFAULT 'Available';
		ALTER TABLE books ADD COLUMN IF NOT EXISTS page_count INTEGER DEFAULT 0;
		ALTER TABLE books ADD COLUMN IF NOT EXISTS language VARCHAR(50) DEFAULT 'English';
		ALTER TABLE books ADD COLUMN IF NOT EXISTS bg_image_url VARCHAR(500) DEFAULT '';

		-- Добавляем колонки для гибких уведомлений, как просит фронтенд
		ALTER TABLE notifications ADD COLUMN IF NOT EXISTS type VARCHAR(20) DEFAULT 'news';
		ALTER TABLE notifications ADD COLUMN IF NOT EXISTS title VARCHAR(255) DEFAULT '';
		ALTER TABLE notifications ADD COLUMN IF NOT EXISTS notification_image_url VARCHAR(500) DEFAULT '';

		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY, name VARCHAR(100) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL, phone VARCHAR(20) NOT NULL,
			password VARCHAR(255) NOT NULL, created_at TIMESTAMP DEFAULT NOW(),
			role VARCHAR(20) DEFAULT 'user', job_position VARCHAR(50) DEFAULT 'Student'
		);
		CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY, title VARCHAR(255) NOT NULL,
			author VARCHAR(255) NOT NULL, description TEXT,
			category VARCHAR(100), year INTEGER,
			available_copies INTEGER DEFAULT 0,
			image_url VARCHAR(500), created_at TIMESTAMP DEFAULT NOW(),
			status VARCHAR(20) DEFAULT 'Available',
			bg_image_url VARCHAR(500) DEFAULT ''
		);
		CREATE TABLE IF NOT EXISTS borrows (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
			borrow_date TIMESTAMP NOT NULL, due_date TIMESTAMP NOT NULL,
			return_date TIMESTAMP, status VARCHAR(20) DEFAULT 'active',
			created_at TIMESTAMP DEFAULT NOW()
		);
		CREATE TABLE IF NOT EXISTS favorites (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
			created_at TIMESTAMP DEFAULT NOW(),
			UNIQUE(user_id, book_id)
		);
		CREATE TABLE IF NOT EXISTS notifications (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			message TEXT NOT NULL, is_read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT NOW(),
			type VARCHAR(20) DEFAULT 'news',
			title VARCHAR(255) DEFAULT '',
			notification_image_url VARCHAR(500) DEFAULT ''
		);
		CREATE TABLE IF NOT EXISTS categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) UNIQUE NOT NULL
		);

		-- НОВАЯ ТАБЛИЦА ДЛЯ ОТЗЫВОВ И РЕЙТИНГА КНИГ
		CREATE TABLE IF NOT EXISTS reviews (
			id SERIAL PRIMARY KEY,
			book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
			review TEXT,
			review_category VARCHAR(100) DEFAULT '',
			created_at TIMESTAMP DEFAULT NOW()
		);
	`
	return db.Exec(query)
}