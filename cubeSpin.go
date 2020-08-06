package main

import "math"

// --------------- HELPER struct and functions ----------------------
type vec3d struct {
	x, y, z float32
}

type triangle struct {
	p [3]vec3d
}

type mesh struct {
	tris []triangle
}

type mat4x4 struct {
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
	fTheta   float32
}

func newCube(container *consoleGraphicEngine) *cube {
	var c cube
	c.graphics = container
	return &c
}

func (c *cube) onCreate() bool {
	c.color = RED
	// c.graphics.drawTriangle(3, 3, 250, 3, 3, 99, FULL_BLOCK, c.color)
	var tri triangle

	// SOUTH
	tri.p[0], tri.p[1], tri.p[2] = vec3d{0.0, 0.0, 0.0}, vec3d{0.0, 1.0, 0.0}, vec3d{1.0, 1.0, 0.0}
	c.meshCube.tris = append(c.meshCube.tris, tri)
	tri.p[0], tri.p[1], tri.p[2] = vec3d{0.0, 0.0, 0.0}, vec3d{1.0, 1.0, 0.0}, vec3d{1.0, 0.0, 0.0}
	c.meshCube.tris = append(c.meshCube.tris, tri)

	// EAST
	tri.p[0], tri.p[1], tri.p[2] = vec3d{1.0, 0.0, 0.0}, vec3d{1.0, 1.0, 0.0}, vec3d{1.0, 1.0, 1.0}
	c.meshCube.tris = append(c.meshCube.tris, tri)
	tri.p[0], tri.p[1], tri.p[2] = vec3d{1.0, 0.0, 0.0}, vec3d{1.0, 1.0, 1.0}, vec3d{1.0, 0.0, 1.0}
	c.meshCube.tris = append(c.meshCube.tris, tri)

	// NORTH
	tri.p[0], tri.p[1], tri.p[2] = vec3d{1.0, 0.0, 1.0}, vec3d{1.0, 1.0, 1.0}, vec3d{0.0, 1.0, 1.0}
	c.meshCube.tris = append(c.meshCube.tris, tri)
	tri.p[0], tri.p[1], tri.p[2] = vec3d{1.0, 0.0, 1.0}, vec3d{0.0, 1.0, 1.0}, vec3d{0.0, 0.0, 1.0}
	c.meshCube.tris = append(c.meshCube.tris, tri)

	// WEST
	tri.p[0], tri.p[1], tri.p[2] = vec3d{0.0, 0.0, 1.0}, vec3d{0.0, 1.0, 1.0}, vec3d{0.0, 1.0, 0.0}
	c.meshCube.tris = append(c.meshCube.tris, tri)
	tri.p[0], tri.p[1], tri.p[2] = vec3d{0.0, 0.0, 1.0}, vec3d{0.0, 1.0, 0.0}, vec3d{0.0, 0.0, 0.0}
	c.meshCube.tris = append(c.meshCube.tris, tri)

	// TOP
	tri.p[0], tri.p[1], tri.p[2] = vec3d{0.0, 1.0, 0.0}, vec3d{0.0, 1.0, 1.0}, vec3d{1.0, 1.0, 1.0}
	c.meshCube.tris = append(c.meshCube.tris, tri)
	tri.p[0], tri.p[1], tri.p[2] = vec3d{0.0, 1.0, 0.0}, vec3d{1.0, 1.0, 1.0}, vec3d{1.0, 1.0, 0.0}
	c.meshCube.tris = append(c.meshCube.tris, tri)

	// BOTTOM
	tri.p[0], tri.p[1], tri.p[2] = vec3d{1.0, 0.0, 1.0}, vec3d{0.0, 0.0, 1.0}, vec3d{0.0, 0.0, 0.0}
	c.meshCube.tris = append(c.meshCube.tris, tri)
	tri.p[0], tri.p[1], tri.p[2] = vec3d{1.0, 0.0, 1.0}, vec3d{0.0, 0.0, 0.0}, vec3d{1.0, 0.0, 0.0}
	c.meshCube.tris = append(c.meshCube.tris, tri)

	fNear := float32(0.1)
	fFar := float32(1000.0)
	fFov := 90.0
	fAspectRatio := float32(c.graphics.screenHeight) / float32(c.graphics.screenWidth)
	fFovRad := float32(1.0) / float32(math.Tan(fFov*0.5/180.0*3.14159))

	c.matProj.m[0][0] = fAspectRatio * fFovRad
	c.matProj.m[1][1] = fFovRad
	c.matProj.m[2][2] = fFar / (fFar - fNear)
	c.matProj.m[3][2] = (-fFar * fNear) / (fFar - fNear)
	c.matProj.m[2][3] = 1.0
	c.matProj.m[3][3] = 0.0
	return true
}

func (c *cube) projectAndDrawTriangle(tri triangle, matRotZ mat4x4, matRotX mat4x4, matScaleX mat4x4) {
	var triProjected, triTranslated, triRotatedZ, triRotatedZX triangle

	// Rotate in Z-Axis
	triRotatedZ.p[0] = MultiplyMatrixVector(tri.p[0], matRotZ)
	triRotatedZ.p[1] = MultiplyMatrixVector(tri.p[1], matRotZ)
	triRotatedZ.p[2] = MultiplyMatrixVector(tri.p[2], matRotZ)

	// Rotate in X-Axis
	triRotatedZX.p[0] = MultiplyMatrixVector(triRotatedZ.p[0], matRotX)
	triRotatedZX.p[1] = MultiplyMatrixVector(triRotatedZ.p[1], matRotX)
	triRotatedZX.p[2] = MultiplyMatrixVector(triRotatedZ.p[2], matRotX)

	triTranslated = triRotatedZX
	triTranslated.p[0].z = triRotatedZX.p[0].z + 2
	triTranslated.p[1].z = triRotatedZX.p[1].z + 2
	triTranslated.p[2].z = triRotatedZX.p[2].z + 2

	triProjected.p[0] = MultiplyMatrixVector(triTranslated.p[0], c.matProj)
	triProjected.p[1] = MultiplyMatrixVector(triTranslated.p[1], c.matProj)
	triProjected.p[2] = MultiplyMatrixVector(triTranslated.p[2], c.matProj)

	// Scale X
	triProjected.p[0] = MultiplyMatrixVector(triProjected.p[0], matScaleX)
	triProjected.p[1] = MultiplyMatrixVector(triProjected.p[1], matScaleX)
	triProjected.p[2] = MultiplyMatrixVector(triProjected.p[2], matScaleX)

	// Scale into view
	triProjected.p[0].x += 1.0
	triProjected.p[0].y += 1.0
	triProjected.p[1].x += 1.0
	triProjected.p[1].y += 1.0
	triProjected.p[2].x += 1.0
	triProjected.p[2].y += 1.0

	triProjected.p[0].x *= 0.5 * float32(c.graphics.screenWidth)
	triProjected.p[0].y *= 0.5 * float32(c.graphics.screenHeight)
	triProjected.p[1].x *= 0.5 * float32(c.graphics.screenWidth)
	triProjected.p[1].y *= 0.5 * float32(c.graphics.screenHeight)
	triProjected.p[2].x *= 0.5 * float32(c.graphics.screenWidth)
	triProjected.p[2].y *= 0.5 * float32(c.graphics.screenHeight)

	c.graphics.drawTriangle(
		int(triProjected.p[0].x), int(triProjected.p[0].y),
		int(triProjected.p[1].x), int(triProjected.p[1].y),
		int(triProjected.p[2].x), int(triProjected.p[2].y),
		FULL_BLOCK, WHITE)
}

func (c *cube) onUpdate() bool {
	c.graphics.fillALL(SPACE_BLOCK, WHITE)
	var matRotZ, matRotX, matScaleX mat4x4
	c.fTheta += 0.2
	var fTheta float64 = float64(c.fTheta)

	// Rotation Z
	matRotZ.m[0][0] = float32(math.Cos(fTheta))
	matRotZ.m[0][1] = float32(math.Sin(fTheta))
	matRotZ.m[1][0] = float32(-math.Sin(fTheta))
	matRotZ.m[1][1] = float32(math.Cos(fTheta))
	matRotZ.m[2][2] = 1
	matRotZ.m[3][3] = 1

	// Rotation X
	matRotX.m[0][0] = 1
	matRotX.m[1][1] = float32(math.Cos(fTheta))
	matRotX.m[1][2] = float32(math.Sin(fTheta))
	matRotX.m[2][1] = float32(-math.Sin(fTheta))
	matRotX.m[2][2] = float32(math.Cos(fTheta))
	matRotX.m[3][3] = 1

	// Scale X to stretch out
	matScaleX.m[0][0] = 2.5
	matScaleX.m[1][1] = 1
	matScaleX.m[2][2] = 1

	// Draw Triangles
	for _, tri := range c.meshCube.tris {
		go c.projectAndDrawTriangle(tri, matRotZ, matRotX, matScaleX)
	}

	return true
}

func main() {
	engine := constructConsoleGraphicEngine(300, 100, WHITE)
	cube := newCube(engine)
	engine.addComponent(cube)
	engine.Start()

	return

	// Zoom terminal out to smallest by [ctrl] + [-]
}
