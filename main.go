package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/IkBenJur/repetition-backend/config"
	userWorkoutExercise "github.com/IkBenJur/repetition-backend/service/UserWorkoutExercise"
	userWorkoutExerciseSet "github.com/IkBenJur/repetition-backend/service/UserWorkoutExerciseSet"
	"github.com/IkBenJur/repetition-backend/service/exercise"
	"github.com/IkBenJur/repetition-backend/service/user"
	"github.com/IkBenJur/repetition-backend/service/userWorkout"
	workouttemplate "github.com/IkBenJur/repetition-backend/service/workoutTemplate"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Addres string
	db     *sql.DB
}

func NewServer(address string, db *sql.DB) *Server {
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

	userController := user.NewController(server.db)
	userHandler := user.NewHandler(userController)
	userHandler.RegisterRoutes(router)

	exerciseController := exercise.NewController(server.db)
	exerciseHandler := exercise.NewHandler(exerciseController)
	exerciseHandler.RegisterRoutes(router)

	userWorkoutController := userWorkout.NewController(server.db)
	userWorkoutHandler := userWorkout.NewHandler(*userWorkoutController, userController)
	userWorkoutHandler.RegisterRoutes(router)

	userWorkoutExerciseController := userWorkoutExercise.NewController(server.db)
	userWorkoutExerciseHandler := userWorkoutExercise.NewHandler(*userWorkoutExerciseController, userController, *userWorkoutController)
	userWorkoutExerciseHandler.RegisterRoutes(router)

	userWorkoutExerciseSetController := userWorkoutExerciseSet.NewController(server.db)
	userWorkoutExerciseSetHandler := userWorkoutExerciseSet.NewHandler(*userWorkoutExerciseSetController, userController, *userWorkoutExerciseController)
	userWorkoutExerciseSetHandler.RegisterRoutes(router)

	templateWorkoutController := workouttemplate.NewController(server.db)
	templateWorkoutHandler := workouttemplate.NewHandler(*templateWorkoutController)
	templateWorkoutHandler.RegisterRoutes(router)

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
