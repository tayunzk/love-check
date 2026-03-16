package model

import "time"

type Item struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateItemRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateItemRequest struct {
	Completed bool `json:"completed"`
}