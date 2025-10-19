package main

import (
	"database/sql"
	"log"

	"github.com/IkBenJur/repetition-backend/config"
	"github.com/IkBenJur/repetition-backend/service/exercise"
	"github.com/IkBenJur/repetition-backend/service/user"
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

	userController := user.NewController(server.db)
	userHandler := user.NewHandler(userController)
	userHandler.RegisterRoutes(router)

	exerciseController := exercise.NewController(server.db)
	exerciseHandler := exercise.NewHandler(exerciseController)
	exerciseHandler.RegisterRoutes(router)

	router.Run(server.Addres)
}

func main() {
	db, err := config.ConnectDatabase(config.Envs)
	if err != nil {
		log.Fatalf("Cannot create database %v", err)
	}

	server := NewServer(":8080", db)
	
	server.Run()
}
