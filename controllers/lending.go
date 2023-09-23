package controllers

import (
	"net/http"
	"strings"

	"github.com/Tus1688/kim-hackathon-2023-api/jsonutil"
	"github.com/Tus1688/kim-hackathon-2023-api/models"
	"github.com/Tus1688/kim-hackathon-2023-api/render"
)

func RegisterAsBorrower(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterAsBorrower
	if err := jsonutil.ShouldBind(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := req.Register()
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			render.HandleError([]string{"user already exists"}, http.StatusConflict, w)
			return
		}

		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UploadDocument(w http.ResponseWriter, r *http.Request) {
	file, m, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	res, err := models.UploadDocument(file, m)
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}

	err = render.JSON(w, http.StatusCreated, res)
}

func CreateLendingProposal(w http.ResponseWriter, r *http.Request) {
	var req models.LendingRequest
	if err := jsonutil.ShouldBind(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req.RequesterUid = r.Context().Value("uid").(string)

	err := req.Create()
	if err != nil {
		if strings.Contains(err.Error(), "uuid_to_bin") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if strings.Contains(err.Error(), "borrower_refer") {
			render.HandleError([]string{"borrower id not found"}, http.StatusConflict, w)
			return
		}

		if strings.Contains(err.Error(), "lender_refer") {
			render.HandleError([]string{"lender id not found"}, http.StatusConflict, w)
			return
		}

		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetLendingProposalUser(w http.ResponseWriter, r *http.Request) {
	uid := r.Context().Value("uid").(string)
	res, err := models.GetLendingAsUser(uid)
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}
	if len(res) == 0 {
		render.HandleError([]string{"no data found"}, http.StatusNotFound, w)
		return
	}

	err = render.JSON(w, http.StatusOK, res)
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
	}
}

func GetLendingProposalAdmin(w http.ResponseWriter, r *http.Request) {
	res, err := models.GetLendingAsAdmin()
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}
	if len(res) == 0 {
		render.HandleError([]string{"no data found"}, http.StatusNotFound, w)
		return
	}

	err = render.JSON(w, http.StatusOK, res)
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
	}
}

func PredictCreditScore(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res, err := models.PredictCreditScore(id)
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}

	err = render.JSON(w, http.StatusOK, res)
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
	}
}
