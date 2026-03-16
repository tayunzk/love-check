package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Item struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	items  []Item
	mu     sync.RWMutex
	nextID int = 1
)

const dataFile = "/app/data/items.json"

func loadItems() {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			items = []Item{}
			return
		}
		log.Printf("加载数据失败: %v", err)
		return
	}
	json.Unmarshal(data, &items)
	if len(items) > 0 {
		maxID := 0
		for _, item := range items {
			if item.ID > maxID {
				maxID = item.ID
			}
		}
		nextID = maxID + 1
	}
}

func saveItems() {
	data, err := json.Marshal(items)
	if err != nil {
		log.Printf("保存数据失败: %v", err)
		return
	}
	os.WriteFile(dataFile, data, 0644)
}

func getItems(c *gin.Context) {
	mu.RLock()
	defer mu.RUnlock()
	c.JSON(http.StatusOK, items)
}

func addItem(c *gin.Context) {
	var newItem struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&newItem); err != nil || newItem.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "内容不能为空"})
		return
	}

	mu.Lock()
	item := Item{
		ID:        nextID,
		Content:   newItem.Content,
		Completed: false,
		CreatedAt: time.Now(),
	}
	items = append(items, item)
	nextID++
	mu.Unlock()

	saveItems()
	c.JSON(http.StatusOK, item)
}

func toggleItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}
	var update struct {
		Completed bool `json:"completed"`
	}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效请求"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i := range items {
		if items[i].ID == id {
			items[i].Completed = update.Completed
			saveItems()
			c.JSON(http.StatusOK, items[i])
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
}

func deleteItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, item := range items {
		if item.ID == id {
			items = append(items[:i], items[i+1:]...)
			saveItems()
			c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
}

func main() {
	// 确保数据目录存在
	os.MkdirAll("/app/data", 0644)
	loadItems()

	r := gin.Default()

	// CORS 中间件
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// API 路由
	r.GET("/api/items", getItems)
	r.POST("/api/items", addItem)
	r.PUT("/api/items/:id", toggleItem)
	r.DELETE("/api/items/:id", deleteItem)

	// 前端静态文件
	r.Static("/static", "/app/frontend")
	r.StaticFile("/", "/app/frontend/index.html")

	log.Println("服务启动: http://localhost:8080")
	r.Run(":8080")
}