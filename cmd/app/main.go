package main

import (
	"raspberry-pi-car/config"
	"raspberry-pi-car/internal/app"

	"periph.io/x/host/v3"
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

	err = app.Run(app.NewContainer(cfg))
	if err != nil {
		panic(err)
	}
}
