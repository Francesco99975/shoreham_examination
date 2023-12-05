package main

import (
	"fmt"
	"os"

	"github.com/Francesco99975/shorehamex/cmd/boot"
	"github.com/Francesco99975/shorehamex/internal/models"
)

func main() {
	err := boot.LoadEnvVariables()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")

	models.Setup(os.Getenv("DSN"))

	e := createRouter()

	fmt.Printf("Running ShoreHamEx on port %s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
