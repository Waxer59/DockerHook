package webhook

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"

	"github.com/Waxer59/DockerHook/packages/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gofiber/fiber/v2"
)

type queryParameters struct {
	Action string `query:"action"`
	Token  string `query:"token"`
}

func Webhook(c *fiber.Ctx, cfg config.ConfigFile, cli client.Client) error {
	queryParams := new(queryParameters)
	serviceName := c.Params("service")
	ctx := context.Background()

	if cfg.Auth.Enable {
		token := queryParams.Token
		registeredTokens := cfg.Auth.Tokens

		if !slices.Contains(registeredTokens, token) {
			return fiber.ErrUnauthorized
		}
	}

	if queryParams.Action == "" {
		queryParams.Action = cfg.Config.DefaultAction
	}

	if err := c.QueryParser(queryParams); err != nil {
		return fiber.ErrInternalServerError
	}

	availableContainers, err := discoverContainers(cli, cfg)

	if err != nil {
		return fiber.ErrInternalServerError
	}

	selectedContainerIdx := slices.IndexFunc(availableContainers, func(c types.Container) bool {
		return slices.Contains(c.Names, "/"+serviceName)
	})

	if selectedContainerIdx == -1 {
		return fiber.ErrNotFound
	}

	selectedContainer := availableContainers[selectedContainerIdx]

	switch queryParams.Action {
	case "start":
		err = cli.ContainerStart(ctx, selectedContainer.ID, container.StartOptions{})
	case "stop":
		err = cli.ContainerStop(ctx, selectedContainer.ID, container.StopOptions{})
	case "restart":
		err = cli.ContainerRestart(ctx, selectedContainer.ID, container.StopOptions{})
	case "pull":
		oldContainer, err := cli.ContainerInspect(ctx, selectedContainer.ID)

		if err != nil {
			fmt.Println(err.Error())
			return fiber.ErrInternalServerError
		}

		pull, err := cli.ImagePull(ctx, selectedContainer.Image, types.ImagePullOptions{})

		if err != nil {
			fmt.Println(err.Error())
			return fiber.ErrInternalServerError
		}

		defer pull.Close()

		// cli.ImagePull is asynchronous.
		// The reader needs to be read completely for the pull operation to complete.
		io.Copy(os.Stdout, pull)

		// Remove old container
		err = cli.ContainerRemove(ctx, oldContainer.ID, container.RemoveOptions{
			Force: true,
		})

		if err != nil {
			fmt.Println(err.Error())
			return fiber.ErrInternalServerError
		}

		if cfg.Config.RemoveOldImage {
			_, err = cli.ImageRemove(ctx, oldContainer.Image, types.ImageRemoveOptions{
				Force: true,
			})

			if err != nil {
				fmt.Println(err.Error())
				return fiber.ErrInternalServerError
			}
		}

		resp, err := cli.ContainerCreate(ctx, oldContainer.Config, oldContainer.HostConfig, nil, nil, oldContainer.Name)

		if err != nil {
			fmt.Println(err.Error())
			return fiber.ErrInternalServerError
		}

		if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
			panic(err)
		}
	default:
		return fiber.ErrNotFound
	}

	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.SendStatus(200)
}

func discoverContainers(cli client.Client, cfg config.ConfigFile) ([]types.Container, error) {
	var availableContainers []types.Container
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{
		All: true,
	})

	if err != nil {
		return nil, err
	}

	for _, c := range containers {
		if cfg.Auth.Enable && containerLabelStatus(c, config.EnableLabel) || !cfg.Auth.Enable {
			availableContainers = append(availableContainers, c)
		}
	}

	return availableContainers, nil
}

func containerLabelStatus(container types.Container, label string) bool {
	return container.Labels[label] == "true"
}
