package main

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// writing custom unmarshallers helps with a ) coverting to more appropriate types and b) validating the data.
// I could have also used a receiptInput struct to unmarshall stuff and then validate it, but that leads to a almost the same amount of code.
// and this feels cleaner.
// I could have converted to string{interface} but with that i have to marshal the items field back to json before passing it off or have a single
// maeshaller. this keeps both separate and does not require remarshalling.

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price"`
}

func (i *Item) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	var shortDescription string
	if err := json.Unmarshal(raw["shortDescription"], &shortDescription); err != nil {
		return fmt.Errorf("item short description is not valid json")
	}
	if matched, err := regexp.MatchString(`^[\w\s\-&]+$`, shortDescription); err != nil || !matched {
		return fmt.Errorf("invalid short description format: %s. want alphanumeric characters, spaces, hyphens, and ampersands", shortDescription)
	}

	var priceStr string
	if err := json.Unmarshal(raw["price"], &priceStr); err != nil {
		return fmt.Errorf("item price is not valid json")
	}
	if matched, err := regexp.MatchString(`^\d+\.\d{2}$`, priceStr); err != nil || !matched {
		return fmt.Errorf("invalid price format: %s. want 0.00 format", priceStr)
	}
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return fmt.Errorf("invalid price value: %s", priceStr)
	}

	i.ShortDescription = shortDescription
	i.Price = price

	return nil
}

type Receipt struct {
	Retailer     string    `json:"retailer"`
	PurchaseDate time.Time `json:"purchaseDate"`
	PurchaseTime time.Time `json:"purchaseTime"`
	Items        []Item    `json:"items"`
	Total        float64   `json:"total"`
}

func (r *Receipt) UnmarshalJSON(b []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	var retailer string
	if err := json.Unmarshal(raw["retailer"], &retailer); err != nil {
		return fmt.Errorf("retailer is not valid json")
	}
	if matched, err := regexp.MatchString(`^[\w\s\-&]+$`, retailer); err != nil || !matched {
		return fmt.Errorf("invalid retailer format: %s. want alphanumeric characters, spaces, hyphens, and ampersands", retailer)
	}

	var dateStr string
	if err := json.Unmarshal(raw["purchaseDate"], &dateStr); err != nil {
		return fmt.Errorf("purchase date is not valid json")
	}
	purchaseDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid purchase date format: %s. want YYYY-MM-DD format", dateStr)
	}

	var timeStr string
	if err := json.Unmarshal(raw["purchaseTime"], &timeStr); err != nil {
		return fmt.Errorf("purchase time is not valid json")
	}
	purchaseTime, err := time.Parse("15:04", timeStr)
	if err != nil {
		return fmt.Errorf("invalid purchase time format: %s. want HH:MM format", timeStr)
	}

	var totalStr string
	if err := json.Unmarshal(raw["total"], &totalStr); err != nil {
		return fmt.Errorf("total is not valid json")
	}
	if matched, err := regexp.MatchString(`^\d+\.\d{2}$`, totalStr); err != nil || !matched {
		return fmt.Errorf("invalid total format: %s. want 0.00 format", totalStr)
	}
	total, err := strconv.ParseFloat(totalStr, 64)
	if err != nil {
		return fmt.Errorf("invalid total value: %s", totalStr)
	}

	var items []Item
	if err := json.Unmarshal(raw["items"], &items); err != nil {
		return fmt.Errorf("invalid items: %s", err)
	}

	if len(items) == 0 {
		return fmt.Errorf("invalid items: must contain at least one item")
	}

	r.Retailer = retailer
	r.PurchaseDate = purchaseDate
	r.PurchaseTime = purchaseTime
	r.Items = items
	r.Total = total

	return nil
}

// writing these separately helps in testing them indepedently.
// making them pointer receivers helps in making less copies of the struct.
func (r *Receipt) calculateRetailerPoints() int {
	points := 0
	for _, char := range r.Retailer {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			points++
		}
	}
	return points
}

func (r *Receipt) calculateTotalPointsForNoCents() int {
	points := 0
	if r.Total == math.Trunc(r.Total) {
		points += 50
	}
	return points
}

func (r *Receipt) calculateTotalPointsForMultipleOf25() int {
	points := 0
	if r.Total/0.25 == math.Trunc(r.Total/0.25) {
		points += 25
	}
	return points
}

func (r *Receipt) calculateTotalPointsForEveryTwoItems() int {
	return len(r.Items) / 2 * 5
}

func (r *Receipt) calculatePointsForItemDescription() int {
	points := 0
	for _, item := range r.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			points += int(math.Ceil(item.Price * 0.2))
		}
	}
	return points
}

func (r *Receipt) calculatePointsForOddDay() int {
	points := 0
	if r.PurchaseDate.Day()%2 != 0 {
		points += 6
	}
	return points
}

func (r *Receipt) calculatePointsForPurchaseTime() int {
	points := 0
	if r.PurchaseTime.Hour() >= 14 && r.PurchaseTime.Hour() <= 16 {
		points += 10
	}
	return points
}

// not making the public function a pointer receiver, otherwise the users get the impression that the /can/ be modified.
func (r Receipt) CalculatePoints() int {
	points := 0
	points += r.calculateRetailerPoints()
	points += r.calculateTotalPointsForNoCents()
	points += r.calculateTotalPointsForMultipleOf25()
	points += r.calculateTotalPointsForEveryTwoItems()
	points += r.calculatePointsForItemDescription()
	points += r.calculatePointsForOddDay()
	points += r.calculatePointsForPurchaseTime()
	return points
}
