package routes

import (
	"context"
	"fmt"

	"github.com/Chetuks/gym-app/config"
	swagger "github.com/arsmn/fiber-swagger/v2"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/rs/zerolog/log"
)

func Start() {
	log.Info().Msg("Starting server")

	configs := config.GetConfig()
	// _ = db.GetMongoDB()

	fiberApp := setupFiberApp(configs)

	mainCtx := context.Background()
	ctx, cancel := context.WithCancel(mainCtx)
	defer cancel()
	log.Debug().Interface("ctx", ctx)
	err := fiberApp.Listen(":" + fmt.Sprint(configs.Server.Port))
	if err == nil {
		log.Error().Interface("Error starting server", err)
	}
}

func setupFiberApp(configs *config.Configurations) *fiber.App {
	log.Info().Msg("Inside setupFiberApp()")
	fiberConfig := fiber.Config{
		Prefork:               configs.Server.PreFork,
		CaseSensitive:         true,
		StrictRouting:         false,
		ServerHeader:          "",
		AppName:               "GYM-APP",
		DisableStartupMessage: false,
		ReduceMemoryUsage:     true,
		ETag:                  true, //enabled expiry tags
	}

	fiberApp := fiber.New(fiberConfig)
	// fiberApp.Use(etag.New())
	// reqLogger := logger.New(logger.Config{
	// 	Format: "[${red}${time}] - ${cyan}${ip}:${port} ${status} - ${method} ${path} ${bytesSent} ${latency}\n",
	// })
	// fiberApp.Use(recover.New())
	fiberApp.Use(compress.New(compress.Config{
		Level: compress.LevelDefault,
	}))

	fiberApp.Get("/dashboard", monitor.New())
	fiberApp.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
	fiberApp.Get("/panic", func(c *fiber.Ctx) error {
		panic("error")
	})
	fiberApp.Get("/swagger/*", swagger.HandlerDefault)
	{
		api := fiberApp.Group("/api/gym")
		{
			api.Post("/login/v1", config.GetLoginDetails)
		}
	}
	return fiberApp
}
