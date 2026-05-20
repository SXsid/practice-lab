package main

import (
	"context"
	"os"

	_ "firebase.google.com/go/auth"
	firebase "firebase.google.com/go/v4"

	"google.golang.org/api/option"
)

func NewFirebaseClient() *firebase.App {
	creds := os.Getenv("FIREBASE_PRIVATE_KEY")
	if creds == "" {
		panic("env is not inserted")
	}
	opt := option.WithAuthCredentialsFile(option.CredentialsType("service_account"), "./pra.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		panic(err)
	}
	return app
}
