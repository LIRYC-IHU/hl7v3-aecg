package types

import (
	"context"
	"strconv"
	"time"
)

// Validate performs validation checks on the HL7AEcg instance.
func (e *HL7AEcg) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate ID
	if e.ID != nil {
		e.ID.Validate(ctx, vctx)
	} else {
		vctx.AddError(ErrMissingID)
	}

	// Validate Code
	if e.Code == nil {
		vctx.AddError(ErrMissingCode)
	} else {
		if e.Code.Code == "" {
			vctx.AddError(ErrMissingCode)
		}
		if e.Code.CodeSystem == "" {
			vctx.AddError(ErrMissingCodeSystem)
		}
	}

	// Validate EffectiveTime
	if e.EffectiveTime == nil {
		vctx.AddError(ErrMissingEffectiveTime)
	} else {
		e.EffectiveTime.Validate(ctx, vctx)
	}

	// Validate optional codes if present
	if e.ConfidentialityCode != nil {
		e.ConfidentialityCode.ValidateCode(ctx, vctx, "ConfidentialityCode")
	}
	if e.ReasonCode != nil {
		e.ReasonCode.ValidateCode(ctx, vctx, "ReasonCode")
	}

	// Validate ComponentOf structure if present
	if e.ComponentOf != nil {
		// Validate TimepointEvent and nested structures
		te := &e.ComponentOf.TimepointEvent
		sa := &te.ComponentOf.SubjectAssignment

		// Validate Subject within ComponentOf
		if err := sa.Subject.Validate(ctx, vctx); err != nil {
			return err
		}

		// Validate ClinicalTrial within ComponentOf
		if err := sa.ComponentOf.ClinicalTrial.Validate(ctx, vctx); err != nil {
			return err
		}
	} else {
		vctx.AddError(ErrMissingSubject)
	}

	// Validate direct Subject if present (alternative structure)
	if e.Subject != nil {
		if err := e.Subject.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	// Validate direct ClinicalTrial if present (alternative structure)
	if e.ClinicalTrial != nil {
		if err := e.ClinicalTrial.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	// Validate all Component (Series) elements
	for i := range e.Component {
		if err := e.Component[i].Series.Validate(ctx, vctx); err != nil {
			return err
		}
	}

	return nil
}

func (id *ID) Validate(ctx context.Context, vctx *ValidationContext) error {
	// Autocomplete empty Root with singleton if available
	if id == nil {
		id = &ID{}
	}
	if id.Root == "" {
		if instanceID != nil && instanceID.ID != "" {
			id.Root = instanceID.ID
		} else {
			vctx.AddError(ErrMissingID)
			return nil
		}
	}
	// Accept any non-empty Root value (OID, custom identifier, etc.)
	// No specific format validation required
	return nil
}

func (e *EffectiveTime) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if e.Low.Value != "" && !isValidTimestamp(e.Low.Value) {
		vctx.AddError(ErrInvalidTimeFormat)
	} else if e.High.Value != "" && !isValidTimestamp(e.High.Value) {
		vctx.AddError(ErrInvalidTimeFormat)
	}
	return nil
}

// isValidTimestamp checks if the string is a valid HL7 TS or Unix timestamp.
func isValidTimestamp(s string) bool {
	// Try HL7 TS format (YYYYMMDDHHmmss or with milliseconds)
	layouts := []string{
		"20060102150405",     // YYYYMMDDHHmmss
		"20060102150405.000", // with milliseconds
		"20060102",           // date only
		"200601",             // year + month
		"2006",               // year only
	}

	for _, layout := range layouts {
		if _, err := time.Parse(layout, s); err == nil {
			return true
		}
	}

	// Try Unix timestamp
	if sec, err := strconv.ParseInt(s, 10, 64); err == nil {
		t := time.Unix(sec, 0)
		min := time.Unix(0, 0)
		max := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
		return t.After(min) && t.Before(max)
	}

	return false
}
