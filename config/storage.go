package config

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Global variable to hold the JSON credentials for Google Cloud Storage
var JSONCreds string
func ConnectStorage() {
	// Load the service account JSON credentials from a file
	ctx := context.Background()
	data, err := os.ReadFile("storage-accessor.json")

	// Check for errors while reading the file
	if err != nil {
		panic(fmt.Sprintf("failed to read service account file: %v", err))
	}

	// Convert the byte slice to a string and assign it to JSONCreds
	JSONCreds = string(data)
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(JSONCreds)))

	// Check for errors while creating the storage client
	if err != nil {
		panic(fmt.Sprintf("failed to create storage client: %v", err))
	}

	// Close the client when done
	defer client.Close()
}
