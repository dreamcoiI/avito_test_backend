package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dreamcoiI/avito_test_backend/internal/service"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(Service *service.Service) *Handler {
	newHandler := new(Handler)
	newHandler.service = Service
	return newHandler
}

func WrapOK(w http.ResponseWriter, m map[string]interface{}) {
	res, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Println(w, string(res))
}

func WrapErrorStatus(w http.ResponseWriter, err error, httpStatus int) {
	var m = map[string]string{
		"result": "Error",
		"data":   err.Error(),
	}
	res, _ := json.Marshal(m)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(httpStatus)
	fmt.Println(w, string(res))
}

func WrapError(w http.ResponseWriter, err error) {
	WrapErrorStatus(w, err, http.StatusBadRequest)
}

func WrapNotFound(w http.ResponseWriter, r *http.Request) {
	WrapErrorStatus(w, errors.New("not found"), http.StatusNotFound)
}
