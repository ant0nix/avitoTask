package repository

import (
	avitotask "github.com/ant0nix/avitoTask"
	"github.com/jmoiron/sqlx"
)

type Start interface {
	CreateUser(user avitotask.User) error
	CreateServices(service avitotask.Service) error
	ShowServices() ([]avitotask.Service, error)
}
type Service interface {
	GetServicesPrice(id int) (int, error)
	CreateOrder(order avitotask.Order) (string, error)
	DoOrder(id int) (string, error)
}

type InternalService interface {
	CheckUser(id int) (avitotask.User, error)
	ChangeBalance(balance avitotask.Balance, user avitotask.User) (string, error)
	P2p(p2p avitotask.P2p) (string, error)
}

type Repository struct {
	Start
	Service
	InternalService
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Start:           NewStartPostgres(db),
		Service:         NewDoServicesRepository(db),
		InternalService: NewInternalServicesPostgres(db),
	}
}
