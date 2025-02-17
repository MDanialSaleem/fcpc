package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFullCycle(t *testing.T) {
	testCases := []struct {
		name           string
		requestBody    string
		wantPointsResp int64
	}{
		{
			name: "readme example 1: not round dollar, not multiple of 0.25, odd day, not special time",
			requestBody: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [
					{
					"shortDescription": "Mountain Dew 12PK",
					"price": "6.49"
					},{
					"shortDescription": "Emils Cheese Pizza",
					"price": "12.25"
					},{
					"shortDescription": "Knorr Creamy Chicken",
					"price": "1.26"
					},{
					"shortDescription": "Doritos Nacho Cheese",
					"price": "3.35"
					},{
					"shortDescription": "   Klarbrunn 12-PK 12 FL OZ  ",
					"price": "12.00"
					}
				],
				"total": "35.35"
			}`,
			wantPointsResp: 28,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := setup()

			req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
				return
			}

			var resp map[string]string
			if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			pointsReq := httptest.NewRequest("GET", "/receipts/"+resp["id"]+"/points", nil)
			pointsRR := httptest.NewRecorder()

			router.ServeHTTP(pointsRR, pointsReq)

			if status := pointsRR.Code; status != http.StatusOK {
				t.Errorf("points handler returned wrong status code: got %v want %v", status, http.StatusOK)
				return
			}

			var pointsResp map[string]int64
			if err := json.Unmarshal(pointsRR.Body.Bytes(), &pointsResp); err != nil {
				t.Fatalf("Failed to parse points response: %v", err)
			}

			if points := pointsResp["points"]; points != tc.wantPointsResp {
				t.Errorf("wrong points calculation: got %v want %v", points, tc.wantPointsResp)
			}
		})
	}
}
func TestInvalidReceipt(t *testing.T) {
	testCases := []struct {
		name        string
		requestBody string
	}{
		{
			name: "invalid receipt: missing retailer",
			requestBody: `{
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [
					{
					"shortDescription": "Mountain Dew 12PK",
					"price": "6.49"
					}
				],
				"total": "6.49"
			}`,
		},
		{
			name: "invalid receipt: missing items",
			requestBody: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"total": "6.49"
			}`,
		},
		{
			name: "invalid receipt: invalid JSON",
			requestBody: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [
					{
					"shortDescription": "Mountain Dew 12PK",
					"price": "6.49"
					}
				],
				"total": "6.49"
			`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			router := setup()

			req := httptest.NewRequest("POST", "/receipts/process", bytes.NewBufferString(tc.requestBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusBadRequest {
				t.Errorf("handler returned wrong status code: got %v expected %v", status, http.StatusBadRequest)
			}

			expectedResponse := "The receipt is invalid.\n"
			if rr.Body.String() != expectedResponse {
				t.Errorf("handler returned unexpected body: got %v expected %v", rr.Body.String(), expectedResponse)
			}
		})
	}
}

func TestNonExistentReceipt(t *testing.T) {
	router := setup()

	req := httptest.NewRequest("GET", "/receipts/whatever/points", nil)
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}

	expectedResponse := "No receipt found for that ID.\n"
	if rr.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v expected %v", rr.Body.String(), expectedResponse)
	}
}
