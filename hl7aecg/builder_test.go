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

			if h.HL7AEcg.ComponentOf == nil {
				t.Errorf("SetSubject() did not set ComponentOf")
				return
			}

			subject := &h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.SubjectAssignment.Subject
			if subject.TrialSubject.ID == nil {
				t.Errorf("SetSubject() did not set Subject.ID")
				return
			}

			if subject.TrialSubject.ID.Extension != tt.extension {
				t.Errorf("SetSubject() Extension = %v, want %v", subject.TrialSubject.ID.Extension, tt.extension)
			}
		})
	}
}

// TestSetSubjectDemographics tests the SetSubjectDemographics method
func TestSetSubjectDemographics(t *testing.T) {
	tests := []struct {
		name      string
		subjName  string
		patientID string
		gender    types.GenderCode
		birthDate string
		race      types.RaceCode
	}{
		{
			name:      "Set complete demographics",
			subjName:  "JDO",
			patientID: "PAT-001",
			gender:    types.GENDER_MALE,
			birthDate: "19800101",
			race:      types.RACE_WHITE,
		},
		{
			name:      "Set demographics with female",
			subjName:  "ABC",
			patientID: "PAT-002",
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

			result := h.SetSubjectDemographics(tt.subjName, tt.patientID, tt.gender, tt.birthDate, tt.race)

			if result == nil {
				t.Errorf("SetSubjectDemographics() returned nil")
				return
			}

			if h.HL7AEcg.ComponentOf == nil ||
				h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.SubjectAssignment.Subject.TrialSubject.SubjectDemographicPerson == nil {
				t.Errorf("SetSubjectDemographics() did not set SubjectDemographicPerson")
				return
			}

			person := h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.SubjectAssignment.Subject.TrialSubject.SubjectDemographicPerson

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
		SetSubjectDemographics("JDO", "PAT-001", types.GENDER_MALE, "19800101", types.RACE_WHITE)

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

	if h.HL7AEcg.ComponentOf == nil {
		t.Errorf("Fluent API did not set ComponentOf")
	} else {
		if h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.SubjectAssignment.Subject.TrialSubject.SubjectDemographicPerson == nil {
			t.Errorf("Fluent API did not set SubjectDemographicPerson")
		}
	}
}

// TestSetResponsibleParty tests the SetResponsibleParty method
func TestSetResponsibleParty(t *testing.T) {
	tests := []struct {
		name              string
		investigatorRoot  string
		investigatorID    string
		prefix            string
		given             string
		family            string
		suffix            string
		wantEmptyName     bool
		wantNameWithValue bool
	}{
		{
			name:              "Set responsible party with full name",
			investigatorRoot:  "2.16.840.1.113883.3.6",
			investigatorID:    "INV_001",
			prefix:            "Dr.",
			given:             "John",
			family:            "Smith",
			suffix:            "MD",
			wantEmptyName:     false,
			wantNameWithValue: true,
		},
		{
			name:              "Set responsible party with empty name",
			investigatorRoot:  "",
			investigatorID:    "trialInvestigator",
			prefix:            "",
			given:             "",
			family:            "",
			suffix:            "",
			wantEmptyName:     true,
			wantNameWithValue: false,
		},
		{
			name:              "Set responsible party with partial name",
			investigatorRoot:  "2.16.840.1.113883.3.6",
			investigatorID:    "INV_002",
			prefix:            "",
			given:             "Jane",
			family:            "Doe",
			suffix:            "",
			wantEmptyName:     false,
			wantNameWithValue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHl7xml("/tmp/test").
				Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "").
				SetSubject("SUBJ-001", "001", types.SUBJECT_ROLE_ENROLLED)

			result := h.SetResponsibleParty(tt.investigatorRoot, tt.investigatorID, tt.prefix, tt.given, tt.family, tt.suffix)

			if result == nil {
				t.Errorf("SetResponsibleParty() returned nil")
				return
			}

			// Check that ComponentOf exists
			if h.HL7AEcg.ComponentOf == nil {
				t.Errorf("SetResponsibleParty() did not initialize ComponentOf")
				return
			}

			// Navigate to ResponsibleParty
			clinicalTrial := &h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.SubjectAssignment.ComponentOf.ClinicalTrial
			if clinicalTrial.Location == nil {
				t.Errorf("SetResponsibleParty() did not initialize Location")
				return
			}

			rp := clinicalTrial.Location.TrialSite.ResponsibleParty
			if rp == nil {
				t.Errorf("SetResponsibleParty() did not set ResponsibleParty")
				return
			}

			// Check investigator ID extension
			if rp.TrialInvestigator.ID.Extension != tt.investigatorID {
				t.Errorf("InvestigatorID Extension = %v, want %v", rp.TrialInvestigator.ID.Extension, tt.investigatorID)
			}

			// Check investigator ID root (when provided)
			if tt.investigatorRoot != "" && rp.TrialInvestigator.ID.Root != tt.investigatorRoot {
				t.Errorf("InvestigatorID Root = %v, want %v", rp.TrialInvestigator.ID.Root, tt.investigatorRoot)
			}

			// Check InvestigatorPerson and Name
			if rp.TrialInvestigator.InvestigatorPerson == nil {
				t.Errorf("SetResponsibleParty() did not initialize InvestigatorPerson")
				return
			}

			if rp.TrialInvestigator.InvestigatorPerson.Name == nil {
				t.Errorf("SetResponsibleParty() did not initialize Name")
				return
			}

			name := rp.TrialInvestigator.InvestigatorPerson.Name

			if tt.wantEmptyName {
				// All components should be nil for empty name
				if name.Prefix != nil || name.Given != nil || name.Family != nil || name.Suffix != nil {
					t.Errorf("Name should be empty but has values")
				}
			}

			if tt.wantNameWithValue {
				if tt.prefix != "" && (name.Prefix == nil || *name.Prefix != tt.prefix) {
					t.Errorf("Name.Prefix = %v, want %v", name.Prefix, tt.prefix)
				}
				if tt.given != "" && (name.Given == nil || *name.Given != tt.given) {
					t.Errorf("Name.Given = %v, want %v", name.Given, tt.given)
				}
				if tt.family != "" && (name.Family == nil || *name.Family != tt.family) {
					t.Errorf("Name.Family = %v, want %v", name.Family, tt.family)
				}
				if tt.suffix != "" && (name.Suffix == nil || *name.Suffix != tt.suffix) {
					t.Errorf("Name.Suffix = %v, want %v", name.Suffix, tt.suffix)
				}
			}
		})
	}
}

// TestSetResponsibleParty_MethodChaining tests method chaining with SetResponsibleParty
func TestSetResponsibleParty_MethodChaining(t *testing.T) {
	h := NewHl7xml("/tmp/test").
		Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "", "").
		SetSubject("SUBJ-001", "001", types.SUBJECT_ROLE_ENROLLED).
		SetLocation("SITE_001", "2.16.840.1.113883.3.5", "Test Site", "Boston", "MA", "USA").
		SetResponsibleParty("", "INV_001", "Dr.", "John", "Smith", "MD")

	if h == nil {
		t.Fatal("Method chaining returned nil")
	}

	// Verify Location and ResponsibleParty are both set
	if h.HL7AEcg.ComponentOf == nil {
		t.Fatal("ComponentOf not initialized")
	}

	clinicalTrial := &h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.SubjectAssignment.ComponentOf.ClinicalTrial
	if clinicalTrial.Location == nil {
		t.Fatal("Location not set")
	}

	// Check Location
	if clinicalTrial.Location.TrialSite.Location == nil || clinicalTrial.Location.TrialSite.Location.Name == nil {
		t.Error("Location name not set")
	}

	// Check ResponsibleParty
	if clinicalTrial.Location.TrialSite.ResponsibleParty == nil {
		t.Error("ResponsibleParty not set")
	}

	rp := clinicalTrial.Location.TrialSite.ResponsibleParty
	if rp.TrialInvestigator.ID.Extension != "INV_001" {
		t.Errorf("Investigator ID not set correctly: %v", rp.TrialInvestigator.ID.Extension)
	}

	if rp.TrialInvestigator.InvestigatorPerson == nil || rp.TrialInvestigator.InvestigatorPerson.Name == nil {
		t.Error("Investigator name not initialized")
	}
}

// TestAddSecondaryPerformer tests adding a secondary performer to a series
func TestAddSecondaryPerformer(t *testing.T) {
	h := NewHl7xml("/tmp")
	h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "CPT-4", "")

	// Add a series first
	leads := map[types.LeadCode][]int{
		types.MDC_ECG_LEAD_I: {1, 2, 3},
	}
	h.AddRhythmSeries("20021122091000.000", "20021122091010.000", 500, leads, 0, 5)

	// Add secondary performer
	h.AddSecondaryPerformer(types.PERFORMER_ECG_TECHNICIAN, "", "", "")

	// Verify secondary performer was added
	if len(h.HL7AEcg.Component) == 0 {
		t.Fatal("No components in document")
	}

	series := &h.HL7AEcg.Component[0].Series
	if len(series.SecondaryPerformer) == 0 {
		t.Fatal("No secondary performers added")
	}

	performer := series.SecondaryPerformer[0]
	if performer.FunctionCode == nil {
		t.Error("FunctionCode should be set")
	} else if performer.FunctionCode.Code != types.PERFORMER_ECG_TECHNICIAN {
		t.Errorf("FunctionCode = %v, want %v", performer.FunctionCode.Code, types.PERFORMER_ECG_TECHNICIAN)
	}

	if performer.SeriesPerformer.AssignedPerson == nil {
		t.Error("AssignedPerson should be set")
	} else if performer.SeriesPerformer.AssignedPerson.Name == nil {
		t.Error("Name should be set")
	} else if *performer.SeriesPerformer.AssignedPerson.Name != "" {
		t.Errorf("Name = %v, want empty string", *performer.SeriesPerformer.AssignedPerson.Name)
	}
}

// TestAddSecondaryPerformer_WithDetails tests adding a secondary performer with full details
func TestAddSecondaryPerformer_WithDetails(t *testing.T) {
	h := NewHl7xml("/tmp")
	h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "CPT-4", "")

	// Add a series first
	leads := map[types.LeadCode][]int{
		types.MDC_ECG_LEAD_I: {1, 2, 3},
	}
	h.AddRhythmSeries("20021122091000.000", "20021122091010.000", 500, leads, 0, 5)

	// Add secondary performer with details
	h.AddSecondaryPerformer(types.PERFORMER_HOLTER_ANALYST, "2.16.840.1.113883.3.4", "TECH-221", "KAB")

	// Verify
	performer := h.HL7AEcg.Component[0].Series.SecondaryPerformer[0]

	if performer.FunctionCode.Code != types.PERFORMER_HOLTER_ANALYST {
		t.Errorf("FunctionCode = %v, want %v", performer.FunctionCode.Code, types.PERFORMER_HOLTER_ANALYST)
	}

	if performer.SeriesPerformer.ID == nil {
		t.Fatal("Performer ID should be set")
	}

	if performer.SeriesPerformer.ID.Root != "2.16.840.1.113883.3.4" {
		t.Errorf("Performer ID Root = %v, want %v", performer.SeriesPerformer.ID.Root, "2.16.840.1.113883.3.4")
	}

	if performer.SeriesPerformer.ID.Extension != "TECH-221" {
		t.Errorf("Performer ID Extension = %v, want %v", performer.SeriesPerformer.ID.Extension, "TECH-221")
	}

	if *performer.SeriesPerformer.AssignedPerson.Name != "KAB" {
		t.Errorf("Performer Name = %v, want %v", *performer.SeriesPerformer.AssignedPerson.Name, "KAB")
	}
}

// TestAddSecondaryPerformer_MultiplePerformers tests adding multiple performers to a series
func TestAddSecondaryPerformer_MultiplePerformers(t *testing.T) {
	h := NewHl7xml("/tmp")
	h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "CPT-4", "")

	// Add a series first
	leads := map[types.LeadCode][]int{
		types.MDC_ECG_LEAD_I: {1, 2, 3},
	}
	h.AddRhythmSeries("20021122091000.000", "20021122091010.000", 500, leads, 0, 5)

	// Add multiple secondary performers
	h.AddSecondaryPerformer(types.PERFORMER_HOLTER_HOOKUP, "", "", "Tech1").
		AddSecondaryPerformer(types.PERFORMER_HOLTER_ANALYST, "", "", "Analyst1")

	// Verify
	series := &h.HL7AEcg.Component[0].Series
	if len(series.SecondaryPerformer) != 2 {
		t.Fatalf("Expected 2 secondary performers, got %d", len(series.SecondaryPerformer))
	}

	// Check first performer
	if series.SecondaryPerformer[0].FunctionCode.Code != types.PERFORMER_HOLTER_HOOKUP {
		t.Errorf("First performer function code = %v, want %v", series.SecondaryPerformer[0].FunctionCode.Code, types.PERFORMER_HOLTER_HOOKUP)
	}

	// Check second performer
	if series.SecondaryPerformer[1].FunctionCode.Code != types.PERFORMER_HOLTER_ANALYST {
		t.Errorf("Second performer function code = %v, want %v", series.SecondaryPerformer[1].FunctionCode.Code, types.PERFORMER_HOLTER_ANALYST)
	}
}

// TestAddSecondaryPerformer_NoSeries tests that adding performer without series doesn't crash
func TestAddSecondaryPerformer_NoSeries(t *testing.T) {
	h := NewHl7xml("/tmp")
	h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "CPT-4", "")

	// Try to add secondary performer without a series - should not crash
	h.AddSecondaryPerformer(types.PERFORMER_ECG_TECHNICIAN, "", "", "")

	// Document should still be valid, just no performers added
	if len(h.HL7AEcg.Component) > 0 {
		t.Error("Should not have any components")
	}
}

