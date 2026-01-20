package types

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper function to create a PhysicalQuantity AnnotationValue for tests
func makePQValue(value, unit string) *AnnotationValue {
	return &AnnotationValue{
		XsiType: "PQ",
		Typed: &PhysicalQuantity{
			XsiType: "PQ",
			Value:   value,
			Unit:    unit,
		},
	}
}

// TestAnnotationSet_Validate tests AnnotationSet validation
func TestAnnotationSet_Validate(t *testing.T) {
	ctx := context.Background()

	t.Run("valid annotation set", func(t *testing.T) {
		annSet := &AnnotationSet{
			ActivityTime: &Time{Value: "20250923103600"},
			Component: []AnnotationComponent{
				{
					Annotation: Annotation{
						Code: &Code[string, string]{
							Code:       "MDC_ECG_HEART_RATE",
							CodeSystem: "2.16.840.1.113883.6.24",
						},
						Value: makePQValue("57", "bpm"),
					},
				},
			},
		}

		vctx := NewValidationContext(false)
		err := annSet.Validate(ctx, vctx)
		assert.NoError(t, err)
	})

	t.Run("nil annotation set", func(t *testing.T) {
		var annSet *AnnotationSet
		vctx := NewValidationContext(false)
		err := annSet.Validate(ctx, vctx)
		assert.NoError(t, err)
	})

	t.Run("invalid activity time", func(t *testing.T) {
		annSet := &AnnotationSet{
			ActivityTime: &Time{Value: "invalid"},
		}

		vctx := NewValidationContext(false)
		err := annSet.Validate(ctx, vctx)
		assert.Error(t, err)
	})

	t.Run("annotation with empty code", func(t *testing.T) {
		annSet := &AnnotationSet{
			Component: []AnnotationComponent{
				{
					Annotation: Annotation{
						Code: &Code[string, string]{
							Code:       "",
							CodeSystem: "2.16.840.1.113883.6.24",
						},
						Value: makePQValue("57", "bpm"),
					},
				},
			},
		}

		vctx := NewValidationContext(false)
		err := annSet.Validate(ctx, vctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "code cannot be empty")
	})

	t.Run("annotation with empty value", func(t *testing.T) {
		annSet := &AnnotationSet{
			Component: []AnnotationComponent{
				{
					Annotation: Annotation{
						Code: &Code[string, string]{
							Code:       "MDC_ECG_HEART_RATE",
							CodeSystem: "2.16.840.1.113883.6.24",
						},
						Value: makePQValue("", "bpm"),
					},
				},
			},
		}

		vctx := NewValidationContext(false)
		err := annSet.Validate(ctx, vctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "value cannot be empty")
	})

	t.Run("annotation with invalid numeric value", func(t *testing.T) {
		annSet := &AnnotationSet{
			Component: []AnnotationComponent{
				{
					Annotation: Annotation{
						Code: &Code[string, string]{
							Code:       "MDC_ECG_HEART_RATE",
							CodeSystem: "2.16.840.1.113883.6.24",
						},
						Value: makePQValue("not-a-number", "bpm"),
					},
				},
			},
		}

		vctx := NewValidationContext(false)
		err := annSet.Validate(ctx, vctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be a valid number")
	})
}

// TestAnnotation_Validate tests Annotation validation
func TestAnnotation_Validate(t *testing.T) {
	ctx := context.Background()

	t.Run("valid simple annotation", func(t *testing.T) {
		ann := &Annotation{
			Code: &Code[string, string]{
				Code:       "MDC_ECG_HEART_RATE",
				CodeSystem: "2.16.840.1.113883.6.24",
			},
			Value: makePQValue("57", "bpm"),
		}

		vctx := NewValidationContext(false)
		err := ann.Validate(ctx, vctx)
		assert.NoError(t, err)
	})

	t.Run("valid vendor-specific annotation", func(t *testing.T) {
		ann := &Annotation{
			Code: &Code[string, string]{
				Code:           "MINDRAY_P_ONSET",
				CodeSystem:     "", // Empty for vendor codes
				CodeSystemName: "MINDRAY",
			},
			Value: makePQValue("234", "ms"),
		}

		vctx := NewValidationContext(false)
		err := ann.Validate(ctx, vctx)
		assert.NoError(t, err) // Should pass even with empty codeSystem
	})

	t.Run("valid nested annotations", func(t *testing.T) {
		ann := &Annotation{
			Code: &Code[string, string]{
				Code:       "MDC_ECG_TIME_PD_QTc",
				CodeSystem: "2.16.840.1.113883.6.24",
			},
			Component: []AnnotationComponent{
				{
					Annotation: Annotation{
						Code: &Code[string, string]{
							Code: "ECG_TIME_PD_QTcH",
						},
						Value: makePQValue("413", "ms"),
					},
				},
			},
		}

		vctx := NewValidationContext(false)
		err := ann.Validate(ctx, vctx)
		assert.NoError(t, err)
	})

	t.Run("nil annotation", func(t *testing.T) {
		var ann *Annotation
		vctx := NewValidationContext(false)
		err := ann.Validate(ctx, vctx)
		assert.NoError(t, err)
	})
}

// TestAnnotationSupportingROI_Validate tests SupportingROI validation
func TestAnnotationSupportingROI_Validate(t *testing.T) {
	ctx := context.Background()

	t.Run("valid supporting ROI", func(t *testing.T) {
		roi := &AnnotationSupportingROI{
			ClassCode: "ROIBND",
			Code: &Code[string, string]{
				Code:       "ROIPS",
				CodeSystem: "2.16.840.1.113883.5.4",
			},
			Component: []AnnotationBoundaryComponent{
				{
					Boundary: AnnotationBoundary{
						Code: Code[string, string]{
							Code:       "MDC_ECG_LEAD_I",
							CodeSystem: "2.16.840.1.113883.6.24",
						},
					},
				},
			},
		}

		vctx := NewValidationContext(false)
		err := roi.Validate(ctx, vctx)
		assert.NoError(t, err)
	})

	t.Run("invalid classCode", func(t *testing.T) {
		roi := &AnnotationSupportingROI{
			ClassCode: "INVALID",
			Code: &Code[string, string]{
				Code:       "ROIPS",
				CodeSystem: "2.16.840.1.113883.5.4",
			},
			Component: []AnnotationBoundaryComponent{
				{
					Boundary: AnnotationBoundary{
						Code: Code[string, string]{
							Code:       "MDC_ECG_LEAD_I",
							CodeSystem: "2.16.840.1.113883.6.24",
						},
					},
				},
			},
		}

		vctx := NewValidationContext(false)
		err := roi.Validate(ctx, vctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "should be 'ROIBND'")
	})

	t.Run("invalid ROI code", func(t *testing.T) {
		roi := &AnnotationSupportingROI{
			ClassCode: "ROIBND",
			Code: &Code[string, string]{
				Code:       "INVALID",
				CodeSystem: "2.16.840.1.113883.5.4",
			},
			Component: []AnnotationBoundaryComponent{
				{
					Boundary: AnnotationBoundary{
						Code: Code[string, string]{
							Code:       "MDC_ECG_LEAD_I",
							CodeSystem: "2.16.840.1.113883.6.24",
						},
					},
				},
			},
		}

		vctx := NewValidationContext(false)
		err := roi.Validate(ctx, vctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ROIPS")
	})

	t.Run("empty boundary lead code", func(t *testing.T) {
		roi := &AnnotationSupportingROI{
			ClassCode: "ROIBND",
			Code: &Code[string, string]{
				Code:       "ROIPS",
				CodeSystem: "2.16.840.1.113883.5.4",
			},
			Component: []AnnotationBoundaryComponent{
				{
					Boundary: AnnotationBoundary{
						Code: Code[string, string]{
							Code:       "",
							CodeSystem: "2.16.840.1.113883.6.24",
						},
					},
				},
			},
		}

		vctx := NewValidationContext(false)
		err := roi.Validate(ctx, vctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "lead code cannot be empty")
	})
}
