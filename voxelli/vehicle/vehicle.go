package vehicle

import (
	"fmt"
	"go-experiments/voxelli/renderer"
	"go-experiments/voxelli/voxel"
	"go-experiments/voxelli/voxelArray"

	"github.com/go-gl/mathgl/mgl32"
)

type Vehicle struct {
	Position    mgl32.Vec2
	Orientation float32

	Velocity mgl32.Vec2

	// -1: Hard left. 1: Hard Right
	SteeringPos float32

	// -1: Slamming the brake. 0: Coast. 1: Flooring it
	AccelPos float32

	Shape    *voxelArray.VoxelArrayObject
	Center   mgl32.Vec2 // Offset for the vehicle center.
	HalfSize mgl32.Vec2
}

func (v *Vehicle) Render(renderer *renderer.VoxelArrayObjectRenderer) {
	offset := mgl32.Translate3D(-v.HalfSize.X(), -v.HalfSize.Y(), 1) // 1 bumps us up to the road level.
	rotation := mgl32.HomogRotate3DZ(v.Orientation)

	model := rotation.Mul4(offset)
	renderer.Render(v.Shape, &model)
}

func (v *Vehicle) Delete() {
	v.Shape.Delete()
}

func NewVehicle(model string) *Vehicle {
	vehicle := Vehicle{
		SteeringPos: 0, AccelPos: 0,
		Velocity: mgl32.Vec2{0, 0},
		Position: mgl32.Vec2{0, 0}}

	carRaw := voxel.NewVoxelObject(model)
	fmt.Printf("Vehicle objects: %v\n", len(carRaw.SubObjects))

	vehicle.Shape = voxelArray.NewVoxelArrayObject(carRaw)
	fmt.Printf("Optimized Vehicle vertices: %v\n\n", vehicle.Shape.Vertices)

	// Center the vehicle
	vehicle.HalfSize = mgl32.Vec2{
		float32(vehicle.Shape.VoxelObject.MaxBounds.X() - vehicle.Shape.VoxelObject.MinBounds.X()),
		float32(vehicle.Shape.VoxelObject.MaxBounds.Y() - vehicle.Shape.VoxelObject.MinBounds.Y())}

	vehicle.Center = mgl32.Vec2{
		float32(vehicle.Shape.VoxelObject.MinBounds.X()) + vehicle.HalfSize.X(),
		float32(vehicle.Shape.VoxelObject.MinBounds.Y()) + vehicle.HalfSize.Y()}

	return &vehicle
}
