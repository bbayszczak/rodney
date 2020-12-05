package switchprocontroller_test

import (
	"fmt"
	"time"

	"github.com/bbayszczak/rodney/pkg/drivers/switchprocontroller"
)

func Example_nonBlocking() {
	controller := switchprocontroller.NewSwitchProController()
	controller.StartListener(0)
	for {
		// displays full controller state
		controller.Display()
		// display button A state
		fmt.Println(controller.ButtonA)
		time.Sleep(100 * time.Millisecond)
	}
}

func Example_blocking() {
	controller := switchprocontroller.NewSwitchProController()
	controller.StartListener(0)
	for {
		select {
		case <-controller.Event:
			// displays full controller state
			controller.Display()
			// display button A state
			fmt.Println(controller.ButtonA)
		}
	}
}
