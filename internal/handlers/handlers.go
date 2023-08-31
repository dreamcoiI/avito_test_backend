package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dreamcoiI/avito_test_backend/internal/service"
	"net/http"
	"os"
	"path/filepath"
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

func (h *Handler) DeleteSegment(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Slug string `json:"slug"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		WrapError(w, err)
		return
	}

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

func (h *Handler) AddAndDeleteSegmentToUser(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Add    []string `json:"add"`
		Delete []string `json:"delete"`
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
	err = h.service.DeleteSegmentToUser(ctx, requestData.Delete, requestData.UserID)
	if err != nil {
		WrapError(w, err)
		return
	}

	response := map[string]interface{}{
		"result": "OK",
		"add":    requestData.Add,
		"delete": requestData.Delete,
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

func (h *Handler) GenerateSegmentHistoryCSVByMonth(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		Year  int `json:"year"`
		Month int `json:"month"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		WrapError(w, err)
		return
	}

	ctx := r.Context()

	filename := "app/static/user_segment_history.csv"

	file, err := os.Create(filename)
	if err != nil {
		WrapError(w, err)
		return
	}
	defer file.Close()

	filePath, err := h.service.GenerateSegmentHistoryCSVByMonth(ctx, requestData.Year, requestData.Month, filename)
	if err != nil {
		WrapError(w, err)
		return
	}

	w.Header().Set("Content-Type", filePath)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "user_segment_history.csv"))

	response := map[string]interface{}{
		"result":          "OK",
		"filepath to CSV": "/" + filePath,
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

	filePath = filepath.Join("static", r.URL.Path)
	fmt.Println("Requested file path:", filePath)
	http.ServeFile(w, r, filePath)
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
