package model

type Advertisement struct {
	Index        int     `db:"id" json:"id"`
	Name         string  `db:"product_name" json:"product_name"`
	Description  string  `db:"description" json:"description"`
	Brand        string  `db:"brand" json:"brand"`
	Category     string  `db:"category" json:"category"`
	Price        float64 `db:"price" json:"price"`
	Currency     string  `db:"currency" json:"currency"`
	Stock        int     `db:"stock" json:"stock"`
	Ean          string  `db:"ean" json:"ean"`
	Color        string  `db:"color" json:"color"`
	Size         string  `db:"size" json:"size"`
	Availability string  `db:"availability" json:"availability"`
}
