package repository

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	userTable       = "users"
	servicesTable   = "services"
	accountingTable = "accounting"
	ordersTable     = "orders"
	companyTable    = "company"
)

type DBConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
	SSL  string
}

func NewPostgresDB(cfg DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Name, cfg.Pass, cfg.SSL))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	log.Println("PostgresDB started")
	return db, err
}
