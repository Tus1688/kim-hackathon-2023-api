package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

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

type GoBlobResponse struct {
	Filename string `json:"filename"`
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

func CreateProductImage(productId string, image multipart.File, header *multipart.FileHeader) (GoBlobResponse, error) {
	// check if product exists
	_, err := database.MysqlInstance.Exec(
		"SELECT id FROM products WHERE id = UUID_TO_BIN(?)", productId,
	)
	if err != nil {
		return GoBlobResponse{}, err
	}

	//	Create a POST request to the image server
	url := goBlobBaseUrl + "/file"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return GoBlobResponse{}, err
	}

	// Set the content type to multipart/form-data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		return GoBlobResponse{}, err
	}
	_, err = io.Copy(part, image)
	if err != nil {
		return GoBlobResponse{}, err
	}
	err = writer.Close()
	if err != nil {
		return GoBlobResponse{}, err
	}

	//	pass file name
	req.Header.Set("File-Name", header.Filename)
	req.Header.Set("Authorization", goBlobAuthorization)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Body = io.NopCloser(body)

	//	Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GoBlobResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return GoBlobResponse{}, fmt.Errorf("image server error")
	}

	//	Extract filename
	//  {filename: target.jpg}
	var res GoBlobResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return GoBlobResponse{}, err
	}

	//	Insert into database
	_, err = database.MysqlInstance.Exec(
		"INSERT INTO product_images (product_refer, file_name) VALUES (UUID_TO_BIN(?), ?)", productId, res.Filename,
	)
	if err != nil {
		return GoBlobResponse{}, err
	}
	return res, nil
}
