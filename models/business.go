package models

import (
	"fmt"

	"github.com/Tus1688/kim-hackathon-2023-api/database"
)

type CreateBusiness struct {
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type BusinessResponse struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	UpdatedOn   string `json:"updated_on"`
}

type ModifyBusiness struct {
	Id          string `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

func (c *CreateBusiness) Create() error {
	if len(c.PhoneNumber) > 15 || len(c.Name) > 32 || len(c.Address) > 255 {
		return fmt.Errorf("invalid input")
	}
	_, err := database.MysqlInstance.Exec(
		`INSERT INTO businesses (name, address, phone_number) VALUES (?, ?, ?)`, c.Name, c.Address, c.PhoneNumber,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetBusinesses(searchQuery string) ([]BusinessResponse, error) {
	var businesses []BusinessResponse
	query := "SELECT BIN_TO_UUID(id), name, address, phone_number, updated_at FROM businesses"
	var args []interface{}
	if searchQuery != "" {
		query += " WHERE name LIKE ?"
		args = append(args, "%"+searchQuery+"%")
	}
	rows, err := database.MysqlInstance.Query(query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var business BusinessResponse
		err = rows.Scan(&business.Id, &business.Name, &business.Address, &business.PhoneNumber, &business.UpdatedOn)
		if err != nil {
			return nil, err
		}
		businesses = append(businesses, business)
	}
	return businesses, nil
}

func (m *ModifyBusiness) ModifyBusiness() error {
	if len(m.PhoneNumber) > 15 || len(m.Name) > 32 || len(m.Address) > 255 {
		return fmt.Errorf("invalid input")
	}
	res, err := database.MysqlInstance.Exec(
		`UPDATE businesses SET name = ?, address = ?, phone_number = ?, updated_at = NOW() WHERE id = UUID_TO_BIN(?)`,
		m.Name, m.Address,
		m.PhoneNumber, m.Id,
	)
	if err != nil {
		return err
	}

	if affected, err := res.RowsAffected(); err != nil || affected == 0 {
		return fmt.Errorf("business not found")
	}

	return nil
}

func DeleteBusiness(id string) error {
	res, err := database.MysqlInstance.Exec(
		`DELETE FROM businesses WHERE id = UUID_TO_BIN(?)`, id,
	)
	if err != nil {
		return err
	}
	if affected, err := res.RowsAffected(); err != nil || affected == 0 {
		return fmt.Errorf("business not found")
	}
	return nil
}
