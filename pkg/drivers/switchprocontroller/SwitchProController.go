package switchprocontroller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xcafed00d/joystick"
	log "github.com/sirupsen/logrus"
)

const (
	stickMax   float32       = 28500
	stickMin   float32       = -29735
	fetchDelay time.Duration = 50 * time.Millisecond
)

// SwitchProController represent the physical controller
type SwitchProController struct {
	StickLeft        Stick
	StickRight       Stick
	StickPad         Stick
	ButtonB          int
	ButtonA          int
	ButtonY          int
	ButtonX          int
	ButtonL          int
	ButtonR          int
	ButtonZL         int
	ButtonZR         int
	ButtonLess       int
	ButtonPlus       int
	ButtonLeftStick  int
	ButtonRightStick int
	ButtonHome       int
	ButtonCapture    int
}

// Stick represent a physical stick
//
// It contains a value for each axis (x and y), this value is contained between -100 and 100
// 0 is the default value when the stick is in the default position
type Stick struct {
	X float32
	Y float32
}

func setAxisValue(axis *float32, value float32) {
	if value > 100 {
		*axis = 100.0
	} else if value < -100 {
		*axis = -100.0
	} else {
		*axis = value
	}
}

func (controller *SwitchProController) updateSticks(axisData []int) {
	// Left stick update
	setAxisValue(&controller.StickLeft.X, (-100.0*float32(axisData[0]))/stickMin)
	setAxisValue(&controller.StickLeft.Y, (100.0*float32(axisData[1]))/stickMax)
	// Right stick update
	setAxisValue(&controller.StickRight.X, (-100.0*float32(axisData[2]))/stickMin)
	setAxisValue(&controller.StickRight.Y, (100.0*float32(axisData[3]))/stickMax)
	// DPad update
	setAxisValue(&controller.StickPad.X, (-100.0*float32(axisData[4]))/stickMin)
	setAxisValue(&controller.StickPad.Y, (100.0*float32(axisData[5]))/stickMax)
}

func (controller *SwitchProController) updateButtons(buttons uint32) {
	values := []uint32{8192, 4096, 2048, 1024, 512, 256, 128, 64, 32, 16, 8, 4, 2, 1}
	for _, val := range values {
		var isPressed int
		isPressed = 0
		if buttons >= val {
			buttons = buttons - val
			isPressed = 1
		}
		switch val {
		case 8192:
			controller.ButtonCapture = isPressed
		case 4096:
			controller.ButtonHome = isPressed
		case 2048:
			controller.ButtonRightStick = isPressed
		case 1024:
			controller.ButtonLeftStick = isPressed
		case 512:
			controller.ButtonPlus = isPressed
		case 256:
			controller.ButtonLess = isPressed
		case 128:
			controller.ButtonZR = isPressed
		case 64:
			controller.ButtonZL = isPressed
		case 32:
			controller.ButtonR = isPressed
		case 16:
			controller.ButtonL = isPressed
		case 8:
			controller.ButtonX = isPressed
		case 4:
			controller.ButtonY = isPressed
		case 2:
			controller.ButtonA = isPressed
		case 1:
			controller.ButtonB = isPressed
		}
	}
}

// NewSwitchProController creates a SwitchProController instance
func NewSwitchProController() *SwitchProController {
	log.Info("creating new SwitchProController")
	var controller SwitchProController
	return &controller
}

// Display pprint the current controller status
func (controller *SwitchProController) Display() {
	marshalled, _ := json.MarshalIndent(*controller, "", "  ")
	fmt.Printf(string(marshalled))
}

// StartListener starts listening for controller inputs and
// keep the controller instance up to date
func (controller *SwitchProController) StartListener(joystickID int) {
	js, err := joystick.Open(joystickID)
	if err != nil {
		panic(err)
	}
	go func() {
		defer js.Close()
		for {
			state, err := js.Read()
			if err != nil {
				panic(err)
			}
			controller.updateSticks(state.AxisData)
			controller.updateButtons(state.Buttons)
			time.Sleep(fetchDelay)
		}
	}()
}
