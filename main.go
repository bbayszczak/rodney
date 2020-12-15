package main

import (
	"fmt"
	"os"

	"github.com/bbayszczak/rodney/pkg/rodney"
	"github.com/stianeikeland/go-rpio/v4"
)

func main() {
	fmt.Println("I'm Rodney !")
	err := rpio.Open()
	if err != nil {
		fmt.Println("impossible to init gpio")
		os.Exit(1)
	}
	r := rodney.NewRodney()
	if err := r.Start(); err != nil {
		fmt.Println("error")
		fmt.Println(err)
	}
	defer rpio.Close()
}
