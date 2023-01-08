package main

import (
	"dumbmerch/database"
	"dumbmerch/pkg/mysql"
	"dumbmerch/routes"
	"fmt"
	"net/http"
	"os"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	// DB CONNECT
	mysql.DatabaseInit()

	// MIGRATION SET UP
	database.RunMigration()

	r := mux.NewRouter()

	var AllowedHeaders = handlers.AllowedHeaders([]string{"X-Requested-Width", "Content-Type", "Authorization"})
	var AllowedMethods = handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "PATCH", "DELETE"})
	var AllowedOrigins = handlers.AllowedOrigins([]string{"*"})

	r.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	routes.RouteInit(r.PathPrefix("/api/v1").Subrouter())
	
	port := os.Getenv("PORT")
	
	fmt.Println("Server is success to Connect")
	http.ListenAndServe(":" +port, handlers.CORS(AllowedHeaders, AllowedMethods, AllowedOrigins)(r))
}
