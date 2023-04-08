package db

import (
	"context"
	"errors"
	"fmt"
	"restapi-lesson/internal/logging"
	"restapi-lesson/internal/prdlist"
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

func (r *repository) Create(ctx context.Context, productList *prdlist.ProductList) error {
	q := `
		INSERT INTO product_list 
		    (note_id, product_id, amount) 
		VALUES 
		       ($1, $2, $3) 
		RETURNING id
	`
	r.logger.Info.Println(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, productList.NoteID, productList.ProductID, productList.Amount).Scan(&productList.ID); err != nil {
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

func (r *repository) FindAll(ctx context.Context) ([]prdlist.ProductList, error) {
	q := `
		SELECT
		    id, note_id, product_id, amount
		FROM
		    public.product_list
	`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	productLists := make([]prdlist.ProductList, 0)

	for rows.Next() {
		var pl prdlist.ProductList

		err = rows.Scan(&pl.ID, &pl.NoteID, &pl.ProductID, &pl.Amount)
		if err != nil {
			return nil, err
		}

		productLists = append(productLists, pl)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return productLists, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (prdlist.ProductList, error) {
	q := `
		SELECT
		    id, note_id, product_id, amount
		FROM
		    public.product_list
		WHERE id = $1
	`

	r.logger.Info.Println(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var pl prdlist.ProductList
	err := r.client.QueryRow(ctx, q, id).Scan(&pl.ID, &pl.NoteID, &pl.ProductID, &pl.Amount)
	if err != nil {
		return prdlist.ProductList{}, err
	}

	return pl, nil
}

func (r *repository) Update(ctx context.Context, productList prdlist.ProductList) error {
	q := `
		UPDATE 
    		public.product_list
		SET
			note_id = $1, product_id = $2, amount = $3
		WHERE
		    id = $4
	`

	commandTag, err := r.client.Exec(ctx, q, productList.NoteID, productList.ProductID, productList.Amount, productList.ID)
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
	q := `DELETE FROM public.product_list WHERE id = $1`
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

func NewRepository(client postgresql.Client, logger *logging.Logger) prdlist.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
