package service

import (
	avitotask "github.com/ant0nix/avitoTask"
	"github.com/ant0nix/avitoTask/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Start interface {
	CreateUser(user avitotask.User) error
	CreateServices(service avitotask.Service) error
	ShowServices() ([]avitotask.Service, error)
}

type Service interface {
	MakeOrder(order avitotask.Order) (string, error)
	DoOrder(id int) (string, error)
}

type InternalServices interface {
	ChangeBalance(balance avitotask.Balance) (string, error)
	ShowBalance(balance avitotask.Balance) (int, error)
	P2p(p2p avitotask.P2p) (string, error)
}

type Services struct {
	Service
	Start
	InternalServices
}

func NewServices(repos *repository.Repository) *Services {
	return &Services{
		Start:            NewStartServices(repos.Start),
		Service:          NewDoServices(repos.Service),
		InternalServices: NewInternalService(repos.InternalService),
	}
}
