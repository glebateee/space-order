package models

import "time"

type Product struct {
	UUID        string    `json:"uid"`
	SKU         string    `json:"sku"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Currency    rune      `json:"currency"`
	BasePrice   int64     `json:"base_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
