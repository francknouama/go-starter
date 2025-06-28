package verify_lib_test

import (
	"context"
	"fmt"
	"log"

	"github.com/verify/lib"
)

func ExampleNew() {
	// Create a new client with default configuration
	client, err := verify-lib.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Process some input
	result, err := client.Process(context.Background(), "hello")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	// Output: Processed: hello
}

func ExampleNew_withConfig() {
	// Create a client with custom configuration
	config := &verify-lib.Config{
		Debug: false,
	}

	client, err := verify-lib.New(config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Process input
	result, err := client.Process(context.Background(), "world")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
	// Output: Processed: world
}

func ExampleClient_Process() {
	client, err := verify-lib.New(nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Process different types of input
	inputs := []string{"hello", "world", "go library"}

	for _, input := range inputs {
		result, err := client.Process(context.Background(), input)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(result)
	}
	// Output:
	// Processed: hello
	// Processed: world
	// Processed: go library
}