package types

import "context"

// Validate validates the Subject structure.
// Subject is required in SubjectAssignment context.
func (s *Subject) Validate(ctx context.Context, vctx *ValidationContext) error {
	// Subject always contains TrialSubject (required)
	s.TrialSubject.Validate(ctx, vctx)
	return nil
}

// Validate validates the TrialSubject structure.
// Validates ID (required), Code (optional), and SubjectDemographicPerson (optional).
func (t *TrialSubject) Validate(ctx context.Context, vctx *ValidationContext) error {
	// ID is required
	t.ID.Validate(ctx, vctx)

	// Code is optional but if present must be valid
	if t.Code != nil {
		t.Code.ValidateCode(ctx, vctx, "TrialSubjectCode")
	}

	// SubjectDemographicPerson is optional
	if t.SubjectDemographicPerson != nil {
		t.SubjectDemographicPerson.Validate(ctx, vctx)
	}

	return nil
}

// Validate validates the SubjectDemographicPerson structure.
// All fields are optional, but if present must be valid.
func (s *SubjectDemographicPerson) Validate(ctx context.Context, vctx *ValidationContext) error {
	// Name is optional (no specific validation beyond presence)

	// AdministrativeGenderCode is optional but must be valid if present
	if s.AdministrativeGenderCode != nil {
		s.AdministrativeGenderCode.ValidateCode(ctx, vctx, "AdministrativeGender")
	}

	// BirthTime is optional (format validation could be added)
	// Could add: validate YYYYMMDD format, reasonable date range, etc.

	// RaceCode is optional but must be valid if present
	if s.RaceCode != nil {
		s.RaceCode.ValidateCode(ctx, vctx, "RaceCode")
	}

	return nil
}

// Validate validates SubjectAssignment structure.
func (s *SubjectAssignment) Validate(ctx context.Context, vctx *ValidationContext) error {
	// Subject is required
	s.Subject.Validate(ctx, vctx)

	// Definition is optional
	if s.Definition != nil {
		s.Definition.Validate(ctx, vctx)
	}

	// ComponentOf is required (links to ClinicalTrial)
	s.ComponentOf.Validate(ctx, vctx)

	return nil
}

// Validate validates SubjectAssignmentDefinition structure.
func (s *SubjectAssignmentDefinition) Validate(ctx context.Context, vctx *ValidationContext) error {
	// TreatmentGroupAssignment is required within Definition
	s.TreatmentGroupAssignment.Validate(ctx, vctx)
	return nil
}

// Validate validates TreatmentGroupAssignment structure.
func (t *TreatmentGroupAssignment) Validate(ctx context.Context, vctx *ValidationContext) error {
	// Code is required
	t.Code.ValidateCode(ctx, vctx, "TreatmentGroupCode")
	return nil
}

// Validate validates ComponentOfClinicalTrial structure.
func (c *ComponentOfClinicalTrial) Validate(ctx context.Context, vctx *ValidationContext) error {
	// ClinicalTrial is required
	c.ClinicalTrial.Validate(ctx, vctx)
	return nil
}
