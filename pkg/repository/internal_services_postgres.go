package repository

import (
	"errors"
	"fmt"
	"log"

	avitotask "github.com/ant0nix/avitoTask"
	"github.com/jmoiron/sqlx"
)

const (
	changeBalanceOk  = "Balance changed successfull"
	changeBalanceBad = "Balance cannot be negative"
	BadUser          = "don't exist this user"
	NegativeAmount   = "Amount cannot be negative"
	BadRequest       = "bad request (400)"
)

type InternalServicesPostgres struct {
	db *sqlx.DB
}

func NewInternalServicesPostgres(db *sqlx.DB) *InternalServicesPostgres {
	return &InternalServicesPostgres{db: db}
}

func (r *InternalServicesPostgres) CheckUser(id int) (avitotask.User, error) {
	var user avitotask.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", userTable)
	err := r.db.Get(&user, query, id)

	if err != nil {
		log.Println(err.Error())
	}
	if avitotask.IsEmptyStruct(user) {
		return user, errors.New(BadUser)
	} else {
		return user, nil
	}
}

func (r *InternalServicesPostgres) ChangeBalance(balance avitotask.Balance, user avitotask.User) (string, error) {

	if user.Balance >= balance.ChangeBalance {
		user.Balance += balance.ChangeBalance
		query := fmt.Sprintf("UPDATE %s SET balance = $1 WHERE id = $2", userTable)
		_, err := r.db.Exec(query, user.Balance, user.Id)
		if err != nil {
			log.Printf("Error with CreateUser in repository/services_postgres. Error: %s", err.Error())
		}
		return changeBalanceOk, nil
	} else {
		return changeBalanceBad, nil
	}
}

func (r *InternalServicesPostgres) P2p(p2p avitotask.P2p) (string, error) {
	query := fmt.Sprint("CALL transaction_p2p($1,$2,$3)")
	_, err := r.db.Exec(query, p2p.SId, p2p.DId, p2p.Amount)
	if err != nil {
		log.Println(err.Error())
	}
	return "transaction have done", err
}
