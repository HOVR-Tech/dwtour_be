package routes

import (
	"dumbmerch/handlers"
	"dumbmerch/pkg/middleware"
	"dumbmerch/pkg/mysql"
	"dumbmerch/repositories"

	"github.com/gorilla/mux"
)

func CountryRoutes(r *mux.Router) {
	countryRepository := repositories.RepositoryCountry(mysql.DB)
	h := handlers.HandlerCountry(countryRepository)

	r.HandleFunc("/countries", h.FindCountries).Methods("GET")
	r.HandleFunc("/country/{id}", h.GetCountry).Methods("GET")
	r.HandleFunc("/addCountry",middleware.Auth(h.AddCountry)).Methods("POST")
	r.HandleFunc("/editCountry/{id}",middleware.Auth(h.EditCountry)).Methods("PATCH")
	r.HandleFunc("/deleteCountry/{id}",middleware.Auth(h.DeleteCountry)).Methods("DELETE")
}