package hl7aecg

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
)

// TestHl7xml_Unmarshal tests the Unmarshal method on Hl7xml
func TestHl7xml_Unmarshal(t *testing.T) {
	tests := []struct {
		name      string
		xmlData   string
		wantError bool
		check     func(t *testing.T, h *Hl7xml)
	}{
		{
			name: "Valid minimal aECG document",
			xmlData: `<?xml version="1.0" encoding="UTF-8"?>
<AnnotatedECG xmlns="urn:hl7-org:v3">
	<id root="728989ec-b8bc-49cd-9a5a-30be5ade1db5"/>
	<code code="93000" codeSystem="2.16.840.1.113883.6.12"/>
	<text>Test ECG</text>
</AnnotatedECG>`,
			wantError: false,
			check: func(t *testing.T, h *Hl7xml) {
				if h.HL7AEcg.ID == nil || h.HL7AEcg.ID.Root != "728989ec-b8bc-49cd-9a5a-30be5ade1db5" {
					t.Errorf("ID.Root = %v, want %v", h.HL7AEcg.ID.Root, "728989ec-b8bc-49cd-9a5a-30be5ade1db5")
				}
				if h.HL7AEcg.Code == nil || h.HL7AEcg.Code.Code != "93000" {
					t.Errorf("Code.Code = %v, want %v", h.HL7AEcg.Code.Code, "93000")
				}
				if h.HL7AEcg.Text != "Test ECG" {
					t.Errorf("Text = %v, want %v", h.HL7AEcg.Text, "Test ECG")
				}
			},
		},
		{
			name: "Valid document with typed Code fields",
			xmlData: `<?xml version="1.0" encoding="UTF-8"?>
<AnnotatedECG xmlns="urn:hl7-org:v3">
	<id root="test-uuid"/>
	<code code="93000" codeSystem="2.16.840.1.113883.6.12" displayName="Routine ECG"/>
	<confidentialityCode code="B" codeSystem="" displayName="Blinded"/>
	<reasonCode code="PER_PROTOCOL" codeSystem=""/>
</AnnotatedECG>`,
			wantError: false,
			check: func(t *testing.T, h *Hl7xml) {
				// Check Code[CPT_CODE, CodeSystemOID]
				if h.HL7AEcg.Code.Code != types.CPT_CODE("93000") {
					t.Errorf("Code.Code = %v, want %v", h.HL7AEcg.Code.Code, types.CPT_CODE("93000"))
				}
				// Check Code[ConfidentialityCode, string]
				if h.HL7AEcg.ConfidentialityCode == nil {
					t.Fatal("ConfidentialityCode is nil")
				}
				if h.HL7AEcg.ConfidentialityCode.Code != types.ConfidentialityCode("B") {
					t.Errorf("ConfidentialityCode.Code = %v, want %v", h.HL7AEcg.ConfidentialityCode.Code, "B")
				}
				// Check Code[ReasonCode, string]
				if h.HL7AEcg.ReasonCode == nil {
					t.Fatal("ReasonCode is nil")
				}
				if h.HL7AEcg.ReasonCode.Code != types.ReasonCode("PER_PROTOCOL") {
					t.Errorf("ReasonCode.Code = %v, want %v", h.HL7AEcg.ReasonCode.Code, "PER_PROTOCOL")
				}
			},
		},
		{
			name: "Valid document with EffectiveTime",
			xmlData: `<?xml version="1.0" encoding="UTF-8"?>
<AnnotatedECG xmlns="urn:hl7-org:v3">
	<id root="test-uuid"/>
	<code code="93000" codeSystem="2.16.840.1.113883.6.12"/>
	<effectiveTime>
		<low value="20021122091000"/>
		<high value="20021122091010"/>
	</effectiveTime>
</AnnotatedECG>`,
			wantError: false,
			check: func(t *testing.T, h *Hl7xml) {
				if h.HL7AEcg.EffectiveTime == nil {
					t.Fatal("EffectiveTime is nil")
				}
				if h.HL7AEcg.EffectiveTime.Low.Value != "20021122091000" {
					t.Errorf("EffectiveTime.Low.Value = %v, want %v", h.HL7AEcg.EffectiveTime.Low.Value, "20021122091000")
				}
				if h.HL7AEcg.EffectiveTime.High.Value != "20021122091010" {
					t.Errorf("EffectiveTime.High.Value = %v, want %v", h.HL7AEcg.EffectiveTime.High.Value, "20021122091010")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("")
			err := h.Unmarshal([]byte(tt.xmlData))

			if (err != nil) != tt.wantError {
				t.Errorf("Unmarshal() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.check != nil && !tt.wantError {
				tt.check(t, h)
			}
		})
	}
}

// TestHl7xml_Unmarshal_Errors tests error handling in Unmarshal
func TestHl7xml_Unmarshal_Errors(t *testing.T) {
	tests := []struct {
		name           string
		xmlData        string
		wantError      bool
		errorContains  string
	}{
		{
			name:          "Invalid XML - malformed tags",
			xmlData:       `<AnnotatedECG><id root="test"</AnnotatedECG>`,
			wantError:     true,
			errorContains: "unmarshal AnnotatedECG",
		},
		{
			name:          "Empty data",
			xmlData:       ``,
			wantError:     true,
			errorContains: "unmarshal AnnotatedECG",
		},
		{
			name:          "Non-XML data",
			xmlData:       `this is not xml at all`,
			wantError:     true,
			errorContains: "unmarshal AnnotatedECG",
		},
		{
			name:          "Invalid XML - unclosed tag",
			xmlData:       `<AnnotatedECG xmlns="urn:hl7-org:v3"><id root="test"/>`,
			wantError:     true,
			errorContains: "unmarshal AnnotatedECG",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("")
			err := h.Unmarshal([]byte(tt.xmlData))

			if (err != nil) != tt.wantError {
				t.Errorf("Unmarshal() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError && tt.errorContains != "" {
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Unmarshal() error = %v, want error containing %v", err, tt.errorContains)
				}
			}
		})
	}
}

// TestHl7xml_Unmarshal_LenientParsing tests that unknown elements are ignored
func TestHl7xml_Unmarshal_LenientParsing(t *testing.T) {
	tests := []struct {
		name      string
		xmlData   string
		wantError bool
		check     func(t *testing.T, h *Hl7xml)
	}{
		{
			name: "Unknown elements are ignored",
			xmlData: `<?xml version="1.0" encoding="UTF-8"?>
<AnnotatedECG xmlns="urn:hl7-org:v3">
	<id root="test-uuid"/>
	<code code="93000" codeSystem="2.16.840.1.113883.6.12"/>
	<unknownElement>This should be ignored</unknownElement>
	<anotherUnknown attr="value">Content</anotherUnknown>
</AnnotatedECG>`,
			wantError: false,
			check: func(t *testing.T, h *Hl7xml) {
				// Known elements should still be parsed
				if h.HL7AEcg.ID == nil || h.HL7AEcg.ID.Root != "test-uuid" {
					t.Errorf("ID.Root = %v, want %v", h.HL7AEcg.ID.Root, "test-uuid")
				}
				if h.HL7AEcg.Code == nil || h.HL7AEcg.Code.Code != "93000" {
					t.Errorf("Code.Code = %v, want %v", h.HL7AEcg.Code.Code, "93000")
				}
			},
		},
		{
			name: "Extra attributes are ignored",
			xmlData: `<?xml version="1.0" encoding="UTF-8"?>
<AnnotatedECG xmlns="urn:hl7-org:v3" unknownAttr="value">
	<id root="test-uuid" extraAttr="ignored"/>
	<code code="93000" codeSystem="2.16.840.1.113883.6.12" customAttr="test"/>
</AnnotatedECG>`,
			wantError: false,
			check: func(t *testing.T, h *Hl7xml) {
				if h.HL7AEcg.ID.Root != "test-uuid" {
					t.Errorf("ID.Root = %v, want %v", h.HL7AEcg.ID.Root, "test-uuid")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("")
			err := h.Unmarshal([]byte(tt.xmlData))

			if (err != nil) != tt.wantError {
				t.Errorf("Unmarshal() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.check != nil && !tt.wantError {
				tt.check(t, h)
			}
		})
	}
}

// TestHl7xml_UnmarshalFromReader tests UnmarshalFromReader method
func TestHl7xml_UnmarshalFromReader(t *testing.T) {
	validXML := `<?xml version="1.0" encoding="UTF-8"?>
<AnnotatedECG xmlns="urn:hl7-org:v3">
	<id root="reader-test-uuid"/>
	<code code="93000" codeSystem="2.16.840.1.113883.6.12"/>
</AnnotatedECG>`

	t.Run("Valid reader", func(t *testing.T) {
		h := NewHl7xml("")
		reader := strings.NewReader(validXML)
		err := h.UnmarshalFromReader(reader)

		if err != nil {
			t.Errorf("UnmarshalFromReader() error = %v, want nil", err)
			return
		}

		if h.HL7AEcg.ID.Root != "reader-test-uuid" {
			t.Errorf("ID.Root = %v, want %v", h.HL7AEcg.ID.Root, "reader-test-uuid")
		}
	})

	t.Run("Reader from bytes.Buffer", func(t *testing.T) {
		h := NewHl7xml("")
		buf := bytes.NewBufferString(validXML)
		err := h.UnmarshalFromReader(buf)

		if err != nil {
			t.Errorf("UnmarshalFromReader() error = %v, want nil", err)
			return
		}

		if h.HL7AEcg.ID.Root != "reader-test-uuid" {
			t.Errorf("ID.Root = %v, want %v", h.HL7AEcg.ID.Root, "reader-test-uuid")
		}
	})

	t.Run("Invalid XML from reader", func(t *testing.T) {
		h := NewHl7xml("")
		reader := strings.NewReader("<invalid>")
		err := h.UnmarshalFromReader(reader)

		if err == nil {
			t.Error("UnmarshalFromReader() expected error, got nil")
			return
		}

		if !strings.Contains(err.Error(), "unmarshal AnnotatedECG") {
			t.Errorf("Error should contain 'unmarshal AnnotatedECG', got: %v", err)
		}
	})
}

// TestHl7xml_Unmarshal_ErrorWrapping tests that errors are properly wrapped
func TestHl7xml_Unmarshal_ErrorWrapping(t *testing.T) {
	h := NewHl7xml("")
	err := h.Unmarshal([]byte("<invalid"))

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that error can be unwrapped (errors.Unwrap)
	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		t.Error("Expected error to be wrapped (errors.Unwrap returned nil)")
	}

	// Error message should contain context
	if !strings.Contains(err.Error(), "unmarshal AnnotatedECG") {
		t.Errorf("Error should contain 'unmarshal AnnotatedECG', got: %v", err)
	}
}

// TestHl7xml_Unmarshal_FullDocument tests parsing a complete HL7 aECG document
// with series, sequences, and SequenceValue polymorphism (GLIST_TS, SLIST_PQ)
func TestHl7xml_Unmarshal_FullDocument(t *testing.T) {
	fullXML := `<?xml version="1.0" encoding="UTF-8"?>
<AnnotatedECG xmlns="urn:hl7-org:v3" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<id root="728989ec-b8bc-49cd-9a5a-30be5ade1db5"/>
	<code code="93000" codeSystem="2.16.840.1.113883.6.12" displayName="Routine ECG"/>
	<text>Test ECG Document</text>
	<effectiveTime>
		<low value="20021122091000"/>
		<high value="20021122091010"/>
	</effectiveTime>
	<confidentialityCode code="B" codeSystem="2.16.840.1.113883.5.25" displayName="Blinded"/>
	<component>
		<series classCode="OBSSER" moodCode="EVN">
			<id root="series-001"/>
			<code code="RHYTHM" codeSystem="2.16.840.1.113883.5.4"/>
			<effectiveTime>
				<low value="20021122091000.000"/>
				<high value="20021122091010.000"/>
			</effectiveTime>
			<component>
				<sequenceSet>
					<component>
						<sequence classCode="OBS" moodCode="EVN">
							<id root="seq-time-001"/>
							<code code="TIME_ABSOLUTE" codeSystem="2.16.840.1.113883.5.4"/>
							<value xsi:type="GLIST_TS">
								<head value="20021122091000.000"/>
								<increment value="0.002" unit="s"/>
							</value>
						</sequence>
					</component>
					<component>
						<sequence classCode="OBS" moodCode="EVN">
							<id root="seq-lead-001"/>
							<code code="MDC_ECG_LEAD_I" codeSystem="2.16.840.1.113883.6.24"/>
							<value xsi:type="SLIST_PQ">
								<origin value="0" unit="uV"/>
								<scale value="5" unit="uV"/>
								<digits>1 2 3 4 5 -1 -2 -3 -4 -5</digits>
							</value>
						</sequence>
					</component>
				</sequenceSet>
			</component>
		</series>
	</component>
</AnnotatedECG>`

	h := NewHl7xml("")
	err := h.Unmarshal([]byte(fullXML))

	if err != nil {
		t.Fatalf("Unmarshal() error = %v, want nil", err)
	}

	// Verify document root
	if h.HL7AEcg.ID == nil || h.HL7AEcg.ID.Root != "728989ec-b8bc-49cd-9a5a-30be5ade1db5" {
		t.Errorf("ID.Root = %v, want %v", h.HL7AEcg.ID.Root, "728989ec-b8bc-49cd-9a5a-30be5ade1db5")
	}
	if h.HL7AEcg.Code == nil || h.HL7AEcg.Code.Code != "93000" {
		t.Errorf("Code.Code = %v, want %v", h.HL7AEcg.Code.Code, "93000")
	}
	if h.HL7AEcg.Text != "Test ECG Document" {
		t.Errorf("Text = %v, want %v", h.HL7AEcg.Text, "Test ECG Document")
	}

	// Verify EffectiveTime
	if h.HL7AEcg.EffectiveTime == nil {
		t.Fatal("EffectiveTime is nil")
	}
	if h.HL7AEcg.EffectiveTime.Low.Value != "20021122091000" {
		t.Errorf("EffectiveTime.Low.Value = %v, want %v", h.HL7AEcg.EffectiveTime.Low.Value, "20021122091000")
	}

	// Verify ConfidentialityCode
	if h.HL7AEcg.ConfidentialityCode == nil {
		t.Fatal("ConfidentialityCode is nil")
	}
	if h.HL7AEcg.ConfidentialityCode.Code != types.ConfidentialityCode("B") {
		t.Errorf("ConfidentialityCode.Code = %v, want B", h.HL7AEcg.ConfidentialityCode.Code)
	}

	// Verify Component (Series) - Series is a value type within Component
	if len(h.HL7AEcg.Component) == 0 {
		t.Fatal("Component is empty, expected at least 1 series")
	}
	series := h.HL7AEcg.Component[0].Series
	if series.ID == nil || series.ID.Root != "series-001" {
		t.Errorf("Series.ID.Root = %v, want series-001", series.ID.Root)
	}
	if series.Code == nil || series.Code.Code != types.SeriesTypeCode("RHYTHM") {
		t.Errorf("Series.Code.Code = %v, want RHYTHM", series.Code.Code)
	}

	// Verify SequenceSet - SequenceSet is a value type within SeriesComponent
	if len(series.Component) == 0 {
		t.Fatal("Series.Component is empty")
	}
	seqSet := series.Component[0].SequenceSet
	if len(seqSet.Component) < 2 {
		t.Fatalf("SequenceSet.Component has %d elements, expected at least 2", len(seqSet.Component))
	}

	// Verify time sequence (GLIST_TS) - Sequence is a value type within SequenceComponent
	timeSeq := seqSet.Component[0].Sequence
	if timeSeq.Value == nil {
		t.Fatal("Time sequence Value is nil")
	}
	if timeSeq.Value.XsiType != "GLIST_TS" {
		t.Errorf("Time sequence XsiType = %v, want GLIST_TS", timeSeq.Value.XsiType)
	}
	glistTs, ok := timeSeq.Value.Typed.(*types.GLIST_TS)
	if !ok {
		t.Errorf("Time sequence Typed is not *GLIST_TS, got %T", timeSeq.Value.Typed)
	} else {
		// Head and Increment are value types, check their Value field
		if glistTs.Head.Value != "20021122091000.000" {
			t.Errorf("GLIST_TS.Head.Value = %v, want 20021122091000.000", glistTs.Head.Value)
		}
		if glistTs.Increment.Value != "0.002" {
			t.Errorf("GLIST_TS.Increment.Value = %v, want 0.002", glistTs.Increment.Value)
		}
	}

	// Verify lead sequence (SLIST_PQ)
	leadSeq := seqSet.Component[1].Sequence
	if leadSeq.Value == nil {
		t.Fatal("Lead sequence Value is nil")
	}
	if leadSeq.Value.XsiType != "SLIST_PQ" {
		t.Errorf("Lead sequence XsiType = %v, want SLIST_PQ", leadSeq.Value.XsiType)
	}
	slistPq, ok := leadSeq.Value.Typed.(*types.SLIST_PQ)
	if !ok {
		t.Errorf("Lead sequence Typed is not *SLIST_PQ, got %T", leadSeq.Value.Typed)
	} else {
		// Origin and Scale are value types (PhysicalQuantity), check their Value field
		if slistPq.Origin.Value != "0" {
			t.Errorf("SLIST_PQ.Origin.Value = %v, want 0", slistPq.Origin.Value)
		}
		if slistPq.Scale.Value != "5" {
			t.Errorf("SLIST_PQ.Scale.Value = %v, want 5", slistPq.Scale.Value)
		}
		if slistPq.Digits != "1 2 3 4 5 -1 -2 -3 -4 -5" {
			t.Errorf("SLIST_PQ.Digits = %v, want '1 2 3 4 5 -1 -2 -3 -4 -5'", slistPq.Digits)
		}
	}
}

// TestHl7xml_UnmarshalFromFile_Errors tests error handling for file-based unmarshalling
func TestHl7xml_UnmarshalFromFile_Errors(t *testing.T) {
	tests := []struct {
		name          string
		filePath      string
		wantError     bool
		errorContains string
	}{
		{
			name:          "Non-existent file",
			filePath:      "/path/to/nonexistent/file.xml",
			wantError:     true,
			errorContains: "read file",
		},
		{
			name:          "Empty file path",
			filePath:      "",
			wantError:     true,
			errorContains: "read file",
		},
		{
			name:          "Directory instead of file",
			filePath:      "/tmp",
			wantError:     true,
			errorContains: "read file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("")
			err := h.UnmarshalFromFile(tt.filePath)

			if (err != nil) != tt.wantError {
				t.Errorf("UnmarshalFromFile() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError && tt.errorContains != "" {
				if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("UnmarshalFromFile() error = %v, want error containing %v", err, tt.errorContains)
				}
			}
		})
	}
}

// TestHl7xml_UnmarshalFromFile_ErrorWrapping tests that file errors are properly wrapped
func TestHl7xml_UnmarshalFromFile_ErrorWrapping(t *testing.T) {
	h := NewHl7xml("")
	err := h.UnmarshalFromFile("/nonexistent/path/file.xml")

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check that error can be unwrapped
	unwrapped := errors.Unwrap(err)
	if unwrapped == nil {
		t.Error("Expected error to be wrapped (errors.Unwrap returned nil)")
	}

	// Verify the underlying error is a PathError (file not found)
	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Error("Expected underlying error to be *os.PathError")
	}

	// Error message should contain the file path
	if !strings.Contains(err.Error(), "/nonexistent/path/file.xml") {
		t.Errorf("Error should contain file path, got: %v", err)
	}
}

// TestHl7xml_UnmarshalFromFile_Success tests successful file-based unmarshalling
func TestHl7xml_UnmarshalFromFile_Success(t *testing.T) {
	// Create a temporary XML file for testing
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test_aecg.xml")

	validXML := `<?xml version="1.0" encoding="UTF-8"?>
<AnnotatedECG xmlns="urn:hl7-org:v3">
	<id root="file-test-uuid"/>
	<code code="93000" codeSystem="2.16.840.1.113883.6.12"/>
	<text>Test from file</text>
</AnnotatedECG>`

	err := os.WriteFile(tempFile, []byte(validXML), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	t.Run("Parse from absolute path", func(t *testing.T) {
		h := NewHl7xml("")
		err := h.UnmarshalFromFile(tempFile)

		if err != nil {
			t.Errorf("UnmarshalFromFile() error = %v, want nil", err)
			return
		}

		if h.HL7AEcg.ID == nil || h.HL7AEcg.ID.Root != "file-test-uuid" {
			t.Errorf("ID.Root = %v, want file-test-uuid", h.HL7AEcg.ID.Root)
		}
		if h.HL7AEcg.Text != "Test from file" {
			t.Errorf("Text = %v, want 'Test from file'", h.HL7AEcg.Text)
		}
	})
}

// TestHl7xml_UnmarshalFromFile_ReferenceFile tests parsing the reference aECG file
func TestHl7xml_UnmarshalFromFile_ReferenceFile(t *testing.T) {
	// Use the reference file in project root
	referenceFile := "../25060897140_23092025103550.xml"

	// Skip if reference file doesn't exist (CI environment)
	if _, err := os.Stat(referenceFile); os.IsNotExist(err) {
		t.Skip("Reference file not found, skipping integration test")
	}

	h := NewHl7xml("")
	err := h.UnmarshalFromFile(referenceFile)

	if err != nil {
		t.Fatalf("UnmarshalFromFile() error = %v, want nil", err)
	}

	// Verify basic document structure was parsed
	if h.HL7AEcg.ID == nil {
		t.Error("ID is nil after parsing reference file")
	}
	if h.HL7AEcg.Code == nil {
		t.Error("Code is nil after parsing reference file")
	}

	// Verify series/sequences were parsed (polymorphism works)
	if len(h.HL7AEcg.Component) == 0 {
		t.Error("No components (series) parsed from reference file")
	} else {
		series := h.HL7AEcg.Component[0].Series
		if len(series.Component) == 0 {
			t.Error("Series has no components (sequence sets)")
		} else {
			seqSet := series.Component[0].SequenceSet
			if len(seqSet.Component) == 0 {
				t.Error("SequenceSet has no sequences")
			} else {
				// Check first sequence has a value with xsi:type
				seq := seqSet.Component[0].Sequence
				if seq.Value != nil && seq.Value.XsiType != "" {
					t.Logf("First sequence XsiType: %s", seq.Value.XsiType)
				}
			}
		}
	}
}
