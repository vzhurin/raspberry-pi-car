package main

import (
	"fmt"
	"raspberry-pi-car/config"
	"raspberry-pi-car/internal/app"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		fmt.Println(err) // TODO logger
	}

	app.Run(app.NewContainer(cfg))
}
