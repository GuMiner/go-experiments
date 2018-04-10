package genetics

import (
	"go-experiments/voxelli/neural"
	"go-experiments/voxelli/renderer"
	"go-experiments/voxelli/roadway"
	"go-experiments/voxelli/vehicle"
	"go-experiments/voxelli/voxelArray"

	"github.com/go-gl/mathgl/mgl32"
)

// If a car is this close to a wall and hits it, this lets it scoot along the wall slowly.
const wiggleDistance = 4.0
const wiggleSpeed = 0.03

type Agent struct {
	startingPosition mgl32.Vec2

	car     *vehicle.Vehicle
	net     *neural.NeuralNet
	isAlive bool
}

func (a *Agent) reset() {
	a.isAlive = true

	a.car.Score = 0
	a.car.Position = a.startingPosition
}

// TODO: Refactor so we don't need this for debug drawing car info.
func (a *Agent) GetCar() *vehicle.Vehicle {
	return a.car
}

func (a *Agent) GetFinalScore() float32 {
	return a.car.Score
}

func getSmallestBoundary(boundaryLengths []float32, boundaryNormals []mgl32.Vec2) (float32, mgl32.Vec2) {
	if boundaryLengths[0] < boundaryLengths[2] {
		if boundaryLengths[0] < boundaryLengths[1] {
			return boundaryLengths[0], boundaryNormals[0]
		} else {
			return boundaryLengths[1], boundaryNormals[1]
		}
	} else {
		if boundaryLengths[2] < boundaryLengths[1] {
			return boundaryLengths[2], boundaryNormals[2]
		} else {
			return boundaryLengths[1], boundaryNormals[1]
		}
	}
}

// Updates the agent, returning true if the agent is alive, false otherwise
func (a *Agent) Update(frameTime float32, roadway *roadway.Roadway) {
	if a.isAlive {
		hitWall := a.car.Update(frameTime, roadway)
		eyePositions, eyeDirections := a.car.GetEyes()
		boundaryLengths, boundaryNormals := roadway.GetBoundaries(eyePositions, eyeDirections)
		if hitWall {
			// Bounce along the direction with the shortest normal, to let cars that just miss turns (and which are moving straight) keep going.
			boundaryLength, boundaryNormal := getSmallestBoundary(boundaryLengths, boundaryNormals)
			if boundaryLength < wiggleDistance {
				a.car.Position = a.car.Position.Add(boundaryNormal.Normalize().Mul(wiggleSpeed))
			} else {
				// We can't wiggle, so we are dead.
				a.isAlive = false
			}
		}

		steeringAndAccel := a.net.Evaluate(append(boundaryLengths, a.car.Velocity))
		a.car.SteeringPos = steeringAndAccel[0]*2 - 1
		a.car.AccelPos = steeringAndAccel[1]*2 - 1
	}
}

func (a *Agent) Render(renderer *renderer.VoxelArrayObjectRenderer) {
	if a.isAlive {
		a.car.Render(renderer)
	}
}

// Modifies this agent by crossbreeding it with the two given agents.
func (a *Agent) CrossBreed(first, second *Agent, crossoverProbability float32) {
	a.net.CrossMerge(first.net, second.net, crossoverProbability)

	a.reset()
}

func NewAgent(id int, carModel *voxelArray.VoxelArrayObject, startingPosition mgl32.Vec2) *Agent {
	agent := Agent{
		car:              vehicle.NewVehicle(id, carModel),
		net:              neural.NewNeuralNet([]int{4, 7, 7, 7, 7}, 2),
		startingPosition: startingPosition}
	agent.reset()

	return &agent
}
