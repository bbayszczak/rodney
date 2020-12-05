package led_test

import (
	"fmt"
	"os"
	"time"

	"github.com/bbayszczak/rodney/pkg/drivers/led"
	"github.com/stianeikeland/go-rpio/v4"
)

func Example_onoff() {
	err := rpio.Open()
	if err != nil {
		os.Exit(1)
	}
	defer rpio.Close()

	led := led.NewLED(18)
	led.On()
	fmt.Println(led.GetState())
	time.Sleep(time.Second)
	led.Off()
	fmt.Println(led.GetState())
	time.Sleep(time.Second)
	led.Toggle()
	fmt.Println(led.GetState())
	time.Sleep(time.Second)
	led.Toggle()
	fmt.Println(led.GetState())
}

func Example_dimmable() {
	err := rpio.Open()
	if err != nil {
		os.Exit(1)
	}
	defer rpio.Close()

	led := led.NewLED(18)
	led.Dimmable()
	for true {
		for i := 0; i <= 100; i++ {
			if err := led.SetBrightness(i); err != nil {
				fmt.Println(err.Error())
			}
			time.Sleep(time.Millisecond * 10)
		}
		for i := 100; i >= 0; i-- {
			if err := led.SetBrightness(i); err != nil {
				fmt.Println(err.Error())
			}
			time.Sleep(time.Millisecond * 10)
		}
	}
}
