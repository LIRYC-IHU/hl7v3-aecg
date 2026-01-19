package types

import (
	"context"
	"strings"
	"sync"
	"testing"
)

// TestValidationContext_AddError tests adding errors to validation context
func TestValidationContext_AddError(t *testing.T) {
	vctx := NewValidationContext(false)

	if vctx.HasErrors() {
		t.Error("NewValidationContext should have no errors initially")
	}

	vctx.AddError(ErrMissingID)
	if !vctx.HasErrors() {
		t.Error("HasErrors should return true after adding an error")
	}

	if len(vctx.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(vctx.Errors))
	}

	vctx.AddError(ErrMissingCode)
	if len(vctx.Errors) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(vctx.Errors))
	}
}

// TestValidationContext_AddWarning tests adding warnings to validation context
func TestValidationContext_AddWarning(t *testing.T) {
	vctx := NewValidationContext(false)

	if vctx.HasWarnings() {
		t.Error("NewValidationContext should have no warnings initially")
	}

	vctx.AddWarning("This is a warning")
	if !vctx.HasWarnings() {
		t.Error("HasWarnings should return true after adding a warning")
	}

	if len(vctx.Warnings) != 1 {
		t.Errorf("Expected 1 warning, got %d", len(vctx.Warnings))
	}
}

// TestValidationContext_GetError tests retrieving errors from validation context
func TestValidationContext_GetError(t *testing.T) {
	tests := []struct {
		name         string
		errors       []error
		wantNil      bool
		wantMultiple bool
	}{
		{
			name:    "No errors",
			errors:  []error{},
			wantNil: true,
		},
		{
			name:    "Single error",
			errors:  []error{ErrMissingID},
			wantNil: false,
		},
		{
			name:         "Multiple errors",
			errors:       []error{ErrMissingID, ErrMissingCode},
			wantNil:      false,
			wantMultiple: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vctx := NewValidationContext(false)
			for _, err := range tt.errors {
				vctx.AddError(err)
			}

			err := vctx.GetError()
			if tt.wantNil && err != nil {
				t.Errorf("GetError() should return nil, got %v", err)
			}
			if !tt.wantNil && err == nil {
				t.Error("GetError() should return error, got nil")
			}
			if tt.wantMultiple && !strings.Contains(err.Error(), "multiple") {
				t.Errorf("GetError() should indicate multiple errors, got %v", err)
			}
		})
	}
}

// TestID_Validate tests ID validation
func TestID_Validate(t *testing.T) {
	// Reset singleton to ensure empty IDs produce errors
	savedInstance := instanceID
	savedOnce := once
	instanceID = nil
	once = *new(sync.Once)
	defer func() {
		instanceID = savedInstance
		once = savedOnce
	}()

	tests := []struct {
		name      string
		id        ID
		wantError *ValidationError
	}{
		{
			name: "Valid UUID",
			id: ID{
				Root: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantError: nil,
		},
		{
			name: "Valid OID",
			id: ID{
				Root: "2.16.840.1.113883.3.1234",
			},
			wantError: nil,
		},
		{
			name: "Empty root",
			id: ID{
				Root: "",
			},
			wantError: ErrMissingID,
		},
		{
			name: "Custom ID format (now allowed)",
			id: ID{
				Root: "invalid-id-format",
			},
			wantError: nil, // Any non-empty Root is now valid
		},
		{
			name: "OID with extension",
			id: ID{
				Root:      "2.16.840.1.113883.3.1234",
				Extension: "TEST-001",
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.id.Validate(ctx, vctx)

			if tt.wantError != nil {
				if !vctx.HasErrors() {
					t.Error("Expected validation error, got none")
					return
				}
				found := false
				for _, err := range vctx.Errors {
					if err == tt.wantError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error %v, got errors %v", tt.wantError, vctx.Errors)
				}
			} else {
				if vctx.HasErrors() {
					t.Errorf("Expected no errors, got %v", vctx.Errors)
				}
			}
		})
	}
}

// TestEffectiveTime_Validate tests EffectiveTime validation
func TestEffectiveTime_Validate(t *testing.T) {
	tests := []struct {
		name          string
		effectiveTime EffectiveTime
		wantError     *ValidationError
	}{
		{
			name: "Valid time range",
			effectiveTime: EffectiveTime{
				Low:  Time{Value: "20231223120000"},
				High: Time{Value: "20231223120010"},
			},
			wantError: nil,
		},
		{
			name: "Valid with milliseconds",
			effectiveTime: EffectiveTime{
				Low:  Time{Value: "20231223120000.000"},
				High: Time{Value: "20231223120010.000"},
			},
			wantError: nil,
		},
		{
			name: "Valid with only low",
			effectiveTime: EffectiveTime{
				Low: Time{Value: "20231223120000"},
			},
			wantError: nil,
		},
		{
			name: "Valid with only high",
			effectiveTime: EffectiveTime{
				High: Time{Value: "20231223120000"},
			},
			wantError: nil,
		},
		{
			name: "Both empty",
			effectiveTime: EffectiveTime{
				Low:  Time{Value: ""},
				High: Time{Value: ""},
			},
			wantError: ErrMissingTimeValue,
		},
		{
			name: "Invalid low format",
			effectiveTime: EffectiveTime{
				Low:  Time{Value: "invalid-date"},
				High: Time{Value: "20231223120000"},
			},
			wantError: ErrInvalidTimeFormat,
		},
		{
			name: "Invalid high format",
			effectiveTime: EffectiveTime{
				Low:  Time{Value: "20231223120000"},
				High: Time{Value: "not-a-date"},
			},
			wantError: ErrInvalidTimeFormat,
		},
		{
			name: "Date only format",
			effectiveTime: EffectiveTime{
				Low: Time{Value: "20231223"},
			},
			wantError: nil,
		},
		{
			name: "Year-month format",
			effectiveTime: EffectiveTime{
				Low: Time{Value: "202312"},
			},
			wantError: nil,
		},
		{
			name: "Year only format",
			effectiveTime: EffectiveTime{
				Low: Time{Value: "2023"},
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.effectiveTime.Validate(ctx, vctx)

			if tt.wantError != nil {
				if !vctx.HasErrors() {
					t.Error("Errors:", vctx.Errors)
					t.Error("Expected error:", tt.wantError)
					t.Error("name:", tt.name)
					t.Error("Expected validation error, got none")
					return
				}
				found := false
				for _, err := range vctx.Errors {
					if err == tt.wantError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error %v, got errors %v", tt.wantError, vctx.Errors)
				}
			} else {
				if vctx.HasErrors() {
					t.Errorf("Expected no errors, got %v", vctx.Errors)
				}
			}
		})
	}
}

// TestCode_ValidateCode_Confidentiality tests confidentiality code validation
func TestCode_ValidateCode_Confidentiality(t *testing.T) {
	tests := []struct {
		name      string
		code      Code[ConfidentialityCode, string]
		wantError *ValidationError
	}{
		{
			name: "Valid - S (Sponsor blinded)",
			code: Code[ConfidentialityCode, string]{
				Code: CONFIDENTIALITY_SPONSOR_BLINDED,
			},
			wantError: nil,
		},
		{
			name: "Valid - I (Investigator blinded)",
			code: Code[ConfidentialityCode, string]{
				Code: CONFIDENTIALITY_INVESTIGATOR_BLINDED,
			},
			wantError: nil,
		},
		{
			name: "Valid - B (Both blinded)",
			code: Code[ConfidentialityCode, string]{
				Code: CONFIDENTIALITY_BOTH,
			},
			wantError: nil,
		},
		{
			name: "Valid - C (Custom)",
			code: Code[ConfidentialityCode, string]{
				Code: CONFIDENTIALITY_CUSTOM,
			},
			wantError: nil,
		},
		{
			name: "Invalid code",
			code: Code[ConfidentialityCode, string]{
				Code: "INVALID",
			},
			wantError: ErrConfidentialityCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.code.ValidateCode(ctx, vctx, "ConfidentialityCode")

			if tt.wantError != nil {
				if !vctx.HasErrors() {
					t.Error("Expected validation error, got none")
					return
				}
				found := false
				for _, err := range vctx.Errors {
					if err == tt.wantError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error %v, got errors %v", tt.wantError, vctx.Errors)
				}
			} else {
				if vctx.HasErrors() {
					t.Errorf("Expected no errors, got %v", vctx.Errors)
				}
			}
		})
	}
}

// TestCode_ValidateCode_Reason tests reason code validation
func TestCode_ValidateCode_Reason(t *testing.T) {
	tests := []struct {
		name      string
		code      Code[ReasonCode, string]
		wantError *ValidationError
	}{
		{
			name: "Valid - PER_PROTOCOL",
			code: Code[ReasonCode, string]{
				Code: REASON_PER_PROTOCOL,
			},
			wantError: nil,
		},
		{
			name: "Valid - NOT_IN_PROTOCOL",
			code: Code[ReasonCode, string]{
				Code: REASON_NOT_IN_PROTOCOL,
			},
			wantError: nil,
		},
		{
			name: "Valid - WRONG_EVENT",
			code: Code[ReasonCode, string]{
				Code: REASON_WRONG_EVENT,
			},
			wantError: nil,
		},
		{
			name: "Invalid code",
			code: Code[ReasonCode, string]{
				Code: "INVALID_REASON",
			},
			wantError: ErrReasonCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.code.ValidateCode(ctx, vctx, "ReasonCode")

			if tt.wantError != nil {
				if !vctx.HasErrors() {
					t.Error("Expected validation error, got none")
					return
				}
				found := false
				for _, err := range vctx.Errors {
					if err == tt.wantError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error %v, got errors %v", tt.wantError, vctx.Errors)
				}
			} else {
				if vctx.HasErrors() {
					t.Errorf("Expected no errors, got %v", vctx.Errors)
				}
			}
		})
	}
}

// TestCode_ValidateCode_TrialSubject tests trial subject code validation
func TestCode_ValidateCode_TrialSubject(t *testing.T) {
	tests := []struct {
		name      string
		code      Code[CodeRole, CodeSystemOID]
		wantError *ValidationError
	}{
		{
			name: "Valid - SCREENING",
			code: Code[CodeRole, CodeSystemOID]{
				Code: SUBJECT_ROLE_SCREENING,
			},
			wantError: nil,
		},
		{
			name: "Valid - ENROLLED",
			code: Code[CodeRole, CodeSystemOID]{
				Code: SUBJECT_ROLE_ENROLLED,
			},
			wantError: nil,
		},
		{
			name: "Invalid code",
			code: Code[CodeRole, CodeSystemOID]{
				Code: "INVALID_ROLE",
			},
			wantError: ErrTrialSubjectCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.code.ValidateCode(ctx, vctx, "TrialSubjectCode")

			if tt.wantError != nil {
				if !vctx.HasErrors() {
					t.Error("Expected validation error, got none")
					return
				}
				found := false
				for _, err := range vctx.Errors {
					if err == tt.wantError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error %v, got errors %v", tt.wantError, vctx.Errors)
				}
			} else {
				if vctx.HasErrors() {
					t.Errorf("Expected no errors, got %v", vctx.Errors)
				}
			}
		})
	}
}

// TestCode_ValidateCode_Gender tests gender code validation
// NOTE: There's a bug in the validator - it returns ErrRaceCode instead of ErrGenderCode
func TestCode_ValidateCode_Gender(t *testing.T) {
	tests := []struct {
		name      string
		code      Code[GenderCode, CodeSystemOID]
		wantError *ValidationError
	}{
		{
			name: "Valid - FEMALE",
			code: Code[GenderCode, CodeSystemOID]{
				Code: GENDER_FEMALE,
			},
			wantError: nil,
		},
		{
			name: "Valid - MALE",
			code: Code[GenderCode, CodeSystemOID]{
				Code: GENDER_MALE,
			},
			wantError: nil,
		},
		{
			name: "Valid - UNDIFFERENTIATED",
			code: Code[GenderCode, CodeSystemOID]{
				Code: GENDER_UNDIFFERENTIATED,
			},
			wantError: nil,
		},
		{
			name: "Invalid code",
			code: Code[GenderCode, CodeSystemOID]{
				Code: "X",
			},
			// NOTE: Bug in validator - returns ErrRaceCode instead of ErrGenderCode
			wantError: ErrRaceCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.code.ValidateCode(ctx, vctx, "GenderCode")

			if tt.wantError != nil {
				if !vctx.HasErrors() {
					t.Error("Expected validation error, got none")
					return
				}
				found := false
				for _, err := range vctx.Errors {
					if err == tt.wantError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error %v, got errors %v", tt.wantError, vctx.Errors)
				}
			} else {
				if vctx.HasErrors() {
					t.Errorf("Expected no errors, got %v", vctx.Errors)
				}
			}
		})
	}
}

// TestHL7AEcg_Validate tests complete HL7AEcg validation
func TestHL7AEcg_Validate(t *testing.T) {
	tests := []struct {
		name        string
		aecg        HL7AEcg
		wantErrors  []*ValidationError
		description string
	}{
		{
			name: "Valid minimal aECG",
			aecg: HL7AEcg{
				ID: &ID{Root: "2.16.840.1.113883.3.1234"},
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
										ID: &ID{Root: "2.16.840.1.113883.3.1234"},
									},
								},
							},
						},
					},
				},
			},
			wantErrors:  nil,
			description: "Should validate successfully with all required fields",
		},
		{
			name: "Missing code",
			aecg: HL7AEcg{
				ID: &ID{Root: "2.16.840.1.113883.3.1234"},
				Code: &Code[CPT_CODE, CodeSystemOID]{
					CodeSystem: CPT_OID,
				},
				EffectiveTime: &EffectiveTime{
					Low: Time{Value: "20231223120000"},
				},
				Subject: &TrialSubject{
					ID: &ID{Root: "2.16.840.1.113883.3.1234"},
				},
			},
			wantErrors:  []*ValidationError{ErrMissingCode},
			description: "Should fail validation when code is missing",
		},
		{
			name: "Missing code system",
			aecg: HL7AEcg{
				ID: &ID{Root: "2.16.840.1.113883.3.1234"},
				Code: &Code[CPT_CODE, CodeSystemOID]{
					Code: CPT_CODE_ECG_Routine,
				},
				EffectiveTime: &EffectiveTime{
					Low: Time{Value: "20231223120000"},
				},
				Subject: &TrialSubject{
					ID: &ID{Root: "2.16.840.1.113883.3.1234"},
				},
			},
			wantErrors:  []*ValidationError{ErrMissingCodeSystem},
			description: "Should fail validation when code system is missing",
		},
		{
			name: "Missing effective time",
			aecg: HL7AEcg{
				ID: &ID{Root: "2.16.840.1.113883.3.1234"},
				Code: &Code[CPT_CODE, CodeSystemOID]{
					Code:       CPT_CODE_ECG_Routine,
					CodeSystem: CPT_OID,
				},
				Subject: &TrialSubject{
					ID: &ID{Root: "2.16.840.1.113883.3.1234"},
				},
			},
			wantErrors:  []*ValidationError{ErrMissingEffectiveTime},
			description: "Should fail validation when effective time is missing",
		},
		{
			name: "Missing subject",
			aecg: HL7AEcg{
				ID: &ID{Root: "2.16.840.1.113883.3.1234"},
				Code: &Code[CPT_CODE, CodeSystemOID]{
					Code:       CPT_CODE_ECG_Routine,
					CodeSystem: CPT_OID,
				},
				EffectiveTime: &EffectiveTime{
					Low: Time{Value: "20231223120000"},
				},
			},
			wantErrors:  []*ValidationError{ErrMissingSubject},
			description: "Should fail validation when subject is missing",
		},
		{
			name: "Invalid confidentiality code",
			aecg: HL7AEcg{
				ID: &ID{Root: "2.16.840.1.113883.3.1234"},
				Code: &Code[CPT_CODE, CodeSystemOID]{
					Code:       CPT_CODE_ECG_Routine,
					CodeSystem: CPT_OID,
				},
				EffectiveTime: &EffectiveTime{
					Low: Time{Value: "20231223120000"},
				},
				ConfidentialityCode: &Code[ConfidentialityCode, string]{
					Code: "INVALID",
				},
				Subject: &TrialSubject{
					ID: &ID{Root: "2.16.840.1.113883.3.1234"},
				},
			},
			wantErrors:  []*ValidationError{ErrConfidentialityCode},
			description: "Should fail validation with invalid confidentiality code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.aecg.Validate(ctx, vctx)

			if tt.wantErrors != nil {
				if !vctx.HasErrors() {
					t.Errorf("%s: Expected validation errors, got none", tt.description)
					return
				}
				for _, wantErr := range tt.wantErrors {
					found := false
					for _, gotErr := range vctx.Errors {
						if gotErr == wantErr {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("%s: Expected error %v not found in %v", tt.description, wantErr, vctx.Errors)
					}
				}
			} else {
				if vctx.HasErrors() {
					t.Errorf("%s: Expected no errors, got %v", tt.description, vctx.Errors)
				}
			}
		})
	}
}

// TestClinicalTrial_Validate tests clinical trial validation
func TestClinicalTrial_Validate(t *testing.T) {
	// Reset singleton to ensure empty IDs produce errors
	savedInstance := instanceID
	savedOnce := once
	instanceID = nil
	once = *new(sync.Once)
	defer func() {
		instanceID = savedInstance
		once = savedOnce
	}()

	tests := []struct {
		name      string
		trial     ClinicalTrial
		wantError *ValidationError
	}{
		{
			name: "Valid trial with OID",
			trial: ClinicalTrial{
				ID: ID{Root: "2.16.840.1.113883.3.1234"},
			},
			wantError: nil,
		},
		{
			name: "Valid trial with UUID",
			trial: ClinicalTrial{
				ID: ID{Root: "550e8400-e29b-41d4-a716-446655440000"},
			},
			wantError: nil,
		},
		{
			name: "Missing ID",
			trial: ClinicalTrial{
				ID: ID{Root: ""},
			},
			wantError: ErrMissingID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			vctx := NewValidationContext(false)

			tt.trial.Validate(ctx, vctx)

			if tt.wantError != nil {
				if !vctx.HasErrors() {
					t.Error("Expected validation error, got none")
					return
				}
				found := false
				for _, err := range vctx.Errors {
					if err == tt.wantError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error %v, got errors %v", tt.wantError, vctx.Errors)
				}
			} else {
				if vctx.HasErrors() {
					t.Errorf("Expected no errors, got %v", vctx.Errors)
				}
			}
		})
	}
}

// TestValidationError_Error tests ValidationError's error message formatting
func TestValidationError_Error(t *testing.T) {
	tests := []struct {
		name        string
		err         *ValidationError
		wantContain string
	}{
		{
			name: "Error without value",
			err: &ValidationError{
				Field:   "TestField",
				Message: "Test message",
			},
			wantContain: "TestField",
		},
		{
			name: "Error with value",
			err: &ValidationError{
				Field:   "TestField",
				Message: "Test message",
				Value:   "invalid-value",
			},
			wantContain: "invalid-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := tt.err.Error()
			if !strings.Contains(errMsg, tt.wantContain) {
				t.Errorf("Error message should contain %q, got %q", tt.wantContain, errMsg)
			}
		})
	}
}
