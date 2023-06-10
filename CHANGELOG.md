# Changelog

## [0.0.2] - 2023-06-10 - PRE-RELEASE

### Features
- Vivado project is builadble (no testbench)
- Petalinux project is buildable, kernel module is working
- Golang application blinks LEDs using kernel module interrupts (if there's no module, uses POLLing)
- Golang application is embedded into the build (`~/blinkapp/blink.run`)

### Fixed
- Proper gitignore for petalinux folder

### Changed
- Golang application simplified in terms of LED blinking

### Removed
- RWMutex for LED register, excessive for LED blinking



## [0.0.1] - 2023-06-08 - PRE-RELEASE

### Features
- Vivado project is builadble (no testbench)
- Petalinux project is buildable
- Kernel module to handle FPGA interrupts working
- Tested interrupts handling and LEDs manually
- Detailed and documented Vivado and petalinux guide

### Fixed
- Pins for LEDs were bidir by default from board support file

### Changed
- Vivado project to be cleaner

### Removed
Nothing

### NOTE
Golang application is not ready. It's builadable but not debugged. 
