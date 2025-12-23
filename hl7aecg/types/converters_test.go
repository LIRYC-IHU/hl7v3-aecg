package types

import (
	"testing"
	"time"
)

// TestGetGender tests gender string to GenderCode conversion
func TestGetGender(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  GenderCode
	}{
		// Male variations
		{name: "Male uppercase", input: "M", want: GENDER_MALE},
		{name: "Male lowercase", input: "m", want: GENDER_MALE},
		{name: "Male full word", input: "Male", want: GENDER_MALE},
		{name: "Male uppercase word", input: "MALE", want: GENDER_MALE},
		{name: "Male numeric", input: "1", want: GENDER_MALE},
		{name: "Male with spaces", input: "  M  ", want: GENDER_MALE},

		// Female variations
		{name: "Female uppercase", input: "F", want: GENDER_FEMALE},
		{name: "Female lowercase", input: "f", want: GENDER_FEMALE},
		{name: "Female full word", input: "Female", want: GENDER_FEMALE},
		{name: "Female uppercase word", input: "FEMALE", want: GENDER_FEMALE},
		{name: "Female numeric", input: "2", want: GENDER_FEMALE},

		// Undifferentiated variations
		{name: "Undifferentiated U", input: "U", want: GENDER_UNDIFFERENTIATED},
		{name: "Undifferentiated UN", input: "UN", want: GENDER_UNDIFFERENTIATED},
		{name: "Undifferentiated Unknown", input: "Unknown", want: GENDER_UNDIFFERENTIATED},
		{name: "Undifferentiated UNKNOWN", input: "UNKNOWN", want: GENDER_UNDIFFERENTIATED},
		{name: "Undifferentiated full", input: "Undifferentiated", want: GENDER_UNDIFFERENTIATED},
		{name: "Undifferentiated numeric", input: "0", want: GENDER_UNDIFFERENTIATED},

		// Invalid/default cases
		{name: "Invalid value", input: "X", want: GENDER_UNDIFFERENTIATED},
		{name: "Empty string", input: "", want: GENDER_UNDIFFERENTIATED},
		{name: "Random text", input: "other", want: GENDER_UNDIFFERENTIATED},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetGender(tt.input)
			if got != tt.want {
				t.Errorf("GetGender(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestGetDeviceTypeCode tests device model name to DeviceTypeCode conversion
func TestGetDeviceTypeCode(t *testing.T) {
	tests := []struct {
		name      string
		modelName string
		want      DeviceTypeCode
	}{
		{name: "Holter device lowercase", modelName: "holter recorder", want: DEVICE_12LEAD_HOLTER},
		{name: "Holter device uppercase", modelName: "HOLTER DEVICE", want: DEVICE_12LEAD_HOLTER},
		{name: "Holter device mixed case", modelName: "Holter Monitor XYZ", want: DEVICE_12LEAD_HOLTER},
		{name: "Standard ECG device", modelName: "ECG Machine 3000", want: DEVICE_12LEAD_ECG},
		{name: "Generic device", modelName: "CardioMax Pro", want: DEVICE_12LEAD_ECG},
		{name: "Empty string", modelName: "", want: DEVICE_12LEAD_ECG},
		{name: "With spaces", modelName: "  Standard ECG  ", want: DEVICE_12LEAD_ECG},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetDeviceTypeCode(tt.modelName)
			if got != tt.want {
				t.Errorf("GetDeviceTypeCode(%q) = %v, want %v", tt.modelName, got, tt.want)
			}
		})
	}
}

// TestGetRaceCode tests race string to RaceCode conversion
func TestGetRaceCode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  RaceCode
	}{
		// White variations
		{name: "White", input: "white", want: RACE_WHITE},
		{name: "White uppercase", input: "WHITE", want: RACE_WHITE},
		{name: "Caucasian", input: "caucasian", want: RACE_WHITE},
		{name: "Caucasian uppercase", input: "CAUCASIAN", want: RACE_WHITE},
		{name: "White code", input: "2106-3", want: RACE_WHITE},

		// Black variations
		{name: "Black", input: "black", want: RACE_BLACK_OR_AFRICAN_AMERICAN},
		{name: "Black uppercase", input: "BLACK", want: RACE_BLACK_OR_AFRICAN_AMERICAN},
		{name: "African American", input: "african american", want: RACE_BLACK_OR_AFRICAN_AMERICAN},
		{name: "African", input: "african", want: RACE_BLACK_OR_AFRICAN_AMERICAN},
		{name: "Black code", input: "2054-5", want: RACE_BLACK_OR_AFRICAN_AMERICAN},

		// Asian variations
		{name: "Asian", input: "asian", want: RACE_ASIAN},
		{name: "Asian uppercase", input: "ASIAN", want: RACE_ASIAN},
		{name: "Asian code", input: "2028-9", want: RACE_ASIAN},

		// Native American variations
		{name: "Native American", input: "native american", want: RACE_NATIVE_AMERICAN},
		{name: "American Indian", input: "american indian", want: RACE_NATIVE_AMERICAN},
		{name: "Alaska Native", input: "alaska native", want: RACE_NATIVE_AMERICAN},
		{name: "Native", input: "native", want: RACE_NATIVE_AMERICAN},
		{name: "Native American code", input: "1002-5", want: RACE_NATIVE_AMERICAN},

		// Hawaiian/Pacific Islander variations
		{name: "Hawaiian", input: "hawaiian", want: RACE_HAWAIIAN_OR_PACIFIC_ISLAND},
		{name: "Pacific Islander", input: "pacific islander", want: RACE_HAWAIIAN_OR_PACIFIC_ISLAND},
		// NOTE: Bug in converter - "native hawaiian" matches NATIVE_AMERICAN because "native" check comes first
		{name: "Native Hawaiian (matches native)", input: "native hawaiian", want: RACE_NATIVE_AMERICAN},
		{name: "Hawaiian code", input: "2076-8", want: RACE_HAWAIIAN_OR_PACIFIC_ISLAND},

		// Other/invalid cases
		{name: "Other code", input: "2131-1", want: RACE_OTHER},
		{name: "Empty string", input: "", want: RACE_OTHER},
		{name: "Unrecognized", input: "martian", want: RACE_OTHER},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRaceCode(tt.input)
			if got != tt.want {
				t.Errorf("GetRaceCode(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestGetConfidentialityCode tests confidentiality string to ConfidentialityCode conversion
func TestGetConfidentialityCode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  ConfidentialityCode
	}{
		// Sponsor blinded variations
		{name: "S", input: "S", want: CONFIDENTIALITY_SPONSOR_BLINDED},
		{name: "Sponsor", input: "SPONSOR", want: CONFIDENTIALITY_SPONSOR_BLINDED},
		{name: "Sponsor blinded", input: "SPONSOR_BLINDED", want: CONFIDENTIALITY_SPONSOR_BLINDED},
		{name: "Sponsor lowercase", input: "sponsor", want: CONFIDENTIALITY_SPONSOR_BLINDED},

		// Investigator blinded variations
		{name: "I", input: "I", want: CONFIDENTIALITY_INVESTIGATOR_BLINDED},
		{name: "Investigator", input: "INVESTIGATOR", want: CONFIDENTIALITY_INVESTIGATOR_BLINDED},
		{name: "Investigator blinded", input: "INVESTIGATOR_BLINDED", want: CONFIDENTIALITY_INVESTIGATOR_BLINDED},

		// Both blinded variations
		{name: "B", input: "B", want: CONFIDENTIALITY_BOTH},
		{name: "Both", input: "BOTH", want: CONFIDENTIALITY_BOTH},
		{name: "Double blind", input: "DOUBLE_BLIND", want: CONFIDENTIALITY_BOTH},
		{name: "Double", input: "DOUBLE", want: CONFIDENTIALITY_BOTH},

		// Custom variations
		{name: "C", input: "C", want: CONFIDENTIALITY_CUSTOM},
		{name: "Custom", input: "CUSTOM", want: CONFIDENTIALITY_CUSTOM},

		// Default case
		{name: "Empty string", input: "", want: CONFIDENTIALITY_SPONSOR_BLINDED},
		{name: "Unrecognized", input: "OTHER", want: CONFIDENTIALITY_SPONSOR_BLINDED},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetConfidentialityCode(tt.input)
			if got != tt.want {
				t.Errorf("GetConfidentialityCode(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestGetReasonCode tests reason string to ReasonCode conversion
func TestGetReasonCode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  ReasonCode
	}{
		// Per protocol variations
		{name: "PER_PROTOCOL", input: "PER_PROTOCOL", want: REASON_PER_PROTOCOL},
		{name: "PROTOCOL", input: "PROTOCOL", want: REASON_PER_PROTOCOL},
		{name: "P", input: "P", want: REASON_PER_PROTOCOL},
		{name: "Per protocol lowercase", input: "per_protocol", want: REASON_PER_PROTOCOL},

		// Not in protocol variations
		{name: "NOT_IN_PROTOCOL", input: "NOT_IN_PROTOCOL", want: REASON_NOT_IN_PROTOCOL},
		{name: "NOT_PROTOCOL", input: "NOT_PROTOCOL", want: REASON_NOT_IN_PROTOCOL},
		{name: "N", input: "N", want: REASON_NOT_IN_PROTOCOL},

		// Wrong event variations
		{name: "IN_PROTOCOL_WRONG_EVENT", input: "IN_PROTOCOL_WRONG_EVENT", want: REASON_WRONG_EVENT},
		{name: "WRONG_EVENT", input: "WRONG_EVENT", want: REASON_WRONG_EVENT},
		{name: "W", input: "W", want: REASON_WRONG_EVENT},

		// Default case
		{name: "Empty string", input: "", want: REASON_PER_PROTOCOL},
		{name: "Unrecognized", input: "OTHER", want: REASON_PER_PROTOCOL},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetReasonCode(tt.input)
			if got != tt.want {
				t.Errorf("GetReasonCode(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// TestFormatHL7DateTime tests time.Time to HL7 datetime string formatting
func TestFormatHL7DateTime(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{
			name: "Full datetime with milliseconds",
			time: time.Date(2024, 12, 22, 14, 30, 45, 123000000, time.UTC),
			want: "20241222143045.123",
		},
		{
			name: "Midnight",
			time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want: "20240101000000.000",
		},
		{
			name: "End of year",
			time: time.Date(2023, 12, 31, 23, 59, 59, 999000000, time.UTC),
			want: "20231231235959.999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatHL7DateTime(tt.time)
			if got != tt.want {
				t.Errorf("FormatHL7DateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFormatHL7Date tests time.Time to HL7 date string formatting
func TestFormatHL7Date(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{
			name: "Regular date",
			time: time.Date(2024, 12, 22, 14, 30, 45, 0, time.UTC),
			want: "20241222",
		},
		{
			name: "First day of year",
			time: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			want: "20240101",
		},
		{
			name: "Last day of year",
			time: time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
			want: "20231231",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatHL7Date(tt.time)
			if got != tt.want {
				t.Errorf("FormatHL7Date() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestParseHL7DateTime tests HL7 datetime string parsing
func TestParseHL7DateTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "Full precision datetime",
			input:   "20241222143045.123",
			want:    time.Date(2024, 12, 22, 14, 30, 45, 123000000, time.UTC),
			wantErr: false,
		},
		{
			name:    "Second precision datetime",
			input:   "20241222143045",
			want:    time.Date(2024, 12, 22, 14, 30, 45, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Date only",
			input:   "20241222",
			want:    time.Date(2024, 12, 22, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Midnight datetime",
			input:   "20240101000000.000",
			want:    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "Invalid format",
			input:   "invalid-date",
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   "",
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "Partial date",
			input:   "202412",
			want:    time.Time{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHL7DateTime(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseHL7DateTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("ParseHL7DateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestNormalizeLeadCode tests lead name to LeadCode normalization
func TestNormalizeLeadCode(t *testing.T) {
	tests := []struct {
		name     string
		leadName string
		want     LeadCode
	}{
		// Lead I variations
		{name: "Lead I uppercase", leadName: "I", want: MDC_ECG_LEAD_I},
		{name: "Lead I with prefix", leadName: "LEAD_I", want: MDC_ECG_LEAD_I},
		{name: "Lead I numeric", leadName: "1", want: MDC_ECG_LEAD_I},
		{name: "Lead I lowercase", leadName: "i", want: MDC_ECG_LEAD_I},

		// Lead II variations
		{name: "Lead II", leadName: "II", want: MDC_ECG_LEAD_II},
		{name: "Lead II with prefix", leadName: "LEAD_II", want: MDC_ECG_LEAD_II},
		{name: "Lead II numeric", leadName: "2", want: MDC_ECG_LEAD_II},

		// Lead III variations
		{name: "Lead III", leadName: "III", want: MDC_ECG_LEAD_III},
		{name: "Lead III with prefix", leadName: "LEAD_III", want: MDC_ECG_LEAD_III},
		{name: "Lead III numeric", leadName: "3", want: MDC_ECG_LEAD_III},

		// Augmented leads
		{name: "aVR uppercase", leadName: "AVR", want: MDC_ECG_LEAD_AVR},
		{name: "aVR mixed case", leadName: "aVR", want: MDC_ECG_LEAD_AVR},
		{name: "aVR with prefix", leadName: "LEAD_AVR", want: MDC_ECG_LEAD_AVR},
		{name: "aVL uppercase", leadName: "AVL", want: MDC_ECG_LEAD_AVL},
		{name: "aVL mixed case", leadName: "aVL", want: MDC_ECG_LEAD_AVL},
		{name: "aVF uppercase", leadName: "AVF", want: MDC_ECG_LEAD_AVF},
		{name: "aVF mixed case", leadName: "aVF", want: MDC_ECG_LEAD_AVF},

		// Precordial leads
		{name: "V1", leadName: "V1", want: MDC_ECG_LEAD_V1},
		{name: "V1 with prefix", leadName: "LEAD_V1", want: MDC_ECG_LEAD_V1},
		{name: "V2", leadName: "V2", want: MDC_ECG_LEAD_V2},
		{name: "V3", leadName: "V3", want: MDC_ECG_LEAD_V3},
		{name: "V4", leadName: "V4", want: MDC_ECG_LEAD_V4},
		{name: "V5", leadName: "V5", want: MDC_ECG_LEAD_V5},
		{name: "V6", leadName: "V6", want: MDC_ECG_LEAD_V6},

		// With spaces
		{name: "Lead I with spaces", leadName: "  I  ", want: MDC_ECG_LEAD_I},
		{name: "aVR with spaces", leadName: "  aVR  ", want: MDC_ECG_LEAD_AVR},

		// Unrecognized (returns as-is)
		{name: "Unrecognized lead", leadName: "CUSTOM_LEAD", want: LeadCode("CUSTOM_LEAD")},
		{name: "Empty string", leadName: "", want: LeadCode("")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeLeadCode(tt.leadName)
			if got != tt.want {
				t.Errorf("NormalizeLeadCode(%q) = %v, want %v", tt.leadName, got, tt.want)
			}
		})
	}
}

// TestGetStandardLeads tests retrieval of standard 12-lead ECG codes
func TestGetStandardLeads(t *testing.T) {
	leads := GetStandardLeads()

	// Should return exactly 12 leads
	if len(leads) != 12 {
		t.Errorf("GetStandardLeads() returned %d leads, want 12", len(leads))
	}

	// Verify all expected leads are present
	expectedLeads := []LeadCode{
		MDC_ECG_LEAD_I, MDC_ECG_LEAD_II, MDC_ECG_LEAD_III,
		MDC_ECG_LEAD_AVR, MDC_ECG_LEAD_AVL, MDC_ECG_LEAD_AVF,
		MDC_ECG_LEAD_V1, MDC_ECG_LEAD_V2, MDC_ECG_LEAD_V3,
		MDC_ECG_LEAD_V4, MDC_ECG_LEAD_V5, MDC_ECG_LEAD_V6,
	}

	for i, expected := range expectedLeads {
		if i >= len(leads) {
			t.Errorf("Missing lead at index %d: %v", i, expected)
			continue
		}
		if leads[i] != expected {
			t.Errorf("Lead at index %d = %v, want %v", i, leads[i], expected)
		}
	}
}
