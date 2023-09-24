package models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/Tus1688/kim-hackathon-2023-api/database"
	"github.com/Tus1688/kim-hackathon-2023-api/midtrans"
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
	IsPaid           bool    `json:"is_paid"`
}

type LendingAdminResponse struct {
	Id               string  `json:"id"`
	UserId           string  `json:"user_id"`
	Username         string  `json:"username"`
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
	IsApproved       bool    `json:"is_approved"`
	IsRejected       bool    `json:"is_rejected"`
}

type LendingPredictRequest struct {
	Age              int `json:"Age"`
	Gender           int `json:"Gender"`
	Income           int `json:"Income"`
	Education        int `json:"Education"`
	MaritalStatus    int `json:"Marital_Status"`
	NumberOfChildren int `json:"Number_of_Children"`
	HomeOwnership    int `json:"Home_Ownership"`
}

type LendingPredictResponse struct {
	Predictions []string `json:"predictions"`
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
		        l.status, COALESCE(l.payment_token, ''), COALESCE(l.payment_url, ''), l.is_paid
		FROM lending l
		WHERE l.user_refer = UUID_TO_BIN(?)
		ORDER BY l.created_at DESC
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
			&temp.PaymentToken, &temp.PaymentUrl, &temp.IsPaid,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, temp)
	}
	return res, nil
}

func GetLendingAsAdmin() ([]LendingAdminResponse, error) {
	rows, err := database.MysqlInstance.Query(
		`
		SELECT BIN_TO_UUID(l.id),
		       BIN_TO_UUID(l.user_refer),
		       u.username,
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
		        l.status, COALESCE(l.payment_token, ''), COALESCE(l.payment_url, ''), is_approved, is_rejected
		FROM lending l
		INNER JOIN users u ON l.user_refer = u.id
		ORDER BY l.created_at DESC
	`,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var res []LendingAdminResponse
	for rows.Next() {
		var temp LendingAdminResponse
		err := rows.Scan(
			&temp.Id, &temp.UserId, &temp.Username, &temp.Amount, &temp.InterestRate, &temp.Tenor, &temp.Age,
			&temp.Gender, &temp.Income, &temp.LastEducation, &temp.MaritalStatus, &temp.NumberOfChildren,
			&temp.HasHouse, &temp.KkUrl, &temp.KtpUrl, &temp.Status, &temp.PaymentToken, &temp.PaymentUrl,
			&temp.IsApproved, &temp.IsRejected,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, temp)
	}
	return res, nil
}

func PredictCreditScore(id string) (LendingPredictResponse, error) {
	var req LendingPredictRequest
	err := database.MysqlInstance.QueryRow(
		`
	SELECT l.age, l.gender, FLOOR(l.income), l.last_education, l.marital_status, l.number_of_children, l.home_ownership FROM lending l
	WHERE l.id = UUID_TO_BIN(?)
	`, id,
	).Scan(
		&req.Age, &req.Gender, &req.Income, &req.Education, &req.MaritalStatus, &req.NumberOfChildren,
		&req.HomeOwnership,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LendingPredictResponse{}, fmt.Errorf("order id not found")
		}
		return LendingPredictResponse{}, err
	}

	url := flaskMLBaseUrl + "/predict"
	body, err := json.Marshal(req)
	if err != nil {
		return LendingPredictResponse{}, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return LendingPredictResponse{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return LendingPredictResponse{}, fmt.Errorf("ml server error")
	}

	var res LendingPredictResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return LendingPredictResponse{}, err
	}
	return res, nil
}

func ApproveLending(id string) error {
	res, err := database.MysqlInstance.Exec(
		`UPDATE lending SET status = 'approved', is_approved = TRUE WHERE id = UUID_TO_BIN(?) AND is_rejected = FALSE`,
		id,
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func RejectLending(id string) error {
	res, err := database.MysqlInstance.Exec(
		`UPDATE lending SET status = 'rejected', is_rejected = TRUE WHERE id = UUID_TO_BIN(?) AND is_approved = FALSE`,
		id,
	)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("not found")
	}
	return nil
}

func MakePayment(id string) (midtrans.ResponseSnap, error) {
	var exits int
	err := database.MysqlInstance.QueryRow(
		`SELECT 1 FROM lending WHERE id = UUID_TO_BIN(?) AND is_approved = TRUE AND payment_token IS NULL`, id,
	).Scan(&exits)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return midtrans.ResponseSnap{}, fmt.Errorf("lending not found")
		}
		return midtrans.ResponseSnap{}, err
	}
	if exits == 0 {
		return midtrans.ResponseSnap{}, fmt.Errorf("lending not found")
	}

	var transDetails midtrans.TransactionDetails
	transDetails.OrderId = midtrans.BaseOrderId + "-" + id
	err = database.MysqlInstance.QueryRow("SELECT FLOOR(amount * (100 + interest_rate) / 100) FROM lending").Scan(&transDetails.GrossAmount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return midtrans.ResponseSnap{}, fmt.Errorf("lending not found")
		}
		return midtrans.ResponseSnap{}, err
	}

	var snapReq midtrans.RequestSnap
	snapReq.TransactionDetails = transDetails
	snapReq.Expiry = midtrans.Expiry{
		//	StartTime: will be time.Now() utc to string with format "2020-06-30 15:07:00 -0700"
		StartTime: time.Now().Format("2006-01-02 15:04:05 -0700"),
		Unit:      "day",
		Duration:  30,
	}

	res, err := snapReq.CreatePayment()
	if err != nil {
		return midtrans.ResponseSnap{}, err
	}
	go func() {
		_, err := database.MysqlInstance.Exec(
			"UPDATE lending SET payment_token = ?, payment_url = ? WHERE id = UUID_TO_BIN(?)", res.Token,
			res.RedirectUrl, id,
		)
		if err != nil {
			log.Print(err)
		}
	}()
	return res, nil
}
