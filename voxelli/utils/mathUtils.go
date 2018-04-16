package utils

// Returns the maximum of a and b
func MaxFloat32(a, b float32) float32 {
	if a > b {
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
