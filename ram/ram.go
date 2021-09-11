package ram

import (
	"errors"
	"io/ioutil"
)

type RAM struct {
	mem [3584]uint8 // first 512 bytes of the 4096 total are reserved for the interpreter on a real machine
	fnt [80]uint8   // font data
}

func New() RAM {
	r := RAM{[3584]uint8{0}, [80]uint8{0}}

	// load font data
	var buff, err = ioutil.ReadFile("./ram/font.bin")

	if err != nil {
		panic("can't read font file at ./ram/font.bin")
	}

	for i := 0; i < 80; i++ {
		r.fnt[i] = buff[i]
	}

	return r
}

func (r *RAM) LoadRom(fileName string) error {
	var buff, err = ioutil.ReadFile(fileName)

	if err != nil {
		return err
	}

	//r.mem = [3584]uint8{0}

	len := len(buff)
	if len > 3584 {
		return errors.New("ROM file is too big")
	}

	for i := 0; i < len; i++ {
		r.mem[i] = buff[i]
	}

	return nil
}

func (r RAM) Read(addr uint16) (uint8, error) {
	if addr >= 0x0 && addr < 0x050 {
		return r.fnt[addr], nil
	} else if addr >= 0x50 && addr < 0x200 {
		return 0, nil
	} else if addr >= 0x200 && addr < 0xFFF {
		return r.mem[addr-512], nil
	} else {
		return 0, errors.New("out of memory bounds")
	}
}

func (r *RAM) Write(addr uint16, data uint8) error {
	if addr >= 0x0 && addr < 0x200 {
		return errors.New("attempt to write to interpreter-reserved memory")
	} else if addr >= 0x200 && addr < 0xFFF {
		r.mem[addr-512] = data
		return nil
	} else {
		return errors.New("out of memory bounds")
	}
}
