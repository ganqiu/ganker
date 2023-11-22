package main

import (
	"go_docker_learning/ganker/cmd"
	"go_docker_learning/ganker/container"
)

func main() {
	cmd.Execute()

	container.CreateRootDir()
}
