package types

import (
	"context"
	"sync"
	"testing"
)

// TestID_Validate_Autocomplete tests that empty Root IDs are autocompleted with singleton
func TestID_Validate_Autocomplete(t *testing.T) {
	// Setup singleton
	aecg := &HL7AEcg{}
	aecg.SetRootID("755.3045256.2025923.103550", "")

	tests := []struct {
		name          string
		id            ID
		wantRoot      string
		wantError     bool
		setupSingleton bool
	}{
		{
			name: "Empty root autocompleted with singleton",
			id: ID{
				Root:      "",
				Extension: "test",
			},
			wantRoot:      "755.3045256.2025923.103550",
			wantError:     false,
			setupSingleton: true,
		},
		{
			name: "Empty root without singleton fails",
			id: ID{
				Root:      "",
				Extension: "test",
			},
			wantRoot:      "",
			wantError:     true,
			setupSingleton: false,
		},
		{
			name: "Non-empty root unchanged",
			id: ID{
				Root:      "2.16.840.1.113883.3.1234",
				Extension: "test",
			},
			wantRoot:      "2.16.840.1.113883.3.1234",
			wantError:     false,
			setupSingleton: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset singleton for tests that don't want it
			if !tt.setupSingleton {
				// Save current singleton
				savedInstance := instanceID
				savedOnce := once
				// Reset for this test
				instanceID = nil
				once = *new(sync.Once)
				// Restore after test
				defer func() {
					instanceID = savedInstance
					once = savedOnce
				}()
			}

			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.id.Validate(ctx, vctx)

			// Check if root was autocompleted
			if tt.id.Root != tt.wantRoot {
				t.Errorf("Root after validation = %v, want %v", tt.id.Root, tt.wantRoot)
			}

			// Check validation errors
			hasError := vctx.HasErrors()
			if hasError != tt.wantError {
				if tt.wantError {
					t.Errorf("Expected validation error but got none")
				} else {
					t.Errorf("Unexpected validation error: %v", vctx.GetError())
				}
			}
		})
	}
}

// TestID_Validate_AutocompletePreservesExtension tests that autocomplete preserves extension
func TestID_Validate_AutocompletePreservesExtension(t *testing.T) {
	// Setup singleton
	aecg := &HL7AEcg{}
	aecg.SetRootID("755.3045256.2025923.103550", "")

	id := ID{
		Root:      "",
		Extension: "myExtension",
	}

	ctx := context.Background()
	vctx := NewValidationContext(false)

	id.Validate(ctx, vctx)

	// Check root was autocompleted
	if id.Root != "755.3045256.2025923.103550" {
		t.Errorf("Root = %v, want %v", id.Root, "755.3045256.2025923.103550")
	}

	// Check extension was preserved
	if id.Extension != "myExtension" {
		t.Errorf("Extension = %v, want %v", id.Extension, "myExtension")
	}

	// Should not have errors
	if vctx.HasErrors() {
		t.Errorf("Unexpected validation error: %v", vctx.GetError())
	}
}

// TestHL7AEcg_Validate_AutocompleteIDs tests that all IDs in document are autocompleted
func TestHL7AEcg_Validate_AutocompleteIDs(t *testing.T) {
	// Setup singleton
	aecg := &HL7AEcg{
		ID: &ID{Root: "", Extension: "annotatedEcg"},
		Code: &Code[CPT_CODE, CodeSystemOID]{
			Code:       CPT_CODE_ECG_Routine,
			CodeSystem: CPT_OID,
		},
		EffectiveTime: &EffectiveTime{
			Low: Time{Value: "20231223120000"},
		},
		ComponentOf: &ComponentOfTimepointEvent{
			TimepointEvent: TimepointEvent{
				ComponentOf: ComponentOfSubjectAssignment{
					SubjectAssignment: SubjectAssignment{
						Subject: Subject{
							TrialSubject: TrialSubject{
								ID: &ID{Root: "", Extension: "subject"},
							},
						},
						ComponentOf: ComponentOfClinicalTrial{
							ClinicalTrial: ClinicalTrial{
								ID: ID{Root: "", Extension: "clinicalTrial"},
							},
						},
					},
				},
			},
		},
	}

	aecg.SetRootID("755.3045256.2025923.103550", "")

	ctx := context.Background()
	vctx := NewValidationContext(false)

	err := aecg.Validate(ctx, vctx)
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}

	if vctx.HasErrors() {
		t.Fatalf("Validation failed: %v", vctx.GetError())
	}

	// Check that all IDs were autocompleted
	if aecg.ID.Root != "755.3045256.2025923.103550" {
		t.Errorf("Document ID Root = %v, want %v", aecg.ID.Root, "755.3045256.2025923.103550")
	}

	subjectID := aecg.ComponentOf.TimepointEvent.ComponentOf.SubjectAssignment.Subject.TrialSubject.ID
	if subjectID.Root != "755.3045256.2025923.103550" {
		t.Errorf("Subject ID Root = %v, want %v", subjectID.Root, "755.3045256.2025923.103550")
	}

	clinicalTrialID := &aecg.ComponentOf.TimepointEvent.ComponentOf.SubjectAssignment.ComponentOf.ClinicalTrial.ID
	if clinicalTrialID.Root != "755.3045256.2025923.103550" {
		t.Errorf("ClinicalTrial ID Root = %v, want %v", clinicalTrialID.Root, "755.3045256.2025923.103550")
	}
}
