package main

import (
	"fmt"
	"os"

	"github.com/Gigalione/golang-chip8/buzzer"
	"github.com/Gigalione/golang-chip8/cpu"
	"github.com/Gigalione/golang-chip8/display"
	"github.com/Gigalione/golang-chip8/keyboard"
	"github.com/Gigalione/golang-chip8/ram"
	"github.com/veandco/go-sdl2/sdl"
)

const CHIP8_CLOCKRATE int = 600

func main() {
	c := cpu.New()
	r := ram.New()

	// Parse command line arguments
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Please provide a path to a ROM")
		return
	}

	// Load ROM into memory
	var romFilepath string = args[0]
	err := r.LoadRom(romFilepath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Initialize SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	// Create display & window
	d := display.New(10, romFilepath)

	kb := keyboard.New()
	bz := buzzer.New()

	// Run program
	for {
		for i := 0; i < CHIP8_CLOCKRATE/60; i++ {
			c.ExecuteInstruction(&r, &d, &kb)
		}

		d.UpdatePicture()
		shouldBuzz := c.DecrementTimers()

		if shouldBuzz && !bz.IsBuzzing {
			bz.Play()
		} else if !shouldBuzz && bz.IsBuzzing {
			bz.Stop()
		}

		kb.PollEvents()

		sdl.Delay(17)
	}
}
