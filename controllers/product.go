package controllers

import (
	"net/http"
	"strings"

	"github.com/Tus1688/kim-hackathon-2023-api/jsonutil"
	"github.com/Tus1688/kim-hackathon-2023-api/models"
	"github.com/Tus1688/kim-hackathon-2023-api/render"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProduct
	if err := jsonutil.ShouldBind(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := req.Create()
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			render.HandleError([]string{"product already exists"}, http.StatusConflict, w)
			return
		}

		if strings.Contains(err.Error(), "uuid_to_bin") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if strings.Contains(err.Error(), "business_refer") {
			render.HandleError([]string{"business id not found"}, http.StatusConflict, w)
			return
		}

		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}

	err = render.JSON(w, http.StatusCreated, res)
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
	}
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	res, err := models.GetProduct(query)
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}
	if len(res) == 0 {
		render.HandleError([]string{"no product found"}, http.StatusNotFound, w)
		return
	}

	err = render.JSON(w, http.StatusOK, res)
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
	}
}

func ModifyProduct(w http.ResponseWriter, r *http.Request) {
	var req models.ModifyProduct
	if err := jsonutil.ShouldBind(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := req.Modify()
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			render.HandleError([]string{"product not found"}, http.StatusNotFound, w)
			return
		}

		if strings.Contains(err.Error(), "uuid_to_bin") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if strings.Contains(err.Error(), "business_refer") {
			render.HandleError([]string{"business id not found"}, http.StatusConflict, w)
			return
		}

		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := models.DeleteProduct(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			render.HandleError([]string{"business not found"}, http.StatusNotFound, w)
			return
		}

		if strings.Contains(err.Error(), "uuid_to_bin") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateProductImage(w http.ResponseWriter, r *http.Request) {
	file, m, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	// validate image by checking extension of the file (jpg, jpeg, png only)
	if m.Filename[len(m.Filename)-3:] != "jpg" && m.Filename[len(m.Filename)-4:] != "jpeg" && m.Filename[len(m.Filename)-3:] != "png" {
		render.HandleError([]string{"invalid image format"}, http.StatusBadRequest, w)
		return
	}

	id := r.FormValue("id")
	res, err := models.CreateProductImage(id, file, m)
	if err != nil {
		if strings.Contains(err.Error(), "uuid_to_bin") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}

	err = render.JSON(w, http.StatusCreated, res)
}

func DeleteProductImage(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("file")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := models.DeleteProductImage(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			render.HandleError([]string{"image not found"}, http.StatusNotFound, w)
			return
		}

		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}
	w.WriteHeader(http.StatusOK)
}
