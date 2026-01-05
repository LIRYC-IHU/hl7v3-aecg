package hl7aecg

import (
	"testing"

	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
)

// TestNewHl7xml tests the creation of a new Hl7xml instance
func TestNewHl7xml(t *testing.T) {
	tests := []struct {
		name      string
		outputDir string
	}{
		{
			name:      "Create with output directory",
			outputDir: "/tmp/test",
		},
		{
			name:      "Create with empty output directory",
			outputDir: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml(tt.outputDir)

			if h == nil {
				t.Errorf("NewHl7xml() returned nil")
				return
			}

			// Note: outputDir is unexported, so we can't test it directly
			// but we can verify the instance was created successfully
		})
	}
}

// TestInitialize tests the Initialize method
func TestInitialize(t *testing.T) {
	tests := []struct {
		name       string
		code       types.CPT_CODE
		codeSystem types.CodeSystemOID
	}{
		{
			name:       "Initialize with routine ECG code",
			code:       types.CPT_CODE_ECG_Routine,
			codeSystem: types.CPT_OID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("/tmp/test")
			result := h.Initialize(tt.code, tt.codeSystem, "", "")

			if result == nil {
				t.Errorf("Initialize() returned nil")
				return
			}

			if h.HL7AEcg.Code == nil {
				t.Errorf("Initialize() did not set Code")
				return
			}

			if h.HL7AEcg.Code.Code != tt.code {
				t.Errorf("Initialize() Code = %v, want %v", h.HL7AEcg.Code.Code, tt.code)
			}

			if h.HL7AEcg.Code.CodeSystem != tt.codeSystem {
				t.Errorf("Initialize() CodeSystem = %v, want %v", h.HL7AEcg.Code.CodeSystem, tt.codeSystem)
			}

			if h.HL7AEcg.ID == nil {
				t.Errorf("Initialize() did not create ID")
			}

			if h.HL7AEcg.XmlnsXsi == "" {
				t.Errorf("Initialize() did not set XmlnsXsi")
			}
		})
	}
}

// TestSetText tests the SetText method
func TestSetText(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{
			name: "Set simple text",
			text: "Test ECG",
		},
		{
			name: "Set complex text",
			text: "12-lead ECG for clinical trial ABC-123",
		},
		{
			name: "Set empty text",
			text: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("/tmp/test").Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "")
			result := h.SetText(tt.text)

			if result == nil {
				t.Errorf("SetText() returned nil")
				return
			}

			if h.HL7AEcg.Text != tt.text {
				t.Errorf("SetText() Text = %v, want %v", h.HL7AEcg.Text, tt.text)
			}
		})
	}
}

// TestSetEffectiveTime tests the SetEffectiveTime method
func TestSetEffectiveTime(t *testing.T) {
	tests := []struct {
		name string
		low  string
		high string
	}{
		{
			name: "Set valid time range",
			low:  "20231223120000",
			high: "20231223120010",
		},
		{
			name: "Set time range with milliseconds",
			low:  "20231223120000.000",
			high: "20231223120010.000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("/tmp/test").Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "")
			result := h.SetEffectiveTime(tt.low, tt.high)

			if result == nil {
				t.Errorf("SetEffectiveTime() returned nil")
				return
			}

			if h.HL7AEcg.EffectiveTime == nil {
				t.Errorf("SetEffectiveTime() did not set EffectiveTime")
				return
			}

			if h.HL7AEcg.EffectiveTime.Low.Value != tt.low {
				t.Errorf("SetEffectiveTime() Low = %v, want %v", h.HL7AEcg.EffectiveTime.Low.Value, tt.low)
			}

			if h.HL7AEcg.EffectiveTime.High.Value != tt.high {
				t.Errorf("SetEffectiveTime() High = %v, want %v", h.HL7AEcg.EffectiveTime.High.Value, tt.high)
			}
		})
	}
}

// TestSetSubject tests the SetSubject method
func TestSetSubject(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		extension string
		code      types.CodeRole
	}{
		{
			name:      "Set enrolled subject",
			id:        "SUBJ-001",
			extension: "001",
			code:      types.SUBJECT_ROLE_ENROLLED,
		},
		{
			name:      "Set screening subject",
			id:        "SUBJ-002",
			extension: "002",
			code:      types.SUBJECT_ROLE_SCREENING,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("/tmp/test").Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "")
			result := h.SetSubject(tt.id, tt.extension, tt.code)

			if result == nil {
				t.Errorf("SetSubject() returned nil")
				return
			}

			if h.HL7AEcg.Subject == nil {
				t.Errorf("SetSubject() did not set Subject")
				return
			}

			if h.HL7AEcg.Subject.ID == nil {
				t.Errorf("SetSubject() did not set Subject.ID")
				return
			}

			if h.HL7AEcg.Subject.ID.Extension != tt.extension {
				t.Errorf("SetSubject() Extension = %v, want %v", h.HL7AEcg.Subject.ID.Extension, tt.extension)
			}
		})
	}
}

// TestSetSubjectDemographics tests the SetSubjectDemographics method
func TestSetSubjectDemographics(t *testing.T) {
	tests := []struct {
		name      string
		subjName  string
		gender    types.GenderCode
		birthDate string
		race      types.RaceCode
	}{
		{
			name:      "Set complete demographics",
			subjName:  "JDO",
			gender:    types.GENDER_MALE,
			birthDate: "19800101",
			race:      types.RACE_WHITE,
		},
		{
			name:      "Set demographics with female",
			subjName:  "ABC",
			gender:    types.GENDER_FEMALE,
			birthDate: "19900615",
			race:      types.RACE_BLACK_OR_AFRICAN_AMERICAN,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("/tmp/test").
				Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "").
				SetSubject("SUBJ-001", "001", types.SUBJECT_ROLE_ENROLLED)

			result := h.SetSubjectDemographics(tt.subjName, tt.gender, tt.birthDate, tt.race)

			if result == nil {
				t.Errorf("SetSubjectDemographics() returned nil")
				return
			}

			if h.HL7AEcg.Subject.SubjectDemographicPerson == nil {
				t.Errorf("SetSubjectDemographics() did not set SubjectDemographicPerson")
				return
			}

			person := h.HL7AEcg.Subject.SubjectDemographicPerson

			if person.Name == nil || *person.Name != tt.subjName {
				t.Errorf("SetSubjectDemographics() Name = %v, want %v", person.Name, tt.subjName)
			}

			if person.AdministrativeGenderCode == nil || person.AdministrativeGenderCode.Code != tt.gender {
				t.Errorf("SetSubjectDemographics() Gender = %v, want %v", person.AdministrativeGenderCode.Code, tt.gender)
			}

			if person.BirthTime == nil || person.BirthTime.Value != tt.birthDate {
				t.Errorf("SetSubjectDemographics() BirthDate = %v, want %v", person.BirthTime.Value, tt.birthDate)
			}
		})
	}
}

// TestAddRhythmSeries tests the AddRhythmSeries method
func TestAddRhythmSeries(t *testing.T) {
	tests := []struct {
		name       string
		startTime  string
		endTime    string
		sampleRate float64
		leads      map[types.LeadCode][]int
		origin     float64
		scale      float64
		wantLeads  int
	}{
		{
			name:       "Add rhythm series with 3 leads",
			startTime:  "20231223120000.000",
			endTime:    "20231223120010.000",
			sampleRate: 500.0,
			leads: map[types.LeadCode][]int{
				types.MDC_ECG_LEAD_I:  {1, 2, 3, 4, 5},
				types.MDC_ECG_LEAD_II: {2, 3, 4, 5, 6},
				types.MDC_ECG_LEAD_V1: {3, 4, 5, 6, 7},
			},
			origin:    0.0,
			scale:     5.0,
			wantLeads: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("/tmp/test").Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "")
			result := h.AddRhythmSeries(tt.startTime, tt.endTime, tt.sampleRate, tt.leads, tt.origin, tt.scale)

			if result == nil {
				t.Errorf("AddRhythmSeries() returned nil")
				return
			}

			if len(h.HL7AEcg.Component) == 0 {
				t.Errorf("AddRhythmSeries() did not add component")
				return
			}

			series := h.HL7AEcg.Component[0].Series

			if series.Code == nil {
				t.Errorf("AddRhythmSeries() did not set series code")
				return
			}

			if series.Code.Code != types.RHYTHM_CODE {
				t.Errorf("AddRhythmSeries() Code = %v, want %v", series.Code.Code, types.RHYTHM_CODE)
			}

			if len(series.Component) == 0 {
				t.Errorf("AddRhythmSeries() did not add sequences")
				return
			}

			sequenceSet := series.Component[0].SequenceSet

			// +1 for time sequence
			expectedSequences := tt.wantLeads + 1
			if len(sequenceSet.Component) != expectedSequences {
				t.Errorf("AddRhythmSeries() sequences = %d, want %d", len(sequenceSet.Component), expectedSequences)
			}
		})
	}
}

// TestAddRepresentativeBeatSeries tests the AddRepresentativeBeatSeries method
func TestAddRepresentativeBeatSeries(t *testing.T) {
	tests := []struct {
		name       string
		startTime  string
		endTime    string
		sampleRate float64
		leads      map[types.LeadCode][]int
		origin     float64
		scale      float64
	}{
		{
			name:       "Add representative beat series",
			startTime:  "20231223120000.000",
			endTime:    "20231223120001.000",
			sampleRate: 500.0,
			leads: map[types.LeadCode][]int{
				types.MDC_ECG_LEAD_I:  {1, 2, 3},
				types.MDC_ECG_LEAD_II: {2, 3, 4},
			},
			origin: 0.0,
			scale:  5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("/tmp/test").Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "")
			result := h.AddRepresentativeBeatSeries(tt.startTime, tt.endTime, tt.sampleRate, tt.leads, tt.origin, tt.scale)

			if result == nil {
				t.Errorf("AddRepresentativeBeatSeries() returned nil")
				return
			}

			if len(h.HL7AEcg.Component) == 0 {
				t.Errorf("AddRepresentativeBeatSeries() did not add component")
				return
			}

			series := h.HL7AEcg.Component[0].Series

			if series.Code == nil {
				t.Errorf("AddRepresentativeBeatSeries() did not set series code")
				return
			}

			if series.Code.Code != types.REPRESENTATIVE_BEAT_CODE {
				t.Errorf("AddRepresentativeBeatSeries() Code = %v, want %v", series.Code.Code, types.REPRESENTATIVE_BEAT_CODE)
			}
		})
	}
}

// TestFluentAPI tests the fluent API (method chaining)
func TestFluentAPI(t *testing.T) {
	h := NewHl7xml("/tmp/test").
		Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "CPT-4", "").
		SetText("Test ECG").
		SetEffectiveTime("20231223120000", "20231223120010").
		SetSubject("SUBJ-001", "001", types.SUBJECT_ROLE_ENROLLED).
		SetSubjectDemographics("JDO", types.GENDER_MALE, "19800101", types.RACE_WHITE)

	if h == nil {
		t.Errorf("Fluent API returned nil")
		return
	}

	// Verify all fields were set
	if h.HL7AEcg.Text != "Test ECG" {
		t.Errorf("Fluent API did not set Text")
	}

	if h.HL7AEcg.EffectiveTime == nil {
		t.Errorf("Fluent API did not set EffectiveTime")
	}

	if h.HL7AEcg.Subject == nil {
		t.Errorf("Fluent API did not set Subject")
	}

	if h.HL7AEcg.Subject.SubjectDemographicPerson == nil {
		t.Errorf("Fluent API did not set SubjectDemographicPerson")
	}
}
