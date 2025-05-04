package main

import (
	"log"
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

	pwmPinRight := gpioreg.ByName("GPIO18")
	if pwmPinRight == nil {
		log.Fatal("Failed to find GPIO18")
	}
	pwmControllerRight := pwm.NewPWM(periph.NewPin(pwmPinRight))

	controlPinRight1 := gpioreg.ByName("GPIO0")
	if controlPinRight1 == nil {
		log.Fatal("Failed to find GPIO0")
	}

	controlPinRight2 := gpioreg.ByName("GPIO1")
	if controlPinRight2 == nil {
		log.Fatal("Failed to find GPIO1")
	}

	pwmPinLeft := gpioreg.ByName("GPIO19")
	if pwmPinLeft == nil {
		log.Fatal("Failed to find GPIO19")
	}
	//pwmControllerLeft := pwm.NewPWM(periph.NewPin(pwmPinLeft))

	controlPinLeft1 := gpioreg.ByName("GPIO2")
	if controlPinLeft1 == nil {
		log.Fatal("Failed to find GPIO2")
	}

	controlPinLeft2 := gpioreg.ByName("GPIO3")
	if controlPinLeft2 == nil {
		log.Fatal("Failed to find GPIO3")
	}

	pwmControllerRight.Start(99, 50)
	//pwmControllerLeft.Start(100, 50)
	time.Sleep(15 * time.Second)
	pwmControllerRight.Stop()
	//pwmControllerLeft.Stop()

	//ch := chassis.NewChassis(
	//	pwmControllerRight,
	//	pwmControllerLeft,
	//	periph.NewPin(controlPinRight1),
	//	periph.NewPin(controlPinRight2),
	//	periph.NewPin(controlPinLeft1),
	//	periph.NewPin(controlPinLeft2),
	//)
	//
	//_ = ch.Move(100, 100, 10*time.Second)
}
