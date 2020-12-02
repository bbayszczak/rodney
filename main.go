package main

import (
	"fmt"
	"time"

	"github.com/0xcafed00d/joystick"
)

func main() {
	js, err := joystick.Open(0)
	if err != nil {
		panic(err)
	}
	defer js.Close()
	fmt.Printf("Joystick Name: %s", js.Name())
	fmt.Printf("   Axis Count: %d", js.AxisCount())
	fmt.Printf(" Button Count: %d", js.ButtonCount())
	for {
		state, err := js.Read()
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Millisecond * 100)
		fmt.Printf("Axis Data: %v -- Buttons Data: %v\n", state.AxisData, state.Buttons)
	}
}
