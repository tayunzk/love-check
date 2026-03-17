package model

import "time"

type Item struct {
	ID        int        `json:"id"`
	Content   string     `json:"content"`
	Completed bool       `json:"completed"`
	ItemDate  *time.Time `json:"item_date"`
	CreatedAt time.Time  `json:"created_at"`
}

type CreateItemRequest struct {
	Content string     `json:"content" binding:"required"`
	ItemDate *string   `json:"item_date"`
}

type UpdateItemRequest struct {
	Completed bool   `json:"completed"`
	ItemDate  *string `json:"item_date"`
}