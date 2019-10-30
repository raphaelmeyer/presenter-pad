package mapper

import (
	"fmt"
	"log"
	"strings"

	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
)

// ListDevices lists all input devices
func ListDevices() error {
	devices, err := evdev.ListInputDevices()
	if err != nil {
		return err
	}

	if len(devices) == 0 {
		fmt.Print("No devices found.\n")
		return nil
	}

	fmt.Printf("Found %d devices:\n", len(devices))
	for _, device := range devices {
		fmt.Printf("- %q\n", strings.TrimSpace(device.Name))
	}

	return nil
}

// Run starts the gamepad to keystroke mapper
func Run(deviceName string) error {

	devnode, err := findDevice(deviceName)
	if err != nil {
		return err
	}

	log.Print(devnode)

	device, err := evdev.Open(devnode)
	if err != nil {
		return err
	}

	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("mapper-kbd"))
	if err != nil {
		return err
	}

	defer keyboard.Close()

	for {
		events, err := device.Read()
		if err != nil {
			return err
		}
		for _, event := range events {
			if int(event.Type) == evdev.EV_ABS {
				switch int(event.Code) {
				case evdev.ABS_X:
					switch int(event.Value) {
					case 0:
						keyboard.KeyPress(uinput.KeyLeft)
					case 255:
						keyboard.KeyPress(uinput.KeyRight)
					}
				case evdev.ABS_Y:
					switch int(event.Value) {
					case 0:
						keyboard.KeyPress(uinput.KeyUp)
					case 255:
						keyboard.KeyPress(uinput.KeyDown)
					}
				}
			}
		}
	}
}

func findDevice(name string) (string, error) {
	devices, err := evdev.ListInputDevices()
	if err != nil {
		return "", err
	}

	if len(name) > 0 {
		for _, device := range devices {
			if strings.TrimSpace(device.Name) == name {
				log.Printf("found device %q\n", device.Name)
				return device.Fn, nil
			}
		}
	} else {
		for _, device := range devices {
			lower := strings.ToLower(device.Name)
			if strings.Contains(lower, "gamepad") {
				log.Printf("found device %q\n", device.Name)
				return device.Fn, nil
			}
		}
	}

	return "", fmt.Errorf("device not found")
}
