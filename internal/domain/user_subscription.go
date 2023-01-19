package domain

import "time"

// SubscriptionStatus type to represent current subscription status
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusPaused    SubscriptionStatus = "paused"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
)

// UserSubscription represent unique subscription for the user
// Note that the price is inclusive of tax amount
type UserSubscription struct {
	ID             string
	CreatedAt      time.Time
	UpdatedAt      *time.Time
	Email          string
	ProductName    string
	StartDate      time.Time
	EndDate        time.Time
	Price          float64
	Tax            float64
	Status         SubscriptionStatus
	PauseStartDate *time.Time
}
