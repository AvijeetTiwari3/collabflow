package unit

import (
	"testing"

	"github.com/sde/collabflow/internal/service"
)

func TestCalculatePositionScenarios(t *testing.T) {
	// Case 1: Empty list
	pos, respace := service.CalculatePosition(nil, nil)
	if pos != service.DefaultInitialPosition || respace {
		t.Fatalf("expected initial pos %f, got %f", service.DefaultInitialPosition, pos)
	}

	// Case 2: Prepend before first item (next = 1000.0)
	next := 1000.0
	pos, _ = service.CalculatePosition(nil, &next)
	if pos >= next {
		t.Fatalf("prepended position %f must be less than next %f", pos, next)
	}

	// Case 3: Append after last item (prev = 1000.0)
	prev := 1000.0
	pos, _ = service.CalculatePosition(&prev, nil)
	if pos <= prev {
		t.Fatalf("appended position %f must be greater than prev %f", pos, prev)
	}

	// Case 4: Insert between 1000.0 and 2000.0 -> should be 1500.0
	p1 := 1000.0
	p2 := 2000.0
	pos, _ = service.CalculatePosition(&p1, &p2)
	if !service.FloatEquals(pos, 1500.0, 1e-6) {
		t.Fatalf("expected midpoint 1500.0, got %f", pos)
	}

	// Case 5: Repeated subdivisions between items
	low := 100.0
	high := 200.0
	for i := 0; i < 20; i++ {
		mid, _ := service.CalculatePosition(&low, &high)
		if mid <= low || mid >= high {
			t.Fatalf("iteration %d: midpoint %f out of bounds (%f, %f)", i, mid, low, high)
		}
		high = mid // divide further
	}
}
