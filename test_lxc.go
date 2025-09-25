package main

import (
	"fmt"
	"lxc-terminal/internal/lxc"
)

func main() {
	manager := lxc.NewManager()
	containers, err := manager.ListContainers()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found %d containers:\n", len(containers))
	for _, c := range containers {
		fmt.Printf("- %s (%s, %s)\n", c.Name, c.Status, c.Type)
	}

	// Test specific container info
	if len(containers) > 0 {
		testContainer := containers[0].Name
		fmt.Printf("\nTesting container: %s\n", testContainer)

		container, err := manager.GetContainer(testContainer)
		if err != nil {
			fmt.Printf("Error getting container info: %v\n", err)
		} else {
			fmt.Printf("Container info: %+v\n", container)
		}

		running, err := manager.IsContainerRunning(testContainer)
		if err != nil {
			fmt.Printf("Error checking if running: %v\n", err)
		} else {
			fmt.Printf("Is running: %v\n", running)
		}
	}
}