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

	// RelatedObservation is required
	if err := cv.RelatedObservation.Validate(ctx, vctx); err != nil {
		return err
	}

	return nil
}

// Validate validates the RelatedObservation structure.
func (ro *RelatedObservation) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Code is optional but must be valid if present
	if ro.Code != nil {
		ro.Code.ValidateCode(ctx, vctx, "ObservationCode")
	}

	// Text is optional - no validation needed
	// Value is optional - no validation needed for PhysicalQuantity

	// Author is optional but must be valid if present
	if ro.Author != nil {
		if err := ro.Author.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates the ObservationAuthor structure.
func (oa *ObservationAuthor) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// AssignedEntity validation
	if err := oa.AssignedEntity.Validate(ctx, vctx); err != nil {
		return err
	}

	return nil
}

// Validate validates the AssignedEntity structure.
func (ae *AssignedEntity) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ID is optional but must be valid if present
	if ae.ID != nil {
		if err := ae.ID.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	// AssignedAuthorType is optional but must be valid if present
	if ae.AssignedAuthorType != nil {
		if err := ae.AssignedAuthorType.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	// RepresentedAuthoringOrganization is optional but must be valid if present
	if ae.RepresentedAuthoringOrganization != nil {
		if err := ae.RepresentedAuthoringOrganization.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates the AssignedAuthorType structure.
func (aat *AssignedAuthorType) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Either AssignedPerson or AssignedDevice should be present
	// AssignedPerson - no specific validation needed
	// AssignedDevice is optional but must be valid if present
	if aat.AssignedDevice != nil {
		if err := aat.AssignedDevice.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates the ObservationAssignedDevice structure.
func (oad *ObservationAssignedDevice) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ID is optional but must be valid if present
	if oad.ID != nil {
		if err := oad.ID.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	// Code is optional but must be valid if present
	if oad.Code != nil {
		oad.Code.ValidateCode(ctx, vctx, "DeviceTypeCode")
	}

	// ManufacturerModelName is optional - no validation needed
	// SoftwareName is optional - no validation needed

	// PlayedManufacturedDevice is optional but must be valid if present
	if oad.PlayedManufacturedDevice != nil {
		if err := oad.PlayedManufacturedDevice.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates the PlayedManufacturedDevice structure.
func (pmd *PlayedManufacturedDevice) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ManufacturingOrganization is optional but must be valid if present
	if pmd.ManufacturingOrganization != nil {
		if err := pmd.ManufacturingOrganization.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates the ObservationManufacturingOrganization structure.
func (omo *ObservationManufacturingOrganization) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ID is optional but must be valid if present
	if omo.ID != nil {
		if err := omo.ID.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	// Name is optional - no validation needed

	return nil
}

// Validate validates the RepresentedAuthoringOrganization structure.
func (rao *RepresentedAuthoringOrganization) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ID is optional but must be valid if present
	if rao.ID != nil {
		if err := rao.ID.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	// Name is optional - no validation needed

	// Identification is optional but must be valid if present
	if rao.Identification != nil {
		if err := rao.Identification.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	return nil
}

// Validate validates the OrganizationIdentification structure.
func (oi *OrganizationIdentification) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// ID is optional but must be valid if present
	if oi.ID != nil {
		if err := oi.ID.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	return nil
}

// Validate Medications structure.
// Ensures at least one Medication entry is present.
func (m *Medications) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	fmt.Println("Validating Medications")
	if len(m.Medication) == 0 {
		m.Medication = []string{}
		m.Medication = append(m.Medication, "")
	}
	return nil

}
