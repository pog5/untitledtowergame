package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 800
	windowHeight = 600
	blockWidth   = 50
	blockHeight  = 50
	craneSpeed   = 3
)

var craneX int32 = platformX * 2
var craneWidth int32 = 100
var comboMultiplier int
var comboExpireTime time.Time
var missedAttempts int
var moveRight bool = true
var droppingBlock bool = false
var platformWidth int32 = 200
var platformHeight int32 = 50
var platformX int32 = windowWidth/2 - platformWidth/2
var platformY int32 = windowHeight - 100
var platformCenterX int32 = platformX + platformWidth/2

var blockX = platformX + platformWidth/2

var renderer *sdl.Renderer

type Block struct {
	ID          int32 `json:"id"`
	PersonCount int32 `json:"personCount"`
	X           int32 `json:"x"`
	Y           int32 `json:"y"`
}

var blocks = []Block{}

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("SDL initialization failed:", err)
		fmt.Println("Expect instability!")
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Untitled Tower Game", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE|sdl.WINDOW_ALLOW_HIGHDPI)
	if err != nil {
		fmt.Println("Window creation failed:", err)
		fmt.Println("Expect instability!")
		return
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		fmt.Println("Renderer creation failed:", err)
		fmt.Println("Expect instability!")
		return
	}
	defer renderer.Destroy()
	surface, err := window.GetSurface()
	if err != nil {
		fmt.Println("Surface creation failed:", err)
		fmt.Println("Expect instability!")
		return
	}
	defer surface.Free()
	gameLoop(surface, renderer, window)
}

func gameLoop(surface *sdl.Surface, renderer *sdl.Renderer, window *sdl.Window) {

	for {
		// TODO: Add title screen
		// TODO: Add city building mode
		// TODO: set existing code as freeplay/quick build mode
		processInput(surface, renderer, window)
		update(surface, renderer, window)

		render(surface, renderer, window)

		if missedAttempts >= 3 {
			blocks = []Block{}
			missedAttempts = 0
		}

		if len(blocks) > 10 {
			fmt.Println("You won!")
			return
		}

		sdl.Delay(5)
	}
}

func processInput(surface *sdl.Surface, renderer *sdl.Renderer, window *sdl.Window) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event := event.(type) {
		case *sdl.QuitEvent:
			sdl.Quit()
		case *sdl.KeyboardEvent:
			keyEvent := event
			if keyEvent.Type == sdl.KEYDOWN {
				if keyEvent.Keysym.Sym == sdl.K_SPACE {

					if abs(blockX-craneX) <= 20 {
						handlePerfectHit()
						buildBlock(surface, renderer, window)
					} else if abs(blockX-craneX) <= 50 {
						buildBlock(surface, renderer, window)
					} else {
						missedAttempts++
					}

				}
			}
		}
	}
}

func update(surface *sdl.Surface, renderer *sdl.Renderer, window *sdl.Window) {
	if time.Now().After(comboExpireTime) {
		comboMultiplier = 1
	}
	if moveRight {
		craneX += craneSpeed
		if craneX >= platformCenterX+200 {
			moveRight = false
		}
	} else {
		craneX -= craneSpeed
		if craneX <= platformCenterX-200 {
			moveRight = true
		}
	}

	sdl.Delay(5)
}

func render(surface *sdl.Surface, renderer *sdl.Renderer, window *sdl.Window) {
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.Clear()
	drawBackground(surface, renderer, window)
	drawBase(surface, renderer, window)
	if droppingBlock {
		rect := sdl.Rect{X: craneX - craneWidth/8, Y: craneWidth, W: 50, H: 50}
		colour := sdl.Color{R: 0, G: 0, B: 0, A: 255}
		pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
		surface.FillRect(&rect, pixel)
	}
	drawCrane(surface, renderer, window)
	drawOldBlocks(surface, renderer, window)
	if comboMultiplier > 1 {
		drawComboBar(surface, renderer, window)
	}
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		fmt.Println("Texture creation from surface failed:", err)
		return
	}
	defer texture.Destroy()
	renderer.Copy(texture, nil, nil)
	renderer.Present()
}

func handlePerfectHit() {
	comboMultiplier++
	comboExpireTime = time.Now().Add(time.Second * time.Duration(10-comboMultiplier))
}

func drawBackground(surface *sdl.Surface, renderer *sdl.Renderer, window *sdl.Window) {
	bg := sdl.Rect{X: 0, Y: 0, W: windowWidth, H: windowHeight}
	bgcolor := sdl.Color{R: 139, G: 233, B: 253, A: 255}
	bgpixel := sdl.MapRGBA(surface.Format, bgcolor.R, bgcolor.G, bgcolor.B, bgcolor.A)
	surface.FillRect(&bg, bgpixel)
}

func drawBase(surface *sdl.Surface, renderer *sdl.Renderer, window *sdl.Window) {
	rect := sdl.Rect{X: platformX, Y: platformY, W: platformWidth, H: platformHeight}
	colour := sdl.Color{R: 80, G: 250, B: 123, A: 255}
	pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	surface.FillRect(&rect, pixel)
}

func drawCrane(surface *sdl.Surface, renderer *sdl.Renderer, window *sdl.Window) {
	rect := sdl.Rect{X: craneX, Y: 0, W: 25, H: craneWidth}
	switch missedAttempts {
	case 0:
		colour := sdl.Color{R: 98, G: 114, B: 164, A: 255}
		pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
		surface.FillRect(&rect, pixel)
	case 1:
		colour := sdl.Color{R: 255, G: 255, B: 0, A: 255}
		pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
		surface.FillRect(&rect, pixel)
	case 2:
		colour := sdl.Color{R: 255, G: 0, B: 0, A: 255}
		pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
		surface.FillRect(&rect, pixel)
	}
	// draw a attached block which is 25x25 on the crane
	rect = sdl.Rect{X: craneX - craneWidth/8, Y: craneWidth, W: 50, H: 50}
	colour := sdl.Color{R: 0, G: 0, B: 0, A: 255}
	pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
	surface.FillRect(&rect, pixel)
}

func drawOldBlocks(surface *sdl.Surface, renderer *sdl.Renderer, window *sdl.Window) {
	for i := 0; i < len(blocks); i++ {
		rect := sdl.Rect{X: blocks[i].X, Y: blocks[i].Y, W: blockWidth, H: blockHeight}
		colour := sdl.Color{R: 0, G: 0, B: 0, A: 255}
		pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
		surface.FillRect(&rect, pixel)
	}
}

func buildBlock(surface *sdl.Surface, renderer *sdl.Renderer, window *sdl.Window) {
	if len(blocks) > 0 {
		for i := 0; i < len(blocks); i++ {
			blocks[i].Y -= blockHeight
		}
	}
	blocks = append(blocks, Block{ID: int32(len(blocks)), PersonCount: rand.Int31n(100), X: craneX - craneWidth/8, Y: platformY - blockHeight})
}

func drawComboBar(surface *sdl.Surface, renderer *sdl.Renderer, window *sdl.Window) {
	combobaroutline := sdl.Rect{X: platformCenterX - 80, Y: 10, W: 200, H: 25}
	combobaroutlinecolor := sdl.Color{R: 255, G: 184, B: 108, A: 255}
	combobaroutlinepixel := sdl.MapRGBA(surface.Format, combobaroutlinecolor.R, combobaroutlinecolor.G, combobaroutlinecolor.B, combobaroutlinecolor.A)
	surface.FillRect(&combobaroutline, combobaroutlinepixel)

	combobar := sdl.Rect{X: platformCenterX - 75, Y: 15, W: int32(comboMultiplier) * 10, H: 15}
	combobarcolor := sdl.Color{R: 241, G: 250, B: 140, A: 255}
	combobarpixel := sdl.MapRGBA(surface.Format, combobarcolor.R, combobarcolor.G, combobarcolor.B, combobarcolor.A)
	surface.FillRect(&combobar, combobarpixel)

	// TODO: Add (xN) white text next to the combo bar using SDL_ttf
}

func abs(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
