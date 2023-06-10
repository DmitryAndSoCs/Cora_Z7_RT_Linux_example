package fpgaLEDs

import (
	"errors"
	"fmt"

	"periph.io/x/host/v3/pmem"
)

type leds struct {
	ledsReg uint32 // Offset 0x00, holds all 6 bits for 2 RGB LEDs
}

var (
	errCreationFailed = errors.New("failed to create a new LEDs GPIO device from FPGA memory")
	//errWriteFailed      = errors.New("failed to write into FPGA's memory")
	errWrongLEDSettings = errors.New("received wrong settings for LEDs")
)

// Creates an array of uint32 values that are addressed at the passed base address
// Offsets are base + 4 per each element in the array.
func NewFpgaLEDs(baseAddr uint64) (*leds, error) {
	var fpga *leds
	if err := pmem.MapAsPOD(baseAddr, &fpga); err != nil {
		return nil, fmt.Errorf("%w:%s \n\r", errCreationFailed, err.Error())
	}
	return fpga, nil
}

// Sets specified (string) color to specified (int) LED.
// Example: .SetColor("R",0); .SetColor("off", 0) to turn off the LED
func (leds *leds) SetColor(color string, ledNum int) error {
	// Check if the desired LED valid
	if ledNum != 0 && ledNum != 1 {
		return fmt.Errorf("%w: ledNum - %d", errWrongLEDSettings, ledNum)
	}

	// Simple handling of the LEDs
	switch ledNum {
	case 0:
		switch color {
		case "R":
			leds.ledsReg |= 0x00000001
		case "G":
			leds.ledsReg |= 0x00000002
		case "B":
			leds.ledsReg |= 0x00000004
		case "off":
			leds.ledsReg &= ^uint32(0x00000007)
		default:
			return fmt.Errorf("%w: unsupported color - %s", errWrongLEDSettings, color)
		}
	case 1:
		switch color {
		case "R":
			leds.ledsReg |= 0x00000008
		case "G":
			leds.ledsReg |= 0x00000010
		case "B":
			leds.ledsReg |= 0x00000020
		case "off":
			leds.ledsReg &= ^uint32(0x00000038)
		default:
			return fmt.Errorf("%w: unsupported color - %s", errWrongLEDSettings, color)
		}
	}

	return nil
}
