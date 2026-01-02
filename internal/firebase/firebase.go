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
	opt := option.WithCredentialsFile("grovia-mobile-app-firebase-adminsdk-fbsvc-5725680d24.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	FirebaseApp = app
}
