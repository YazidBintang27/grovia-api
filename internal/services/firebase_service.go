package services

import (
	"context"
	"fmt"

	"grovia/internal/firebase"

	"firebase.google.com/go/v4/auth"
)

func VerifyFirebaseToken(idToken string) (*auth.Token, error) {
	ctx := context.Background()

	client, err := firebase.FirebaseApp.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting auth client: %v", err)
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("invalid firebase token: %v", err)
	}

	return token, nil
}