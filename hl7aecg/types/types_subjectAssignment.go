package types

// SubjectAssignment represents the act of associating a subject to a clinical trial.
//
// This is a required element that links the subject, trial, and timepoint together.
// It may optionally include treatment group assignment information.
//
// XML Structure:
//
//	<subjectAssignment>
//	  <subject>...</subject>
//	  <definition>
//	    <treatmentGroupAssignment>
//	      <code code="GRP_003" codeSystem="2.16.840.1.113883.3.1"/>
//	    </treatmentGroupAssignment>
//	  </definition>
//	  <componentOf>
//	    <clinicalTrial>...</clinicalTrial>
//	  </componentOf>
//	</subjectAssignment>
//
// Cardinality: Required (in TimepointEvent context)
// Reference: HL7 aECG Implementation Guide, Page 16
type SubjectAssignment struct {
	// Subject identifies the subject from which the ECG waveforms were obtained.
	//
	// XML Tag: <subject>...</subject>
	// Cardinality: Required
	Subject Subject `xml:"subject"`

	// Definition defines the subject's association with the trial by naming
	// the treatment group they are assigned to.
	//
	// XML Tag: <definition>...</definition>
	// Cardinality: Optional
	Definition *SubjectAssignmentDefinition `xml:"definition,omitempty"`

	// ComponentOf links this subject assignment to the clinical trial.
	//
	// XML Tag: <componentOf>...</componentOf>
	// Cardinality: Required
	ComponentOf ComponentOfClinicalTrial `xml:"componentOf"`
}

// SubjectAssignmentDefinition defines the subject's association with a trial
// by naming the treatment group they are assigned to.
//
// XML Structure:
//
//	<definition>
//	  <treatmentGroupAssignment>
//	    <code code="GRP_003" codeSystem="2.16.840.1.113883.3.1"/>
//	  </treatmentGroupAssignment>
//	</definition>
//
// Cardinality: Optional (within SubjectAssignment)
type SubjectAssignmentDefinition struct {
	// TreatmentGroupAssignment identifies the treatment group the subject belongs to.
	//
	// XML Tag: <treatmentGroupAssignment>...</treatmentGroupAssignment>
	// Cardinality: Required (within Definition)
	TreatmentGroupAssignment TreatmentGroupAssignment `xml:"treatmentGroupAssignment"`
}

// TreatmentGroupAssignment identifies a group of subjects that went through
// the trial in the same way.
//
// All subjects in a treatment group receive the same drug dosages in the same order.
// This concept is similar to the CDISC SDTM "arm" concept.
//
// XML Structure:
//
//	<treatmentGroupAssignment>
//	  <code code="GRP_003" codeSystem="2.16.840.1.113883.3.1"/>
//	</treatmentGroupAssignment>
//
// Cardinality: Required (within SubjectAssignmentDefinition)
// Reference: HL7 aECG Implementation Guide, Page 16
type TreatmentGroupAssignment struct {
	// Code is the protocol or trial-specific code identifying the treatment group.
	//
	// Code System Rules:
	//   - If it's a trial-specific code: codeSystem = trial's OID
	//     Example: Trial OID "2.16.840.1.113883.3.5" -> use as codeSystem
	//   - If it's a protocol-specific code: codeSystem = protocol's OID
	//     Example: Protocol OID "2.16.840.1.113883.3.1" -> use as codeSystem
	//
	// Example: code="GRP_003" codeSystem="2.16.840.1.113883.3.1"
	//
	// Common Treatment Group Codes (sponsor-defined):
	//   - "PLACEBO": Placebo group
	//   - "TREATMENT_A": First treatment group
	//   - "TREATMENT_B": Second treatment group
	//   - "GRP_001", "GRP_002": Numbered groups
	//
	// XML Tag: <code code="..." codeSystem="..."/>
	// Cardinality: Required
	Code Code[TreatmentGroupCode, CodeSystemOID] `xml:"code"`
}

// ComponentOfClinicalTrial links the subject assignment to the clinical trial.
//
// XML Structure:
//
//	<componentOf>
//	  <clinicalTrial>...</clinicalTrial>
//	</componentOf>
//
// Cardinality: Required (within SubjectAssignment)
type ComponentOfClinicalTrial struct {
	// ClinicalTrial contains the trial identification and metadata.
	//
	// XML Tag: <clinicalTrial>...</clinicalTrial>
	// Cardinality: Required
	ClinicalTrial ClinicalTrial `xml:"clinicalTrial"`
}
