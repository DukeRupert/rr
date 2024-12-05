package models

// Product represents a product in the Orderspace API
type Product struct {
	ID                 string           `json:"id" db:"id" form:"id"`
	Code               string           `json:"code" db:"code" form:"code"`
	Name               string           `json:"name" db:"name" form:"name"`
	Description        string           `json:"description" db:"description" form:"description"`
	Active             bool             `json:"active" db:"active" form:"active"`
	Minimum            *int             `json:"minimum,omitempty" db:"minimum" form:"minimum"`
	TariffCode         *string          `json:"tariff_code,omitempty" db:"tariff_code" form:"tariff_code"`
	CountryOfOrigin    *string          `json:"country_of_origin,omitempty" db:"country_of_origin" form:"country_of_origin"`
	Composition        *string          `json:"composition,omitempty" db:"composition" form:"composition"`
	VariantOptions     []string         `json:"variant_options" db:"variant_options" form:"variant_options"`
	ProductVariants    []ProductVariant `json:"product_variants" db:"product_variants" form:"product_variants"`
	Categories         []Category       `json:"categories" db:"categories" form:"categories"`
	GroupingCategoryID *string          `json:"grouping_category_id,omitempty" db:"grouping_category_id" form:"grouping_category_id"`
	Images             []string         `json:"images" db:"images" form:"images"`
}

// ProductVariant represents a specific version of a product
type ProductVariant struct {
	ID              string            `json:"id" db:"id" form:"variant_id"`
	SKU             string            `json:"sku" db:"sku" form:"variant_sku"`
	Barcode         string            `json:"barcode" db:"barcode" form:"variant_barcode"`
	Options         map[string]string `json:"options" db:"options" form:"variant_options"`
	UnitPrice       float64           `json:"unit_price" db:"unit_price" form:"variant_unit_price"`
	PriceListPrices []PriceListPrice  `json:"price_list_prices" db:"price_list_prices" form:"variant_price_list_prices"`
	RRP             float64           `json:"rrp" db:"rrp" form:"variant_rrp"`
	Backorder       bool              `json:"backorder" db:"backorder" form:"variant_backorder"`
	Minimum         *int              `json:"minimum,omitempty" db:"minimum" form:"variant_minimum"`
	Multiple        *int              `json:"multiple,omitempty" db:"multiple" form:"variant_multiple"`
	Weight          float64           `json:"weight" db:"weight" form:"variant_weight"`
	TaxRateID       *string           `json:"tax_rate_id,omitempty" db:"tax_rate_id" form:"variant_tax_rate_id"`
	Location        *string           `json:"location,omitempty" db:"location" form:"variant_location"`
}

// PriceListPrice represents a specific price for a price list
type PriceListPrice struct {
	ID        string  `json:"id" db:"id" form:"price_list_id"`
	UnitPrice float64 `json:"unit_price" db:"unit_price" form:"price_list_unit_price"`
}

// Category represents a product category
type Category struct {
	ID   string `json:"id" db:"id" form:"category_id"`
	Name string `json:"name" db:"name" form:"category_name"`
}
