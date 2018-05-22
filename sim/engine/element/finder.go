package element

// Defines how to quickly add, remove, and find points on our gameboard.
import (
	"github.com/go-gl/mathgl/mgl32"
)

type ElementFinder struct {
	elements []Element
}

type ElementWithDistance struct {
	Distance    float32
	Element     Element
	SnapNodeIdx int
}

type LineElementWithDistance struct {
	ElementAndDist ElementWithDistance
	LinePercent    float32 // % from start to end.
	LinePoint      mgl32.Vec2
}

func NewElementFinder() *ElementFinder {
	return &ElementFinder{
		elements: make([]Element, 0)}
}

func (e *ElementFinder) Add(elem Element) {
	e.elements = append(e.elements, elem)
}

// Returns true if any points are in range, false otherwise. TODO -- in the long term, this needs to support Region-based intersections
func (e *ElementFinder) IntersectsWithElement(pos mgl32.Vec2, distance float32) bool {
	for _, element := range e.elements {
		region := element.GetRegion()
		if region != nil { // Line-based elements can never interesect. TODO, this may be revisited.
			distanceToPoint := region.Position.Sub(pos).Len()
			if distanceToPoint < distance {
				return true
			}
		}
	}

	return false
}

// Returns the K-nearest elements, via snap ndoes
func (e *ElementFinder) KNearest(pos mgl32.Vec2, count int) []ElementWithDistance {
	elementsWithDistances := make([]ElementWithDistance, count)
	elementsFound := 0

	for _, element := range e.elements {
		snapNodes := element.GetSnapNodes()
		for n, snapNode := range snapNodes {
			distance := snapNode.Sub(pos).Len()

			addedToResultSet := false
			shouldAddElement := elementsFound < count
			for i := 0; i < elementsFound; i++ {
				// Determine if the distance is less than the number of elements found. If so, insert and push everything down.
				if distance < elementsWithDistances[i].Distance {
					addedToResultSet = true

					// This lets us add the new element (instead of popping the furthest away one off of the stack) if the stack is not full.
					if shouldAddElement {
						elementsFound++
					}

					for j := count - 1; j >= i; j-- {
						if j != i {
							elementsWithDistances[j] = elementsWithDistances[j-1]
						} else {
							elementsWithDistances[j] = ElementWithDistance{
								Distance:    distance,
								Element:     element,
								SnapNodeIdx: n}
						}
					}
					break
				}
			}

			// If we have empty space, add this element by-default to the end if not already added.
			if !addedToResultSet && shouldAddElement {
				elementsWithDistances[elementsFound] = ElementWithDistance{
					Distance:    distance,
					Element:     element,
					SnapNodeIdx: n}
				elementsFound++
			}
		}
	}

	return elementsWithDistances[0:elementsFound]
}
