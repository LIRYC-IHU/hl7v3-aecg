package types

// SetActivityTime sets the activity time range for the clinical trial.
//
// Parameters:
//   - low: Start date of the trial (e.g., "20010509")
//   - high: End date of the trial (e.g., "20020316")
//
// Format: HL7 TS (Timestamp)
//   - YYYYMMDDHHmmss (with time)
//   - YYYYMMDD (date only)
//   - YYYYMM (year-month)
//   - YYYY (year only)
//
// Example:
//
//	ct.SetActivityTime("20010509", "20020316")
//
// Returns the ClinicalTrial for method chaining.
func (ct *ClinicalTrial) SetActivityTime(low, high string) *ClinicalTrial {
	if ct.ActivityTime == nil {
		ct.ActivityTime = &EffectiveTime{}
	}
	ct.ActivityTime.Low = Time{Value: low}
	ct.ActivityTime.High = Time{Value: high}
	return ct
}

// SetActivityTimeLow sets only the low (start) activity time.
//
// Parameters:
//   - low: Start date of the trial (e.g., "20010509")
//
// Returns the ClinicalTrial for method chaining.
func (ct *ClinicalTrial) SetActivityTimeLow(low string) *ClinicalTrial {
	if ct.ActivityTime == nil {
		ct.ActivityTime = &EffectiveTime{}
	}
	ct.ActivityTime.Low = Time{Value: low}
	return ct
}

// SetActivityTimeHigh sets only the high (end) activity time.
//
// Parameters:
//   - high: End date of the trial (e.g., "20020316")
//
// Returns the ClinicalTrial for method chaining.
func (ct *ClinicalTrial) SetActivityTimeHigh(high string) *ClinicalTrial {
	if ct.ActivityTime == nil {
		ct.ActivityTime = &EffectiveTime{}
	}
	ct.ActivityTime.High = Time{Value: high}
	return ct
}
