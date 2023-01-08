package handlers

import (
	"context"
	dto "dumbmerch/dto/result"
	transactiondto "dumbmerch/dto/transaction"
	"dumbmerch/models"
	"dumbmerch/repositories"
	"encoding/json"
	"fmt"

	// "github.com/cloudinary/cloudinary-go/v2/api/admin"

	"net/http"
	// "os"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	// "gopkg.in/gomail.v2"
)

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
}

func HandlerTransaction(TransactionRepository repositories.TransactionRepository) *handlerTransaction {
	return &handlerTransaction{TransactionRepository}
}

func (h *handlerTransaction) BookTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	qty, _ := strconv.Atoi(r.FormValue("qty"))
	trip_id, _ := strconv.Atoi(r.FormValue("trip_id"))
	total, _ := strconv.Atoi(r.FormValue("total"))
	request := transactiondto.BookRequest{
		Qty:    qty,
		Total:  total,
		Status: r.FormValue("status"),
		TripID: trip_id,
		UserID: userID,
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	validation := validator.New()
	err = validation.Struct(request)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	transaction := models.Transaction{
		Qty:    request.Qty,
		Total:  request.Total,
		Status: request.Status,
		TripID: request.TripID,
		UserID: request.UserID,
	}

	newTransaction, err := h.TransactionRepository.BookTransaction(transaction)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: newTransaction}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) Snap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/form-data")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	dataTransaction, err := h.TransactionRepository.GetTransaction(int(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	request := transactiondto.TransactionRequest{
		Status: r.FormValue("status"),
	}

	if request.Status != "" {
		dataTransaction.Status = request.Status
	}

	data, err := h.TransactionRepository.CheckTransaction(dataTransaction, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	trans, _ := h.TransactionRepository.GetTransaction(data.ID)
	fmt.Println()

	var s = snap.Client{}
	s.New("SB-Mid-server-yabfvfwIWlF_17N59lI1WuJ0", midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(trans.ID),
			GrossAmt: int64(trans.Total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: trans.User.Name,
			Email: trans.User.Email,
		},
	}

	snapResp, _ := s.CreateTransaction(req)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: snapResp}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "multipart/form-data")
	dataContext := r.Context().Value("dataFile")
	filepath := dataContext.(string)

	ctx := context.Background()
	CLOUD_NAME := "cloudme19"
	API_KEY := "281282276658861"
	API_SECRET := "uxlAfli9ExpwY2o6j8qTS0gRJ9g"

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	dataTransaction, err := h.TransactionRepository.GetTransaction(int(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	cld, _ := cloudinary.NewFromParams(CLOUD_NAME, API_KEY, API_SECRET)
	resp, err := cld.Upload.Upload(ctx, filepath, uploader.UploadParams{Folder: "dewe_tour"})
	if err != nil {
		fmt.Println("Upload Gagal!", err.Error())
	}

	request := transactiondto.TransactionRequest{
		Image: resp.SecureURL,
	}

	if request.Image != "" {
		dataTransaction.Image = request.Image
	}

	data, err := h.TransactionRepository.UpdateTransaction(dataTransaction, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	trans, _ := h.TransactionRepository.GetTransaction(data.ID)

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: trans}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) Notification(w http.ResponseWriter, r *http.Request) {
	var notificationPayload map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&notificationPayload)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderId := notificationPayload["order_id"].(string)

	transaction, _ := h.TransactionRepository.GetOneTransaction(orderId)
	fmt.Println(transactionStatus, fraudStatus, orderId, transaction)
	fmt.Println(notificationPayload)
	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			SendEmail("pending", transaction)
			h.TransactionRepository.Notification("pending", transaction.ID)
		} else if fraudStatus == "accept" {
			SendEmail("Waiting Approval", transaction)
			h.TransactionRepository.Notification("Waiting Approval", transaction.ID)
		} else if transactionStatus == "settlement" {
			SendEmail("Waiting Approval", transaction)
			h.TransactionRepository.Notification("Waiting Approval", transaction.ID)
		} else if transactionStatus == "deny" {
			SendEmail("failed", transaction)
			h.TransactionRepository.Notification("failed", transaction.ID)
		} else if transactionStatus == "cancel" {
			SendEmail("failed", transaction)
			h.TransactionRepository.Notification("failed", transaction.ID)
		} else if transactionStatus == "pending" {
			SendEmail("pending", transaction)
			h.TransactionRepository.Notification("pending", transaction.ID)
		}

		w.WriteHeader(http.StatusOK)
	}
}

func SendEmail(status string, transaction models.Transaction) {
	// var CONFIG_SMTP_HOST = "smtp.gmail.com"
	// var CONFIG_SMTP_PORT = 587
	// var CONFIG_SENDER_NAME = "DeweTour <hydrilla.salim@gmail.com>"
	// var CONFIG_AUTH_EMAIL = os.Getenv("EMAIL_SYSTEM")
	// var CONFIG_AUTH_PASSWORD = os.Getenv("PASSWORD_SYSTEM")

	// var tripTitle = transaction.Trip.Title
	// var price = strconv.Itoa(transaction.Total)

	// mailer := gomail.NewMessage()
	// mailer.SetHeader("from", CONFIG_SENDER_NAME)
	// mailer.SetHeader("To", transaction.User.Email)
	// mailer.SetHeader("Subject", "Status Transaksi AMAAAN")
	// mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
	//  <html lang="en">
	//    <head>
	//    <meta charset="UTF-8" />
	//    <meta http-equiv="X-UA-Compatible" content="IE=edge" />
	//    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
	//    <title>Document</title>
	//    <style>
	//      h1 {
	//      color: brown;
	//      }
	//    </style>
	//    </head>
	//    <body>
	//    <h2>Product payment :</h2>
	//    <ul style="list-style-type:none;">
	//      <li>Name : %s</li>
	//      <li>Total payment: Rp.%s</li>
	//      <li>Status : <b>%s</b></li>
	//      <li>Iklan : <b>%s</b></li>
	//    </ul>
	//    </body>
	//  </html>`, tripTitle, price, status, "TOUR SURGA DISKON 99%"))

	// dialer := gomail.NewDialer(
	// 	CONFIG_SMTP_HOST,
	// 	CONFIG_SMTP_PORT,
	// 	CONFIG_AUTH_EMAIL,
	// 	CONFIG_AUTH_PASSWORD,
	// )

	// err := dialer.DialAndSend(mailer)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// log.Println("Pesan Terkirim")
}

func (h *handlerTransaction) FindTransactions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	transaction, err := h.TransactionRepository.FindTransactions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: transaction}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) GetTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	transaction, err := h.TransactionRepository.GetTransaction(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: transaction}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) GetTransactionByUserID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	transaction, err := h.TransactionRepository.GetTransactionByUserID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := dto.ErrorResult{Code: http.StatusInternalServerError, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := dto.SuccessResult{Code: http.StatusOK, Data: transaction}
	json.NewEncoder(w).Encode(response)
}

func (h *handlerTransaction) CheckTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// userInfo := r.Context().Value("userInfo").(jwt.MapClaims)
	// userId := int(userInfo["id"].(float64))

	// dataContext := r.Context().Value("dataFile")
	// filename := Path_File + dataContext.(string)

	var status string
	json.NewDecoder(r.Body).Decode(&status)

	fmt.Println(status)

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	transaction, err := h.TransactionRepository.GetTransaction(int(id))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	// if request.UserID != 0 {
	// 	transaction.UserID = request.UserID
	// }
	if status != "" {
		transaction.Status = status
	}

	data, err := h.TransactionRepository.CheckTransaction(transaction, transaction.ID)
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

func (h *handlerTransaction) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	trip, err := h.TransactionRepository.GetTransaction(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()}
		json.NewEncoder(w).Encode(response)
		return
	}

	data, err := h.TransactionRepository.DeleteTransaction(trip, id)
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
