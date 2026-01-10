package services

import (
	"context"

	"grovia/internal/firebase"
	"grovia/pkg"

	"firebase.google.com/go/v4/auth"
)

func VerifyFirebaseToken(idToken string) (*auth.Token, error) {
	ctx := context.Background()

	client, err := firebase.FirebaseApp.Auth(ctx)
	if err != nil {
		return nil, pkg.NewInternalServerError("Error")
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, pkg.NewBadRequestError("Firebase token invalid")
	}

	return token, nil
}
