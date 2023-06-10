package main

import (
	"blink/fpgaLEDs"
	"blink/fpgaTimer"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// Function checks the passed parameters from the main app,
// checks if interrupts are avaiable through FPGA timer interrupt kernel module
// and blinks until CTRL+C is issued by the user.
func blink(timerPeriod, timerBaseAddress, ledsBaseAddress string) (err error) {
	// Check for the correctness of the input

	timerBase, err := strconv.ParseUint(timerBaseAddress, 0, 64) // parse timer base address for hex value
	if err != nil {
		fmt.Printf("invalid timer base address (must be hex): %v", err)
		fmt.Println("Using default: 0x43c00000")
		timerBase = 0x43c00000
	}
	ledsBase, err := strconv.ParseUint(ledsBaseAddress, 0, 64) // parse leds base address for hex value
	if err != nil {
		fmt.Printf("invalid leds base address (must be hex): %v", err)
		fmt.Println("Using default: 0x41210000")
		ledsBase = 0x41210000
	}

	// Parse the period parameter
	cycles, err := strconv.ParseInt(timerPeriod, 10, 32)
	if err != nil {
		fmt.Printf("Failed to parse timer option: %v\n", err)
		fmt.Println("Using default 500 ms period.")
		cycles = 49999999
	}

	timer, koPreset, err := fpgaTimer.NewFpgaTimer(timerBase) // create new timer handle
	if err != nil {
		return fmt.Errorf("failed to create new FPGA timer: %v", err)
	}
	leds, err := fpgaLEDs.NewFpgaLEDs(ledsBase) // create new LEDs handle
	if err != nil {
		return fmt.Errorf("failed to create new FPGA LEDs: %v", err)
	}

	err = timer.SetFpgaTimerThreshold(int(cycles)) // apply the new period from the user
	if err != nil {
		return fmt.Errorf("failed to update timer threshold:%v", err)
	}

	flip := false // make a signal that would act as a flip-flop for color switching
	// Define the blink function for numerous future uses
	blinkLeds := func() (err error) {
		if !flip {
			err = leds.SetColor("R", 0)
			if err != nil {
				return fmt.Errorf("color switching: %w", err)
			}
			err = leds.SetColor("off", 1)
			if err != nil {
				return fmt.Errorf("color switching: %w", err)
			}
		} else {
			err = leds.SetColor("off", 0)
			if err != nil {
				return fmt.Errorf("color switching: %w", err)
			}
			err = leds.SetColor("G", 1)
			if err != nil {
				return fmt.Errorf("color switching: %w", err)
			}
		}
		flip = !flip
		return nil
	}

	// Wrapping context for the infinite cycle
	ctx, ctxCancel := context.WithCancel(context.Background())

	// Create a channel to receive signals from the operating system
	cancelSig := make(chan os.Signal, 1) // for CTRL+C
	signal.Notify(cancelSig, syscall.SIGINT, syscall.SIGTERM)

	intSig := make(chan os.Signal, 1) // for interrupts
	signal.Notify(intSig, syscall.SIGUSR1)

	// Start a goroutine to wait for the context to be cancelled
	go func() {
		<-ctx.Done()
		fmt.Println("Blinking stopped")
	}()

	// Work with the context
	if koPreset { // if the interrupt from the kernel module is avaiable, use it
		fmt.Println("Using kernel module interrupt...")
		for {
			select {
			case <-cancelSig: // if user pressed CTRL+C
				fmt.Printf("Received CTRL+C\n")
				ctxCancel() // cancel the context
				return nil
			case <-intSig:
				fmt.Printf("Got FPGA interrupt from kernel module!\n\n")
				err = blinkLeds()
				if err != nil {
					fmt.Printf("%s\n", err)
					fmt.Println("Cancelling context...")
					ctxCancel()
				}
			}
		}
	} else {
		fmt.Println("Using POLLing...")
		for {
			select {
			case <-cancelSig: // if user pressed CTRL+C
				fmt.Printf("Received CTRL+C\n")
				ctxCancel()
				return nil
			default:
				if timer.SoftIntHappened() {
					fmt.Printf("POLLing: found period happened bit")
					err = blinkLeds()
					if err != nil {
						fmt.Printf("%s\n", err)
						ctxCancel()
					}
				}
			}
		}
	}
}
