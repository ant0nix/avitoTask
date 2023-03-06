package repository

import (
	"errors"
	"fmt"
	"log"

	avitotask "github.com/ant0nix/avitoTask"
	"github.com/jmoiron/sqlx"
)

const (
	BadService    = "don't exist this service"
	InternalError = "internal server error (500)"
	OrderOk       = "Order created"
	OrderBad      = "Ordaer didn't create (order's price higher user balance)"
)

type DoServicesStructRepository struct {
	db *sqlx.DB
}

func NewDoServicesRepository(db *sqlx.DB) *DoServicesStructRepository {
	return &DoServicesStructRepository{db: db}
}

func (r *DoServicesStructRepository) GetServicesPrice(id int) (int, error) {
	var service avitotask.Service
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", servicesTable)
	err := r.db.Get(&service, query, id)
	if err != nil {
		log.Println(err.Error())
	}
	if avitotask.IsEmptyStruct(service) {
		return service.Price, errors.New(BadService)
	} else {
		return service.Price, nil
	}
}

func (r *DoServicesStructRepository) CreateOrder(order avitotask.Order) (string, error) {
	query := fmt.Sprintf("SELECT balance > $1 AS result FROM %s WHERE id = $2", userTable)
	var tmp []bool
	if err := r.db.Select(&tmp, query, order.Amount, order.UId); err != nil {
		log.Println(err)
	}
	if len(tmp) == 0 {
		return BadUser, errors.New(BadUser)
	} else {
		if tmp[0] == true {
			query = fmt.Sprint("CALL make_order ($1,$2)")
			_, err := r.db.Exec(query, order.SId, order.UId)
			if err != nil {
				log.Printf("Error with CreateServices in repository/start_postgres. Error: %s", err.Error())
				return InternalError, err
			}
			return OrderOk, nil
		} else {
			return OrderBad, nil
		}
	}
}

func (r *DoServicesStructRepository) DoOrder(id int) (string, error) {
	return "", nil
}
