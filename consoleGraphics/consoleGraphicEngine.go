package consoleGraphics

import (
	"fmt"
)

const (
	// COLOUR
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
	BLUE   = "\033[34m"
	PURPLE = "\033[35m"
	CYAN   = "\033[36m"
	WHITE  = "\033[37m"
)

// Unicode table:
// https://en.wikipedia.org/wiki/List_of_Unicode_characters
// https://www.branah.com/unicode-converter
const (
	// PIXEL_TYPE
	FULL_BLOCK   = "\u2588"
	DARK_BLOCK   = "\u2593"
	MEDIUM_BLOCK = "\u2592"
	LIGHT_BLOCK  = "\u2591"
	SPACE_BLOCK  = " "
)

type pixel struct {
	pixel_type string
	color      string
}

type consoleComponent interface {
	onCreate() bool
	onUpdate() bool
}

type consoleGraphicEngine struct {
	screenWidth  int
	screenHeight int
	pixels       []pixel
	output       string
	component    consoleComponent
}

func constructConsoleGraphicEngine(width int, height int, color string) *consoleGraphicEngine {
	CGE := &consoleGraphicEngine{}
	CGE.screenWidth = width
	CGE.screenHeight = height
	pix_len := width * height
	pixels := make([]pixel, pix_len)
	for index := 0; index < pix_len; index++ {
		pixels[index].pixel_type = SPACE_BLOCK
		pixels[index].color = color
	}
	CGE.pixels = pixels
	return CGE
}

func (CGE *consoleGraphicEngine) fillALL(pixel_type string, color string) {
	for index := 0; index < len(CGE.pixels); index++ {
		CGE.pixels[index].pixel_type = pixel_type
		CGE.pixels[index].color = color
	}
}

func (CGE *consoleGraphicEngine) drawPixel(x int, y int, pix_type string, pix_color string) {
	if x >= 0 && x < CGE.screenWidth && y >= 0 && y < CGE.screenHeight {
		target := y*CGE.screenWidth + x
		CGE.pixels[target].color = pix_color
		CGE.pixels[target].pixel_type = pix_type
	}
}

func (CGE *consoleGraphicEngine) drawLine(x1 int, y1 int, x2 int, y2 int, pix_type string, pix_color string) {
	var x, y, xe, ye int
	dx := x2 - x1
	dy := y2 - y1
	dx1 := abs(dx)
	dy1 := abs(dy)
	px := 2*dy1 - dx1
	py := 2*dx1 - dy1
	if dy1 <= dx1 {
		if dx >= 0 {
			x = x1
			y = y1
			xe = x2
		} else {
			x = x2
			y = y2
			xe = x1
		}
		go CGE.drawPixel(x, y, pix_type, pix_color)

		for i := 0; x < xe; i++ {
			x++
			if px < 0 {
				px += 2 * dy1
			} else {
				if (dx < 0 && dy < 0) || (dx > 0 && dy > 0) {
					y++
				} else {
					y--
				}
				px += 2 * (dy1 - dx1)
			}
			go CGE.drawPixel(x, y, pix_type, pix_color)
		}
	} else {
		if dy >= 0 {
			x = x1
			y = y1
			ye = y2
		} else {
			x = x2
			y = y2
			ye = y1
		}
		go CGE.drawPixel(x, y, pix_type, pix_color)
		for i := 0; y < ye; i++ {
			y++
			if py <= 0 {
				py += 2 * dx1
			} else {
				if (dx < 0 && dy < 0) || (dx > 0 && dy > 0) {
					x++
				} else {
					x--
				}
				py += 2 * (dx1 - dy1)
			}
			go CGE.drawPixel(x, y, pix_type, pix_color)
		}
	}
}

func abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

func (CGE *consoleGraphicEngine) drawTriangle(x1 int, y1 int, x2 int, y2 int, x3 int, y3 int, pix_type string, pix_color string) {
	go CGE.drawLine(x1, y1, x2, y2, pix_type, pix_color)
	go CGE.drawLine(x2, y2, x3, y3, pix_type, pix_color)
	go CGE.drawLine(x3, y3, x1, y1, pix_type, pix_color)
}

func (CGE *consoleGraphicEngine) computeGraphics() {
	var out string = "\033[0;0H"
	for index, pix := range CGE.pixels {
		out += pix.color + pix.pixel_type
		if (index+1)%CGE.screenWidth == 0 {
			out += "\n"
		}
	}
	out += "\n"
	CGE.output = out
}

func (CGE *consoleGraphicEngine) render() {
	fmt.Print(CGE.output)
}

func (CGE *consoleGraphicEngine) Start() {
	CGE.component.onCreate()

	for {
		CGE.component.onUpdate()
		CGE.computeGraphics()
		CGE.render()
	}

}

func (CGE *consoleGraphicEngine) addComponent(new consoleComponent) {
	CGE.component = new
}

// func main() {
// 	conRen := constructConsoleGraphicEngine(300, 100, WHITE)

// 	conRen.computeGraphics()
// 	conRen.render()

// 	conRen.drawTriangle(3, 3, 250, 3, 3, 99, FULL_BLOCK, RED)

// 	conRen.computeGraphics()
// 	conRen.render()

// 	return

// 	// Zoom terminal out to smallest by [ctrl] + [-]
// }
