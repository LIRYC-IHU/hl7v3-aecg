package types

import "testing"

func TestClinicalTrial_SetActivityTime(t *testing.T) {
	tests := []struct {
		name     string
		low      string
		high     string
		wantLow  string
		wantHigh string
	}{
		{
			name:     "Set full activity time range",
			low:      "20010509",
			high:     "20020316",
			wantLow:  "20010509",
			wantHigh: "20020316",
		},
		{
			name:     "Set activity time with timestamps",
			low:      "20010509120000",
			high:     "20020316153000",
			wantLow:  "20010509120000",
			wantHigh: "20020316153000",
		},
		{
			name:     "Set activity time with empty high",
			low:      "20010509",
			high:     "",
			wantLow:  "20010509",
			wantHigh: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := &ClinicalTrial{}

			result := ct.SetActivityTime(tt.low, tt.high)

			// Check method chaining
			if result != ct {
				t.Error("SetActivityTime should return the ClinicalTrial for chaining")
			}

			// Check ActivityTime was created
			if ct.ActivityTime == nil {
				t.Fatal("ActivityTime should not be nil")
			}

			// Check low value
			if ct.ActivityTime.Low.Value != tt.wantLow {
				t.Errorf("Low = %q, want %q", ct.ActivityTime.Low.Value, tt.wantLow)
			}

			// Check high value
			if ct.ActivityTime.High.Value != tt.wantHigh {
				t.Errorf("High = %q, want %q", ct.ActivityTime.High.Value, tt.wantHigh)
			}
		})
	}
}

func TestClinicalTrial_SetActivityTimeLow(t *testing.T) {
	tests := []struct {
		name    string
		low     string
		wantLow string
	}{
		{
			name:    "Set low activity time",
			low:     "20010509",
			wantLow: "20010509",
		},
		{
			name:    "Set low with timestamp",
			low:     "20010509120000",
			wantLow: "20010509120000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := &ClinicalTrial{}

			result := ct.SetActivityTimeLow(tt.low)

			// Check method chaining
			if result != ct {
				t.Error("SetActivityTimeLow should return the ClinicalTrial for chaining")
			}

			// Check ActivityTime was created
			if ct.ActivityTime == nil {
				t.Fatal("ActivityTime should not be nil")
			}

			// Check low value
			if ct.ActivityTime.Low.Value != tt.wantLow {
				t.Errorf("Low = %q, want %q", ct.ActivityTime.Low.Value, tt.wantLow)
			}
		})
	}
}

func TestClinicalTrial_SetActivityTimeHigh(t *testing.T) {
	tests := []struct {
		name     string
		high     string
		wantHigh string
	}{
		{
			name:     "Set high activity time",
			high:     "20020316",
			wantHigh: "20020316",
		},
		{
			name:     "Set high with timestamp",
			high:     "20020316153000",
			wantHigh: "20020316153000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := &ClinicalTrial{}

			result := ct.SetActivityTimeHigh(tt.high)

			// Check method chaining
			if result != ct {
				t.Error("SetActivityTimeHigh should return the ClinicalTrial for chaining")
			}

			// Check ActivityTime was created
			if ct.ActivityTime == nil {
				t.Fatal("ActivityTime should not be nil")
			}

			// Check high value
			if ct.ActivityTime.High.Value != tt.wantHigh {
				t.Errorf("High = %q, want %q", ct.ActivityTime.High.Value, tt.wantHigh)
			}
		})
	}
}

func TestClinicalTrial_SetActivityTime_MethodChaining(t *testing.T) {
	ct := &ClinicalTrial{}

	// Test method chaining
	result := ct.
		SetActivityTimeLow("20010509").
		SetActivityTimeHigh("20020316")

	if result != ct {
		t.Error("Method chaining should return the same ClinicalTrial instance")
	}

	if ct.ActivityTime == nil {
		t.Fatal("ActivityTime should not be nil")
	}

	if ct.ActivityTime.Low.Value != "20010509" {
		t.Errorf("Low = %q, want %q", ct.ActivityTime.Low.Value, "20010509")
	}

	if ct.ActivityTime.High.Value != "20020316" {
		t.Errorf("High = %q, want %q", ct.ActivityTime.High.Value, "20020316")
	}
}

func TestClinicalTrial_SetActivityTime_UpdateExisting(t *testing.T) {
	ct := &ClinicalTrial{
		ActivityTime: &EffectiveTime{
			Low:  Time{Value: "20000101"},
			High: Time{Value: "20001231"},
		},
	}

	// Update existing ActivityTime
	ct.SetActivityTime("20010509", "20020316")

	if ct.ActivityTime.Low.Value != "20010509" {
		t.Errorf("Low = %q, want %q", ct.ActivityTime.Low.Value, "20010509")
	}

	if ct.ActivityTime.High.Value != "20020316" {
		t.Errorf("High = %q, want %q", ct.ActivityTime.High.Value, "20020316")
	}
}
