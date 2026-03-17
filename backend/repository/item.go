package repository

import (
	"database/sql"
	"love-check/model"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type ItemRepository struct {
	db *sql.DB
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) GetAll() ([]model.Item, error) {
	rows, err := r.db.Query("SELECT id, content, completed, item_date, created_at FROM items ORDER BY item_date DESC, created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		var itemDate sql.NullTime
		if err := rows.Scan(&item.ID, &item.Content, &item.Completed, &itemDate, &item.CreatedAt); err != nil {
			continue
		}
		if itemDate.Valid {
			item.ItemDate = &itemDate.Time
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ItemRepository) GetByDate(date time.Time) ([]model.Item, error) {
	rows, err := r.db.Query("SELECT id, content, completed, item_date, created_at FROM items WHERE DATE(item_date) = ? ORDER BY created_at DESC", date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		var itemDate sql.NullTime
		if err := rows.Scan(&item.ID, &item.Content, &item.Completed, &itemDate, &item.CreatedAt); err != nil {
			continue
		}
		if itemDate.Valid {
			item.ItemDate = &itemDate.Time
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ItemRepository) Create(content string, itemDate *time.Time) (*model.Item, error) {
	var result sql.Result
	var err error

	if itemDate != nil {
		result, err = r.db.Exec("INSERT INTO items (content, completed, item_date) VALUES (?, ?, ?)", content, false, itemDate)
	} else {
		result, err = r.db.Exec("INSERT INTO items (content, completed) VALUES (?, ?)", content, false)
	}
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()

	var item model.Item
	var date sql.NullTime
	err = r.db.QueryRow("SELECT id, content, completed, item_date, created_at FROM items WHERE id = ?", id).
		Scan(&item.ID, &item.Content, &item.Completed, &date, &item.CreatedAt)
	if err != nil {
		return nil, err
	}
	if date.Valid {
		item.ItemDate = &date.Time
	}
	return &item, nil
}

func (r *ItemRepository) Update(id int, completed bool, itemDate *time.Time) (*model.Item, error) {
	var err error
	if itemDate != nil {
		_, err = r.db.Exec("UPDATE items SET completed = ?, item_date = ? WHERE id = ?", completed, itemDate, id)
	} else {
		_, err = r.db.Exec("UPDATE items SET completed = ? WHERE id = ?", completed, id)
	}
	if err != nil {
		return nil, err
	}

	var item model.Item
	var date sql.NullTime
	err = r.db.QueryRow("SELECT id, content, completed, item_date, created_at FROM items WHERE id = ?", id).
		Scan(&item.ID, &item.Content, &item.Completed, &date, &item.CreatedAt)
	if err != nil {
		return nil, err
	}
	if date.Valid {
		item.ItemDate = &date.Time
	}
	return &item, nil
}

func (r *ItemRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM items WHERE id = ?", id)
	return err
}