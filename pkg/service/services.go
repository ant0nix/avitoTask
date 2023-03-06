package service

import (
	avitotask "github.com/ant0nix/avitoTask"
	"github.com/ant0nix/avitoTask/pkg/repository"
)

type DoServicesStruct struct {
	repo repository.Service
}

func NewDoServices(repo repository.Service) *DoServicesStruct {
	return &DoServicesStruct{repo: repo}
}

func (s *DoServicesStruct) GetServicesPrice(id int) (int, error) {
	return s.repo.GetServicesPrice(id)
}

func (s *DoServicesStruct) MakeOrder(order avitotask.Order) (string, error) {
	return s.repo.CreateOrder(order)
}

func (s *DoServicesStruct) DoOrder(id int) (string, error) {
	return s.repo.DoOrder(id)
}
