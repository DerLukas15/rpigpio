package rpigpio

import (
	"time"

	"github.com/DerLukas15/rpimemmap"
	"github.com/pkg/errors"
)

//Pin defines a GPIO pin
type Pin struct {
	pinNum uint32
}

//NewPin creates a new GPIO pin. Do not user physical pin number!
func NewPin(pinNum uint32) (*Pin, error) {
	if pinNum > 53 {
		return nil, errors.Wrap(ErrWrongPin, "NewPin")
	}
	return &Pin{
		pinNum: pinNum,
	}, nil
}

//Mode sets the mode of the GPIO pin
func (p *Pin) Mode(targetMode Mode) error {
	if gpioRegisterMem == nil {
		return errors.Wrap(ErrNotInitialized, "pin mode")
	}
	//Need to shift as there are only 10 GPIOs per Register. Each GPIO has 3 bits per Mode
	shift := p.pinNum % 10 * 3
	memoryOffset := registerOffsetMode + p.pinNum/10*4                                 // 10 GPIOs per Register, 4 bytes per Register
	*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) &= ^(7 << shift)                 // Set all 000
	*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) |= (uint32(targetMode) << shift) // Set actual mode
	return nil
}

//Set sets the value of the GPIO pin
func (p *Pin) Set(targetValue int) error {
	if gpioRegisterMem == nil {
		return errors.Wrap(ErrNotInitialized, "pin set")
	}
	memoryOffset := registerOffsetClr + p.pinNum/32*4 // 32 GPIOs per Register, 4 bytes per Register
	if targetValue != 0 {
		memoryOffset = registerOffsetSet + p.pinNum/32*4 // 32 GPIOs per Register, 4 bytes per Register
	}
	*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) |= (1 << p.pinNum % 32)
	return nil
}

//Get returns if the pin is high
func (p *Pin) Get() (bool, error) {
	if gpioRegisterMem == nil {
		return false, errors.Wrap(ErrNotInitialized, "pin get")
	}
	memoryOffset := registerOffsetLevel + p.pinNum/32*4 // 32 GPIOs per Register, 4 bytes per Register
	res := *(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) >> (p.pinNum % 32)
	return (res & 1) != 0, nil
}

//Pull sets the Pull-Up or Pull-Down resistors of the GPIO pin
func (p *Pin) Pull(targetMode PullMode) error {
	if gpioRegisterMem == nil {
		return errors.Wrap(ErrNotInitialized, "pin pull")
	}
	memoryOffsetGppudClk := registerOffsetGppudClk + p.pinNum/32*4                // 32 GPIOs per Register, 4 bytes per Register
	*(rpimemmap.Reg32(gpioRegisterMem, registerOffsetGppud)) = uint32(targetMode) // Set type of pull mode
	time.Sleep(2 * time.Microsecond)
	*(rpimemmap.Reg32(gpioRegisterMem, memoryOffsetGppudClk)) = 1 << (p.pinNum % 32) // Set which pin should use the defined pull mode
	time.Sleep(2 * time.Microsecond)
	*(rpimemmap.Reg32(gpioRegisterMem, registerOffsetGppud)) = 0 // Set type of pull mode to 0
	return nil
}

//Event returns if there was / is an event
func (p *Pin) Event() (bool, error) {
	if gpioRegisterMem == nil {
		return false, errors.Wrap(ErrNotInitialized, "pin event")
	}
	memoryOffset := registerOffsetEventDetect + p.pinNum/32*4 // 32 GPIOs per Register, 4 bytes per Register
	res := *(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) >> (p.pinNum % 32)
	if !NoEventClearing {
		*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) = 1 << (p.pinNum % 32)
	}
	return (res & 1) != 0, nil
}

//ClearEvent clears a possible event
func (p *Pin) ClearEvent() error {
	if gpioRegisterMem == nil {
		return errors.Wrap(ErrNotInitialized, "pin clear event")
	}
	memoryOffset := registerOffsetEventDetect + p.pinNum/32*4 // 32 GPIOs per Register, 4 bytes per Register
	*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) = 1 << (p.pinNum % 32)
	return nil
}

//RisingEdgeDetect sets the rising edge detection of the GPIO pin
func (p *Pin) RisingEdgeDetect(targetStatus bool) error {
	if gpioRegisterMem == nil {
		return errors.Wrap(ErrNotInitialized, "pin rising detect")
	}
	memoryOffset := registerOffsetRisingEdge + p.pinNum/32*4 // 32 GPIOs per Register, 4 bytes per Register
	if targetStatus {
		*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) |= (1 << p.pinNum % 32)
	} else {
		*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) &= ^(1 << p.pinNum % 32)
	}
	return nil
}

//FallingEdgeDetect sets the rising edge detection of the GPIO pin
func (p *Pin) FallingEdgeDetect(targetStatus bool) error {
	if gpioRegisterMem == nil {
		return errors.Wrap(ErrNotInitialized, "pin falling detect")
	}
	memoryOffset := registerOffsetFallingEdge + p.pinNum/32*4 // 32 GPIOs per Register, 4 bytes per Register
	if targetStatus {
		*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) |= (1 << p.pinNum % 32)
	} else {
		*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) &= ^(1 << p.pinNum % 32)
	}
	return nil
}

//ARisingEdgeDetect sets the rising edge detection of the GPIO pin. Async
func (p *Pin) ARisingEdgeDetect(targetStatus bool) error {
	if gpioRegisterMem == nil {
		return errors.Wrap(ErrNotInitialized, "pin async rising detect")
	}
	memoryOffset := registerOffsetARisingEdge + p.pinNum/32*4 // 32 GPIOs per Register, 4 bytes per Register
	if targetStatus {
		*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) |= (1 << p.pinNum % 32)
	} else {
		*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) &= ^(1 << p.pinNum % 32)
	}
	return nil
}

//AFallingEdgeDetect sets the rising edge detection of the GPIO pin. Async
func (p *Pin) AFallingEdgeDetect(targetStatus bool) error {
	if gpioRegisterMem == nil {
		return errors.Wrap(ErrNotInitialized, "pin async falling detect")
	}
	memoryOffset := registerOffsetAFallingEdge + p.pinNum/32*4 // 32 GPIOs per Register, 4 bytes per Register
	if targetStatus {
		*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) |= (1 << p.pinNum % 32)
	} else {
		*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) &= ^(1 << p.pinNum % 32)
	}
	return nil
}

//HighDetect sets the high detection of the GPIO pin
func (p *Pin) HighDetect(targetStatus bool) error {
	if gpioRegisterMem == nil {
		return errors.Wrap(ErrNotInitialized, "pin high detect")
	}
	memoryOffset := registerOffsetHighDetect + p.pinNum/32*4 // 32 GPIOs per Register, 4 bytes per Register
	if targetStatus {
		*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) |= (1 << p.pinNum % 32)
	} else {
		*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) &= ^(1 << p.pinNum % 32)
	}
	return nil
}

//LowDetect sets the low detection of the GPIO pin
func (p *Pin) LowDetect(targetStatus bool) error {
	if gpioRegisterMem == nil {
		return errors.Wrap(ErrNotInitialized, "pin low detect")
	}
	memoryOffset := registerOffsetLowDetect + p.pinNum/32*4 // 32 GPIOs per Register, 4 bytes per Register
	if targetStatus {
		*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) |= (1 << p.pinNum % 32)
	} else {
		*(rpimemmap.Reg32(gpioRegisterMem, memoryOffset)) &= ^(1 << p.pinNum % 32)
	}
	return nil
}
