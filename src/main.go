package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var receiptStore = make(map[string]int64)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")
	router.HandleFunc("/receipts/process", processReceipt).Methods("POST")

	fmt.Println("Starting server on port 8000...")
	http.ListenAndServe(":8000", router)
}

func processReceipt(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Processing receipt...")
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(receipt)

	receiptID := uuid.New().String()
	fmt.Println("Generated UUID:", receiptID)
	response := map[string]string{"id": receiptID}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	receiptStore[receiptID] = int64(receipt.CalculatePoints())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	receipt, ok := receiptStore[id]
	if !ok {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}
	response, err := json.Marshal(receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
