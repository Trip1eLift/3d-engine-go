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
	graphics *openglGraphicsEngine
	color    RGB
	meshCube mesh
	matProj  mat4x4
	fTheta   float32
	vCamera  vec3d
}

func newCube(container *openglGraphicsEngine) *cube {
	var c cube
	c.graphics = container
	return &c
}

func (c *cube) onCreate() bool {
	c.color = RED
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

func (c *cube) projectAndDrawTriangle(tri triangle, matRotZ mat4x4, matRotX mat4x4) {
	var triProjected, triTranslated, triRotatedZ, triRotatedZX triangle

	// Rotate in Z-Axis
	triRotatedZ.p[0] = MultiplyMatrixVector(tri.p[0], matRotZ)
	triRotatedZ.p[1] = MultiplyMatrixVector(tri.p[1], matRotZ)
	triRotatedZ.p[2] = MultiplyMatrixVector(tri.p[2], matRotZ)

	// Rotate in X-Axis
	triRotatedZX.p[0] = MultiplyMatrixVector(triRotatedZ.p[0], matRotX)
	triRotatedZX.p[1] = MultiplyMatrixVector(triRotatedZ.p[1], matRotX)
	triRotatedZX.p[2] = MultiplyMatrixVector(triRotatedZ.p[2], matRotX)

	// Offset into the screen
	triTranslated = triRotatedZX
	triTranslated.p[0].z = triRotatedZX.p[0].z + 3
	triTranslated.p[1].z = triRotatedZX.p[1].z + 3
	triTranslated.p[2].z = triRotatedZX.p[2].z + 3

	// Hide lines behind object
	var normal, line1, line2 vec3d
	line1.x = triTranslated.p[1].x - triTranslated.p[0].x
	line1.y = triTranslated.p[1].y - triTranslated.p[0].y
	line1.z = triTranslated.p[1].z - triTranslated.p[0].z

	line2.x = triTranslated.p[2].x - triTranslated.p[0].x
	line2.y = triTranslated.p[2].y - triTranslated.p[0].y
	line2.z = triTranslated.p[2].z - triTranslated.p[0].z

	normal.x = line1.y*line2.z - line1.z*line2.y
	normal.y = line1.z*line2.x - line1.x*line2.z
	normal.z = line1.x*line2.y - line1.y*line2.x

	var l float32 = float32(math.Sqrt(float64(normal.x*normal.x + normal.y*normal.y + normal.z*normal.z)))
	normal.x /= l
	normal.y /= l
	normal.z /= l

	if (normal.x*(triTranslated.p[0].x-c.vCamera.x) +
		normal.y*(triTranslated.p[0].y-c.vCamera.y) +
		normal.z*(triTranslated.p[0].z-c.vCamera.z)) < 0 {
		// Project triangles from 3D -> 2D
		triProjected.p[0] = MultiplyMatrixVector(triTranslated.p[0], c.matProj)
		triProjected.p[1] = MultiplyMatrixVector(triTranslated.p[1], c.matProj)
		triProjected.p[2] = MultiplyMatrixVector(triTranslated.p[2], c.matProj)

		tri := c.graphics.SixP2Triangle(triProjected.p[0].x, triProjected.p[0].y, triProjected.p[1].x, triProjected.p[1].y, triProjected.p[2].x, triProjected.p[2].y)

		c.graphics.DrawTriangle(tri, c.graphics.TriangleColorEvenly(WHITE))
	}

}

func (c *cube) onUpdate() bool {
	var matRotZ, matRotX mat4x4
	c.fTheta += 0.01 * float32(c.graphics.delta)
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

	// Draw Triangles
	for _, tri := range c.meshCube.tris {
		c.projectAndDrawTriangle(tri, matRotZ, matRotX)
	}

	return true
}

func main() {
	engine := constructOpenglGraphicsEngine(500, 500, "Cube spin", 75)
	cube := newCube(engine)
	engine.addElement(cube)
	engine.Start()

	return
}
