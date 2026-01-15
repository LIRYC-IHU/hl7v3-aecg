package types

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
//	  <code code="MINDRAY_MEASUREMENT_MATRIX" codeSystem="" codeSystemName="MINDRAY"/>
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
