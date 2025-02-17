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

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// DTOs are used to handle the raw JSON input, followed by validation and conversion to proper types
// the validators help for debugging even if they are yet not sent to the user.
type ItemDTO struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

func (r ItemDTO) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ShortDescription,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[\w\s\-&]+$`)).Error("want alphanumeric characters, spaces, hyphens, and ampersands")),
		validation.Field(&r.Price,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\d+\.\d{2}$`)).Error("want 0.00 format")),
	)
}

func (r ItemDTO) ToItem() (Item, error) {
	if err := r.Validate(); err != nil {
		return Item{}, err
	}

	price, err := strconv.ParseFloat(r.Price, 64)
	if err != nil {
		return Item{}, fmt.Errorf("invalid price value: %s", r.Price)
	}

	// making an assumption here.
	if price < 0 {
		return Item{}, fmt.Errorf("price must be a positive number")
	}

	return Item{
		ShortDescription: r.ShortDescription,
		Price:            price,
	}, nil
}

type ReceiptDTO struct {
	Retailer     string    `json:"retailer"`
	PurchaseDate string    `json:"purchaseDate"`
	PurchaseTime string    `json:"purchaseTime"`
	Items        []ItemDTO `json:"items"`
	Total        string    `json:"total"`
}

func (r ReceiptDTO) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Retailer,
			validation.Required,
			validation.Match(regexp.MustCompile(`^[\w\s\-&]+$`)).Error("only alphanumeric characters, spaces, hyphens, and ampersands are allowed")),
		validation.Field(&r.PurchaseDate,
			validation.Required,
			validation.Date("2006-01-02").Error("want YYYY-MM-DD format")),
		validation.Field(&r.PurchaseTime,
			validation.Required,
			validation.Date("15:04").Error("want HH:MM format")),
		validation.Field(&r.Items,
			validation.Required,
			validation.Length(1, 0).Error("must contain at least one item")),
		validation.Field(&r.Total,
			validation.Required,
			validation.Match(regexp.MustCompile(`^\d+\.\d{2}$`)).Error("want 0.00 format")),
	)
}

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price"`
}

type Receipt struct {
	Retailer     string    `json:"retailer"`
	PurchaseDate time.Time `json:"purchaseDate"`
	PurchaseTime time.Time `json:"purchaseTime"`
	Items        []Item    `json:"items"`
	Total        float64   `json:"total"`
}

func (r ReceiptDTO) ToReceipt() (Receipt, error) {
	// these errors are unlikely to happen - and should signify some internal server error.
	purchaseDate, err := time.Parse("2006-01-02", r.PurchaseDate)
	if err != nil {
		return Receipt{}, validation.Errors{"purchaseDate": validation.NewError("purchaseDate", err.Error())}
	}

	purchaseTime, err := time.Parse("15:04", r.PurchaseTime)
	if err != nil {
		return Receipt{}, validation.Errors{"purchaseTime": validation.NewError("purchaseTime", err.Error())}
	}

	total, err := strconv.ParseFloat(r.Total, 64)
	if err != nil {
		return Receipt{}, validation.Errors{"total": validation.NewError("total", err.Error())}
	}

	// making an assumption here.
	if total < 0 {
		return Receipt{}, validation.Errors{"total": validation.NewError("total", "must be a positive number")}
	}

	items := make([]Item, len(r.Items))
	for i, itemDTO := range r.Items {
		item, err := itemDTO.ToItem()
		if err != nil {
			return Receipt{}, validation.Errors{fmt.Sprintf("items.%d", i): validation.NewError(fmt.Sprintf("items.%d", i), err.Error())}
		}
		items[i] = item
	}

	return Receipt{
		Retailer:     r.Retailer,
		PurchaseDate: purchaseDate,
		PurchaseTime: purchaseTime,
		Items:        items,
		Total:        total,
	}, nil
}

func (r *Receipt) UnmarshalJSON(b []byte) error {
	var dto ReceiptDTO
	if err := json.Unmarshal(b, &dto); err != nil {
		return err
	}

	if err := dto.Validate(); err != nil {
		return err
	}

	receipt, err := dto.ToReceipt()
	if err != nil {
		return err
	}

	*r = receipt
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
