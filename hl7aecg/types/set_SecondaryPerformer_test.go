package types

import "testing"

func TestSecondaryPerformer_SetFunctionCode(t *testing.T) {
	sp := &SecondaryPerformer{}

	sp.SetFunctionCode(PERFORMER_ECG_TECHNICIAN, CodeSystemOID(""), "", "")

	if sp.FunctionCode == nil {
		t.Fatal("FunctionCode should not be nil")
	}

	if sp.FunctionCode.Code != PERFORMER_ECG_TECHNICIAN {
		t.Errorf("FunctionCode.Code = %v, want %v", sp.FunctionCode.Code, PERFORMER_ECG_TECHNICIAN)
	}
}

func TestSecondaryPerformer_SetTime(t *testing.T) {
	sp := &SecondaryPerformer{}
	low := "20021122091000"
	high := "20021122091010"

	sp.SetTime(low, high)

	if sp.Time == nil {
		t.Fatal("Time should not be nil")
	}

	if sp.Time.Low.Value != low {
		t.Errorf("Time.Low.Value = %v, want %v", sp.Time.Low.Value, low)
	}

	if sp.Time.High.Value != high {
		t.Errorf("Time.High.Value = %v, want %v", sp.Time.High.Value, high)
	}
}

func TestSecondaryPerformer_SetPerformerID(t *testing.T) {
	sp := &SecondaryPerformer{}
	root := "2.16.840.1.113883.3.4"
	extension := "TECH-221"

	sp.SetPerformerID(root, extension)

	if sp.SeriesPerformer.ID == nil {
		t.Fatal("SeriesPerformer.ID should not be nil")
	}

	if sp.SeriesPerformer.ID.Root != root {
		t.Errorf("SeriesPerformer.ID.Root = %v, want %v", sp.SeriesPerformer.ID.Root, root)
	}

	if sp.SeriesPerformer.ID.Extension != extension {
		t.Errorf("SeriesPerformer.ID.Extension = %v, want %v", sp.SeriesPerformer.ID.Extension, extension)
	}
}

func TestSecondaryPerformer_SetPerformerName(t *testing.T) {
	sp := &SecondaryPerformer{}
	name := "KAB"

	sp.SetPerformerName(name)

	if sp.SeriesPerformer.AssignedPerson == nil {
		t.Fatal("SeriesPerformer.AssignedPerson should not be nil")
	}

	if sp.SeriesPerformer.AssignedPerson.Name == nil {
		t.Fatal("AssignedPerson.Name should not be nil")
	}

	if *sp.SeriesPerformer.AssignedPerson.Name != name {
		t.Errorf("AssignedPerson.Name = %v, want %v", *sp.SeriesPerformer.AssignedPerson.Name, name)
	}
}

func TestSecondaryPerformer_SetEmptyPerformerName(t *testing.T) {
	sp := &SecondaryPerformer{}

	sp.SetEmptyPerformerName()

	if sp.SeriesPerformer.AssignedPerson == nil {
		t.Fatal("SeriesPerformer.AssignedPerson should not be nil")
	}

	if sp.SeriesPerformer.AssignedPerson.Name == nil {
		t.Fatal("AssignedPerson.Name should not be nil")
	}

	if *sp.SeriesPerformer.AssignedPerson.Name != "" {
		t.Errorf("AssignedPerson.Name = %v, want empty string", *sp.SeriesPerformer.AssignedPerson.Name)
	}
}

func TestSecondaryPerformer_MethodChaining(t *testing.T) {
	sp := &SecondaryPerformer{}

	// Test that all methods return the same instance for chaining
	result := sp.
		SetFunctionCode(PERFORMER_ECG_TECHNICIAN, CodeSystemOID(""), "", "").
		SetTime("20021122091000", "20021122091010").
		SetPerformerID("2.16.840.1.113883.3.4", "TECH-221").
		SetPerformerName("KAB")

	if result != sp {
		t.Error("Methods should return the same instance for chaining")
	}

	// Verify all fields were set correctly
	if sp.FunctionCode == nil || sp.FunctionCode.Code != PERFORMER_ECG_TECHNICIAN {
		t.Error("FunctionCode not set correctly")
	}

	if sp.Time == nil || sp.Time.Low.Value != "20021122091000" {
		t.Error("Time not set correctly")
	}

	if sp.SeriesPerformer.ID == nil || sp.SeriesPerformer.ID.Root != "2.16.840.1.113883.3.4" {
		t.Error("Performer ID not set correctly")
	}

	if sp.SeriesPerformer.AssignedPerson == nil || *sp.SeriesPerformer.AssignedPerson.Name != "KAB" {
		t.Error("Performer name not set correctly")
	}
}

func TestSeriesPerformer_SetPerformer(t *testing.T) {
	tests := []struct {
		name               string
		performerID        string
		performerExtension string
		personName         string
		wantIDSet          bool
	}{
		{
			name:               "Complete performer",
			performerID:        "2.16.840.1.113883.3.4",
			performerExtension: "TECH-221",
			personName:         "KAB",
			wantIDSet:          true,
		},
		{
			name:               "Performer with ID only",
			performerID:        "2.16.840.1.113883.3.4",
			performerExtension: "",
			personName:         "KAB",
			wantIDSet:          true,
		},
		{
			name:               "Performer with name only",
			performerID:        "",
			performerExtension: "",
			personName:         "KAB",
			wantIDSet:          false,
		},
		{
			name:               "Empty name element",
			performerID:        "",
			performerExtension: "",
			personName:         "",
			wantIDSet:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spe := &SeriesPerformer{}

			spe.SetPerformer(tt.performerID, tt.performerExtension, tt.personName)

			// Check ID
			if tt.wantIDSet {
				if spe.ID == nil {
					t.Error("ID should be set")
				} else {
					if spe.ID.Root != tt.performerID {
						t.Errorf("ID.Root = %v, want %v", spe.ID.Root, tt.performerID)
					}
					if spe.ID.Extension != tt.performerExtension {
						t.Errorf("ID.Extension = %v, want %v", spe.ID.Extension, tt.performerExtension)
					}
				}
			}

			// Check name
			if spe.AssignedPerson == nil {
				t.Fatal("AssignedPerson should not be nil")
			}

			if spe.AssignedPerson.Name == nil {
				t.Fatal("AssignedPerson.Name should not be nil")
			}

			if *spe.AssignedPerson.Name != tt.personName {
				t.Errorf("AssignedPerson.Name = %v, want %v", *spe.AssignedPerson.Name, tt.personName)
			}
		})
	}
}
