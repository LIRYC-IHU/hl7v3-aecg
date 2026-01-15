package hl7aecg

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAnnotations_ReferenceFile tests annotation unmarshalling with the real reference file
func TestAnnotations_ReferenceFile(t *testing.T) {
	// Read the reference file
	filePath := "../25060897140_23092025103550.xml"
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Skipf("Reference file not found at %s: %v", filePath, err)
	}

	h := NewHl7xml("")
	err = h.Unmarshal(data)
	require.NoError(t, err)

	// Verify document metadata
	require.NotNil(t, h.HL7AEcg.ID)
	assert.Equal(t, "755.3045256.2025923.103550", h.HL7AEcg.ID.Root)

	require.NotNil(t, h.HL7AEcg.Code)
	assert.Equal(t, "93000", string(h.HL7AEcg.Code.Code))

	// Verify series exists
	require.NotEmpty(t, h.HL7AEcg.Component, "Expected at least one component (series)")

	series := h.HL7AEcg.Component[0].Series
	require.NotNil(t, series)
	require.NotNil(t, series.Code)
	assert.Equal(t, "RHYTHM", string(series.Code.Code))

	// Verify subjectOf (annotation set) exists
	require.NotEmpty(t, series.SubjectOf, "Expected at least one SubjectOf with annotations")

	subjectOf := series.SubjectOf[0]
	require.NotNil(t, subjectOf.AnnotationSet)

	annSet := subjectOf.AnnotationSet

	// Verify activityTime
	require.NotNil(t, annSet.ActivityTime)
	assert.Equal(t, "20250923103600", annSet.ActivityTime.Value)

	// Verify we have annotations
	require.NotEmpty(t, annSet.Component, "Expected annotation components")

	// Test global annotations using helper methods
	t.Run("Heart Rate", func(t *testing.T) {
		hrAnn := annSet.GetAnnotationByCode("MDC_ECG_HEART_RATE")
		require.NotNil(t, hrAnn, "MDC_ECG_HEART_RATE annotation not found")
		val, ok := hrAnn.GetValueFloat()
		assert.True(t, ok)
		assert.Equal(t, float64(57), val)
		assert.Equal(t, "bpm", hrAnn.GetValueUnit())
	})

	t.Run("PR Interval", func(t *testing.T) {
		prAnn := annSet.GetAnnotationByCode("MDC_ECG_TIME_PD_PR")
		require.NotNil(t, prAnn, "MDC_ECG_TIME_PD_PR annotation not found")
		val, ok := prAnn.GetValueFloat()
		assert.True(t, ok)
		assert.Equal(t, float64(192), val)
		assert.Equal(t, "ms", prAnn.GetValueUnit())
	})

	t.Run("QRS Duration", func(t *testing.T) {
		qrsAnn := annSet.GetAnnotationByCode("MDC_ECG_TIME_PD_QRS")
		require.NotNil(t, qrsAnn, "MDC_ECG_TIME_PD_QRS annotation not found")
		val, ok := qrsAnn.GetValueFloat()
		assert.True(t, ok)
		assert.Equal(t, float64(88), val)
		assert.Equal(t, "ms", qrsAnn.GetValueUnit())
	})

	t.Run("QT Interval", func(t *testing.T) {
		qtAnn := annSet.GetAnnotationByCode("MDC_ECG_TIME_PD_QT")
		require.NotNil(t, qtAnn, "MDC_ECG_TIME_PD_QT annotation not found")
		val, ok := qtAnn.GetValueFloat()
		assert.True(t, ok)
		assert.Equal(t, float64(418), val)
		assert.Equal(t, "ms", qtAnn.GetValueUnit())
	})

	t.Run("QTc with nested annotation", func(t *testing.T) {
		qtcAnn := annSet.GetAnnotationByCode("MDC_ECG_TIME_PD_QTc")
		require.NotNil(t, qtcAnn, "MDC_ECG_TIME_PD_QTc annotation not found")

		// QTc should have a nested MINDRAY annotation
		require.NotEmpty(t, qtcAnn.Component, "QTc should have nested annotations")

		mindrayQtc := qtcAnn.GetNestedAnnotationByCode("MINDRAY_ECG_TIME_PD_QTcH")
		require.NotNil(t, mindrayQtc, "MINDRAY_ECG_TIME_PD_QTcH nested annotation not found")
		val, ok := mindrayQtc.GetValueFloat()
		assert.True(t, ok)
		assert.Equal(t, float64(413), val)
	})

	// Test lead-specific annotations
	t.Run("Lead I annotations", func(t *testing.T) {
		leadIAnn := annSet.GetLeadAnnotations("MDC_ECG_LEAD_I")
		require.NotNil(t, leadIAnn, "Lead I annotations not found")

		// Verify support/supportingROI structure
		require.NotNil(t, leadIAnn.Support)
		assert.Equal(t, "ROIBND", leadIAnn.Support.SupportingROI.ClassCode)

		// Verify we have nested measurements
		require.NotEmpty(t, leadIAnn.Component)

		// Test a specific nested annotation
		pOnset := leadIAnn.GetNestedAnnotationByCode("MINDRAY_P_ONSET")
		require.NotNil(t, pOnset, "MINDRAY_P_ONSET not found for Lead I")
		val, ok := pOnset.GetValueFloat()
		assert.True(t, ok)
		assert.Equal(t, float64(234), val)

		rAmp := leadIAnn.GetNestedAnnotationByCode("MINDRAY_R_AMP")
		require.NotNil(t, rAmp, "MINDRAY_R_AMP not found for Lead I")
		val, ok = rAmp.GetValueFloat()
		assert.True(t, ok)
		assert.Equal(t, float64(535), val)
		assert.Equal(t, "uV", rAmp.GetValueUnit())
	})

	t.Run("Lead II annotations", func(t *testing.T) {
		leadIIAnn := annSet.GetLeadAnnotations("MDC_ECG_LEAD_II")
		require.NotNil(t, leadIIAnn, "Lead II annotations not found")

		rAmp := leadIIAnn.GetNestedAnnotationByCode("MINDRAY_R_AMP")
		require.NotNil(t, rAmp, "MINDRAY_R_AMP not found for Lead II")
		val, ok := rAmp.GetValueFloat()
		assert.True(t, ok)
		assert.Equal(t, float64(752), val)
	})

	t.Run("Lead V1 annotations", func(t *testing.T) {
		leadV1Ann := annSet.GetLeadAnnotations("MDC_ECG_LEAD_V1")
		require.NotNil(t, leadV1Ann, "Lead V1 annotations not found")

		sAmp := leadV1Ann.GetNestedAnnotationByCode("MINDRAY_S_AMP")
		require.NotNil(t, sAmp, "MINDRAY_S_AMP not found for Lead V1")
		val, ok := sAmp.GetValueFloat()
		assert.True(t, ok)
		assert.Equal(t, float64(-565), val) // S wave is negative
	})
}

// TestAnnotations_AllLeads verifies all 12 standard leads are accessible
func TestAnnotations_AllLeads(t *testing.T) {
	filePath := "../25060897140_23092025103550.xml"
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Skipf("Reference file not found at %s: %v", filePath, err)
	}

	h := NewHl7xml("")
	err = h.Unmarshal(data)
	require.NoError(t, err)

	series := h.HL7AEcg.Component[0].Series
	annSet := series.SubjectOf[0].AnnotationSet

	// Standard 12 leads
	standardLeads := []string{
		"MDC_ECG_LEAD_I", "MDC_ECG_LEAD_II", "MDC_ECG_LEAD_III",
		"MDC_ECG_LEAD_AVR", "MDC_ECG_LEAD_AVL", "MDC_ECG_LEAD_AVF",
		"MDC_ECG_LEAD_V1", "MDC_ECG_LEAD_V2", "MDC_ECG_LEAD_V3",
		"MDC_ECG_LEAD_V4", "MDC_ECG_LEAD_V5", "MDC_ECG_LEAD_V6",
	}

	for _, lead := range standardLeads {
		t.Run(lead, func(t *testing.T) {
			leadAnn := annSet.GetLeadAnnotations(lead)
			require.NotNil(t, leadAnn, "%s annotations not found", lead)
			require.NotEmpty(t, leadAnn.Component, "%s should have measurements", lead)
		})
	}
}
