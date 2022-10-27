package repo

import (
	"strconv"
	"strings"

	"github.com/gocomerse/internal/service/product/model"
)

const (
	getProduct    = `SELECT prouct_id,product_name,product_price,product_rating from product`
	insertProduct = `Insert into product 
				   (product_id,product_name,product_price,product_rating)
				   Values($1,$2,$3,$4) 
				   RETURNING  prouct_id,product_name,product_price,product_rating 
				   `
	deleteProduct = `DELETE FROM product where id =$1`
)

func buildUpdate(product model.Product) string {
	var query string
	if product.ProductName != "" {
		query += createQuery("product_name", product.ProductName)
	}
	if product.ProductPrice != 0 {
		query += createQuery("product_price", strconv.Itoa(product.ProductPrice))

	}

	query = strings.TrimSuffix(query, ",")
	return query
}
func createQuery(tag, value string) string {
	return tag + `=` + "'" + value + "'" + ","
}
