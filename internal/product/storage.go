package product

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, product *Product) error
	FindAll(ctx context.Context) ([]Product, error)
	FindOne(ctx context.Context, id string) (Product, error)
	Update(ctx context.Context, product Product) error
	Delete(ctx context.Context, id string) error
}
