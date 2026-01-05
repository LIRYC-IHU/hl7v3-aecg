package types

// Subject identifies the subject from which the ECG waveforms were obtained.
//
// XML Structure:
//
//	<subject>
//	  <trialSubject>
//	    <id root="2.16.840.1.113883.3.55" extension="SUBJ_1"/>
//	    <code code="ENROLLED" codeSystem="2.16.840.1.113883.5.111"/>
//	    <subjectDemographicPerson>
//	      <name>BDB</name>
//	      <administrativeGenderCode code="M" codeSystem="2.16.840.1.113883.5.1"/> <birthTime value="19530508"/>
//	      <raceCode code="2106-3" codeSystem="2.16.840.1.113883.5.104"/>
//	    </subjectDemographicPerson>
//	  </trialSubject>
//	</subject>
//
// Cardinality: Required (in SubjectAssignment context)
// Reference: HL7 aECG Implementation Guide, Page 14-15
type Subject struct {
	// TrialSubject contains the subject identification and demographics.
	//
	// XML Tag: <trialSubject>...</trialSubject>
	// Cardinality: Required (within Subject)
	TrialSubject TrialSubject `xml:"trialSubject"`
}

// TrialSubject represents a subject participating in a clinical trial.
//
// This contains the subject's unique identifier, their role in the trial,
// and optional demographic information.
//
// XML Structure:
//
//	<trialSubject>
//	  <id root="2.16.840.1.113883.3.55" extension="SUBJ_1"/>
//	  <code code="ENROLLED" codeSystem="2.16.840.1.113883.5.111"/>
//	  <subjectDemographicPerson>...</subjectDemographicPerson>
//	</trialSubject>
//
// Cardinality: Required (within Subject)
type TrialSubject struct {
	// ID is the unique identifier for the subject.
	//
	// Structure:
	//   - Root: Required, must be a UID (OID or UUID)
	//   - Extension: Optional, traditional subject identifier
	//
	// Best Practice:
	//   The sponsor (or its vendor) should assign a unique OID to every subject.
	//   This OID goes into the root part. The traditional subject identifier
	//   (e.g., "SUBJ_1", "PATIENT_001") goes into the extension.
	//
	// The combination of root and extension must be universally unique
	// across all subjects globally.
	//
	// Example: root="2.16.840.1.113883.3.55" extension="SUBJ_1"
	//
	// XML Tag: <id root="..." extension="..."/>
	// Cardinality: Required
	ID *ID `xml:"id"`

	// Code indicates the role the subject was in at the time of ECG collection.
	//
	// Vocabulary: ResearchSubjectRoleBasis
	// Code System: 2.16.840.1.113883.5.111
	//
	// Suggested Values (not formally defined by HL7):
	//   - "SCREENING": Subject being screened but not yet enrolled
	//   - "ENROLLED": Subject enrolled in the trial
	//
	// Example: code="ENROLLED" codeSystem="2.16.840.1.113883.5.111"
	//
	// XML Tag: <code code="..." codeSystem="..."/>
	// Cardinality: Optional
	Code *Code[CodeRole, CodeSystemOID] `xml:"code,omitempty"`

	// SubjectDemographicPerson contains demographic information about the subject.
	//
	// This includes name (often just initials), gender, birth date, and race.
	//
	// XML Tag: <subjectDemographicPerson>...</subjectDemographicPerson>
	// Cardinality: Optional
	SubjectDemographicPerson *SubjectDemographicPerson `xml:"subjectDemographicPerson,omitempty"`
}

// SubjectDemographicPerson represents demographic information about a trial subject.
//
// Contains basic demographic data including name (often just initials for privacy),
// gender, birth date, and race.
//
// XML Structure:
//
//	<subjectDemographicPerson>
//	  <name>BDB</name>
//	  <administrativeGenderCode code="M" codeSystem="2.16.840.1.113883.5.1"/>
//	  <birthTime value="19530508"/>
//	  <raceCode code="2106-3" codeSystem="2.16.840.1.113883.5.104"/>
//	</subjectDemographicPerson>
//
// Cardinality: Optional (within TrialSubject)
// Reference: HL7 aECG Implementation Guide, Page 14-15
type SubjectDemographicPerson struct {
	// Name represents the subject's name.
	//
	// Privacy Note: Often only the subject's initials are used to protect privacy.
	// When using initials, put them directly in the name field without additional
	// structure (e.g., just "BDB" instead of separate given/family fields).
	//
	// For full names, you can use structured PersonName with given, family, etc.
	//
	// Examples:
	//   - Initials only: "BDB"
	//   - Full structured name: see PersonName type
	//
	// XML Tag: <name>...</name>
	// Cardinality: Optional
	Name *string `xml:"name,omitempty"`

	// AdministrativeGenderCode indicates the subject's administrative gender.
	//
	// Vocabulary: AdministrativeGender
	// Code System: 2.16.840.1.113883.5.1
	//
	// Defined Codes:
	//   - "F": Female
	//   - "M": Male
	//   - "UN": Undifferentiated
	//
	// Example: code="M" codeSystem="2.16.840.1.113883.5.1"
	//
	// XML Tag: <administrativeGenderCode code="..." codeSystem="..."/>
	// Cardinality: Optional
	AdministrativeGenderCode *Code[GenderCode, CodeSystemOID] `xml:"administrativeGenderCode,omitempty"`

	// BirthTime is the subject's date of birth.
	//
	// Format: YYYYMMDD
	// Example: "19530508" = May 8, 1953
	//
	// Note: Time precision is typically date only (no hours/minutes).
	//
	// XML Tag: <birthTime value="..."/>
	// Cardinality: Optional
	BirthTime *Time `xml:"birthTime"`

	// RaceCode indicates the subject's race.
	//
	// Vocabulary: Race
	// Code System: 2.16.840.1.113883.5.104
	//
	// Common Codes (see HL7 documentation for complete list):
	//   - "1002-5": Native American
	//   - "2028-9": Asian
	//   - "2054-5": Black or African American
	//   - "2076-8": Hawaiian or Pacific Islander
	//   - "2106-3": White
	//   - "2131-1": Other Race
	//
	// Example: code="2106-3" codeSystem="2.16.840.1.113883.5.104"
	//
	// XML Tag: <raceCode code="..." codeSystem="..."/>
	// Cardinality: Optional
	RaceCode *Code[RaceCode, CodeSystemOID] `xml:"raceCode,omitempty"`

	// PatientID is the primary patient identifier.
	//
	// This is typically the hospital or institution's internal patient ID.
	//
	// Example: "25060897140"
	//
	// XML Tag: <PatientID>...</PatientID>
	// Cardinality: Optional
	PatientID string `xml:"PatientID"`

	// SecondPatientID is an optional secondary patient identifier.
	//
	// Used when the patient has multiple identification numbers
	// (e.g., different hospital systems).
	//
	// XML Tag: <SecondPatientID>...</SecondPatientID>
	// Cardinality: Optional
	SecondPatientID string `xml:"SecondPatientID"`

	// Age is the subject's age at the time of ECG acquisition.
	//
	// Can be represented as a number (years) or other format.
	//
	// XML Tag: <Age>...</Age>
	// Cardinality: Optional
	Age string `xml:"Age"`

	// Paced indicates whether the patient has a cardiac pacemaker.
	//
	// True if the patient is paced, false otherwise.
	//
	// Example: true
	//
	// XML Tag: <Paced>...</Paced>
	// Cardinality: Optional
	Paced bool `xml:"Paced"`

	// Medications contains the list of medications the patient is taking.
	//
	// Each medication can include name, dosage, etc.
	//
	// XML Tag: <Medications>...</Medications>
	// Cardinality: Optional
	Medications Medications `xml:"Medications"`

	// ClinicalClassifications contains clinical classification information.
	//
	// Used to categorize the patient's clinical status or conditions.
	//
	// XML Tag: <ClinicalClassifications>...</ClinicalClassifications>
	// Cardinality: Optional
	ClinicalClassifications ClinicalClassifications `xml:"ClinicalClassifications"`

	// Bed is the patient's bed location within the facility.
	//
	// Example: "12A"
	//
	// XML Tag: <Bed>...</Bed>
	// Cardinality: Optional
	Bed string `xml:"Bed"`

	// Room is the patient's room number or identifier.
	//
	// Example: "302"
	//
	// XML Tag: <Room>...</Room>
	// Cardinality: Optional
	Room string `xml:"Room"`

	// PointOfCare identifies the care unit or department.
	//
	// Example: "Cardiology ICU", "Emergency Department"
	//
	// XML Tag: <PointOfCare>...</PointOfCare>
	// Cardinality: Optional
	PointOfCare string `xml:"PointOfCare"`
}

// Medications represents a list of medications the patient is taking.
//
// XML Structure:
//
//	<Medications>
//	  <Medication>Aspirin</Medication>
//	  <Medication>Metoprolol</Medication>
//	</Medications>
//
// Cardinality: Optional (within SubjectDemographicPerson)
type Medications struct {
	// Medication is a list of individual medications.
	//
	// Each entry can be a medication name, dosage information, or structured data.
	//
	// XML Tag: <Medication>...</Medication>
	// Cardinality: 0..* (multiple medications allowed)
	Medication []string `xml:"Medication"`
}

// ClinicalClassifications represents clinical classification information.
//
// Used to categorize the patient's clinical status, conditions, or diagnoses.
//
// XML Structure:
//
//	<ClinicalClassifications>
//	  <ClinicalClassification>Hypertension</ClinicalClassification>
//	  <ClinicalClassification>Diabetes Type 2</ClinicalClassification>
//	</ClinicalClassifications>
//
// Cardinality: Optional (within SubjectDemographicPerson)
type ClinicalClassifications struct {
	// ClinicalClassification is a list of clinical classifications.
	//
	// Each entry represents a clinical category, diagnosis, or classification.
	//
	// XML Tag: <ClinicalClassification>...</ClinicalClassification>
	// Cardinality: 0..* (multiple classifications allowed)
	ClinicalClassification []string `xml:"ClinicalClassification"`
}
