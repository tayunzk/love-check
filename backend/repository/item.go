package repository

import (
	"database/sql"
	"love-check/model"

	_ "github.com/go-sql-driver/mysql"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) GetAll() ([]model.Item, error) {
	rows, err := r.db.Query("SELECT id, content, completed, created_at FROM items ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		if err := rows.Scan(&item.ID, &item.Content, &item.Completed, &item.CreatedAt); err != nil {
			continue
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ItemRepository) Create(content string) (*model.Item, error) {
	result, err := r.db.Exec("INSERT INTO items (content, completed) VALUES (?, ?)", content, false)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return &model.Item{
		ID:        int(id),
		Content:   content,
		Completed: false,
	}, nil
}

func (r *ItemRepository) Update(id int, completed bool) (*model.Item, error) {
	_, err := r.db.Exec("UPDATE items SET completed = ? WHERE id = ?", completed, id)
	if err != nil {
		return nil, err
	}

	var item model.Item
	err = r.db.QueryRow("SELECT id, content, completed, created_at FROM items WHERE id = ?", id).
		Scan(&item.ID, &item.Content, &item.Completed, &item.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM items WHERE id = ?", id)
	return err
}