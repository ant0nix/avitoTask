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

func (r *DoServicesStructRepository) CreateOrder(order avitotask.Order) (string, error) {
	query := fmt.Sprintf("SELECT balance > (SELECT price FROM %s WHERE id = $1) AS result FROM %s WHERE id = $2", servicesTable, userTable)
	var tmp []bool
	if err := r.db.Select(&tmp, query, order.SId, order.UId); err != nil {
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

func (r *DoServicesStructRepository) DoOrders(id int) (string, error) {
	query := fmt.Sprint("CALL do_order ($1)")
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Printf("Error with CreateServices in repository/start_postgres. Error: %s", err.Error())
		return InternalError, err
	}
	return "order sent", nil
}
