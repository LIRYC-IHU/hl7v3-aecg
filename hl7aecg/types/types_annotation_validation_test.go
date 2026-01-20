package types

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAddAnnotation_InputValidation tests input validation for AddAnnotation
func TestAddAnnotation_InputValidation(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}

	t.Run("valid annotation", func(t *testing.T) {
		idx := annSet.AddAnnotation("MDC_ECG_HEART_RATE", "2.16.840.1.113883.6.24", 57, "bpm")
		assert.NotEqual(t, -1, idx)
	})

	t.Run("empty code", func(t *testing.T) {
		idx := annSet.AddAnnotation("", "2.16.840.1.113883.6.24", 57, "bpm")
		assert.Equal(t, -1, idx, "Should reject empty code")
	})

	t.Run("empty codeSystem", func(t *testing.T) {
		idx := annSet.AddAnnotation("MDC_ECG_HEART_RATE", "", 57, "bpm")
		assert.Equal(t, -1, idx, "Should reject empty codeSystem")
	})

	t.Run("empty unit", func(t *testing.T) {
		idx := annSet.AddAnnotation("MDC_ECG_HEART_RATE", "2.16.840.1.113883.6.24", 57, "")
		assert.Equal(t, -1, idx, "Should reject empty unit")
	})

	t.Run("NaN value", func(t *testing.T) {
		idx := annSet.AddAnnotation("MDC_ECG_HEART_RATE", "2.16.840.1.113883.6.24", math.NaN(), "bpm")
		assert.Equal(t, -1, idx, "Should reject NaN value")
	})

	t.Run("positive infinity value", func(t *testing.T) {
		idx := annSet.AddAnnotation("MDC_ECG_HEART_RATE", "2.16.840.1.113883.6.24", math.Inf(1), "bpm")
		assert.Equal(t, -1, idx, "Should reject positive infinity")
	})

	t.Run("negative infinity value", func(t *testing.T) {
		idx := annSet.AddAnnotation("MDC_ECG_HEART_RATE", "2.16.840.1.113883.6.24", math.Inf(-1), "bpm")
		assert.Equal(t, -1, idx, "Should reject negative infinity")
	})

	t.Run("zero value is valid", func(t *testing.T) {
		idx := annSet.AddAnnotation("MDC_ECG_TIME_PD_QTc", "2.16.840.1.113883.6.24", 0, "ms")
		assert.NotEqual(t, -1, idx, "Should accept zero as valid value")
	})

	t.Run("negative value is valid", func(t *testing.T) {
		idx := annSet.AddAnnotation("AMPLITUDE", "test", -123, "uV")
		assert.NotEqual(t, -1, idx, "Should accept negative values")
	})
}

// TestAddAnnotationWithCodeSystemName_InputValidation tests input validation
func TestAddAnnotationWithCodeSystemName_InputValidation(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}

	t.Run("valid annotation", func(t *testing.T) {
		idx := annSet.AddAnnotationWithCodeSystemName("MINDRAY_ECG_HEART_RATE", "MINDRAY", 57, "bpm")
		assert.NotEqual(t, -1, idx)
	})

	t.Run("empty code", func(t *testing.T) {
		idx := annSet.AddAnnotationWithCodeSystemName("", "MINDRAY", 57, "bpm")
		assert.Equal(t, -1, idx, "Should reject empty code")
	})

	t.Run("empty codeSystemName", func(t *testing.T) {
		idx := annSet.AddAnnotationWithCodeSystemName("MINDRAY_ECG_HEART_RATE", "", 57, "bpm")
		assert.Equal(t, -1, idx, "Should reject empty codeSystemName")
	})

	t.Run("empty unit", func(t *testing.T) {
		idx := annSet.AddAnnotationWithCodeSystemName("MINDRAY_ECG_HEART_RATE", "MINDRAY", 57, "")
		assert.Equal(t, -1, idx, "Should reject empty unit")
	})

	t.Run("invalid float", func(t *testing.T) {
		idx := annSet.AddAnnotationWithCodeSystemName("MINDRAY_ECG_HEART_RATE", "MINDRAY", math.NaN(), "bpm")
		assert.Equal(t, -1, idx, "Should reject NaN value")
	})
}

// TestAddLeadAnnotation_InputValidation tests input validation for AddLeadAnnotation
func TestAddLeadAnnotation_InputValidation(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}

	t.Run("valid with codeSystem", func(t *testing.T) {
		idx := annSet.AddLeadAnnotation("MDC_ECG_LEAD_I", "MEASUREMENT_MATRIX", "test", "")
		assert.NotEqual(t, -1, idx)
	})

	t.Run("valid with codeSystemName", func(t *testing.T) {
		idx := annSet.AddLeadAnnotation("MDC_ECG_LEAD_I", "MINDRAY_MEASUREMENT_MATRIX", "", "MINDRAY")
		assert.NotEqual(t, -1, idx)
	})

	t.Run("empty leadCode", func(t *testing.T) {
		idx := annSet.AddLeadAnnotation("", "MEASUREMENT_MATRIX", "test", "")
		assert.Equal(t, -1, idx, "Should reject empty leadCode")
	})

	t.Run("empty matrixCode", func(t *testing.T) {
		idx := annSet.AddLeadAnnotation("MDC_ECG_LEAD_I", "", "test", "")
		assert.Equal(t, -1, idx, "Should reject empty matrixCode")
	})

	t.Run("both codeSystem and codeSystemName empty", func(t *testing.T) {
		idx := annSet.AddLeadAnnotation("MDC_ECG_LEAD_I", "MEASUREMENT_MATRIX", "", "")
		assert.Equal(t, -1, idx, "Should reject when both are empty")
	})
}

// TestAddNestedAnnotation_InputValidation tests input validation for nested annotations
func TestAddNestedAnnotation_InputValidation(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}
	parentIdx := annSet.AddAnnotation("MDC_ECG_TIME_PD_QTc", "2.16.840.1.113883.6.24", 0, "ms")
	parent := annSet.GetAnnotation(parentIdx)

	t.Run("valid nested annotation", func(t *testing.T) {
		idx := parent.AddNestedAnnotation("ECG_TIME_PD_QTcH", "", 413, "ms")
		assert.NotEqual(t, -1, idx)
	})

	t.Run("empty code", func(t *testing.T) {
		idx := parent.AddNestedAnnotation("", "", 413, "ms")
		assert.Equal(t, -1, idx, "Should reject empty code")
	})

	t.Run("empty unit", func(t *testing.T) {
		idx := parent.AddNestedAnnotation("ECG_TIME_PD_QTcH", "", 413, "")
		assert.Equal(t, -1, idx, "Should reject empty unit")
	})

	t.Run("invalid float", func(t *testing.T) {
		idx := parent.AddNestedAnnotation("ECG_TIME_PD_QTcH", "", math.NaN(), "ms")
		assert.Equal(t, -1, idx, "Should reject NaN value")
	})
}

// TestAddNestedAnnotationWithCodeSystemName_InputValidation tests input validation
func TestAddNestedAnnotationWithCodeSystemName_InputValidation(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}
	parentIdx := annSet.AddLeadAnnotation("MDC_ECG_LEAD_I", "MINDRAY_MEASUREMENT_MATRIX", "", "MINDRAY")
	parent := annSet.GetAnnotation(parentIdx)

	t.Run("valid nested annotation", func(t *testing.T) {
		idx := parent.AddNestedAnnotationWithCodeSystemName("MINDRAY_P_ONSET", "MINDRAY", 234, "ms")
		assert.NotEqual(t, -1, idx)
	})

	t.Run("empty code", func(t *testing.T) {
		idx := parent.AddNestedAnnotationWithCodeSystemName("", "MINDRAY", 234, "ms")
		assert.Equal(t, -1, idx, "Should reject empty code")
	})

	t.Run("empty codeSystemName", func(t *testing.T) {
		idx := parent.AddNestedAnnotationWithCodeSystemName("MINDRAY_P_ONSET", "", 234, "ms")
		assert.Equal(t, -1, idx, "Should reject empty codeSystemName")
	})

	t.Run("empty unit", func(t *testing.T) {
		idx := parent.AddNestedAnnotationWithCodeSystemName("MINDRAY_P_ONSET", "MINDRAY", 234, "")
		assert.Equal(t, -1, idx, "Should reject empty unit")
	})

	t.Run("invalid float", func(t *testing.T) {
		idx := parent.AddNestedAnnotationWithCodeSystemName("MINDRAY_P_ONSET", "MINDRAY", math.Inf(1), "ms")
		assert.Equal(t, -1, idx, "Should reject infinity value")
	})
}

// TestConvenienceMethods_InputValidation tests validation for convenience methods
func TestConvenienceMethods_InputValidation(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}

	t.Run("AddHeartRate with valid value", func(t *testing.T) {
		idx := annSet.AddHeartRate(60)
		assert.NotEqual(t, -1, idx)
	})

	t.Run("AddHeartRate with NaN", func(t *testing.T) {
		idx := annSet.AddHeartRate(math.NaN())
		assert.Equal(t, -1, idx, "Should reject NaN")
	})

	t.Run("AddPRInterval with valid value", func(t *testing.T) {
		idx := annSet.AddPRInterval(192)
		assert.NotEqual(t, -1, idx)
	})

	t.Run("AddQRSDuration with valid value", func(t *testing.T) {
		idx := annSet.AddQRSDuration(88)
		assert.NotEqual(t, -1, idx)
	})

	t.Run("AddQTInterval with valid value", func(t *testing.T) {
		idx := annSet.AddQTInterval(418)
		assert.NotEqual(t, -1, idx)
	})

	t.Run("AddQTcInterval with zero (for parent with nested)", func(t *testing.T) {
		idx := annSet.AddQTcInterval(0)
		assert.NotEqual(t, -1, idx, "Zero is valid for parent annotations")
	})
}
