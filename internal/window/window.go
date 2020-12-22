package window

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Color struct {
	R byte
	G byte
	B byte
}

type Window struct {
	width           int32
	height          int32
	window          *sdl.Window
	texture         *sdl.Texture
	renderer        *sdl.Renderer
	pixels          []byte
	backGroundColor Color
	foreGroundColor Color
	font            *ttf.Font
	fontTexture     *sdl.Texture
	text            *sdl.Surface
}

func (window *Window) Setup(width int32, height int32) {
	window.width = width
	window.height = height

	var flags uint32 = sdl.WINDOW_SHOWN

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}

	window.window, err = sdl.CreateWindow("chip-8-go",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		window.width, window.height, flags)
	if err != nil {
		panic(err)
	}

	window.renderer, err = sdl.CreateRenderer(window.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	window.texture, err = window.renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING,
		window.width, window.height)
	if err != nil {
		panic(err)
	}

	window.pixels = make([]byte, window.width*window.height*4)
}

func (window *Window) SetPixel(x int, y int, c Color) {
	index := (y*int(window.width) + x) * 4

	if index < len(window.pixels)-4 && index >= 0 {
		window.pixels[index] = c.R
		window.pixels[index+1] = c.G
		window.pixels[index+2] = c.B
	}
}

func (window *Window) Update() {
	window.texture.Update(nil, window.pixels, int(window.width*4))
}

func (window *Window) Render() {
	window.renderer.Clear()

	window.renderer.Copy(window.texture, nil, nil)

	window.renderer.Present()
}
