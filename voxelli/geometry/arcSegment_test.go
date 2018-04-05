package geometry

import (
	"math"
	"testing"

	"github.com/go-gl/mathgl/mgl32"
)

func TestEdgeArcIntersections(t *testing.T) {
	seg := NewArcSegment(mgl32.Vec2{2, 2}, 2.0, math.Pi, 2*math.Pi)

	vec := NewVector(mgl32.Vec2{-1, -1}, mgl32.Vec2{0, 1})
	doesIntersect, intersectPos := seg.Intersects(vec)
	if doesIntersect {
		t.Error("Vector did intersect")
	}

	vec = NewVector(mgl32.Vec2{-1, -1}, mgl32.Vec2{1, 0})
	doesIntersect, _ = seg.Intersects(vec)
	if doesIntersect {
		t.Error("Vector did intersect")
	}

	vec = NewVector(mgl32.Vec2{-1, -1}, mgl32.Vec2{1, -1})
	doesIntersect, _ = seg.Intersects(vec)
	if doesIntersect {
		t.Error("Vector did intersect")
	}

	vec = NewVector(mgl32.Vec2{0, 0}, mgl32.Vec2{1, 0})
	doesIntersect, intersectPos = seg.Intersects(vec)
	if !doesIntersect {
		t.Error("Vector did not intersect")
	}

	verifyEffectivelyEqual(t, mgl32.Vec2{2, 0}, intersectPos)
}

func TestMoreArcIntersections(t *testing.T) {
	seg := NewArcSegment(mgl32.Vec2{2, 2}, 2.0, math.Pi, 2*math.Pi)

	vec := NewVector(mgl32.Vec2{0, 0}, mgl32.Vec2{1, 1})
	doesIntersect, intersectPos := seg.Intersects(vec)
	if !doesIntersect {
		t.Error("Vector did not intersect")
	}

	verifyEffectivelyEqual(t, mgl32.Vec2{0.5857865, 0.5857865}, intersectPos)

	vec = NewVector(mgl32.Vec2{2, -2}, mgl32.Vec2{0, 1})
	doesIntersect, intersectPos = seg.Intersects(vec)
	if !doesIntersect {
		t.Error("Vector did not intersect")
	}

	verifyEffectivelyEqual(t, mgl32.Vec2{2, 0}, intersectPos)
}

func TestSimpleArcIntersections(t *testing.T) {
	// Full circle
	seg := NewArcSegment(mgl32.Vec2{0, 0}, 1.0, 0.0, 2*math.Pi)
	vec := NewVector(mgl32.Vec2{0, 0}, mgl32.Vec2{0, 1})

	doesIntersect, intersectPos := seg.Intersects(vec)
	if !doesIntersect {
		t.Error("Vector did not intersect")
	}

	verifyEffectivelyEqual(t, mgl32.Vec2{0, 1}, intersectPos)

	// Bottom of circle
	vec = NewVector(mgl32.Vec2{0, -10}, mgl32.Vec2{0, 1})

	doesIntersect, intersectPos = seg.Intersects(vec)
	if !doesIntersect {
		t.Error("Vector did not intersect")
	}

	verifyEffectivelyEqual(t, mgl32.Vec2{0, -1}, intersectPos)

	// Half circle
	seg = NewArcSegment(mgl32.Vec2{0, 0}, 1.0, 0.0, math.Pi)
	vec = NewVector(mgl32.Vec2{0, 0}, mgl32.Vec2{0, 1})

	doesIntersect, intersectPos = seg.Intersects(vec)
	if !doesIntersect {
		t.Error("Vector did not intersect")
	}

	verifyEffectivelyEqual(t, mgl32.Vec2{0, 1}, intersectPos)

	// Top of circle
	vec = NewVector(mgl32.Vec2{0, -2}, mgl32.Vec2{0, 1})

	doesIntersect, intersectPos = seg.Intersects(vec)
	if !doesIntersect {
		t.Error("Vector did not intersect")
		return
	}

	verifyEffectivelyEqual(t, mgl32.Vec2{0, 1}, intersectPos)

	// Misses half circle
	vec = NewVector(mgl32.Vec2{0, 0}, mgl32.Vec2{0, -1})

	doesIntersect, intersectPos = seg.Intersects(vec)
	if doesIntersect {
		t.Error("Vector did intersect")
	}
}
