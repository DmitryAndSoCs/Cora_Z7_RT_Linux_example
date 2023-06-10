# Golang application that blinks LEDs using FPGA IP's interrupts

This application is intended to show an approach to interface with FPGA IPs created by user.
It utilizes the interrupt from the kernle module created at petalinux project creation. 

It checks for the presense of the kernel module and runs it using the interrupt if present, 
or by polling if the module is not present. 

## How to build

On a linux machine (WSL2 will suffice as well):

```bash
cd src
GOOS="linux" GOARCH="arm" go build -o ../bin/blink.run
```

This would create an executable in bin folder. Copy it to the board (and\or embed it into the petalinux image) and run it.

Usage:

```bash
./blink.run # Run with default parameters and base addresses
./blink.run --timer 20999999 # increase the blinking speed by reducing the interrupt period
./blink.run --timer 15999999 --timer-base 0x4c310000 --leds-base 0x41220000 # if you need to specify different base addresses of the IP modules
```