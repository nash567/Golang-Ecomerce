package model

type Product struct {
	ProductID     int    `json:"product_id"`
	ProductName   string `json:"product_name"`
	ProductPrice  int    `json:"product_price"`
	ProductRating int    `json:"product_rating"`
}
