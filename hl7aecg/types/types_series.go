package types

// =============================================================================
// Series Types
// =============================================================================

// Series contains all sequences, regions of interest, and annotations sharing
// a common frame of reference.
//
// All physical quantities having the same code must be comparable. For example,
// if the series contains 2 different voltage sequences for lead II, the voltages
// must come from the same set of electrodes. If electrodes were changed or
// additional filtering was used, those sequences must appear in different series.
//
// Typically, a series contains all waveforms and annotations for a single ECG.
// If multiple ECGs are in a single aECG file, each uses a different series.
//
// A series can be derived from another series. For example:
//   - Representative beat waveforms derived from rhythm series
//   - Waveforms with special filtering derived from "raw" rhythm waveforms
//
// XML Structure:
//
//	<component>
//	  <series>
//	    <id root="b65deea0-078e-11d9-9669-0800200c9a66"/>
//	    <code code="RHYTHM" codeSystem="2.16.840.1.113883.5.4"/>
//	    <effectiveTime>
//	      <low value="20021122091000"/>
//	      <high value="20021122091010"/>
//	    </effectiveTime>
//	    <author>...</author>
//	    <component>
//	      <sequenceSet>...</sequenceSet>
//	    </component>
//	  </series>
//	</component>
//
// Cardinality: Optional (in AnnotatedECG component context)
// Reference: HL7 aECG Implementation Guide, Page 26-27
type Series struct {
	// ID is the unique identifier of the series.
	//
	// This is typically assigned by the ECG waveform collection device
	// or the ECG management system that exported the aECG.
	//
	// Format: May be OID or UUID, but typically UUID
	//
	// Example: root="b65deea0-078e-11d9-9669-0800200c9a66"
	//
	// XML Tag: <id root="..."/>
	// Cardinality: Optional
	ID *ID `xml:"id,omitempty"`

	// Code identifies the type of series.
	//
	// Vocabulary: HL7 ActCode
	// Code System: 2.16.840.1.113883.5.4
	//
	// Defined Codes:
	//   - "RHYTHM": Series contains rhythm waveforms collected by the device.
	//     Voltage samples are related to each other in real time (wall time).
	//   - "REPRESENTATIVE_BEAT": Series contains waveforms of a representative
	//     beat derived from rhythm waveforms. Voltage samples are related in
	//     time relative to the beginning of cardiac cycle, not real time.
	//
	// Example: code="RHYTHM" codeSystem="2.16.840.1.113883.5.4"
	//
	// XML Tag: <code code="..." codeSystem="..."/>
	// Cardinality: Required
	Code *Code[SeriesTypeCode, CodeSystemOID] `xml:"code"`

	// EffectiveTime is the physiologically relevant time range for the ECG waveforms.
	//
	// This is typically the "acquisition time" determined by the device that
	// collected the waveforms.
	//
	// Usage:
	//   - If device recorded only start of acquisition: use Low
	//   - If device recorded only end of acquisition: use High
	//   - If both start and end are known: use both Low and High
	//
	// Example: 10-second ECG from 09:10:00 to 09:10:10
	//   <effectiveTime>
	//     <low value="20021122091000"/>
	//     <high value="20021122091010"/>
	//   </effectiveTime>
	//
	// XML Tag: <effectiveTime>...</effectiveTime>
	// Cardinality: Required
	EffectiveTime EffectiveTime `xml:"effectiveTime"`

	// Author describes the device that "authored" (recorded) the series waveforms.
	//
	// This typically describes an electrocardiograph or Holter recorder.
	//
	// XML Tag: <author>...</author>
	// Cardinality: Optional
	Author *Author `xml:"author,omitempty"`

	// SecondaryPerformer describes technician(s) operating the device.
	//
	// May include multiple performers for different functions:
	//   - Holter hookup technician
	//   - Waveform analysis technician
	//   - ECG recording technician
	//
	// XML Tag: <secondaryPerformer>...</secondaryPerformer>
	// Cardinality: Optional (0..*)
	SecondaryPerformer []SecondaryPerformer `xml:"secondaryPerformer,omitempty"`

	// Support defines the region of interest if this series is derived from another.
	//
	// Used when a series is derived from part of a parent series.
	// Example: Representative beats derived from 30-second segments of rhythm waveforms
	//
	// XML Tag: <support>...</support>
	// Cardinality: Optional
	Support *SeriesSupport `xml:"support,omitempty"`

	// ControlVariable captures related information about the subject or ECG collection conditions.
	//
	// Control variables can capture information such as:
	//   - Subject's age (independent of birth date recorded elsewhere)
	//   - Fasting status
	//   - Other clinical information relevant to ECG interpretation
	//
	// XML Tag: <controlVariable>...</controlVariable>
	// Cardinality: Optional (0..*)
	ControlVariable []ControlVariable `xml:"controlVariable,omitempty"`

	// Component contains the sequence sets with actual waveform data.
	//
	// XML Tag: <component>...</component>
	// Cardinality: Required (1..*)
	Component []SeriesComponent `xml:"component"`
}

func NewSeries() *Series {
	return &Series{
		ID:   &ID{},
		Code: &Code[SeriesTypeCode, CodeSystemOID]{},
	}
}

// Author wraps the SeriesAuthor to provide the correct XML structure.
//
// XML Structure:
//
//	<author>
//	  <seriesAuthor>...</seriesAuthor>
//	</author>
//
// Cardinality: Optional (within Series)
type Author struct {
	// SeriesAuthor contains the device and manufacturer information.
	//
	// XML Tag: <seriesAuthor>...</seriesAuthor>
	// Cardinality: Required (within Author)
	SeriesAuthor SeriesAuthor `xml:"seriesAuthor"`
}

// SeriesAuthor describes the device that authored (recorded) the series waveforms.
//
// This typically describes an electrocardiograph or Holter recorder used to
// capture the ECG data.
//
// XML Structure:
//
//	<author>
//	  <seriesAuthor>
//	    <id root="2.16.840.1.113883.3.5" extension="45"/>
//	    <manufacturedSeriesDevice>
//	      <id root="1.3.6.1.4.1.57054" extension="SN234-AR9-102993"/>
//	      <code code="12LEAD_ELECTROCARDIOGRAPH" codeSystem=""/>
//	      <manufacturerModelName>Electrograph 250</manufacturerModelName>
//	      <softwareName>Rx 5.3</softwareName>
//	    </manufacturedSeriesDevice>
//	    <manufacturerOrganization>
//	      <id root="1.3.6.1.4.1.57054"/>
//	      <n>ECG Devices By Smith, Inc.</n>
//	    </manufacturerOrganization>
//	  </seriesAuthor>
//	</author>
//
// Cardinality: Optional (within Series)
// Reference: HL7 aECG Implementation Guide, Page 28-29
type SeriesAuthor struct {
	// ID is the unique identifier of the device in its role as waveform author
	// in this clinical trial.
	//
	// This could be assigned by the investigator, healthcare organization,
	// CRO, or trial sponsor.
	//
	// Best Practice: Put traditional identifier in extension
	//
	// Example: root="2.16.840.1.113883.3.5" extension="45"
	//
	// XML Tag: <id root="..." extension="..."/>
	// Cardinality: Optional
	ID *ID `xml:"id,omitempty"`

	// ManufacturedSeriesDevice contains device identification and specifications.
	//
	// XML Tag: <manufacturedSeriesDevice>...</manufacturedSeriesDevice>
	// Cardinality: Required (within SeriesAuthor)
	ManufacturedSeriesDevice ManufacturedSeriesDevice `xml:"manufacturedSeriesDevice"`

	// ManufacturerOrganization identifies the organization that manufactured the device.
	//
	// XML Tag: <manufacturerOrganization>...</manufacturerOrganization>
	// Cardinality: Optional
	ManufacturerOrganization *ManufacturerOrganization `xml:"manufacturerOrganization,omitempty"`
}

// ManufacturedSeriesDevice represents the ECG device specifications.
//
// Contains device identification, type, model, and software information.
//
// XML Structure:
//
//	<manufacturedSeriesDevice>
//	  <id root="1.3.6.1.4.1.57054" extension="SN234-AR9-102993"/>
//	  <code code="12LEAD_ELECTROCARDIOGRAPH" codeSystem=""/>
//	  <manufacturerModelName>Electrograph 250</manufacturerModelName>
//	  <softwareName>Rx 5.3</softwareName>
//	</manufacturedSeriesDevice>
//
// Cardinality: Required (within SeriesAuthor)
// Reference: HL7 aECG Implementation Guide, Page 28-29
type ManufacturedSeriesDevice struct {
	// ID is the unique identifier of the device, independent of its role.
	//
	// This is typically the serial number assigned by the device manufacturer.
	//
	// Best Practice:
	//   - Root: OID identifying the manufacturing organization
	//   - Extension: Serial number
	//
	// Example: root="1.3.6.1.4.1.57054" extension="SN234-AR9-102993"
	//
	// XML Tag: <id root="..." extension="..."/>
	// Cardinality: Optional
	ID *ID `xml:"id,omitempty"`

	// Code identifies the type of device.
	//
	// Note: As of 2005, no formal vocabulary existed for device types.
	//
	// Suggested Codes:
	//   - "12LEAD_ELECTROCARDIOGRAPH": Standard 12-lead ECG device
	//   - "12LEAD_HOLTER": Holter monitor device
	//   - "3LEAD_MONITOR": 3-lead continuous monitor
	//
	// Example: code="12LEAD_ELECTROCARDIOGRAPH" codeSystem=""
	//
	// XML Tag: <code code="..." codeSystem=""/>
	// Cardinality: Optional
	Code *Code[DeviceTypeCode, CodeSystemOID] `xml:"code,omitempty"`

	// ManufacturerModelName is the model name of the device.
	//
	// Example: "Electrograph 250", "CardioMax Pro 3000"
	//
	// XML Tag: <manufacturerModelName>...</manufacturerModelName>
	// Cardinality: Optional
	ManufacturerModelName *string `xml:"manufacturerModelName,omitempty"`

	// SerialNumber is the serial number of the device.
	// Example: "SN234-AR9-102993"
	//
	// XML Tag: <SerialNumber>...</SerialNumber>
	// Cardinality: Optional
	SerialNumber *string `xml:"SerialNumber,omitempty"`

	// SoftwareName is the name and/or version of the software in the device.
	//
	// Example: "Rx 5.3", "CardioSoft v2.1.4"
	//
	// XML Tag: <softwareName>...</softwareName>
	// Cardinality: Optional
	SoftwareName *string `xml:"softwareName,omitempty"`
}

// ManufacturerOrganization represents the organization that manufactured the device.
//
// XML Structure:
//
//	<manufacturerOrganization>
//	  <id root="1.3.6.1.4.1.57054"/>
//	  <n>ECG Devices By Smith, Inc.</n>
//	</manufacturerOrganization>
//
// Cardinality: Optional (within SeriesAuthor)
type ManufacturerOrganization struct {
	// ID is the unique identifier of the manufacturing organization.
	//
	// This is typically in the form of an OID.
	//
	// Example: root="1.3.6.1.4.1.57054"
	//
	// XML Tag: <id root="..."/>
	// Cardinality: Optional
	ID *ID `xml:"id,omitempty"`

	// Name is the name of the manufacturing organization.
	//
	// Example: "ECG Devices By Smith, Inc.", "PhilipsHealthcare"
	//
	// XML Tag: <n>...</n>
	// Cardinality: Optional
	Name *string `xml:"name,omitempty"`
}

// SecondaryPerformer describes a technician operating the device that captured
// the ECG waveforms.
//
// Multiple performers may be involved in ECG acquisition:
//   - Holter hookup technician
//   - Holter analyst who produces summary report
//   - Electrocardiograph technician
//
// XML Structure:
//
//	<secondaryPerformer>
//	  <functionCode code="ELECTROCARDIOGRAPH_TECH" codeSystem=""/>
//	  <time>
//	    <low value="20021122091000"/>
//	    <high value="20021122091010"/>
//	  </time>
//	  <seriesPerformer>
//	    <id root="2.16.840.1.113883.3.4" extension="TECH-221"/>
//	    <assignedPerson>
//	      <n>KAB</n>
//	    </assignedPerson>
//	  </seriesPerformer>
//	</secondaryPerformer>
//
// Cardinality: Optional (within Series, 0..*)
// Reference: HL7 aECG Implementation Guide, Page 30
type SecondaryPerformer struct {
	// FunctionCode describes the function the technician was performing.
	//
	// Note: As of 2005, no formal vocabulary existed for performer functions.
	//
	// Suggested Codes:
	//   - "HOLTER_HOOKUP_TECH": Holter recorder hookup technician
	//   - "HOLTER_ANALYST": Holter waveform analysis technician
	//   - "ELECTROCARDIOGRAPH_TECH": Standard ECG technician
	//
	// Example: code="ELECTROCARDIOGRAPH_TECH" codeSystem=""
	//
	// XML Tag: <functionCode code="..." codeSystem=""/>
	// Cardinality: Optional
	FunctionCode *Code[PerformerFunctionCode, CodeSystemOID] `xml:"functionCode,omitempty"`

	// Time is the time or period during which the technician performed
	// the indicated function for this series.
	//
	// Example: Same as series effective time for standard ECG
	//
	// XML Tag: <time>...</time>
	// Cardinality: Optional
	Time *EffectiveTime `xml:"time,omitempty"`

	// SeriesPerformer contains the performer identification and name.
	//
	// XML Tag: <seriesPerformer>...</seriesPerformer>
	// Cardinality: Required (within SecondaryPerformer)
	SeriesPerformer SeriesPerformer `xml:"seriesPerformer"`
}

// SeriesPerformer represents a technician who performed ECG acquisition.
//
// XML Structure:
//
//	<seriesPerformer>
//	  <id root="2.16.840.1.113883.3.4" extension="TECH-221"/>
//	  <assignedPerson>
//	    <n>KAB</n>
//	  </assignedPerson>
//	</seriesPerformer>
//
// Cardinality: Required (within SecondaryPerformer)
type SeriesPerformer struct {
	// ID is the role-specific unique identifier of the secondary performer.
	//
	// This is the identifier assigned to this technician in their role within
	// the trial. Could be assigned by investigator, employer, CRO, or sponsor.
	//
	// Example: root="2.16.840.1.113883.3.4" extension="TECH-221"
	//
	// XML Tag: <id root="..." extension="..."/>
	// Cardinality: Optional
	ID *ID `xml:"id,omitempty"`

	// AssignedPerson contains the name of the technician.
	//
	// XML Tag: <assignedPerson>...</assignedPerson>
	// Cardinality: Optional
	AssignedPerson *SeriesAssignedPerson `xml:"assignedPerson,omitempty"`
}

// SeriesAssignedPerson represents a person assigned to perform ECG tasks.
//
// XML Structure:
//
//	<assignedPerson>
//	  <n>KAB</n>
//	</assignedPerson>
//
// Cardinality: Optional (within SeriesPerformer)
type SeriesAssignedPerson struct {
	// Name is the name of the technician.
	//
	// Can be as simple as initials or include full name components
	// (first, middle, last, prefix, suffix, etc.).
	//
	// Example: "KAB" (initials), "Karen A. Brown"
	//
	// XML Tag: <n>...</n>
	// Cardinality: Optional
	Name *string `xml:"name,omitempty"`
}

// SeriesSupport defines the region of interest if this series is derived
// from another series.
//
// A series can be derived by applying an algorithm to transform waveform data.
// If derived from part of a parent series, the supporting ROI specifies which part.
//
// Example: Representative beats derived from 30-second segments of 3-minute
// rhythm waveforms. Each derived series has a supporting ROI specifying which
// 30-second segment it came from.
//
// XML Structure:
//
//	<support>
//	  <supportingROI>...</supportingROI>
//	</support>
//
// Cardinality: Optional (within Series)
// Reference: HL7 aECG Implementation Guide, Page 31-32
type SeriesSupport struct {
	// SupportingROI identifies the region of interest from the parent series.
	//
	// XML Tag: <supportingROI>...</supportingROI>
	// Cardinality: Required (within Support)
	SupportingROI SupportingROI `xml:"supportingROI"`
}

// SupportingROI represents a region of interest in a parent series.
//
// Defines which part of a parent series was used to derive this series.
//
// ROI can be:
//   - Fully specified: All boundaries defined, unnamed sequences not included
//   - Partially specified: All sequences included by default unless boundary specifies otherwise
//
// Cardinality: Required (within SeriesSupport)
type SupportingROI struct {
	// Code specifies if the ROI is fully or partially specified.
	//
	// XML Tag: <code code="..." codeSystem="..."/>
	// Cardinality: Required
	Code Code[RegionOfInterestType, CodeSystemOID] `xml:"code"`

	// Component contains boundary definitions for the ROI.
	//
	// XML Tag: <component>...</component>
	// Cardinality: Optional (0..*)
	Component []ROIComponent `xml:"component,omitempty"`
}

// ROIComponent contains a boundary definition for a region of interest.
//
// Cardinality: Optional (within SupportingROI, 0..*)
type ROIComponent struct {
	// Boundary defines the limits of the region of interest.
	//
	// XML Tag: <boundary>...</boundary>
	// Cardinality: Required (within Component)
	Boundary Boundary `xml:"boundary"`
}

// Boundary defines the limits of a region of interest.
//
// Used to specify which portion of sequences are included in the ROI.
// For example, time boundary for first 30 seconds, or specific leads.
//
// Cardinality: Required (within ROIComponent)
type Boundary struct {
	// Code identifies which sequence dimension this boundary applies to.
	//
	// Examples:
	//   - "TIME_ABSOLUTE": Time boundary
	//   - "MDC_ECG_LEAD_II": Lead II boundary
	//
	// XML Tag: <code code="..." codeSystem="..."/>
	// Cardinality: Required
	Code Code[TimeSequenceCode, CodeSystemOID] `xml:"code"`

	// Value defines the range of the boundary.
	//
	// Can be an interval (IVL_INT) or other value type.
	//
	// XML Tag: <value>...</value>
	// Cardinality: Optional
	Value *BoundaryValue `xml:"value,omitempty"`
}

// BoundaryValue represents a range value for a boundary.
//
// Typically an interval of integers (IVL_INT) defining start and end indices.
//
// Cardinality: Optional (within Boundary)
type BoundaryValue struct {
	// Low is the starting index of the interval.
	//
	// XML Tag: <low value="..."/>
	// Cardinality: Optional
	Low *int `xml:"low,attr,omitempty"`

	// High is the ending index of the interval.
	//
	// XML Tag: <high value="..."/>
	// Cardinality: Optional
	High *int `xml:"high,attr,omitempty"`
}

// SeriesComponent contains a sequence set with waveform data.
//
// XML Structure:
//
//	<component>
//	  <sequenceSet>...</sequenceSet>
//	</component>
//
// Cardinality: Required (within Series, 1..*)
type SeriesComponent struct {
	// SequenceSet contains related sequences sharing the same length.
	//
	// XML Tag: <sequenceSet>...</sequenceSet>
	// Cardinality: Required
	SequenceSet SequenceSet `xml:"sequenceSet"`
}
