package models

import (
	"time"
)

type Order struct {
	ID               string         `json:"id"`
	Number           int            `json:"number"`
	Created          time.Time      `json:"created"`
	Status           string         `json:"status"`
	CustomerID       string         `json:"customer_id"`
	CompanyName      string         `json:"company_name"`
	Phone            string         `json:"phone"`
	EmailAddresses   EmailAddresses `json:"email_addresses"`
	CreatedBy        string         `json:"created_by"`
	DeliveryDate     string         `json:"delivery_date"`
	Reference        string         `json:"reference"`
	InternalNote     string         `json:"internal_note"`
	CustomerPONumber string         `json:"customer_po_number"`
	CustomerNote     string         `json:"customer_note"`
	StandingOrderID  *string        `json:"standing_order_id"`
	ShippingType     string         `json:"shipping_type"`
	ShippingAddress  Address        `json:"shipping_address"`
	BillingAddress   Address        `json:"billing_address"`
	OrderLines       []OrderLine    `json:"order_lines"`
	Currency         string         `json:"currency"`
	NetTotal         float64        `json:"net_total"`
	GrossTotal       float64        `json:"gross_total"`
}

type OrderLine struct {
	ID               string           `json:"id"`
	SKU              string           `json:"sku"`
	Name             string           `json:"name"`
	Options          string           `json:"options"`
	GroupingCategory GroupingCategory `json:"grouping_category"`
	Shipping         bool             `json:"shipping"`
	Quantity         int              `json:"quantity"`
	UnitPrice        float64          `json:"unit_price"`
	SubTotal         float64          `json:"sub_total"`
	TaxRateID        string           `json:"tax_rate_id"`
	TaxName          string           `json:"tax_name"`
	TaxRate          float64          `json:"tax_rate"`
	TaxAmount        float64          `json:"tax_amount"`
	PreorderWindowID string           `json:"preorder_window_id"`
	OnHold           bool             `json:"on_hold"`
	Invoiced         int              `json:"invoiced"`
	Paid             int              `json:"paid"`
	Dispatched       int              `json:"dispatched"`
}

// GroupingCategory represents a category for grouping products
type GroupingCategory struct {
	ID   string `json:"id" db:"id" form:"category_id"`
	Name string `json:"name" db:"name" form:"category_name"`
}

// OrderStatus represents the possible status values for an order
type OrderStatus string

const (
	OrderStatusNew           OrderStatus = "new"
	OrderStatusInvoiced      OrderStatus = "invoiced"
	OrderStatusReleased      OrderStatus = "released"
	OrderStatusPartFulfilled OrderStatus = "part_fulfilled"
	OrderStatusPreorder      OrderStatus = "preorder"
	OrderStatusFulfilled     OrderStatus = "fulfilled"
	OrderStatusStandingOrder OrderStatus = "standing_order"
	OrderStatusCancelled     OrderStatus = "cancelled"
)

// Validate checks if an order status is valid
func (s OrderStatus) Validate() bool {
	switch s {
	case OrderStatusNew, OrderStatusInvoiced, OrderStatusReleased,
		OrderStatusPartFulfilled, OrderStatusPreorder, OrderStatusFulfilled,
		OrderStatusStandingOrder, OrderStatusCancelled:
		return true
	default:
		return false
	}
}
