package main

import (
	"fmt"
	"go-experiments/voxelli/input"
	"go-experiments/voxelli/renderer"
	"go-experiments/voxelli/roadway"
	"go-experiments/voxelli/vehicle"

	"github.com/go-gl/glfw/v3.2/glfw"
)

type ControllableCar struct {
	Car *vehicle.Vehicle
}

const accelScaleFactor = 0.1
const steerScaleFactor = 0.1

func (c *ControllableCar) Update(frameTime float32, roadway *roadway.Roadway) {

	pressedKey := false
	if input.PressedKeys[glfw.KeyI] {
		c.Car.AccelPos += accelScaleFactor * frameTime
		pressedKey = true
	}

	if input.PressedKeys[glfw.KeyK] {
		c.Car.AccelPos -= accelScaleFactor * frameTime
		pressedKey = true
	}

	if input.PressedKeys[glfw.KeyJ] {
		c.Car.SteeringPos += steerScaleFactor * frameTime
		pressedKey = true
	}

	if input.PressedKeys[glfw.KeyL] {
		c.Car.SteeringPos -= steerScaleFactor * frameTime
		pressedKey = true
	}

	if pressedKey {
		fmt.Printf("A: %.2f S: %.2f\n", c.Car.AccelPos, c.Car.SteeringPos)
	}

	c.Car.Update(frameTime, roadway)
}

func (c *ControllableCar) Render(voxelArrayObjectRenderer *renderer.VoxelArrayObjectRenderer, roadway *roadway.Roadway) {
	c.Car.Render(voxelArrayObjectRenderer)
}
