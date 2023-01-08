package handlers

import (
	dto "dumbmerch/dto/result"
	tripsdto "dumbmerch/dto/trip"
	"dumbmerch/models"
	"dumbmerch/repositories"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var Path_File = "http://localhost:5000/uploads/"

type handlerTrip struct {
	TripRepository repositories.TripRepository
}

func HandlerTrip(TripRepository repositories.TripRepository) *handlerTrip {
	return &handlerTrip{TripRepository}
}

func (h *handlerTrip) FindTrips(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	trips, err := h.TripRepository.FindTrips()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: trips}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTrip) GetTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/form-data")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	trip, err := h.TripRepository.GetTrip(id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: trip}
	json.NewEncoder(w).Encode(response)

}

func (h *handlerTrip) AddTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dataContext := r.Context().Value("dataFile")
	filename := Path_File + dataContext.(string)

	country_id, _ := strconv.Atoi(r.FormValue("country_id"))
	price, _ := strconv.Atoi(r.FormValue("price"))
	quota, _ := strconv.Atoi(r.FormValue("quota"))
	request := tripsdto.TripRequest{
		Title:          r.FormValue("title"),
		CountryID:      country_id,
		Accomodation:   r.FormValue("accomodation"),
		Transportation: r.FormValue("transportation"),
		Eat:            r.FormValue("eat"),
		Duration:       r.FormValue("duration"),
		Date_Trip:      r.FormValue("date_trip"),
		Price:          price,
		Quota:          quota,
		Description:    r.FormValue("description"),
	}
	validation := validator.New()
	err := validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	trip := models.Trip{

		Title:          request.Title,
		Accomodation:   request.Accomodation,
		Transportation: request.Transportation,
		Eat:            request.Eat,
		Duration:       request.Duration,
		Date_Trip:      request.Date_Trip,
		Price:          request.Price,
		Quota:          request.Quota,
		Description:    request.Description,
		Image:          filename,
		CountryID:      request.CountryID,
	}
	trip, err = h.TripRepository.AddTrip(trip)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	trip, _ = h.TripRepository.GetTrip(trip.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: trip}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTrip) EditTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dataContext := r.Context().Value("dataFile")
	filename := Path_File + dataContext.(string)

	country_id, _ := strconv.Atoi(r.FormValue("country_id"))
	price, _ := strconv.Atoi(r.FormValue("price"))
	quota, _ := strconv.Atoi(r.FormValue("quota"))
	request := tripsdto.TripRequest{
		Title:          r.FormValue("title"),
		CountryID:      country_id,
		Accomodation:   r.FormValue("accomodation"),
		Transportation: r.FormValue("transportation"),
		Eat:            r.FormValue("eat"),
		Duration:       r.FormValue("duration"),
		Date_Trip:      r.FormValue("date_trip"),
		Price:          price,
		Quota:          quota,
		Description:    r.FormValue("description"),
		Image:          filename,
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	trip, err := h.TripRepository.GetTrip(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	if request.Title != "" {
		trip.Title = request.Title
	}
	if request.CountryID != 0 {
		trip.CountryID = request.CountryID
	}
	if request.Accomodation != "" {
		trip.Accomodation = request.Accomodation
	}
	if request.Transportation != "" {
		trip.Transportation = request.Transportation
	}
	if request.Eat != "" {
		trip.Eat = request.Eat
	}
	if request.Duration != "" {
		trip.Duration = request.Duration
	}
	if request.Date_Trip != "" {
		trip.Date_Trip = request.Date_Trip
	}
	if request.Price != 0 {
		trip.Price = request.Price
	}
	if request.Quota != 0 {
		trip.Quota = request.Quota
	}
	if request.Description != "" {
		trip.Description = request.Description
	}
	if request.Image != "" {
		trip.Image = filename
	}

	data, err := h.TripRepository.EditTrip(trip, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTrip) DeleteTrip(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	trip, err := h.TripRepository.GetTrip(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.TripRepository.DeleteTrip(trip, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: data}
	json.NewEncoder(w).Encode(response)
}
