package rodney

import (
	"time"

	"github.com/bbayszczak/raspberrypi-go-drivers/l293d"
	"github.com/bbayszczak/raspberrypi-go-drivers/switchprocontroller"
	"github.com/bbayszczak/rodney/pkg/statusled"
)

const (
	runLEDPin       uint8 = 23
	bluetoothLEDPin uint8 = 24
	issueLEDPin     uint8 = 25
	motor1EnPin     uint8 = 12
	motor2EnPin     uint8 = 13
	motor1aPin      uint8 = 17
	motor1bPin      uint8 = 27
	motor2aPin      uint8 = 5
	motor2bPin      uint8 = 6
)

// Rodney represent the robot
type Rodney struct {
	runLED       *statusled.StatusLED
	bluetoothLED *statusled.StatusLED
	issueLED     *statusled.StatusLED
	rightMotor   *l293d.Motor
	leftMotor    *l293d.Motor
	controller   *switchprocontroller.SwitchProController
}

// NewRodney Creqtes a Rodney instance
func NewRodney() *Rodney {
	rodney := Rodney{
		runLED: statusled.NewStatusLED(23),
	}
	rodney.runLED.On()
	rodney.bluetoothLED.Off()
	rodney.issueLED.Off()
	return &rodney
}

// Start rodney
func (rodney *Rodney) Start() error {
	time.Sleep(time.Second)
	rodney.runLED.Off()
	time.Sleep(600 * time.Millisecond)
	return nil
}
