package buyer

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, buyer *Buyer) error
	FindAll(ctx context.Context) (u []Buyer, err error)
	FindOne(ctx context.Context, id string) (Buyer, error)
	Update(ctx context.Context, buyer Buyer) error
	Delete(ctx context.Context, id string) error
}
