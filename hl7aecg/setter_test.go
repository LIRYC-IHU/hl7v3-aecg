package hl7aecg

import (
	"testing"

	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
)

// TestHl7xml_SetSeriesCode tests the SetSeriesCode method.
func TestHl7xml_SetSeriesCode(t *testing.T) {
	t.Run("Success - updates last series code", func(t *testing.T) {
		h := NewHl7xml("/tmp")
		h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "CPT-4", "")

		// Add a rhythm series
		leads := map[types.LeadCode][]int{
			types.MDC_ECG_LEAD_I: {1, 2, 3, 4, 5},
		}
		h.AddRhythmSeries("20021122091000", "20021122091010", nil, nil, 500, leads, 0, 5)

		// Update the series code
		err := h.SetSeriesCode(
			types.RHYTHM_CODE,
			types.HL7_ActCode_OID,
			"ActCode",
			"Rhythm Waveforms",
		)

		if err != nil {
			t.Fatalf("SetSeriesCode returned error: %v", err)
		}

		// Verify the code was updated
		if len(h.HL7AEcg.Component) == 0 {
			t.Fatal("No components found after adding series")
		}

		series := &h.HL7AEcg.Component[0].Series
		if series.Code == nil {
			t.Fatal("Series code is nil")
		}

		if series.Code.Code != types.RHYTHM_CODE {
			t.Errorf("Code = %v, want %v", series.Code.Code, types.RHYTHM_CODE)
		}

		if series.Code.CodeSystem != types.HL7_ActCode_OID {
			t.Errorf("CodeSystem = %v, want %v", series.Code.CodeSystem, types.HL7_ActCode_OID)
		}

		if series.Code.CodeSystemName != "ActCode" {
			t.Errorf("CodeSystemName = %v, want %v", series.Code.CodeSystemName, "ActCode")
		}

		if series.Code.DisplayName != "Rhythm Waveforms" {
			t.Errorf("DisplayName = %v, want %v", series.Code.DisplayName, "Rhythm Waveforms")
		}
	})

	t.Run("Error - no series in Component array", func(t *testing.T) {
		h := NewHl7xml("/tmp")
		h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "CPT-4", "")

		// Try to set series code without adding any series
		err := h.SetSeriesCode(
			types.RHYTHM_CODE,
			types.HL7_ActCode_OID,
			"ActCode",
			"Rhythm Waveforms",
		)

		if err == nil {
			t.Fatal("SetSeriesCode should return error when no series exists")
		}

		expectedMsg := "no series found: Component array is empty"
		if err.Error() != expectedMsg {
			t.Errorf("Error message = %v, want %v", err.Error(), expectedMsg)
		}
	})

	t.Run("Success - updates last of multiple series", func(t *testing.T) {
		h := NewHl7xml("/tmp")
		h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "CPT-4", "")

		// Add two series
		leads := map[types.LeadCode][]int{
			types.MDC_ECG_LEAD_I: {1, 2, 3, 4, 5},
		}
		h.AddRhythmSeries("20021122091000", "20021122091010", nil, nil, 500, leads, 0, 5)
		h.AddRepresentativeBeatSeries("20021122091000", "20021122091010", 500, leads, 0, 5)

		// Update the last series code (representative beat)
		err := h.SetSeriesCode(
			types.REPRESENTATIVE_BEAT_CODE,
			types.HL7_ActCode_OID,
			"ActCode",
			"Representative Beat Waveforms",
		)

		if err != nil {
			t.Fatalf("SetSeriesCode returned error: %v", err)
		}

		// Verify the last series was updated
		if len(h.HL7AEcg.Component) != 2 {
			t.Fatalf("Expected 2 components, got %d", len(h.HL7AEcg.Component))
		}

		lastSeries := &h.HL7AEcg.Component[1].Series
		if lastSeries.Code.DisplayName != "Representative Beat Waveforms" {
			t.Errorf("Last series DisplayName = %v, want %v",
				lastSeries.Code.DisplayName, "Representative Beat Waveforms")
		}

		// Verify the first series was not changed
		firstSeries := &h.HL7AEcg.Component[0].Series
		if firstSeries.Code.DisplayName != "" {
			t.Errorf("First series DisplayName should be empty, got %v",
				firstSeries.Code.DisplayName)
		}
	})
}
