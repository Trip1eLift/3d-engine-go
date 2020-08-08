package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

// Main Reference:
// https://kylewbanks.com/blog/tutorial-opengl-with-golang-part-1-hello-opengl
// https://www.youtube.com/watch?v=iyDJ_1lElms

type RGB struct {
	red, green, blue float32
}

const (
	// input vec3 as position
	vertexShaderSource = `
		#version 410
		in vec2 vertice_point;
		in vec4 vertice_color;
		out vec4 color;
		void main() {
			color = vertice_color;
			gl_Position = vec4(vertice_point, 0.0, 1.0);
		}
	` + "\x00"
	// output colour
	fragmentShaderSource = `
		#version 410
		in vec4 color;
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(color);
		}
	` + "\x00"
)

var (
	WHITE = RGB{1.0, 1.0, 1.0}
	RED   = RGB{1.0, 0.0, 0.0}
	GREEN = RGB{0.0, 1.0, 0.0}
	BLUE  = RGB{0.0, 0.0, 1.0}
)

type element interface {
	onCreate() bool
	onUpdate() bool
}

type openglGraphicsEngine struct {
	screenWidth  int
	screenHeight int
	title        string
	window       *glfw.Window
	program      uint32
	elements     []element
	delta        float64
	fps          float64
}

func constructOpenglGraphicsEngine(width int, height int, title string, targetFPS float64) *openglGraphicsEngine {
	OGE := &openglGraphicsEngine{}
	OGE.screenWidth = width
	OGE.screenHeight = height
	OGE.title = title
	OGE.fps = targetFPS

	runtime.LockOSThread()

	OGE.window = initGlfw(OGE.screenWidth, OGE.screenHeight, OGE.title)
	OGE.program = initOpenGL()
	return OGE
}

func (OGE *openglGraphicsEngine) destructor() {
	defer glfw.Terminate()
}

func (OGE *openglGraphicsEngine) DrawTriangle(vertices []float32, colors []float32) {
	vao := makeVao(vertices, colors)
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
}

func (OGE *openglGraphicsEngine) FillTriangle(vertices []float32, colors []float32) {
	vao := makeVao(vertices, colors)
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
}

func (OGE *openglGraphicsEngine) ScaleRGB(color RGB, factor float32) RGB {
	color.red *= factor
	if color.red > 1.0 {
		color.red = 1.0
	}
	color.green *= factor
	if color.green > 1.0 {
		color.green = 1.0
	}
	color.blue *= factor
	if color.blue > 1.0 {
		color.blue = 1.0
	}
	return color
}

func (OGE *openglGraphicsEngine) TriangleColorEvenly(color RGB) []float32 {
	return []float32{
		color.red, color.green, color.blue, 1,
		color.red, color.green, color.blue, 1,
		color.red, color.green, color.blue, 1,
	}
}

func (OGE *openglGraphicsEngine) SixP2Triangle(x1, y1, x2, y2, x3, y3 float32) []float32 {
	return []float32{
		x1, y1,
		x2, y2,
		x3, y3,
	}
}

func (OGE *openglGraphicsEngine) Start() {
	for _, ele := range OGE.elements {
		ele.onCreate()
	}

	for !OGE.window.ShouldClose() {
		frameStartTime := time.Now()

		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		gl.UseProgram(OGE.program)

		for _, ele := range OGE.elements {
			ele.onUpdate()
		}

		glfw.PollEvents()
		OGE.window.SwapBuffers()

		// Control framerate
		frameTime := 1.0 / float64(OGE.fps)
		for time.Since(frameStartTime).Seconds() < frameTime {
			time.Sleep(time.Duration(frameTime/4) * time.Second)
		}
		OGE.window.SetTitle(OGE.title + fmt.Sprintf(" FPS: %.1f", 1.0/time.Since(frameStartTime).Seconds()))
		OGE.delta = time.Since(frameStartTime).Seconds() * OGE.fps
	}
}

func (OGE *openglGraphicsEngine) addElement(new element) {
	OGE.elements = append(OGE.elements, new)
}

// --------------------- test code -----------------------------
// func main() {

// 	// OGE one frame test
// 	var (
// 		width        = 500
// 		height       = 500
// 		triangleTest = []float32{
// 			0, 0.5, // top
// 			-0.5, -0.5, // left
// 			0.5, -0.5, // right
// 		}
// 		triangleTest2 = []float32{
// 			0.0, -0.5, // bottom
// 			0.5, 0.5, // right
// 			-0.5, 0.5, // letf
// 		}
// 	)
// 	oge := constructOpenglGraphicsEngine(width, height, "OGE TEST", 60)

// 	for !oge.window.ShouldClose() {

// 		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
// 		gl.UseProgram(oge.program)

// 		oge.FillTriangle(triangleTest2, oge.TriangleColorEvenly(WHITE))
// 		oge.DrawTriangle(triangleTest, oge.TriangleColorEvenly(oge.ScaleRGB(BLUE, 0.5)))

// 		glfw.PollEvents()
// 		oge.window.SwapBuffers()
// 	}

// 	oge.destructor()
// }

// ------------------------  Helper funcsions -----------------------------------

func initGlfw(screenWidth int, screenHeight int, name string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(screenWidth, screenHeight, name, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	point_ptr := []uint8("vertice_point")
	color_ptr := []uint8("vertice_color")
	gl.BindAttribLocation(prog, 0, &point_ptr[0])
	gl.BindAttribLocation(prog, 1, &color_ptr[0])
	gl.LinkProgram(prog)
	gl.ValidateProgram(prog)
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)
	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}
	return shader, nil
}

func makeVao(vertices []float32, colors []float32) uint32 {
	var vao uint32    // Vertex Array Object
	var vbo [2]uint32 // Vertex Buffer Object
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.GenBuffers(2, &vbo[0])

	// Vertices
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[0])
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(vertices), gl.Ptr(vertices), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)

	// Color
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[1])
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(colors), gl.Ptr(colors), gl.STATIC_DRAW)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, nil)

	return vao
}
