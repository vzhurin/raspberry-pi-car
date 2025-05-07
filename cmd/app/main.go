package main

import (
	"periph.io/x/host/v3"
	"raspberry-pi-car/config"
	"raspberry-pi-car/internal/app"
)

func main() {
	_, err := host.Init()
	if err != nil {
		panic(err)
	}

	cfg, err := config.Init()
	if err != nil {
		panic(err)
	}

	app.Run(app.NewContainer(cfg))
}
