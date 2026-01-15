package types

import (
	"context"
	"fmt"
)

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
		if err := s.Author.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	// SecondaryPerformer is optional but must be valid if present
	for i := range s.SecondaryPerformer {
		if err := s.SecondaryPerformer[i].Validate(ctx, vctx); err != nil {
			return err
		}
	}

	// ControlVariable is optional but must be valid if present
	for i := range s.ControlVariable {
		if err := s.ControlVariable[i].Validate(ctx, vctx); err != nil {
			return err
		}
	}

	// Derivation is optional but must be valid if present
	for i, deriv := range s.Derivation {
		if err := deriv.Validate(ctx, vctx); err != nil {
			vctx.AddError(NewValidationError(
				fmt.Sprintf("Series.Derivation[%d]", i),
				fmt.Sprintf("Derivation validation failed: %v", err),
			))
		}
	}

	// TODO: Component (SequenceSet) validation if needed

	return nil
}

// Validate validates the Author structure.
// Validates SeriesAuthor (required).
func (a *Author) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// SeriesAuthor is required
	if err := a.SeriesAuthor.Validate(ctx, vctx); err != nil {
		return err
	}

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

// Validate validates the SecondaryPerformer structure.
func (sp *SecondaryPerformer) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// FunctionCode is optional but must be valid if present
	if sp.FunctionCode != nil {
		sp.FunctionCode.ValidateCode(ctx, vctx, "PerformerFunctionCode")
	}

	// Time is optional but must be valid if present
	if sp.Time != nil {
		if err := sp.Time.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	// SeriesPerformer validation
	if err := sp.SeriesPerformer.Validate(ctx, vctx); err != nil {
		return err
	}

	return nil
}

// Validate validates the SeriesPerformer structure.
func (sp *SeriesPerformer) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ID is optional but must be valid if present
	if sp.ID != nil {
		if err := sp.ID.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	// AssignedPerson is optional - no specific validation needed for name

	return nil
}

// Validate validates the ControlVariable structure.
func (cv *ControlVariable) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ControlVariable inner structure is optional but must be valid if present
	if cv.ControlVariable != nil {
		if err := cv.ControlVariable.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates the ControlVariableInner structure (recursive).
func (cvi *ControlVariableInner) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Code is optional but must be valid if present
	if cvi.Code != nil {
		cvi.Code.ValidateCode(ctx, vctx, "ControlVariableCode")
	}

	// Value is optional - no validation needed for PhysicalQuantity
	// Text is optional - no validation needed for strings

	// Component is optional but must be valid if present
	for i := range cvi.Component {
		if err := cvi.Component[i].Validate(ctx, vctx); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates the ControlVariableComponent structure.
func (cvc *ControlVariableComponent) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ControlVariable is required within a component
	if cvc.ControlVariable != nil {
		if err := cvc.ControlVariable.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates the Derivation structure.
//
// Checks:
//   - DerivedSeries is valid
//   - DerivedSeries code is appropriate (REPRESENTATIVE_BEAT or MEDIAN_BEAT)
//   - DerivedSeries does not have nested derivation (no recursive derivation)
//   - DerivedSeries uses TIME_RELATIVE for time sequences
func (d *Derivation) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate derived series structure
	if err := d.DerivedSeries.Validate(ctx, vctx); err != nil {
		return err
	}

	// Validate code is appropriate for derived series
	if d.DerivedSeries.Code != nil {
		code := d.DerivedSeries.Code.Code
		if code != REPRESENTATIVE_BEAT_CODE && code != MEDIAN_BEAT_CODE {
			vctx.AddError(NewValidationError(
				"Derivation.DerivedSeries.Code",
				fmt.Sprintf("Derived series code must be REPRESENTATIVE_BEAT or MEDIAN_BEAT, got: %s", code),
			))
		}
	}

	// Validate no nested derivation (derived series cannot have derivation)
	if len(d.DerivedSeries.Derivation) > 0 {
		vctx.AddError(NewValidationError(
			"Derivation.DerivedSeries.Derivation",
			"Derived series cannot have nested derivation (recursive derivation not allowed)",
		))
	}

	// Validate time sequences use TIME_RELATIVE
	if len(d.DerivedSeries.Component) > 0 {
		for i, comp := range d.DerivedSeries.Component {
			seqSet := comp.SequenceSet
			if len(seqSet.Component) > 0 {
				// First sequence should be time
				firstSeq := seqSet.Component[0].Sequence
				if firstSeq.Code.Time != nil {
					if firstSeq.Code.Time.Code != TIME_RELATIVE_CODE {
						vctx.AddError(NewValidationError(
							fmt.Sprintf("Derivation.DerivedSeries.Component[%d].SequenceSet", i),
							fmt.Sprintf("Derived series must use TIME_RELATIVE for time sequences, got: %s", firstSeq.Code.Time.Code),
						))
					}
				}
			}
		}
	}

	return nil
}
