package main

import (
	"log"
	"time"

	"github.com/IkBenJur/repetition-backend/config"
	"github.com/IkBenJur/repetition-backend/service/userWorkout"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Addres string
	db     *pgxpool.Pool
}

func NewServer(address string, db *pgxpool.Pool) *Server {
	return &Server{
		Addres: address,
		db:     db,
	}
}

func (server *Server) Run() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.Envs.FrontEndUrl},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// userController := user.NewController(server.db)
	// userHandler := user.NewHandler(userController)
	// userHandler.RegisterRoutes(router)

	// exerciseController := exercise.NewController(server.db)
	// exerciseHandler := exercise.NewHandler(exerciseController)
	// exerciseHandler.RegisterRoutes(router)

	// userWorkoutController := userWorkout.NewController(server.db)
	newController := userWorkout.NewUserWorkoutController(server.db)
	userWorkoutHandler := userWorkout.NewHandler(*newController)
	userWorkoutHandler.RegisterRoutes(router)

	// userWorkoutExerciseController := userWorkoutExercise.NewController(server.db)
	// userWorkoutExerciseHandler := userWorkoutExercise.NewHandler(*userWorkoutExerciseController, userController, *userWorkoutController)
	// userWorkoutExerciseHandler.RegisterRoutes(router)

	// userWorkoutExerciseSetController := userWorkoutExerciseSet.NewController(server.db)
	// userWorkoutExerciseSetHandler := userWorkoutExerciseSet.NewHandler(*userWorkoutExerciseSetController, userController, *userWorkoutExerciseController)
	// userWorkoutExerciseSetHandler.RegisterRoutes(router)

	// templateWorkoutController := workouttemplate.NewController(server.db)
	// templateWorkoutHandler := workouttemplate.NewHandler(*templateWorkoutController)
	// templateWorkoutHandler.RegisterRoutes(router)

	router.Run(server.Addres)
}

func main() {
	db, err := config.ConnectDatabase(config.Envs)
	if err != nil {
		log.Fatalf("Cannot create database %v", err)
		return
	}

	server := NewServer(":8080", db)

	server.Run()
}
