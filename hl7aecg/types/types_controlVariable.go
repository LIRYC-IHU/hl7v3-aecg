package types

// =============================================================================
// Control Variable Types
// =============================================================================

// ControlVariable captures related information about the subject or ECG collection conditions.
//
// This structure supports nested control variables, where a parent control variable
// can have child control variables in components. This is commonly used for filter
// specifications where a filter type (e.g., Low Pass Filter) has nested parameters
// (e.g., Cutoff Frequency).
//
// XML Structure (nested example):
//
//	<controlVariable>
//	  <controlVariable>
//	    <code code="MDC_ECG_CTL_VBL_ATTR_FILTER_LOW_PASS"
//	          codeSystem="2.16.840.1.113883.6.24"
//	          codeSystemName="MDC"
//	          displayName="Low Pass Filter"/>
//	    <component>
//	      <controlVariable>
//	        <code code="MDC_ECG_CTL_VBL_ATTR_FILTER_CUTOFF_FREQ"
//	              codeSystem="2.16.840.1.113883.6.24"
//	              codeSystemName="MDC"
//	              displayName="Cutoff Frequency"/>
//	        <value xsi:type="PQ" value="35" unit="Hz"/>
//	      </controlVariable>
//	    </component>
//	  </controlVariable>
//	</controlVariable>
//
// XML Structure (simple observation):
//
//	<controlVariable>
//	  <controlVariable>
//	    <code code="21612-7"
//	          codeSystem="2.16.840.1.113883.6.1"
//	          displayName="Reported Age"/>
//	    <value xsi:type="PQ" value="34" unit="a"/>
//	  </controlVariable>
//	</controlVariable>
//
// Cardinality: Optional (0..* in Series)
//
// Reference: HL7 aECG Implementation Guide
type ControlVariable struct {
	// ControlVariable contains the actual control variable data.
	//
	// This nested structure is required by the HL7 aECG standard.
	// The outer <controlVariable> is a wrapper, while the inner one
	// contains the code, value, and optional components.
	//
	// XML Tag: <controlVariable>...</controlVariable>
	// Cardinality: Optional (but typically present when ControlVariable is used)
	ControlVariable *ControlVariableInner `xml:"controlVariable,omitempty"`
}

// ControlVariableInner represents the actual control variable data with code, value, and nested components.
//
// This structure can be used recursively - a control variable can contain
// components, each of which contains another control variable.
//
// XML Structure:
//
//	<controlVariable>
//	  <code code="..." codeSystem="..." codeSystemName="..." displayName="..."/>
//	  <text>Optional description</text>
//	  <value xsi:type="PQ" value="..." unit="..."/>
//	  <component>
//	    <controlVariable>...</controlVariable>
//	  </component>
//	</controlVariable>
//
// Cardinality: Required (within ControlVariable wrapper)
type ControlVariableInner struct {
	// Code identifies the type of control variable.
	//
	// Common codes from MDC (2.16.840.1.113883.6.24):
	//   - "MDC_ECG_CTL_VBL_ATTR_FILTER_LOW_PASS": Low pass filter
	//   - "MDC_ECG_CTL_VBL_ATTR_FILTER_HIGH_PASS": High pass filter
	//   - "MDC_ECG_CTL_VBL_ATTR_FILTER_NOTCH": Notch filter
	//   - "MDC_ECG_CTL_VBL_ATTR_FILTER_CUTOFF_FREQ": Cutoff frequency
	//   - "MDC_ECG_CTL_VBL_ATTR_FILTER_NOTCH_FREQ": Notch filter frequency
	//
	// Common codes from LOINC (2.16.840.1.113883.6.1):
	//   - "21612-7": Reported Age
	//   - "49541-6": Fasting status
	//
	// XML Tag: <code code="..." codeSystem="..." codeSystemName="..." displayName="..."/>
	// Cardinality: Optional
	Code *Code[string, CodeSystemOID] `xml:"code,omitempty"`

	// Text provides a textual description of the control variable.
	//
	// This can be used when additional explanation is needed beyond
	// what the code and value provide.
	//
	// XML Tag: <text>...</text>
	// Cardinality: Optional
	Text *string `xml:"text,omitempty"`

	// Value contains the value of the control variable.
	//
	// For example:
	//   - Age: value="34" unit="a" (years)
	//   - Cutoff frequency: value="35" unit="Hz"
	//   - Temperature: value="20" unit="Cel"
	//
	// Note: The xsi:type attribute is typically "PQ" (Physical Quantity)
	//
	// XML Tag: <value xsi:type="PQ" value="..." unit="..."/>
	// Cardinality: Optional
	Value *PhysicalQuantity `xml:"value,omitempty"`

	// Component contains nested control variables.
	//
	// This is used when a control variable has sub-parameters. For example,
	// a filter specification (Low Pass Filter) can have nested components
	// for its parameters (Cutoff Frequency, Filter Order, etc.).
	//
	// XML Tag: <component>...</component>
	// Cardinality: Optional (0..*)
	Component []ControlVariableComponent `xml:"component,omitempty"`
}

// ControlVariableComponent represents a component that contains a nested control variable.
//
// This structure allows for recursive nesting of control variables, enabling
// complex hierarchical specifications like filter configurations with multiple
// parameters.
//
// XML Structure:
//
//	<component>
//	  <controlVariable>
//	    <code code="..." codeSystem="..." displayName="..."/>
//	    <value xsi:type="PQ" value="..." unit="..."/>
//	  </controlVariable>
//	</component>
//
// Cardinality: Optional (0..* in ControlVariableInner)
type ControlVariableComponent struct {
	// ControlVariable contains the nested control variable data.
	//
	// XML Tag: <controlVariable>...</controlVariable>
	// Cardinality: Required (within ControlVariableComponent)
	ControlVariable *ControlVariableInner `xml:"controlVariable,omitempty"`
}
