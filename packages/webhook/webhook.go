package webhook

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/Waxer59/DockerHook/packages/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gofiber/fiber/v2"
)

type queryParameters struct {
	Action  string `query:"action"`
	Version string `query:"version"`
}

func Webhook(c *fiber.Ctx, cfg config.ConfigFile, cli client.Client) error {
	queryParams := new(queryParameters)
	serviceName := c.Params("service")
	ctx := context.Background()

	if queryParams.Action == "" {
		queryParams.Action = cfg.Config.DefaultAction
	}

	if err := c.QueryParser(queryParams); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	availableContainers, err := discoverContainers(cli, cfg)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	selectedContainerIdx := slices.IndexFunc(availableContainers, func(c types.Container) bool {
		return slices.Contains(c.Names, "/"+serviceName)
	})

	if selectedContainerIdx == -1 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	selectedContainer := availableContainers[selectedContainerIdx]

	switch queryParams.Action {
	case "start":
		fmt.Println("Starting container: " + selectedContainer.ID)
		err = cli.ContainerStart(ctx, selectedContainer.ID, container.StartOptions{})
	case "stop":
		fmt.Println("Stopping container: " + selectedContainer.ID)
		err = cli.ContainerStop(ctx, selectedContainer.ID, container.StopOptions{})
	case "restart":
		fmt.Println("Restarting container: " + selectedContainer.ID)
		err = cli.ContainerRestart(ctx, selectedContainer.ID, container.StopOptions{})
	case "pull":
		oldContainer, err := cli.ContainerInspect(ctx, selectedContainer.ID)

		if err != nil {
			fmt.Println(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if queryParams.Version == "" {
			queryParams.Version = "latest"
		}

		imageNameWithoutVersion := strings.Split(selectedContainer.Image, ":")[0]
		selectedContainer.Image = fmt.Sprintf("%s:%s", imageNameWithoutVersion, queryParams.Version)

		fmt.Println("Pulling image: " + selectedContainer.Image)

		pull, err := cli.ImagePull(ctx, selectedContainer.Image, types.ImagePullOptions{})

		if err != nil {
			fmt.Println(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
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
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		fmt.Println("Creating new container with image: " + selectedContainer.Image)

		// Update container with new image
		oldContainer.Config.Image = selectedContainer.Image

		// Create container with the updated image
		resp, err := cli.ContainerCreate(ctx, oldContainer.Config, oldContainer.HostConfig, nil, nil, oldContainer.Name)

		if err != nil {
			fmt.Println(err.Error())
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
			panic(err)
		}
	default:
		return c.SendStatus(fiber.StatusNotFound)
	}

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
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
		if cfg.Config.LabelBased && containerLabelStatus(c, config.EnableLabel) || !cfg.Config.LabelBased {
			availableContainers = append(availableContainers, c)
		}
	}

	return availableContainers, nil
}

func containerLabelStatus(container types.Container, label string) bool {
	return container.Labels[label] == "true"
}
