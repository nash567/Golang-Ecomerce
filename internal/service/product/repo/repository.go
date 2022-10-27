package repo

import (
	"context"
	"database/sql"
	"fmt"

	logModel "github.com/gocomerse/internal/logger/model"
	"github.com/gocomerse/internal/service/product/model"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}
func (r *Repository) Get(ctx context.Context, log logModel.Logger) ([]*model.Product, error) {

	var Products []*model.Product

	stmt, err := r.db.PrepareContext(ctx, getProduct)
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)

	if err != nil {
		log.WithError(err).Error("failed to get products")
		return nil, fmt.Errorf("failed to get poducts: %w", err)
	}
	for rows.Next() {
		var product model.Product
		err = rows.Scan(&product.ProductID, &product.ProductName, &product.ProductPrice, &product.ProductRating)

		if err != nil {
			log.WithError(err).Error("failed to get Products:")

			return nil, fmt.Errorf("failed to get Products: %w", err)
		}
		Products = append(Products, &product)
	}
	if err = rows.Err(); err != nil {
		log.WithError(err).Error("failed to get all the Products")

		return Products, fmt.Errorf("failed to get all Products: %w", err)
	}
	if len(Products) == 0 {
		log.WithError(err).Error("no Product found")

		return nil, fmt.Errorf("%w", model.ErrNoRecordFound)
	}
	defer rows.Close()

	return Products, nil
}

func (r *Repository) Create(ctx context.Context, log logModel.Logger, product model.Product) (*model.Product, error) {
	var Product model.Product
	stmt, err := r.db.PrepareContext(ctx, insertProduct)
	if err != nil {
		log.WithError(err).Error("failed to prepare context with query")
		return nil, fmt.Errorf("failed to prepare context for Addiying product: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, Product.ProductID, Product.ProductName, Product.ProductPrice, Product.ProductRating)
	err = row.Scan(product.ProductID, &product.ProductName, &product.ProductPrice, &product.ProductRating)
	if err != nil {
		log.WithError(err).Error("failed to scan while adding product")
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &Product, nil
}

func (r *Repository) Update(ctx context.Context, log logModel.Logger, product model.Product) (*model.Product, error) {
	var Product model.Product
	query := `Update "Product" set ` + buildUpdate(Product) + ` where id=$1 RETURNING prouct_id,product_name,product_price,product_rating `
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.WithError(err).Error("failed to prepare context for update Product")
		return nil, fmt.Errorf("failed to prepare context for update Product: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, product.ProductID)
	err = row.Scan(&Product.ProductID, &Product.ProductName, &Product.ProductPrice, &Product.ProductRating)
	if err != nil {
		log.WithError(err).Error("failed to scan while adding Product")
		return nil, fmt.Errorf("failed to create Product: %w", err)
	}

	return &Product, nil
}

func (r *Repository) Delete(ctx context.Context, log logModel.Logger, id int) error {
	res, err := r.db.ExecContext(ctx, deleteProduct, id)

	if err != nil {
		log.WithError(err).Error("failed to delete product with id")
		return fmt.Errorf("failed to delete product :%w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.WithError(err).Error("failed to delete product with id")
		return fmt.Errorf("failed to delete product :%w", err)
	}

	if count == 1 {
		return nil
	} else if count >= 0 {
		return fmt.Errorf("error occurred: %w", err)
	}
	return nil
}
