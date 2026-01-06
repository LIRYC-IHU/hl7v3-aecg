package types

import "context"

// Validate validates the ClinicalTrial structure.
// Validates ID (required), ActivityTime (optional), and Location (optional).
func (ct *ClinicalTrial) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ID is required
	ct.ID.Validate(ctx, vctx)

	// ActivityTime is optional but must be valid if present
	// Note: ActivityTime is not a pointer, so we check if it has values
	if ct.ActivityTime.Low.Value != "" || ct.ActivityTime.High.Value != "" {
		ct.ActivityTime.Validate(ctx, vctx)
	}

	// Location is optional but must be valid if present
	if ct.Location != nil {
		ct.Location.Validate(ctx, vctx)
	}

	return nil
}

// Validate validates the Location structure.
// Location always contains TrialSite (required).
func (l *Location) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// TrialSite is required within Location
	l.TrialSite.Validate(ctx, vctx)

	return nil
}

// Validate validates the TrialSite structure.
// Validates ID (required), Location (optional), and ResponsibleParty (optional).
func (ts *TrialSite) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ID is required
	ts.ID.Validate(ctx, vctx)

	// Location (SiteLocation) is optional - no validation needed for simple strings
	// ResponsibleParty is optional but must be valid if present
	if ts.ResponsibleParty != nil {
		ts.ResponsibleParty.Validate(ctx, vctx)
	}

	return nil
}

// Validate validates the ResponsibleParty structure.
// ResponsibleParty always contains TrialInvestigator (required).
func (rp *ResponsibleParty) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// TrialInvestigator is required within ResponsibleParty
	rp.TrialInvestigator.Validate(ctx, vctx)

	return nil
}

// Validate validates the TrialInvestigator structure.
// Validates ID (required), InvestigatorPerson (optional).
func (ti *TrialInvestigator) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ID is required
	ti.ID.Validate(ctx, vctx)

	// InvestigatorPerson is optional - no additional validation needed for names

	return nil
}
