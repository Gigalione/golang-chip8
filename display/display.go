package display

import (
	"path/filepath"

	"github.com/veandco/go-sdl2/sdl"
)

type Display struct {
	displayData [32][64]bool // 64 * 32
	scale       int32
	window      *sdl.Window
	renderer    *sdl.Renderer
}

func New(scale int32, romFilepath string) Display {
	filename := filepath.Base(romFilepath) // displayed in window title

	// Create a window
	window, err := sdl.CreateWindow(
		"Chip-8 Emulator | "+filename,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		64*scale, 32*scale,
		sdl.WINDOW_OPENGL,
	)
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	renderer.Clear()

	d := Display{[32][64]bool{}, scale, window, renderer}
	return d
}

func (d *Display) ClearDisplay() {
	d.displayData = [32][64]bool{{false}}
}

func (d *Display) Draw(coordX uint8, coordY uint8, spriteData []uint8) bool {
	var collision bool = false

	height := len(spriteData)
	for i := 0; i < height; i++ {
		x := (coordX + uint8(i)) % 32
		mask := uint8(0x80)

		for j := 0; j < 8; j++ {
			y := (coordY + uint8(j)) % 64

			if spriteData[i]&mask != 0 {
				d.displayData[x][y] = !d.displayData[x][y]
				if !d.displayData[x][y] {
					collision = true
				}
			}

			mask = mask >> 1
		}
	}

	return collision
}

func (d Display) UpdatePicture() {
	d.renderer.SetDrawColor(0, 0, 0, 0)
	d.renderer.Clear()

	for i := 0; i < 32; i++ {
		for j := 0; j < 64; j++ {
			if d.displayData[i][j] {
				d.renderer.SetDrawColor(255, 255, 255, 255)
				d.renderer.FillRect(&sdl.Rect{
					Y: int32(i) * d.scale,
					X: int32(j) * d.scale,
					W: d.scale,
					H: d.scale,
				})
			}
		}
	}

	d.renderer.Present()
}
