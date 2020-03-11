package mapper

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/bendahl/uinput"
	evdev "github.com/gvalkov/golang-evdev"
)

type shutdown struct {
	initiate time.Time
	count    int
}

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
	keyboard, err := uinput.CreateKeyboard("/dev/uinput", []byte("mapper-kbd"))
	if err != nil {
		return err
	}

	defer keyboard.Close()

	for {
		device, err := connectDevice(deviceName)
		if err != nil {
			return err
		}

		s := &shutdown{time.Now(), 0}

		for {
			err := processEvent(device, keyboard, s)
			if err != nil {
				log.Printf("Error: %v", err)
				break
			}
		}
	}
}

func findDevice(name string) (string, error) {
	log.Println("looking for device ...")
	for {
		devices, err := evdev.ListInputDevices()
		if err != nil {
			return "", err
		}

		for _, device := range devices {
			if len(name) > 0 {
				if strings.TrimSpace(device.Name) == name {
					log.Printf("found device %q\n", device.Name)
					return device.Fn, nil
				}
			} else {
				lower := strings.ToLower(device.Name)
				if strings.Contains(lower, "gamepad") {
					log.Printf("found device %q\n", device.Name)
					return device.Fn, nil
				}
			}
		}
		time.Sleep(time.Second)
	}
}

func connectDevice(deviceName string) (*evdev.InputDevice, error) {
	devnode, err := findDevice(deviceName)
	if err != nil {
		return nil, err
	}

	log.Printf("connect to %q\n", devnode)
	return evdev.Open(devnode)
}

func processEvent(device *evdev.InputDevice, keyboard uinput.Keyboard, s *shutdown) error {
	events, err := device.Read()
	if err != nil {
		return err
	}

	for _, event := range events {
		if int(event.Type) == evdev.EV_KEY {
			if event.Value == 1 {
				switch int(event.Code) {
				case evdev.BTN_A, evdev.BTN_TR, evdev.BTN_THUMB, evdev.BTN_PINKIE:
					keyboard.KeyPress(uinput.KeyRight)
				case evdev.BTN_B, evdev.BTN_THUMB2:
					keyboard.KeyPress(uinput.KeyDown)
				case evdev.BTN_X, evdev.BTN_TRIGGER:
					keyboard.KeyPress(uinput.KeyUp)
				case evdev.BTN_Y, evdev.BTN_TL, evdev.BTN_TOP, evdev.BTN_TOP2:
					keyboard.KeyPress(uinput.KeyLeft)
				case evdev.BTN_START:
					s.trigger()
				}
			}
		}
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

	return nil
}

func (s *shutdown) trigger() {
	if time.Since(s.initiate) < time.Second {
		s.count++
		if s.count >= 5 {
			cmd := exec.Command("systemctl", "poweroff")
			err := cmd.Start()
			if err != nil {
				log.Printf("Error: %v", err)
			}
		}
	} else {
		s.initiate = time.Now()
		s.count = 0
	}
}
