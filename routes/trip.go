package routes

import (
	"dumbmerch/handlers"
	"dumbmerch/pkg/middleware"
	"dumbmerch/pkg/mysql"
	"dumbmerch/repositories"

	"github.com/gorilla/mux"
)

func TripRoutes(r *mux.Router) {
	tripRepository := repositories.RepositoryTrip(mysql.DB)
	h := handlers.HandlerTrip(tripRepository)

	r.HandleFunc("/trip", h.FindTrips).Methods("GET")
	r.HandleFunc("/trip/{id}", h.GetTrip).Methods("GET")
	r.HandleFunc("/trip", middleware.Auth(middleware.UploadFile(h.AddTrip))).Methods("POST")
	r.HandleFunc("/trip/{id}", middleware.Auth(middleware.UploadFile(h.EditTrip))).Methods("PATCH")
	r.HandleFunc("/trip/{id}", middleware.Auth(h.DeleteTrip)).Methods("DELETE")
}
