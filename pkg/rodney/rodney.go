package rodney

import (
	"math"
	"os/exec"
	"time"

	"github.com/bbayszczak/raspberrypi-go-drivers/l293d"
	"github.com/bbayszczak/rodney/pkg/statusled"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	"github.com/raspberrypi-go-drivers/switchprocontroller"
	log "github.com/sirupsen/logrus"
)

const (
	runLEDPin       uint8 = 23
	bluetoothLEDPin uint8 = 24
	issueLEDPin     uint8 = 25
	motor1EnPin     uint8 = 13
	motor2EnPin     uint8 = 12
	motor1aPin      uint8 = 17
	motor1bPin      uint8 = 27
	motor2aPin      uint8 = 5
	motor2bPin      uint8 = 6
)

// Rodney represent the robot
type Rodney struct {
	runLED           *statusled.StatusLED
	bluetoothLED     *statusled.StatusLED
	issueLED         *statusled.StatusLED
	rightMotor       *l293d.Motor
	leftMotor        *l293d.Motor
	controller       *switchprocontroller.SwitchProController
	controllerDevice *device.Device1
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
					rodney.rightMotor.SetSpeed(-rSpeed)
				}
			}
		}
	}
}

func (rodney *Rodney) gracefulExit() {
	rodney.disconnectController()
	rodney.runLED.Off()
	rodney.bluetoothLED.Off()
	time.Sleep(500 * time.Millisecond)
	_, err := exec.Command("/sbin/shutdown", "-h", "now").Output()
	if err != nil {
		log.WithField("error", err).Error("cannot shutdown host")
		rodney.handleFatal()
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

func getMotorsSpeedFromStick(x float32, y float32) (int, int) {
	speedRatio := 80
	baseSpeed := 100 - speedRatio
	xAbs := math.Abs(float64(x))
	yAbs := math.Abs(float64(y))
	xRelative := int(math.RoundToEven(xAbs / 100 * float64(speedRatio)))
	yRelative := int(math.RoundToEven(yAbs / 100 * float64(speedRatio)))
	if x == 0 && y == 0 {
		return 0, 0
	}
	if x == 0 {
		if y > 0 {
			return baseSpeed + yRelative, baseSpeed + yRelative
		} else if y < 0 {
			return -(baseSpeed + yRelative), -(baseSpeed + yRelative)
		}
	}
	if y == 0 {
		if x > 0 {
			return baseSpeed + xRelative, -(baseSpeed + xRelative)
		} else if x < 0 {
			return -(baseSpeed + xRelative), baseSpeed + yRelative
		}
	}
	speedRatioReduced := 0
	baseSpeedReduced := 100 - speedRatioReduced
	yRelativeReduced := int(math.RoundToEven(yAbs / 100 * float64(speedRatioReduced)))
	if y > 0 {
		if x > 0 {
			return baseSpeedReduced + yRelativeReduced, -(xRelative - 100)
		} else if x < 0 {
			return -(xRelative - 100), baseSpeedReduced + yRelativeReduced
		}
	} else if y < 0 {
		if x > 0 {
			return -(baseSpeedReduced - yRelativeReduced), -(-xRelative + 100)
		} else if x < 0 {
			return -(-xRelative + 100), -(baseSpeedReduced - yRelativeReduced)
		}
	}
	return 0, 0
}
