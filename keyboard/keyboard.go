package keyboard

import (
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type Keyboard struct {
	keys     uint16
	mappings [16]sdl.Keycode
}

func New() Keyboard {
	mappings := [16]sdl.Keycode{
		sdl.K_x, sdl.K_1, sdl.K_2, sdl.K_3,
		sdl.K_q, sdl.K_w, sdl.K_e, sdl.K_a,
		sdl.K_s, sdl.K_d, sdl.K_z, sdl.K_c,
		sdl.K_4, sdl.K_r, sdl.K_f, sdl.K_v,
	}
	kb := Keyboard{0, mappings}
	return kb
}

func (kb *Keyboard) PollEvents() (uint8, bool) {
	for sdlEvent := sdl.PollEvent(); sdlEvent != nil; sdlEvent = sdl.PollEvent() {
		switch t := sdlEvent.(type) {
		case *sdl.QuitEvent:
			os.Exit(0)
		case *sdl.KeyboardEvent:
			keyCode := t.Keysym.Sym
			for i := 0; i < 16; i++ {
				if keyCode == kb.mappings[i] {
					mask := uint16(0x1) << i
					if t.Type == sdl.KEYDOWN {
						kb.keys = kb.keys | mask
						return uint8(i), true
					} else if t.Type == sdl.KEYUP {
						kb.keys = kb.keys & (^mask)
					}
				}
			}
		}
	}
	return 0, false
}

func (kb Keyboard) IsPressed(key uint8) bool {
	key = key & 0x0F
	mask := uint16(0x1 << key)

	return kb.keys&mask != 0
}
