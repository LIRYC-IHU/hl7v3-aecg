package types

// SetInvestigatorID sets the ID of the trial investigator.
//
// Parameters:
//   - root: The root OID for the investigator
//   - extension: The investigator identifier extension
//
// Returns the ResponsibleParty for method chaining.
func (rp *ResponsibleParty) SetInvestigatorID(root, extension string) *ResponsibleParty {
	rp.TrialInvestigator.ID.SetID(root, extension)
	return rp
}

// SetInvestigatorName sets the name of the trial investigator.
//
// Parameters:
//   - prefix: Title/honorific (e.g., "Dr.", "Prof.")
//   - given: Given/first name
//   - family: Family/last name
//   - suffix: Suffix/credentials (e.g., "MD", "PhD")
//
// Use empty string for any component you don't want to set.
//
// Returns the ResponsibleParty for method chaining.
func (rp *ResponsibleParty) SetInvestigatorName(prefix, given, family, suffix string) *ResponsibleParty {
	if rp.TrialInvestigator.InvestigatorPerson == nil {
		rp.TrialInvestigator.InvestigatorPerson = &InvestigatorPerson{}
	}
	if rp.TrialInvestigator.InvestigatorPerson.Name == nil {
		rp.TrialInvestigator.InvestigatorPerson.Name = &PersonName{}
	}

	name := rp.TrialInvestigator.InvestigatorPerson.Name
	if prefix != "" {
		name.Prefix = &prefix
	}
	if given != "" {
		name.Given = &given
	}
	if family != "" {
		name.Family = &family
	}
	if suffix != "" {
		name.Suffix = &suffix
	}

	return rp
}

// SetEmptyInvestigatorName initializes an empty name element for the investigator.
// This is useful when you want the <name/> element without any content.
//
// Returns the ResponsibleParty for method chaining.
func (rp *ResponsibleParty) SetEmptyInvestigatorName() *ResponsibleParty {
	if rp.TrialInvestigator.InvestigatorPerson == nil {
		rp.TrialInvestigator.InvestigatorPerson = &InvestigatorPerson{}
	}
	rp.TrialInvestigator.InvestigatorPerson.Name = &PersonName{}
	return rp
}
