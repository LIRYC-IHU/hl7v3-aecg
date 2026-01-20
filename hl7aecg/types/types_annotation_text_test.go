package types

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTextAnnotation_Marshal tests marshalling text annotations (xsi:type="ST")
func TestTextAnnotation_Marshal(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}

	// Add text interpretation annotation
	annSet.AddTextAnnotation("MDC_ECG_INTERPRETATION", "2.16.840.1.113883.6.24", "Rythme sinusal avec ESA")

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(annSet, "", "  ")
	require.NoError(t, err)

	// Verify XML contains xsi:type="ST" and text content
	xmlStr := string(xmlData)
	assert.Contains(t, xmlStr, `xsi:type="ST"`)
	assert.Contains(t, xmlStr, "Rythme sinusal avec ESA")
	assert.Contains(t, xmlStr, "MDC_ECG_INTERPRETATION")
}

// TestTextAnnotation_Unmarshal tests unmarshalling text annotations
func TestTextAnnotation_Unmarshal(t *testing.T) {
	xmlData := `<annotationSet xmlns="urn:hl7-org:v3" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
  <activityTime value="20250923103600"/>
  <component>
    <annotation>
      <code code="MDC_ECG_INTERPRETATION" codeSystem="2.16.840.1.113883.6.24"/>
      <value xsi:type="ST">Rythme sinusal avec ESA</value>
    </annotation>
  </component>
</annotationSet>`

	var annSet AnnotationSet
	err := xml.Unmarshal([]byte(xmlData), &annSet)
	require.NoError(t, err)

	// Verify annotation
	require.Len(t, annSet.Component, 1)
	ann := &annSet.Component[0].Annotation

	// Verify code
	require.NotNil(t, ann.Code)
	assert.Equal(t, "MDC_ECG_INTERPRETATION", ann.Code.Code)

	// Verify value is ST type
	require.NotNil(t, ann.Value)
	assert.True(t, ann.Value.IsST(), "Value should be StringValue")
	assert.False(t, ann.Value.IsPQ(), "Value should not be PhysicalQuantity")
	assert.Equal(t, "ST", ann.Value.XsiType)

	// Verify text content
	text, ok := ann.Value.GetText()
	assert.True(t, ok)
	assert.Equal(t, "Rythme sinusal avec ESA", text)

	// GetValueFloat should return false for ST
	_, ok = ann.Value.GetValueFloat()
	assert.False(t, ok, "GetValueFloat should return false for ST type")
}

// TestTextAnnotation_NestedStatements tests nested text annotations
func TestTextAnnotation_NestedStatements(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}

	// Add interpretation annotation with nested statements
	interpIdx := annSet.AddTextAnnotation("MDC_ECG_INTERPRETATION", "2.16.840.1.113883.6.24", "")

	// Add nested interpretation statements
	interp := annSet.GetAnnotation(interpIdx)
	require.NotNil(t, interp)

	interp.AddNestedTextAnnotation("MDC_ECG_INTERPRETATION_STATEMENT", "2.16.840.1.113883.6.24", "Rythme sinusal avec ESA")
	interp.AddNestedTextAnnotation("MDC_ECG_INTERPRETATION_STATEMENT", "2.16.840.1.113883.6.24", "--- Interprétation sans connaître le sexe/l'âge du patient ---")

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(annSet, "", "  ")
	require.NoError(t, err)

	// Unmarshal to verify round-trip
	var unmarshaledSet AnnotationSet
	err = xml.Unmarshal(xmlData, &unmarshaledSet)
	require.NoError(t, err)

	// Verify parent annotation
	interpResult := unmarshaledSet.GetAnnotationByCode("MDC_ECG_INTERPRETATION")
	require.NotNil(t, interpResult)
	require.Len(t, interpResult.Component, 2)

	// Verify first nested statement
	stmt1 := &interpResult.Component[0].Annotation
	require.NotNil(t, stmt1.Value)
	assert.True(t, stmt1.Value.IsST())
	text1, ok := stmt1.Value.GetText()
	assert.True(t, ok)
	assert.Equal(t, "Rythme sinusal avec ESA", text1)

	// Verify second nested statement
	stmt2 := &interpResult.Component[1].Annotation
	require.NotNil(t, stmt2.Value)
	assert.True(t, stmt2.Value.IsST())
	text2, ok := stmt2.Value.GetText()
	assert.True(t, ok)
	assert.Equal(t, "--- Interprétation sans connaître le sexe/l'âge du patient ---", text2)
}

// TestTextAnnotation_MixedTypes tests mixing PQ and ST annotations
func TestTextAnnotation_MixedTypes(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}

	// Add numeric annotation (PQ)
	annSet.AddAnnotation("MDC_ECG_HEART_RATE", "2.16.840.1.113883.6.24", 57, "bpm")

	// Add text annotation (ST)
	annSet.AddTextAnnotation("MDC_ECG_INTERPRETATION", "2.16.840.1.113883.6.24", "Rythme sinusal normal")

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(annSet, "", "  ")
	require.NoError(t, err)

	// Unmarshal to verify round-trip
	var unmarshaledSet AnnotationSet
	err = xml.Unmarshal(xmlData, &unmarshaledSet)
	require.NoError(t, err)

	require.Len(t, unmarshaledSet.Component, 2)

	// Verify first annotation is PQ
	ann1 := &unmarshaledSet.Component[0].Annotation
	require.NotNil(t, ann1.Value)
	assert.True(t, ann1.Value.IsPQ())
	val, ok := ann1.Value.GetValueFloat()
	assert.True(t, ok)
	assert.Equal(t, float64(57), val)
	assert.Equal(t, "bpm", ann1.Value.GetValueUnit())

	// Verify second annotation is ST
	ann2 := &unmarshaledSet.Component[1].Annotation
	require.NotNil(t, ann2.Value)
	assert.True(t, ann2.Value.IsST())
	text, ok := ann2.Value.GetText()
	assert.True(t, ok)
	assert.Equal(t, "Rythme sinusal normal", text)
}

// TestTextAnnotation_VendorSpecific tests vendor-specific text annotations
func TestTextAnnotation_VendorSpecific(t *testing.T) {
	annSet := &AnnotationSet{
		ActivityTime: &Time{Value: "20250923103600"},
	}

	// Add vendor-specific text annotation
	annSet.AddTextAnnotationWithCodeSystemName("MINDRAY_INTERPRETATION", "MINDRAY", "Tracé normal pour l'âge")

	// Marshal and unmarshal
	xmlData, err := xml.MarshalIndent(annSet, "", "  ")
	require.NoError(t, err)

	var unmarshaledSet AnnotationSet
	err = xml.Unmarshal(xmlData, &unmarshaledSet)
	require.NoError(t, err)

	// Verify annotation
	ann := unmarshaledSet.GetAnnotationByCode("MINDRAY_INTERPRETATION")
	require.NotNil(t, ann)
	assert.Equal(t, "MINDRAY", ann.Code.CodeSystemName)
	assert.Equal(t, "", ann.Code.CodeSystem) // Empty for vendor codes

	// Verify text value
	require.NotNil(t, ann.Value)
	assert.True(t, ann.Value.IsST())
	text, ok := ann.Value.GetText()
	assert.True(t, ok)
	assert.Equal(t, "Tracé normal pour l'âge", text)
}
