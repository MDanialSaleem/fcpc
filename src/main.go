package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// using sync.Map instead of map+mutex because the requirements for this app fall specifically into what sync.Map
// is recommended for: https://pkg.go.dev/sync#Map
var receiptStore = sync.Map{}
var logger *zap.Logger

func main() {

	router := setup()
	defer logger.Sync()

	logger.Info("Starting server on port 8000")
	http.ListenAndServe(":8000", router)
}

func setup() *mux.Router {
	logLevel := os.Getenv("LOG_LEVEL")
	var err error

	if logLevel == "DEBUG" {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic("failed to initialize logger")
	}

	router := mux.NewRouter()

	router.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")
	router.HandleFunc("/receipts/process", processReceipt).Methods("POST")

	return router
}

func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)

	if err != nil {
		logger.Debug("Failed to decode receipt", zap.Error(err))
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}
	logger.Debug("Received receipt", zap.Any("receipt", receipt))

	receiptID := uuid.New().String()
	logger.Debug("Generated UUID", zap.String("receiptID", receiptID))

	// very unlikely, but just in case.
	if _, ok := receiptStore.Load(receiptID); ok {
		logger.Error("Duplicate UUID generated", zap.String("receiptID", receiptID))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	points := receipt.CalculatePoints()
	receiptStore.Store(receiptID, int64(points))
	logger.Debug("Stored receipt points", zap.String("receiptID", receiptID), zap.Int("points", points))

	jsonResponse, err := json.Marshal(map[string]string{"id": receiptID})
	if err != nil {
		logger.Error("Failed to marshal response", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	logger.Debug("Getting points for receipt", zap.String("receiptID", id))

	points, ok := receiptStore.Load(id)
	if !ok {
		http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
		return
	}

	response := map[string]int64{"points": points.(int64)}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
