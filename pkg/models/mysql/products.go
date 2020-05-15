package mysql

import (
	"database/sql"
	"errors"

	"github.com/ahmedavid/gocommerce/pkg/models"
)

type ProductModel struct {
	DB *sql.DB
}

func (m *ProductModel) DeleteProduct(id int) (int, error) {
	stmt := `DELETE FROM products WHERE id = ?`
	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		return 0, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rows), nil
}

func (m *ProductModel) CreateProduct(cat_id int, name, description, imgURL string, price float64, stock int) (int, error) {
	stmt := `INSERT INTO products (cat_id,name,description,img_url,price,stock,created) VALUES(?, ?, ?, ?, ?, ?, UTC_TIMESTAMP())`
	result, err := m.DB.Exec(stmt, cat_id, name, description, imgURL, price, stock)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *ProductModel) GetByID(product_id int) (*models.Product, error) {
	stmt := `SELECT * from products WHERE id = ?`
	row := m.DB.QueryRow(stmt, product_id)
	p := &models.Product{}
	err := row.Scan(&p.ID, &p.CatID, &p.Name, &p.Price, &p.Stock, &p.Created, &p.ImgURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return p, nil
}

func (m *ProductModel) GetByCategory(cat_id int) ([]models.Product, error) {
	var stmt string
	var rows *sql.Rows
	var err error
	if cat_id == 0 {
		stmt = `SELECT p.id,p.cat_id,p.name,p.price,p.stock,p.img_url,p.created FROM products p INNER JOIN categories c ON p.cat_id = c.id`
		rows, err = m.DB.Query(stmt)
	} else {
		stmt = `SELECT p.id,p.cat_id,p.name,p.price,p.stock,p.img_url,p.created FROM products p INNER JOIN categories c ON p.cat_id = c.id WHERE p.cat_id = ?`
		rows, err = m.DB.Query(stmt, cat_id)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	products := []models.Product{}
	for rows.Next() {
		p := &models.Product{}
		err := rows.Scan(&p.ID, &p.CatID, &p.Name, &p.Price, &p.Stock, &p.ImgURL, &p.Created)
		if err != nil {
			return nil, err
		}
		products = append(products, *p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
