package handler

import (
	"net/http"
	"strconv"
	"time"

	"love-check/model"
	"love-check/repository"

	"github.com/gin-gonic/gin"
)

type ItemHandler struct {
	repo *repository.ItemRepository
}

func NewItemHandler(repo *repository.ItemRepository) *ItemHandler {
	return &ItemHandler{repo: repo}
}

func (h *ItemHandler) GetItems(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "无效日期格式"})
			return
		}
		items, err := h.repo.GetByDate(date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
			return
		}
		c.JSON(http.StatusOK, items)
		return
	}

	items, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *ItemHandler) AddItem(c *gin.Context) {
	var req model.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "内容不能为空"})
		return
	}

	var itemDate *time.Time
	if req.ItemDate != nil {
		t, err := time.Parse("2006-01-02", *req.ItemDate)
		if err == nil {
			itemDate = &t
		}
	}

	item, err := h.repo.Create(req.Content, itemDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加失败"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) ToggleItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}

	var req model.UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效请求"})
		return
	}

	var itemDate *time.Time
	if req.ItemDate != nil {
		t, err := time.Parse("2006-01-02", *req.ItemDate)
		if err == nil {
			itemDate = &t
		}
	}

	item, err := h.repo.Update(id, req.Completed, itemDate)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) DeleteItem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效ID"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}