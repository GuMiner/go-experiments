package geometry

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func verifyEffectivelyEqual(t *testing.T, expected, actual mgl32.Vec2) {
	if math.Abs(float64(expected.X()-actual.X())) > 0.00001 ||
		math.Abs(float64(expected.Y()-actual.Y())) > 0.00001 {
		t.Errorf("Vectors do not match. Actual: %v. Expected: %v", actual, expected)
	}
}

func TestSimpleIntersections(t *testing.T) {
	seg := NewLineSegment(mgl32.Vec2{0, 0}, mgl32.Vec2{20, 0})
	vec := NewVector(mgl32.Vec2{10, 1}, mgl32.Vec2{0, -1})

	doesIntersect, intersectPos := seg.Intersects(vec)
	if !doesIntersect {
		t.Error("Vector did not intersect")
	}

	verifyEffectivelyEqual(t, mgl32.Vec2{10, 0}, intersectPos)

	// Swap axis
	seg = NewLineSegment(mgl32.Vec2{0, 0}, mgl32.Vec2{0, 20})
	vec = NewVector(mgl32.Vec2{1, 10}, mgl32.Vec2{-1, 0})

	doesIntersect, intersectPos = seg.Intersects(vec)
	if !doesIntersect {
		t.Error("Vector did not intersect")
	}

	verifyEffectivelyEqual(t, mgl32.Vec2{0, 10}, intersectPos)
}

func TestVectorDirectionIntersection(t *testing.T) {
	seg := NewLineSegment(mgl32.Vec2{0, 0}, mgl32.Vec2{20, 0})
	vec := NewVector(mgl32.Vec2{10, 1}, mgl32.Vec2{0, 1})

	doesIntersect, intersectPos := seg.Intersects(vec)
	if doesIntersect {
		t.Errorf("Vector should not intersect, did at %v", intersectPos)
	}

	// Swap axis
	seg = NewLineSegment(mgl32.Vec2{0, 0}, mgl32.Vec2{0, 20})
	vec = NewVector(mgl32.Vec2{1, 10}, mgl32.Vec2{1, 0})

	doesIntersect, intersectPos = seg.Intersects(vec)
	if doesIntersect {
		t.Errorf("Vector should not intersect, did at %v", intersectPos)
	}
}

func TestCoincidenceIntersection(t *testing.T) {
	seg := NewLineSegment(mgl32.Vec2{0, 0}, mgl32.Vec2{20, 0})
	vec := NewVector(mgl32.Vec2{0, 0}, mgl32.Vec2{1, 0})

	doesIntersect, intersectPos := seg.Intersects(vec)
	if doesIntersect {
		t.Errorf("Vector should not intersect, did at %v", intersectPos)
	}

	// Test with different positions
	vec = NewVector(mgl32.Vec2{10, 0}, mgl32.Vec2{-1, 0})

	doesIntersect, intersectPos = seg.Intersects(vec)
	if doesIntersect {
		t.Errorf("Vector should not intersect, did at %v", intersectPos)
	}

	vec = NewVector(mgl32.Vec2{15, 0}, mgl32.Vec2{1, 0})

	doesIntersect, intersectPos = seg.Intersects(vec)
	if doesIntersect {
		t.Errorf("Vector should not intersect, did at %v", intersectPos)
	}
}

func TestAngledIntersection(t *testing.T) {
	seg := NewLineSegment(mgl32.Vec2{0, 40}, mgl32.Vec2{40, 0})
	vec := NewVector(mgl32.Vec2{2, 2}, mgl32.Vec2{1, 1})

	doesIntersect, intersectPos := seg.Intersects(vec)
	if !doesIntersect {
		t.Error("Vector did not intersect")
	}

	verifyEffectivelyEqual(t, mgl32.Vec2{20, 20}, intersectPos)
}
