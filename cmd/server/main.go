package main

import (
	"fmt"

	"github.com/Francesco99975/shorehamex/cmd/boot"
	"github.com/labstack/gommon/log"
)

func main() {
	cfg, err := boot.LoadConfig()
	if err != nil {
		panic(err)
	}
	log.Infof("Starting UI server: %s", cfg.CurrentDirectory)

	e := createRouter(cfg)

	fmt.Printf("Running ShoreHamEx on port %s\n", cfg.ServerAddress())
	log.Fatal(e.Start(cfg.ServerAddress()))
}
