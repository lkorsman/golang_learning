package product

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"lukekorsman.com/store/internal/database"
)

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore(connStr string) (*MySQLStore, error) {
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

    if err := database.RunMigrations(db); err != nil {
        return nil, fmt.Errorf("migration failed: %w", err)
    }

	return &MySQLStore{db: db}, nil
}

func (s *MySQLStore) List(ctx context.Context) ([]Product, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT id, name, price FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			continue 
		}
		products = append(products, p)
	}

	return products, nil
}

func (s *MySQLStore) GetByID(ctx context.Context, id int) (Product, error) {
	var p Product
	err := s.db.QueryRowContext(ctx, "SELECT id, name, price FROM products WHERE id = ?", id).
		Scan(&p.ID, &p.Name, &p.Price)

	if err == sql.ErrNoRows {
		return Product{}, fmt.Errorf("product %d not found", id)
	}
	if err != nil {
		return Product{}, nil
	}

	return p, nil
}

func (s *MySQLStore) Create(ctx context.Context, p Product) (Product, error) {
	result, err := s.db.ExecContext( ctx,
		"INSERT INTO products (name, price) VALUES (?,?)",
		p.Name, p.Price,
	)

	if err != nil {
		return Product{}, err
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)
	return p, nil
}

func (s *MySQLStore) Update(ctx context.Context, id int, p Product) (Product, error) {
	result, err := s.db.ExecContext(ctx, 
		"UPDATE products SET name = ?, price = ?, WHERE id = ?",
		p.Name, p.Price, id, 
	)

	if err != nil {
		return Product{}, err 
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return Product{}, fmt.Errorf("product %d not found", id)
	}

	p.ID = id
	return p, nil
}

func (s *MySQLStore) Delete(ctx context.Context, id int) error {
	result, err := s.db.ExecContext(ctx,
		"DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err 
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("product %d not found", id)
	}

	return nil
}

func (s *MySQLStore) Close() error {
    return s.db.Close()
}