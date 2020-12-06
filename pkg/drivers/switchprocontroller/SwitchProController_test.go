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
		// display button A state
		aState, _ := controller.GetButtonState("a")
		fmt.Printf("A:%d\n", aState)
		// display left stick position
		leftStick, _ := controller.GetStick("left")
		fmt.Printf("x:%f - y:%f\n", leftStick.X, leftStick.Y)
		time.Sleep(100 * time.Millisecond)
	}
	// Output:
	// A:0
	// x:0.000000 - y:0.000000
	// [...]
}

func Example_blocking() {
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
	// Output:
	// A:0
	// x:0.000000 - y:0.000000
}
