// Package verify_lib provides verify-lib functionality.
//
// This package offers a clean and simple API for verify-lib operations.
// It includes structured logging, comprehensive error handling, and
// production-ready features.
//
// Basic usage:
//
//	client, err := verify-lib.New(nil)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
//
//	result, err := client.Process(context.Background(), "input")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// For advanced usage with custom configuration:
//
//	config := &verify-lib.Config{
//		Debug: true,
//	}
//	config.Logger.Level = "debug"
//
//	client, err := verify-lib.New(config)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer client.Close()
//
// See the examples directory for more detailed usage examples.
package verify_lib