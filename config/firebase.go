package config

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

type ServiceAccountInfo struct {
	ProjectID string `json:"project_id"`
}

func ConnectFirebase() {
	opt := option.WithCredentialsFile("firebase-services.json")
	config := &firebase.Config{ProjectID: "presence-8f010"}
	FirebaseApp, _ = firebase.NewApp(context.Background(), config, opt)

	if FirebaseApp == nil {
		log.Println("FirebaseApp is nil")
	} else {
		log.Println("FirebaseApp initialized successfully", FirebaseApp)
	}

}
