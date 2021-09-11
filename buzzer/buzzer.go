package buzzer

// typedef unsigned char Uint8;
// void SquareWave(void *userdata, Uint8 *stream, int len);
import "C"

import (
	"math"
	"reflect"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	toneFreq   = 440
	sampleRate = 48000
)

//export SquareWave
func SquareWave(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length)
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Uint8)(unsafe.Pointer(&hdr))

	const delta = 2 * math.Pi * toneFreq / sampleRate
	var phase float64
	for i := 0; i < n; i++ {
		phase += delta

		sample := C.Uint8(0)
		if !math.Signbit(math.Sin(phase)) {
			sample = C.Uint8(127)
		}

		buf[i] = sample
	}
}

type Buzzer struct {
	IsBuzzing   bool
	audioDevice sdl.AudioDeviceID
}

func New() Buzzer {
	spec := &sdl.AudioSpec{
		Freq:     sampleRate,
		Format:   sdl.AUDIO_U8,
		Channels: 1,
		Samples:  128,
		Callback: sdl.AudioCallback(C.SquareWave),
	}

	device, err := sdl.OpenAudioDevice("", false, spec, nil, 0)
	if err != nil {
		panic(err)
	}

	bz := Buzzer{false, device}
	return bz
}

func (bz *Buzzer) Play() {
	sdl.PauseAudioDevice(bz.audioDevice, false)
	bz.IsBuzzing = true
}

func (bz *Buzzer) Stop() {
	sdl.PauseAudioDevice(bz.audioDevice, true)
	bz.IsBuzzing = false
}
