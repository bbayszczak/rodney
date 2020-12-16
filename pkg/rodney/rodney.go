package rodney

import (
	"fmt"
	"time"

	"github.com/bbayszczak/raspberrypi-go-drivers/l293d"
	"github.com/bbayszczak/raspberrypi-go-drivers/switchprocontroller"
	"github.com/bbayszczak/rodney/pkg/statusled"
	log "github.com/sirupsen/logrus"
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
		runLED:       statusled.NewStatusLED(runLEDPin),
		bluetoothLED: statusled.NewStatusLED(bluetoothLEDPin),
		issueLED:     statusled.NewStatusLED(issueLEDPin),
	}
	rodney.runLED.On()
	rodney.bluetoothLED.Off()
	rodney.issueLED.Off()
	return &rodney
}

func (rodney *Rodney) initMotors() error {
	log.Info("initializing motors")
	var err error
	chip := l293d.NewL293D()
	if rodney.rightMotor, err = chip.NewMotor(1, motor1EnPin, motor1aPin, motor1bPin); err != nil {
		log.WithField("error", err).Error("impossible to initialize motor 1")
		return err
	}
	log.Info("motor 1 initialized")
	if rodney.leftMotor, err = chip.NewMotor(2, motor2EnPin, motor2aPin, motor2bPin); err != nil {
		log.WithField("error", err).Error("impossible to initialize motor 2")
		return err
	}
	log.Info("motor 2 initialized")
	log.Info("all motors successfully initialized")
	return nil
}

func (rodney *Rodney) handleFatal() {
	rodney.issueLED.On()
	rodney.runLED.Off()
	time.Sleep(500 * time.Millisecond)
}

func (rodney *Rodney) mainLoop() {
	log.Info("starting input listening")
	rodney.controller.StartListener(0)
	for {
		select {
		case <-rodney.controller.Event:
			// display button A state
			aState, _ := rodney.controller.GetButtonState("a")
			bState, _ := rodney.controller.GetButtonState("b")
			xState, _ := rodney.controller.GetButtonState("x")
			yState, _ := rodney.controller.GetButtonState("y")
			rState, _ := rodney.controller.GetButtonState("r")
			lState, _ := rodney.controller.GetButtonState("l")
			zrState, _ := rodney.controller.GetButtonState("zr")
			zlState, _ := rodney.controller.GetButtonState("zl")
			fmt.Printf(
				"A:%d B:%d X:%d Y:%d R:%d L:%d ZR:%d ZL%d\n",
				aState,
				bState,
				xState,
				yState,
				rState,
				lState,
				zrState,
				zlState,
			)
			// display left stick position
			leftStick, _ := rodney.controller.GetStick("left")
			fmt.Printf("x:%f - y:%f\n", leftStick.X, leftStick.Y)
		}
	}
}

// Start rodney
func (rodney *Rodney) Start() error {
	log.Info("I'm Rodney !")
	if err := rodney.initMotors(); err != nil {
		rodney.handleFatal()
		return err
	}
	// for now, the controller have to be paired each time
	if err := rodney.getController(); err != nil {
		rodney.handleFatal()
		return err
	}
	rodney.controller = switchprocontroller.NewSwitchProController()
	rodney.mainLoop()
	rodney.runLED.Off()
	time.Sleep(500 * time.Millisecond)
	return nil
}
