package rodney

import (
	"time"

	"github.com/bbayszczak/rodney/pkg/statusled"
)

// Rodney represent the robot
type Rodney struct {
	runLED *statusled.StatusLED
}

// NewRodney Creqtes a Rodney instance
func NewRodney() *Rodney {
	rodney := Rodney{
		runLED: statusled.NewStatusLED(17),
	}
	rodney.runLED.On()
	return &rodney
}

// Start rodney
func (rodney *Rodney) Start() error {
	time.Sleep(time.Second)
	rodney.runLED.Off()
	time.Sleep(600 * time.Millisecond)
	return nil
}
