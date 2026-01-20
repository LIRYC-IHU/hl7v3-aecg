package types

import "math"

// =============================================================================
// Helper Methods
// =============================================================================

// IsEmpty returns true if the EffectiveTime is empty (i.e., both Low and High are empty).
func (et EffectiveTime) IsEmpty() bool {
	return et.Low.Value == "" && et.High.Value == ""
}

// IsEmpty returns true if the Time is empty (i.e., has no Value).
func (t Time) IsEmpty() bool {
	return t.Value == ""
}

// isInvalidFloat returns true if the float64 value is NaN or Inf, which are invalid for medical measurements.
func isInvalidFloat(value float64) bool {
	return math.IsNaN(value) || math.IsInf(value, 0)
}
