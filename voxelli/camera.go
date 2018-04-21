package main

import (
	"math"

	"go-experiments/voxelli/cache"
	"go-experiments/voxelli/config"
	"go-experiments/voxelli/input"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	Position mgl32.Vec3

	// Normalized
	Forwards mgl32.Vec3
	Up       mgl32.Vec3

	// Computed & normalized
	Right  mgl32.Vec3
	Target mgl32.Vec3

	wasMouseDown bool
	lastMousePos mgl32.Vec2
}

const cacheName = "camera"

func (c *Camera) normalize() {
	c.Up = c.Up.Normalize()
	c.Forwards = c.Forwards.Normalize()

	c.Target = c.Position.Add(c.Forwards)
	c.Right = c.Up.Cross(c.Forwards)

	// fmt.Printf("Updated camera: %+v\n", c)
}

func (c *Camera) handleLinearMotion(positive input.KeyAssignment, negative input.KeyAssignment, direction mgl32.Vec3, scale float32) bool {
	Updated := false
	if input.IsPressed(positive) {
		c.Position = c.Position.Add(direction.Mul(scale))
		Updated = true
	}

	if input.IsPressed(negative) {
		c.Position = c.Position.Sub(direction.Mul(scale))
		Updated = true
	}

	return Updated
}

func (c *Camera) Update(frameTime float32) {
	moveSpeed := frameTime * config.Config.Camera.MotionSpeed
	rotateSpeed := frameTime * config.Config.Camera.RotationSpeed

	Updated := c.handleLinearMotion(input.MoveForwards, input.MoveBackwards, c.Forwards, moveSpeed)
	Updated = c.handleLinearMotion(input.MoveLeft, input.MoveRight, c.Right, moveSpeed) || Updated
	Updated = c.handleLinearMotion(input.MoveUp, input.MoveDown, c.Up, moveSpeed) || Updated

	// Key-based rotation
	if input.IsPressed(input.RotateClockwise) {
		c.Up = mgl32.HomogRotate3D(rotateSpeed, c.Forwards).Mul4x1(c.Up.Vec4(1.0)).Vec3()
		Updated = true
	}

	if input.IsPressed(input.RotateCounterClockwise) {
		c.Up = mgl32.HomogRotate3D(-rotateSpeed, c.Forwards).Mul4x1(c.Up.Vec4(1.0)).Vec3()
		Updated = true
	}

	if input.IsPressed(input.LookLeft) {
		c.Forwards = mgl32.HomogRotate3D(rotateSpeed, c.Up).Mul4x1(c.Forwards.Vec4(1.0)).Vec3()
		Updated = true
	}

	if input.IsPressed(input.LookRight) {
		c.Forwards = mgl32.HomogRotate3D(-rotateSpeed, c.Up).Mul4x1(c.Forwards.Vec4(1.0)).Vec3()
		Updated = true
	}

	if input.IsPressed(input.LookUp) {
		c.Forwards = mgl32.HomogRotate3D(rotateSpeed, c.Right).Mul4x1(c.Forwards.Vec4(1.0)).Vec3()
		c.Up = mgl32.HomogRotate3D(rotateSpeed, c.Right).Mul4x1(c.Up.Vec4(1.0)).Vec3()
		Updated = true
	}

	if input.IsPressed(input.LookDown) {
		c.Forwards = mgl32.HomogRotate3D(-rotateSpeed, c.Right).Mul4x1(c.Forwards.Vec4(1.0)).Vec3()
		c.Up = mgl32.HomogRotate3D(-rotateSpeed, c.Right).Mul4x1(c.Up.Vec4(1.0)).Vec3()
		Updated = true
	}

	if Updated {
		c.normalize()
	}

	// Mouse-based rotation
	if input.PressedButtons[glfw.MouseButtonRight] {
		if !c.wasMouseDown {
			c.wasMouseDown = true
			c.lastMousePos = input.MousePos
		}

		difference := input.MousePos.Sub(c.lastMousePos)
		if math.Abs(float64(difference.X())) > 1 {
			mouseRotationSpeed := config.Config.Camera.RotationSpeed * difference.X() / 1200.0

			c.Forwards = mgl32.HomogRotate3D(mouseRotationSpeed, c.Up).Mul4x1(c.Forwards.Vec4(1.0)).Vec3()
			Updated = true
		}

		if math.Abs(float64(difference.Y())) > 1 {
			mouseRotationSpeed := -config.Config.Camera.RotationSpeed * difference.Y() / 1200.0

			c.Forwards = mgl32.HomogRotate3D(mouseRotationSpeed, c.Right).Mul4x1(c.Forwards.Vec4(1.0)).Vec3()
			c.Up = mgl32.HomogRotate3D(mouseRotationSpeed, c.Right).Mul4x1(c.Up.Vec4(1.0)).Vec3()
			Updated = true
		}

		c.lastMousePos = input.MousePos
	} else {
		c.wasMouseDown = false
	}

	if Updated {
		c.normalize()
	}
}

func (c *Camera) GetLookAtMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Target, c.Up)
}

func (c *Camera) CachePosition() {
	cache.SaveToCache(cacheName, c)
}

func NewCamera(Position mgl32.Vec3, Forwards mgl32.Vec3, Up mgl32.Vec3) *Camera {
	camera := Camera{wasMouseDown: false, lastMousePos: mgl32.Vec2{-1.0, -1.0}}

	cacheMiss := cache.LoadFromCache(cacheName, true, &camera)
	if cacheMiss {
		// TODO: This original logic from Fractal / etc doesn't properly Update
		//  Forwards / Up / Right to be at proper Right angles to each other
		camera.Position = Position
		camera.Forwards = Forwards
		camera.Up = Up
		camera.normalize()
	}

	return &camera
}
