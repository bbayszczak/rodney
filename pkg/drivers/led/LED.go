package led

import (
	"fmt"

	"github.com/stianeikeland/go-rpio/v4"
)

// LED instance
type LED struct {
	pinID    uint8
	pin      rpio.Pin
	dimmable bool
}

// NewLED creates a new LED instance
func NewLED(pinID uint8) *LED {
	led := LED{
		pinID:    pinID,
		dimmable: false,
	}
	led.pin = rpio.Pin(led.pinID)
	led.pin.Mode(rpio.Output)
	return &led
}

// Toggle change the state of the LED
func (led *LED) Toggle() error {
	led.pin.Toggle()
	return nil
}

// On lights on the LED
func (led *LED) On() error {
	led.pin.High()
	return nil
}

// Off lights off the LED
func (led *LED) Off() error {
	led.pin.Low()
	return nil
}

// GetState return the current state of the lED. on: true, off: false
func (led *LED) GetState() (bool, error) {
	state := led.pin.Read()
	if state == 1 {
		return true, nil
	} else if state == 0 {
		return false, nil
	}
	return false, fmt.Errorf("unknown state '%d'", state)
}

// Dimmable make the LED dimmable (only available on PWM pins)
func (led *LED) Dimmable() error {
	led.dimmable = true
	led.pin.Mode(rpio.Pwm)
	led.pin.Freq(10000)
	return nil
}

// NonDimmable make the dimmable LED non dimmable
func (led *LED) NonDimmable() error {
	led.dimmable = false
	led.pin.Mode(rpio.Output)
	return nil
}

// SetBrightness set the brightness of the LED (only available on PWM pins)
func (led *LED) SetBrightness(percentage int) error {
	if !led.dimmable {
		return fmt.Errorf("LED is not setup as dimmable")
	}
	if percentage < 0 || percentage > 100 {
		return fmt.Errorf("percentage value must be >= 0 <= 100")
	}
	led.pin.DutyCycle(uint32(percentage), 100)
	return nil
}
