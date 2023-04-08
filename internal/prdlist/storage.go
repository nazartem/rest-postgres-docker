package prdlist

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, productList *ProductList) error
	FindAll(ctx context.Context) ([]ProductList, error)
	FindOne(ctx context.Context, id string) (ProductList, error)
	Update(ctx context.Context, productList ProductList) error
	Delete(ctx context.Context, id string) error
}
