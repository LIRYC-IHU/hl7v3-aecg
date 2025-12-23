package types

import "slices"

// =============================================================================
// HL7 Code Systems OIDs
// =============================================================================
// Reference: https://terminology.hl7.org/external_terminologies.html

type CodeSystemOID string

const (
	// Primary Code Systems for aECG
	CPT_OID                          CodeSystemOID = "2.16.840.1.113883.6.12"  // http://www.ama-assn.org/go/cpt
	MDC_OID                          CodeSystemOID = "2.16.840.1.113883.6.24"  // urn:iso:std:iso:11073:10101
	HL7_ActCode_OID                  CodeSystemOID = "2.16.840.1.113883.5.4"   // http://terminology.hl7.org/CodeSystem/v3-ActCode
	HL7_ActAdministrativeGender_OID  CodeSystemOID = "2.16.840.1.113883.5.1"   // http://terminology.hl7.org/CodeSystem/v3-AdministrativeGender
	HL7_Race_OID                     CodeSystemOID = "2.16.840.1.113883.5.104" // http://terminology.hl7.org/CodeSystem/v3-Race
	HL7_ResearchSubjectRoleBasis_OID CodeSystemOID = "2.16.840.1.113883.5.111" // http://terminology.hl7.org/CodeSystem/v3-ResearchSubjectRoleBasis

	// Additional Code Systems
	LOINC_OID     CodeSystemOID = "2.16.840.1.113883.6.1"  // http://loinc.org
	SNOMED_CT_OID CodeSystemOID = "2.16.840.1.113883.6.96" // http://snomed.info/sct
	UCUM_OID      CodeSystemOID = "2.16.840.1.113883.6.8"  // http://unitsofmeasure.org
)

// =============================================================================
// CPT Codes for ECG Procedures
// =============================================================================

type CPT_CODE string

const (
	CPT_CODE_ECG_Routine    CPT_CODE = "93000" // Electrocardiogram, routine ECG with at least 12 leads
	CPT_CODE_ECG_Tracing    CPT_CODE = "93005" // Electrocardiogram, tracing only, without interpretation and report
	CPT_CODE_Interpretation CPT_CODE = "93010" // Electrocardiogram, interpretation and report only
)

// =============================================================================
// HL7 ActCode Values
// =============================================================================

type SeriesTypeCode string

type TimeSequenceCode string

const (
	// Series Types
	RHYTHM_CODE              SeriesTypeCode = "RHYTHM"              // Rhythm waveform
	REPRESENTATIVE_BEAT_CODE SeriesTypeCode = "REPRESENTATIVE_BEAT" // Representative beat waveform
	MEDIAN_BEAT_CODE         SeriesTypeCode = "MEDIAN_BEAT"         // Median beat waveform

	// Time Sequence Types
	TIME_ABSOLUTE_CODE TimeSequenceCode = "TIME_ABSOLUTE" // Absolute time reference
	TIME_RELATIVE_CODE TimeSequenceCode = "TIME_RELATIVE" // Relative time reference
)

// =============================================================================
// MDC Lead Codes (Standard 12-Lead ECG)
// =============================================================================

type LeadCode string

const (
	// Limb Leads
	MDC_ECG_LEAD_I   LeadCode = "MDC_ECG_LEAD_I"   // Lead I
	MDC_ECG_LEAD_II  LeadCode = "MDC_ECG_LEAD_II"  // Lead II
	MDC_ECG_LEAD_III LeadCode = "MDC_ECG_LEAD_III" // Lead III

	// Augmented Leads
	MDC_ECG_LEAD_AVR LeadCode = "MDC_ECG_LEAD_AVR" // Lead AVR (Augmented Vector Right)
	MDC_ECG_LEAD_AVL LeadCode = "MDC_ECG_LEAD_AVL" // Lead AVL (Augmented Vector Left)
	MDC_ECG_LEAD_AVF LeadCode = "MDC_ECG_LEAD_AVF" // Lead AVF (Augmented Vector Foot)

	// Precordial Leads
	MDC_ECG_LEAD_V1 LeadCode = "MDC_ECG_LEAD_V1" // Lead V1
	MDC_ECG_LEAD_V2 LeadCode = "MDC_ECG_LEAD_V2" // Lead V2
	MDC_ECG_LEAD_V3 LeadCode = "MDC_ECG_LEAD_V3" // Lead V3
	MDC_ECG_LEAD_V4 LeadCode = "MDC_ECG_LEAD_V4" // Lead V4
	MDC_ECG_LEAD_V5 LeadCode = "MDC_ECG_LEAD_V5" // Lead V5
	MDC_ECG_LEAD_V6 LeadCode = "MDC_ECG_LEAD_V6" // Lead V6
)

// =============================================================================
// MDC Annotation Codes
// =============================================================================

type WaveformAnnotationCode string
type MeasurementCode string
type IntervalCode string

const (
	// Waveform Components
	MDC_ECG_WAVC_PWAVE   WaveformAnnotationCode = "MDC_ECG_WAVC_PWAVE"   // P wave
	MDC_ECG_WAVC_QRSWAVE WaveformAnnotationCode = "MDC_ECG_WAVC_QRSWAVE" // QRS complex
	MDC_ECG_WAVC_TWAVE   WaveformAnnotationCode = "MDC_ECG_WAVC_TWAVE"   // T wave
	MDC_ECG_WAVC_UWAVE   WaveformAnnotationCode = "MDC_ECG_WAVC_UWAVE"   // U wave

	// Time Intervals
	MDC_ECG_TIME_PD_QT  IntervalCode = "MDC_ECG_TIME_PD_QT"  // QT interval (ms)
	MDC_ECG_TIME_PD_QTC IntervalCode = "MDC_ECG_TIME_PD_QTC" // QT interval corrected (ms)
	MDC_ECG_TIME_PD_RR  IntervalCode = "MDC_ECG_TIME_PD_RR"  // RR interval (ms)
	MDC_ECG_TIME_PD_PP  IntervalCode = "MDC_ECG_TIME_PD_PP"  // PP interval (ms)
	MDC_ECG_TIME_PD_PR  IntervalCode = "MDC_ECG_TIME_PD_PR"  // PR interval (ms)
	MDC_ECG_TIME_PD_QRS IntervalCode = "MDC_ECG_TIME_PD_QRS" // QRS duration (ms)

	// Measurements
	MDC_ECG_HEART_RATE MeasurementCode = "MDC_ECG_HEART_RATE" // Heart rate (bpm)
	MDC_ECG_AMPL_QRS   MeasurementCode = "MDC_ECG_AMPL_QRS"   // QRS amplitude (µV)
	MDC_ECG_AMPL_P     MeasurementCode = "MDC_ECG_AMPL_P"     // P wave amplitude (µV)
	MDC_ECG_AMPL_T     MeasurementCode = "MDC_ECG_AMPL_T"     // T wave amplitude (µV)
	MDC_ECG_AMPL_ST    MeasurementCode = "MDC_ECG_AMPL_ST"    // ST segment amplitude (µV)
)

// =============================================================================
// HL7 AdministrativeGender Codes
// =============================================================================

type GenderCode string

const (
	GENDER_MALE             GenderCode = "M"  // Male
	GENDER_FEMALE           GenderCode = "F"  // Female
	GENDER_UNDIFFERENTIATED GenderCode = "UN" // Undifferentiated
)

// =============================================================================
// HL7 Race Codes
// =============================================================================

type RaceCode string

const (
	RACE_NATIVE_AMERICAN            RaceCode = "1002-5" // American Indian or Alaska Native
	RACE_ASIAN                      RaceCode = "2028-9" // Asian
	RACE_BLACK_OR_AFRICAN_AMERICAN  RaceCode = "2054-5" // Black or African American
	RACE_HAWAIIAN_OR_PACIFIC_ISLAND RaceCode = "2076-8" // Native Hawaiian or Other Pacific Islander
	RACE_WHITE                      RaceCode = "2106-3" // White
	RACE_OTHER                      RaceCode = "2131-1" // Other Race
)

// =============================================================================
// Suggested Codes (Not Formally Defined)
// =============================================================================

type CodeRole string

// Subject Role Codes (suggested)
const (
	SUBJECT_ROLE_SCREENING CodeRole = "SCREENING" // Subject being screened
	SUBJECT_ROLE_ENROLLED  CodeRole = "ENROLLED"  // Subject enrolled in trial
)

type DeviceTypeCode string

// Device Types (suggested, no formal vocabulary)
const (
	DEVICE_12LEAD_ECG    DeviceTypeCode = "12LEAD_ELECTROCARDIOGRAPH"
	DEVICE_12LEAD_HOLTER DeviceTypeCode = "12LEAD_HOLTER"
)

// Performer Functions (suggested, no formal vocabulary)

type PerformerFunctionCode string

const (
	PERFORMER_HOLTER_HOOKUP  PerformerFunctionCode = "HOLTER_HOOKUP_TECH"
	PERFORMER_HOLTER_ANALYST PerformerFunctionCode = "HOLTER_ANALYST"
	PERFORMER_ECG_TECHNICIAN PerformerFunctionCode = "ELECTROCARDIOGRAPH_TECH"
)

type RegionOfInterestType string

const (
	// 	– Partially specified region of interest. A partially specified bounded Region
	// of Interest (ROI) specifies a ROI in which at least all values in the dimensions
	// specified by the boundary criteria participate. For example, if an episode of
	// ventricular fibrillations (Vfib) is observed, it usually doesn’t make sense to exclude
	// any ECG leads from the observation and the partially specified ROI would contain
	// only one boundary for time indicating the time interval where Vfib was observed.
	ROIPS RegionOfInterestType = "ROIPS"
	// 	– Fully specified region of interest. A fully specified bounded Region of
	// Interest (ROI) delineates a ROI in which only those dimensions participate that are
	// specified by boundary criteria, whereas all other dimensions are excluded. For
	// example a ROI to mark an episode of “ST elevation” in a subset of the ECG leads
	// V2, V3, and V4 would include 4 boundaries, one each for time, V2, V3, and V4.
	ROIFS RegionOfInterestType = "ROIFS"
)

// Visit Type (suggested from CDISC)
type VisitTypeCode string

const (
	VISIT_TYPE_SCHEDULED   = "S" // Scheduled visit per protocol
	VISIT_TYPE_UNSCHEDULED = "U" // Unscheduled visit
)

type ConfidentialityCode string

// Confidentiality Codes (suggested from CDISC)
const (
	CONFIDENTIALITY_SPONSOR_BLINDED      ConfidentialityCode = "S" // Blinded to sponsor only
	CONFIDENTIALITY_INVESTIGATOR_BLINDED ConfidentialityCode = "I" // Blinded to investigator only
	CONFIDENTIALITY_BOTH                 ConfidentialityCode = "B" // Double-blind
	CONFIDENTIALITY_CUSTOM               ConfidentialityCode = "C" // Custom blinding
)

type ReasonCode string

// Reason Codes (suggested)
const (
	REASON_PER_PROTOCOL    ReasonCode = "PER_PROTOCOL"            // Per protocol acquisition
	REASON_NOT_IN_PROTOCOL ReasonCode = "NOT_IN_PROTOCOL"         // Not per protocol
	REASON_WRONG_EVENT     ReasonCode = "IN_PROTOCOL_WRONG_EVENT" // Wrong timepoint
)

type TreatmentGroupCode string

const (
	GRP_001 TreatmentGroupCode = "GRP_001" // Treatment Group 1
	GRP_002 TreatmentGroupCode = "GRP_002" // Treatment Group 2
	GRP_003 TreatmentGroupCode = "GRP_003" // Treatment Group 3
	GRP_004 TreatmentGroupCode = "GRP_004" // Treatment Group 4
)

// =============================================================================
// UCUM Units
// =============================================================================

const (
	// Time Units
	UNIT_SECOND      = "s"   // Second
	UNIT_MILLISECOND = "ms"  // Millisecond
	UNIT_MINUTE      = "min" // Minute
	UNIT_HOUR        = "h"   // Hour

	// Voltage Units
	UNIT_MICROVOLT = "uV" // Microvolt (most common for ECG)
	UNIT_MILLIVOLT = "mV" // Millivolt
	UNIT_VOLT      = "V"  // Volt

	// Frequency
	UNIT_HERTZ = "Hz" // Hertz

	// Heart Rate
	UNIT_BPM = "bpm" // Beats per minute
)

// =============================================================================
// Helper Functions
// =============================================================================

// IsStandardLead checks if a lead code is part of standard 12-lead ECG.
func IsStandardLead(leadCode LeadCode) bool {
	standardLeads := []LeadCode{
		MDC_ECG_LEAD_I, MDC_ECG_LEAD_II, MDC_ECG_LEAD_III,
		MDC_ECG_LEAD_AVR, MDC_ECG_LEAD_AVL, MDC_ECG_LEAD_AVF,
		MDC_ECG_LEAD_V1, MDC_ECG_LEAD_V2, MDC_ECG_LEAD_V3,
		MDC_ECG_LEAD_V4, MDC_ECG_LEAD_V5, MDC_ECG_LEAD_V6,
	}
	return slices.Contains(standardLeads, leadCode)
}

// IsTimeSequenceCode checks if code represents a time sequence.
func IsTimeSequenceCode(code TimeSequenceCode) bool {
	return code == TIME_ABSOLUTE_CODE || code == TIME_RELATIVE_CODE
}

// IsSeriesTypeCode checks if code is a valid series type.
func IsSeriesTypeCode(code SeriesTypeCode) bool {
	return code == RHYTHM_CODE ||
		code == REPRESENTATIVE_BEAT_CODE ||
		code == MEDIAN_BEAT_CODE
}

// IsValidGender checks if code is a valid HL7 gender.
func IsValidGender(code string) bool {
	return code == string(GENDER_MALE) ||
		code == string(GENDER_FEMALE) ||
		code == string(GENDER_UNDIFFERENTIATED)
}

// GetStandardLeads returns all standard 12-lead codes.
func GetStandardLeads() []LeadCode {
	return []LeadCode{
		MDC_ECG_LEAD_I, MDC_ECG_LEAD_II, MDC_ECG_LEAD_III,
		MDC_ECG_LEAD_AVR, MDC_ECG_LEAD_AVL, MDC_ECG_LEAD_AVF,
		MDC_ECG_LEAD_V1, MDC_ECG_LEAD_V2, MDC_ECG_LEAD_V3,
		MDC_ECG_LEAD_V4, MDC_ECG_LEAD_V5, MDC_ECG_LEAD_V6,
	}
}

// GetCodeSystemName returns a human-readable name for an OID.
func GetCodeSystemName(oid CodeSystemOID) string {
	names := map[CodeSystemOID]string{
		CPT_OID:                          "CPT (Current Procedural Terminology)",
		MDC_OID:                          "MDC (Medical Device Communications)",
		HL7_ActCode_OID:                  "HL7 ActCode",
		HL7_ActAdministrativeGender_OID:  "HL7 AdministrativeGender",
		HL7_Race_OID:                     "HL7 Race",
		HL7_ResearchSubjectRoleBasis_OID: "HL7 ResearchSubjectRoleBasis",
		LOINC_OID:                        "LOINC",
		SNOMED_CT_OID:                    "SNOMED CT",
		UCUM_OID:                         "UCUM",
	}
	if name, ok := names[oid]; ok {
		return name
	}
	return "Unknown Code System"
}
