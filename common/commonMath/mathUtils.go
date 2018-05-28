package commonMath

import "github.com/go-gl/mathgl/mgl32"

func MaxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}

	return b
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func MinFloat32(a, b float32) float32 {
	if a < b {
		return a
	}

	return b
}

func MinInt32(a, b int32) int32 {
	if a > b {
		return b
	}

	return a
}

// Takes in parameters to generate a look-at matrix, returning the inverse of the rotation
func InverseLookAtRotationMatrix(eye, center, up mgl32.Vec3) mgl32.Mat4 {
	f := center.Sub(eye).Normalize()
	s := f.Cross(up.Normalize()).Normalize()
	u := s.Cross(f)

	M := mgl32.Mat4{
		s[0], u[0], -f[0], 0,
		s[1], u[1], -f[1], 0,
		s[2], u[2], -f[2], 0,
		0, 0, 0, 1,
	}

	return M.Transpose()
}
