package main

import (
	"fmt"
	"os"

	"github.com/bbayszczak/rodney/pkg/rodney"
	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
)

func initLogger() {
	fd, err := os.OpenFile("./rodney.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Could not open log file : " + err.Error())
	}
	log.SetOutput(fd)
	log.SetLevel(log.DebugLevel)
}
func main() {
	initLogger()
	err := rpio.Open()
	if err != nil {
		log.WithField("error", err).Fatal("impossible to init GPIO")
	}
	r := rodney.NewRodney()
	if err := r.Start(); err != nil {
		log.WithField("error", err).Fatal("a fatal error occured")
	}
	defer rpio.Close()
}
