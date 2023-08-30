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

func (h *Handler) GetUserSegment(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		UserID int `json:"user_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		WrapError(w, err)
		return
	}

	fmt.Println(requestData.UserID)
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

func (h *Handler) CreateSegment(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Slug string `json:"slug"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		WrapError(w, err)
		return
	}

	fmt.Println(requestData.Slug + "ABOBA")
	ctx := r.Context()
	err = h.service.CreateSegment(ctx, requestData.Slug)
	if err != nil {
		WrapError(w, err)
		return
	}

	response := map[string]interface{}{
		"result": "OK",
		"data":   requestData.Slug,
	}

	WrapOK(w, response)
}

func (h *Handler) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Slug string `json:"slug"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		WrapError(w, err)
		return
	}

	fmt.Println(requestData.Slug + "ABOBA")

	ctx := r.Context()
	err = h.service.DeleteSegment(ctx, requestData.Slug)
	if err != nil {
		WrapError(w, err)
		return
	}

	response := map[string]interface{}{
		"result": "OK",
		"data":   requestData.Slug,
	}

	WrapOK(w, response)
}

func (h *Handler) AddSegmentToUser(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Add    []string `json:"add"`
		UserID int      `json:"user_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		WrapError(w, err)
		return
	}

	ctx := r.Context()

	err = h.service.AddSegmentToUser(ctx, requestData.Add, requestData.UserID)
	if err != nil {
		WrapError(w, err)
		return
	}

	response := map[string]interface{}{
		"result": "OK",
		"add":    requestData.Add,
	}

	WrapOK(w, response)
}

func (h *Handler) DeleteSegmentToUser(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Delete []string `json:"delete"`
		UserID int      `json:"user_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		WrapError(w, err)
		return
	}

	ctx := r.Context()

	err = h.service.DeleteSegmentToUser(ctx, requestData.Delete, requestData.UserID)
	if err != nil {
		WrapError(w, err)
		return
	}

	response := map[string]interface{}{
		"result": "OK",
		"add":    requestData.Delete,
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
