package product

import (
	"database/sql"

	"github.com/phildehovre/go-complete-api/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetProducts() ([]types.Product, error) {
	rows, err := s.db.Query("SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	products := make([]types.Product, 0)
	for rows.Next() {
		p, err := scanRowsIntoProducts(rows)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}
	return products, nil
}

func scanRowsIntoProducts(rows *sql.Rows) (*types.Product, error) {
	p := new(types.Product)
	err := rows.Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Quantity,
		&p.Image,
		&p.CreatedAt,
		&p.Price,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Store) CreateProduct(p types.Product) error {
	_, err := s.db.Exec(`INSERT INTO products (name, description, quantity, image, price)
VALUES (?,?,?,?,?)
`, p.Name, p.Description, p.Quantity, p.Image, p.Price)
	if err != nil {
		return err
	}
	return nil
}
