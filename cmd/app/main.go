package main

import (
	"log"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
	"raspberry-pi-car/internal/periph"
	"raspberry-pi-car/internal/pwm"
	"time"
)

func main() {
	_, err := host.Init()
	if err != nil {
		log.Fatal(err)
	}

	pwmPin := gpioreg.ByName("GPIO19")
	if pwmPin == nil {
		log.Fatal("Failed to find GPIO19")
	}

	pwmController := pwm.NewPWM(periph.NewPin(pwmPin))

	p2 := gpioreg.ByName("GPIO2")
	if p2 == nil {
		log.Fatal("Failed to find GPIO2")
	}

	if err := p2.Out(gpio.High); err != nil {
		log.Fatal(err)
	}

	p3 := gpioreg.ByName("GPIO3")
	if p3 == nil {
		log.Fatal("Failed to find GPIO3")
	}

	if err := p3.Out(gpio.Low); err != nil {
		log.Fatal(err)
	}

	err = pwmController.Start(5, 50)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	err = pwmController.Start(10, 50)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(5 * time.Second)

	pwmController.Stop()
}
