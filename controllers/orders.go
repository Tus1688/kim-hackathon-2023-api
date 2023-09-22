package controllers

import (
	"net/http"
	"strings"

	"github.com/Tus1688/kim-hackathon-2023-api/jsonutil"
	"github.com/Tus1688/kim-hackathon-2023-api/models"
	"github.com/Tus1688/kim-hackathon-2023-api/render"
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.CreateOrder
	if err := jsonutil.ShouldBind(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := req.Create()
	if err != nil {
		if strings.Contains(err.Error(), "uuid_to_bin") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if strings.Contains(err.Error(), "product_refer") {
			render.HandleError([]string{"product id not found"}, http.StatusConflict, w)
			return
		}

		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetOrder(w http.ResponseWriter, r *http.Request) {
	res, err := models.GetOrder()
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}

	if len(res) == 0 {
		render.HandleError([]string{"no order found"}, http.StatusNotFound, w)
		return
	}

	err = render.JSON(w, http.StatusOK, res)
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
	}
}

func ModifyOrder(w http.ResponseWriter, r *http.Request) {
	var req models.ModifyOrder
	if err := jsonutil.ShouldBind(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := req.Modify()
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			render.HandleError([]string{"order not found"}, http.StatusNotFound, w)
			return
		}

		if strings.Contains(err.Error(), "uuid_to_bin") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if strings.Contains(err.Error(), "product_refer") {
			render.HandleError([]string{"product id not found"}, http.StatusConflict, w)
			return
		}

		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
