package main

import (
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "aws-waf-ip-sets-management/backend/routes"
)

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/api/create-ip-set", routes.CreateIPSet).Methods("POST")
    r.HandleFunc("/api/list-ip-sets", routes.ListIPSets).Methods("GET")
    r.HandleFunc("/api/add-ip-address", routes.AddIPAddress).Methods("POST")
    r.HandleFunc("/api/remove-ip-address", routes.RemoveIPAddress).Methods("POST")
    r.HandleFunc("/api/delete-ip-set", routes.DeleteIPSet).Methods("POST")

    r.PathPrefix("/").Handler(http.FileServer(http.Dir("../frontend")))

    log.Println("Server is running on port 3000")
    log.Fatal(http.ListenAndServe(":3000", r))
}
