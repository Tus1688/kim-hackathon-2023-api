package models

import (
	"fmt"

	"github.com/Tus1688/kim-hackathon-2023-api/database"
)

type CreateProduct struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	BusinessId  string  `json:"business_id" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}

type ProductResponse struct {
	Id           string  `json:"id"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	BusinessId   string  `json:"business_id"`
	BusinessName string  `json:"business_name"`
	Price        float64 `json:"price"`
	UpdatedOn    string  `json:"updated_on"`
}

type ModifyProduct struct {
	Id          string  `json:"id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	BusinessId  string  `json:"business_id" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}

func (c *CreateProduct) Create() error {
	if c.Price <= 0 {
		return fmt.Errorf("invalid input")
	}
	_, err := database.MysqlInstance.Exec(
		"INSERT INTO products (name, description, business_refer, price) VALUES (?, ?, UUID_TO_BIN(?), ?)", c.Name,
		c.Description, c.BusinessId, c.Price,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetProduct(searchQuery string) ([]ProductResponse, error) {
	query := `
		SELECT BIN_TO_UUID(p.id), p.name, p.description, BIN_TO_UUID(p.business_refer), b.name, price, p.updated_at
		FROM products p
		INNER JOIN businesses b ON b.id = p.business_refer
	`
	var args []interface{}

	if searchQuery != "" {
		query += " WHERE p.name LIKE ?"
		args = append(args, "%"+searchQuery+"%")
	}

	rows, err := database.MysqlInstance.Query(query, args...)
	if err != nil {
		return nil, err
	}

	var res []ProductResponse
	for rows.Next() {
		var temp ProductResponse
		err := rows.Scan(
			&temp.Id, &temp.Name, &temp.Description, &temp.BusinessId, &temp.BusinessName, &temp.Price, &temp.UpdatedOn,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, temp)
	}
	return res, nil
}

func (m *ModifyProduct) Modify() error {
	res, err := database.MysqlInstance.Exec(
		"UPDATE products SET name = ?, description = ?, business_refer = UUID_TO_BIN(?), price = ?, updated_at = NOW() WHERE id = UUID_TO_BIN(?)",
		m.Name, m.Description, m.BusinessId, m.Price, m.Id,
	)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func DeleteProduct(id string) error {
	res, err := database.MysqlInstance.Exec(
		"DELETE FROM products WHERE id = UUID_TO_BIN(?)", id,
	)
	if err != nil {
		return err
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}
