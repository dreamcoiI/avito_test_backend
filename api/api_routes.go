package api

import (
	"github.com/dreamcoiI/avito_test_backend/internal/handlers"
	"github.com/gorilla/mux"
)

func ConfigureRouters(h *handlers.Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/user_segment/find/{user_id:[0-9]+}", h.GetUserSegment).Methods("GET")
	router.HandleFunc("/user_segment", h.CreateUserSegment).Methods("POST")
	//router.HandleFunc()
	return router
}