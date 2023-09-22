package models

import (
	"fmt"

	"github.com/Tus1688/kim-hackathon-2023-api/database"
)

type CreateOrder struct {
	ProductId  string `json:"product_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required"`
	Commission int    `json:"commission" binding:"required"`
}

type OrderResponse struct {
	Id          string  `json:"id"`
	ProductId   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Commission  int     `json:"commission"`
	Profit      float64 `json:"profit"`
	UpdatedOn   string  `json:"updated_on"`
}

type ModifyOrder struct {
	Id         string `json:"id" binding:"required"`
	ProductId  string `json:"product_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required"`
	Commission int    `json:"commission" binding:"required"`
}

func (c *CreateOrder) Create() error {
	_, err := database.MysqlInstance.Exec(
		"INSERT INTO orders (product_refer, quantity, commission_rate) VALUES (UUID_TO_BIN(?), ?, ?)", c.ProductId,
		c.Quantity, c.Commission,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetOrder() ([]OrderResponse, error) {
	rows, err := database.MysqlInstance.Query(
		`
	SELECT BIN_TO_UUID(o.id), BIN_TO_UUID(o.product_refer), p.name, o.quantity, o.commission_rate, (p.price * o.quantity * o.commission_rate / 100), o.updated_at
	FROM orders o
	INNER JOIN products p ON p.id = o.product_refer
	`,
	)
	if err != nil {
		return nil, err
	}

	var res []OrderResponse
	for rows.Next() {
		var order OrderResponse
		err = rows.Scan(
			&order.Id, &order.ProductId, &order.ProductName, &order.Quantity, &order.Commission, &order.Profit,
			&order.UpdatedOn,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, order)
	}
	return res, nil
}

func (m *ModifyOrder) Modify() error {
	res, err := database.MysqlInstance.Exec(
		"UPDATE orders SET product_refer = UUID_TO_BIN(?), quantity = ?, commission_rate = ? WHERE id = UUID_TO_BIN(?)",
		m.ProductId, m.Quantity, m.Commission, m.Id,
	)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("order not found")
	}
	return nil
}
