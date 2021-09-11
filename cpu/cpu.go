package cpu

import (
	"encoding/binary"
	"math/rand"
	"time"

	"github.com/Gigalione/golang-chip8/display"
	"github.com/Gigalione/golang-chip8/keyboard"
	"github.com/Gigalione/golang-chip8/ram"
)

type cpu struct {
	vx    [16]uint8 // general purpose registers
	stack []uint16  // stack (stack pointer register is not needed)
	i     uint16    // index register
	delay uint8     // delay timer
	sound uint8     // sound timer
	pc    uint16    // program counter
	rng   rand.Rand // pseudo-random number generator
}

func New() cpu {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // Create the pseudorandom number generator
	cpu := cpu{[16]uint8{}, []uint16{}, 0, 0, 0, 0x200, *rng}
	return cpu
}

// Decrements timers and returns if a buzz should be played
func (c *cpu) DecrementTimers() bool {
	if c.delay > 0 {
		c.delay -= 1
	}

	if c.sound > 0 {
		c.sound -= 1
		return true
	}

	return false
}

func (c *cpu) ExecuteInstruction(r *ram.RAM, d *display.Display, kb *keyboard.Keyboard) {
	var cmd [2]uint8
	var err error

	cmd[0], err = r.Read(c.pc)
	if err != nil {
		panic(err)
	}

	cmd[1], err = r.Read(c.pc + 1)
	if err != nil {
		panic(err)
	}

	c.pc = c.pc + 2
	cmdBin := binary.BigEndian.Uint16(cmd[:]) // this method takes a slice and not an array

	switch cmdBin & 0xF000 {
	case 0x0000:
		// SYS is ignored
		switch cmdBin & 0x00FF {
		case 0x00E0:
			// CLS - Clear the display
			d.ClearDisplay()
		case 0x00EE:
			// RET - Return from a subroutine
			c.ret(cmdBin)
		}
	case 0x1000:
		// JP - Set the program counter register
		c.jp(cmdBin, false)
	case 0x2000:
		// CALL - Call subroutine
		c.call(cmdBin)
	case 0x3000:
		// SE Vx, byte - Compare register value to byte, skip next command if they are equal
		c.se(cmdBin, false)
	case 0x4000:
		// SNE Vx, byte - Compare register value to byte, skip next command if they are NOT equal
		c.sne(cmdBin, false)
	case 0x5000:
		// SE Vx, Vy - Compare register values, skip next command if they are equal
		c.se(cmdBin, true)
	case 0x6000:
		// LD Vx, byte - Load byte into register
		c.ld(cmdBin, false)
	case 0x7000:
		// ADD Vx, byte - Add byte to register
		c.add(cmdBin, false)
	case 0x8000:
		switch cmdBin & 0x000F {
		case 0x0:
			// LD Vx, Vy - Store value of register Vy in Vx
			c.ld(cmdBin, true)
		case 0x1:
			// OR Vx, Vy - Store result of bitwise OR between Vx and Vy registers in Vx
			c.or(cmdBin)
		case 0x2:
			// AND Vx, Vy - Store result of bitwise AND between Vx and Vy registers in Vx
			c.and(cmdBin)
		case 0x3:
			// XOR Vx, Vy - Store result of bitwise XOR between Vx and Vy registers in Vx
			c.xor(cmdBin)
		case 0x4:
			// ADD Vx, Vy - Add Vx and Vy and store in Vx, store carry in VF
			c.add(cmdBin, true)
		case 0x5:
			// SUB Vx, Vy - Subtract Vy from Vx and store in Vx, store NOT borrow in VF
			c.sub(cmdBin)
		case 0x6:
			// SHR Vx - Shift Vx right, store least significant bit in VF
			c.shr(cmdBin)
		case 0x7:
			// SUBN - Subtract Vx from Vy and store in Vx, store NOT borrow in VF
			c.subn(cmdBin)
		case 0xE:
			// SHL - Shift Vx left, store most significant bit in VF
			c.shl(cmdBin)
		}
	case 0x9000:
		// SNE Vx, Vy - Skip next instruction if Vx != Vy
		c.sne(cmdBin, true)
	case 0xA000:
		// LD I, addr - Set value of I to addr
		c.ldi(cmdBin)
	case 0xB000:
		// JP V0, addr - Jump to location addr + V0
		c.jp(cmdBin, true)
	case 0xC000:
		// RND Vx, byte - Generate a pseudo-random number from 0 to 255, bitwise AND the value to the provided byte and store in Vx
		c.rnd(cmdBin)
	case 0xD000:
		// DRW Vx, Vy, nibble - Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision
		c.drw(cmdBin, d, r)
	case 0xE000:
		switch cmdBin & 0x00FF {
		case 0x009E:
			// SKP Vx - Skip next instruction if key with the value Vx is pressed
			c.skp(cmdBin, kb)
		case 0x00A1:
			// SKNP Vx - Skip next instruction if key with the value Vx is not pressed
			c.sknp(cmdBin, kb)
		}
	case 0xF000:
		switch cmdBin & 0x00FF {
		case 0x7:
			// LDD Vx, DT - Load delay timer value into Vx
			c.ldd(cmdBin, false)
		case 0xA:
			// LDK Vx - Wait for a key press and store the key number in Vx
			c.ldk(cmdBin, kb)
		case 0x15:
			// LDD DT, Vx - Store value of Vx into delay timer
			c.ldd(cmdBin, true)
		case 0x18:
			// LDS ST, Vx - Store value of Vx into sound timer
			c.lds(cmdBin)
		case 0x1E:
			// ADDI Vx - Set I = I + Vx
			c.addi(cmdBin)
		case 0x29:
			// FNT Vx - Put address for digit Vx's sprite in I
			c.fnt(cmdBin)
		case 0x33:
			// BCD Vx - Store BCD representation of Vx in memory, starting with address stored in I
			c.bcd(cmdBin, r)
		case 0x55:
			// STR Vx - Store registers V0 - Vx in memory starting from memory address I
			c.str(cmdBin, r)
		case 0x65:
			// LDR Vx - Load registers V0 - Vx from memory starting at address I
			c.ldr(cmdBin, r)
		}
	}

}
