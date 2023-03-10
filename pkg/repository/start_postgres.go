package repository

import (
	"fmt"
	"log"

	avitotask "github.com/ant0nix/avitoTask"
	"github.com/jmoiron/sqlx"
)

type StartPostgres struct {
	db *sqlx.DB
}

func NewStartPostgres(db *sqlx.DB) *StartPostgres {
	return &StartPostgres{db: db}
}

func (r *StartPostgres) CreateUser(user avitotask.User) error {
	query := fmt.Sprintf("INSERT INTO %s (balance,uname,reserved) VALUES ($1,$2,$3)", userTable)
	_, err := r.db.Exec(query, user.Balance, user.UName, user.Reserved)
	if err != nil {
		log.Printf("Error with CreateUser in repository/start_postgres. Error: %s", err.Error())
		return err
	}
	return nil
}

func (r *StartPostgres) CreateServices(service avitotask.Service) error {
	query := fmt.Sprintf("INSERT INTO %s (id,price) VALUES ($1,$2)", servicesTable)
	_, err := r.db.Exec(query, service.Id, service.Price)
	if err != nil {
		log.Printf("Error with CreateServices in repository/start_postgres. Error: %s", err.Error())
		return err
	}
	return nil
}

func (r *StartPostgres) ShowServices() ([]avitotask.Service, error) {
	var services []avitotask.Service
	query := fmt.Sprintf("SELECT * FROM %s", servicesTable)
	err := r.db.Select(&services, query)
	if err != nil {
		log.Printf("Error with ShowServices in repository/start_postgres. Error: %s", err.Error())
		return services, err
	}
	return services, nil

}
