package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/Tus1688/kim-hackathon-2023-api/database"
	"golang.org/x/crypto/bcrypt"
)

type RegisterAsBorrower struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LendingRequest struct {
	RequesterUid     string  // we get this from the context
	Amount           float64 `json:"amount" binding:"required"`
	InterestRate     int     `json:"interest_rate" binding:"required"`
	Tenor            int     `json:"tenor" binding:"required"`
	Age              int     `json:"age" binding:"required"`
	Gender           bool    `json:"gender" binding:"required"`
	Income           float64 `json:"income" binding:"required"`
	LastEducation    int     `json:"last_education" binding:"required"`
	MaritalStatus    bool    `json:"marital_status" binding:"required"`
	NumberOfChildren int     `json:"number_of_children" binding:"required"`
	HasHouse         bool    `json:"has_house" binding:"required"`
	KkUrl            string  `json:"kk_url" binding:"required"`
	KtpUrl           string  `json:"ktp_url" binding:"required"`
}

type LendingResponse struct {
	Id               string  `json:"id"`
	Amount           float64 `json:"amount"`
	InterestRate     int     `json:"interest_rate"`
	Tenor            int     `json:"tenor"`
	Age              int     `json:"age"`
	Gender           string  `json:"gender"`
	Income           float64 `json:"income"`
	LastEducation    string  `json:"last_education"`
	MaritalStatus    string  `json:"marital_status"`
	NumberOfChildren int     `json:"number_of_children"`
	HasHouse         string  `json:"has_house"`
	KkUrl            string  `json:"kk_url"`
	KtpUrl           string  `json:"ktp_url"`
	Status           string  `json:"status"`
	PaymentToken     string  `json:"payment_token,omitempty"`
	PaymentUrl       string  `json:"payment_url,omitempty"`
}

func (r *RegisterAsBorrower) Register() error {
	passBytes, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = database.MysqlInstance.Exec(
		`INSERT INTO users (username, hashed_password, is_user) VALUES (?, ?, ?)`, r.Username, string(passBytes), true,
	)
	if err != nil {
		return err
	}
	return nil
}

func UploadDocument(image multipart.File, header *multipart.FileHeader) (GoBlobResponse, error) {
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
	return res, nil
}

func (l *LendingRequest) Create() error {
	_, err := database.MysqlInstance.Exec(
		`INSERT INTO lending(user_refer, amount, interest_rate, tenor, age, income, last_education, number_of_children, kk_url, ktp_url, status) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		l.RequesterUid, l.Amount, l.InterestRate, l.Tenor, l.Age, l.Income, l.LastEducation, l.NumberOfChildren,
		l.KkUrl, l.KtpUrl, "pending",
	)
	if err != nil {
		return err
	}
	return nil
}

func GetLendingAsUser(uid string) ([]LendingResponse, error) {
	rows, err := database.MysqlInstance.Query(
		`
		SELECT BIN_TO_UUID(l.id),
			   l.amount,
			   l.interest_rate,
			   l.tenor,
			   l.age,
			   CASE l.gender
				   WHEN 0 THEN 'Laki-Laki'
				   WHEN 1 THEN 'Perempuan'
				   END AS GENDER,
			   l.income,
			   CASE l.last_education
				   WHEN 0 THEN 'SMA'
				   WHEN 1 THEN 'D3'
				   WHEN 2 THEN 'S1'
				   WHEN 3 THEN 'S2'
				   WHEN 4 THEN 'S3'
				   END AS last_education,
			   CASE l.marital_status
				   WHEN 0 THEN 'Lajang'
				   WHEN 1 THEN 'Menikah'
				   END AS marital_status,
			   l.number_of_children,
			   CASE l.home_ownership
				   WHEN 0 THEN 'Menyewa'
				   WHEN 1 THEN 'Memiliki'
				   END AS home_ownership,
				l.kk_url, l.ktp_url,
		        l.status, COALESCE(l.payment_token, ''), COALESCE(l.payment_url, '')
		FROM lending l
		WHERE l.user_refer = UUID_TO_BIN(?)
	`, uid,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var res []LendingResponse
	for rows.Next() {
		var temp LendingResponse
		err := rows.Scan(
			&temp.Id, &temp.Amount, &temp.InterestRate, &temp.Tenor, &temp.Age, &temp.Gender, &temp.Income,
			&temp.LastEducation, &temp.MaritalStatus, &temp.NumberOfChildren, &temp.HasHouse, &temp.KkUrl, &temp.KtpUrl,
			&temp.Status,
			&temp.PaymentToken, &temp.PaymentUrl,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, temp)
	}
	return res, nil
}
