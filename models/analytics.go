package models

import "github.com/Tus1688/kim-hackathon-2023-api/database"

type TotalLending struct {
	Amount       float64 `json:"amount"`
	InterestRate float64 `json:"interest_rate"`
}

type TotalHelper struct {
	Total int `json:"total"`
}

type TopProduct struct {
	ProductId   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Amount      float64 `json:"amount"`
}

func GetTotalUser() (TotalHelper, error) {
	var total TotalHelper
	query := `SELECT COUNT(*) FROM users`
	err := database.MysqlInstance.QueryRow(query).Scan(&total.Total)
	if err != nil {
		return total, err
	}

	return total, nil
}

func GetTotalSME() (TotalHelper, error) {
	var total TotalHelper
	query := `SELECT COUNT(*) FROM businesses`
	err := database.MysqlInstance.QueryRow(query).Scan(&total.Total)
	if err != nil {
		return total, err
	}
	return total, nil
}

func GetTotalAwaitingApproval() (TotalHelper, error) {
	var total TotalHelper
	query := `SELECT COUNT(*) FROM kim.lending WHERE is_approved = false AND is_rejected = false`
	err := database.MysqlInstance.QueryRow(query).Scan(&total.Total)
	if err != nil {
		return total, err
	}
	return total, nil
}

func GetTopProduct() ([]TopProduct, error) {
	var top []TopProduct
	query := `
	SELECT MAX(BIN_TO_UUID(o.product_refer)) AS product_id, p.name, SUM(p.price * o.quantity * o.commission_rate) AS amount
	FROM orders o
	LEFT JOIN products p ON p.id = o.product_refer
	GROUP BY p.id
	ORDER BY amount DESC
	LIMIT 5;
	`
	rows, err := database.MysqlInstance.Query(query)
	if err != nil {
		return top, err
	}
	for rows.Next() {
		var t TopProduct
		err = rows.Scan(&t.ProductId, &t.ProductName, &t.Amount)
		if err != nil {
			return top, err
		}
		top = append(top, t)
	}
	return top, nil
}
