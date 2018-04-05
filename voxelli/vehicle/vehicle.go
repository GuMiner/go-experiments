package vehicle

import (
	"go-experiments/voxelli/renderer"
	"go-experiments/voxelli/roadway"
	"go-experiments/voxelli/voxelArray"
	"math"
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
)

type Vehicle struct {
	Position    mgl32.Vec2
	Orientation float32

	Velocity float32 // TODO -- make this a vector to allow skidding

	// 1: Hard left. -1: Hard Right
	SteeringPos float32

	// -1: Slamming the brake. 0: Coast. 1: Flooring it
	AccelPos float32

	// We expect the vehicle to be oriented such that Y- == forwards
	// X+ is then the left-side of the vehicle.
	Shape    *voxelArray.VoxelArrayObject
	Center   mgl32.Vec2 // Offset for the vehicle center.
	HalfSize mgl32.Vec2

	Id    int
	Score float32

	RandomizeOnWallHit bool
}

const MaxVelocity = 200.0
const MinVelocity = -50.0

const AccelScaleFactor = 20
const SteeringScaleFactor = 30

func boundValue(value float32, min float32, max float32) float32 {
	if value > max {
		value = max
	} else if value < min {
		value = min
	}

	return value
}

// Returns the eye position and vectors of the vehicles 'eyes'
func (v *Vehicle) GetEyes() ([]mgl32.Vec2, []mgl32.Vec2) {
	rotation := mgl32.Rotate2D(v.Orientation - math.Pi)

	eyePositions := []mgl32.Vec2{
		v.Position.Add(rotation.Mul2x1(mgl32.Vec2{-v.HalfSize.X(), v.HalfSize.Y() * 2})),
		v.Position.Add(rotation.Mul2x1(mgl32.Vec2{v.HalfSize.X(), v.HalfSize.Y() * 2})),
	}

	eyeDirections := []mgl32.Vec2{
		rotation.Mul2x1(mgl32.Vec2{-1, 1}.Normalize()),
		rotation.Mul2x1(mgl32.Vec2{1, 1}.Normalize())}

	return eyePositions, eyeDirections
}

// Updates the vehicle,
func (v *Vehicle) Update(frameTime float32, roadway *roadway.Roadway) {
	v.AccelPos = boundValue(v.AccelPos, -1.0, 1.0)
	v.SteeringPos = boundValue(v.SteeringPos, -1.0, 1.0)

	v.Velocity += v.AccelPos * AccelScaleFactor * frameTime
	v.Velocity = boundValue(v.Velocity, MinVelocity, MaxVelocity)

	oldOrientation := v.Orientation
	v.Orientation += v.SteeringPos * SteeringScaleFactor * (v.Velocity / MaxVelocity) * frameTime
	for v.Orientation > 2*math.Pi {
		v.Orientation -= 2 * math.Pi
	}

	for v.Orientation < 0 {
		v.Orientation += 2 * math.Pi
	}

	xStep := float32(math.Cos(float64(v.Orientation-math.Pi/2))) * v.Velocity * frameTime
	yStep := float32(math.Sin(float64(v.Orientation-math.Pi/2))) * v.Velocity * frameTime

	step := mgl32.Vec2{xStep, yStep}
	oldPosition := v.Position
	v.Position = v.Position.Add(step)

	// Score == distance moved, significantly prioritizing straight-motion travel.
	v.Score += step.Len() * float32(math.Pow(1.0-math.Abs(float64(v.SteeringPos)), 8))

	// Stop at the wall if we hit a wall.
	if !roadway.InAllBounds(v.GetBounds()) {

		// Randomize the accelerator and steering after hitting a wall
		// Remove once the neural net starts driving
		if v.RandomizeOnWallHit {
			v.AccelPos = rand.Float32()*2 - 1
			v.SteeringPos = rand.Float32()*2 - 1
		}

		v.Velocity = 0
		v.Position = oldPosition
		v.Orientation = oldOrientation
		v.Score = 0
	}
}

// Returns the bounds of the vehicle in CW order, starting from the left bumper
func (v *Vehicle) GetBounds() []mgl32.Vec2 {
	rotation := mgl32.Rotate2D(v.Orientation - math.Pi)

	bounds := []mgl32.Vec2{
		v.Position.Add(rotation.Mul2x1(mgl32.Vec2{-v.HalfSize.X(), 0})),
		v.Position.Add(rotation.Mul2x1(mgl32.Vec2{-v.HalfSize.X(), v.HalfSize.Y() * 2})),
		v.Position.Add(rotation.Mul2x1(mgl32.Vec2{v.HalfSize.X(), v.HalfSize.Y() * 2})),
		v.Position.Add(rotation.Mul2x1(mgl32.Vec2{v.HalfSize.X(), 0})),
	}

	return bounds
}

func (v *Vehicle) Render(emphasize bool, renderer *renderer.VoxelArrayObjectRenderer) {
	var height float32 = 1.0
	if emphasize {
		height = 5
	}

	offset := mgl32.Translate3D(-v.HalfSize.X(), -v.HalfSize.Y()*2, height) // 1 bumps us up to the road level, *2 means we rotate from the back and appear to steer.
	rotation := mgl32.HomogRotate3DZ(v.Orientation)
	translation := mgl32.Translate3D(v.Position.X(), v.Position.Y(), 0)

	model := translation.Mul4(rotation.Mul4(offset))
	renderer.Render(v.Shape, &model)
}

func NewVehicle(id int, shape *voxelArray.VoxelArrayObject) *Vehicle {
	vehicle := Vehicle{
		Id:          id,
		Score:       0.0,
		SteeringPos: 0, AccelPos: 0,
		Velocity: 0.0, Orientation: 0.0, Position: mgl32.Vec2{0, 0}}

	vehicle.Shape = shape

	// Center the vehicle
	vehicle.HalfSize = mgl32.Vec2{
		float32(vehicle.Shape.VoxelObject.MaxBounds.X()-vehicle.Shape.VoxelObject.MinBounds.X()) / 2,
		float32(vehicle.Shape.VoxelObject.MaxBounds.Y()-vehicle.Shape.VoxelObject.MinBounds.Y()) / 2}

	vehicle.Center = mgl32.Vec2{
		float32(vehicle.Shape.VoxelObject.MinBounds.X()) + vehicle.HalfSize.X(),
		float32(vehicle.Shape.VoxelObject.MinBounds.Y()) + vehicle.HalfSize.Y()}

	vehicle.RandomizeOnWallHit = true
	return &vehicle
}
