package render

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, statusCode int, v any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

type Err struct {
	ID    string   `json:"id,omitempty"`
	Error []string `json:"error"`
}

func HandleError(errMessage []string, statusCode int, w http.ResponseWriter) {
	err := JSON(
		w, statusCode, Err{
			Error: errMessage,
		},
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
