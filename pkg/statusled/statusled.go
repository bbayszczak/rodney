package statusled

import (
	"time"

	driver "github.com/bbayszczak/raspberrypi-go-drivers/led"
)

// StatusLED represent a status LED
type StatusLED struct {
	led   *driver.LED
	state string
}

// NewStatusLED creates a StatusLED instance
func NewStatusLED(pinID uint8) *StatusLED {
	statusLED := StatusLED{
		led:   driver.NewLED(pinID),
		state: "off",
	}
	go statusLED.loop()
	return &statusLED
}

func (statusLED *StatusLED) loop() {
	for {
		if statusLED.state == "on" {
			statusLED.led.On()
		} else if statusLED.state == "off" {
			statusLED.led.Off()
		} else if statusLED.state == "blink" {
			statusLED.led.Toggle()
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// On switch on the LED
func (statusLED *StatusLED) On() {
	statusLED.state = "on"
}

// Off switch off the LED
func (statusLED *StatusLED) Off() {
	statusLED.state = "off"
}

// Blink make the LED blink
func (statusLED *StatusLED) Blink() {
	statusLED.state = "blink"
}
