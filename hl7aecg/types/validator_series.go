package types

import "context"

// Validate validates the Series structure.
// Validates Code (required), EffectiveTime (required), ID (optional), and Author (optional).
func (s *Series) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ID is optional but must be valid if present
	if s.ID != nil {
		s.ID.Validate(ctx, vctx)
	}

	// Code is required
	if s.Code == nil {
		vctx.AddError(ErrMissingCode)
	} else {
		s.Code.ValidateCode(ctx, vctx, "SeriesCode")
	}

	// EffectiveTime is required
	s.EffectiveTime.Validate(ctx, vctx)

	// Author is optional but must be valid if present
	if s.Author != nil {
		s.Author.Validate(ctx, vctx)
	}

	// SecondaryPerformer is optional - no validation needed for now
	// Support is optional - no validation needed for now
	// Component is optional - validated separately if needed

	return nil
}

// Validate validates the SeriesAuthor structure.
// Validates ID (optional) and ManufacturedSeriesDevice (required).
func (sa *SeriesAuthor) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ID is optional but must be valid if present
	if sa.ID != nil {
		sa.ID.Validate(ctx, vctx)
	}

	// ManufacturedSeriesDevice is required within SeriesAuthor
	sa.ManufacturedSeriesDevice.Validate(ctx, vctx)

	// ManufacturerOrganization is optional but must be valid if present
	if sa.ManufacturerOrganization != nil {
		sa.ManufacturerOrganization.Validate(ctx, vctx)
	}

	return nil
}

// Validate validates the ManufacturedSeriesDevice structure.
// Validates ID (optional) and Code (optional).
func (msd *ManufacturedSeriesDevice) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ID is optional but must be valid if present
	if msd.ID != nil {
		msd.ID.Validate(ctx, vctx)
	}

	// Code is optional but must be valid if present
	if msd.Code != nil {
		msd.Code.ValidateCode(ctx, vctx, "DeviceTypeCode")
	}

	// ManufacturerModelName is optional - no validation needed for strings
	// SoftwareName is optional - no validation needed for strings

	return nil
}

// Validate validates the ManufacturerOrganization structure.
// Validates ID (optional) and Name (optional).
func (mo *ManufacturerOrganization) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ID is optional but must be valid if present
	if mo.ID != nil {
		mo.ID.Validate(ctx, vctx)
	}

	// Name is optional - no validation needed for strings

	return nil
}
