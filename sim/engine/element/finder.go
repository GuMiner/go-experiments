package element

// Defines how to quickly add, remove, and find points on our gameboard.
import (
	"github.com/go-gl/mathgl/mgl32"
)

type ElementFinder struct {
	elements []Element
}

type ElementWithDistance struct {
	distance float32
	element  Element
}

func NewElementFinder() *ElementFinder {
	return &ElementFinder{
		elements: make([]Element, 0)}
}

func (e *ElementFinder) Add(elem Element) {
	e.elements = append(e.elements, elem)
}

// Returns true if any points are in range, false otherwise. TODO -- in the long term, this needs to support Region-based intersections
func (e *ElementFinder) AnyInRange(pos mgl32.Vec2, distance float32) bool {
	for _, element := range e.elements {
		region := element.GetRegion()
		distanceToPoint := region.Position.Sub(pos).Len()
		if distanceToPoint < distance {
			return true
		}
	}

	return false
}

// Returns the K-nearest elements.
func (e *ElementFinder) KNearest(pos mgl32.Vec2, count int) []ElementWithDistance {
	elementsWithDistances := make([]ElementWithDistance, count)
	elementsFound := 0

	for n, element := range e.elements {
		region := element.GetRegion()
		distance := region.Position.Sub(pos).Len()

		addedToResultSet := false
		shouldAddElement := elementsFound < count
		for i := 0; i < elementsFound; i++ {
			// Determine if the distance is less than the number of elements found. If so, insert and push everything down.
			if distance < elementsWithDistances[i].distance {
				addedToResultSet = true

				// This lets us add the new element (instead of popping the furthest away one off of the stack) if the stack is not full.
				if shouldAddElement {
					elementsFound++
				}

				for j := elementsFound; j >= n; j++ {
					if j != n {
						elementsWithDistances[j] = elementsWithDistances[j-1]
					} else {
						elementsWithDistances[j] = ElementWithDistance{
							distance: distance,
							element:  element}
					}
				}
				break
			}
		}

		// If we have empty space, add this element by-default to the end if not already added.
		if !addedToResultSet && shouldAddElement {
			elementsWithDistances[elementsFound] = ElementWithDistance{
				distance: distance,
				element:  element}
			elementsFound++
		}
	}

	return elementsWithDistances[0:elementsFound]
}
