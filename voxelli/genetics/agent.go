package genetics

import (
	"go-experiments/voxelli/neural"
	"go-experiments/voxelli/renderer"
	"go-experiments/voxelli/roadway"
	"go-experiments/voxelli/vehicle"
	"go-experiments/voxelli/voxelArray"
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

type Agent struct {
	car        *vehicle.Vehicle
	net        *neural.NeuralNet
	FinalScore float32
	isAlive    bool
}

// TODO: Refactor so we don't need this for debug drawing car info.
func (a *Agent) GetCar() *vehicle.Vehicle {
	return a.car
}

// Updates the agent, returning true if the agent is alive, false otherwise
func (a *Agent) Update(frameTime float32, roadway *roadway.Roadway) bool {
	if !a.isAlive {
		return false
	}

	hitWall := a.car.Update(frameTime, roadway)
	if hitWall {
		a.FinalScore = a.car.Score
		a.isAlive = false
	} else {
		eyePositions, eyeDirections := a.car.GetEyes()
		boundaryLengths := roadway.GetBoundaries(eyePositions, eyeDirections)

		steeringAndAccel := a.net.Evaluate(append(boundaryLengths, a.car.Velocity))
		a.car.SteeringPos = steeringAndAccel[0]*2 - 1
		a.car.AccelPos = steeringAndAccel[1]*2 - 1
	}

	return a.isAlive
}

func (a *Agent) Render(renderer *renderer.VoxelArrayObjectRenderer) {
	if a.isAlive {
		a.car.Render(renderer)
	}
}

func NewAgent(id int, carModel *voxelArray.VoxelArrayObject, startingPosition mgl32.Vec2) *Agent {
	agent := Agent{
		car:     vehicle.NewVehicle(id, carModel),
		net:     neural.NewNeuralNet([]int{3, 6, 5, 4}, 2),
		isAlive: true}
	agent.car.Position = startingPosition
	agent.car.AccelPos = rand.Float32()*2 - 1
	agent.car.SteeringPos = rand.Float32()*2 - 1

	return &agent
}
