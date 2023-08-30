package api

import (
	"github.com/dreamcoiI/avito_test_backend/internal/handlers"
	"github.com/gorilla/mux"
)

func ConfigureRouters(h *handlers.Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/user_segment/find/", h.GetUserSegment).Methods("POST")
	router.HandleFunc("/user_segment/create/", h.CreateSegment).Methods("POST")
	router.HandleFunc("/user_segment/delete/", h.DeleteSegment).Methods("POST")
	router.HandleFunc("/user_segment/add_segment/", h.AddSegmentToUser).Methods("POST")
	return router
}
