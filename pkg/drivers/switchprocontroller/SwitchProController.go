package switchprocontroller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xcafed00d/joystick"
	log "github.com/sirupsen/logrus"
)

const (
	stickPeakValue float32       = 20000
	fetchDelta     time.Duration = 10 * time.Millisecond
)

// SwitchProController represent the physical controller
type SwitchProController struct {
	FetchDelta time.Duration
	// each time a new event is received, true is sent to this channel
	Event            chan bool `json:"-"`
	Sticks           []*Stick
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
// It contains a value (in %) for each axis (x and y), this value is contained between -100 and 100
//
// 0 is the default value when the stick is in the default position
type Stick struct {
	Name string
	X    float32
	Y    float32
	xMin float32
	xMax float32
	yMin float32
	yMax float32
}

// GetStick returns a pointer to a Stick instance with specified name
//
// err not nil if stick name not found
func (controller *SwitchProController) GetStick(name string) (*Stick, error) {
	for _, stick := range controller.Sticks {
		if stick.Name == name {
			return stick, nil
		}
	}
	return nil, fmt.Errorf("impossible to find stick")
}

func (controller *SwitchProController) updateStick(name string, x float32, y float32) {
	stick, err := controller.GetStick(name)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"name":  name,
		}).Error("stick not found")
	}
	if x < stick.xMin {
		stick.xMin = x
		log.WithField("stick.xMin", x).Info("stick peak value changed")
	} else if x > stick.xMax {
		stick.xMax = x
		log.WithField("stick.xMax", x).Info("stick peak value changed")
	}
	if y < stick.yMin {
		stick.yMin = y
		log.WithField("stick.yMin", y).Info("stick peak value changed")
	} else if y > stick.yMax {
		stick.yMax = y
		log.WithField("stick.yMax", y).Info("stick peak value changed")
	}
	var newX float32
	var newY float32
	if x > 0 {
		newX = (100.0 * x) / stick.xMax
	} else if x < 0 {
		newX = (-100.0 * x) / stick.xMin
	} else if x == 0 {
		newX = 0
	}
	if newX != stick.X {
		stick.X = newX
		controller.eventChange()
	}
	if y > 0 {
		newY = (-100.0 * y) / stick.yMax
	} else if y < 0 {
		newY = (100.0 * y) / stick.yMin
	} else if y == 0 {
		newY = 0
	}
	if newY != stick.Y {
		stick.Y = newY
		controller.eventChange()
	}
}

func (controller *SwitchProController) updateSticks(axisData []int) {
	controller.updateStick("left", float32(axisData[0]), float32(axisData[1]))
	controller.updateStick("right", float32(axisData[2]), float32(axisData[3]))
	controller.updateStick("pad", float32(axisData[4]), float32(axisData[5]))
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

func initStick(name string) *Stick {
	return &Stick{
		Name: name,
		X:    0,
		Y:    0,
		xMin: -stickPeakValue,
		xMax: stickPeakValue,
		yMin: -stickPeakValue,
		yMax: stickPeakValue,
	}
}

// NewSwitchProController creates and initialize a SwitchProController instance
func NewSwitchProController() *SwitchProController {
	log.Info("creating new SwitchProController")
	controller := SwitchProController{
		FetchDelta: fetchDelta,
		Event:      make(chan bool, 1),
		Sticks: []*Stick{
			initStick("left"),
			initStick("right"),
			initStick("pad"),
		},
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
