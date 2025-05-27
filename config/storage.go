package config

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var JSONCreds string
func ConnectStorage() {

	ctx := context.Background()
	data, err := os.ReadFile("storage-accessor.json")
	if err != nil {
		panic(fmt.Sprintf("failed to read service account file: %v", err))
	}
	JSONCreds = string(data)
	client, err := storage.NewClient(ctx, option.WithCredentialsJSON([]byte(JSONCreds)))
	if err != nil {
		panic(fmt.Sprintf("failed to create storage client: %v", err))
	}
	defer client.Close()
}
