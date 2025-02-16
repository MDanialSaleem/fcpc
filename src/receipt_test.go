package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestReceiptUnmarshalJSON(t *testing.T) {
	testCases := []struct {
		name       string
		json       string
		want       Receipt
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "valid receipt",
			json: `
					{
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
					}
				`,
			want: Receipt{
				Retailer:     "Target",
				PurchaseDate: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				PurchaseTime: time.Date(0, 1, 1, 13, 1, 0, 0, time.UTC),
				Items: []Item{
					{
						ShortDescription: "Mountain Dew 12PK",
						Price:            6.49,
					},
					{
						ShortDescription: "Emils Cheese Pizza",
						Price:            12.25,
					},
					{
						ShortDescription: "Knorr Creamy Chicken",
						Price:            1.26,
					},
					{
						ShortDescription: "Doritos Nacho Cheese",
						Price:            3.35,
					},
					{
						ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
						Price:            12.00,
					},
				},
				Total: 35.35,
			},
			wantErr: false,
		},
		{
			name: "invalid retailer",
			json: `{
				"retailer": "Target!!!",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [],
				"total": "0.00"
			}`,
			wantErr:    true,
			wantErrMsg: "invalid retailer format: Target!!!. want alphanumeric characters, spaces, hyphens, and ampersands",
		},
		{
			name: "invalid date format",
			json: `{
				"retailer": "Target",
				"purchaseDate": "01-01-2022",
				"purchaseTime": "13:01",
				"items": [],
				"total": "0.00"
			}`,
			wantErr:    true,
			wantErrMsg: "invalid purchase date format: 01-01-2022. want YYYY-MM-DD format",
		},
		{
			name: "invalid time format",
			json: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "1:01 PM",
				"items": [],
				"total": "0.00"
			}`,
			wantErr:    true,
			wantErrMsg: "invalid purchase time format: 1:01 PM. want HH:MM format",
		},
		{
			name: "invalid total format",
			json: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [],
				"total": "0"
			}`,
			wantErr:    true,
			wantErrMsg: "invalid total format: 0. want 0.00 format",
		},
		{
			name: "invalid item description",
			json: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [{
					"shortDescription": "Mountain Dew!!!",
					"price": "1.25"
				}],
				"total": "1.25"
			}`,
			wantErr:    true,
			wantErrMsg: "invalid items: invalid short description format: Mountain Dew!!!. want alphanumeric characters, spaces, hyphens, and ampersands",
		},
		{
			name: "invalid item price format",
			json: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [{
					"shortDescription": "Mountain Dew",
					"price": "1.2"
				}],
				"total": "1.20"
			}`,
			wantErr:    true,
			wantErrMsg: "invalid items: invalid price format: 1.2. want 0.00 format",
		},
		{
			name: "invalid items length",
			json: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [],
				"total": "0.00"
			}`,
			wantErr:    true,
			wantErrMsg: "invalid items: must contain at least one item",
		},
		{
			name: "missing retailer",
			json: `{
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [{
					"shortDescription": "Mountain Dew",
					"price": "1.25"
				}],
				"total": "1.25"
			}`,
			wantErr:    true,
			wantErrMsg: "retailer is not valid json",
		},
		{
			name: "missing purchase date",
			json: `{
				"retailer": "Target",
				"purchaseTime": "13:01",
				"items": [{
					"shortDescription": "Mountain Dew",
					"price": "1.25"
				}],
				"total": "1.25"
			}`,
			wantErr:    true,
			wantErrMsg: "purchase date is not valid json",
		},
		{
			name: "missing purchase time",
			json: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"items": [{
					"shortDescription": "Mountain Dew",
					"price": "1.25"
				}],
				"total": "1.25"
			}`,
			wantErr:    true,
			wantErrMsg: "purchase time is not valid json",
		},
		{
			name: "missing items",
			json: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"total": "1.25"
			}`,
			wantErr:    true,
			wantErrMsg: "invalid items: unexpected end of JSON input",
		},
		{
			name: "missing total",
			json: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [{
					"shortDescription": "Mountain Dew",
					"price": "1.25"
				}]
			}`,
			wantErr:    true,
			wantErrMsg: "total is not valid json",
		},
		{
			name: "missing item short description",
			json: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [{
					"price": "1.25"
				}],
				"total": "1.25"
			}`,
			wantErr:    true,
			wantErrMsg: "invalid items: item short description is not valid json",
		},
		{
			name: "missing item price",
			json: `{
				"retailer": "Target",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "13:01",
				"items": [{
					"shortDescription": "Mountain Dew"
				}],
				"total": "1.25"
			}`,
			wantErr:    true,
			wantErrMsg: "invalid items: item price is not valid json",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var got Receipt
			err := json.Unmarshal([]byte(tc.json), &got)

			if (err != nil) != tc.wantErr {
				t.Errorf("%s: error = %v, wantErr %v", tc.name, err, tc.wantErr)
				return
			}

			if tc.wantErr {
				if err.Error() != tc.wantErrMsg {
					t.Errorf("%s: error message = %v, expected %v", tc.name, err.Error(), tc.wantErrMsg)
				}
				return
			}

			if !tc.wantErr {
				if got.Retailer != tc.want.Retailer {
					t.Errorf("%s: Retailer = %v, expected %v", tc.name, got.Retailer, tc.want.Retailer)
				}
				if !got.PurchaseDate.Equal(tc.want.PurchaseDate) {
					t.Errorf("%s: PurchaseDate = %v, expected %v", tc.name, got.PurchaseDate, tc.want.PurchaseDate)
				}
				if !got.PurchaseTime.Equal(tc.want.PurchaseTime) {
					t.Errorf("%s: PurchaseTime = %v, expected %v", tc.name, got.PurchaseTime, tc.want.PurchaseTime)
				}
				if got.Total != tc.want.Total {
					t.Errorf("%s: Total = %v, want %v", tc.name, got.Total, tc.want.Total)
				}
				if len(got.Items) != len(tc.want.Items) {
					t.Errorf("%s: Items length = %v, expected %v", tc.name, len(got.Items), len(tc.want.Items))
				}
				for i := range got.Items {
					if got.Items[i].ShortDescription != tc.want.Items[i].ShortDescription {
						t.Errorf("%s: Item[%d] ShortDescription = %v, expected %v", tc.name, i, got.Items[i].ShortDescription, tc.want.Items[i].ShortDescription)
					}
					if got.Items[i].Price != tc.want.Items[i].Price {
						t.Errorf("%s: Item[%d] Price = %v, expected %v", tc.name, i, got.Items[i].Price, tc.want.Items[i].Price)
					}
				}
			}
		})
	}
}

func TestReceiptPoints(t *testing.T) {
	testCases := []struct {
		name                   string
		receipt                Receipt
		want                   int
		wantRetailerPoints     int
		wantNoCentsPoints      int
		wantMultipleOf25Points int
		wantItemPairsPoints    int
		wantDescriptionPoints  int
		wantOddDayPoints       int
		wantTimePoints         int
	}{
		{
			name: "readme example 1: not round dollar, not multiple of 0.25, odd day, not special time",
			receipt: Receipt{
				Retailer:     "Target",
				PurchaseDate: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				PurchaseTime: time.Date(0, 1, 1, 13, 1, 0, 0, time.UTC),
				Items: []Item{
					{
						ShortDescription: "Mountain Dew 12PK",
						Price:            6.49,
					},
					{
						ShortDescription: "Emils Cheese Pizza",
						Price:            12.25,
					},
					{
						ShortDescription: "Knorr Creamy Chicken",
						Price:            1.26,
					},
					{
						ShortDescription: "Doritos Nacho Cheese",
						Price:            3.35,
					},
					{
						ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
						Price:            12.00,
					},
				},
				Total: 35.35,
			},
			want:                   28,
			wantRetailerPoints:     6,
			wantNoCentsPoints:      0,
			wantMultipleOf25Points: 0,
			wantItemPairsPoints:    10,
			wantDescriptionPoints:  6,
			wantOddDayPoints:       6,
			wantTimePoints:         0,
		},
		{
			name: "readme example 2: round dollar, multiple of 0.25, non-alphanumeric retailer name, not odd day, special time",
			receipt: Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: time.Date(2022, 3, 20, 0, 0, 0, 0, time.UTC),
				PurchaseTime: time.Date(0, 1, 1, 14, 33, 0, 0, time.UTC),
				Items: []Item{
					{
						ShortDescription: "Gatorade",
						Price:            2.25,
					},
					{
						ShortDescription: "Gatorade",
						Price:            2.25,
					},
					{
						ShortDescription: "Gatorade",
						Price:            2.25,
					},
					{
						ShortDescription: "Gatorade",
						Price:            2.25,
					},
				},
				Total: 9.00,
			},
			want:                   109,
			wantRetailerPoints:     14,
			wantNoCentsPoints:      50,
			wantMultipleOf25Points: 25,
			wantItemPairsPoints:    10,
			wantDescriptionPoints:  0,
			wantOddDayPoints:       0,
			wantTimePoints:         10,
		},
		{
			name: "multiple item descriptions having length multiple of 3 with differente prices",
			receipt: Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: time.Date(2022, 3, 20, 0, 0, 0, 0, time.UTC),
				PurchaseTime: time.Date(0, 1, 1, 14, 33, 0, 0, time.UTC),
				Items: []Item{
					{
						ShortDescription: "Gat",
						Price:            2.25,
					},
					{
						ShortDescription: "Gat",
						Price:            6.25,
					},
					{
						ShortDescription: "Gat",
						Price:            8.25,
					},
				},
				Total: 9.00,
			},
			want:                   109,
			wantRetailerPoints:     14,
			wantNoCentsPoints:      50,
			wantMultipleOf25Points: 25,
			wantItemPairsPoints:    5,
			wantDescriptionPoints:  1 + 2 + 2,
			wantOddDayPoints:       0,
			wantTimePoints:         10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("retailer points", func(t *testing.T) {
				got := tc.receipt.calculateRetailerPoints()
				if got != tc.wantRetailerPoints {
					t.Errorf("calculateRetailerPoints() = %v, expected %v", got, tc.wantRetailerPoints)
				}
			})

			t.Run("no cents points", func(t *testing.T) {
				got := tc.receipt.calculateTotalPointsForNoCents()
				if got != tc.wantNoCentsPoints {
					t.Errorf("calculateTotalPointsForNoCents() = %v, expected %v", got, tc.wantNoCentsPoints)
				}
			})

			t.Run("multiple of 0.25 points", func(t *testing.T) {
				got := tc.receipt.calculateTotalPointsForMultipleOf25()
				if got != tc.wantMultipleOf25Points {
					t.Errorf("calculateTotalPointsForMultipleOf25() = %v, expected %v", got, tc.wantMultipleOf25Points)
				}
			})

			t.Run("items pair points", func(t *testing.T) {
				got := tc.receipt.calculateTotalPointsForEveryTwoItems()
				if got != tc.wantItemPairsPoints {
					t.Errorf("calculateTotalPointsForEveryTwoItems() = %v, expected %v", got, tc.wantItemPairsPoints)
				}
			})

			t.Run("description length points", func(t *testing.T) {
				got := tc.receipt.calculatePointsForItemDescription()
				if got != tc.wantDescriptionPoints {
					t.Errorf("calculatePointsForItemDescription() = %v, expected %v", got, tc.wantDescriptionPoints)
				}
			})

			t.Run("odd day points", func(t *testing.T) {
				got := tc.receipt.calculatePointsForOddDay()
				if got != tc.wantOddDayPoints {
					t.Errorf("calculatePointsForOddDay() = %v, expected %v", got, tc.wantOddDayPoints)
				}
			})

			t.Run("time points", func(t *testing.T) {
				got := tc.receipt.calculatePointsForPurchaseTime()
				if got != tc.wantTimePoints {
					t.Errorf("calculatePointsForPurchaseTime() = %v, expected %v", got, tc.wantTimePoints)
				}
			})

			t.Run("total points", func(t *testing.T) {
				got := tc.receipt.CalculatePoints()
				if got != tc.want {
					t.Errorf("CalculatePoints() = %v, expected %v", got, tc.want)
				}
			})
		})
	}
}
