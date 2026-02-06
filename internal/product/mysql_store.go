package product

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
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

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			price DECIMAL(10,2) NOT NULL
		)	
	`)

	if err != nil {
		return nil, err
	}

	return &MySQLStore{db: db}, nil
}

func (s *MySQLStore) List() []Product {
	rows, err := s.db.Query("SELECT id, name, price FROM products")
	if err != nil {
		return []Product{}
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

	return products
}

func (s *MySQLStore) GetByID(id int) (Product, error) {
	var p Product
	err := s.db.QueryRow("SELECT id, name, price FROM products WHERE id = ?", id).
		Scan(&p.ID, &p.Name, &p.Price)

	if err == sql.ErrNoRows {
		return Product{}, fmt.Errorf("product %d not found", id)
	}
	if err != nil {
		return Product{}, nil
	}

	return p, nil
}

func (s *MySQLStore) Create(p Product) Product {
	result, err := s.db.Exec(
		"INSERT INTO products (name, price) VALUES (?,?)",
		p.Name, p.Price,
	)

	if err != nil {
		return Product{}
	}

	id, _ := result.LastInsertId()
	p.ID = int(id)
	return p 
}

func (s *MySQLStore) Update(id int, p Product) (Product, error) {
	result, err := s.db.Exec(
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

func (s *MySQLStore) Delete(id int) error {
	result, err := s.db.Exec("DELETE FROM products WHERE id = ?", id)
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