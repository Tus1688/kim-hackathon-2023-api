package midtrans

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Tus1688/kim-hackathon-2023-api/database"
	"github.com/Tus1688/kim-hackathon-2023-api/jsonutil"
)

var ServerKey string
var ServerKeyEncoded string
var BaseUrlSnap string
var BaseUrlCoreApi string

// BaseOrderId is used to prefix the order id in database
// for example if the order id is 1, then the order id in midtrans is "something-1"
var BaseOrderId string

func (r *RequestSnap) CreatePayment() (ResponseSnap, error) {
	url := BaseUrlSnap + "/snap/v1/transactions"
	body, err := json.Marshal(r)
	if err != nil {
		return ResponseSnap{}, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return ResponseSnap{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Basic "+ServerKeyEncoded)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return ResponseSnap{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 201 {
		var result ResponseErrorSnap
		err = json.NewDecoder(res.Body).Decode(&result)
		if err != nil {
			return ResponseSnap{}, err
		}
		return ResponseSnap{}, fmt.Errorf("%v", result.ErrorMessages)
	}
	var result ResponseSnap
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return ResponseSnap{}, err
	}
	return result, nil
}

func HandleNotifications(w http.ResponseWriter, r *http.Request) {
	var request WebhookNotification
	if err := jsonutil.ShouldBind(r, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// verify the signature
	// SHA512(order_id+status_code+gross_amount+ServerKey)
	hash := sha512.New()
	hash.Write([]byte(request.OrderId + request.StatusCode + request.GrossAmount + ServerKey))
	signature := fmt.Sprintf("%x", hash.Sum(nil))
	// if not verified, return 401
	if signature != request.SignatureKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// strip the BaseOrderId+"-" from the order id
	// for example if the order id is "something-1", then the order id in database is 1
	OrderId := request.OrderId[len(BaseOrderId)+1:]
	// if there is FraudStatus, always check if it is "accept"
	// if transaction_status value is settlement or capture change the is_paid to true
	if request.TransactionStatus == "settlement" || request.TransactionStatus == "capture" {
		if request.FraudStatus != "deny" && request.FraudStatus != "challenge" {
			_, err := database.MysqlInstance.Exec(
				`UPDATE lending set is_paid = true, status = ? WHERE id = UUID_TO_BIN(?)`, request.TransactionStatus,
				OrderId,
			)
			if err != nil {
				log.Print(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	} else {
		_, err := database.MysqlInstance.Exec(
			`UPDATE lending set status = ? WHERE id = UUID_TO_BIN(?)`, request.TransactionStatus, OrderId,
		)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
