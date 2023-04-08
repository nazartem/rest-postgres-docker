package db

import (
	"context"
	"errors"
	"fmt"
	"restapi-lesson/internal/logging"
	"restapi-lesson/internal/product"
	"restapi-lesson/pkg/client/postgresql"
	"strings"

	"github.com/jackc/pgconn"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) Create(ctx context.Context, product *product.Product) error {
	q := `
		INSERT INTO product 
		    (name, description, price, amount) 
		VALUES 
		       ($1, $2, $3, $4) 
		RETURNING id
	`
	r.logger.Info.Println(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, product.Name, product.Description, product.Price, product.Amount).Scan(&product.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Err.Println(newErr)
			return newErr
		}
		return err
	}

	return nil
}

func (r *repository) FindAll(ctx context.Context) ([]product.Product, error) {
	q := `
		SELECT
		    id, name, description, price, amount
		FROM
		    public.product
	`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	products := make([]product.Product, 0)

	for rows.Next() {
		var prd product.Product

		err = rows.Scan(&prd.ID, &prd.Name, &prd.Description, &prd.Price, &prd.Amount)
		if err != nil {
			return nil, err
		}

		products = append(products, prd)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (product.Product, error) {
	q := `
		SELECT
		    id, name, description, price, amount
		FROM
		    public.product
		WHERE id = $1
	`

	r.logger.Info.Println(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var prd product.Product
	err := r.client.QueryRow(ctx, q, id).Scan(&prd.ID, &prd.Name, &prd.Description, &prd.Price, &prd.Amount)
	if err != nil {
		return product.Product{}, err
	}

	return prd, nil
}

func (r *repository) Update(ctx context.Context, product product.Product) error {
	q := `
		UPDATE 
    		public.product
		SET
			name = $1, description = $2, price = $3, amount = $4
		WHERE
		    id = $5
	`

	commandTag, err := r.client.Exec(ctx, q, product.Name, product.Description, product.Price, product.Amount, product.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Err.Println(newErr)
			return newErr
		}
		return err
	}
	if commandTag.RowsAffected() != 1 {
		newErr := errors.New("no row found to update")
		r.logger.Err.Println(newErr)
		return newErr
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM product WHERE id = $1`
	commandTag, err := r.client.Exec(ctx, q, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Err.Println(newErr)
			return newErr
		}
		return err
	}

	if commandTag.RowsAffected() != 1 {
		newErr := errors.New("no row found to delete")
		r.logger.Err.Println(newErr)
		return newErr
	}

	return nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) product.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
