package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

var FirebaseApp *firebase.App

func InitFirebase() {
	ctx := context.Background()
	opt := option.WithCredentialsFile("grovia-api-firebase-adminsdk-fbsvc-53faecaa0f.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	FirebaseApp = app
}
