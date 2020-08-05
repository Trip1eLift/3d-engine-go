package main

import (
	"fmt"
	"reflect"
)

// --------------- HELPER struct and functions ----------------------
type vec3d struct {
	x, y, z float32
}

type triangle struct {
	p0 vec3d
	p1 vec3d
	p2 vec3d
}

type mesh {
	tris []triangle
}

type mat4x4 {
	m [4][4]float32
}

func MultiplyMatrixVector(in vec3d, m mat4x4) vec3d {
	var out vec3d
	out.x = in.x*m.m[0][0] + in.y*m.m[1][0] + in.z*m.m[2][0] + m.m[3][0]
	out.y = in.x*m.m[0][1] + in.y*m.m[1][1] + in.z*m.m[2][1] + m.m[3][1]
	out.z = in.x*m.m[0][2] + in.y*m.m[1][2] + in.z*m.m[2][2] + m.m[3][2]
	var w float32 = in.x*m.m[0][3] + in.y*m.m[1][3] + in.z*m.m[2][3] + m.m[3][3]

	if w != 0.0 {
		out.x /= w
		out.y /= w
		out.z /= w
	}
	return out
}

// ---------------------------- 3D cube --------------------------

type cube struct {
	graphics *consoleGraphicEngine
	color    string
	meshCube mesh
	matProj  mat4x4
	fTheta float32
}

func newCube(container *consoleGraphicEngine) *cube {
	var c cube
	c.graphics = container
	return &c
}

func (c *cube) onCreate() {
	c.color = RED
	// c.graphics.drawTriangle(3, 3, 250, 3, 3, 99, FULL_BLOCK, c.color)
	// c.meshCube.tris = {
	// 	// SOUTH
	// 	triangle{ vec3d{0.0, 0.0, 0.0},    vec3d{0.0, 1.0, 0.0},    vec3d{1.0, 1.0, 0.0} },
	// 	triangle{ vec3d{0.0, 0.0, 0.0},    vec3d{1.0, 1.0, 0.0},    vec3d{1.0, 0.0, 0.0} },

	// 	// EAST
	// 	triangle{ vec3d{1.0, 0.0, 0.0},    vec3d{1.0, 1.0, 0.0},    vec3d{1.0, 1.0, 1.0} },
	// 	triangle{ vec3d{1.0, 0.0, 0.0},    vec3d{1.0, 1.0, 1.0},    vec3d{1.0, 0.0, 1.0} },

	// 	// NORTH
	// 	triangle{ vec3d{1.0, 0.0, 1.0},    vec3d{1.0, 1.0, 1.0},    vec3d{0.0, 1.0, 1.0} },
	// 	triangle{ vec3d{1.0, 0.0, 1.0},    vec3d{0.0, 1.0, 1.0},    vec3d{0.0, 0.0, 1.0} },

	// 	// WEST
	// 	triangle{ vec3d{0.0, 0.0, 1.0},    vec3d{0.0, 1.0, 1.0},    vec3d{0.0, 1.0, 0.0} },
	// 	triangle{ vec3d{0.0, 0.0, 1.0},    vec3d{0.0, 1.0, 0.0},    vec3d{0.0, 0.0, 0.0} },

	// 	// TOP
	// 	triangle{ vec3d{0.0, 1.0, 0.0},    vec3d{0.0, 1.0, 1.0},    vec3d{1.0, 1.0, 1.0} },
	// 	triangle{ vec3d{0.0, 1.0, 0.0},    vec3d{1.0, 1.0, 1.0},    vec3d{1.0, 1.0, 0.0} },

	// 	// BOTTOM
	// 	triangle{ vec3d{1.0, 0.0, 1.0},    vec3d{0.0, 0.0, 1.0},    vec3d{0.0, 0.0, 0.0} },
	// 	triangle{ vec3d{1.0, 0.0, 1.0},    vec3d{0.0, 0.0, 0.0},    vec3d{1.0, 0.0, 0.0} },
	// }
	// fNear := 0.1
	// fFar := 1000.0
	// fFov := 90.0
	// fAspectRatio := (float32)c.graphics.screenHeight / (float32)c.graphics.screenWidth
	// fFovRad := 1.0 / math.Tan(fFov * 0.5 / 180.0 * 3.14159)

	// c.matProj.m[0][0] = fAspectRatio * fFovRad
	// c.matProj.m[1][1] = fFovRad
	// c.matProj.m[2][2] = fFar / (fFar - fNear)
	// c.matProj.m[3][2] = (-fFar * fNear) / (fFar - fNear)
	// c.matProj.m[2][3] = 1.0
	// c.matProj.m[3][3] = 0.0
	return
}

func (c *cube) onUpdate() {
	if c.color == RED {
		c.color = BLUE
	} else {
		c.color = RED
	}
	c.graphics.drawTriangle(3, 3, 250, 3, 3, 99, FULL_BLOCK, c.color)
	return
}

func main() {
	// engine := constructConsoleGraphicEngine(300, 100, WHITE)
	// cube := newCube(engine)
	// engine.addComponent(cube)
	// engine.Start()

	// vec0 := vec3d{0.0, 0.0, 0.0}
	// vec1 := vec3d{0.0, 1.0, 0.0}
	// vec2 := vec3d{1.0, 1.0, 0.0}
	// test := triangle{vec0, vec1, vec2}
	test1 := triangle{vec3d{0.0, 0.0, 0.0}, vec3d{0.0, 1.0, 0.0}, vec3d{1.0, 1.0, 0.0}}
	test2 := triangle{vec3d{0.0, 0.0, 0.0}, vec3d{0.0, 1.0, 0.0}, vec3d{1.0, 1.0, 0.0}}
	test3 := triangle{vec3d{0.0, 0.0, 0.0}, vec3d{0.0, 1.0, 0.0}, vec3d{1.0, 1.0, 0.0}}
	test4 := mesh{test1, test2, test3}
	fmt.Println(reflect.TypeOf(test1.p1.y))
	fmt.Println(mesh.tris[0].test1.p1.y)

	return

	// Zoom terminal out to smallest by [ctrl] + [-]
}
