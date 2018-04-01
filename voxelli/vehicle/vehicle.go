package vehicle

import (
	"fmt"
	"go-experiments/voxelli/renderer"
	"go-experiments/voxelli/voxelArray"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Vehicle struct {
	Position    mgl32.Vec2
	Orientation float32

	Velocity float32 // TODO -- make this a vector to allow skidding

	// -1: Hard left. 1: Hard Right
	SteeringPos float32

	// -1: Slamming the brake. 0: Coast. 1: Flooring it
	AccelPos float32

	// We expect the vehicle to be oriented such that Y- == forwards
	// X+ is then the left-side of the vehicle.
	Shape    *voxelArray.VoxelArrayObject
	Center   mgl32.Vec2 // Offset for the vehicle center.
	HalfSize mgl32.Vec2
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

func (v *Vehicle) Update(frameTime float32) (mgl32.Vec2, float32) {
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

	oldPosition := v.Position
	v.Position = v.Position.Add(mgl32.Vec2{xStep, yStep})

	return oldPosition, oldOrientation
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

func (v *Vehicle) Render(renderer *renderer.VoxelArrayObjectRenderer) {
	offset := mgl32.Translate3D(-v.HalfSize.X(), -v.HalfSize.Y()*2, 1) // 1 bumps us up to the road level, *2 means we rotate from the back and appear to steer.
	rotation := mgl32.HomogRotate3DZ(v.Orientation)
	translation := mgl32.Translate3D(v.Position.X(), v.Position.Y(), 0)

	model := translation.Mul4(rotation.Mul4(offset))
	renderer.Render(v.Shape, &model)
}

func NewVehicle(shape *voxelArray.VoxelArrayObject) *Vehicle {
	vehicle := Vehicle{
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

	fmt.Printf("Vehicle half size: %v. Center: %v.\n\n", vehicle.HalfSize, vehicle.Center)
	return &vehicle
}
