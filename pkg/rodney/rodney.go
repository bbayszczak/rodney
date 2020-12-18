package rodney

import (
	"math"
	"time"

	"github.com/bbayszczak/raspberrypi-go-drivers/fs90r"
	"github.com/bbayszczak/rodney/pkg/statusled"
	"github.com/raspberrypi-go-drivers/switchprocontroller"
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
	// rightMotor   *l293d.Motor
	// leftMotor    *l293d.Motor
	rightMotor *fs90r.FS90R
	leftMotor  *fs90r.FS90R
	controller *switchprocontroller.SwitchProController
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

// func (rodney *Rodney) initMotors() error {
// 	log.Info("initializing motors")
// 	var err error
// 	chip := l293d.NewL293D()
// 	if rodney.rightMotor, err = chip.NewMotor(1, motor1EnPin, motor1aPin, motor1bPin); err != nil {
// 		log.WithField("error", err).Error("impossible to initialize motor 1")
// 		return err
// 	}
// 	log.Info("motor 1 initialized")
// 	if rodney.leftMotor, err = chip.NewMotor(2, motor2EnPin, motor2aPin, motor2bPin); err != nil {
// 		log.WithField("error", err).Error("impossible to initialize motor 2")
// 		return err
// 	}
// 	log.Info("motor 2 initialized")
// 	log.Info("all motors successfully initialized")
// 	return nil
// }

func (rodney *Rodney) initServos() error {
	log.Info("initializing motors")
	rodney.rightMotor = fs90r.NewFS90R(12)
	rodney.leftMotor = fs90r.NewFS90R(13)
	log.Info("all motors successfully initialized")
	return nil
}

func (rodney *Rodney) handleFatal() {
	rodney.issueLED.On()
	rodney.runLED.Off()
	time.Sleep(500 * time.Millisecond)
}

func (rodney *Rodney) mainLoop() error {
	log.Info("starting input listening")
	rodney.controller.StartListener(0)
	for {
		select {
		case ev := <-rodney.controller.Events:
			if ev.Button != nil {
				if ev.Button.Name == "home" {
					rodney.gracefulExit()
					return nil
				}
			} else if ev.Stick != nil {
				if ev.Stick.Name == "left" {
					lSpeed, rSpeed := getMotorsSpeedFromStick(ev.Stick.X, ev.Stick.Y)
					rodney.leftMotor.SetSpeed(lSpeed)
					rodney.rightMotor.SetSpeed(rSpeed)
				}
			}
		}
	}
}

func (rodney *Rodney) gracefulExit() {
	rodney.runLED.Off()
	rodney.bluetoothLED.Off()
	disconnectController()
}

// Start rodney
func (rodney *Rodney) Start() error {
	log.Info("I'm Rodney !")
	if err := rodney.initServos(); err != nil {
		rodney.handleFatal()
		return err
	}
	// for now, the controller have to be paired each time
	if err := rodney.getController(); err != nil {
		rodney.handleFatal()
		return err
	}
	// Avoid a wrong first input
	time.Sleep(500 * time.Millisecond)
	rodney.controller = switchprocontroller.NewSwitchProController()
	if err := rodney.mainLoop(); err != nil {
		rodney.handleFatal()
		return err
	}
	rodney.runLED.Off()
	time.Sleep(500 * time.Millisecond)
	return nil
}

func getMotorsSpeedFromStick(x float32, y float32) (int8, int8) {
	var straightSpeed int8
	// minSpeedPct := 50
	if x == 0.0 && y == 0.0 {
		return 0, 0
	}
	if y == 0 {
		turnSpeed := int8(math.RoundToEven(float64(x)))
		return turnSpeed, turnSpeed
	}
	straightSpeed = int8(math.RoundToEven(float64(y)))
	if x > 0 {
		lSpeed := straightSpeed
		rSpeed := straightSpeed - int8(math.RoundToEven(float64(straightSpeed)/100.0*float64(x)))
		return lSpeed, -rSpeed
	} else if x < 0 {
		lSpeed := straightSpeed + int8(math.RoundToEven(float64(straightSpeed)/100.0*float64(x)))
		rSpeed := straightSpeed
		return lSpeed, -rSpeed
	} else {
		return straightSpeed, -straightSpeed
	}
}

/*
0     100

50    100
*/
