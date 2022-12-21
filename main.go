package main

import (
	"context"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(NewApp())
	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}
