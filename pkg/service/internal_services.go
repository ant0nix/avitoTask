package service

import (
	"log"

	avitotask "github.com/ant0nix/avitoTask"
	"github.com/ant0nix/avitoTask/pkg/repository"
)

type InternalService struct {
	repo repository.InternalService
}

func NewInternalService(repo repository.InternalService) *InternalService {
	return &InternalService{repo: repo}
}

func (s *InternalService) ChangeBalance(balance avitotask.Balance) (string, error) {

	user, err := s.repo.CheckUser(balance.UserID)
	if err != nil && err.Error() != repository.BadUser {
		log.Printf("Error service/internal_services. Error: %s", err.Error())
	}
	if err != nil && err.Error() == repository.BadUser {
		return repository.BadUser, err
	} else {
		return s.repo.ChangeBalance(balance, user)
	}
}

func (s *InternalService) ShowBalance(balance avitotask.Balance) (int, error) {

	user, err := s.repo.CheckUser(balance.UserID)
	if err != nil && err.Error() != repository.BadUser {
		log.Printf("Error service/internal_services. Error: %s", err.Error())
	}
	if err != nil && err.Error() == repository.BadUser {
		return 0, err

	} else {
		return user.Balance, err
	}
}

func (s *InternalService) P2p(p2p avitotask.P2p) (string, error) {
	return s.repo.P2p(p2p)
}
