package models

import (
	"time"
)

// Customer represents a customer in the Orderspace API
type Customer struct {
	ID              string         `json:"id" db:"id" form:"id"`
	CompanyName     string         `json:"company_name" db:"company_name" form:"company_name"`
	CreatedAt       time.Time      `json:"created_at" db:"created_at" form:"created_at"`
	Status          string         `json:"status" db:"status" form:"status"`
	Reference       string         `json:"reference" db:"reference" form:"reference"`
	InternalNote    string         `json:"internal_note" db:"internal_note" form:"internal_note"`
	Buyers          []Buyer        `json:"buyers" db:"buyers" form:"buyers"`
	Phone           string         `json:"phone" db:"phone" form:"phone"`
	EmailAddresses  EmailAddresses `json:"email_addresses" db:"email_addresses" form:"email_addresses"`
	TaxNumber       string         `json:"tax_number" db:"tax_number" form:"tax_number"`
	TaxRateID       *string        `json:"tax_rate_id" db:"tax_rate_id" form:"tax_rate_id"`
	Addresses       []Address      `json:"addresses" db:"addresses" form:"addresses"`
	MinimumSpend    *float64       `json:"minimum_spend" db:"minimum_spend" form:"minimum_spend"`
	PaymentTermsID  *string        `json:"payment_terms_id" db:"payment_terms_id" form:"payment_terms_id"`
	CustomerGroupID *string        `json:"customer_group_id" db:"customer_group_id" form:"customer_group_id"`
	PriceListID     *string        `json:"price_list_id" db:"price_list_id" form:"price_list_id"`

	// Additional fields for our application
	OrderInterval *int `json:"order_interval,omitempty" db:"order_interval" form:"order_interval"`
}

// Buyer represents a user that can log in and access the ordering site
type Buyer struct {
	Name         string `json:"name" db:"name" form:"buyer_name"`
	EmailAddress string `json:"email_address" db:"email_address" form:"buyer_email"`
}

// EmailAddresses represents the different email addresses for different types of communications
type EmailAddresses struct {
	Orders     string `json:"orders" db:"orders_email" form:"orders_email"`
	Dispatches string `json:"dispatches" db:"dispatches_email" form:"dispatches_email"`
	Invoices   string `json:"invoices" db:"invoices_email" form:"invoices_email"`
}

// Address represents a customer's address
type Address struct {
	CompanyName string `json:"company_name" db:"company_name" form:"address_company_name"`
	ContactName string `json:"contact_name" db:"contact_name" form:"address_contact_name"`
	Line1       string `json:"line1" db:"line1" form:"address_line1"`
	Line2       string `json:"line2" db:"line2" form:"address_line2"`
	City        string `json:"city" db:"city" form:"address_city"`
	State       string `json:"state" db:"state" form:"address_state"`
	PostalCode  string `json:"postal_code" db:"postal_code" form:"address_postal_code"`
	Country     string `json:"country" db:"country" form:"address_country"`
}

// CustomerStatus represents the possible status values for a customer
type CustomerStatus string

const (
	CustomerStatusNew    CustomerStatus = "new"
	CustomerStatusActive CustomerStatus = "active"
	CustomerStatusClosed CustomerStatus = "closed"
)

// Validate checks if a customer status is valid
func (s CustomerStatus) Validate() bool {
	switch s {
	case CustomerStatusNew, CustomerStatusActive, CustomerStatusClosed:
		return true
	default:
		return false
	}
}
