package main

import (
	"go-experiments/voxelli/input"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	position mgl32.Vec3

	// Normalized
	forwards mgl32.Vec3
	up       mgl32.Vec3

	// Computed & normalized
	right  mgl32.Vec3
	target mgl32.Vec3
}

const motionSpeed = 20.0
const rotationSpeed = 1.0

func (c *Camera) normalize() {
	c.up = c.up.Normalize()
	c.forwards = c.forwards.Normalize()

	c.target = c.position.Add(c.forwards)
	c.right = c.up.Cross(c.forwards)
}

func (c *Camera) handleLinearMotion(positive glfw.Key, negative glfw.Key, direction mgl32.Vec3, scale float32) bool {
	updated := false
	if input.PressedKeys[positive] {
		c.position = c.position.Add(direction.Mul(scale))
		updated = true
	}

	if input.PressedKeys[negative] {
		c.position = c.position.Sub(direction.Mul(scale))
		updated = true
	}

	return updated
}

func (c *Camera) Update(frameTime float32, cameraMatrix *mgl32.Mat4) bool {
	moveSpeed := frameTime * motionSpeed
	rotateSpeed := frameTime * rotationSpeed

	updated := false
	updated = c.handleLinearMotion(glfw.KeyA, glfw.KeyZ, c.forwards, moveSpeed) || updated
	updated = c.handleLinearMotion(glfw.KeyQ, glfw.KeyW, c.right, moveSpeed) || updated
	updated = c.handleLinearMotion(glfw.KeyS, glfw.KeyX, c.up, moveSpeed) || updated

	if input.PressedKeys[glfw.KeyE] {
		c.up = mgl32.HomogRotate3D(rotateSpeed, c.forwards).Mul4x1(c.up.Vec4(1.0)).Vec3()
	}

	if updated {
		c.normalize()
		*cameraMatrix = c.GetLookAtMatrix()
	}

	return updated
}

func (c *Camera) GetLookAtMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.position, c.target, c.up)
}

func NewCamera(position mgl32.Vec3, forwards mgl32.Vec3, up mgl32.Vec3) *Camera {
	var camera Camera

	// TODO: This original logic from Fractal / etc doesn't properly update
	//  forwards / up / right to be at proper right angles to each other
	camera.position = position
	camera.forwards = forwards
	camera.up = up
	camera.normalize()

	return &camera
}
