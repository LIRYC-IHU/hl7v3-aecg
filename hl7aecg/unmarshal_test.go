package hl7aecg

import (
	"testing"
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
			name: "clinical trial with minimal fields",
			xmlData: `<?xml version="1.0" encoding="UTF-8"?>
<AnnotatedECG xmlns="urn:hl7-org:v3">
	<id root="minimal-trial-doc"/>
	<code code="93000" codeSystem="2.16.840.1.113883.6.12"/>
	<clinicalTrial>
		<id root="trial-001"/>
	</clinicalTrial>
</AnnotatedECG>`,
			check: func(t *testing.T, h *Hl7xml) {
				if h.HL7AEcg.ClinicalTrial == nil {
					t.Fatal("ClinicalTrial is nil")
				}
				if h.HL7AEcg.ClinicalTrial.ID.Root != "trial-001" {
					t.Errorf("ClinicalTrial.ID.Root = %v, want trial-001", h.HL7AEcg.ClinicalTrial.ID.Root)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("")
			err := h.Unmarshal([]byte(tt.xmlData))
			if err != nil {
				t.Fatalf("Unmarshal() error = %v, want nil", err)
			}
			tt.check(t, h)
		})
	}
}

// TestHl7xml_AccessClinicalTrialInformation tests accessing clinical trial information after unmarshalling
func TestHl7xml_AccessClinicalTrialInformation(t *testing.T) {
	tests := []struct {
		name    string
		xmlData string
		check   func(t *testing.T, h *Hl7xml)
	}{
		{
			name: "clinical trial with all fields",
			xmlData: `<?xml version="1.0" encoding="UTF-8"?>
<AnnotatedECG xmlns="urn:hl7-org:v3">
	<id root="doc-with-trial"/>
	<code code="93000" codeSystem="2.16.840.1.113883.6.12"/>
	<clinicalTrial>
		<id root="trial-123" extension="protocol"/>
		<activityTime>
			<low value="20240101100000"/>
			<high value="20240131235959"/>
		</activityTime>
		<title>Test Clinical Trial</title>
	</clinicalTrial>
</AnnotatedECG>`,
			check: func(t *testing.T, h *Hl7xml) {
				if h.HL7AEcg.ClinicalTrial == nil {
					t.Fatal("ClinicalTrial is nil")
				}
				if h.HL7AEcg.ClinicalTrial.ID.Root != "trial-123" {
					t.Errorf("ClinicalTrial.ID.Root = %v, want trial-123", h.HL7AEcg.ClinicalTrial.ID.Root)
				}
				if h.HL7AEcg.ClinicalTrial.ActivityTime == nil {
					t.Fatal("ActivityTime is nil")
				}
				if h.HL7AEcg.ClinicalTrial.ActivityTime.Low.Value != "20240101100000" {
					t.Errorf("ActivityTime.Low.Value = %v, want 20240101100000", h.HL7AEcg.ClinicalTrial.ActivityTime.Low.Value)
				}
				if h.HL7AEcg.ClinicalTrial.Title == nil || *h.HL7AEcg.ClinicalTrial.Title != "Test Clinical Trial" {
					t.Errorf("Title = %v, want Test Clinical Trial", h.HL7AEcg.ClinicalTrial.Title)
				}
			},
		},
		{
			name: "clinical trial with minimal fields",
			xmlData: `<?xml version="1.0" encoding="UTF-8"?>
<AnnotatedECG xmlns="urn:hl7-org:v3">
	<id root="minimal-trial-doc"/>
	<code code="93000" codeSystem="2.16.840.1.113883.6.12"/>
	<clinicalTrial>
		<id root="trial-001"/>
	</clinicalTrial>
</AnnotatedECG>`,
			check: func(t *testing.T, h *Hl7xml) {
				if h.HL7AEcg.ClinicalTrial == nil {
					t.Fatal("ClinicalTrial is nil")
				}
				if h.HL7AEcg.ClinicalTrial.ID.Root != "trial-001" {
					t.Errorf("ClinicalTrial.ID.Root = %v, want trial-001", h.HL7AEcg.ClinicalTrial.ID.Root)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("")
			err := h.Unmarshal([]byte(tt.xmlData))
			if err != nil {
				t.Fatalf("Unmarshal() error = %v, want nil", err)
			}
			tt.check(t, h)
		})
	}
}

// TestHl7xml_AccessSubjectInformation tests accessing subject information after unmarshalling
func TestHl7xml_AccessSubjectInformation(t *testing.T) {
	tests := []struct {
		name    string
		xmlData string
		check   func(t *testing.T, h *Hl7xml)
	}{
		{
			name: "subject with minimal fields",
			xmlData: `<?xml version="1.0" encoding="UTF-8"?>
<AnnotatedECG xmlns="urn:hl7-org:v3">
	<id root="minimal-subject-doc"/>
	<code code="93000" codeSystem="2.16.840.1.113883.6.12"/>
	<subject>
		<id root="subject-001"/>
		<code code="ENROLLED" codeSystem="2.16.840.1.113883.5.110"/>
	</subject>
</AnnotatedECG>`,
			check: func(t *testing.T, h *Hl7xml) {
				if h.HL7AEcg.Subject == nil {
					t.Fatal("Subject is nil")
				}
				// Subject is *TrialSubject, access fields directly
				if h.HL7AEcg.Subject.ID == nil || h.HL7AEcg.Subject.ID.Root != "subject-001" {
					t.Errorf("Subject.ID.Root = %v, want subject-001", h.HL7AEcg.Subject.ID)
				}
				if h.HL7AEcg.Subject.Code == nil || string(h.HL7AEcg.Subject.Code.Code) != "ENROLLED" {
					t.Errorf("Subject.Code.Code = %v, want ENROLLED", h.HL7AEcg.Subject.Code.Code)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("")
			err := h.Unmarshal([]byte(tt.xmlData))
			if err != nil {
				t.Fatalf("Unmarshal() error = %v, want nil", err)
			}
			tt.check(t, h)
		})
	}
}
