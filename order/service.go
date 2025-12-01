package order

import (
	"context"
	"time"

	"github.com/segmentio/ksuid"
)

type service interface {
	PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error)
	GetOrdersFromAccount(ctx context.Context, accountID string) ([]Order, error)
}

type Order struct {
	ID         string
	CreatedAt  time.Time
	TotalPrice float64
	AccountID  string
	Products   []OrderedProduct
}
type OrderedProduct struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    uint32
}
type orderService struct {
	repository Repository
}

func NewService(r Repository) {
	return &orderService{r}
}

func (s orderService) PostOrder(ctx context.Context, accountID string, products []OrderedProduct) (*Order, error) {
	o := &Order{
		ID:        ksuid.New().String(),
		CreatedAt: time.Now().UTC(),
		AccountID: accountID,
		Products:  products,
	}
	o.TotalPrice := 0.0
	for _, p := range products {
		o.Total += p.Price * float64(p.Qunatity)
	}
	err := s.repository.PutOrder(ctx, *o)
	if err != nil {
		return nil, err
	}

	return o, nil
}
func (s orderService) GetOrdersFromAccount(ctx context.Context, accountID string) ([]Order, error) {
	return s.repository.GetOrdersFromAccount(ctx, accountID)
}
