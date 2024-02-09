package webhook

import (
	"context"
	"fmt"
	"github.com/Waxer59/DockerHook/packages/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gofiber/fiber/v2"
	"slices"
)

type queryParameters struct {
	Action string `query:"action"`
}

func Webhook(c *fiber.Ctx, cfg config.ConfigFile, cli client.Client) error {
	queryParams := new(queryParameters)
	serviceName := c.Params("service")
	ctx := context.Background()

	if queryParams.Action == "" {
		queryParams.Action = cfg.Config.DefaultAction
	}

	if err := c.QueryParser(queryParams); err != nil {
		return fiber.ErrInternalServerError
	}

	availableContainers, err := discoverContainers(cli, cfg)

	fmt.Println(availableContainers)

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
		// TODO
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
		if cfg.Config.Enable && containerLabelStatus(c, config.EnableLabel) || !cfg.Config.Enable {
			availableContainers = append(availableContainers, c)
		}
	}

	return availableContainers, nil
}

func containerLabelStatus(container types.Container, label string) bool {
	return container.Labels[label] == "true"
}
