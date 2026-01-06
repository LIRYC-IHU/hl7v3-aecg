package types

// SetFunctionCode sets the function code for the secondary performer.
//
// Parameters:
//   - code: The function code (e.g., PERFORMER_ECG_TECHNICIAN)
//   - codeSystem: The code system OID (typically empty string for suggested codes)
//   - displayName: Optional display name for the code
//   - codeSystemName: Optional code system name
//
// Returns the SecondaryPerformer for method chaining.
func (sp *SecondaryPerformer) SetFunctionCode(code PerformerFunctionCode, codeSystem CodeSystemOID, displayName, codeSystemName string) *SecondaryPerformer {
	if sp.FunctionCode == nil {
		sp.FunctionCode = &Code[PerformerFunctionCode, CodeSystemOID]{}
	}
	sp.FunctionCode.SetCode(code, codeSystem, displayName, codeSystemName)
	return sp
}

// SetTime sets the time period for the secondary performer's function.
//
// Parameters:
//   - low: Start time in HL7 TS format (e.g., "20021122091000")
//   - high: End time in HL7 TS format
//
// Returns the SecondaryPerformer for method chaining.
func (sp *SecondaryPerformer) SetTime(low, high string) *SecondaryPerformer {
	if sp.Time == nil {
		sp.Time = &EffectiveTime{}
	}
	sp.Time.Low.Value = low
	sp.Time.High.Value = high
	return sp
}

// SetPerformerID sets the ID for the series performer.
//
// Parameters:
//   - root: The root OID or UUID
//   - extension: Optional extension (traditional identifier)
//
// Returns the SecondaryPerformer for method chaining.
func (sp *SecondaryPerformer) SetPerformerID(root, extension string) *SecondaryPerformer {
	if sp.SeriesPerformer.ID == nil {
		sp.SeriesPerformer.ID = &ID{}
	}
	sp.SeriesPerformer.ID.SetID(root, extension)
	return sp
}

// SetPerformerName sets the name for the assigned person.
//
// Parameters:
//   - name: The name of the technician (can be initials, full name, etc.)
//
// Returns the SecondaryPerformer for method chaining.
func (sp *SecondaryPerformer) SetPerformerName(name string) *SecondaryPerformer {
	if sp.SeriesPerformer.AssignedPerson == nil {
		sp.SeriesPerformer.AssignedPerson = &SeriesAssignedPerson{}
	}
	sp.SeriesPerformer.AssignedPerson.Name = &name
	return sp
}

// SetEmptyPerformerName sets an empty name element for the assigned person.
// This creates the <name/> XML element without content.
//
// Returns the SecondaryPerformer for method chaining.
func (sp *SecondaryPerformer) SetEmptyPerformerName() *SecondaryPerformer {
	if sp.SeriesPerformer.AssignedPerson == nil {
		sp.SeriesPerformer.AssignedPerson = &SeriesAssignedPerson{}
	}
	emptyName := ""
	sp.SeriesPerformer.AssignedPerson.Name = &emptyName
	return sp
}

// SetSeriesPerformer sets the series performer details.
//
// Parameters:
//   - performerID: The role-specific identifier root
//   - performerExtension: The role-specific identifier extension
//   - name: The technician's name (empty string for <name/> element)
//
// Returns the SeriesPerformer for method chaining.
func (spe *SeriesPerformer) SetPerformer(performerID, performerExtension, name string) *SeriesPerformer {
	// Set ID if provided
	if performerID != "" || performerExtension != "" {
		if spe.ID == nil {
			spe.ID = &ID{}
		}
		spe.ID.SetID(performerID, performerExtension)
	}

	// Set name
	if spe.AssignedPerson == nil {
		spe.AssignedPerson = &SeriesAssignedPerson{}
	}
	spe.AssignedPerson.Name = &name

	return spe
}
