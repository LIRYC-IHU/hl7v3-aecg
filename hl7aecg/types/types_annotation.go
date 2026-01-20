package types

import (
	"encoding/xml"
	"strconv"
)

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

	// Value contains the measurement value (either Physical Quantity or String).
	//
	// The value type depends on the annotation type:
	//   - Numeric measurements: xsi:type="PQ" (PhysicalQuantity)
	//   - Text statements: xsi:type="ST" (StringValue)
	//
	// XML Tag: <value xsi:type="...">...</value>
	// Cardinality: Optional
	Value *AnnotationValue `xml:"value,omitempty"`

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

// StringValue represents a string text value for annotations.
//
// Used for interpretation statements and other textual annotations.
//
// XML Structure:
//
//	<value xsi:type="ST">Rythme sinusal avec ESA</value>
//
// Cardinality: Optional (within Annotation)
// XML Attribute: xsi:type="ST"
type StringValue struct {
	// XsiType specifies the type as "ST" for String.
	//
	// XML Tag: xsi:type="ST"
	// Cardinality: Required
	XsiType string `xml:"xsi:type,attr"`

	// Value is the text content.
	//
	// XML Tag: (text content)
	// Cardinality: Required
	Value string `xml:",chardata"`
}

// AnnotationValue represents a polymorphic value that can be either a PhysicalQuantity or StringValue.
//
// This type handles the xsi:type discrimination for annotation values.
//
// Supported types:
//   - xsi:type="PQ": PhysicalQuantity (numeric with unit)
//   - xsi:type="ST": StringValue (text content)
//
// XML Structure (PQ):
//
//	<value xsi:type="PQ" value="57" unit="bpm"/>
//
// XML Structure (ST):
//
//	<value xsi:type="ST">Rythme sinusal avec ESA</value>
type AnnotationValue struct {
	XMLName xml.Name `xml:"value"`
	XsiType string   `xml:"xsi:type,attr,omitempty"`

	RawXML []byte `xml:",innerxml"`
	// Typed holds the decoded value as one of:
	//   - *PhysicalQuantity (xsi:type="PQ")
	//   - *StringValue (xsi:type="ST")
	Typed any `xml:"-"`
}

// UnmarshalXML decodes the annotation value based on xsi:type.
func (av *AnnotationValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// Extract xsi:type attribute
	for _, attr := range start.Attr {
		// Check for xsi:type attribute (namespace can be empty or xsi)
		if attr.Name.Local == "type" || (attr.Name.Space == "xsi" && attr.Name.Local == "type") {
			av.XsiType = attr.Value
			break
		}
	}

	// Decode based on type
	switch av.XsiType {
	case "PQ":
		// For PhysicalQuantity, attributes are on the element itself
		pq := &PhysicalQuantity{
			XsiType: "PQ",
		}
		for _, attr := range start.Attr {
			switch attr.Name.Local {
			case "value":
				pq.Value = attr.Value
			case "unit":
				pq.Unit = attr.Value
			}
		}
		av.Typed = pq
		// Consume the end element
		return d.Skip()

	case "ST":
		// For String, content is character data
		var text string
		if err := d.DecodeElement(&text, &start); err != nil {
			return err
		}
		st := &StringValue{
			XsiType: "ST",
			Value:   text,
		}
		av.Typed = st
		return nil

	default:
		// Unknown type - skip the element
		return d.Skip()
	}
}

// MarshalXML encodes the annotation value based on its typed content.
func (av *AnnotationValue) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "value"}

	// Add xsi:type attribute
	if av.XsiType != "" {
		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: "xsi:type"},
			Value: av.XsiType,
		})
	}

	// Encode based on typed value
	switch typed := av.Typed.(type) {
	case *PhysicalQuantity:
		// For PQ, we need to encode as attributes
		if typed.Value != "" {
			start.Attr = append(start.Attr, xml.Attr{
				Name:  xml.Name{Local: "value"},
				Value: typed.Value,
			})
		}
		if typed.Unit != "" {
			start.Attr = append(start.Attr, xml.Attr{
				Name:  xml.Name{Local: "unit"},
				Value: typed.Unit,
			})
		}
		return e.EncodeElement("", start)

	case *StringValue:
		// For ST, encode as character data
		return e.EncodeElement(typed.Value, start)

	default:
		// If no typed value, encode empty element
		return e.EncodeElement("", start)
	}
}

// GetValueFloat returns the numeric value if this is a PhysicalQuantity.
// Returns (value, true) if successful, (0, false) otherwise.
func (av *AnnotationValue) GetValueFloat() (float64, bool) {
	if pq, ok := av.Typed.(*PhysicalQuantity); ok {
		return pq.GetValueFloat()
	}
	return 0, false
}

// GetValueUnit returns the unit if this is a PhysicalQuantity.
// Returns empty string if this is not a PQ or if unit is not set.
func (av *AnnotationValue) GetValueUnit() string {
	if pq, ok := av.Typed.(*PhysicalQuantity); ok {
		return pq.Unit
	}
	return ""
}

// GetText returns the text content if this is a StringValue.
// Returns (text, true) if successful, ("", false) otherwise.
func (av *AnnotationValue) GetText() (string, bool) {
	if st, ok := av.Typed.(*StringValue); ok {
		return st.Value, true
	}
	return "", false
}

// IsPQ returns true if this annotation value is a PhysicalQuantity.
func (av *AnnotationValue) IsPQ() bool {
	_, ok := av.Typed.(*PhysicalQuantity)
	return ok
}

// IsST returns true if this annotation value is a StringValue.
func (av *AnnotationValue) IsST() bool {
	_, ok := av.Typed.(*StringValue)
	return ok
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
// Returns empty string if the annotation has no value or is not a PQ.
func (a *Annotation) GetValueUnit() string {
	if a == nil || a.Value == nil {
		return ""
	}
	return a.Value.GetValueUnit()
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
//   - int: Index of the newly created annotation (use GetAnnotation to retrieve safely)
//
// Example:
//
//	idx := annotationSet.AddAnnotation("MDC_ECG_HEART_RATE", "2.16.840.1.113883.6.24", 57, "bpm")
//	ann := annotationSet.GetAnnotation(idx)
func (as *AnnotationSet) AddAnnotation(code, codeSystem string, value float64, unit string) int {
	if as == nil {
		return -1
	}

	// Input validation
	if code == "" {
		return -1
	}
	if codeSystem == "" {
		return -1
	}
	if unit == "" {
		return -1
	}
	if isInvalidFloat(value) {
		return -1
	}

	ann := Annotation{
		Code: &Code[string, string]{
			Code:       code,
			CodeSystem: codeSystem,
		},
		Value: &AnnotationValue{
			XsiType: "PQ",
			Typed: &PhysicalQuantity{
				XsiType: "PQ",
				Value:   formatFloat(value),
				Unit:    unit,
			},
		},
	}

	as.Component = append(as.Component, AnnotationComponent{Annotation: ann})
	return len(as.Component) - 1
}

// AddAnnotationWithCodeSystemName adds a global annotation with codeSystemName instead of codeSystem.
//
// Used for vendor-specific codes (e.g., MINDRAY, GE, Philips) that don't have an OID.
//
// Parameters:
//   - code: The annotation code
//   - codeSystemName: The code system name (e.g., "MINDRAY")
//   - value: The numeric value
//   - unit: The unit of measurement
//
// Returns:
//   - int: Index of the newly created annotation (use GetAnnotation to retrieve safely)
func (as *AnnotationSet) AddAnnotationWithCodeSystemName(code, codeSystemName string, value float64, unit string) int {
	if as == nil {
		return -1
	}

	// Input validation
	if code == "" {
		return -1
	}
	if codeSystemName == "" {
		return -1
	}
	if unit == "" {
		return -1
	}
	if isInvalidFloat(value) {
		return -1
	}

	ann := Annotation{
		Code: &Code[string, string]{
			Code:           code,
			CodeSystem:     "", // Empty for vendor-specific codes
			CodeSystemName: codeSystemName,
		},
		Value: &AnnotationValue{
			XsiType: "PQ",
			Typed: &PhysicalQuantity{
				XsiType: "PQ",
				Value:   formatFloat(value),
				Unit:    unit,
			},
		},
	}

	as.Component = append(as.Component, AnnotationComponent{Annotation: ann})
	return len(as.Component) - 1
}

// AddLeadAnnotation adds a lead-specific annotation with supportingROI.
//
// Creates an annotation that applies to a specific ECG lead. The annotation includes
// a supportingROI structure that identifies which lead the measurements apply to.
//
// Parameters:
//   - leadCode: The lead code (e.g., "MDC_ECG_LEAD_I", "MDC_ECG_LEAD_V1")
//   - matrixCode: The measurement matrix code (typically "MEASUREMENT_MATRIX")
//   - codeSystem: The code system OID (empty for vendor codes)
//   - codeSystemName: The code system name (e.g., "MINDRAY")
//
// Returns:
//   - int: Index of the lead annotation (use GetAnnotation to retrieve safely)
//
// Example:
//
//	idx := annotationSet.AddLeadAnnotation("MDC_ECG_LEAD_I", "MEASUREMENT_MATRIX", "", "MINDRAY")
//	ann := annotationSet.GetAnnotation(idx)
//	ann.AddNestedAnnotation("P_ONSET", "", 234, "ms")
func (as *AnnotationSet) AddLeadAnnotation(leadCode, matrixCode, codeSystem, codeSystemName string) int {
	if as == nil {
		return -1
	}

	// Input validation
	if leadCode == "" {
		return -1
	}
	if matrixCode == "" {
		return -1
	}
	// Note: codeSystem can be empty for vendor codes, but codeSystemName should be provided in that case
	if codeSystem == "" && codeSystemName == "" {
		return -1
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
	return len(as.Component) - 1
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
//   - int: Index of the nested annotation (use GetNestedAnnotation to retrieve safely)
//
// Example:
//
//	idx := annotationSet.AddAnnotation("MDC_ECG_TIME_PD_QTc", MDC_OID, 0, "")
//	ann := annotationSet.GetAnnotation(idx)
//	nestedIdx := ann.AddNestedAnnotation("ECG_TIME_PD_QTcH", "", 413, "ms")
func (a *Annotation) AddNestedAnnotation(code, codeSystem string, value float64, unit string) int {
	if a == nil {
		return -1
	}

	// Input validation
	if code == "" {
		return -1
	}
	// Note: codeSystem can be empty for vendor codes
	if unit == "" {
		return -1
	}
	if isInvalidFloat(value) {
		return -1
	}

	nested := Annotation{
		Code: &Code[string, string]{
			Code:       code,
			CodeSystem: codeSystem,
		},
		Value: &AnnotationValue{
			XsiType: "PQ",
			Typed: &PhysicalQuantity{
				XsiType: "PQ",
				Value:   formatFloat(value),
				Unit:    unit,
			},
		},
	}

	a.Component = append(a.Component, AnnotationComponent{Annotation: nested})
	return len(a.Component) - 1
}

// AddNestedAnnotationWithCodeSystemName adds a nested annotation with codeSystemName.
//
// Parameters:
//   - code: The nested annotation code
//   - codeSystemName: The code system name (e.g., "MINDRAY")
//   - value: The numeric value
//   - unit: The unit of measurement
//
// Returns:
//   - int: Index of the nested annotation (use GetNestedAnnotation to retrieve safely)
func (a *Annotation) AddNestedAnnotationWithCodeSystemName(code, codeSystemName string, value float64, unit string) int {
	if a == nil {
		return -1
	}

	// Input validation
	if code == "" {
		return -1
	}
	if codeSystemName == "" {
		return -1
	}
	if unit == "" {
		return -1
	}
	if isInvalidFloat(value) {
		return -1
	}

	nested := Annotation{
		Code: &Code[string, string]{
			Code:           code,
			CodeSystem:     "",
			CodeSystemName: codeSystemName,
		},
		Value: &AnnotationValue{
			XsiType: "PQ",
			Typed: &PhysicalQuantity{
				XsiType: "PQ",
				Value:   formatFloat(value),
				Unit:    unit,
			},
		},
	}

	a.Component = append(a.Component, AnnotationComponent{Annotation: nested})
	return len(a.Component) - 1
}

// =============================================================================
// Text Annotation Builder Methods
// =============================================================================

// AddTextAnnotation adds a global text annotation (interpretation statement, etc).
//
// Used for textual annotations like ECG interpretation statements.
//
// Parameters:
//   - code: The annotation code (e.g., "MDC_ECG_INTERPRETATION")
//   - codeSystem: The code system OID (e.g., "2.16.840.1.113883.6.24" for MDC)
//   - text: The text content
//
// Returns:
//   - int: Index of the newly created annotation (use GetAnnotation to retrieve safely)
//
// Example:
//
//	idx := annotationSet.AddTextAnnotation("MDC_ECG_INTERPRETATION", "2.16.840.1.113883.6.24", "Rythme sinusal")
//	ann := annotationSet.GetAnnotation(idx)
func (as *AnnotationSet) AddTextAnnotation(code, codeSystem, text string) int {
	if as == nil {
		return -1
	}

	// Input validation
	if code == "" {
		return -1
	}
	if codeSystem == "" {
		return -1
	}

	ann := Annotation{
		Code: &Code[string, string]{
			Code:       code,
			CodeSystem: codeSystem,
		},
	}

	// Only add Value if text is not empty (allows container annotations)
	if text != "" {
		ann.Value = &AnnotationValue{
			XsiType: "ST",
			Typed: &StringValue{
				XsiType: "ST",
				Value:   text,
			},
		}
	}

	as.Component = append(as.Component, AnnotationComponent{Annotation: ann})
	return len(as.Component) - 1
}

// AddTextAnnotationWithCodeSystemName adds a text annotation with codeSystemName.
//
// Used for vendor-specific textual annotations.
//
// Parameters:
//   - code: The annotation code
//   - codeSystemName: The code system name (e.g., "MINDRAY")
//   - text: The text content
//
// Returns:
//   - int: Index of the newly created annotation (use GetAnnotation to retrieve safely)
func (as *AnnotationSet) AddTextAnnotationWithCodeSystemName(code, codeSystemName, text string) int {
	if as == nil {
		return -1
	}

	// Input validation
	if code == "" {
		return -1
	}
	if codeSystemName == "" {
		return -1
	}

	ann := Annotation{
		Code: &Code[string, string]{
			Code:           code,
			CodeSystem:     "", // Empty for vendor-specific codes
			CodeSystemName: codeSystemName,
		},
	}

	// Only add Value if text is not empty (allows container annotations)
	if text != "" {
		ann.Value = &AnnotationValue{
			XsiType: "ST",
			Typed: &StringValue{
				XsiType: "ST",
				Value:   text,
			},
		}
	}

	as.Component = append(as.Component, AnnotationComponent{Annotation: ann})
	return len(as.Component) - 1
}

// AddNestedTextAnnotation adds a nested text annotation to an existing annotation.
//
// Parameters:
//   - code: The nested annotation code
//   - codeSystem: The code system OID (can be empty for vendor codes)
//   - text: The text content
//
// Returns:
//   - int: Index of the nested annotation (use GetNestedAnnotation to retrieve safely)
//
// Example:
//
//	interpretIdx := annotationSet.AddTextAnnotation("MDC_ECG_INTERPRETATION", MDC_OID, "")
//	interp := annotationSet.GetAnnotation(interpretIdx)
//	interp.AddNestedTextAnnotation("MDC_ECG_INTERPRETATION_STATEMENT", MDC_OID, "Rythme sinusal avec ESA")
func (a *Annotation) AddNestedTextAnnotation(code, codeSystem, text string) int {
	if a == nil {
		return -1
	}

	// Input validation
	if code == "" {
		return -1
	}
	// Note: codeSystem can be empty for vendor codes

	nested := Annotation{
		Code: &Code[string, string]{
			Code:       code,
			CodeSystem: codeSystem,
		},
	}

	// Only add Value if text is not empty (allows container annotations)
	if text != "" {
		nested.Value = &AnnotationValue{
			XsiType: "ST",
			Typed: &StringValue{
				XsiType: "ST",
				Value:   text,
			},
		}
	}

	a.Component = append(a.Component, AnnotationComponent{Annotation: nested})
	return len(a.Component) - 1
}

// AddNestedTextAnnotationWithCodeSystemName adds a nested text annotation with codeSystemName.
//
// Parameters:
//   - code: The nested annotation code
//   - codeSystemName: The code system name (e.g., "MINDRAY")
//   - text: The text content
//
// Returns:
//   - int: Index of the nested annotation (use GetNestedAnnotation to retrieve safely)
func (a *Annotation) AddNestedTextAnnotationWithCodeSystemName(code, codeSystemName, text string) int {
	if a == nil {
		return -1
	}

	// Input validation
	if code == "" {
		return -1
	}
	if codeSystemName == "" {
		return -1
	}

	nested := Annotation{
		Code: &Code[string, string]{
			Code:           code,
			CodeSystem:     "",
			CodeSystemName: codeSystemName,
		},
	}

	// Only add Value if text is not empty (allows container annotations)
	if text != "" {
		nested.Value = &AnnotationValue{
			XsiType: "ST",
			Typed: &StringValue{
				XsiType: "ST",
				Value:   text,
			},
		}
	}

	a.Component = append(a.Component, AnnotationComponent{Annotation: nested})
	return len(a.Component) - 1
}

// =============================================================================
// Safe Accessor Methods
// =============================================================================

// GetAnnotation safely retrieves an annotation by index.
// Returns nil if index is out of bounds.
func (as *AnnotationSet) GetAnnotation(index int) *Annotation {
	if as == nil || index < 0 || index >= len(as.Component) {
		return nil
	}
	return &as.Component[index].Annotation
}

// GetNestedAnnotation safely retrieves a nested annotation by index.
// Returns nil if index is out of bounds.
func (a *Annotation) GetNestedAnnotation(index int) *Annotation {
	if a == nil || index < 0 || index >= len(a.Component) {
		return nil
	}
	return &a.Component[index].Annotation
}

// =============================================================================
// Convenience Methods for Common Annotations
// =============================================================================

// AddHeartRate adds a heart rate annotation in beats per minute.
// Returns the index of the added annotation.
func (as *AnnotationSet) AddHeartRate(value float64) int {
	return as.AddAnnotation(string(MDC_ECG_HEART_RATE), string(MDC_OID), value, "bpm")
}

// AddPRInterval adds a PR interval annotation in milliseconds.
// Returns the index of the added annotation.
func (as *AnnotationSet) AddPRInterval(value float64) int {
	return as.AddAnnotation(string(MDC_ECG_TIME_PD_PR), string(MDC_OID), value, "ms")
}

// AddQRSDuration adds a QRS duration annotation in milliseconds.
// Returns the index of the added annotation.
func (as *AnnotationSet) AddQRSDuration(value float64) int {
	return as.AddAnnotation(string(MDC_ECG_TIME_PD_QRS), string(MDC_OID), value, "ms")
}

// AddQTInterval adds a QT interval annotation in milliseconds.
// Returns the index of the added annotation.
func (as *AnnotationSet) AddQTInterval(value float64) int {
	return as.AddAnnotation(string(MDC_ECG_TIME_PD_QT), string(MDC_OID), value, "ms")
}

// AddQTcInterval adds a QTc interval annotation in milliseconds.
// Returns the index of the added annotation so nested correction formulas can be added.
func (as *AnnotationSet) AddQTcInterval(value float64) int {
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
