package main

import (
	"context"
	"fmt"
	"github.com/Waxer59/DockerHook/packages/config"
	"github.com/Waxer59/DockerHook/packages/webhook"
	"github.com/docker/docker/client"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

func main() {
	envVariables, err := config.LoadEnvVariables()
	fmt.Println(envVariables)
	if err != nil {
		log.Fatal(err)
		return
	}

	cfg, err := config.LoadConfig(envVariables.ConfigPath)

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Conecting to docker cli...")

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)

	if err != nil {
		fmt.Println("Error conecting to docker cli")
		return
	}

	fmt.Println("Connected to docker cli")

	cli.NegotiateAPIVersion(ctx)

	app := fiber.New(fiber.Config{
		AppName:       "DockerHook",
		CaseSensitive: true,
		GETOnly:       true,
	})

	app.Use(logger.New())

	app.Get("/:service", func(c *fiber.Ctx) error {
		return webhook.Webhook(c, *cfg, *cli)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%s", envVariables.Port)))
}
