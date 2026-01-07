package types

// =============================================================================
// ControlVariable Setter Methods
// =============================================================================

// NewControlVariable creates a new ControlVariable with the specified code.
//
// This is the primary constructor for creating control variables. The inner
// ControlVariable structure is automatically initialized.
//
// Parameters:
//   - code: The control variable code (e.g., "MDC_ECG_CTL_VBL_ATTR_FILTER_LOW_PASS")
//   - codeSystem: The OID of the code system (e.g., MDC_OID for MDC codes)
//   - codeSystemName: Human-readable name of the code system (e.g., "MDC")
//   - displayName: Human-readable name of the code (e.g., "Low Pass Filter")
//
// Example:
//
//	cv := types.NewControlVariable(
//	    "MDC_ECG_CTL_VBL_ATTR_FILTER_LOW_PASS",
//	    types.MDC_OID,
//	    "MDC",
//	    "Low Pass Filter",
//	)
//
// Returns a pointer to the created ControlVariable for method chaining.
func NewControlVariable(code string, codeSystem CodeSystemOID, codeSystemName, displayName string) *ControlVariable {
	return &ControlVariable{
		ControlVariable: &ControlVariableInner{
			Code: &Code[string, CodeSystemOID]{
				Code:           code,
				CodeSystem:     codeSystem,
				CodeSystemName: codeSystemName,
				DisplayName:    displayName,
			},
		},
	}
}

// SetCode sets or updates the code of the control variable.
//
// Parameters:
//   - code: The control variable code
//   - codeSystem: The OID of the code system
//   - codeSystemName: Human-readable name of the code system
//   - displayName: Human-readable name of the code
//
// Returns the ControlVariableInner pointer for method chaining.
func (cvi *ControlVariableInner) SetCode(code string, codeSystem CodeSystemOID, codeSystemName, displayName string) *ControlVariableInner {
	cvi.Code = &Code[string, CodeSystemOID]{
		Code:           code,
		CodeSystem:     codeSystem,
		CodeSystemName: codeSystemName,
		DisplayName:    displayName,
	}
	return cvi
}

// SetValue sets the physical quantity value of the control variable.
//
// Parameters:
//   - value: The numeric value as a string (e.g., "35", "0.56")
//   - unit: The unit of measure (e.g., "Hz", "a" for years, "Cel" for Celsius)
//
// Example:
//
//	cvi.SetValue("35", "Hz")  // 35 Hertz
//	cvi.SetValue("34", "a")   // 34 years (age)
//
// Returns the ControlVariableInner pointer for method chaining.
func (cvi *ControlVariableInner) SetValue(value, unit string) *ControlVariableInner {
	cvi.Value = &PhysicalQuantity{
		XsiType: "PQ",
		Value:   value,
		Unit:    unit,
	}
	return cvi
}

// SetText sets the textual description of the control variable.
//
// Parameters:
//   - text: The descriptive text
//
// Returns the ControlVariableInner pointer for method chaining.
func (cvi *ControlVariableInner) SetText(text string) *ControlVariableInner {
	cvi.Text = &text
	return cvi
}

// AddComponent adds a nested control variable as a component.
//
// This method allows building hierarchical control variable structures,
// such as a filter specification with nested parameters.
//
// Parameters:
//   - code: The component control variable code
//   - codeSystem: The OID of the code system
//   - codeSystemName: Human-readable name of the code system
//   - displayName: Human-readable name of the code
//   - value: Optional value (use "" to skip)
//   - unit: Optional unit (use "" to skip)
//
// Example:
//
//	// Add cutoff frequency component to a low pass filter
//	lowPassFilter.ControlVariable.AddComponent(
//	    "MDC_ECG_CTL_VBL_ATTR_FILTER_CUTOFF_FREQ",
//	    types.MDC_OID,
//	    "MDC",
//	    "Cutoff Frequency",
//	    "35",
//	    "Hz",
//	)
//
// Returns the ControlVariableInner pointer for method chaining.
func (cvi *ControlVariableInner) AddComponent(code string, codeSystem CodeSystemOID, codeSystemName, displayName, value, unit string) *ControlVariableInner {
	component := &ControlVariableInner{
		Code: &Code[string, CodeSystemOID]{
			Code:           code,
			CodeSystem:     codeSystem,
			CodeSystemName: codeSystemName,
			DisplayName:    displayName,
		},
	}

	if value != "" {
		component.Value = &PhysicalQuantity{
			XsiType: "PQ",
			Value:   value,
			Unit:    unit,
		}
	}

	cvi.Component = append(cvi.Component, ControlVariableComponent{
		ControlVariable: component,
	})

	return cvi
}

// AddComponentWithText adds a nested control variable component with text instead of value.
//
// This is useful when the component contains textual information rather than
// a numeric value.
//
// Parameters:
//   - code: The component control variable code
//   - codeSystem: The OID of the code system
//   - codeSystemName: Human-readable name of the code system
//   - displayName: Human-readable name of the code
//   - text: The textual content
//
// Returns the ControlVariableInner pointer for method chaining.
func (cvi *ControlVariableInner) AddComponentWithText(code string, codeSystem CodeSystemOID, codeSystemName, displayName, text string) *ControlVariableInner {
	component := &ControlVariableInner{
		Code: &Code[string, CodeSystemOID]{
			Code:           code,
			CodeSystem:     codeSystem,
			CodeSystemName: codeSystemName,
			DisplayName:    displayName,
		},
		Text: &text,
	}

	cvi.Component = append(cvi.Component, ControlVariableComponent{
		ControlVariable: component,
	})

	return cvi
}

// =============================================================================
// Helper Functions for Common Control Variables
// =============================================================================

// NewLowPassFilter creates a low pass filter control variable with cutoff frequency.
//
// Parameters:
//   - cutoffFreq: The cutoff frequency value as a string (e.g., "35")
//   - unit: The frequency unit (typically "Hz")
//
// Example:
//
//	lowPass := types.NewLowPassFilter("35", "Hz")
//
// This creates:
//
//	<controlVariable>
//	  <controlVariable>
//	    <code code="MDC_ECG_CTL_VBL_ATTR_FILTER_LOW_PASS" .../>
//	    <component>
//	      <controlVariable>
//	        <code code="MDC_ECG_CTL_VBL_ATTR_FILTER_CUTOFF_FREQ" .../>
//	        <value xsi:type="PQ" value="35" unit="Hz"/>
//	      </controlVariable>
//	    </component>
//	  </controlVariable>
//	</controlVariable>
func NewLowPassFilter(cutoffFreq, unit string) *ControlVariable {
	cv := NewControlVariable(
		"MDC_ECG_CTL_VBL_ATTR_FILTER_LOW_PASS",
		MDC_OID,
		"MDC",
		"Low Pass Filter",
	)

	cv.ControlVariable.AddComponent(
		"MDC_ECG_CTL_VBL_ATTR_FILTER_CUTOFF_FREQ",
		MDC_OID,
		"MDC",
		"Cutoff Frequency",
		cutoffFreq,
		unit,
	)

	return cv
}

// NewHighPassFilter creates a high pass filter control variable with cutoff frequency.
//
// Parameters:
//   - cutoffFreq: The cutoff frequency value as a string (e.g., "0.56")
//   - unit: The frequency unit (typically "Hz")
//
// Example:
//
//	highPass := types.NewHighPassFilter("0.56", "Hz")
func NewHighPassFilter(cutoffFreq, unit string) *ControlVariable {
	cv := NewControlVariable(
		"MDC_ECG_CTL_VBL_ATTR_FILTER_HIGH_PASS",
		MDC_OID,
		"MDC",
		"High Pass Filter",
	)

	cv.ControlVariable.AddComponent(
		"MDC_ECG_CTL_VBL_ATTR_FILTER_CUTOFF_FREQ",
		MDC_OID,
		"MDC",
		"Cutoff Frequency",
		cutoffFreq,
		unit,
	)

	return cv
}

// NewNotchFilter creates a notch filter control variable with notch frequency.
//
// Parameters:
//   - notchFreq: The notch frequency value as a string (e.g., "50" or "60")
//   - unit: The frequency unit (typically "Hz")
//
// Example:
//
//	notch := types.NewNotchFilter("50", "Hz")  // 50 Hz for Europe
//	notch := types.NewNotchFilter("60", "Hz")  // 60 Hz for North America
func NewNotchFilter(notchFreq, unit string) *ControlVariable {
	cv := NewControlVariable(
		"MDC_ECG_CTL_VBL_ATTR_FILTER_NOTCH",
		MDC_OID,
		"MDC",
		"Notch Filter",
	)

	cv.ControlVariable.AddComponent(
		"MDC_ECG_CTL_VBL_ATTR_FILTER_NOTCH_FREQ",
		MDC_OID,
		"MDC",
		"Notch filter frequency",
		notchFreq,
		unit,
	)

	return cv
}

// NewAgeObservation creates an age observation control variable.
//
// Parameters:
//   - age: The age value as a string (e.g., "34")
//   - unit: The unit (typically "a" for years)
//
// Example:
//
//	age := types.NewAgeObservation("34", "a")
func NewAgeObservation(age, unit string) *ControlVariable {
	cv := NewControlVariable(
		"21612-7",
		CodeSystemOID("2.16.840.1.113883.6.1"), // LOINC
		"LOINC",
		"Reported Age",
	)

	cv.ControlVariable.SetValue(age, unit)

	return cv
}
