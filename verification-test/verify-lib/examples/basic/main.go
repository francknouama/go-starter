package main

import (
	"context"
	"fmt"
	"log"

	verify_lib "github.com/verify/lib"
)

func main() {
	fmt.Println("verify-lib Basic Example")
	fmt.Println("==============================")

	// Create a client with default configuration
	client, err := verify_lib.New(nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Process some input
	inputs := []string{
		"hello world",
		"go library example",
		"verify-lib is awesome",
	}

	for i, input := range inputs {
		fmt.Printf("\nExample %d: Processing '%s'\n", i+1, input)
		
		result, err := client.Process(context.Background(), input)
		if err != nil {
			log.Printf("Error processing input: %v", err)
			continue
		}
		
		fmt.Printf("Result: %s\n", result)
	}

	fmt.Println("\nBasic example completed!")
}