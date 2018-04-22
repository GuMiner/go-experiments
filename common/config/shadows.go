package commonConfig

type OrthoProjection struct {
	Left   float32
	Right  float32
	Bottom float32
	Top    float32
	Near   float32
	Far    float32
}

type Shadows struct {
	Projection OrthoProjection
	Position   SerializableVec3
	Forwards   SerializableVec3
	Up         SerializableVec3
}
