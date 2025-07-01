package main

import (
	"legiskuy-backend/internal/auth"
	"legiskuy-backend/internal/candidate"
	"legiskuy-backend/internal/election"
	"legiskuy-backend/internal/voter"
	"legiskuy-backend/pkg/database"
	"legiskuy-backend/pkg/middleware"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	_ "legiskuy-backend/docs"

	"github.com/gofiber/swagger"
)

// @title LegisKuy API
// @version 1.0
// @description This is the API for the LegisKuy (Pemilu) application.
// @contact.name Rozhak
// @contact.url https://github.com/RozhakXD
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:3000
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	database.ConnectDB()

	app := fiber.New()

	app.Use(logger.New())

	api := app.Group("/api")
	v1 := api.Group("/v1")

	voterRepo := voter.NewRepository()
	authRepo := auth.NewRepository()
	authService := auth.NewService(authRepo, voterRepo)
	authHandler := auth.NewHandler(authService)
	v1.Post("/register", authHandler.Register)
	v1.Post("/login", authHandler.Login)

	protected := v1.Group("/", middleware.Protected())

	candidateRepo := candidate.NewRepository()
	candidateService := candidate.NewService(candidateRepo)
	candidateHandler := candidate.NewHandler(candidateService)

	petugasOnly := protected.Group("/", middleware.RequireRole("petugas"))

	petugasOnly.Post("/candidates", candidateHandler.CreateCandidate)
	petugasOnly.Put("/candidates/:id", candidateHandler.UpdateCandidate)
	petugasOnly.Delete("/candidates/:id", candidateHandler.DeleteCandidate)

	voterService := voter.NewService(voterRepo)
	voterHandler := voter.NewHandler(voterService)

	petugasOnly.Post("/voters", voterHandler.CreateVoter)
	petugasOnly.Put("/voters/:id", voterHandler.UpdateVoter)
	petugasOnly.Delete("/voters/:id", voterHandler.DeleteVoter)

	electionRepo := election.NewRepository()

	electionService := election.NewService(electionRepo, voterRepo, candidateRepo)
	electionHandler := election.NewHandler(electionService)

	petugasOnly.Post("/election/time", electionHandler.SetElectionTime)
	petugasOnly.Post("/election/threshold", electionHandler.SetThreshold)

	protected.Get("/candidates", candidateHandler.GetAllCandidates)
	protected.Get("/candidates/:id", candidateHandler.GetCandidateByID)
	protected.Get("/voters", voterHandler.GetAllVoters)
	protected.Get("/voters/:id", voterHandler.GetVoterByID)
	protected.Get("/results", electionHandler.GetResults)
	protected.Post("/votes", electionHandler.CastVote)

	app.Get("/swagger/*", swagger.HandlerDefault)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
