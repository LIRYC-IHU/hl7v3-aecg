package types

// This file contains helper methods for setting IDs on various structures
// with automatic default extensions based on the structure type.

// SetID sets the ID of the ClinicalTrial with automatic default extension.
// If extension is empty, uses "clinicalTrial" as default.
func (ct *ClinicalTrial) SetID(root, extension string) *ClinicalTrial {
	ct.ID.SetID(root, extension, "clinicalTrial")
	return ct
}

// SetID sets the ID of the TrialSubject with automatic default extension.
// If extension is empty, uses "trialSubject" as default.
func (ts *TrialSubject) SetID(root, extension string) *TrialSubject {
	ts.ID.SetID(root, extension, "trialSubject")
	return ts
}

// SetID sets the ID of the Subject with automatic default extension.
// If extension is empty, uses "subject" as default.
func (s *Subject) SetID(root, extension string) *Subject {
	// Subject doesn't have its own ID, it uses TrialSubject.ID
	s.TrialSubject.SetID(root, extension)
	return s
}

// SetID sets the ID of the TrialSite with automatic default extension.
// If extension is empty, uses "trialSite" as default.
func (ts *TrialSite) SetID(root, extension string) *TrialSite {
	ts.ID.SetID(root, extension, "trialSite")
	return ts
}

// SetID sets the ID of the HL7AEcg document with automatic default extension.
// If extension is empty, uses "annotatedEcg" as default.
func (h *HL7AEcg) SetID(root, extension string) *HL7AEcg {
	h.ID.SetID(root, extension, "annotatedEcg")
	return h
}

// SetID sets the ID of the Series with automatic default extension.
// If extension is empty, uses "series" as default.
func (s *Series) SetID(root, extension string) *Series {
	s.ID.SetID(root, extension, "series")
	return s
}

// SetID sets the ID of the TrialInvestigator with automatic default extension.
// If extension is empty, uses "trialInvestigator" as default.
func (ti *TrialInvestigator) SetID(root, extension string) *TrialInvestigator {
	ti.ID.SetID(root, extension, "trialInvestigator")
	return ti
}

// SetID sets the ID of the ResponsibleParty's TrialInvestigator with automatic default extension.
// This is a convenience method that delegates to TrialInvestigator.SetID.
func (rp *ResponsibleParty) SetID(root, extension string) *ResponsibleParty {
	rp.TrialInvestigator.SetID(root, extension)
	return rp
}
