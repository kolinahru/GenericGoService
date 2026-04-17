package repository

import (
	"database/sql"
	"errors"

	"go-day3/models"
)

type ItemRepository interface {
	GetAll() ([]models.Item, error)
	GetByID(id int) (models.Item, error)
	Create(name string) (models.Item, error)
	Update(id int, name string) (models.Item, error)
	Delete(id int) error
}

type PostgresItemRepository struct {
	db *sql.DB
}

func NewPostgresItemRepository(db *sql.DB) *PostgresItemRepository {
	return &PostgresItemRepository{db: db}
}

func (r *PostgresItemRepository) GetAll() ([]models.Item, error) {
	rows, err := r.db.Query("SELECT id, name FROM items ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.Item, 0)

	for rows.Next() {
		var item models.Item
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *PostgresItemRepository) GetByID(id int) (models.Item, error) {
	var item models.Item

	err := r.db.QueryRow(
		"SELECT id, name FROM items WHERE id = $1",
		id,
	).Scan(&item.ID, &item.Name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Item{}, sql.ErrNoRows
		}
		return models.Item{}, err
	}

	return item, nil
}

func (r *PostgresItemRepository) Create(name string) (models.Item, error) {
	var item models.Item

	err := r.db.QueryRow(
		"INSERT INTO items (name) VALUES ($1) RETURNING id, name",
		name,
	).Scan(&item.ID, &item.Name)

	if err != nil {
		return models.Item{}, err
	}

	return item, nil
}

func (r *PostgresItemRepository) Update(id int, name string) (models.Item, error) {
	var item models.Item

	err := r.db.QueryRow(
		"UPDATE items SET name = $1 WHERE id = $2 RETURNING id, name",
		name,
		id,
	).Scan(&item.ID, &item.Name)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Item{}, sql.ErrNoRows
		}
		return models.Item{}, err
	}

	return item, nil
}

func (r *PostgresItemRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM items WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
