package types

import (
	"fmt"
	"strings"
	"time"
)

// =============================================================================
// Converter Functions
// =============================================================================
// These functions convert extracted data into HL7 aECG code values

// GetGender converts a gender string to HL7 GenderCode
//
// Accepts various formats:
//   - "M", "Male", "MALE", "m" → GENDER_MALE
//   - "F", "Female", "FEMALE", "f" → GENDER_FEMALE
//   - "U", "Unknown", "UNKNOWN", "Undifferentiated" → GENDER_UNDIFFERENTIATED
//
// Returns GENDER_UNDIFFERENTIATED for unrecognized values.
func GetGender(s string) GenderCode {
	normalized := strings.ToUpper(strings.TrimSpace(s))

	switch normalized {
	case "M", "MALE", "1":
		return GENDER_MALE
	case "F", "FEMALE", "2":
		return GENDER_FEMALE
	case "U", "UN", "UNKNOWN", "UNDIFFERENTIATED", "0":
		return GENDER_UNDIFFERENTIATED
	default:
		return GENDER_UNDIFFERENTIATED
	}
}

// GetDeviceTypeCode converts a manufacturer model name to a DeviceTypeCode
//
// This is a best-effort conversion as HL7 aECG doesn't define a formal
// vocabulary for device types. Adjust as needed for your use case.
func GetDeviceTypeCode(modelName string) DeviceTypeCode {
	normalized := strings.ToLower(strings.TrimSpace(modelName))

	// Check for holter keywords
	if strings.Contains(normalized, "holter") {
		return DEVICE_12LEAD_HOLTER
	}

	// Default to standard 12-lead ECG
	return DEVICE_12LEAD_ECG
}

// GetRaceCode converts a race string to HL7 RaceCode
//
// Accepts common race descriptions and maps them to HL7 codes.
// Returns RACE_OTHER for unrecognized values.
func GetRaceCode(s string) RaceCode {
	normalized := strings.ToLower(strings.TrimSpace(s))

	// Check for keywords
	if strings.Contains(normalized, "white") || strings.Contains(normalized, "caucasian") {
		return RACE_WHITE
	}
	if strings.Contains(normalized, "black") || strings.Contains(normalized, "african") {
		return RACE_BLACK_OR_AFRICAN_AMERICAN
	}
	if strings.Contains(normalized, "asian") {
		return RACE_ASIAN
	}
	if strings.Contains(normalized, "native") || strings.Contains(normalized, "american indian") ||
		strings.Contains(normalized, "alaska") {
		return RACE_NATIVE_AMERICAN
	}
	if strings.Contains(normalized, "hawaiian") || strings.Contains(normalized, "pacific") {
		return RACE_HAWAIIAN_OR_PACIFIC_ISLAND
	}

	// Check for direct codes
	switch s {
	case "1002-5":
		return RACE_NATIVE_AMERICAN
	case "2028-9":
		return RACE_ASIAN
	case "2054-5":
		return RACE_BLACK_OR_AFRICAN_AMERICAN
	case "2076-8":
		return RACE_HAWAIIAN_OR_PACIFIC_ISLAND
	case "2106-3":
		return RACE_WHITE
	case "2131-1":
		return RACE_OTHER
	}

	return RACE_OTHER
}

// GetConfidentialityCode converts a string to ConfidentialityCode
//
// Maps common confidentiality descriptions to codes.
func GetConfidentialityCode(s string) ConfidentialityCode {
	normalized := strings.ToUpper(strings.TrimSpace(s))

	switch normalized {
	case "S", "SPONSOR", "SPONSOR_BLINDED":
		return CONFIDENTIALITY_SPONSOR_BLINDED
	case "I", "INVESTIGATOR", "INVESTIGATOR_BLINDED":
		return CONFIDENTIALITY_INVESTIGATOR_BLINDED
	case "B", "BOTH", "DOUBLE_BLIND", "DOUBLE":
		return CONFIDENTIALITY_BOTH
	case "C", "CUSTOM":
		return CONFIDENTIALITY_CUSTOM
	default:
		return CONFIDENTIALITY_SPONSOR_BLINDED // Default
	}
}

// GetReasonCode converts a string to ReasonCode
//
// Maps common reason descriptions to codes.
func GetReasonCode(s string) ReasonCode {
	normalized := strings.ToUpper(strings.TrimSpace(s))

	switch normalized {
	case "PER_PROTOCOL", "PROTOCOL", "P":
		return REASON_PER_PROTOCOL
	case "NOT_IN_PROTOCOL", "NOT_PROTOCOL", "N":
		return REASON_NOT_IN_PROTOCOL
	case "IN_PROTOCOL_WRONG_EVENT", "WRONG_EVENT", "W":
		return REASON_WRONG_EVENT
	default:
		return REASON_PER_PROTOCOL // Default
	}
}

// FormatHL7DateTime converts a time.Time to HL7 datetime format (YYYYMMDDHHmmss.SSS)
//
// Example: 2024-12-22 14:30:45.123 → "20241222143045.123"
func FormatHL7DateTime(t time.Time) string {
	return t.Format("20060102150405.000")
}

// FormatHL7Date converts a time.Time to HL7 date format (YYYYMMDD)
//
// Example: 2024-12-22 → "20241222"
func FormatHL7Date(t time.Time) string {
	return t.Format("20060102")
}

// ParseHL7DateTime parses an HL7 datetime string (YYYYMMDDHHmmss.SSS)
//
// Accepts formats:
//   - YYYYMMDDHHmmss.SSS (full precision)
//   - YYYYMMDDHHmmss (second precision)
//   - YYYYMMDD (date only)
func ParseHL7DateTime(s string) (time.Time, error) {
	formats := []string{
		"20060102150405.000", // Full precision
		"20060102150405",     // Second precision
		"20060102",           // Date only
	}

	for _, format := range formats {
		if t, err := time.Parse(format, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("cannot parse HL7 datetime: %s", s)
}

// NormalizeLeadCode normalizes lead names to MDC codes
//
// Accepts various lead name formats and returns the appropriate MDC_ECG_LEAD_* code.
func NormalizeLeadCode(leadName string) LeadCode {
	normalized := strings.ToUpper(strings.TrimSpace(leadName))

	switch normalized {
	case "I", "LEAD_I", "1":
		return MDC_ECG_LEAD_I
	case "II", "LEAD_II", "2":
		return MDC_ECG_LEAD_II
	case "III", "LEAD_III", "3":
		return MDC_ECG_LEAD_III
	case "AVR", "aVR", "LEAD_AVR":
		return MDC_ECG_LEAD_AVR
	case "AVL", "aVL", "LEAD_AVL":
		return MDC_ECG_LEAD_AVL
	case "AVF", "aVF", "LEAD_AVF":
		return MDC_ECG_LEAD_AVF
	case "V1", "LEAD_V1":
		return MDC_ECG_LEAD_V1
	case "V2", "LEAD_V2":
		return MDC_ECG_LEAD_V2
	case "V3", "LEAD_V3":
		return MDC_ECG_LEAD_V3
	case "V4", "LEAD_V4":
		return MDC_ECG_LEAD_V4
	case "V5", "LEAD_V5":
		return MDC_ECG_LEAD_V5
	case "V6", "LEAD_V6":
		return MDC_ECG_LEAD_V6
	default:
		return LeadCode(normalized) // Return as-is if not recognized
	}
}
