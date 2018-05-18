package commonConfig

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
	Samples         int
}
