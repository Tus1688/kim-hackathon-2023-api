package models

import "github.com/Tus1688/kim-hackathon-2023-api/database"

type TotalLending struct {
	Amount       float64 `json:"amount"`
	InterestRate float64 `json:"interest_rate"`
}

type TotalHelper struct {
	Total int `json:"total"`
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
