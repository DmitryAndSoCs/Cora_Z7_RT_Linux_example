package fpgaTimer

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"periph.io/x/host/v3/pmem"
)

var (
	errCreationFailed = errors.New("failed to create a new timer device from FPGA memory")
	errWriteFailed    = errors.New("failed to write into FPGA's memory")
	errKoRegister     = errors.New("failed to register the application for interrupts fetching")
)

type timer struct {
	threshold   uint32 // Offset 0x00, maximum timer value
	softwareInt uint32 // Offset 0x04, Register that has the bit for polling option
}

// Creates an array of uint32 values that are addressed at the passed base address
// Offsets are base + 4 per each element in the array.
// It also checks for the kernel module interrupt handler and tells if it's present.
func NewFpgaTimer(baseAddr uint64) (timerStruct *timer, koPresent bool, err error) {
	var fpga *timer
	if err := pmem.MapAsPOD(baseAddr, &fpga); err != nil {
		return nil, false, fmt.Errorf("%w:%s \n\r", errCreationFailed, err.Error())
	}
	err = registerForInterrupts("/sys/kernel/fpgatimer/pid") // hardcoded, it can only work with this specific kernel module
	if err != nil {
		fmt.Printf("%s, OS soft interrupt from FPGA IP is not available", errKoRegister)
		return fpga, false, nil
	} else {
		return fpga, true, nil
	}

}

// Set new threshold value for timer, int is for user's convenience
func (timer *timer) SetFpgaTimerThreshold(newThreshold int) (err error) {
	if newThreshold <= 0 { // check that the value is not negative to avoid unpredicted period values
		return fmt.Errorf("threshold value must be > 0")
	}
	timer.threshold = uint32(newThreshold)       // Set the value
	if timer.threshold != uint32(newThreshold) { // Check written value
		return fmt.Errorf("%w", errWriteFailed)
	}
	return nil
}

// Function for polling the software interrupt bit
func (timer *timer) SoftIntHappened() bool {
	return (timer.softwareInt&0x00000001 == 1)
}

// Register for the interrupt handler that resides in the kernel.
// Everything is hardcoded, it works only if the driver is present in the kernel.
func registerForInterrupts(pidPath string) error {
	// Check if the pid file exists
	if _, err := os.Stat(pidPath); os.IsNotExist(err) {
		return fmt.Errorf("pid file not found: %s", pidPath)
	}

	// Get own application's PID
	pid := syscall.Getpid()

	// Convert the PID to a string
	pidStr := strconv.Itoa(pid)

	// Write the PID to the sysfs file
	err := os.WriteFile(pidPath, []byte(pidStr), 0644)
	if err != nil {
		return fmt.Errorf("failed to write PID to file: %v", err)
	}
	return nil
}
