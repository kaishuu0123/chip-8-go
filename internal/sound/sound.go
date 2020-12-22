package sound

// typedef signed short Int16;
// typedef unsigned char Uint8;
// void SineWave(void *userdata, Uint8 *stream, int len);
import "C"
import (
	"encoding/binary"
	"math"
	"reflect"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	toneHz   = 440.0
	sampleHz = 64 * 60
)

type Audio struct {
	audioDevise  sdl.AudioDeviceID
	desiredSpec  *sdl.AudioSpec
	obtainedSpec *sdl.AudioSpec
}

func wave(t float64) int {
	if t-math.Floor(t) < 0.5 {
		return 1
	} else {
		return -1
	}
}

//export SineWave
func SineWave(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	n := int(length) / 2
	hdr := reflect.SliceHeader{Data: uintptr(unsafe.Pointer(stream)), Len: n, Cap: n}
	buf := *(*[]C.Int16)(unsafe.Pointer(&hdr))
	step := 0

	for i := 0; i < n; i++ {
		buf[i] = C.Int16(wave(float64(step*400)/sampleHz) * 3000)
		step++
	}
}

func InitAudio() *Audio {
	audio := &Audio{}

	audio.desiredSpec = &sdl.AudioSpec{
		Freq:     sampleHz,
		Format:   sdl.AUDIO_F32LSB,
		Channels: 1,
		Samples:  128,
	}

	audio.obtainedSpec = &sdl.AudioSpec{}

	var err error
	if sdl.GetNumAudioDevices(false) > 0 {
		if audio.audioDevise, err = sdl.OpenAudioDevice("", false, audio.desiredSpec, audio.obtainedSpec, sdl.AUDIO_ALLOW_ANY_CHANGE); err != nil {
			panic(err)
		}

		sdl.PauseAudioDevice(audio.audioDevise, false)
	}

	return audio
}

func (audio *Audio) UpdateSound() {
	if audio.audioDevise != 0 {
		sdl.ClearQueuedAudio(audio.audioDevise)

		n := int(audio.obtainedSpec.Channels) * int(audio.obtainedSpec.Samples) * 4
		data := make([]byte, n)

		a := 0.5

		sample := make([]byte, 4)
		step := 1
		for i := 0; i < n; i += 4 {
			// Create sine wave
			s := a * math.Sin(2.0*math.Pi*toneHz*float64(step)/float64(audio.obtainedSpec.Freq))
			binary.LittleEndian.PutUint32(sample, math.Float32bits(float32(s)))
			copy(data[i:], sample)
			step++
		}

		if err := sdl.QueueAudio(audio.audioDevise, data); err != nil {
			println(err)
		}
	}
}

func StopSound() {
	sdl.CloseAudio()
}
