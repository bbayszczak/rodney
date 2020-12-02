package switchprocontroller_test

import (
	"time"

	"github.com/bbayszczak/rodney/pkg/drivers/switchprocontroller"
)

func Example() {
	controller := switchprocontroller.NewSwitchProController()
	controller.StartListener(0)
	for {
		controller.Display()
		time.Sleep(100 * time.Millisecond)
	}
}
