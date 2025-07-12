package model

type Advertisement struct {
	Index        int    `db:"id" json:"id"`
	Name         string `db:"product_name" json:"name"`
	Description  string
	Brand        string
	Category     string
	Price        float64
	Currency     string
	Stock        int
	Ean          string
	Color        string
	Size         string
	Availability string
}
