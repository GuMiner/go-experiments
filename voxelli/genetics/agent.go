package genetics

import (
	"go-experiments/voxelli/neural"
	"go-experiments/voxelli/renderer"
	"go-experiments/voxelli/roadway"
	"go-experiments/voxelli/vehicle"
	"go-experiments/voxelli/voxelArray"

	"github.com/go-gl/mathgl/mgl32"
)

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

// Updates the agent, returning true if the agent is alive, false otherwise
func (a *Agent) Update(frameTime float32, roadway *roadway.Roadway) {
	hitWall := a.car.Update(frameTime, roadway)
	if hitWall {
		a.isAlive = false
	} else {
		eyePositions, eyeDirections := a.car.GetEyes()
		boundaryLengths := roadway.GetBoundaries(eyePositions, eyeDirections)

		steeringAndAccel := a.net.Evaluate(append(boundaryLengths, a.car.Velocity))
		a.car.SteeringPos = steeringAndAccel[0]*2 - 1
		a.car.AccelPos = steeringAndAccel[1]*2 - 1
	}
}

func (a *Agent) Render(renderer *renderer.VoxelArrayObjectRenderer) {
	a.car.Render(renderer)
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
