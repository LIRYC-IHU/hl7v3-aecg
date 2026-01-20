# Wiki - HL7 aECG Implementation Status

## Project Overview

This Go library implements the HL7 v3 standard for generating Annotated ECG (aECG) XML files compliant with the HL7 aECG Implementation Guide (Final 21-March-2005) for FDA clinical trial submissions.

**Main Package**: `github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg`

---

## ✅ Implemented Features

### 1. Basic aECG Document Structure

#### Root Document (AnnotatedECG)

- ✅ **ID** - Unique document identifier (UUID or OID)
- ✅ **Code** - CPT procedure code (e.g., ECG_Routine)
- ✅ **Text** - Narrative description
- ✅ **EffectiveTime** - Acquisition time interval (Low/High)
- ✅ **ConfidentialityCode** - Blinding status (S/I/B/C)
- ✅ **ReasonCode** - Protocol compliance
- ✅ **ClinicalTrial** - Clinical trial metadata
- ✅ **Subject** - Subject demographics
- ✅ **Component[]** - ECG series

### 2. Clinical Trial Information

#### ClinicalTrial

- ✅ **ID** - Trial identifier (OID + extension)
- ✅ **Title** - Human-readable trial name
- ✅ **ActivityTime** - Trial execution period
- ✅ Helper methods: `GetIdentifier()`, `GetDisplayName()`

#### ClinicalTrialProtocol

- ✅ **ID** - Unique protocol identifier
- ✅ **Title** - Protocol name

#### ClinicalTrialSponsor / SponsorOrganization

- ✅ **ID** - Organization identifier
- ✅ **Extension** - Additional identifier
- ✅ **Name** - Organization name

### 3. Subject Information

#### Subject / TrialSubject

- ✅ **ID** - Subject identifier
- ✅ **Code** - Role (SCREENING, ENROLLED)
- ✅ **SubjectDemographicPerson** - Demographics data

#### SubjectDemographicPerson

- ✅ **Name** - Initials or full name
- ✅ **AdministrativeGenderCode** - Gender (M/F/UN)
- ✅ **BirthTime** - Birth date (YYYYMMDD format)
- ✅ **RaceCode** - HL7 race codes

#### SubjectAssignment

- ✅ **Subject** - Subject identification
- ✅ **Definition** - Treatment group assignment
- ✅ **ComponentOf** - Reference to clinical trial

#### TreatmentGroupAssignment

- ✅ **Code** - Group code (e.g., GRP_001, GRP_002)

### 4. ECG Series and Waveforms

#### Series

- ✅ **ID** - Series identifier
- ✅ **Code** - Series type (RHYTHM, REPRESENTATIVE_BEAT, MEDIAN_BEAT)
- ✅ **EffectiveTime** - Acquisition time interval
- ✅ **Author** - Device information (SeriesAuthor)
- ✅ **SecondaryPerformer[]** - Technician information
- ✅ **Support** - Region of interest (for derived series)
- ✅ **Component[]** - Series components with sequences

#### SeriesAuthor / ManufacturedSeriesDevice

- ✅ **ID** - Device identifier (serial number)
- ✅ **Code** - Device type (12LEAD_ELECTROCARDIOGRAPH, 12LEAD_HOLTER)
- ✅ **ManufacturerModelName** - Model
- ✅ **SoftwareName** - Software version
- ✅ **ManufacturerOrganization** - Manufacturer information

#### SecondaryPerformer / SeriesPerformer

- ✅ **FunctionCode** - Technician role (HOLTER_HOOKUP_TECH, HOLTER_ANALYST, etc.)
- ✅ **Time** - Performance period
- ✅ **SeriesPerformer** - Technician identification
- ✅ **AssignedPerson** - Technician name

#### SeriesSupport / SupportingROI

- ✅ **Code** - ROI type (ROIPS - partially specified, ROIFS - fully specified)
- ✅ **Component[]** - ROI boundaries (ROIComponent)

#### Boundary / BoundaryValue

- ✅ **Code** - TIME_ABSOLUTE or lead code
- ✅ **Low/High** - Start/end indices

### 5. Sequences and Waveform Data

#### SequenceSet / SequenceComponent / Sequence

- ✅ **Code** - Dimension identifier (time or lead)
- ✅ **Value** - Waveform data (polymorphic)

#### Waveform Data Types

**GLIST_TS** - Generated List of Timestamps

- ✅ **Head** - First timestamp
- ✅ **Increment** - Time interval between samples

**SLIST_PQ** - Scaled List of Physical Quantities

- ✅ **Origin** - Baseline value (typically 0)
- ✅ **Scale** - Resolution per digit (e.g., 5 µV)
- ✅ **Digits** - Raw integer values separated by spaces
- ✅ Methods: `GetDigits()`, `GetLength()`, `GetActualValues()`

**SLIST_INT** - Scaled List of Integers

- ✅ **Origin** - Base integer
- ✅ **Scale** - Multiplication factor
- ✅ **Digits** - Integer string
- ✅ Methods: `GetDigits()`, `GetLength()`, `GetActualValues()`

**PhysicalQuantity / Increment**

- ✅ **Value** - Numeric string
- ✅ **Unit** - Unit of measure

### 6. Study Events and Timepoints

#### TimepointEvent

- ✅ **Code** - Visit code (e.g., VISIT_1, VISIT_2)
- ✅ **EffectiveTime** - Event time interval
- ✅ **ReasonCode** - S (scheduled), U (unscheduled)
- ✅ **Performer** - Event performer (optional)
- ✅ **ComponentOf** - Reference to subject assignment

#### RelativeTimepoint

- ✅ **Code** - Timepoint code
- ✅ **ComponentOf** - Protocol timing
- ✅ **PauseQuantity** - Delay from reference
- ✅ **ProtocolTimepointEvent** - Protocol definition
- ✅ **ReferenceEvent** - Reference event

### 7. Trial Site and Location

#### Location / TrialSite / SiteLocation

- ✅ **ID** - Site identifier
- ✅ **Name** - Site name
- ✅ **Addr** - Address (Address)
- ✅ **City/State/Country** - Address components

### 8. Investigator Information

#### ResponsibleParty / TrialInvestigator

- ✅ **ID** - Investigator identifier
- ✅ **InvestigatorPerson** - Personal information
- ✅ **PersonName** - Name components (Prefix, Given, Family, Suffix)

### 9. Code Systems

#### Implemented OIDs (100+ codes)

- ✅ **CPT_OID** - Procedure codes (2.16.840.1.113883.6.12)
- ✅ **MDC_OID** - Medical device codes / ECG leads (2.16.840.1.113883.6.24)
- ✅ **HL7_ActCode_OID** - HL7 action codes (2.16.840.1.113883.5.4)
- ✅ **HL7_ActAdministrativeGender_OID** - Gender codes (2.16.840.1.113883.5.1)
- ✅ **HL7_Race_OID** - Race codes (2.16.840.1.113883.5.104)

#### Standard 12-Lead ECG

- ✅ Limb leads: `MDC_ECG_LEAD_I`, `II`, `III`
- ✅ Augmented leads: `MDC_ECG_LEAD_AVR`, `AVL`, `AVF`
- ✅ Precordial leads: `MDC_ECG_LEAD_V1` through `V6`
- ✅ Helper function: `GetStandardLeads()`

### 10. Builder and Setter Methods

#### Main Methods (hl7aecg package)

- ✅ `NewHl7xml(outputDir)` - Create aECG instance
- ✅ `Initialize(code, codeSystem)` - Set CPT code
- ✅ `SetText(text)` - Set narrative text
- ✅ `SetEffectiveTime(low, high)` - Set time interval
- ✅ `SetSubject(id, ext, role)` - Set subject ID and role
- ✅ `SetSubjectDemographics(name, gender, birth, race)` - Set demographics
- ✅ `AddRhythmSeries(...)` - Add rhythm series
- ✅ `AddRepresentativeBeatSeries(...)` - Add representative beat series
- ✅ `SetSeriesAuthor(...)` - Set device information
- ✅ `Test()` - Write XML to /tmp/hl7aecg_example.xml

#### Types Package Methods

- ✅ `SetID(id, extension)` - Set ID (generates UUID if empty)
- ✅ `SetCode(code, system, display)` - Set code
- ✅ `SetName()`, `SetGender()`, `SetBirthDate()`, `SetRace()` - Subject demographics
- ✅ Helper functions for code systems

### 11. Annotations - ✅ FULLY IMPLEMENTED

**The annotation system is now completely implemented**, allowing you to add measurements, interpretations, and observations to ECG data.

#### AnnotationSet - ✅ IMPLEMENTED

Container for annotations created by a single author at a specific time.

**Implemented structures:**

- ✅ `AnnotationSet` with ID, Code, EffectiveTime, ActivityTime, Author, Component[]
- ✅ `AnnotationSetAuthor` (identical to SeriesAuthor)
- ✅ Support for device and person authors

**Usage:**

- Automatic annotations from ECG device
- Manual annotations from lab technician
- Annotations from cardiologist

#### Annotation - ✅ IMPLEMENTED

Individual observation on a series (e.g., P wave, R peak, QT interval).

**Implemented structures:**

- ✅ `Annotation` with Code, Value, Support, Component[]
- ✅ `AnnotationComponent` for nested annotations
- ✅ `AnnotationValue` - polymorphic value supporting PQ and ST types
- ✅ `AnnotationSupport` with SupportingROI for lead-specific annotations
- ✅ `AnnotationSupportingROI` with Boundary components
- ✅ `AnnotationBoundary` for lead identification

#### AnnotationValue - ✅ IMPLEMENTED (Polymorphic)

**Physical Quantities (PQ)** - Numeric values with units:

- ✅ XsiType: "PQ"
- ✅ Value: float64 (e.g., 72.0)
- ✅ Unit: string (e.g., "bpm", "ms", "mV")
- ✅ Methods: `IsPQ()`, `GetValueFloat()`, `GetValueUnit()`

**String Values (ST)** - Textual content:

- ✅ XsiType: "ST"
- ✅ Value: string (e.g., "Normal sinus rhythm")
- ✅ Methods: `IsST()`, `GetText()`

#### Builder Methods - ✅ IMPLEMENTED

**Global Measurements:**

```go
func (as *AnnotationSet) AddHeartRate(bpm float64) int
func (as *AnnotationSet) AddQRSDuration(durationMs float64) int
func (as *AnnotationSet) AddQTInterval(intervalMs float64) int
func (as *AnnotationSet) AddQTcInterval(intervalMs float64) int
```

**Numeric Annotations:**

```go
func (as *AnnotationSet) AddAnnotation(code, codeSystem string, value float64, unit string) int
func (as *AnnotationSet) AddAnnotationWithCodeSystemName(code, codeSystemName string, value float64, unit string) int
func (a *Annotation) AddNestedAnnotation(code, codeSystem string, value float64, unit string) int
func (a *Annotation) AddNestedAnnotationWithCodeSystemName(code, codeSystemName string, value float64, unit string) int
```

**Text Annotations:**

```go
func (as *AnnotationSet) AddTextAnnotation(code, codeSystem, text string) int
func (as *AnnotationSet) AddTextAnnotationWithCodeSystemName(code, codeSystemName, text string) int
func (a *Annotation) AddNestedTextAnnotation(code, codeSystem, text string) int
func (a *Annotation) AddNestedTextAnnotationWithCodeSystemName(code, codeSystemName, text string) int
```

**Lead-Specific Annotations:**

```go
func (as *AnnotationSet) AddLeadAnnotation(leadCode, code, codeSystem, codeSystemName string) int
```

**Accessors:**

```go
func (as *AnnotationSet) GetAnnotation(idx int) *Annotation
func (as *AnnotationSet) GetAnnotationByCode(code string) *Annotation
func (a *Annotation) GetNestedAnnotation(idx int) *Annotation
func (a *Annotation) GetValueFloat() (float64, bool)
```

#### Example Usage:

```go
// Numeric annotation (heart rate)
annSet.AddHeartRate(72)

// Text annotation (ECG interpretation)
interpIdx := annSet.AddTextAnnotation("MDC_ECG_INTERPRETATION", MDC_OID, "")
interp := annSet.GetAnnotation(interpIdx)
interp.AddNestedTextAnnotation("MDC_ECG_INTERPRETATION_STATEMENT", MDC_OID, "Normal sinus rhythm")

// Lead-specific annotation
leadIdx := annSet.AddLeadAnnotation("MDC_ECG_LEAD_II", "VENDOR_MATRIX", "VENDOR_MATRIX", "VENDOR")
leadAnn := annSet.GetAnnotation(leadIdx)
leadAnn.AddNestedAnnotationWithCodeSystemName("VENDOR_P_ONSET", "VENDOR", 234, "ms")
```

### 12. Validation

#### Implemented Validation

- ✅ ID validation (UUID or OID format)
- ✅ Code and CodeSystem validation
- ✅ EffectiveTime validation (at least Low or High must exist)
- ✅ Time format validation (HL7 TS or Unix timestamp)
- ✅ ConfidentialityCode validation (S/I/B/C)
- ✅ ReasonCode validation (PER_PROTOCOL/NOT_IN_PROTOCOL/IN_PROTOCOL_WRONG_EVENT)
- ✅ Subject presence validation (required)
- ✅ GenderCode validation (M/F/UN)
- ✅ TrialSubjectCode validation (SCREENING/ENROLLED)
- ✅ Basic SubjectAssignment and ClinicalTrial validation
- ✅ Annotation validation (Code, Value for PQ and ST types)
- ✅ AnnotationSet validation (Code, Component)
- ✅ SupportingROI validation (ClassCode, boundary codes)

---

## ❌ Missing Features

### 1. 🔴 CRITICAL: RelatedObservation (High Priority)

Related observations to the main document (e.g., subject age, medications, etc.).

**Missing structures:**

```go
type RelatedObservation struct {
    XMLName   xml.Name                  `xml:"relatedObservation"`
    Code      *Code[ObservationCode]    `xml:"code,omitempty"`
    Value     *ObservationValue         `xml:"value,omitempty"`
    Author    *RelatedObservationAuthor `xml:"author,omitempty"`
}

type RelatedObservationAuthor struct {
    // Can be a person or device
    Time                      *EffectiveTime              `xml:"time,omitempty"`
    AssignedAuthor            *AssignedAuthor             `xml:"assignedAuthor,omitempty"`
    ManufacturedObservationDevice *ManufacturedObservationDevice `xml:"assignedObservationDevice/manufacturingObservationDevice,omitempty"`
}

type ObservationValue struct {
    XsiType string      `xml:"xsi:type,attr"`
    Value   interface{} `xml:",chardata"`
    Unit    *string     `xml:"unit,attr,omitempty"`
}
```

**Missing observation codes:**

```go
// Demographic observations
OBS_AGE          = "21612-7" // Age (LOINC)
OBS_HEIGHT       = "8302-2"  // Height
OBS_WEIGHT       = "29463-7" // Weight

// Clinical observations
OBS_HEART_RATE   = "8867-4"  // Heart rate
OBS_BP_SYSTOLIC  = "8480-6"  // Systolic blood pressure
OBS_BP_DIASTOLIC = "8462-4"  // Diastolic blood pressure
```

**Missing builder methods:**

```go
func (h *Hl7xml) AddRelatedObservation(code ObservationCode, value interface{}, unit string) *Hl7xml
func (h *Hl7xml) AddSubjectAge(ageYears int) *Hl7xml
func (h *Hl7xml) AddVitalSign(code ObservationCode, value float64, unit string) *Hl7xml
```

---

**Missing validation to implement:**

```go
// In validator_sequence.go
func (s *SequenceSet) ValidateLengthConsistency(ctx context.Context, vctx *ValidationContext)
func (s *SLIST_PQ) ValidateVoltageRange(ctx context.Context, vctx *ValidationContext)
func (s *SLIST_PQ) ValidateDigitsFormat(ctx context.Context, vctx *ValidationContext)
func (g *GLIST_TS) ValidateTimestampFormat(ctx context.Context, vctx *ValidationContext)
```

---

### 2. 🟢 Minor Improvements (Low Priority)

#### Identified Minor Bugs

- ⚠️ **SponsorOrganization** has redundant fields `ID.Extension` and `Extension` (`types_clinicalTrialSponsor.go:66-73`)
- ⚠️ **Test()** writes to `/tmp/hl7aecg_example.xml` instead of using configured `outputDir`
- ⚠️ XML namespaces not automatically generated for nested elements

#### Additional Series Types

- ✅ RHYTHM - implemented
- ✅ REPRESENTATIVE_BEAT - implemented
- ⚠️ MEDIAN_BEAT - defined but incomplete support

#### Improved Error Reporting

- ❌ Path context in validation errors (e.g., `"HL7AEcg.Subject.TrialSubject.ID"`)
- ❌ Multilingual error messages (currently English only)
- ❌ Correction suggestions in errors

#### Additional Helper Methods

```go
// Conversion helpers
func ConvertSampleRateToIncrement(sampleRate float64) string
func ConvertTimestampToHL7(t time.Time) string
func ParseHL7Timestamp(hl7ts string) (time.Time, error)

// Validation helpers
func IsValidHL7Timestamp(ts string) bool
func IsValidUUID(uuid string) bool
func IsValidOID(oid string) bool
```

---

## 📊 Implementation Statistics

### Data Structures

- ✅ **Implemented**: 51+ structures
- ❌ **Missing**: 5 structures (RelatedObservation, etc.)
- 📈 **Completion Rate**: ~95%

### Methods

- ✅ **Builder/Setter**: 45+ implemented methods
- ❌ **Missing**: ~10 methods (observations, parsing)
- 📈 **Completion Rate**: ~90%

### Code Systems

- ✅ **Defined OIDs**: 5+ main code systems
- ✅ **Defined Codes**: 100+ codes (CPT, MDC, HL7 Act, Gender, Race)
- ❌ **Missing**: ~20 MDC codes for wave annotations
- 📈 **Completion Rate**: ~85%

### Validation

- ✅ **Implemented Checks**: ~35 validations
- ❌ **Missing Checks**: ~15+ validations
- ⚠️ **Bugs**: 2 identified bugs
- 📈 **Completion Rate**: ~70%

---

## 🎯 Recommended Roadmap

### Phase 1: Critical Features (COMPLETED ✅)

1. ✅ Implement **AnnotationSet** and **Annotation**
2. ✅ Add annotation builder methods
3. ✅ Create helper methods for common annotations
4. ✅ Implement polymorphic AnnotationValue (PQ and ST)
5. ✅ Add comprehensive annotation validation
6. ✅ Create complete test suite for annotations

### Phase 2: Remaining Critical Features (1-2 weeks)

1. ✨ Implement **RelatedObservation**
2. ✨ Add LOINC observation codes
3. 🐛 Fix identified validation bugs

### Phase 3: Robust Validation (1-2 weeks)

1. ✨ Sequence length consistency validation
2. ✨ Range validation (voltage, heart rate)
3. ✨ Strict format validation (timestamps, OID, UUID)
4. ✨ Improve error context and messages

### Phase 4: Additional Features (1-2 weeks)

1. ✨ Implement parsing/reading of aECG XML files
2. ✨ Full support for MEDIAN_BEAT
3. ✨ Conversion helper methods
4. ✨ Extended documentation and examples

### Phase 5: Polish (1 week)

1. 🐛 Fix minor bugs (outputDir, redundant fields)
2. ✨ Exhaustive unit tests
3. ✨ Performance benchmarks
4. ✨ Complete API documentation

---

## 📚 References

- **Implementation Guide**: HL7 aECG Implementation Guide (Final 21-March-2005)
- **HL7 v3 Standard**: Clinical Document Architecture (CDA)
- **MDC Codes**: ISO/IEEE 11073-10101 Medical Device Codes
- **LOINC**: Logical Observation Identifiers Names and Codes
- **CPT**: Current Procedural Terminology

---

## 🤝 Contributing

To contribute to this project:

1. Check this wiki for current implementation status
2. Choose a missing feature from the roadmap
3. Follow existing code conventions
4. Add unit tests
5. Update this wiki

---

_Last updated: 2026-01-20_
_Version: 1.1_
_Major update: Annotations fully implemented (numeric and text)_
