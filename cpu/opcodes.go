package cpu

import (
	"github.com/Gigalione/golang-chip8/display"
	"github.com/Gigalione/golang-chip8/keyboard"
	"github.com/Gigalione/golang-chip8/ram"
	"github.com/veandco/go-sdl2/sdl"
)

func (c *cpu) ret(cmdBin uint16) {
	if len(c.stack) <= 0 {
		return
	}
	stackTop := len(c.stack) - 1
	c.pc = c.stack[stackTop]
	c.stack = c.stack[:stackTop]
}

func (c *cpu) jp(cmdBin uint16, mode bool) {
	jumpAddress := cmdBin & 0x0FFF
	if mode {
		// Add value of V0 to jump address
		jumpAddress += uint16(c.vx[0])
	}
	c.pc = jumpAddress
}

func (c *cpu) call(cmdBin uint16) {
	c.stack = append(c.stack, c.pc)
	jumpAddress := cmdBin & 0x0FFF
	c.pc = jumpAddress
}

func (c *cpu) se(cmdBin uint16, mode bool) {
	if !mode {
		// Compare to a constant
		registerNumber := (cmdBin & 0x0F00) >> 8
		numberToCompare := uint8(cmdBin & 0x00FF)
		if c.vx[registerNumber] == numberToCompare {
			c.pc += 2
		}
	} else {
		// Compare to another register's value
		firstRegisterNumber := (cmdBin & 0x0F00) >> 8
		secondRegisterNumber := (cmdBin & 0x00F0) >> 4
		if c.vx[firstRegisterNumber] == c.vx[secondRegisterNumber] {
			c.pc += 2
		}
	}
}

func (c *cpu) sne(cmdBin uint16, mode bool) {
	if !mode {
		// Compare with constant
		registerNumber := (cmdBin & 0x0F00) >> 8
		numberToCompare := uint8(cmdBin & 0x00FF)
		if c.vx[registerNumber] != numberToCompare {
			c.pc += 2
		}
	} else {
		// Compare with another register
		firstRegisterNumber := (cmdBin & 0x0F00) >> 8
		secondRegisterNumber := (cmdBin & 0x00F0) >> 4
		if c.vx[firstRegisterNumber] != c.vx[secondRegisterNumber] {
			c.pc += 2
		}
	}
}

func (c *cpu) ld(cmdBin uint16, mode bool) {
	if !mode {
		// Load constant
		registerNumber := (cmdBin & 0x0F00) >> 8
		valueToPut := uint8(cmdBin & 0x00FF)
		c.vx[registerNumber] = valueToPut
	} else {
		// Load value from another register
		firstRegisterNumber := (cmdBin & 0x0F00) >> 8
		secondRegisterNumber := (cmdBin & 0x00F0) >> 4
		c.vx[firstRegisterNumber] = c.vx[secondRegisterNumber]
	}

}

func (c *cpu) add(cmdBin uint16, mode bool) {
	if !mode {
		// Add constant
		registerNumber := (cmdBin & 0x0F00) >> 8
		valueToAdd := uint8(cmdBin & 0x00FF)
		c.vx[registerNumber] += valueToAdd
	} else {
		// Add registers together
		reg1 := (cmdBin & 0x0F00) >> 8 // first register number
		reg2 := (cmdBin & 0x00F0) >> 4 // second register number

		sum := uint16(c.vx[reg1]) + uint16(c.vx[reg2])
		result := uint8(sum)
		carry := uint8(sum >> 8)

		c.vx[reg1] = result
		c.vx[0xF] = carry
	}
}

func (c *cpu) or(cmdBin uint16) {
	reg1 := (cmdBin & 0x0F00) >> 8 // first register number
	reg2 := (cmdBin & 0x00F0) >> 4 // second register number
	result := c.vx[reg1] | c.vx[reg2]
	c.vx[reg1] = result
}

func (c *cpu) and(cmdBin uint16) {
	reg1 := (cmdBin & 0x0F00) >> 8 // first register number
	reg2 := (cmdBin & 0x00F0) >> 4 // second register number
	result := c.vx[reg1] & c.vx[reg2]
	c.vx[reg1] = result
}

func (c *cpu) xor(cmdBin uint16) {
	reg1 := (cmdBin & 0x0F00) >> 8 // first register number
	reg2 := (cmdBin & 0x00F0) >> 4 // second register number
	result := c.vx[reg1] ^ c.vx[reg2]
	c.vx[reg1] = result
}

func (c *cpu) sub(cmdBin uint16) {
	reg1 := (cmdBin & 0x0F00) >> 8 // first register number
	reg2 := (cmdBin & 0x00F0) >> 4 // second register number

	result := c.vx[reg1] - c.vx[reg2]

	var borrow uint8
	if c.vx[reg1] > c.vx[reg2] {
		borrow = 1
	} else {
		borrow = 0
	}

	c.vx[reg1] = result
	c.vx[0xF] = borrow
}

func (c *cpu) shr(cmdBin uint16) {
	reg := (cmdBin & 0x0F00) >> 8 // register number

	c.vx[0xF] = c.vx[reg] & 0x0001
	c.vx[reg] = c.vx[reg] >> 1
}

func (c *cpu) subn(cmdBin uint16) {
	reg1 := (cmdBin & 0x0F00) >> 8 // first register number
	reg2 := (cmdBin & 0x00F0) >> 4 // second register number

	result := c.vx[reg2] - c.vx[reg1]

	var borrow uint8
	if c.vx[reg2] > c.vx[reg1] {
		borrow = 1
	} else {
		borrow = 0
	}

	c.vx[reg1] = result
	c.vx[0xF] = borrow
}

func (c *cpu) shl(cmdBin uint16) {
	reg := (cmdBin & 0x0F00) >> 8 // register number

	c.vx[0xF] = (c.vx[reg] & 0x80) >> 7
	c.vx[reg] = c.vx[reg] << 1
}

func (c *cpu) ldi(cmdBin uint16) {
	constantToLoad := cmdBin & 0x0FFF
	c.i = constantToLoad
}

func (c *cpu) rnd(cmdBin uint16) {
	randNumber := uint8(c.rng.Uint32()) // Use the CPU's pseudorandom number generator to get a uint8
	result := randNumber & uint8((cmdBin & 0x00FF))

	reg := (cmdBin & 0x0F00) >> 8 // register number
	c.vx[reg] = result
}

func (c *cpu) ldd(cmdBin uint16, mode bool) {
	reg := (cmdBin & 0x0F00) >> 8 // register number
	if !mode {
		c.vx[reg] = c.delay
	} else {
		c.delay = c.vx[reg]
	}
}

func (c *cpu) lds(cmdBin uint16) {
	reg := (cmdBin & 0x0F00) >> 8 // register number
	c.sound = c.vx[reg]
}

func (c *cpu) addi(cmdBin uint16) {
	reg := (cmdBin & 0x0F00) >> 8 // register number
	c.i = c.i + uint16(c.vx[reg])
}

func (c *cpu) drw(cmdBin uint16, d *display.Display, r *ram.RAM) {
	height := cmdBin & 0x000F
	reg1 := (cmdBin & 0x0F00) >> 8 // first register number
	reg2 := (cmdBin & 0x00F0) >> 4 // second register number
	coordY := c.vx[reg1]
	coordX := c.vx[reg2]

	// read sprite data
	var spriteData []uint8
	for offset := uint16(0); offset < height; offset++ {
		spriteByte, err := r.Read(c.i + offset)
		if err != nil {
			panic(err)
		}
		spriteData = append(spriteData, spriteByte)
	}

	collision := d.Draw(coordX, coordY, spriteData)
	if collision {
		c.vx[0xF] = 1
	} else {
		c.vx[0xF] = 0
	}
}

func (c *cpu) fnt(cmdBin uint16) {
	reg := (cmdBin & 0x0F00) >> 8
	num := c.vx[reg]
	c.i = uint16(num) * 5
}

func (c cpu) bcd(cmdBin uint16, r *ram.RAM) {
	reg := (cmdBin & 0x0F00) >> 8 // register number
	regValue := c.vx[reg]
	var bcdDigit uint8 // current digit only

	for i := 0; i < 3; i++ {
		bcdDigit = regValue % 10
		memAddr := c.i + 2 - uint16(i) // we're moving from the least-significant digit to most-significant
		err := r.Write(memAddr, bcdDigit)
		if err != nil {
			panic(err)
		}
		regValue /= 10 // work with next digit on next iteration
	}
}

func (c cpu) str(cmdBin uint16, r *ram.RAM) {
	lastReg := (cmdBin & 0x0F00) >> 8

	for offset := 0; offset <= int(lastReg); offset++ {
		memAddr := c.i + uint16(offset)
		err := r.Write(memAddr, c.vx[offset])
		if err != nil {
			panic(err)
		}
	}
}

func (c *cpu) ldr(cmdBin uint16, r *ram.RAM) {
	lastReg := (cmdBin & 0x0F00) >> 8

	for reg := 0; reg <= int(lastReg); reg++ {
		memAddr := c.i + uint16(reg)

		var err error
		c.vx[reg], err = r.Read(memAddr)

		if err != nil {
			panic(err)
		}
	}
}

func (c *cpu) skp(cmdBin uint16, kb *keyboard.Keyboard) {
	reg := (cmdBin & 0x0F00) >> 8
	keyToCheck := c.vx[reg]

	if kb.IsPressed(keyToCheck) {
		c.pc += 2
	}
}

func (c *cpu) sknp(cmdBin uint16, kb *keyboard.Keyboard) {
	reg := (cmdBin & 0x0F00) >> 8
	keyToCheck := c.vx[reg]

	if !kb.IsPressed(keyToCheck) {
		c.pc += 2
	}
}

func (c *cpu) ldk(cmdBin uint16, kb *keyboard.Keyboard) {
	reg := (cmdBin & 0x0F00) >> 8

	var keypressRegistered bool = false
	var btn uint8
	for !keypressRegistered {
		btn, keypressRegistered = kb.PollEvents()
		sdl.Delay(17)
	}

	c.vx[reg] = btn
}
