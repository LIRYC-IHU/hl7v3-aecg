package types

import "strconv"

// =============================================================================
// Annotation Types
// =============================================================================
// Reference: HL7 aECG Implementation Guide, Section 6: Annotations

// AnnotationSet contains ECG annotations for a series.
//
// Annotations include measurements such as heart rate, PR interval, QRS duration,
// QT interval, and lead-specific measurements. The annotationSet groups all
// annotations that share the same activity time.
//
// XML Structure:
//
//	<subjectOf>
//	  <annotationSet>
//	    <activityTime value="20250923103600"/>
//	    <component>
//	      <annotation>
//	        <code code="MDC_ECG_HEART_RATE" codeSystem="2.16.840.1.113883.6.24"/>
//	        <value xsi:type="PQ" value="57" unit="bpm"/>
//	      </annotation>
//	    </component>
//	  </annotationSet>
//	</subjectOf>
//
// Cardinality: Optional (within SubjectOf)
// Reference: HL7 aECG Implementation Guide
type AnnotationSet struct {
	// ActivityTime is the time when the annotations were made or the measurements were taken.
	//
	// XML Tag: <activityTime value="..."/>
	// Cardinality: Optional
	ActivityTime *Time `xml:"activityTime,omitempty"`

	// Component contains the individual annotations.
	//
	// XML Tag: <component>...</component>
	// Cardinality: Optional (0..*)
	Component []AnnotationComponent `xml:"component,omitempty"`
}

// AnnotationComponent wraps an Annotation to provide the correct XML structure.
//
// XML Structure:
//
//	<component>
//	  <annotation>...</annotation>
//	</component>
//
// Cardinality: Optional (0..* within AnnotationSet or Annotation)
type AnnotationComponent struct {
	// Annotation contains the measurement or observation.
	//
	// XML Tag: <annotation>...</annotation>
	// Cardinality: Required (within AnnotationComponent)
	Annotation Annotation `xml:"annotation"`
}

// Annotation represents an ECG measurement or observation.
//
// Annotations can be:
// - Global measurements (heart rate, PR, QRS, QT intervals)
// - Lead-specific measurements (amplitude, timing per lead)
// - Nested annotations (e.g., QTc with correction method sub-annotations)
//
// XML Structure (simple):
//
//	<annotation>
//	  <code code="MDC_ECG_HEART_RATE" codeSystem="2.16.840.1.113883.6.24"/>
//	  <value xsi:type="PQ" value="57" unit="bpm"/>
//	</annotation>
//
// XML Structure (lead-specific with support):
//
//	<annotation>
//	  <code code="MEASUREMENT_MATRIX" codeSystem="" codeSystemName=""/>
//	  <support>
//	    <supportingROI classCode="ROIBND">
//	      <code code="ROIPS" codeSystem="2.16.840.1.113883.5.4"/>
//	      <component>
//	        <boundary>
//	          <code code="MDC_ECG_LEAD_I" codeSystem="2.16.840.1.113883.6.24"/>
//	        </boundary>
//	      </component>
//	    </supportingROI>
//	  </support>
//	  <component>
//	    <annotation>...</annotation>
//	  </component>
//	</annotation>
//
// Cardinality: Required (within AnnotationComponent)
// Reference: HL7 aECG Implementation Guide
type Annotation struct {
	// Code identifies the type of annotation or measurement.
	//
	// Common MDC Codes:
	//   - "MDC_ECG_HEART_RATE": Heart rate (bpm)
	//   - "MDC_ECG_TIME_PD_PR": PR interval (ms)
	//   - "MDC_ECG_TIME_PD_QRS": QRS duration (ms)
	//   - "MDC_ECG_TIME_PD_QT": QT interval (ms)
	//   - "MDC_ECG_TIME_PD_QTc": Corrected QT interval (ms)
	//
	// XML Tag: <code code="..." codeSystem="..." codeSystemName="..." displayName="..."/>
	// Cardinality: Optional
	Code *Code[string, string] `xml:"code,omitempty"`

	// Value contains the measurement value as a Physical Quantity.
	//
	// XML Tag: <value xsi:type="PQ" value="..." unit="..."/>
	// Cardinality: Optional
	Value *PhysicalQuantity `xml:"value,omitempty"`

	// Support defines the region of interest for lead-specific annotations.
	//
	// When present, the annotation applies to a specific lead or set of leads
	// identified by the supporting ROI.
	//
	// XML Tag: <support>...</support>
	// Cardinality: Optional
	Support *AnnotationSupport `xml:"support,omitempty"`

	// Component contains nested annotations.
	//
	// Used for complex annotations that have sub-measurements.
	// Example: QTc can have nested annotations for different correction formulas.
	//
	// XML Tag: <component>...</component>
	// Cardinality: Optional (0..*)
	Component []AnnotationComponent `xml:"component,omitempty"`
}

// AnnotationSupport contains supporting region of interest information for annotations.
//
// Used to specify which lead(s) an annotation applies to.
//
// XML Structure:
//
//	<support>
//	  <supportingROI classCode="ROIBND">
//	    <code code="ROIPS" codeSystem="2.16.840.1.113883.5.4"/>
//	    <component>
//	      <boundary>
//	        <code code="MDC_ECG_LEAD_I" codeSystem="2.16.840.1.113883.6.24"/>
//	      </boundary>
//	    </component>
//	  </supportingROI>
//	</support>
//
// Cardinality: Optional (within Annotation)
type AnnotationSupport struct {
	// SupportingROI identifies the region (typically a lead) for the annotation.
	//
	// XML Tag: <supportingROI classCode="...">...</supportingROI>
	// Cardinality: Required (within AnnotationSupport)
	SupportingROI AnnotationSupportingROI `xml:"supportingROI"`
}

// AnnotationSupportingROI defines the region of interest for lead-specific annotations.
//
// This is distinct from the series-level SupportingROI and is specifically
// used within annotation contexts to identify which lead(s) the measurements apply to.
//
// XML Structure:
//
//	<supportingROI classCode="ROIBND">
//	  <code code="ROIPS" codeSystem="2.16.840.1.113883.5.4"/>
//	  <component>
//	    <boundary>
//	      <code code="MDC_ECG_LEAD_I" codeSystem="2.16.840.1.113883.6.24"/>
//	    </boundary>
//	  </component>
//	</supportingROI>
//
// Cardinality: Required (within AnnotationSupport)
type AnnotationSupportingROI struct {
	// ClassCode specifies the ROI class, typically "ROIBND" (Region of Interest Bounded).
	//
	// XML Tag: classCode="..."
	// Cardinality: Optional
	ClassCode string `xml:"classCode,attr,omitempty"`

	// Code specifies whether the ROI is fully or partially specified.
	//
	// Values:
	//   - "ROIPS": Partially specified ROI
	//   - "ROIFS": Fully specified ROI
	//
	// XML Tag: <code code="..." codeSystem="..."/>
	// Cardinality: Optional
	Code *Code[string, string] `xml:"code,omitempty"`

	// Component contains the boundary definitions for the ROI.
	//
	// XML Tag: <component>...</component>
	// Cardinality: Optional (0..*)
	Component []AnnotationBoundaryComponent `xml:"component,omitempty"`
}

// AnnotationBoundaryComponent wraps an AnnotationBoundary.
//
// XML Structure:
//
//	<component>
//	  <boundary>
//	    <code code="MDC_ECG_LEAD_I" codeSystem="2.16.840.1.113883.6.24"/>
//	  </boundary>
//	</component>
//
// Cardinality: Optional (0..* within AnnotationSupportingROI)
type AnnotationBoundaryComponent struct {
	// Boundary identifies the lead for lead-specific annotations.
	//
	// XML Tag: <boundary>...</boundary>
	// Cardinality: Required (within AnnotationBoundaryComponent)
	Boundary AnnotationBoundary `xml:"boundary"`
}

// AnnotationBoundary identifies a specific lead for lead-specific annotations.
//
// XML Structure:
//
//	<boundary>
//	  <code code="MDC_ECG_LEAD_I" codeSystem="2.16.840.1.113883.6.24" codeSystemName="MDC"/>
//	</boundary>
//
// Cardinality: Required (within AnnotationBoundaryComponent)
type AnnotationBoundary struct {
	// Code identifies the lead.
	//
	// Common MDC Lead Codes:
	//   - "MDC_ECG_LEAD_I", "MDC_ECG_LEAD_II", "MDC_ECG_LEAD_III"
	//   - "MDC_ECG_LEAD_AVR", "MDC_ECG_LEAD_AVL", "MDC_ECG_LEAD_AVF"
	//   - "MDC_ECG_LEAD_V1" through "MDC_ECG_LEAD_V6"
	//
	// XML Tag: <code code="..." codeSystem="..." codeSystemName="..."/>
	// Cardinality: Required
	Code Code[string, string] `xml:"code"`
}

// =============================================================================
// Helper Functions for Annotations
// =============================================================================

// GetAnnotationByCode finds an annotation with the given code in the annotation set.
// Returns nil if not found.
func (as *AnnotationSet) GetAnnotationByCode(code string) *Annotation {
	if as == nil {
		return nil
	}
	for i := range as.Component {
		if as.Component[i].Annotation.Code != nil &&
			as.Component[i].Annotation.Code.Code == code {
			return &as.Component[i].Annotation
		}
	}
	return nil
}

// GetLeadAnnotations finds the annotation with measurements for a specific lead.
// Returns nil if not found.
func (as *AnnotationSet) GetLeadAnnotations(leadCode string) *Annotation {
	if as == nil {
		return nil
	}
	for i := range as.Component {
		ann := &as.Component[i].Annotation
		if ann.Support != nil {
			for _, comp := range ann.Support.SupportingROI.Component {
				if string(comp.Boundary.Code.Code) == leadCode {
					return ann
				}
			}
		}
	}
	return nil
}

// GetValueFloat parses the annotation value as a float64.
// Returns 0 and false if the annotation has no value or parsing fails.
func (a *Annotation) GetValueFloat() (float64, bool) {
	if a == nil || a.Value == nil {
		return 0, false
	}
	return a.Value.GetValueFloat()
}

// GetValueUnit returns the unit of the annotation value.
// Returns empty string if the annotation has no value.
func (a *Annotation) GetValueUnit() string {
	if a == nil || a.Value == nil {
		return ""
	}
	return a.Value.Unit
}

// GetNestedAnnotationByCode finds a nested annotation with the given code.
// Returns nil if not found.
func (a *Annotation) GetNestedAnnotationByCode(code string) *Annotation {
	if a == nil {
		return nil
	}
	for i := range a.Component {
		nested := &a.Component[i].Annotation
		if nested.Code != nil && nested.Code.Code == code {
			return nested
		}
	}
	return nil
}

// =============================================================================
// Builder Methods for Creating Annotations
// =============================================================================

// AddAnnotation adds a global annotation (without supportingROI) to the annotation set.
//
// Parameters:
//   - code: The annotation code (e.g., "MDC_ECG_HEART_RATE")
//   - codeSystem: The code system OID (e.g., "2.16.840.1.113883.6.24" for MDC)
//   - value: The numeric value as a float64
//   - unit: The unit of measurement (e.g., "bpm", "ms", "uV")
//
// Returns:
//   - *Annotation: Pointer to the newly created annotation for further customization
//
// Example:
//
//	ann := annotationSet.AddAnnotation("MDC_ECG_HEART_RATE", "2.16.840.1.113883.6.24", 57, "bpm")
func (as *AnnotationSet) AddAnnotation(code, codeSystem string, value float64, unit string) *Annotation {
	if as == nil {
		return nil
	}

	ann := Annotation{
		Code: &Code[string, string]{
			Code:       code,
			CodeSystem: codeSystem,
		},
		Value: &PhysicalQuantity{
			XsiType: "PQ",
			Value:   formatFloat(value),
			Unit:    unit,
		},
	}

	as.Component = append(as.Component, AnnotationComponent{Annotation: ann})
	return &as.Component[len(as.Component)-1].Annotation
}

// AddAnnotationWithCodeSystemName adds a global annotation with codeSystemName instead of codeSystem.
//
// Used for vendor-specific codes (e.g., ) that don't have an OID.
//
// Parameters:
//   - code: The annotation code
//   - codeSystemName: The code system name (e.g., "")
//   - value: The numeric value
//   - unit: The unit of measurement
//
// Returns:
//   - *Annotation: Pointer to the newly created annotation
func (as *AnnotationSet) AddAnnotationWithCodeSystemName(code, codeSystemName string, value float64, unit string) *Annotation {
	if as == nil {
		return nil
	}

	ann := Annotation{
		Code: &Code[string, string]{
			Code:           code,
			CodeSystem:     "", // Empty for vendor-specific codes
			CodeSystemName: codeSystemName,
		},
		Value: &PhysicalQuantity{
			XsiType: "PQ",
			Value:   formatFloat(value),
			Unit:    unit,
		},
	}

	as.Component = append(as.Component, AnnotationComponent{Annotation: ann})
	return &as.Component[len(as.Component)-1].Annotation
}

// AddLeadAnnotation adds a lead-specific annotation with supportingROI.
//
// Creates an annotation that applies to a specific ECG lead. The annotation includes
// a supportingROI structure that identifies which lead the measurements apply to.
//
// Parameters:
//   - leadCode: The lead code (e.g., "MDC_ECG_LEAD_I", "MDC_ECG_LEAD_V1")
//   - matrixCode: The measurement matrix code (typically "MMEASUREMENT_MATRIX")
//   - codeSystemName: The code system name (e.g., "")
//
// Returns:
//   - *Annotation: Pointer to the lead annotation for adding nested measurements
//
// Example:
//
//	leadAnn := annotationSet.AddLeadAnnotation("MDC_ECG_LEAD_I", "MMEASUREMENT_MATRIX", "")
//	leadAnn.AddNestedAnnotation("P_ONSET", "", 234, "ms")
//	leadAnn.AddNestedAnnotation("R_AMP", "", 535, "uV")
func (as *AnnotationSet) AddLeadAnnotation(leadCode, matrixCode, codeSystem, codeSystemName string) *Annotation {
	if as == nil {
		return nil
	}

	ann := Annotation{
		Code: &Code[string, string]{
			Code:           matrixCode,
			CodeSystem:     codeSystem, // Empty for vendor-specific codes
			CodeSystemName: codeSystemName,
		},
		Support: &AnnotationSupport{
			SupportingROI: AnnotationSupportingROI{
				ClassCode: "ROIBND",
				Code: &Code[string, string]{
					Code:       string(ROIPS),
					CodeSystem: string(HL7_ActCode_OID),
				},
				Component: []AnnotationBoundaryComponent{{
					Boundary: AnnotationBoundary{
						Code: Code[string, string]{
							Code:       leadCode,
							CodeSystem: string(MDC_OID),
						},
					},
				}},
			},
		},
	}

	as.Component = append(as.Component, AnnotationComponent{Annotation: ann})
	return &as.Component[len(as.Component)-1].Annotation
}

// AddNestedAnnotation adds a nested annotation to an existing annotation.
//
// Used for complex annotations that have sub-measurements, such as:
// - QTc with different correction formulas
// - Lead-specific measurements within a MEASUREMENT_MATRIX
//
// Parameters:
//   - code: The nested annotation code
//   - codeSystem: The code system OID (can be empty for vendor codes)
//   - value: The numeric value
//   - unit: The unit of measurement
//
// Returns:
//   - *Annotation: Pointer to the nested annotation for further nesting if needed
//
// Example:
//
//	qtcAnn := annotationSet.AddAnnotation("MDC_ECG_TIME_PD_QTc", MDC_OID, 0, "")
//	qtcAnn.AddNestedAnnotation("ECG_TIME_PD_QTcH", "", 413, "ms")
func (a *Annotation) AddNestedAnnotation(code, codeSystem string, value float64, unit string) *Annotation {
	if a == nil {
		return nil
	}

	nested := Annotation{
		Code: &Code[string, string]{
			Code:       code,
			CodeSystem: codeSystem,
		},
		Value: &PhysicalQuantity{
			XsiType: "PQ",
			Value:   formatFloat(value),
			Unit:    unit,
		},
	}

	a.Component = append(a.Component, AnnotationComponent{Annotation: nested})
	return &a.Component[len(a.Component)-1].Annotation
}

// AddNestedAnnotationWithCodeSystemName adds a nested annotation with codeSystemName.
func (a *Annotation) AddNestedAnnotationWithCodeSystemName(code, codeSystemName string, value float64, unit string) *Annotation {
	if a == nil {
		return nil
	}

	nested := Annotation{
		Code: &Code[string, string]{
			Code:           code,
			CodeSystem:     "",
			CodeSystemName: codeSystemName,
		},
		Value: &PhysicalQuantity{
			XsiType: "PQ",
			Value:   formatFloat(value),
			Unit:    unit,
		},
	}

	a.Component = append(a.Component, AnnotationComponent{Annotation: nested})
	return &a.Component[len(a.Component)-1].Annotation
}

// =============================================================================
// Convenience Methods for Common Annotations
// =============================================================================

// AddHeartRate adds a heart rate annotation in beats per minute.
func (as *AnnotationSet) AddHeartRate(value float64) *Annotation {
	return as.AddAnnotation(string(MDC_ECG_HEART_RATE), string(MDC_OID), value, "bpm")
}

// AddPRInterval adds a PR interval annotation in milliseconds.
func (as *AnnotationSet) AddPRInterval(value float64) *Annotation {
	return as.AddAnnotation(string(MDC_ECG_TIME_PD_PR), string(MDC_OID), value, "ms")
}

// AddQRSDuration adds a QRS duration annotation in milliseconds.
func (as *AnnotationSet) AddQRSDuration(value float64) *Annotation {
	return as.AddAnnotation(string(MDC_ECG_TIME_PD_QRS), string(MDC_OID), value, "ms")
}

// AddQTInterval adds a QT interval annotation in milliseconds.
func (as *AnnotationSet) AddQTInterval(value float64) *Annotation {
	return as.AddAnnotation(string(MDC_ECG_TIME_PD_QT), string(MDC_OID), value, "ms")
}

// AddQTcInterval adds a QTc interval annotation in milliseconds.
// Returns the annotation so nested correction formulas can be added.
func (as *AnnotationSet) AddQTcInterval(value float64) *Annotation {
	return as.AddAnnotation(string(MDC_ECG_TIME_PD_QTc), string(MDC_OID), value, "ms")
}

// formatFloat converts a float64 to a string for PhysicalQuantity.Value.
// Removes trailing zeros and decimal point if not needed.
func formatFloat(f float64) string {
	// Check if it's an integer value
	if f == float64(int64(f)) {
		return strconv.FormatInt(int64(f), 10)
	}
	// Use standard formatting for floats
	return strconv.FormatFloat(f, 'f', -1, 64)
}
