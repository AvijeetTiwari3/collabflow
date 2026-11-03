package service

import (
	"math"
)

const (
	DefaultInitialPosition = 1000.0
	DefaultStep            = 1000.0
	MinEpsilon             = 1e-9 // Threshold before re-spacing is triggered
)

// CalculatePosition computes an O(1) fractional index between two tasks
func CalculatePosition(prevPos, nextPos *float64) (float64, bool) {
	// Case 1: First item in empty list
	if prevPos == nil && nextPos == nil {
		return DefaultInitialPosition, false
	}

	// Case 2: Inserted at beginning of list (before first item)
	if prevPos == nil && nextPos != nil {
		if *nextPos <= 1.0 {
			return *nextPos / 2.0, false
		}
		return *nextPos - DefaultStep, false
	}

	// Case 3: Inserted at end of list (after last item)
	if prevPos != nil && nextPos == nil {
		return *prevPos + DefaultStep, false
	}

	// Case 4: Inserted between two items
	prev := *prevPos
	next := *nextPos

	// Sanity check: prev must be strictly less than next
	if prev >= next {
		return prev + 1.0, true // Flag requires re-spacing
	}

	diff := next - prev
	if diff < MinEpsilon {
		// Floating point precision limit reached, trigger re-spacing
		return (prev + next) / 2.0, true
	}

	// Midpoint
	return prev + (diff / 2.0), false
}

// RespacePositions generates cleanly spaced positions for an array of items
func RespacePositions(count int) []float64 {
	positions := make([]float64, count)
	for i := 0; i < count; i++ {
		positions[i] = float64(i+1) * DefaultStep
	}
	return positions
}

// FloatEquals helper for unit tests
func FloatEquals(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}
