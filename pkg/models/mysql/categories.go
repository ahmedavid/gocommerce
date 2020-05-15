package mysql

import (
	"database/sql"

	"github.com/ahmedavid/gocommerce/pkg/models"
)

type CategoryModel struct {
	DB *sql.DB
}

func (m *CategoryModel) GetAll() ([]models.Category, error) {
	stmt := `SELECT * FROM categories`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	categories := []models.Category{}
	for rows.Next() {
		c := &models.Category{}
		err := rows.Scan(&c.ID, &c.Name, &c.Created)
		if err != nil {
			return nil, err
		}
		categories = append(categories, *c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
