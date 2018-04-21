package main

// Defines the data and data flow for the main simulation engine
import (
	"fmt"
	"go-experiments/voxelli/config"
	"go-experiments/voxelli/diagnostics"
	"go-experiments/voxelli/genetics"
	"go-experiments/voxelli/renderer"
	"go-experiments/voxelli/roadway"
	"go-experiments/voxelli/vehicle"
	"go-experiments/voxelli/voxel"
	"go-experiments/voxelli/voxelArray"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func debugDrawCarInfo(car *vehicle.Vehicle, boundaries []float32) {
	eyePositions, eyeDirections := car.GetEyes()

	// Debug draw where we are looking, assuming two eyes only
	model := mgl32.Translate3D(eyePositions[0].X()+2*eyeDirections[0].X(), eyePositions[0].Y()+2*eyeDirections[0].Y(), 8)
	diagnostics.Render(mgl32.Vec4{0.0, 1.0, 0.0, 1.0}, &model)

	model = mgl32.Translate3D(eyePositions[1].X()+2*eyeDirections[1].X(), eyePositions[1].Y()+2*eyeDirections[1].Y(), 8)
	diagnostics.Render(mgl32.Vec4{0.0, 1.0, 1.0, 1.0}, &model) // Cyan

	model = mgl32.Translate3D(eyePositions[2].X()+2*eyeDirections[2].X(), eyePositions[2].Y()+2*eyeDirections[2].Y(), 8)
	diagnostics.Render(mgl32.Vec4{1.0, 1.0, 1.0, 1.0}, &model) // Cyan

	// Debug draw where the boundaries are.
	eyePositions[0] = eyePositions[0].Add(eyeDirections[0].Mul(boundaries[0]))
	eyePositions[1] = eyePositions[1].Add(eyeDirections[1].Mul(boundaries[1]))
	eyePositions[2] = eyePositions[2].Add(eyeDirections[2].Mul(boundaries[2]))

	model = mgl32.Translate3D(eyePositions[0].X(), eyePositions[0].Y(), 8)
	diagnostics.Render(mgl32.Vec4{1.0, 0.0, 0.0, 1.0}, &model)

	model = mgl32.Translate3D(eyePositions[1].X(), eyePositions[1].Y(), 8)
	diagnostics.Render(mgl32.Vec4{1.0, 1.0, 0.0, 1.0}, &model) // Yellow

	model = mgl32.Translate3D(eyePositions[2].X(), eyePositions[2].Y(), 8)
	diagnostics.Render(mgl32.Vec4{1.0, 1.0, 1.0, 1.0}, &model) // Yellow
}

var simpleRoadway *roadway.Roadway
var roadwayDisplayer *roadway.RoadwayDisplayer
var carModel *voxelArray.VoxelArrayObject

var agentEvolver *genetics.Population

func InitSimulation(voxelArrayObjectRenderer *renderer.VoxelArrayObjectRenderer) {
	// Rendering
	simpleRoadway = roadway.NewRoadway("./data/roadways/straight_with_s-curve.txt")
	roadwayDisplayer = roadway.NewRoadwayDisplayer(voxelArrayObjectRenderer)

	carVoxelObject := voxel.NewVoxelObject("./data/models/car.vox")
	fmt.Printf("Vehicle objects: %v\n", len(carVoxelObject.SubObjects))

	carModel = voxelArray.NewVoxelArrayObject(carVoxelObject)
	fmt.Printf("Optimized Vehicle vertices: %v\n\n", carModel.Vertices)

	// Simulation
	agentEvolver = genetics.NewPopulation(
		config.Config.Simulation.PopulationSize,
		func(id int) *genetics.Agent {
			return genetics.NewAgent(id, carModel, float32(math.Pi*0.5), mgl32.Vec2{8, 38})
		})
}

func UpdateSimulation(frameTime, elapsedTime float32) {
	// Update and render the agents
	agentEvolver.Update(frameTime, func(agent *genetics.Agent) {
		agent.Update(frameTime, simpleRoadway)
	})
}

func RenderSimulation(voxelArrayObjectRenderer *renderer.VoxelArrayObjectRenderer) {
	roadwayDisplayer.Render(simpleRoadway)

	agentEvolver.Render(func(agent *genetics.Agent, maxScore float32) {
		if diagnostics.IsDebug() {
			eyePositions, eyeDirections := agent.GetCar().GetEyes()
			boundaryLengths, _ := simpleRoadway.GetBoundaries(eyePositions, eyeDirections)
			debugDrawCarInfo(agent.GetCar(), boundaryLengths)
		}

		agent.Render(voxelArrayObjectRenderer, maxScore)
	})
}

func DeleteSimulation() {
	roadwayDisplayer.Delete()
	carModel.Delete()
}
