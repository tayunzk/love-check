package main

import (
	"database/sql"
	"fmt"
	"log"

	"love-check/config"
	"love-check/handler"
	"love-check/middleware"
	"love-check/repository"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load()

	db := initDB(cfg)
	defer db.Close()

	repo := repository.NewItemRepository(db)
	itemHandler := handler.NewItemHandler(repo)

	r := gin.Default()
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		api.GET("/items", itemHandler.GetItems)
		api.POST("/items", itemHandler.AddItem)
		api.PUT("/items/:id", itemHandler.ToggleItem)
		api.DELETE("/items/:id", itemHandler.DeleteItem)
	}

	r.Static("/static", "/app/frontend")
	r.StaticFile("/", "/app/frontend/index.html")

	log.Printf("服务启动: http://localhost:%s", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}

func initDB(cfg *config.Config) *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Ping数据库失败: %v", err)
	}

	db.Exec(`
		CREATE TABLE IF NOT EXISTS items (
			id INT AUTO_INCREMENT PRIMARY KEY,
			content VARCHAR(255) NOT NULL,
			completed BOOLEAN DEFAULT FALSE,
			item_date DATE DEFAULT (DATE(created_at)),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)

	var colCount int
	db.QueryRow("SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'items' AND COLUMN_NAME = 'item_date'").Scan(&colCount)
	if colCount == 0 {
		db.Exec("ALTER TABLE items ADD COLUMN item_date DATE DEFAULT (DATE(created_at))")
	}

	log.Println("数据库初始化成功")
	return db
}