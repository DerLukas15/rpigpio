//Package rpigpio grants access to the GPIOs on a Raspberry Pi completly written in GO.
//All actions write directly to the registers.
package rpigpio

import (
	"os"

	"github.com/DerLukas15/rpihardware"
	"github.com/DerLukas15/rpimemmap"
	"github.com/pkg/errors"
)

//See https://www.raspberrypi.org/app/uploads/2012/02/BCM2835-ARM-Peripherals.pdf for a description of the registers.

const (
	registerBusOffset uint32 = 0x00200000

	registerOffsetMode         uint32 = 0x00
	registerOffsetSet          uint32 = 0x1c
	registerOffsetClr          uint32 = 0x28
	registerOffsetLevel        uint32 = 0x34
	registerOffsetEventDetect  uint32 = 0x40
	registerOffsetRisingEdge   uint32 = 0x4c
	registerOffsetFallingEdge  uint32 = 0x58
	registerOffsetHighDetect   uint32 = 0x64
	registerOffsetLowDetect    uint32 = 0x70
	registerOffsetARisingEdge  uint32 = 0x7c
	registerOffsetAFallingEdge uint32 = 0x88
	registerOffsetGppud        uint32 = 0x94
	registerOffsetGppudClk     uint32 = 0x98
)

//Mode is used to set the mode of a pin
type Mode uint32

//Defined modes
const (
	ModeIn         Mode = 0b000
	ModeOut        Mode = 0b001
	ModeAlternate0 Mode = 0b100
	ModeAlternate1 Mode = 0b101
	ModeAlternate2 Mode = 0b110
	ModeAlternate3 Mode = 0b111
	ModeAlternate4 Mode = 0b011
	ModeAlternate5 Mode = 0b010
)

//PullMode is used to set the Pull-Up and Pull-Down resistors
type PullMode uint32

//Defined pull modes
const (
	PullOff  PullMode = 0b00
	PullDown PullMode = 0b01
	PullUp   PullMode = 0b10
)

// Errors
var (
	ErrNotInitialized = errors.New("Package not initialized or unsupported hardware. Check error during initialize.")
	ErrWrongPin       = errors.New("Pin is not defined. Max 53")
)

var NoEventClearing bool
var curHardware *rpihardware.Hardware
var gpioRegisterMem rpimemmap.MemMap

func Initialize() error {
	if gpioRegisterMem != nil {
		//Already initialized
		return nil
	}
	var err error
	NoEventClearing = false
	curHardware, err = rpihardware.Check()
	if err != nil {
		return err
	}
	gpioRegisterMem = rpimemmap.NewPeripheral(uint32(os.Getpagesize()))
	err = gpioRegisterMem.Map(registerBusOffset, rpimemmap.MemDevDefault, 0)
	if err != nil {
		return err
	}
	return nil
}
