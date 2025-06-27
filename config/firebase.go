package config

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"

	"google.golang.org/api/option"
)

// Global variable to hold the Firebase app instance
var FirebaseApp *firebase.App

// ConnectFirebase initializes the Firebase app with the provided credentials
func ConnectFirebase() {

	// Load Firebase credentials from a JSON file
	opt := option.WithCredentialsFile("firebase-services.json")

	// Initialize the Firebase app
	app, err := firebase.NewApp(context.Background(), nil, opt)

	// Check for errors during initialization
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
	}

	// Assign the initialized app to the global variable
	FirebaseApp = app

}
