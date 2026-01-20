# HL7 aECG Library for Go

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.20-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A comprehensive Go library for generating HL7 v3 Annotated Electrocardiogram (aECG) XML files compliant with the HL7 aECG Implementation Guide (Final 21-March-2005) for FDA clinical trial submissions.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Documentation](#documentation)
  - [Creating an aECG Document](#creating-an-aecg-document)
  - [Adding ECG Waveform Data](#adding-ecg-waveform-data)
  - [Adding Annotations](#adding-annotations)
  - [Subject Demographics](#subject-demographics)
  - [Clinical Trial Information](#clinical-trial-information)
- [API Reference](#api-reference)
- [Code Systems](#code-systems)
- [Examples](#examples)
- [Testing](#testing)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## Features

### Core Functionality

✅ **Complete HL7 v3 aECG Implementation**

- Full support for AnnotatedECG document structure
- Clinical trial metadata (protocol, sponsor, sites)
- Subject demographics and assignment
- ECG series (rhythm, representative beat, median beat)
- Device and technician information

✅ **ECG Waveform Data**

- Standard 12-lead ECG support
- GLIST_TS (Generated List of Timestamps) for time sequences
- SLIST_PQ (Scaled List of Physical Quantities) for voltage data
- SLIST_INT (Scaled List of Integers) for integer sequences
- Configurable sample rates and scaling

✅ **Comprehensive Annotations**

- Numeric annotations (Physical Quantities with units)
- Text annotations (interpretations, comments, statements)
- Lead-specific annotations with supportingROI
- Nested annotations for complex structures
- Global measurements (heart rate, intervals, amplitudes)

✅ **Code Systems**

- CPT (Current Procedural Terminology) codes
- MDC (Medical Device Codes) - ISO/IEEE 11073-10101
- HL7 Act codes (series types, confidentiality, reasons)
- HL7 Administrative Gender codes
- HL7 Race codes

✅ **Validation**

- Comprehensive document validation
- Format checking (timestamps, IDs, codes)
- Required field validation
- Type-safe generic code validation

✅ **Type Safety**

- Generic types for compile-time code validation
- Strong typing for all HL7 structures
- Fluent API with method chaining

## Installation

```bash
go get github.com/LIRYC-IHU/hl7v3-aecg
```

## Quick Start

```go
package main

import (
    "log"
    "github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg"
    "github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
)

func main() {
    // 1. Create a new aECG instance
    h := hl7aecg.NewHl7xml("./output")

    // 2. Initialize with ECG procedure code
    h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "CPT-4", "")

    // 3. Set document metadata
    h.SetEffectiveTime("20250923103550", "20250923103600", nil, nil)
    h.SetText("12-lead ECG for clinical trial")

    // 4. Configure subject demographics
    h.SetSubjectDemographics(
        "JDO",                         // name/initials
        "SUBJ_001",                    // patient ID
        types.GENDER_MALE,             // gender
        "19530508",                    // birth date
        types.RACE_WHITE,              // race code
    )

    // 5. Add rhythm series with waveform data
    sampleRate := 500.0 // Hz
    leads := map[types.LeadCode][]int{
        types.MDC_ECG_LEAD_I:   {/* voltage samples */},
        types.MDC_ECG_LEAD_II:  {/* voltage samples */},
        types.MDC_ECG_LEAD_III: {/* voltage samples */},
        // ... other leads
    }

    h.AddRhythmSeries(
        "20250923103550.000",  // start time
        "20250923103600.000",  // end time
        sampleRate,            // sample rate (Hz)
        leads,                 // lead data
        0,                     // origin (µV)
        5,                     // scale (µV per digit)
    )

    // 6. Add device information
    h.SetSeriesAuthor(
        "DEVICE_001",                      // device ID
        types.DEVICE_12LEAD_ECG,          // device type
        "ECG Model X",                     // model name
        "v1.0",                            // software version
        "2.16.840.1.113883.3.example",    // manufacturer OID
        "ACME Medical Devices",            // manufacturer name
    )

    // 7. Add annotations
    annSet := h.HL7AEcg.Component[0].Series.SubjectOf.AnnotationSet

    // Add heart rate
    annSet.AddHeartRate(72)

    // Add QT interval
    annSet.AddQTInterval(400)

    // Add ECG interpretation
    interpIdx := annSet.AddTextAnnotation(
        "MDC_ECG_INTERPRETATION",
        "2.16.840.1.113883.6.24",
        "",
    )
    interp := annSet.GetAnnotation(interpIdx)
    interp.AddNestedTextAnnotation(
        "MDC_ECG_INTERPRETATION_STATEMENT",
        "2.16.840.1.113883.6.24",
        "Normal sinus rhythm",
    )

    // 8. Validate
    if err := h.Validate(); err != nil {
        log.Fatalf("Validation failed: %v", err)
    }

    // 9. Export to XML
    if _, err := h.Test(); err != nil {
        log.Fatalf("Export failed: %v", err)
    }

    log.Println("aECG XML generated successfully!")
}
```

## Documentation

### Creating an aECG Document

#### Initialize Document

```go
h := hl7aecg.NewHl7xml("./output")

// Set procedure code (CPT)
h.Initialize(
    types.CPT_CODE_ECG_Routine,  // procedure code
    types.CPT_OID,                // code system
    "CPT-4",                      // code system name
    "",                           // display name (optional)
)

// Set document ID (optional - auto-generated if not set)
h.HL7AEcg.SetRootID("1.2.3.4.5.6.7", "annotatedEcg")

// Set effective time (acquisition period)
h.SetEffectiveTime(
    "20250923103550",  // start time (YYYYMMDDHHmmss)
    "20250923103600",  // end time
    nil,               // low (alternative format)
    nil,               // high (alternative format)
)

// Set narrative text
h.SetText("12-lead resting ECG")

// Set confidentiality code (optional)
h.AddConfidentialityCode(types.CONFIDENTIALITY_NORMAL)

// Set reason code (optional)
h.AddReasonCode(types.REASON_PER_PROTOCOL)
```

### Adding ECG Waveform Data

#### Rhythm Series

```go
// Prepare waveform data (voltage samples for each lead)
leads := map[types.LeadCode][]int{
    types.MDC_ECG_LEAD_I:   leadIData,
    types.MDC_ECG_LEAD_II:  leadIIData,
    types.MDC_ECG_LEAD_III: leadIIIData,
    types.MDC_ECG_LEAD_AVR: leadAVRData,
    types.MDC_ECG_LEAD_AVL: leadAVLData,
    types.MDC_ECG_LEAD_AVF: leadAVFData,
    types.MDC_ECG_LEAD_V1:  leadV1Data,
    types.MDC_ECG_LEAD_V2:  leadV2Data,
    types.MDC_ECG_LEAD_V3:  leadV3Data,
    types.MDC_ECG_LEAD_V4:  leadV4Data,
    types.MDC_ECG_LEAD_V5:  leadV5Data,
    types.MDC_ECG_LEAD_V6:  leadV6Data,
}

h.AddRhythmSeries(
    "20250923103550.000",  // start time (YYYYMMDDHHmmss.SSS)
    "20250923103600.000",  // end time
    500.0,                 // sample rate (Hz)
    leads,                 // voltage data
    0,                     // origin (µV)
    5,                     // scale (µV per digit)
)
```

#### Representative Beat Series

```go
// For derived representative beats
h.AddRepresentativeBeatSeries(
    "20250923103555.000",
    "20250923103555.800",
    500.0,
    representativeBeatLeads,
    0,
    5,
)
```

#### Device Information

```go
h.SetSeriesAuthor(
    "SN123456",                        // serial number
    types.DEVICE_12LEAD_ECG,          // device type code
    "CardioMax 3000",                  // model name
    "v2.1.5",                          // software version
    "2.16.840.1.113883.3.acme",       // manufacturer OID
    "ACME Medical Devices Inc.",       // manufacturer name
)
```

### Adding Annotations

Annotations allow you to add measurements, interpretations, and observations to the ECG data.

#### Global Measurements

```go
// Access the annotation set for the most recent series
annSet := h.HL7AEcg.Component[0].Series.SubjectOf.AnnotationSet

// Heart rate (bpm)
annSet.AddHeartRate(72)

// QRS duration (ms)
annSet.AddQRSDuration(88)

// QT interval (ms)
annSet.AddQTInterval(400)

// QTc interval with nested correction formula
qtcIdx := annSet.AddQTcInterval(0)  // parent has no value
qtcAnn := annSet.GetAnnotation(qtcIdx)
qtcAnn.AddNestedAnnotationWithCodeSystemName(
    "VENDOR_QTcH",    // vendor-specific correction
    "VENDOR",
    413,
    "ms",
)
```

#### Numeric Annotations

```go
// Standard annotation (with OID)
annSet.AddAnnotation(
    "MDC_ECG_HEART_RATE",              // code
    "2.16.840.1.113883.6.24",          // MDC OID
    72,                                 // value
    "bpm",                              // unit
)

// Vendor-specific annotation (with code system name)
annSet.AddAnnotationWithCodeSystemName(
    "VENDOR_P_AXIS",                   // vendor code
    "VENDOR",                          // code system name
    60,                                 // value
    "deg",                              // unit
)
```

#### Text Annotations

```go
// ECG interpretation with nested statements
interpIdx := annSet.AddTextAnnotation(
    "MDC_ECG_INTERPRETATION",
    "2.16.840.1.113883.6.24",
    "",  // empty value - container for nested statements
)

interp := annSet.GetAnnotation(interpIdx)

// Add nested interpretation statements
interp.AddNestedTextAnnotation(
    "MDC_ECG_INTERPRETATION_STATEMENT",
    "2.16.840.1.113883.6.24",
    "Sinus rhythm with occasional PVCs",
)

interp.AddNestedTextAnnotation(
    "MDC_ECG_INTERPRETATION_SUMMARY",
    "2.16.840.1.113883.6.24",
    "Borderline ECG",
)

interp.AddNestedTextAnnotation(
    "MDC_ECG_INTERPRETATION_COMMENT",
    "2.16.840.1.113883.6.24",
    "Recommend follow-up in 6 months",
)
```

#### Lead-Specific Annotations

```go
// Annotation with supportingROI for specific lead
leadIdx := annSet.AddLeadAnnotation(
    "MDC_ECG_LEAD_II",            // lead code
    "VENDOR_MEASUREMENT_MATRIX",   // annotation code
    "VENDOR_MEASUREMENT_MATRIX",   // code system
    "VENDOR",                      // code system name
)

leadAnn := annSet.GetAnnotation(leadIdx)

// Add nested measurements for this lead
leadAnn.AddNestedAnnotationWithCodeSystemName("VENDOR_P_ONSET", "VENDOR", 234, "ms")
leadAnn.AddNestedAnnotationWithCodeSystemName("VENDOR_P_DURATION", "VENDOR", 110, "ms")
leadAnn.AddNestedAnnotationWithCodeSystemName("VENDOR_QRS_ONSET", "VENDOR", 428, "ms")
leadAnn.AddNestedAnnotationWithCodeSystemName("VENDOR_R_AMP", "VENDOR", 1.5, "mV")
```

### Subject Demographics

```go
// Basic subject information
h.SetSubject(
    "SUBJ_001",                    // subject ID
    "",                            // extension (optional)
    types.SUBJECT_ROLE_ENROLLED,   // role in trial
)

// Detailed demographics
h.SetSubjectDemographics(
    "JDO",                         // name (initials or full name)
    "SUBJ_001",                    // patient ID
    types.GENDER_MALE,             // gender (M/F/UN)
    "19530508",                    // birth date (YYYYMMDD)
    types.RACE_WHITE,              // race code
)

// Multiple race codes
h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.
    SubjectAssignment.Subject.TrialSubject.SubjectDemographicPerson.
    AddRaceCode(types.RACE_ASIAN)
```

### Clinical Trial Information

```go
// Clinical trial
h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.
    SubjectAssignment.ComponentOf.ClinicalTrial.
    SetID("2.16.840.1.113883.3.example.trial", "TRIAL_001")

h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.
    SubjectAssignment.ComponentOf.ClinicalTrial.
    SetTitle("Phase III Cardiovascular Safety Study")

// Protocol
protocol := &h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.
    SubjectAssignment.ComponentOf.ClinicalTrial.Protocol

protocol.SetID("2.16.840.1.113883.3.example.protocol", "PROTO_001")
protocol.SetTitle("ECG Monitoring Protocol v2.1")

// Sponsor
sponsor := &h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.
    SubjectAssignment.ComponentOf.ClinicalTrial.Sponsor

sponsor.SetOrganizationID(
    "2.16.840.1.113883.3.example.sponsor",
    "SPONSOR_001",
)
sponsor.SetOrganizationName("Pharmaceutical Research Corp.")
```

## API Reference

### Main Package (`hl7aecg`)

#### Hl7xml

The main structure for building aECG documents.

**Constructor:**

```go
func NewHl7xml(outputDir string) *Hl7xml
```

**Initialization:**

```go
func (h *Hl7xml) Initialize(code, codeSystem, codeSystemName, displayName string) *Hl7xml
```

**Document Metadata:**

```go
func (h *Hl7xml) SetText(text string) *Hl7xml
func (h *Hl7xml) SetEffectiveTime(low, high string, lowTime, highTime *types.Time) *Hl7xml
func (h *Hl7xml) AddConfidentialityCode(code types.ConfidentialityCode) *Hl7xml
func (h *Hl7xml) AddReasonCode(code types.ReasonCode) *Hl7xml
```

**Subject Information:**

```go
func (h *Hl7xml) SetSubject(id, extension string, role types.SubjectRoleCode) *Hl7xml
func (h *Hl7xml) SetSubjectDemographics(name, patientID string, gender types.GenderCode, birthTime string, race types.RaceCode) *Hl7xml
```

**ECG Series:**

```go
func (h *Hl7xml) AddRhythmSeries(startTime, endTime string, sampleRate float64, leads map[types.LeadCode][]int, origin int, scale int) *Hl7xml
func (h *Hl7xml) AddRepresentativeBeatSeries(startTime, endTime string, sampleRate float64, leads map[types.LeadCode][]int, origin int, scale int) *Hl7xml
func (h *Hl7xml) SetSeriesAuthor(deviceID string, deviceType types.DeviceCode, modelName, softwareVersion, manufacturerOID, manufacturerName string) *Hl7xml
```

**Validation:**

```go
func (h *Hl7xml) Validate() error
```

**Export:**

```go
func (h *Hl7xml) Test() (string, error)
```

### Types Package (`hl7aecg/types`)

#### AnnotationSet

Container for annotations.

**Methods:**

```go
// Global measurements
func (as *AnnotationSet) AddHeartRate(bpm float64) int
func (as *AnnotationSet) AddQRSDuration(durationMs float64) int
func (as *AnnotationSet) AddQTInterval(intervalMs float64) int
func (as *AnnotationSet) AddQTcInterval(intervalMs float64) int

// Numeric annotations
func (as *AnnotationSet) AddAnnotation(code, codeSystem string, value float64, unit string) int
func (as *AnnotationSet) AddAnnotationWithCodeSystemName(code, codeSystemName string, value float64, unit string) int

// Text annotations
func (as *AnnotationSet) AddTextAnnotation(code, codeSystem, text string) int
func (as *AnnotationSet) AddTextAnnotationWithCodeSystemName(code, codeSystemName, text string) int

// Lead-specific annotations
func (as *AnnotationSet) AddLeadAnnotation(leadCode, code, codeSystem, codeSystemName string) int

// Accessors
func (as *AnnotationSet) GetAnnotation(idx int) *Annotation
func (as *AnnotationSet) GetAnnotationByCode(code string) *Annotation
```

#### Annotation

Individual annotation.

**Methods:**

```go
// Nested numeric annotations
func (a *Annotation) AddNestedAnnotation(code, codeSystem string, value float64, unit string) int
func (a *Annotation) AddNestedAnnotationWithCodeSystemName(code, codeSystemName string, value float64, unit string) int

// Nested text annotations
func (a *Annotation) AddNestedTextAnnotation(code, codeSystem, text string) int
func (a *Annotation) AddNestedTextAnnotationWithCodeSystemName(code, codeSystemName, text string) int

// Accessors
func (a *Annotation) GetNestedAnnotation(idx int) *Annotation
func (a *Annotation) GetValueFloat() (float64, bool)
```

#### AnnotationValue

Polymorphic value type supporting Physical Quantities (PQ) and String Values (ST).

**Methods:**

```go
func (av *AnnotationValue) IsPQ() bool
func (av *AnnotationValue) IsST() bool
func (av *AnnotationValue) GetValueFloat() (float64, bool)
func (av *AnnotationValue) GetValueUnit() string
func (av *AnnotationValue) GetText() (string, bool)
```

## Code Systems

### CPT Codes (Current Procedural Terminology)

```go
const (
    CPT_OID = "2.16.840.1.113883.6.12"

    CPT_CODE_ECG_Routine           = "93000"  // 12-lead ECG
    CPT_CODE_ECG_Rhythm            = "93040"  // Rhythm ECG
    CPT_CODE_ECG_Signal_Averaged   = "93278"  // Signal-averaged ECG
    // ... more codes
)
```

### MDC Codes (Medical Device Codes - ISO/IEEE 11073)

```go
const (
    MDC_OID = "2.16.840.1.113883.6.24"

    // Standard 12 leads
    MDC_ECG_LEAD_I   LeadCode = "2:1"
    MDC_ECG_LEAD_II  LeadCode = "2:2"
    MDC_ECG_LEAD_III LeadCode = "2:61"
    MDC_ECG_LEAD_AVR LeadCode = "2:62"
    MDC_ECG_LEAD_AVL LeadCode = "2:63"
    MDC_ECG_LEAD_AVF LeadCode = "2:64"
    MDC_ECG_LEAD_V1  LeadCode = "2:3"
    MDC_ECG_LEAD_V2  LeadCode = "2:4"
    MDC_ECG_LEAD_V3  LeadCode = "2:5"
    MDC_ECG_LEAD_V4  LeadCode = "2:6"
    MDC_ECG_LEAD_V5  LeadCode = "2:7"
    MDC_ECG_LEAD_V6  LeadCode = "2:8"
)
```

### HL7 Act Codes

```go
const (
    HL7_ActCode_OID = "2.16.840.1.113883.5.4"

    // Series types
    SERIES_RHYTHM              = "RHYTHM"
    SERIES_REPRESENTATIVE_BEAT = "REPRESENTATIVE_BEAT"
    SERIES_MEDIAN_BEAT         = "MEDIAN_BEAT"

    // Confidentiality codes
    CONFIDENTIALITY_NORMAL              ConfidentialityCode = "N"
    CONFIDENTIALITY_INVESTIGATOR_BLINDED ConfidentialityCode = "I"
    CONFIDENTIALITY_SUBJECT_BLINDED     ConfidentialityCode = "S"
    CONFIDENTIALITY_BOTH_BLINDED        ConfidentialityCode = "B"

    // Reason codes
    REASON_PER_PROTOCOL         ReasonCode = "PER_PROTOCOL"
    REASON_NOT_IN_PROTOCOL      ReasonCode = "NOT_IN_PROTOCOL"
    REASON_IN_PROTOCOL_WRONG_EVENT ReasonCode = "IN_PROTOCOL_WRONG_EVENT"
)
```

### Gender Codes

```go
const (
    HL7_ActAdministrativeGender_OID = "2.16.840.1.113883.5.1"

    GENDER_MALE            GenderCode = "M"
    GENDER_FEMALE          GenderCode = "F"
    GENDER_UNDIFFERENTIATED GenderCode = "UN"
)
```

### Race Codes

```go
const (
    HL7_Race_OID = "2.16.840.1.113883.5.104"

    RACE_WHITE                  RaceCode = "2106-3"
    RACE_BLACK                  RaceCode = "2054-5"
    RACE_ASIAN                  RaceCode = "2028-9"
    RACE_NATIVE_HAWAIIAN        RaceCode = "2076-8"
    RACE_AMERICAN_INDIAN        RaceCode = "1002-5"
    RACE_OTHER                  RaceCode = "2131-1"
)
```

## Examples

See the [main.go](/example/main.go) file for a complete working example that generates a full aECG XML document with:

- 12-lead rhythm series with realistic waveform data
- Device and manufacturer information
- 14+ annotations including:
  - Global measurements (heart rate, QRS duration, QT interval, axes)
  - Lead-specific measurements with nested annotations
  - Text annotations for ECG interpretation

Run the example:

```bash
go run main.go
```

Output will be generated at `/tmp/hl7aecg_example.xml`

## Testing

Run all tests:

```bash
go test ./...
```

Run with verbose output:

```bash
go test -v ./hl7aecg/types/...
```

Run with coverage:

```bash
go test -cover ./...
```

Generate coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Coverage

The library includes comprehensive tests:

- **Unit tests** for all core types and functions
- **Integration tests** for complete workflows
- **Validation tests** for HL7 compliance
- **Marshal/Unmarshal tests** for XML serialization

Current test coverage: **~85%**

## Project Structure

```
hl7v3-aecg/
├── hl7aecg/              # Main package
│   ├── builder.go        # Series builder methods
│   ├── init.go          # Main Hl7xml structure
│   ├── setter.go        # Simple setter methods
│   ├── subject.go       # Subject configuration
│   └── validation.go    # Validation entry point
│
├── hl7aecg/types/       # Type definitions and validation
│   ├── types_*.go       # HL7 data structures
│   ├── validator_*.go   # Validation logic
│   ├── code_systems.go  # OID constants and codes
│   ├── error.go         # Error definitions
│   ├── helper.go        # Helper functions
│   └── converters.go    # Type converters
│
├── main.go              # Complete example
└── README.md            # This file
```

## Standards Compliance

This library implements:

- **HL7 Version 3** Clinical Document Architecture (CDA)
- **HL7 aECG Implementation Guide** (Final 21-March-2005)
- **ISO/IEEE 11073-10101** Medical Device Codes (MDC)
- **CPT** Current Procedural Terminology codes
- **LOINC** Logical Observation Identifiers (for observations)

### XML Schema

The generated XML conforms to the HL7 aECG schema:

- Namespace: `urn:hl7-org:v3`

## Contributing

Contributions are welcome! Please follow these guidelines:

1. **Check the WIKI.md** for current implementation status
2. **Follow existing code conventions** and patterns
3. **Add comprehensive tests** for new features
4. **Update documentation** (README, WIKI, CLAUDE.md)
5. **Validate XML output** against HL7 schema

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Run `go test ./...` and `go vet ./...`
5. Update documentation
6. Submit a pull request

## Roadmap

### Completed ✅

- Complete HL7 v3 aECG document structure
- ECG waveform data (GLIST_TS, SLIST_PQ, SLIST_INT)
- Comprehensive annotation system (numeric and text)
- Subject demographics and clinical trial metadata
- Device and technician information
- Code systems (CPT, MDC, HL7 Act, Gender, Race)
- Validation framework
- ~85% test coverage

### In Progress 🚧

- Waveform data auto-generation from raw samples
- XML schema validation
- Additional code systems

### Planned 📋

- Related observations (medications, vital signs)
- Derivation support for derived series
- Additional series types (Holter, stress test)
- XML parsing (unmarshalling from XML files)
- Enhanced validation (cross-field consistency)

## License

[MIT License](LICENSE)

## Authors

- Jonathan Milhas - Initial work - LIRYC-IHU

## Acknowledgments

- HL7 International for the aECG specification
- FDA for clinical trial submission guidelines
- ISO/IEEE for Medical Device Codes

## Support

For questions, issues, or contributions:

- Open an issue on GitHub
- Refer to the [WIKI.md](WIKI.md) for implementation details

---

**Note**: This library is designed for research and clinical trial use. Ensure compliance with relevant regulations (FDA) for your specific use case.
