package main

import (
	"fmt"
	"os"

	"github.com/bbayszczak/rodney/pkg/drivers/switchprocontroller"
	log "github.com/sirupsen/logrus"
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
	controller := switchprocontroller.NewSwitchProController()
	controller.StartListener(0)
	for {
		select {
		case <-controller.Event:
			// display button A state
			aState, _ := controller.GetButtonState("a")
			fmt.Printf("A:%d\n", aState)
			// display left stick position
			leftStick, _ := controller.GetStick("left")
			fmt.Printf("x:%f - y:%f\n", leftStick.X, leftStick.Y)
		}
	}
}
