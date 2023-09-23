package controllers

import (
	"net/http"

	"github.com/Tus1688/kim-hackathon-2023-api/models"
	"github.com/Tus1688/kim-hackathon-2023-api/render"
)

func GetTotalUser(w http.ResponseWriter, r *http.Request) {
	res, err := models.GetTotalUser()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = render.JSON(w, http.StatusOK, res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetTotalSME(w http.ResponseWriter, r *http.Request) {
	res, err := models.GetTotalSME()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = render.JSON(w, http.StatusOK, res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetAwaitingApproval(w http.ResponseWriter, r *http.Request) {
	res, err := models.GetTotalAwaitingApproval()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = render.JSON(w, http.StatusOK, res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
