package commonConfig

type SerializableVec3 struct {
	X float32
	Y float32
	Z float32
}

type Projection struct {
	Left   float32
	Right  float32
	Bottom float32
	Top    float32
	Near   float32
	Far    float32
}

type Perspective struct {
	FovY float32
	Near float32
	Far  float32
}

type Window struct {
	Width           int
	Height          int
	Title           string
	BackgroundColor SerializableVec3
	OpenGlMajor     int
	OpenGlMinor     int
}
