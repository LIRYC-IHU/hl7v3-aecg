package types

import "fmt"

// ClinicalTrial represents the clinicalTrial element in an HL7 aECG document.
// It contains identifying information about the clinical trial used in the study.
//
// According to HL7 aECG Implementation Guide:
//   - This element is required within the SubjectAssignment context
//   - It uniquely identifies the clinical trial using OID-based identifiers
//   - The trial identifier should remain consistent across all subjects in the trial
//
// XML Structure:
//
//	<clinicalTrial>
//	  <id root="2.16.840.1.113883.3.4" extension="PUK-123-TRL-1"/>
//	  <title>Cardiac Safety Trial 1 for Compound PUK-123</title>
//	  <activityTime>
//	    <low value="20010509"/>
//	    <high value="20020316"/>
//	  </activityTime>
//	</clinicalTrial>
//
// Cardinality: Required
// Reference: HL7 aECG Implementation Guide, Final 21-March-2005, Page 10-11
type ClinicalTrial struct {
	// ID is the unique identifier for the clinical trial.
	//
	// Structure:
	//   - Root: Required, must be a UID (OID or UUID)
	//   - Extension: Optional, may be any string value
	//
	// Best Practice:
	//   Sponsors should assign a unique OID to every clinical trial.
	//   The OID goes into the root part, and the traditional trial
	//   identifier (e.g., "PUK-123-TRL-1") goes into the extension.
	//
	// The combination of root and extension must be universally unique
	// across all trials globally.
	//
	// Example:
	//   root="2.16.840.1.113883.3.4" extension="PUK-123-TRL-1"
	//
	// XML Tag: <id root="..." extension="..."/>
	// Cardinality: Required
	ID ID `xml:"id"`

	// Title is the human-readable name of the clinical trial.
	//
	// This field provides a descriptive name for the trial that is
	// easier for humans to understand than the technical identifier.
	//
	// Example:
	//   "Cardiac Safety Trial 1 for Compound PUK-123"
	//
	// XML Tag: <title>...</title>
	// Cardinality: Optional
	Title *string `xml:"title,omitempty"`

	// ActivityTime represents the date range during which the trial took place.
	//
	// This time range helps identify the particular trial instance,
	// similar to how a birth date helps identify a particular patient.
	//
	// Usage Guidelines:
	//   - If start date is known but end date unknown: include only Low
	//   - If end date is known but start date unknown: include only High
	//   - If both dates are known: include both Low and High
	//
	// Note: This is an "activity time" representing when the trial
	// administratively took place, not the "effective time" of
	// individual measurements or ECG recordings.
	//
	// Time values should follow HL7 timestamp format:
	//   - YYYYMMDD (date only)
	//   - YYYYMMDDHHMMSS (date and time)
	//
	// Example:
	//   <activityTime>
	//     <low value="20010509"/>
	//     <high value="20020316"/>
	//   </activityTime>
	//
	// XML Tag: <activityTime>...</activityTime>
	// Cardinality: Optional
	ActivityTime EffectiveTime `xml:"activityTime"`

	// Location identifies the trial site location where ECG waveforms were acquired.
	//
	// This provides information about the physical location and facility
	// where the subject's data was collected.
	//
	// XML Tag: <location>...</location>
	// Cardinality: Optional
	Location *Location `xml:"location,omitempty"`
}

// GetIdentifier returns a human-readable string representation of the trial identifier.
//
// If both root and extension are present, returns "root:extension".
// If only root is present, returns just the root.
//
// Example outputs:
//   - "2.16.840.1.113883.3.4:PUK-123-TRL-1"
//   - "2.16.840.1.113883.3.4"
func (ct *ClinicalTrial) GetIdentifier() string {
	if ct.ID.Extension != "" {
		return fmt.Sprintf("%s:%s", ct.ID.Root, ct.ID.Extension)
	}
	return ct.ID.Root
}

// GetDisplayName returns a display name for the trial.
//
// If Title is set, returns the title.
// Otherwise, returns the identifier string.
func (ct *ClinicalTrial) GetDisplayName() string {
	if ct.Title != nil && *ct.Title != "" {
		return *ct.Title
	}
	return ct.GetIdentifier()
}
