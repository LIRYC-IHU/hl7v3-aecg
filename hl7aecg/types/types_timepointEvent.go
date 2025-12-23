package types

// TimepointEvent represents the timepoint or study event during which ECG waveforms were collected.
//
// Commonly referred to as a "visit", but can also represent the "element" concept
// from CDISC's SDTM standard.
//
// SDTM Mapping:
//   - TimepointEvent maps to SDTM "visit" concept
//   - Code maps to VISITNUM or VISITDY
//
// XML Structure:
//
//	<timepointEvent>
//	  <code code="VISIT_2" codeSystem="2.16.840.1.113883.3.1"/>
//	  <effectiveTime>
//	    <low value="20020509100700"/>
//	    <high value="20020509134600"/>
//	  </effectiveTime>
//	  <reasonCode code="S" codeSystem=""/>
//	  <performer>...</performer>
//	  <componentOf>
//	    <subjectAssignment>...</subjectAssignment>
//	  </componentOf>
//	</timepointEvent>
//
// Cardinality: Required (in AnnotatedECG componentOf context)
// Reference: HL7 aECG Implementation Guide, Page 17-18
type TimepointEvent struct {
	// Code is the code naming this timepoint event.
	//
	// This could be a visit number, study day, etc. The set of timepoint event
	// codes is usually defined by the protocol. Therefore, the codeSystem UID
	// usually names the protocol.
	//
	// Examples:
	//   - "VISIT_1", "VISIT_2": Visit numbers
	//   - "V1", "V2": Short visit codes
	//   - "DAY_1", "DAY_2": Study days
	//   - "1", "2", "3": Simple numeric codes
	//
	// Example: code="VISIT_2" codeSystem="2.16.840.1.113883.3.1"
	//
	// XML Tag: <code code="..." codeSystem="..."/>
	// Cardinality: Optional
	Code *Code[string, CodeSystemOID] `xml:"code,omitempty"`

	// EffectiveTime is the time or range of time this timepoint event occurred.
	//
	// This typically spans the duration of the visit or study element.
	//
	// Example:
	//   <effectiveTime>
	//     <low value="20020509100700"/>   (Visit started 10:07 AM)
	//     <high value="20020509134600"/>  (Visit ended 1:46 PM)
	//   </effectiveTime>
	//
	// XML Tag: <effectiveTime>...</effectiveTime>
	// Cardinality: Optional
	EffectiveTime *EffectiveTime `xml:"effectiveTime,omitempty"`

	// ReasonCode indicates the reason for the event.
	//
	// Suggested Values (from CDISC lab model VisitType):
	//   - "S": Scheduled - planned visit according to protocol
	//   - "U": Unscheduled - visit not planned
	//
	// Note: As of 2005, no formal HL7 vocabulary existed for these codes,
	// so codeSystem is typically empty.
	//
	// Example: code="S" codeSystem=""
	//
	// XML Tag: <reasonCode code="..." codeSystem=""/>
	// Cardinality: Optional
	ReasonCode *Code[VisitTypeCode, string] `xml:"reasonCode,omitempty"`

	// Performer identifies the person primarily responsible for this timepoint event,
	// other than the trial site investigator.
	//
	// XML Tag: <performer>...</performer>
	// Cardinality: Optional
	Performer *TimepointEventPerformer `xml:"performer,omitempty"`

	// ComponentOf links this timepoint event to the subject assignment.
	//
	// XML Tag: <componentOf>...</componentOf>
	// Cardinality: Required
	ComponentOf ComponentOfSubjectAssignment `xml:"componentOf"`
}

// TimepointEventPerformer identifies the person primarily responsible for this timepoint event.
//
// This is typically a technician or healthcare worker, not the principal investigator
// (who is named elsewhere in the AnnotatedECG).
//
// XML Structure:
//
//	<performer>
//	  <studyEventPerformer>
//	    <id root="2.16.840.1.113883.3.5" extension="TECH_1"/>
//	    <assignedPerson>
//	      <name>Julie Tech</name>
//	    </assignedPerson>
//	  </studyEventPerformer>
//	</performer>
//
// Cardinality: Optional (within TimepointEvent)
// Reference: HL7 aECG Implementation Guide, Page 17
type TimepointEventPerformer struct {
	// StudyEventPerformer contains the identification and name of the performer.
	//
	// XML Tag: <studyEventPerformer>...</studyEventPerformer>
	// Cardinality: Required (within Performer)
	StudyEventPerformer StudyEventPerformer `xml:"studyEventPerformer"`
}

// StudyEventPerformer represents a person involved in conducting the study event.
//
// XML Structure:
//
//	<studyEventPerformer>
//	  <id root="2.16.840.1.113883.3.5" extension="TECH_1"/>
//	  <assignedPerson>
//	    <name>Julie Tech</name>
//	  </assignedPerson>
//	</studyEventPerformer>
//
// Cardinality: Required (within TimepointEventPerformer)
type StudyEventPerformer struct {
	// ID is the unique identifier for this person.
	//
	// Structure:
	//   - Root: Required, must be a UID (OID or UUID)
	//   - Extension: Optional, traditional person identifier
	//
	// Example: root="2.16.840.1.113883.3.5" extension="TECH_1"
	//
	// XML Tag: <id root="..." extension="..."/>
	// Cardinality: Optional
	ID *ID `xml:"id,omitempty"`

	// AssignedPerson contains the name of the performer.
	//
	// XML Tag: <assignedPerson>...</assignedPerson>
	// Cardinality: Optional
	AssignedPerson *AssignedPerson `xml:"assignedPerson,omitempty"`
}

// AssignedPerson represents a person assigned to perform a study event.
//
// XML Structure:
//
//	<assignedPerson>
//	  <name>Julie Tech</name>
//	</assignedPerson>
//
// Cardinality: Optional (within StudyEventPerformer)
type AssignedPerson struct {
	// Name is the name of this person.
	//
	// Can be a simple string (e.g., "Julie Tech") or structured PersonName.
	//
	// XML Tag: <name>...</name>
	// Cardinality: Optional
	Name *string `xml:"name,omitempty"`
}

// ComponentOfSubjectAssignment links the timepoint event to the subject assignment.
//
// XML Structure:
//
//	<componentOf>
//	  <subjectAssignment>...</subjectAssignment>
//	</componentOf>
//
// Cardinality: Required (within TimepointEvent)
type ComponentOfSubjectAssignment struct {
	// SubjectAssignment contains the subject-trial association.
	//
	// XML Tag: <subjectAssignment>...</subjectAssignment>
	// Cardinality: Required
	SubjectAssignment SubjectAssignment `xml:"subjectAssignment"`
}

// =============================================================================
// Relative Timepoint Types
// =============================================================================

// RelativeTimepoint identifies a timepoint relative to a reference event.
//
// For example, if the protocol specifies an ECG assessment 30 minutes after
// the first dosage, this identifies the "30 minutes post dosage" timepoint.
//
// SDTM Mapping:
//   - RelativeTimepoint maps to SDTM "planned time point" concept
//   - Code maps to EGTPT, EGTPTNUM, or EGELTM
//
// XML Structure:
//
//	<relativeTimepoint>
//	  <code code="POST_DOSAGE_30" codeSystem="16.840.1.113883.3.1"/>
//	  <componentOf>
//	    <pauseQuantity value="1800" unit="s"/>
//	    <protocolTimepointEvent>...</protocolTimepointEvent>
//	  </componentOf>
//	</relativeTimepoint>
//
// Cardinality: Required (within AnnotatedECG Definition)
// Reference: HL7 aECG Implementation Guide, Page 19-20
type RelativeTimepoint struct {
	// Code identifies this relative timepoint.
	//
	// This code is most likely defined by the protocol, so the codeSystem UID
	// should name the protocol.
	//
	// Common Examples:
	//   - "0": Baseline (at reference event)
	//   - "0.5": 30 minutes (0.5 hours) after reference event
	//   - "1": 1 hour after reference event
	//   - "POST_DOSAGE_30": 30 minutes post dosage
	//   - "PRE_MEAL": Before meal
	//
	// Example: code="POST_DOSAGE_30" codeSystem="16.840.1.113883.3.1"
	//
	// XML Tag: <code code="..." codeSystem="..."/>
	// Cardinality: Required
	Code Code[string, CodeSystemOID] `xml:"code"`

	// ComponentOf links this relative timepoint to the protocol timepoint event.
	//
	// XML Tag: <componentOf>...</componentOf>
	// Cardinality: Required
	ComponentOf RelativeTimepointComponentOf `xml:"componentOf"`
}

// RelativeTimepointComponentOf links the relative timepoint to protocol information.
//
// XML Structure:
//
//	<componentOf>
//	  <pauseQuantity value="1800" unit="s"/>
//	  <protocolTimepointEvent>...</protocolTimepointEvent>
//	</componentOf>
//
// Cardinality: Required (within RelativeTimepoint)
type RelativeTimepointComponentOf struct {
	// PauseQuantity is the time delay between the ReferenceEvent and RelativeTimepoint.
	//
	// For example, if the reference event is a dosage and the relative timepoint
	// is "30 minutes post dosage", then pauseQuantity = 30 minutes (1800 seconds).
	//
	// Common Units:
	//   - "s": Seconds
	//   - "min": Minutes
	//   - "h": Hours
	//
	// Example: value="1800" unit="s" (1800 seconds = 30 minutes)
	//
	// XML Tag: <pauseQuantity value="..." unit="..."/>
	// Cardinality: Optional
	PauseQuantity *Quantity `xml:"pauseQuantity,omitempty"`

	// ProtocolTimepointEvent identifies the timepoint event as defined in the protocol.
	//
	// XML Tag: <protocolTimepointEvent>...</protocolTimepointEvent>
	// Cardinality: Required
	ProtocolTimepointEvent ProtocolTimepointEvent `xml:"protocolTimepointEvent"`
}

// Quantity represents a physical quantity with value and unit.
//
// XML Structure:
//
//	<pauseQuantity value="1800" unit="s"/>
//
// Cardinality: Optional (in various contexts)
type Quantity struct {
	// Value is the numeric value.
	//
	// Example: "1800" (for 1800 seconds)
	//
	// XML Tag: value="..."
	// Cardinality: Required
	Value string `xml:"value,attr"`

	// Unit is the unit of measure.
	//
	// Examples: "s" (seconds), "min" (minutes), "h" (hours)
	//
	// XML Tag: unit="..."
	// Cardinality: Required
	Unit string `xml:"unit,attr"`
}

// ProtocolTimepointEvent identifies the timepoint event as defined in the protocol.
//
// This is almost always the same as the TimepointEvent (the actual timepoint
// of the ECG assessment). However, if the timepoint event is not per protocol,
// this names the intended protocol timepoint while TimepointEvent names the
// actual timepoint.
//
// XML Structure:
//
//	<protocolTimepointEvent>
//	  <code code="VISIT_2" codeSystem="16.840.1.113883.3.1"/>
//	  <component>
//	    <referenceEvent>...</referenceEvent>
//	  </component>
//	</protocolTimepointEvent>
//
// Cardinality: Required (within RelativeTimepointComponentOf)
// Reference: HL7 aECG Implementation Guide, Page 19-20
type ProtocolTimepointEvent struct {
	// Code names the timepoint event as defined in the protocol.
	//
	// This could be a visit number, study day, etc. The codeSystem UID
	// usually names the protocol defining this code.
	//
	// Example: code="VISIT_2" codeSystem="16.840.1.113883.3.1"
	//
	// XML Tag: <code code="..." codeSystem="..."/>
	// Cardinality: Required
	Code Code[string, CodeSystemOID] `xml:"code"`

	// Component contains the reference event for this timepoint.
	//
	// XML Tag: <component>...</component>
	// Cardinality: Optional
	Component *ProtocolTimepointEventComponent `xml:"component,omitempty"`
}

// ProtocolTimepointEventComponent contains the reference event.
//
// XML Structure:
//
//	<component>
//	  <referenceEvent>...</referenceEvent>
//	</component>
//
// Cardinality: Optional (within ProtocolTimepointEvent)
type ProtocolTimepointEventComponent struct {
	// ReferenceEvent identifies the benchmark event.
	//
	// XML Tag: <referenceEvent>...</referenceEvent>
	// Cardinality: Required (within Component)
	ReferenceEvent ReferenceEvent `xml:"referenceEvent"`
}

// ReferenceEvent identifies the reference or benchmark event for the ECG assessment
// within the timepoint event.
//
// This could be a dosing, meal, exercise, etc. that serves as the reference point
// for relative timepoints.
//
// SDTM Mapping:
//   - ReferenceEvent maps to SDTM "time point reference" concept
//   - Code maps to EGTPTREF
//
// XML Structure:
//
//	<referenceEvent>
//	  <code code="DOSAGE_1" codeSystem="16.840.1.113883.3.1"/>
//	</referenceEvent>
//
// Cardinality: Required (within ProtocolTimepointEventComponent)
// Reference: HL7 aECG Implementation Guide, Page 19-20
type ReferenceEvent struct {
	// Code names the reference event.
	//
	// This could be a dosing, meal, exercise, etc. The reference event codes
	// are most likely defined by the protocol, so the codeSystem UID usually
	// names the protocol.
	//
	// Common Examples:
	//   - "DOSAGE_1": First dosage
	//   - "DOSAGE_2": Second dosage
	//   - "MEAL": Meal time
	//   - "EXERCISE": Exercise event
	//   - "BASELINE": Baseline measurement
	//
	// Example: code="DOSAGE_1" codeSystem="16.840.1.113883.3.1"
	//
	// XML Tag: <code code="..." codeSystem="..."/>
	// Cardinality: Required
	Code Code[string, CodeSystemOID] `xml:"code"`
}
