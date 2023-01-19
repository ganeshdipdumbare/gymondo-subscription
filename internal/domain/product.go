package domain

// Product represents the product details
// ID unique product id
// Name product name
// SubscriptionPeriod is a period in terms of months
// Price is total price including tax
// TaxPercentage is % of tax applied on base price to get Price
type Product struct {
	ID                 string
	Name               string
	SubscriptionPeriod uint
	Price              float64
	TaxPercentage      float64
}
