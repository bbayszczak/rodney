package rodney

import (
	"errors"

	"github.com/muka/go-bluetooth/api"
	"github.com/muka/go-bluetooth/bluez/profile/adapter"
	"github.com/muka/go-bluetooth/bluez/profile/agent"
	"github.com/muka/go-bluetooth/bluez/profile/device"
	log "github.com/sirupsen/logrus"
)

func (rodney *Rodney) getController() error {
	log.Info("getting controller")
	rodney.bluetoothLED.Blink()
	defer api.Exit()
	controllerDevice, err := discoverController()
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"name":    controllerDevice.Properties.Name,
		"address": controllerDevice.Properties.Address,
		"rssi":    controllerDevice.Properties.RSSI,
		"paired":  controllerDevice.Properties.Paired,
	}).Info("bluetooth: controller found")
	if !controllerDevice.Properties.Paired {
		if err := pairController(controllerDevice); err != nil {
			log.WithFields(log.Fields{
				"name":    controllerDevice.Properties.Name,
				"address": controllerDevice.Properties.Address,
				"rssi":    controllerDevice.Properties.RSSI,
				"error":   err,
			}).Error("impossible to pair controller")
			return err
		}
	}
	if err := connectController(controllerDevice); err != nil {
		log.WithFields(log.Fields{
			"name":  controllerDevice.Properties.Name,
			"error": err,
		}).Error("impossible to connect to controller")
		return err
	}
	rodney.bluetoothLED.On()
	return nil
}

func connectController(controllerDevice *device.Device1) error {
	if err := controllerDevice.Connect(); err != nil {
		return err
	}
	log.WithField("name", controllerDevice.Properties.Name).Error("connection to controller successfull")
	return nil
}

func pairController(controllerDevice *device.Device1) error {
	log.WithFields(log.Fields{
		"name":   controllerDevice.Properties.Name,
		"paired": controllerDevice.Properties.Paired,
	}).Info("pairing controller")
	err := controllerDevice.Pair()
	if err != nil {
		return err
	}
	if err := agent.SetTrusted(adapter.GetDefaultAdapterID(), controllerDevice.Path()); err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"name":   controllerDevice.Properties.Name,
		"paired": controllerDevice.Properties.Paired,
	}).Info("controller paired and trusted successfully")
	return nil
}

func discoverController() (*device.Device1, error) {
	a, err := adapter.GetDefaultAdapter()
	if err != nil {
		return nil, err
	}
	log.Debug("bluetooth: flush cached devices")
	err = a.FlushDevices()
	if err != nil {
		return nil, err
	}

	log.Debug("bluetooth: start discovery")
	discovery, cancel, err := api.Discover(a, nil)
	if err != nil {
		return nil, err
	}
	for ev := range discovery {
		if ev.Type == adapter.DeviceRemoved {
			continue
		}
		dev, err := device.NewDevice1(ev.Path)
		if err != nil {
			log.WithFields(log.Fields{
				"ev.Path": ev.Path,
				"error":   err,
			}).Error("bluetooth: impossible to get device")
			continue
		}
		if dev == nil {
			log.WithFields(log.Fields{
				"ev.Path": ev.Path,
				"error":   err,
			}).Error("bluetooth: device not found")
			continue
		}
		log.WithFields(log.Fields{
			"name":    dev.Properties.Name,
			"address": dev.Properties.Address,
			"rssi":    dev.Properties.RSSI,
		}).Info("bluetooth: device found")
		if dev.Properties.Name == "Pro Controller" {
			cancel()
			return dev, nil
		}
	}
	return nil, errors.New("bluetooth: no controller found")
}

func disconnectController() error {
	return nil
}
