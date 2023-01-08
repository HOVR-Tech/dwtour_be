package routes

import (
	"dumbmerch/handlers"
	"dumbmerch/pkg/middleware"
	"dumbmerch/pkg/mysql"
	"dumbmerch/repositories"

	"github.com/gorilla/mux"
)

func TransactionRoutes(r *mux.Router) {
	transactionRepository := repositories.RepositoryTransaction(mysql.DB)
	h := handlers.HandlerTransaction(transactionRepository)

	r.HandleFunc("/transactions", middleware.Auth(h.FindTransactions)).Methods("GET")
	r.HandleFunc("/transaction/{id}", middleware.Auth(h.GetTransaction)).Methods("GET")
	r.HandleFunc("/transactions-user/{id}", middleware.Auth(h.GetTransactionByUserID)).Methods("GET")
	r.HandleFunc("/transaction", middleware.Auth(h.BookTransaction)).Methods("POST")
	r.HandleFunc("/snap/{id}", middleware.Auth(h.Snap)).Methods("GET")
	r.HandleFunc("/notification", h.Notification).Methods("POST")
	r.HandleFunc("/transaction-upload/{id}", middleware.Auth(middleware.UploadFile(h.UpdateTransaction))).Methods("PATCH")
	r.HandleFunc("/transaction/{id}", middleware.Auth(h.CheckTransaction)).Methods("PATCH")
	r.HandleFunc("/transaction/{id}", middleware.Auth(h.DeleteTransaction)).Methods("DELETE")
}
