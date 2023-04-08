package db

import (
	"context"
	"errors"
	"fmt"
	"restapi-lesson/internal/logging"
	"restapi-lesson/internal/note"
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

func (r *repository) Create(ctx context.Context, note *note.Note) error {
	q := `
		INSERT INTO public.note 
		    (date, buyer_id) 
		VALUES 
		       ($1, $2) 
		RETURNING number
	`
	r.logger.Info.Println(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, note.Date, note.BuyerID).Scan(&note.Number); err != nil {
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

func (r *repository) FindAll(ctx context.Context) ([]note.NoteWithPrdList, error) {
	qNote := `
		SELECT
		    number, date, buyer_id
		FROM
		    public.note
	`

	rowsQNote, err := r.client.Query(ctx, qNote)
	if err != nil {
		return nil, err
	}

	notes := make([]note.NoteWithPrdList, 0)

	for rowsQNote.Next() {
		var nt note.NoteWithPrdList

		err = rowsQNote.Scan(&nt.Number, &nt.Date, &nt.BuyerID)
		if err != nil {
			return nil, err
		}

		qPrdList := `
		SELECT
    		product.name, product.price, product_list.amount
		FROM
    		product_list INNER JOIN product ON product_list.product_id = product.id
		WHERE note_id = $1;
	`

		rows, err := r.client.Query(ctx, qPrdList, nt.Number)
		if err != nil {
			return nil, err
		}

		lists := make([]note.PrdList, 0)

		for rows.Next() {
			var list note.PrdList

			err = rows.Scan(&list.Name, &list.Price, &list.Amount)
			if err != nil {
				return nil, err
			}

			list.TotalCount = list.Price * float64(list.Amount)
			lists = append(lists, list)
		}

		if err = rows.Err(); err != nil {
			return nil, err
		}

		nt.PrdLists = lists

		notes = append(notes, nt)
	}

	if err = rowsQNote.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *repository) FindOne(ctx context.Context, number string) (note.NoteWithPrdList, error) {
	qNote := `
		SELECT
		    number, date, buyer_id
		FROM
		    public.note
		WHERE number = $1
	`
	r.logger.Info.Println(fmt.Sprintf("SQL Query: %s", formatQuery(qNote)))

	var nt note.NoteWithPrdList
	err := r.client.QueryRow(ctx, qNote, number).Scan(&nt.Number, &nt.Date, &nt.BuyerID)
	if err != nil {
		return note.NoteWithPrdList{}, err
	}

	qPrdList := `
		SELECT
    		product.name, product.price, product_list.amount
		FROM
    		product_list INNER JOIN product ON product_list.product_id = product.id
		WHERE note_id = $1;
	`

	rows, err := r.client.Query(ctx, qPrdList, nt.Number)
	if err != nil {
		return note.NoteWithPrdList{}, err
	}

	lists := make([]note.PrdList, 0)

	for rows.Next() {
		var list note.PrdList

		err = rows.Scan(&list.Name, &list.Price, &list.Amount)
		if err != nil {
			return note.NoteWithPrdList{}, err
		}

		list.TotalCount = list.Price * float64(list.Amount)
		lists = append(lists, list)
	}

	if err = rows.Err(); err != nil {
		return note.NoteWithPrdList{}, err
	}

	nt.PrdLists = lists

	return nt, nil
}

func (r *repository) Update(ctx context.Context, note note.Note) error {
	q := `
		UPDATE 
    		public.note
		SET
			date = $1, buyer_id = $2
		WHERE
		    number = $3
	`

	commandTag, err := r.client.Exec(ctx, q, note.Date, note.BuyerID, note.Number)
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

func (r *repository) Delete(ctx context.Context, number string) error {
	q := `DELETE FROM public.note WHERE number = $1`
	commandTag, err := r.client.Exec(ctx, q, number)
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

func NewRepository(client postgresql.Client, logger *logging.Logger) note.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
