package api

import (
	schedule "empshift-csp/internal/core"
	"empshift-csp/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleScheduleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.SchedulePackageRequest
	fmt.Printf("NUMBER OF EMPLOYEES: %d\n", len(req.Staffs))
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResponse := map[string]string{
			"message": "Invalid request body",
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	result, err := schedule.ComputeSchedule(req)
	if err != nil {
		errorResponse := map[string]string{
			"message": "Failed to compute schedule",
			"error":   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
