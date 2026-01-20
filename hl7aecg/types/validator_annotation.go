package types

import (
	"context"
	"fmt"
)

// Validate validates an AnnotationSet and all its components.
func (as *AnnotationSet) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if as == nil {
		return nil
	}

	// Validate activityTime if present (basic format check)
	if as.ActivityTime != nil {
		if as.ActivityTime.Value == "" {
			vctx.AddError(NewValidationError(
				"annotationSet.activityTime",
				"Activity time value cannot be empty",
			))
		}
		// Basic format check - should be YYYYMMDDHHmmss or YYYYMMDDHHmmss.SSS
		if len(as.ActivityTime.Value) < 14 {
			vctx.AddError(NewValidationError(
				"annotationSet.activityTime",
				"Activity time must be in format YYYYMMDDHHmmss or YYYYMMDDHHmmss.SSS",
			))
		}
	}

	// Validate each annotation component
	for i := range as.Component {
		if err := as.Component[i].Annotation.Validate(ctx, vctx); err != nil {
			vctx.AddError(fmt.Errorf("annotationSet.component[%d]: %w", i, err))
		}
	}

	return vctx.GetError()
}

// Validate validates an Annotation and all its nested components.
func (a *Annotation) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	if a == nil {
		return nil
	}

	// Validate code if present
	if a.Code != nil {
		if a.Code.Code == "" {
			vctx.AddError(NewValidationError(
				"annotation.code",
				"Annotation code cannot be empty",
			))
		}
		// Note: CodeSystem can be empty for vendor-specific codes
	}

	// Validate value if present
	if a.Value != nil {
		// Validate based on value type (PQ or ST)
		if a.Value.IsPQ() {
			// Physical Quantity validation
			pq, _ := a.Value.Typed.(*PhysicalQuantity)
			if pq != nil {
				if pq.Value == "" {
					vctx.AddError(NewValidationError(
						"annotation.value",
						"Annotation value cannot be empty",
					))
				} else if _, ok := a.Value.GetValueFloat(); !ok {
					vctx.AddError(NewValidationError(
						"annotation.value",
						"Annotation PQ value must be a valid number",
					))
				}
			}
		} else if a.Value.IsST() {
			// String value validation
			if text, ok := a.Value.GetText(); !ok || text == "" {
				vctx.AddError(NewValidationError(
					"annotation.value",
					"Annotation ST value cannot be empty",
				))
			}
		}
	}

	// Validate support/supportingROI if present
	if a.Support != nil {
		if err := a.Support.SupportingROI.Validate(ctx, vctx); err != nil {
			vctx.AddError(fmt.Errorf("annotation.support.supportingROI: %w", err))
		}
	}

	// Validate nested annotation components
	for i := range a.Component {
		if err := a.Component[i].Annotation.Validate(ctx, vctx); err != nil {
			vctx.AddError(fmt.Errorf("annotation.component[%d]: %w", i, err))
		}
	}

	return vctx.GetError()
}

// Validate validates an AnnotationSupportingROI.
func (roi *AnnotationSupportingROI) Validate(ctx context.Context, vctx *ValidationContext) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Validate classCode (typically "ROIBND")
	if roi.ClassCode != "" && roi.ClassCode != "ROIBND" {
		vctx.AddError(NewValidationError(
			"supportingROI.classCode",
			"SupportingROI classCode should be 'ROIBND'",
		))
	}

	// Validate code if present
	if roi.Code != nil {
		if roi.Code.Code != "ROIPS" && roi.Code.Code != "ROIFS" {
			vctx.AddError(NewValidationError(
				"supportingROI.code",
				"SupportingROI code should be 'ROIPS' (partially specified) or 'ROIFS' (fully specified)",
			))
		}
	}

	// Validate boundary components
	for i := range roi.Component {
		if roi.Component[i].Boundary.Code.Code == "" {
			vctx.AddError(NewValidationError(
				fmt.Sprintf("supportingROI.component[%d].boundary.code", i),
				"Boundary lead code cannot be empty",
			))
		}
	}

	return vctx.GetError()
}
