package main

import (
	"fmt"
	"lxc-terminal/internal/lxc"
)

func main() {
	manager := lxc.NewManager()

	// Test mock containers directly
	fmt.Println("Testing mock containers:")
	containers := []lxc.Container{
		{Name: "dev-ubuntu", Status: "running", Type: "container"},
		{Name: "test-alpine", Status: "stopped", Type: "container"},
		{Name: "web-nginx", Status: "running", Type: "container"},
	}

	fmt.Printf("Mock containers (%d):\n", len(containers))
	for _, c := range containers {
		fmt.Printf("- %s (%s, %s)\n", c.Name, c.Status, c.Type)
	}

	// Test with mock - force the getMockContainers fallback
	fmt.Println("\nTesting via ListContainers (should use mock):")
	realContainers, err := manager.ListContainers()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Returned %d containers:\n", len(realContainers))
	for _, c := range realContainers {
		fmt.Printf("- %s (%s, %s)\n", c.Name, c.Status, c.Type)
	}
}