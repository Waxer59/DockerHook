package main

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/Waxer59/DockerHook/packages/config"
	"github.com/Waxer59/DockerHook/packages/webhook"
	"github.com/docker/docker/client"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	envVariables, err := config.LoadEnvVariables()

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

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		fmt.Println("Error conecting to docker cli")
		return
	}

	fmt.Println("Connected to docker cli")

	app := fiber.New(fiber.Config{
		AppName:       "DockerHook",
		CaseSensitive: true,
		GETOnly:       true,
	})

	// middlewares
	app.Use(logger.New())
	app.Use(func(c *fiber.Ctx) error {
		if !cfg.Auth.Enable {
			return c.Next()
		}

		registeredTokens := cfg.Auth.Tokens
		token := c.Query("token")
		action := c.Query("action")
		tokenAccess := cfg.GetTokenActions(token)

		if action == "" {
			action = cfg.Config.DefaultAction
		}

		if !slices.Contains(tokenAccess, action) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		if !slices.ContainsFunc(registeredTokens, func(t string) bool {
			return strings.Contains(t, token)
		}) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		return c.Next()
	})

	app.Get("/:service", func(c *fiber.Ctx) error {
		return webhook.Webhook(c, *cfg, *cli)
	})

	app.Post("/:service", func(c *fiber.Ctx) error {
		return webhook.Webhook(c, *cfg, *cli)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%s", envVariables.Port)))
}
