package vehicle

import (
	"go-experiments/common/commoncolor"
	"go-experiments/voxelli/renderer"
	"go-experiments/voxelli/roadway"
	"go-experiments/voxelli/voxelArray"
	"math"

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
}

const MaxVelocity = 50.0
const MinVelocity = -5.0

const AccelScaleFactor = 25
const SteeringScaleFactor = 10

func boundValue(value float32, min float32, max float32) float32 {
	if value > max {
		value = max
	} else if value < min {
		value = min
	}

	return value
}

func (v *Vehicle) Reset(orientation float32, position mgl32.Vec2) {
	v.Position = position
	v.Orientation = orientation
	v.Velocity = 0
	v.SteeringPos = 0.0
	v.AccelPos = 0.0

	v.Score = 0
}

// Returns the eye position and vectors of the vehicles 'eyes'
func (v *Vehicle) GetEyes() ([]mgl32.Vec2, []mgl32.Vec2) {
	rotation := mgl32.Rotate2D(v.Orientation - math.Pi)

	eyePositions := []mgl32.Vec2{
		v.Position.Add(rotation.Mul2x1(mgl32.Vec2{-v.HalfSize.X(), v.HalfSize.Y() * 2})),
		v.Position.Add(rotation.Mul2x1(mgl32.Vec2{v.HalfSize.X(), v.HalfSize.Y() * 2})),
		v.Position.Add(rotation.Mul2x1(mgl32.Vec2{0, v.HalfSize.Y() * 2}))}

	eyeDirections := []mgl32.Vec2{
		rotation.Mul2x1(mgl32.Vec2{-1, 1}.Normalize()),
		rotation.Mul2x1(mgl32.Vec2{1, 1}.Normalize()),
		rotation.Mul2x1(mgl32.Vec2{0, 1}.Normalize())}

	return eyePositions, eyeDirections
}

// Updates the vehicle, returning true if it hit a wall
func (v *Vehicle) Update(frameTime float32, roadway *roadway.Roadway) bool {
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

	// Stop at the wall if we hit a wall.
	if !roadway.InAllBounds(v.GetBounds()) {
		v.Velocity = 0
		v.Position = oldPosition
		v.Orientation = oldOrientation

		return true
	}

	// Score == distance moved, prioritizing forwards straight-motion travel
	// We also score here to avoid giving points for hitting walls.
	absSteering := math.Abs(float64(v.SteeringPos))
	if v.Velocity > 0 {
		v.Score += step.Len() * float32(math.Pow(1.0-absSteering, 5))
	}

	return false
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

func (v *Vehicle) Render(renderer *renderer.VoxelArrayObjectRenderer, maxScore float32) {
	height := float32(1.0)

	offset := mgl32.Translate3D(-v.HalfSize.X(), -v.HalfSize.Y()*2, height) // 1 bumps us up to the road level, *2 means we rotate from the back and appear to steer.
	rotation := mgl32.HomogRotate3DZ(v.Orientation)
	translation := mgl32.Translate3D(v.Position.X(), v.Position.Y(), 0)

	model := translation.Mul4(rotation.Mul4(offset))

	overlayColor := mgl32.Vec3{1, 1, 1}
	if isColorOverlayEnabled {
		overlayColor = commonColor.LookupColor(v.Score / maxScore)
	}

	renderer.Render(v.Shape, &model, overlayColor)
}

func NewVehicle(id int, shape *voxelArray.VoxelArrayObject) *Vehicle {
	vehicle := Vehicle{Id: id}

	vehicle.Reset(0, mgl32.Vec2{0, 0})
	vehicle.Shape = shape

	// Center the vehicle
	vehicle.HalfSize = mgl32.Vec2{
		float32(vehicle.Shape.VoxelObject.MaxBounds.X()-vehicle.Shape.VoxelObject.MinBounds.X()) / 2,
		float32(vehicle.Shape.VoxelObject.MaxBounds.Y()-vehicle.Shape.VoxelObject.MinBounds.Y()) / 2}

	vehicle.Center = mgl32.Vec2{
		float32(vehicle.Shape.VoxelObject.MinBounds.X()) + vehicle.HalfSize.X(),
		float32(vehicle.Shape.VoxelObject.MinBounds.Y()) + vehicle.HalfSize.Y()}

	return &vehicle
}
