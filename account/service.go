package account

import(
	"context"
	"github.com/segmentio/ksuid"
)

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type accountService struct {
	repository Repository
}

func newService(r Repository) Service {

	return &accountService{r}
}

func (s *accountService) PostAccount(ctx context.Context, name string) (*Account, error) {
	a := &Account{
		Name: name,
		ID:   ksuid.New().String(),
	}
	if s.repository.PutAccount()

}

func (s *accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	s.repository.GetAccountByID()
}
func (s *accountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	s.repository.ListAccounts()
}
