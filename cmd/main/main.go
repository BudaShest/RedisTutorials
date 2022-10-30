package main

import (
	"Redis/internal/app"
	"Redis/pkg/helpers"
)

func main() {
	var myApp *app.App = app.New()

	err := myApp.Run()
	helpers.FailOnError(err)
}
