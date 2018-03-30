package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"

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
}

const motionSpeed = 20.0
const rotationSpeed = 1.0
const cacheFilename = "./data/cache/camera.gob"

func (c *Camera) normalize() {
	c.Up = c.Up.Normalize()
	c.Forwards = c.Forwards.Normalize()

	c.Target = c.Position.Add(c.Forwards)
	c.Right = c.Up.Cross(c.Forwards)
}

func (c *Camera) handleLinearMotion(positive glfw.Key, negative glfw.Key, direction mgl32.Vec3, scale float32) bool {
	Updated := false
	if input.PressedKeys[positive] {
		c.Position = c.Position.Add(direction.Mul(scale))
		Updated = true
	}

	if input.PressedKeys[negative] {
		c.Position = c.Position.Sub(direction.Mul(scale))
		Updated = true
	}

	return Updated
}

func (c *Camera) Update(frameTime float32, cameraMatrix *mgl32.Mat4) bool {
	moveSpeed := frameTime * motionSpeed
	rotateSpeed := frameTime * rotationSpeed

	Updated := false
	Updated = c.handleLinearMotion(glfw.KeyA, glfw.KeyZ, c.Forwards, moveSpeed) || Updated
	Updated = c.handleLinearMotion(glfw.KeyQ, glfw.KeyW, c.Right, moveSpeed) || Updated
	Updated = c.handleLinearMotion(glfw.KeyS, glfw.KeyX, c.Up, moveSpeed) || Updated

	if input.PressedKeys[glfw.KeyE] {
		c.Up = mgl32.HomogRotate3D(rotateSpeed, c.Forwards).Mul4x1(c.Up.Vec4(1.0)).Vec3()
	}

	if Updated {
		c.normalize()
		*cameraMatrix = c.GetLookAtMatrix()
	}

	return Updated
}

func (c *Camera) GetLookAtMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(c.Position, c.Target, c.Up)
}

func (c *Camera) CachePosition() {
	byteBuffer := new(bytes.Buffer)
	err := gob.NewEncoder(byteBuffer).Encode(c)
	if err != nil {
		fmt.Printf("Unable to encode the camera Position: %v\n", err)
	} else {
		err = ioutil.WriteFile(cacheFilename, byteBuffer.Bytes(), os.ModePerm)
		if err != nil {
			fmt.Printf("Unable to cache the encoded camera Position: %v\n", err)
		}
	}
}

func NewCamera(Position mgl32.Vec3, Forwards mgl32.Vec3, Up mgl32.Vec3) *Camera {
	var camera Camera

	cacheFileAsBytes, err := ioutil.ReadFile(cacheFilename)
	cacheMiss := true
	if err == nil {
		// Decode from cache, ignore failures
		err = gob.NewDecoder(bytes.NewBuffer(cacheFileAsBytes)).Decode(&camera)
		if err == nil {
			cacheMiss = false
			fmt.Printf("Loaded cached camera at: %+v\n", camera)
		}
	}

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
