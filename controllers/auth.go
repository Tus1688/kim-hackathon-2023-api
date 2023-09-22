package controllers

import (
	"net/http"
	"strings"

	"github.com/Tus1688/kim-hackathon-2023-api/jsonutil"
	"github.com/Tus1688/kim-hackathon-2023-api/models"
	"github.com/Tus1688/kim-hackathon-2023-api/render"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := jsonutil.ShouldBind(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jwt, ref, err, res := req.Login()
	if err != nil {
		if strings.Contains(err.Error(), "invalid") {
			render.HandleError([]string{err.Error()}, http.StatusUnauthorized, w)
			return
		}
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}

	access := http.Cookie{
		Name:     "access",
		Value:    jwt,
		Path:     "/api/v1",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	refresh := http.Cookie{
		Name:     "refresh",
		Value:    ref,
		Path:     "/api/v1/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &access)
	http.SetCookie(w, &refresh)

	err = render.JSON(w, http.StatusOK, res)
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
	}
}

func GetRefreshToken(w http.ResponseWriter, r *http.Request) {
	refToken, err := r.Cookie("refresh")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	newToken, err := models.GetRefreshToken(refToken.Value)
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusUnauthorized, w)
		return
	}
	access := http.Cookie{
		Name:     "access",
		Value:    newToken,
		Path:     "/api/v1",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &access)

	w.WriteHeader(http.StatusOK)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUser
	if err := jsonutil.ShouldBind(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := req.Create()
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			render.HandleError([]string{"username already exists"}, http.StatusForbidden, w)
			return
		}

		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	res, err := models.GetAllUsers()
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}
	err = render.JSON(w, http.StatusOK, res)
	if err != nil {
		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := models.DeleteUser(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			render.HandleError([]string{"user not found"}, http.StatusNotFound, w)
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

func ModifyUser(w http.ResponseWriter, r *http.Request) {
	var req models.ModifyUser
	if err := jsonutil.ShouldBind(r, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := req.Modify()
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			render.HandleError([]string{"user not found"}, http.StatusNotFound, w)
			return
		}

		if strings.Contains(err.Error(), "cannot") {
			render.HandleError([]string{err.Error()}, http.StatusForbidden, w)
			return
		}

		render.HandleError([]string{err.Error()}, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
