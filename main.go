package main

import (
	"database/sql"
	"log"

	"github.com/IkBenJur/repetition-backend/config"
	controller "github.com/IkBenJur/repetition-backend/controllers/user"
	"github.com/IkBenJur/repetition-backend/routes/user"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Addres string
	db *sql.DB
}

func NewServer(address string, db *sql.DB) *Server {
	return &Server{
		Addres: address,
		db: db,
	}
}

func (server *Server) Run() {
	router := gin.Default()

	userController := controller.NewController(server.db)
	userHandler := user.NewHandler(userController)
	userHandler.RegisterRoutes(router)

	router.Run(server.Addres)
}

func main() {
	envs := config.InitConfig()
	db, err := config.ConnectDatabase(envs)
	if err != nil {
		log.Fatalf("Cannot create database %v", err)
	}

	server := NewServer(":8080", db)
	
	server.Run()
}
