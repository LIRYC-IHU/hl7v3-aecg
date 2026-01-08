package types

import (
	"testing"
)

// TestSeries_SetSeriesCode tests the SetSeriesCode method.
func TestSeries_SetSeriesCode(t *testing.T) {
	tests := []struct {
		name           string
		code           SeriesTypeCode
		codeSystem     CodeSystemOID
		codeSystemName string
		displayName    string
		wantCode       SeriesTypeCode
		wantSystem     CodeSystemOID
		wantSysName    string
		wantDisplay    string
	}{
		{
			name:           "Set RHYTHM code with all attributes",
			code:           RHYTHM_CODE,
			codeSystem:     HL7_ActCode_OID,
			codeSystemName: "ActCode",
			displayName:    "Rhythm Waveforms",
			wantCode:       RHYTHM_CODE,
			wantSystem:     HL7_ActCode_OID,
			wantSysName:    "ActCode",
			wantDisplay:    "Rhythm Waveforms",
		},
		{
			name:           "Set REPRESENTATIVE_BEAT code with all attributes",
			code:           REPRESENTATIVE_BEAT_CODE,
			codeSystem:     HL7_ActCode_OID,
			codeSystemName: "ActCode",
			displayName:    "Representative Beat Waveforms",
			wantCode:       REPRESENTATIVE_BEAT_CODE,
			wantSystem:     HL7_ActCode_OID,
			wantSysName:    "ActCode",
			wantDisplay:    "Representative Beat Waveforms",
		},
		{
			name:           "Set code with empty codeSystemName",
			code:           RHYTHM_CODE,
			codeSystem:     HL7_ActCode_OID,
			codeSystemName: "",
			displayName:    "Rhythm Waveforms",
			wantCode:       RHYTHM_CODE,
			wantSystem:     HL7_ActCode_OID,
			wantSysName:    "",
			wantDisplay:    "Rhythm Waveforms",
		},
		{
			name:           "Set code with empty displayName",
			code:           RHYTHM_CODE,
			codeSystem:     HL7_ActCode_OID,
			codeSystemName: "ActCode",
			displayName:    "",
			wantCode:       RHYTHM_CODE,
			wantSystem:     HL7_ActCode_OID,
			wantSysName:    "ActCode",
			wantDisplay:    "",
		},
		{
			name:           "Set code with both names empty",
			code:           RHYTHM_CODE,
			codeSystem:     HL7_ActCode_OID,
			codeSystemName: "",
			displayName:    "",
			wantCode:       RHYTHM_CODE,
			wantSystem:     HL7_ActCode_OID,
			wantSysName:    "",
			wantDisplay:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			series := &Series{}

			// Call SetSeriesCode
			result := series.SetSeriesCode(tt.code, tt.codeSystem, tt.codeSystemName, tt.displayName)

			// Verify it returns the series pointer for chaining
			if result != series {
				t.Error("SetSeriesCode should return the same Series pointer for method chaining")
			}

			// Verify the code was set correctly
			if series.Code == nil {
				t.Fatal("SetSeriesCode should set the Code field")
			}

			if series.Code.Code != tt.wantCode {
				t.Errorf("Code.Code = %v, want %v", series.Code.Code, tt.wantCode)
			}

			if series.Code.CodeSystem != tt.wantSystem {
				t.Errorf("Code.CodeSystem = %v, want %v", series.Code.CodeSystem, tt.wantSystem)
			}

			if series.Code.CodeSystemName != tt.wantSysName {
				t.Errorf("Code.CodeSystemName = %v, want %v", series.Code.CodeSystemName, tt.wantSysName)
			}

			if series.Code.DisplayName != tt.wantDisplay {
				t.Errorf("Code.DisplayName = %v, want %v", series.Code.DisplayName, tt.wantDisplay)
			}
		})
	}
}

// TestSeries_SetSeriesCode_OverwritesExisting tests that SetSeriesCode
// properly overwrites an existing code.
func TestSeries_SetSeriesCode_OverwritesExisting(t *testing.T) {
	series := &Series{
		Code: &Code[SeriesTypeCode, CodeSystemOID]{
			Code:           RHYTHM_CODE,
			CodeSystem:     HL7_ActCode_OID,
			CodeSystemName: "",
			DisplayName:    "",
		},
	}

	// Update the code with additional attributes
	series.SetSeriesCode(
		RHYTHM_CODE,
		HL7_ActCode_OID,
		"ActCode",
		"Rhythm Waveforms",
	)

	if series.Code.CodeSystemName != "ActCode" {
		t.Errorf("Code.CodeSystemName = %v, want %v", series.Code.CodeSystemName, "ActCode")
	}

	if series.Code.DisplayName != "Rhythm Waveforms" {
		t.Errorf("Code.DisplayName = %v, want %v", series.Code.DisplayName, "Rhythm Waveforms")
	}
}

// TestSeries_SetSeriesCode_MethodChaining tests that SetSeriesCode can be
// chained with other Series methods.
func TestSeries_SetSeriesCode_MethodChaining(t *testing.T) {
	series := &Series{}

	// Test method chaining
	result := series.
		SetSeriesCode(RHYTHM_CODE, HL7_ActCode_OID, "ActCode", "Rhythm Waveforms")

	if result != series {
		t.Error("Method chaining should return the same Series pointer")
	}

	// Verify all fields were set
	if series.Code == nil {
		t.Fatal("Code should be set")
	}

	if series.Code.Code != RHYTHM_CODE {
		t.Errorf("Code.Code = %v, want %v", series.Code.Code, RHYTHM_CODE)
	}

	if series.Code.CodeSystemName != "ActCode" {
		t.Errorf("Code.CodeSystemName = %v, want %v", series.Code.CodeSystemName, "ActCode")
	}

	if series.Code.DisplayName != "Rhythm Waveforms" {
		t.Errorf("Code.DisplayName = %v, want %v", series.Code.DisplayName, "Rhythm Waveforms")
	}
}
