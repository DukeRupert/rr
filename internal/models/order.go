package models

import (
	"time"
)

// Order represents an order in the Orderspace API
type Order struct {
	ID               string         `json:"id" db:"id" form:"id"`
	Number           int            `json:"number" db:"number" form:"number"`
	Created          time.Time      `json:"created" db:"created" form:"created"`
	Status           string         `json:"status" db:"status" form:"status"`
	CustomerID       string         `json:"customer_id" db:"customer_id" form:"customer_id"`
	CompanyName      string         `json:"company_name" db:"company_name" form:"company_name"`
	Phone            string         `json:"phone" db:"phone" form:"phone"`
	EmailAddresses   EmailAddresses `json:"email_addresses" db:"email_addresses" form:"email_addresses"`
	CreatedBy        string         `json:"created_by" db:"created_by" form:"created_by"`
	DeliveryDate     time.Time      `json:"delivery_date" db:"delivery_date" form:"delivery_date"`
	Reference        string         `json:"reference" db:"reference" form:"reference"`
	InternalNote     string         `json:"internal_note" db:"internal_note" form:"internal_note"`
	CustomerPONumber string         `json:"customer_po_number" db:"customer_po_number" form:"customer_po_number"`
	CustomerNote     string         `json:"customer_note" db:"customer_note" form:"customer_note"`
	StandingOrderID  *string        `json:"standing_order_id" db:"standing_order_id" form:"standing_order_id"`
	ShippingType     string         `json:"shipping_type" db:"shipping_type" form:"shipping_type"`
	ShippingAddress  Address        `json:"shipping_address" db:"shipping_address" form:"shipping_address"`
	BillingAddress   Address        `json:"billing_address" db:"billing_address" form:"billing_address"`
	OrderLines       []OrderLine    `json:"order_lines" db:"order_lines" form:"order_lines"`
	Currency         string         `json:"currency" db:"currency" form:"currency"`
	NetTotal         float64        `json:"net_total" db:"net_total" form:"net_total"`
	GrossTotal       float64        `json:"gross_total" db:"gross_total" form:"gross_total"`
}

// OrderLine represents a line item in an order
type OrderLine struct {
	ID               string            `json:"id" db:"id" form:"line_id"`
	SKU              string            `json:"sku" db:"sku" form:"line_sku"`
	Name             string            `json:"name" db:"name" form:"line_name"`
	Options          string            `json:"options" db:"options" form:"line_options"`
	GroupingCategory *GroupingCategory `json:"grouping_category,omitempty" db:"grouping_category" form:"line_grouping_category"`
	Shipping         bool              `json:"shipping" db:"shipping" form:"line_shipping"`
	Quantity         int               `json:"quantity" db:"quantity" form:"line_quantity"`
	UnitPrice        float64           `json:"unit_price" db:"unit_price" form:"line_unit_price"`
	SubTotal         float64           `json:"sub_total" db:"sub_total" form:"line_sub_total"`
	TaxRateID        string            `json:"tax_rate_id" db:"tax_rate_id" form:"line_tax_rate_id"`
	TaxName          string            `json:"tax_name" db:"tax_name" form:"line_tax_name"`
	TaxRate          float64           `json:"tax_rate" db:"tax_rate" form:"line_tax_rate"`
	TaxAmount        float64           `json:"tax_amount" db:"tax_amount" form:"line_tax_amount"`
	PreorderWindowID *string           `json:"preorder_window_id" db:"preorder_window_id" form:"line_preorder_window_id"`
	OnHold           bool              `json:"on_hold" db:"on_hold" form:"line_on_hold"`
	Invoiced         int               `json:"invoiced" db:"invoiced" form:"line_invoiced"`
	Paid             int               `json:"paid" db:"paid" form:"line_paid"`
	Dispatched       int               `json:"dispatched" db:"dispatched" form:"line_dispatched"`
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
