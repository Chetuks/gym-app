package main

import (
	"github.com/Chetuks/gym-app/routes"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	log "github.com/rs/zerolog/log"
)

func main() {
	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault)
	log.Info().Msg("Starting gym-app service")
	routes.Start()
}
