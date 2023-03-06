package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	avitotask "github.com/ant0nix/avitoTask"
	"github.com/ant0nix/avitoTask/pkg/handler"
	"github.com/ant0nix/avitoTask/pkg/repository"
	"github.com/ant0nix/avitoTask/pkg/service"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("Error with config initializing! Error:%s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error with ENV load! Error: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.DBConfig{
		Host: viper.GetString("db.host"),
		Port: viper.GetString("db.port"),
		User: viper.GetString("db.user"),
		Pass: os.Getenv("DB_PASSWORD"),
		SSL:  viper.GetString("db.ssl"),
		Name: viper.GetString("db.name"),
	})

	if err != nil {
		log.Fatalf("Error with new DB! Error: %s", err.Error())
	}
	repos := repository.NewRepository(db)
	servs := service.NewServices(repos)
	handlers := handler.NewHandler(servs)
	srv := new(avitotask.Server)
	go func() {
		if err := srv.Start(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			log.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	log.Print("Server Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Print("Server Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Fatalf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
