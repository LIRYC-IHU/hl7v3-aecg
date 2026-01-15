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
						<code code="MINDRAY_ECG_TIME_PD_QTcH" codeSystem="" codeSystemName="MINDRAY"/>
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
	assert.Equal(t, "MINDRAY_ECG_TIME_PD_QTcH", nestedAnn.Code.Code)

	val, ok := nestedAnn.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(413), val)

	// Test GetNestedAnnotationByCode helper
	foundNested := qtcAnn.GetNestedAnnotationByCode("MINDRAY_ECG_TIME_PD_QTcH")
	require.NotNil(t, foundNested)
	assert.Equal(t, "MINDRAY_ECG_TIME_PD_QTcH", foundNested.Code.Code)
}

// TestAnnotation_LeadSpecific tests lead-specific annotations with supportingROI
func TestAnnotation_LeadSpecific(t *testing.T) {
	xmlData := `<annotationSet xmlns="urn:hl7-org:v3" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
		<activityTime value="20250923103600"/>
		<component>
			<annotation>
				<code code="MINDRAY_MEASUREMENT_MATRIX" codeSystem="" codeSystemName="MINDRAY"/>
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
						<code code="MINDRAY_P_ONSET" codeSystem="" codeSystemName="MINDRAY"/>
						<value xsi:type="PQ" value="234" unit="ms"/>
					</annotation>
				</component>
				<component>
					<annotation>
						<code code="MINDRAY_R_AMP" codeSystem="" codeSystemName="MINDRAY"/>
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
	assert.Equal(t, "MINDRAY_MEASUREMENT_MATRIX", leadAnn.Code.Code)

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

	pOnset := leadAnn.GetNestedAnnotationByCode("MINDRAY_P_ONSET")
	require.NotNil(t, pOnset)
	val, ok := pOnset.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(234), val)
	assert.Equal(t, "ms", pOnset.GetValueUnit())

	rAmp := leadAnn.GetNestedAnnotationByCode("MINDRAY_R_AMP")
	require.NotNil(t, rAmp)
	val, ok = rAmp.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(535), val)
	assert.Equal(t, "uV", rAmp.GetValueUnit())

	// Test GetLeadAnnotations helper
	leadI := annSet.GetLeadAnnotations("MDC_ECG_LEAD_I")
	require.NotNil(t, leadI)
	assert.Equal(t, "MINDRAY_MEASUREMENT_MATRIX", leadI.Code.Code)
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
