package db

import (
	"context"
	"errors"
	"fmt"
	"restapi-lesson/internal/buyer"
	"restapi-lesson/internal/logging"
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

func (r *repository) Create(ctx context.Context, buyer *buyer.Buyer) error {
	q := `
		INSERT INTO public.buyer 
		    (name, surname) 
		VALUES 
		       ($1, $2) 
		RETURNING id
	`
	r.logger.Info.Println(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, buyer.Name, buyer.Surname).Scan(&buyer.ID); err != nil {
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

func (r *repository) FindAll(ctx context.Context) (u []buyer.Buyer, err error) {
	q := `
		SELECT
		    id, name, surname
		FROM
		    public.buyer
	`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	buyers := make([]buyer.Buyer, 0)

	for rows.Next() {
		var buyer buyer.Buyer

		err = rows.Scan(&buyer.ID, &buyer.Name, &buyer.Surname)
		if err != nil {
			return nil, err
		}

		buyers = append(buyers, buyer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return buyers, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (buyer.Buyer, error) {
	q := `
		SELECT
		    id, name, surname
		FROM
		    public.buyer
		WHERE id = $1
	`
	r.logger.Info.Println(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var br buyer.Buyer
	err := r.client.QueryRow(ctx, q, id).Scan(&br.ID, &br.Name, &br.Surname)
	if err != nil {
		return buyer.Buyer{}, err
	}

	return br, nil
}

func (r *repository) Update(ctx context.Context, buyer buyer.Buyer) error {
	q := `
		UPDATE 
    		public.buyer
		SET
			name = $1, surname = $2
		WHERE
		    id = $3
	`

	commandTag, err := r.client.Exec(ctx, q, buyer.Name, buyer.Surname, buyer.ID)
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
	q := `DELETE FROM public.buyer WHERE id = $1`
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

func NewRepository(client postgresql.Client, logger *logging.Logger) buyer.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
