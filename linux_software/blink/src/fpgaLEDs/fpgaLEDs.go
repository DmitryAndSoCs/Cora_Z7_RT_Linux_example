package fpgaLEDs

import (
	"errors"
	"fmt"
	"sync"

	"periph.io/x/host/v3/pmem"
)

type leds struct {
	ledsReg uint32 // Offset 0x00, holds all 6 bits for 2 RGB LEDs

	regSync sync.RWMutex // To make the possible concurrent access safe
}

var (
	errCreationFailed   = errors.New("failed to create a new LEDs GPIO device from FPGA memory")
	errWriteFailed      = errors.New("failed to write into FPGA's memory")
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
// Example: .SetColor("R",0)
func (leds *leds) SetColor(color string, ledNum int) error {
	if ledNum != 0 && ledNum != 1 {
		return fmt.Errorf("%w: ledNum - %d", errWrongLEDSettings, ledNum)
	}

	leds.regSync.Lock()         // Keep the register safe from concurrent access
	defer leds.regSync.Unlock() // Unlock upon exiting the function

	before := leds.ledsReg

	switch color {
	case "R":
		leds.ledsReg &= ^(7 << (ledNum * 3)) // Clear the RGB bits for the specified LED
		leds.ledsReg |= (1 << (ledNum * 3))  // Set the R bit for the specified LED
	case "G":
		leds.ledsReg &= ^(7 << (ledNum * 3))      // Clear the RGB bits for the specified LED
		leds.ledsReg |= (1 << ((ledNum * 3) + 1)) // Set the G bit for the specified LED
	case "B":
		leds.ledsReg &= ^(7 << (ledNum * 3))      // Clear the RGB bits for the specified LED
		leds.ledsReg |= (1 << ((ledNum * 3) + 2)) // Set the B bit for the specified LED
	default:
		return fmt.Errorf("%w: color - %s", errWrongLEDSettings, color) // if something besides "R", "G" and "B" was passed
	}

	// The following check may be excessive for just the LED and can be safely omitted,
	// but it's generally good to check if the action was properly executed.

	// Read back the value from the register
	readValue := (leds.ledsReg >> (ledNum * 3)) & 7

	// Calculate the expected value
	expected := (before & ^(7 << (ledNum * 3))) | (1 << (ledNum * 3))

	if readValue != expected {
		// Attempt to restore the original register value
		leds.ledsReg = before
		return fmt.Errorf("%w: failed to set LED %d to color %s", errWriteFailed, ledNum, color)
	}

	return nil
}
