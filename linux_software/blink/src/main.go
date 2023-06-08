package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Application that is intended to run on Cora dev board with custom
// timer IP with software or "hardware" interrupts
// If the base addresses are not passed, the hardcoded defaults are used.
func main() {
	// Parsing the arguments from launch command
	var timerOption string      // User defined FPGA timer threshold value
	var timerBaseAddress string // User defined base address of the FPGA Timer IP
	var ledsBaseAddress string  // User defined base address of the FPGA LED GPIOs
	flag.StringVar(&timerOption, "timer", "49999999", "Enter amount of 10 ns cycles; period = 10ns * cycles")
	flag.StringVar(&timerBaseAddress, "timer-base", "0x43c00000", "Base address of the timer FPGA IP in hex")
	flag.StringVar(&ledsBaseAddress, "leds-base", "0x41210000", "Base address of the LEDS GPIO FPGA IP in hex")
	flag.Parse()

	err := setHighestPriority()
	if err != nil {
		fmt.Printf("Failed to set highest priority: %s", err)
	}

	err = blink(timerOption, timerBaseAddress, timerBaseAddress)
	if err != nil {
		fmt.Printf("Blink failed: %v\n", err)
	}
}

// This function takes advantage of the CPU isolation and RT-linux features.
func setHighestPriority() error {
	// Set the process's nice value to the highest priority (-20)
	err := unix.Setpriority(unix.PRIO_PROCESS, 0, -20)
	if err != nil {
		return fmt.Errorf("failed to set nice value: %v", err)
	}

	runtime.GOMAXPROCS(1) // Let the application run on only one processor (the isolated one)

	// Check if there are isolated processors in the system
	data, err := os.ReadFile("/sys/devices/system/cpu/isolated")
	if err != nil {
		return fmt.Errorf("failed to read isolated CPUs:%s", err)
	}

	isolated := strings.TrimSpace(string(data))
	if len(isolated) == 0 {
		return fmt.Errorf("no isolated CPUs found")
	}

	cpuNumbers := strings.Split(isolated, ",")
	var mask uint64 = 0
	for _, cpu := range cpuNumbers {
		cpuNum, err := strconv.Atoi(cpu)
		if err != nil {
			return fmt.Errorf("failed to parse isolated CPU number:%s", err)
		}
		fmt.Println("Found isolated CPUs:", cpuNum)
		mask |= 1 << uint64(cpuNum)
	}

	// If the check was successfull, update CPU affinity to the isolated processor
	if _, _, err := syscall.RawSyscall(syscall.SYS_SCHED_SETAFFINITY, 0, uintptr(unsafe.Sizeof(mask)), uintptr(unsafe.Pointer(&mask))); err != 0 {
		return fmt.Errorf("failed to set CPU affinity:%s", err)
	}

	// Scheduling policy SCHED_FIFO is necessary to utilize the NOHZ_FULL kernel bootarg for the isolated CPU
	const SCHED_FIFO = 1
	var priority = uint32(9) // May choose another value (>0 , <100)
	if _, _, err := syscall.RawSyscall(syscall.SYS_SCHED_SETSCHEDULER, uintptr(0), uintptr(SCHED_FIFO), uintptr(unsafe.Pointer(&priority))); err != 0 {
		return fmt.Errorf("failed to set scheduling policy SCHED_FIFO: %s; scheduler is unchanged", err)
	}
	fmt.Println("Updated scheduling policy to SCHED_FIFO")

	return nil
}
