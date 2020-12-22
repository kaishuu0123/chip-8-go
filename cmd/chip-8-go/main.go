package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/kaishuu0123/chip-8-go/internal/sound"
	"github.com/kaishuu0123/chip-8-go/internal/virtualmachine"
	"github.com/kaishuu0123/chip-8-go/internal/window"
	"github.com/veandco/go-sdl2/sdl"
)

const WIDTH int32 = 640
const HEIGHT int32 = 320

var KeyMap = map[sdl.Scancode]uint{
	sdl.SCANCODE_1: 0x1,
	sdl.SCANCODE_2: 0x2,
	sdl.SCANCODE_3: 0x3,
	sdl.SCANCODE_4: 0xC,
	sdl.SCANCODE_Q: 0x4,
	sdl.SCANCODE_W: 0x5,
	sdl.SCANCODE_E: 0x6,
	sdl.SCANCODE_R: 0xD,
	sdl.SCANCODE_A: 0x7,
	sdl.SCANCODE_S: 0x8,
	sdl.SCANCODE_D: 0x9,
	sdl.SCANCODE_F: 0xE,
	sdl.SCANCODE_Z: 0xA,
	sdl.SCANCODE_X: 0x0,
	sdl.SCANCODE_C: 0xB,
	sdl.SCANCODE_V: 0xF,
}

func main() {
	var w window.Window
	var running bool = true
	// var borderColor = window.Color{R: 50, G: 50, B: 50}
	var foreGroundColor = window.Color{R: 156, G: 220, B: 254}
	var backGroundColor = window.Color{R: 50, G: 50, B: 54}

	var vm *virtualmachine.VirtualMachine
	flag.Parse()
	if filePath := flag.Arg(0); filePath != "" {
		vm, _ = virtualmachine.LoadFromFile(filePath)
	} else {
		vm, _ = virtualmachine.LoadROM(virtualmachine.Boot, false)
	}

	w.Setup(WIDTH, HEIGHT)

	clockTimer := time.NewTicker(time.Second / 500)
	videoTimer := time.NewTicker(time.Second / 60)
	delayTimer := time.NewTicker(time.Second / 60)
	soundTimer := time.NewTicker(time.Second / 60)

	audio := sound.InitAudio()

	for running {
		running = processEvents(vm)

		select {
		case <-videoTimer.C:
			for y := 0; y < int(HEIGHT); y++ {
				for x := 0; x < int(WIDTH); x++ {
					videoX := x / 10
					videoY := y / 10

					if vm.Video[videoY][videoX] == 0 {
						w.SetPixel(x, y, backGroundColor)
					} else {
						w.SetPixel(x, y, foreGroundColor)
					}
				}
			}
			w.Update()
			w.Render()
		case <-clockTimer.C:
			vm.Step()
			// debugPrint(vm)
		case <-delayTimer.C:
			if vm.DT > 0 {
				vm.DT--
			}
		case <-soundTimer.C:
			if vm.ST > 0 {
				vm.ST--

				audio.UpdateSound()
			}
		}
	}
}

func processEvents(vm *virtualmachine.VirtualMachine) bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			if t.Type == sdl.KEYUP {
				if key, ok := KeyMap[t.Keysym.Scancode]; ok {
					vm.ReleasedKey(key)
				}
			} else {
				if key, ok := KeyMap[t.Keysym.Scancode]; ok {
					vm.PressKey(key)
				}
			}
		}
	}

	return true
}

func debugPrint(vm *virtualmachine.VirtualMachine) {
	for i := 0; i < 16; i++ {
		fmt.Printf("V%X = %02X ", i, vm.V[i])
	}
	fmt.Println("")
	fmt.Printf("DT = %0X ", vm.DT)
	fmt.Printf("ST = %0X ", vm.ST)
	fmt.Printf("SP = %0X ", vm.SP)
	fmt.Printf("PC = %0X ", vm.PC)
	fmt.Println("")
}
