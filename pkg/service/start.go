package service

import (
	avitotask "github.com/ant0nix/avitoTask"
	"github.com/ant0nix/avitoTask/pkg/repository"
)

type StartService struct {
	repo repository.Start
}

func NewStartServices(repo repository.Start) *StartService {
	return &StartService{repo: repo}
}

func (s *StartService) CreateUser(user avitotask.User) error {
	return s.repo.CreateUser(user)
}

func (s *StartService) CreateServices(service avitotask.Service) error {
	return s.repo.CreateServices(service)
}

func (s *StartService) ShowServices() ([]avitotask.Service, error) {
	return s.repo.ShowServices()
}
