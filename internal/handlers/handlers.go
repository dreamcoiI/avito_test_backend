package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dreamcoiI/avito_test_backend/internal/model"
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

func (h *Handler) GetUserSegment(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UserID int `json:"user_id"`
	}
	fmt.Println("Денис Абоба")
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		WrapError(w, err)
		return
	}

	ctx := r.Context()
	segment, err := h.service.GetUserSegment(ctx, requestData.UserID)
	if err != nil {
		WrapError(w, err)
		return
	}

	response := map[string]interface{}{
		"result": "OK",
		"data":   segment,
	}

	resp, err := json.Marshal(response)
	if err != nil {
		WrapError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resp)
	if err != nil {
		WrapError(w, err)
		return
	}

	WrapOK(w, response)
}

//func (h *Handler) GetUserSegment(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	if vars["user_id"] == "" {
//		WrapError(w, errors.New("missing id"))
//		return
//	}
//
//	userID, err := strconv.Atoi(vars["user_id"])
//	if err != nil {
//		WrapError(w, errors.New("wrong id"))
//		return
//	}
//
//	ctx := r.Context()
//	segment, err := h.service.GetUserSegment(ctx, userID)
//	if err != nil {
//		WrapError(w, err)
//		return
//	}
//
//	response := map[string]interface{}{
//		"result": "OK",
//		"data":   segment,
//	}
//
//	resp, err := json.Marshal(response)
//	if err != nil {
//		WrapError(w, err)
//		return
//	}
//
//	_, err = w.Write(resp)
//	if err != nil {
//		WrapError(w, err)
//		return
//	}
//
//	WrapOK(w, response)
//
//}

func (h *Handler) CreateUserSegment(w http.ResponseWriter, r *http.Request) {
	var newUserSegment model.UserSegment

	err := json.NewDecoder(r.Body).Decode(&newUserSegment)
	if err != nil {
		WrapError(w, err)
		return
	}

	err = h.service.CreateUserSegment(&newUserSegment)
	if err != nil {
		WrapError(w, err)
		return
	}

	response := map[string]interface{}{
		"result": "OK",
		"data":   newUserSegment,
	}

	WrapOK(w, response)
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
