package types

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
