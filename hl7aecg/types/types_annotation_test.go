package types

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAnnotationSet_Unmarshal tests basic annotation set unmarshalling
func TestAnnotationSet_Unmarshal(t *testing.T) {
	xmlData := `<annotationSet xmlns="urn:hl7-org:v3" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
		<activityTime value="20250923103600"/>
		<component>
			<annotation>
				<code code="MDC_ECG_HEART_RATE" codeSystem="2.16.840.1.113883.6.24" codeSystemName="MDC"/>
				<value xsi:type="PQ" value="57" unit="bpm"/>
			</annotation>
		</component>
	</annotationSet>`

	var annSet AnnotationSet
	err := xml.Unmarshal([]byte(xmlData), &annSet)
	require.NoError(t, err)

	// Verify activityTime
	require.NotNil(t, annSet.ActivityTime)
	assert.Equal(t, "20250923103600", annSet.ActivityTime.Value)

	// Verify component
	require.Len(t, annSet.Component, 1)
	ann := annSet.Component[0].Annotation

	// Verify code
	require.NotNil(t, ann.Code)
	assert.Equal(t, "MDC_ECG_HEART_RATE", ann.Code.Code)
	assert.Equal(t, "2.16.840.1.113883.6.24", ann.Code.CodeSystem)
	assert.Equal(t, "MDC", ann.Code.CodeSystemName)

	// Verify value
	require.NotNil(t, ann.Value)
	assert.Equal(t, "57", ann.Value.Value)
	assert.Equal(t, "bpm", ann.Value.Unit)

	// Test GetValueFloat
	val, ok := ann.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(57), val)
}

// TestAnnotationSet_Unmarshal_MultipleAnnotations tests multiple global annotations
func TestAnnotationSet_Unmarshal_MultipleAnnotations(t *testing.T) {
	xmlData := `<annotationSet xmlns="urn:hl7-org:v3" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
		<activityTime value="20250923103600"/>
		<component>
			<annotation>
				<code code="MDC_ECG_HEART_RATE" codeSystem="2.16.840.1.113883.6.24" codeSystemName="MDC"/>
				<value xsi:type="PQ" value="57" unit="bpm"/>
			</annotation>
		</component>
		<component>
			<annotation>
				<code code="MDC_ECG_TIME_PD_PR" codeSystem="2.16.840.1.113883.6.24" codeSystemName="MDC"/>
				<value xsi:type="PQ" value="192" unit="ms"/>
			</annotation>
		</component>
		<component>
			<annotation>
				<code code="MDC_ECG_TIME_PD_QRS" codeSystem="2.16.840.1.113883.6.24" codeSystemName="MDC"/>
				<value xsi:type="PQ" value="88" unit="ms"/>
			</annotation>
		</component>
		<component>
			<annotation>
				<code code="MDC_ECG_TIME_PD_QT" codeSystem="2.16.840.1.113883.6.24" codeSystemName="MDC"/>
				<value xsi:type="PQ" value="418" unit="ms"/>
			</annotation>
		</component>
	</annotationSet>`

	var annSet AnnotationSet
	err := xml.Unmarshal([]byte(xmlData), &annSet)
	require.NoError(t, err)

	require.Len(t, annSet.Component, 4)

	// Test GetAnnotationByCode helper
	hrAnn := annSet.GetAnnotationByCode("MDC_ECG_HEART_RATE")
	require.NotNil(t, hrAnn)
	val, ok := hrAnn.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(57), val)
	assert.Equal(t, "bpm", hrAnn.GetValueUnit())

	prAnn := annSet.GetAnnotationByCode("MDC_ECG_TIME_PD_PR")
	require.NotNil(t, prAnn)
	val, ok = prAnn.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(192), val)
	assert.Equal(t, "ms", prAnn.GetValueUnit())

	qrsAnn := annSet.GetAnnotationByCode("MDC_ECG_TIME_PD_QRS")
	require.NotNil(t, qrsAnn)
	val, ok = qrsAnn.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(88), val)

	qtAnn := annSet.GetAnnotationByCode("MDC_ECG_TIME_PD_QT")
	require.NotNil(t, qtAnn)
	val, ok = qtAnn.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(418), val)

	// Test GetAnnotationByCode with non-existent code
	nilAnn := annSet.GetAnnotationByCode("NON_EXISTENT_CODE")
	assert.Nil(t, nilAnn)
}

// TestAnnotation_NestedAnnotations tests annotations with nested components
func TestAnnotation_NestedAnnotations(t *testing.T) {
	xmlData := `<annotationSet xmlns="urn:hl7-org:v3" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
		<activityTime value="20250923103600"/>
		<component>
			<annotation>
				<code code="MDC_ECG_TIME_PD_QTc" codeSystem="2.16.840.1.113883.6.24" codeSystemName="MDC"/>
				<component>
					<annotation>
						<code code="ECG_TIME_PD_QTcH" codeSystem="" codeSystemName=""/>
						<value xsi:type="PQ" value="413" unit="ms"/>
					</annotation>
				</component>
			</annotation>
		</component>
	</annotationSet>`

	var annSet AnnotationSet
	err := xml.Unmarshal([]byte(xmlData), &annSet)
	require.NoError(t, err)

	require.Len(t, annSet.Component, 1)
	qtcAnn := annSet.Component[0].Annotation

	assert.Equal(t, "MDC_ECG_TIME_PD_QTc", qtcAnn.Code.Code)

	// Test nested annotation
	require.Len(t, qtcAnn.Component, 1)
	nestedAnn := qtcAnn.Component[0].Annotation
	assert.Equal(t, "ECG_TIME_PD_QTcH", nestedAnn.Code.Code)

	val, ok := nestedAnn.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(413), val)

	// Test GetNestedAnnotationByCode helper
	foundNested := qtcAnn.GetNestedAnnotationByCode("ECG_TIME_PD_QTcH")
	require.NotNil(t, foundNested)
	assert.Equal(t, "ECG_TIME_PD_QTcH", foundNested.Code.Code)
}

// TestAnnotation_LeadSpecific tests lead-specific annotations with supportingROI
func TestAnnotation_LeadSpecific(t *testing.T) {
	xmlData := `<annotationSet xmlns="urn:hl7-org:v3" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
		<activityTime value="20250923103600"/>
		<component>
			<annotation>
				<code code="MEASUREMENT_MATRIX" codeSystem="" codeSystemName=""/>
				<support>
					<supportingROI classCode="ROIBND">
						<code code="ROIPS" codeSystem="2.16.840.1.113883.5.4" codeSystemName="HL7V3"/>
						<component>
							<boundary>
								<code code="MDC_ECG_LEAD_I" codeSystem="2.16.840.1.113883.6.24" codeSystemName="MDC"/>
							</boundary>
						</component>
					</supportingROI>
				</support>
				<component>
					<annotation>
						<code code="P_ONSET" codeSystem="" codeSystemName=""/>
						<value xsi:type="PQ" value="234" unit="ms"/>
					</annotation>
				</component>
				<component>
					<annotation>
						<code code="R_AMP" codeSystem="" codeSystemName=""/>
						<value xsi:type="PQ" value="535" unit="uV"/>
					</annotation>
				</component>
			</annotation>
		</component>
	</annotationSet>`

	var annSet AnnotationSet
	err := xml.Unmarshal([]byte(xmlData), &annSet)
	require.NoError(t, err)

	require.Len(t, annSet.Component, 1)
	leadAnn := annSet.Component[0].Annotation

	// Verify the annotation code
	assert.Equal(t, "MEASUREMENT_MATRIX", leadAnn.Code.Code)

	// Verify support/supportingROI
	require.NotNil(t, leadAnn.Support)
	roi := leadAnn.Support.SupportingROI
	assert.Equal(t, "ROIBND", roi.ClassCode)
	assert.Equal(t, "ROIPS", roi.Code.Code)

	// Verify boundary (lead identification)
	require.Len(t, roi.Component, 1)
	boundary := roi.Component[0].Boundary
	assert.Equal(t, "MDC_ECG_LEAD_I", string(boundary.Code.Code))

	// Verify nested measurements
	require.Len(t, leadAnn.Component, 2)

	pOnset := leadAnn.GetNestedAnnotationByCode("P_ONSET")
	require.NotNil(t, pOnset)
	val, ok := pOnset.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(234), val)
	assert.Equal(t, "ms", pOnset.GetValueUnit())

	rAmp := leadAnn.GetNestedAnnotationByCode("R_AMP")
	require.NotNil(t, rAmp)
	val, ok = rAmp.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(535), val)
	assert.Equal(t, "uV", rAmp.GetValueUnit())

	// Test GetLeadAnnotations helper
	leadI := annSet.GetLeadAnnotations("MDC_ECG_LEAD_I")
	require.NotNil(t, leadI)
	assert.Equal(t, "MEASUREMENT_MATRIX", leadI.Code.Code)
}

// TestSubjectOf_Unmarshal tests SubjectOf with annotationSet
func TestSubjectOf_Unmarshal(t *testing.T) {
	xmlData := `<subjectOf xmlns="urn:hl7-org:v3" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
		<annotationSet>
			<activityTime value="20250923103600"/>
			<component>
				<annotation>
					<code code="MDC_ECG_HEART_RATE" codeSystem="2.16.840.1.113883.6.24"/>
					<value xsi:type="PQ" value="72" unit="bpm"/>
				</annotation>
			</component>
		</annotationSet>
	</subjectOf>`

	var subjectOf SubjectOf
	err := xml.Unmarshal([]byte(xmlData), &subjectOf)
	require.NoError(t, err)

	require.NotNil(t, subjectOf.AnnotationSet)
	require.NotNil(t, subjectOf.AnnotationSet.ActivityTime)
	assert.Equal(t, "20250923103600", subjectOf.AnnotationSet.ActivityTime.Value)

	hrAnn := subjectOf.AnnotationSet.GetAnnotationByCode("MDC_ECG_HEART_RATE")
	require.NotNil(t, hrAnn)
	val, ok := hrAnn.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(72), val)
}

// TestPhysicalQuantity_GetValueFloat tests the GetValueFloat method
func TestPhysicalQuantity_GetValueFloat(t *testing.T) {
	tests := []struct {
		name     string
		pq       *PhysicalQuantity
		expected float64
		ok       bool
	}{
		{
			name:     "integer value",
			pq:       &PhysicalQuantity{Value: "57", Unit: "bpm"},
			expected: 57,
			ok:       true,
		},
		{
			name:     "float value",
			pq:       &PhysicalQuantity{Value: "0.008", Unit: "mV"},
			expected: 0.008,
			ok:       true,
		},
		{
			name:     "negative value",
			pq:       &PhysicalQuantity{Value: "-165", Unit: "uV"},
			expected: -165,
			ok:       true,
		},
		{
			name:     "empty value",
			pq:       &PhysicalQuantity{Value: "", Unit: "ms"},
			expected: 0,
			ok:       false,
		},
		{
			name:     "nil pointer",
			pq:       nil,
			expected: 0,
			ok:       false,
		},
		{
			name:     "invalid value",
			pq:       &PhysicalQuantity{Value: "not_a_number", Unit: "ms"},
			expected: 0,
			ok:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, ok := tt.pq.GetValueFloat()
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.expected, val)
		})
	}
}

// TestAnnotationSet_NilSafety tests nil safety of helper methods
func TestAnnotationSet_NilSafety(t *testing.T) {
	var nilSet *AnnotationSet
	assert.Nil(t, nilSet.GetAnnotationByCode("MDC_ECG_HEART_RATE"))
	assert.Nil(t, nilSet.GetLeadAnnotations("MDC_ECG_LEAD_I"))

	var nilAnn *Annotation
	val, ok := nilAnn.GetValueFloat()
	assert.False(t, ok)
	assert.Equal(t, float64(0), val)
	assert.Equal(t, "", nilAnn.GetValueUnit())
	assert.Nil(t, nilAnn.GetNestedAnnotationByCode("test"))
}

// =============================================================================
// Marshalling Tests (Story U2)
// =============================================================================

// TestAnnotationSet_Marshal_GlobalAnnotation tests marshalling a single global annotation
func TestAnnotationSet_Marshal_GlobalAnnotation(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}

	// Add heart rate annotation
	annSet.AddAnnotation("MDC_ECG_HEART_RATE", "2.16.840.1.113883.6.24", 57, "bpm")

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(annSet, "", "  ")
	require.NoError(t, err)

	// Unmarshal to verify round-trip
	var unmarshaledSet AnnotationSet
	err = xml.Unmarshal(xmlData, &unmarshaledSet)
	require.NoError(t, err)

	// Verify activityTime
	require.NotNil(t, unmarshaledSet.ActivityTime)
	assert.Equal(t, "20250923103600", unmarshaledSet.ActivityTime.Value)

	// Verify annotation
	hrAnn := unmarshaledSet.GetAnnotationByCode("MDC_ECG_HEART_RATE")
	require.NotNil(t, hrAnn)
	val, ok := hrAnn.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(57), val)
	assert.Equal(t, "bpm", hrAnn.GetValueUnit())
}

// TestAnnotationSet_Marshal_MultipleAnnotations tests marshalling multiple global annotations
func TestAnnotationSet_Marshal_MultipleAnnotations(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}

	// Add multiple annotations using convenience methods
	annSet.AddHeartRate(57)
	annSet.AddPRInterval(192)
	annSet.AddQRSDuration(88)
	annSet.AddQTInterval(418)

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(annSet, "", "  ")
	require.NoError(t, err)

	// Unmarshal to verify round-trip
	var unmarshaledSet AnnotationSet
	err = xml.Unmarshal(xmlData, &unmarshaledSet)
	require.NoError(t, err)

	// Verify all annotations
	require.Len(t, unmarshaledSet.Component, 4)

	hrAnn := unmarshaledSet.GetAnnotationByCode("MDC_ECG_HEART_RATE")
	require.NotNil(t, hrAnn)
	val, _ := hrAnn.GetValueFloat()
	assert.Equal(t, float64(57), val)

	prAnn := unmarshaledSet.GetAnnotationByCode("MDC_ECG_TIME_PD_PR")
	require.NotNil(t, prAnn)
	val, _ = prAnn.GetValueFloat()
	assert.Equal(t, float64(192), val)

	qrsAnn := unmarshaledSet.GetAnnotationByCode("MDC_ECG_TIME_PD_QRS")
	require.NotNil(t, qrsAnn)
	val, _ = qrsAnn.GetValueFloat()
	assert.Equal(t, float64(88), val)

	qtAnn := unmarshaledSet.GetAnnotationByCode("MDC_ECG_TIME_PD_QT")
	require.NotNil(t, qtAnn)
	val, _ = qtAnn.GetValueFloat()
	assert.Equal(t, float64(418), val)
}

// TestAnnotationSet_Marshal_NestedAnnotation tests marshalling with nested annotations
func TestAnnotationSet_Marshal_NestedAnnotation(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}

	// Add QTc with nested correction formula
	qtcAnn := annSet.AddQTcInterval(0)
	qtcAnn.AddNestedAnnotationWithCodeSystemName("ECG_TIME_PD_QTcH", "", 413, "ms")

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(annSet, "", "  ")
	require.NoError(t, err)

	// Unmarshal to verify round-trip
	var unmarshaledSet AnnotationSet
	err = xml.Unmarshal(xmlData, &unmarshaledSet)
	require.NoError(t, err)

	// Verify parent annotation
	qtcAnn = unmarshaledSet.GetAnnotationByCode("MDC_ECG_TIME_PD_QTc")
	require.NotNil(t, qtcAnn)

	// Verify nested annotation
	nestedAnn := qtcAnn.GetNestedAnnotationByCode("ECG_TIME_PD_QTcH")
	require.NotNil(t, nestedAnn)
	val, ok := nestedAnn.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(413), val)
	assert.Equal(t, "ms", nestedAnn.GetValueUnit())
}

// TestAnnotationSet_Marshal_LeadSpecific tests marshalling lead-specific annotations
func TestAnnotationSet_Marshal_LeadSpecific(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}

	// Add lead-specific annotation for Lead I
	leadIAnn := annSet.AddLeadAnnotation("MDC_ECG_LEAD_I", "MEASUREMENT_MATRIX", "test", "")
	leadIAnn.AddNestedAnnotationWithCodeSystemName("P_ONSET", "", 234, "ms")
	leadIAnn.AddNestedAnnotationWithCodeSystemName("R_AMP", "", 535, "uV")

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(annSet, "", "  ")
	require.NoError(t, err)

	// Unmarshal to verify round-trip
	var unmarshaledSet AnnotationSet
	err = xml.Unmarshal(xmlData, &unmarshaledSet)
	require.NoError(t, err)

	// Verify lead annotation
	leadAnn := unmarshaledSet.GetLeadAnnotations("MDC_ECG_LEAD_I")
	require.NotNil(t, leadAnn)

	// Verify supportingROI structure
	require.NotNil(t, leadAnn.Support)
	assert.Equal(t, "ROIBND", leadAnn.Support.SupportingROI.ClassCode)

	// Verify nested measurements
	pOnset := leadAnn.GetNestedAnnotationByCode("P_ONSET")
	require.NotNil(t, pOnset)
	val, _ := pOnset.GetValueFloat()
	assert.Equal(t, float64(234), val)

	rAmp := leadAnn.GetNestedAnnotationByCode("R_AMP")
	require.NotNil(t, rAmp)
	val, _ = rAmp.GetValueFloat()
	assert.Equal(t, float64(535), val)
}

// TestAnnotationSet_Marshal_VendorSpecific tests vendor-specific annotations with codeSystemName
func TestAnnotationSet_Marshal_VendorSpecific(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}

	// Add vendor-specific annotations
	annSet.AddAnnotationWithCodeSystemName("ECG_P_AXIS", "", 60, "deg")
	annSet.AddAnnotationWithCodeSystemName("ECG_QRS_AXIS", "", 64, "deg")

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(annSet, "", "  ")
	require.NoError(t, err)

	// Unmarshal to verify round-trip
	var unmarshaledSet AnnotationSet
	err = xml.Unmarshal(xmlData, &unmarshaledSet)
	require.NoError(t, err)

	// Verify vendor-specific annotations
	pAxis := unmarshaledSet.GetAnnotationByCode("ECG_P_AXIS")
	require.NotNil(t, pAxis)
	assert.Equal(t, "", pAxis.Code.CodeSystemName)
	assert.Equal(t, "", pAxis.Code.CodeSystem) // Empty for vendor codes
	val, _ := pAxis.GetValueFloat()
	assert.Equal(t, float64(60), val)
}

// TestSeries_InitAnnotationSet tests Series.InitAnnotationSet method
func TestSeries_InitAnnotationSet(t *testing.T) {
	series := &Series{
		Code: &Code[SeriesTypeCode, CodeSystemOID]{
			Code:       RHYTHM_CODE,
			CodeSystem: HL7_ActCode_OID,
		},
	}

	// Initialize annotation set
	annSet := series.InitAnnotationSet("20250923103600")
	require.NotNil(t, annSet)

	// Verify the annotation set was added to SubjectOf
	require.Len(t, series.SubjectOf, 1)
	assert.NotNil(t, series.SubjectOf[0].AnnotationSet)
	assert.Equal(t, "20250923103600", series.SubjectOf[0].AnnotationSet.ActivityTime.Value)

	// Add annotations
	annSet.AddHeartRate(72)

	// Verify annotation was added
	hrAnn := annSet.GetAnnotationByCode("MDC_ECG_HEART_RATE")
	require.NotNil(t, hrAnn)
	val, _ := hrAnn.GetValueFloat()
	assert.Equal(t, float64(72), val)
}

// TestFormatFloat tests the formatFloat helper function
func TestFormatFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected string
	}{
		{"integer", 57, "57"},
		{"float with decimals", 1.008, "1.008"},
		{"negative integer", -165, "-165"},
		{"negative float", -0.565, "-0.565"},
		{"zero", 0, "0"},
		{"small decimal", 0.008, "0.008"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatFloat(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
