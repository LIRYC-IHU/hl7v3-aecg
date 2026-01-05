package hl7aecg

import (
	"context"
	"encoding/xml"
	"os"
	"strings"
	"testing"

	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
)

// TestCompleteWorkflow_MinimalECG tests creating a minimal valid aECG document
func TestCompleteWorkflow_MinimalECG(t *testing.T) {
	// Create temporary directory for output
	tmpDir := t.TempDir()

	// Create new aECG document
	h := NewHl7xml(tmpDir)
	if h == nil {
		t.Fatal("NewHl7xml() returned nil")
	}

	// Initialize with routine ECG code
	h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "")

	// Set global root ID for this test
	h.HL7AEcg.SetRootID("2.16.840.1.113883.3.1", "")

	// Set document ID (required)
	h.HL7AEcg.ID.SetID("", "TEST-DOC-001")

	// Set required codes (initialized by NewHl7xml, must be valid)
	h.HL7AEcg.ConfidentialityCode.SetCode(types.CONFIDENTIALITY_SPONSOR_BLINDED, "", "", "")
	h.HL7AEcg.ReasonCode.SetCode(types.REASON_PER_PROTOCOL, "", "", "")

	// Set basic metadata
	h.SetText("Test ECG Document").
		SetEffectiveTime("20231223120000", "20231223120010")

	// Set subject information (empty ID uses global root ID)
	h.SetSubject("", "SUBJ-001", types.SUBJECT_ROLE_ENROLLED)

	// Validate the document
	if err := h.Validate(); err != nil {
		t.Errorf("Validate() failed: %v", err)
	}

	// Verify required fields are set
	if h.HL7AEcg.Code == nil || h.HL7AEcg.Code.Code != types.CPT_CODE_ECG_Routine {
		t.Error("Code not set correctly")
	}
	if h.HL7AEcg.EffectiveTime == nil {
		t.Error("EffectiveTime not set")
	}
	if h.HL7AEcg.ComponentOf == nil {
		t.Error("Subject not set (ComponentOf is nil)")
	}
}

// TestCompleteWorkflow_FullECG tests creating a complete aECG document with all features
func TestCompleteWorkflow_FullECG(t *testing.T) {
	tmpDir := t.TempDir()

	h := NewHl7xml(tmpDir)
	h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "")

	// Set global root ID for this test
	h.HL7AEcg.SetRootID("2.16.840.1.113883.3.1", "")

	// Set document ID (required)
	h.HL7AEcg.ID.SetID("", "TEST-FULL-DOC-001")

	h.SetText("Complete 12-lead ECG for clinical trial").
		SetEffectiveTime("20231223120000.000", "20231223120010.000")

	// Set confidentiality and reason codes
	h.HL7AEcg.ConfidentialityCode = &types.Code[types.ConfidentialityCode, string]{}
	h.HL7AEcg.ConfidentialityCode.SetCode(types.CONFIDENTIALITY_SPONSOR_BLINDED, "", "Blinded to sponsor", "")

	h.HL7AEcg.ReasonCode = &types.Code[types.ReasonCode, string]{}
	h.HL7AEcg.ReasonCode.SetCode(types.REASON_PER_PROTOCOL, "", "Per protocol", "")

	// Set subject with demographics (empty ID uses global root ID)
	h.SetSubject("", "SUBJ-12345", types.SUBJECT_ROLE_ENROLLED).
		SetSubjectDemographics("JDO", "PAT-12345", types.GENDER_MALE, "19530508", types.RACE_WHITE)

	// Add rhythm series with multiple leads
	leads := map[types.LeadCode][]int{
		types.MDC_ECG_LEAD_I:  {0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		types.MDC_ECG_LEAD_II: {0, 2, 4, 6, 8, 10, 12, 14, 16, 18},
		types.MDC_ECG_LEAD_V1: {0, -1, -2, -3, -4, -5, -6, -7, -8, -9},
	}

	h.AddRhythmSeries("20231223120000.000", "20231223120010.000", 500.0, leads, 0.0, 5.0)

	// Add representative beat series
	repLeads := map[types.LeadCode][]int{
		types.MDC_ECG_LEAD_I:  {0, 1, 2, 3, 4},
		types.MDC_ECG_LEAD_II: {0, 2, 4, 6, 8},
	}
	h.AddRepresentativeBeatSeries("20231223120002.000", "20231223120003.000", 500.0, repLeads, 0.0, 5.0)

	// Validate
	if err := h.Validate(); err != nil {
		t.Errorf("Validate() failed for complete ECG: %v", err)
	}

	// Verify series were added
	if len(h.HL7AEcg.Component) != 2 {
		t.Errorf("Expected 2 components (series), got %d", len(h.HL7AEcg.Component))
	}

	// Verify rhythm series
	if len(h.HL7AEcg.Component) > 0 {
		rhythmSeries := h.HL7AEcg.Component[0].Series
		if rhythmSeries.Code == nil || rhythmSeries.Code.Code != types.RHYTHM_CODE {
			t.Error("First series should be RHYTHM")
		}
	}

	// Verify representative beat series
	if len(h.HL7AEcg.Component) > 1 {
		beatSeries := h.HL7AEcg.Component[1].Series
		if beatSeries.Code == nil || beatSeries.Code.Code != types.REPRESENTATIVE_BEAT_CODE {
			t.Error("Second series should be REPRESENTATIVE_BEAT")
		}
	}
}

// TestXMLMarshaling tests that the generated XML is valid
func TestXMLMarshaling(t *testing.T) {
	tmpDir := t.TempDir()

	h := NewHl7xml(tmpDir)
	h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "")
	h.HL7AEcg.SetRootID("2.16.840.1.113883.3.1", "")
	h.HL7AEcg.ID.SetID("", "TEST-XML-001")
	h.HL7AEcg.ConfidentialityCode.SetCode(types.CONFIDENTIALITY_SPONSOR_BLINDED, "", "", "")
	h.HL7AEcg.ReasonCode.SetCode(types.REASON_PER_PROTOCOL, "", "", "")
	h.SetText("XML Marshaling Test").
		SetEffectiveTime("20231223120000", "20231223120010").
		SetSubject("", "SUBJ-001", types.SUBJECT_ROLE_ENROLLED)

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(&h.HL7AEcg, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal XML: %v", err)
	}

	// Verify XML contains expected elements
	xmlString := string(xmlData)

	expectedElements := []string{
		"<AnnotatedECG",
		"xmlns=\"urn:hl7-org:v3\"",
		"<id",
		"<code",
		"<text>",
		"<effectiveTime>",
		"<subject>",
	}

	for _, element := range expectedElements {
		if !strings.Contains(xmlString, element) {
			t.Errorf("XML missing expected element: %s", element)
		}
	}
}

// TestXMLUnmarshaling tests reading back XML data
func TestXMLUnmarshaling(t *testing.T) {
	// Create a sample XML document
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<AnnotatedECG xmlns="urn:hl7-org:v3" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <id root="550e8400-e29b-41d4-a716-446655440000"/>
  <code code="93000" codeSystem="2.16.840.1.113883.6.12"/>
  <text>Test ECG</text>
  <effectiveTime>
    <low value="20231223120000"/>
    <high value="20231223120010"/>
  </effectiveTime>
</AnnotatedECG>`

	var aecg types.HL7AEcg
	err := xml.Unmarshal([]byte(xmlData), &aecg)
	if err != nil {
		t.Fatalf("Failed to unmarshal XML: %v", err)
	}

	// Verify unmarshaled data
	if aecg.ID == nil || aecg.ID.Root != "550e8400-e29b-41d4-a716-446655440000" {
		t.Error("ID not unmarshaled correctly")
	}
	if aecg.Code == nil || aecg.Code.Code != types.CPT_CODE_ECG_Routine {
		t.Error("Code not unmarshaled correctly")
	}
	if aecg.Text != "Test ECG" {
		t.Errorf("Text = %q, want %q", aecg.Text, "Test ECG")
	}
	if aecg.EffectiveTime == nil {
		t.Error("EffectiveTime not unmarshaled")
	}
}

// TestValidation_ErrorCases tests validation catches errors
func TestValidation_ErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*Hl7xml)
		wantError bool
	}{
		{
			name: "Missing code",
			setupFunc: func(h *Hl7xml) {
				h.HL7AEcg.Code = &types.Code[types.CPT_CODE, types.CodeSystemOID]{}
				h.HL7AEcg.EffectiveTime = &types.EffectiveTime{
					Low: types.Time{Value: "20231223120000"},
				}
				h.HL7AEcg.Subject = &types.TrialSubject{
					ID: &types.ID{Root: "2.16.840.1.113883.3.1234"},
				}
			},
			wantError: true,
		},
		{
			name: "Missing effective time",
			setupFunc: func(h *Hl7xml) {
				h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "")
				h.HL7AEcg.Subject = &types.TrialSubject{
					ID: &types.ID{Root: "2.16.840.1.113883.3.1234"},
				}
				h.HL7AEcg.EffectiveTime = nil
			},
			wantError: true,
		},
		{
			name: "Missing subject",
			setupFunc: func(h *Hl7xml) {
				h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "").
					SetEffectiveTime("20231223120000", "20231223120010")
				h.HL7AEcg.Subject = nil
			},
			wantError: true,
		},
		{
			name: "Valid complete document",
			setupFunc: func(h *Hl7xml) {
				h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "")
				h.HL7AEcg.SetRootID("2.16.840.1.113883.3.1", "")
				h.HL7AEcg.ID.SetID("", "VALID-DOC-001")
				h.HL7AEcg.ConfidentialityCode.SetCode(types.CONFIDENTIALITY_SPONSOR_BLINDED, "", "", "")
				h.HL7AEcg.ReasonCode.SetCode(types.REASON_PER_PROTOCOL, "", "", "")
				h.SetEffectiveTime("20231223120000", "20231223120010").
					SetSubject("", "SUBJ-001", types.SUBJECT_ROLE_ENROLLED)
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			h := NewHl7xml(tmpDir)
			tt.setupFunc(h)

			err := h.Validate()
			if (err != nil) != tt.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestFileOutput tests writing to file system
func TestFileOutput(t *testing.T) {
	tmpDir := t.TempDir()

	h := NewHl7xml(tmpDir)
	h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "")
	h.HL7AEcg.SetRootID("2.16.840.1.113883.3.1", "")
	h.HL7AEcg.ID.SetID("", "TEST-FILE-001")
	h.HL7AEcg.ConfidentialityCode.SetCode(types.CONFIDENTIALITY_SPONSOR_BLINDED, "", "", "")
	h.HL7AEcg.ReasonCode.SetCode(types.REASON_PER_PROTOCOL, "", "", "")
	h.SetText("File Output Test").
		SetEffectiveTime("20231223120000", "20231223120010").
		SetSubject("", "SUBJ-001", types.SUBJECT_ROLE_ENROLLED)

	// Generate XML file (using Test() method which writes to /tmp/hl7aecg_example.xml)
	_, err := h.Test()
	if err != nil {
		t.Fatalf("Test() failed: %v", err)
	}

	// The Test() method writes to a fixed location
	filename := "/tmp/hl7aecg_example.xml"

	// Verify file was created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("Output file was not created: %s", filename)
	}

	// Read and verify file contents
	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	xmlString := string(data)
	if !strings.Contains(xmlString, "<AnnotatedECG") {
		t.Error("Output file does not contain AnnotatedECG element")
	}
	if !strings.Contains(xmlString, "File Output Test") {
		t.Error("Output file does not contain expected text")
	}

	// Clean up
	os.Remove(filename)
}

// TestFluentAPI_Chaining tests method chaining works correctly
func TestFluentAPI_Chaining(t *testing.T) {
	tmpDir := t.TempDir()

	// All methods should return *Hl7xml for chaining
	result := NewHl7xml(tmpDir).
		Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "")

	result.HL7AEcg.SetRootID("2.16.840.1.113883.3.1", "")
	result.HL7AEcg.ID.SetID("", "TEST-FLUENT-001")
	result.HL7AEcg.ConfidentialityCode.SetCode(types.CONFIDENTIALITY_SPONSOR_BLINDED, "", "", "")
	result.HL7AEcg.ReasonCode.SetCode(types.REASON_PER_PROTOCOL, "", "", "")

	result.SetText("Fluent API Test").
		SetEffectiveTime("20231223120000", "20231223120010").
		SetSubject("", "SUBJ-001", types.SUBJECT_ROLE_ENROLLED).
		SetSubjectDemographics("JDO", "PAT-001", types.GENDER_MALE, "19530508", types.RACE_WHITE)

	if result == nil {
		t.Fatal("Fluent API chain returned nil")
	}

	// Verify all operations were applied
	if result.HL7AEcg.Text != "Fluent API Test" {
		t.Error("SetText did not work in chain")
	}
	if result.HL7AEcg.EffectiveTime == nil {
		t.Error("SetEffectiveTime did not work in chain")
	}
	if result.HL7AEcg.ComponentOf == nil {
		t.Error("SetSubject did not work in chain")
	}
	if result.HL7AEcg.ComponentOf != nil &&
		result.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.SubjectAssignment.Subject.TrialSubject.SubjectDemographicPerson == nil {
		t.Error("SetSubjectDemographics did not work in chain")
	}
}

// TestContextCancellation tests that context cancellation works
func TestContextCancellation(t *testing.T) {
	tmpDir := t.TempDir()

	h := NewHl7xml(tmpDir)
	h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "")
	h.HL7AEcg.SetRootID("2.16.840.1.113883.3.1", "")
	h.HL7AEcg.ID.SetID("", "TEST-CTX-001")
	h.HL7AEcg.ConfidentialityCode.SetCode(types.CONFIDENTIALITY_SPONSOR_BLINDED, "", "", "")
	h.HL7AEcg.ReasonCode.SetCode(types.REASON_PER_PROTOCOL, "", "", "")
	h.SetEffectiveTime("20231223120000", "20231223120010").
		SetSubject("", "SUBJ-001", types.SUBJECT_ROLE_ENROLLED)

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	vctx := types.NewValidationContext(false)
	err := h.HL7AEcg.Validate(ctx, vctx)

	// Validation should return context.Canceled error
	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}

// TestMultipleSeries tests adding multiple ECG series
func TestMultipleSeries(t *testing.T) {
	tmpDir := t.TempDir()

	h := NewHl7xml(tmpDir)
	h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "")
	h.HL7AEcg.SetRootID("2.16.840.1.113883.3.1", "")
	h.HL7AEcg.ID.SetID("", "TEST-MULTI-001")
	h.HL7AEcg.ConfidentialityCode.SetCode(types.CONFIDENTIALITY_SPONSOR_BLINDED, "", "", "")
	h.HL7AEcg.ReasonCode.SetCode(types.REASON_PER_PROTOCOL, "", "", "")
	h.SetEffectiveTime("20231223120000", "20231223120100").
		SetSubject("", "SUBJ-001", types.SUBJECT_ROLE_ENROLLED)

	// Add multiple rhythm series
	leads1 := map[types.LeadCode][]int{
		types.MDC_ECG_LEAD_I: {1, 2, 3, 4, 5},
	}
	leads2 := map[types.LeadCode][]int{
		types.MDC_ECG_LEAD_II: {6, 7, 8, 9, 10},
	}
	leads3 := map[types.LeadCode][]int{
		types.MDC_ECG_LEAD_V1: {11, 12, 13, 14, 15},
	}

	h.AddRhythmSeries("20231223120000.000", "20231223120010.000", 500.0, leads1, 0.0, 5.0).
		AddRhythmSeries("20231223120030.000", "20231223120040.000", 500.0, leads2, 0.0, 5.0).
		AddRepresentativeBeatSeries("20231223120050.000", "20231223120051.000", 500.0, leads3, 0.0, 5.0)

	// Should have 3 series
	if len(h.HL7AEcg.Component) != 3 {
		t.Errorf("Expected 3 series, got %d", len(h.HL7AEcg.Component))
	}

	// Verify series types
	if h.HL7AEcg.Component[0].Series.Code.Code != types.RHYTHM_CODE {
		t.Error("First series should be RHYTHM")
	}
	if h.HL7AEcg.Component[1].Series.Code.Code != types.RHYTHM_CODE {
		t.Error("Second series should be RHYTHM")
	}
	if h.HL7AEcg.Component[2].Series.Code.Code != types.REPRESENTATIVE_BEAT_CODE {
		t.Error("Third series should be REPRESENTATIVE_BEAT")
	}
}
