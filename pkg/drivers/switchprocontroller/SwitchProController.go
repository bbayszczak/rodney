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
	fetchDelta time.Duration = 10 * time.Millisecond
)

// SwitchProController represent the physical controller
type SwitchProController struct {
	FetchDelta time.Duration
	// each time a new event is received, true is sent to this channel
	Event            chan bool `json:"-"`
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
//
// 0 is the default value when the stick is in the default position
type Stick struct {
	X float32
	Y float32
}

func setAxisValue(axis *float32, value float32, controller *SwitchProController) {
	if value > 100 {
		if *axis != 100.0 {
			controller.eventChange()
		}
		*axis = 100.0
	} else if value < -100 {
		if *axis != -100.0 {
			controller.eventChange()
		}
		*axis = -100.0
	} else {
		if *axis != value {
			controller.eventChange()
		}
		*axis = value
	}
}

func (controller *SwitchProController) updateSticks(axisData []int) {
	// Left stick update
	setAxisValue(&controller.StickLeft.X, (-100.0*float32(axisData[0]))/stickMin, controller)
	setAxisValue(&controller.StickLeft.Y, (100.0*float32(axisData[1]))/stickMax, controller)
	// Right stick update
	setAxisValue(&controller.StickRight.X, (-100.0*float32(axisData[2]))/stickMin, controller)
	setAxisValue(&controller.StickRight.Y, (100.0*float32(axisData[3]))/stickMax, controller)
	// DPad update
	setAxisValue(&controller.StickPad.X, (-100.0*float32(axisData[4]))/stickMin, controller)
	setAxisValue(&controller.StickPad.Y, (100.0*float32(axisData[5]))/stickMax, controller)
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
			if controller.ButtonCapture != isPressed {
				controller.eventChange()
				controller.ButtonCapture = isPressed
			}
		case 4096:
			if controller.ButtonHome != isPressed {
				controller.eventChange()
				controller.ButtonHome = isPressed
			}
		case 2048:
			if controller.ButtonRightStick != isPressed {
				controller.eventChange()
				controller.ButtonRightStick = isPressed
			}
		case 1024:
			if controller.ButtonLeftStick != isPressed {
				controller.eventChange()
				controller.ButtonLeftStick = isPressed
			}
		case 512:
			if controller.ButtonPlus != isPressed {
				controller.eventChange()
				controller.ButtonPlus = isPressed
			}
		case 256:
			if controller.ButtonLess != isPressed {
				controller.eventChange()
				controller.ButtonLess = isPressed
			}
		case 128:
			if controller.ButtonZR != isPressed {
				controller.eventChange()
				controller.ButtonZR = isPressed
			}
		case 64:
			if controller.ButtonZL != isPressed {
				controller.eventChange()
				controller.ButtonZL = isPressed
			}
		case 32:
			if controller.ButtonR != isPressed {
				controller.eventChange()
				controller.ButtonR = isPressed
			}
		case 16:
			if controller.ButtonL != isPressed {
				controller.eventChange()
				controller.ButtonL = isPressed
			}
		case 8:
			if controller.ButtonX != isPressed {
				controller.eventChange()
				controller.ButtonX = isPressed
			}
		case 4:
			if controller.ButtonY != isPressed {
				controller.eventChange()
				controller.ButtonY = isPressed
			}
		case 2:
			if controller.ButtonA != isPressed {
				controller.eventChange()
				controller.ButtonA = isPressed
			}
		case 1:
			if controller.ButtonB != isPressed {
				controller.eventChange()
				controller.ButtonB = isPressed
			}
		}
	}
}

// NewSwitchProController creates and init a SwitchProController instance
func NewSwitchProController() *SwitchProController {
	log.Info("creating new SwitchProController")
	controller := SwitchProController{
		FetchDelta: fetchDelta,
		Event:      make(chan bool, 1),
	}
	return &controller
}

// Display pprint the current controller status
func (controller *SwitchProController) Display() {
	marshalled, err := json.MarshalIndent(*controller, "", "  ")
	if err != nil {
		log.WithField("error", err).Error("impossible to marshal controller")
		fmt.Println(err)
	}
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
			time.Sleep(fetchDelta)
		}
	}()
}

func (controller *SwitchProController) eventChange() {
	if len(controller.Event) < cap(controller.Event) {
		controller.Event <- true
	}
}
