package note

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, note *Note) error
	FindAll(ctx context.Context) ([]NoteWithPrdList, error)
	FindOne(ctx context.Context, id string) (NoteWithPrdList, error)
	Update(ctx context.Context, note Note) error
	Delete(ctx context.Context, id string) error
}
