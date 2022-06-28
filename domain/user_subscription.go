package domain

import "time"

// UserSubscription represent unique subscription for the user
// Note that the price is inclusive of tax amount
type UserSubscription struct {
	ID          string
	Email       string
	ProductName string
	StartDate   time.Time
	EndDate     time.Time
	Price       float64
	Tax         float64
}
