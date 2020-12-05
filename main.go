package main

import (
	"math"
	"os"

	"github.com/bbayszczak/rodney/pkg/drivers/led"
	"github.com/bbayszczak/rodney/pkg/drivers/switchprocontroller"
	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	log.SetLevel(log.DebugLevel)
	file, err := os.OpenFile("rodney.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
		defer file.Close()
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	if rpio.Open() != nil {
		os.Exit(1)
	}
	defer rpio.Close()
	red := led.NewLED(26)
	green := led.NewLED(17)
	blue := led.NewLED(27)
	yellow := led.NewLED(22)
	white := led.NewLED(18)
	white.Dimmable()
	controller := switchprocontroller.NewSwitchProController()
	controller.StartListener(0)
	for {
		select {
		case <-controller.Event:
			switch controller.ButtonA {
			case 1:
				red.On()
			default:
				red.Off()
			}
			switch controller.ButtonB {
			case 1:
				green.On()
			default:
				green.Off()
			}
			switch controller.ButtonX {
			case 1:
				blue.On()
			default:
				blue.Off()
			}
			switch controller.ButtonY {
			case 1:
				yellow.On()
			default:
				yellow.Off()
			}
			brightnessVal := (int(math.Abs(float64(controller.StickLeft.X)) + math.Abs(float64(controller.StickLeft.Y))))
			if brightnessVal > 100 {
				brightnessVal = 100
			}
			white.SetBrightness(brightnessVal)
		}
	}
}
