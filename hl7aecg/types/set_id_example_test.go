package types_test

import (
	"fmt"

	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
)

// Example demonstrating automatic default extensions for IDs
func ExampleID_SetID_automaticExtension() {
	// Setup global root ID
	aecg := &types.HL7AEcg{}
	aecg.SetRootID("755.3045256.2025923.103550", "")

	// Example 1: Using wrapper method on ClinicalTrial
	// The extension will automatically be "clinicalTrial" if not provided
	ct := &types.ClinicalTrial{ID: types.ID{}}
	ct.SetID("", "") // Uses singleton root + default extension "clinicalTrial"
	fmt.Printf("ClinicalTrial ID: root=%s, extension=%s\n", ct.ID.Root, ct.ID.Extension)

	// Example 2: Using direct ID.SetID with explicit default
	id := &types.ID{}
	id.SetID("", "", "customDefault") // Uses singleton root + "customDefault"
	fmt.Printf("Custom ID: root=%s, extension=%s\n", id.Root, id.Extension)

	// Example 3: Override default extension
	ts := &types.TrialSite{ID: types.ID{}}
	ts.SetID("", "customSite") // Uses singleton root + "customSite" (overrides default)
	fmt.Printf("TrialSite ID: root=%s, extension=%s\n", ts.ID.Root, ts.ID.Extension)

	// Output:
	// ClinicalTrial ID: root=755.3045256.2025923.103550, extension=clinicalTrial
	// Custom ID: root=755.3045256.2025923.103550, extension=customDefault
	// TrialSite ID: root=755.3045256.2025923.103550, extension=customSite
}

// Example showing all available wrapper methods with automatic extensions
func ExampleID_SetID_wrapperMethods() {
	// Setup global root ID
	aecg := &types.HL7AEcg{}
	aecg.SetRootID("2.16.840.1.113883.3.1", "")

	// Each structure has a wrapper method that automatically applies the correct extension
	examples := []struct {
		name      string
		extension string
	}{
		{"HL7AEcg", "annotatedEcg"},
		{"ClinicalTrial", "clinicalTrial"},
		{"TrialSubject", "trialSubject"},
		{"TrialSite", "trialSite"},
		{"Series", "series"},
		{"TrialInvestigator", "trialInvestigator"},
	}

	for _, ex := range examples {
		fmt.Printf("%s uses default extension: %s\n", ex.name, ex.extension)
	}

	// Output:
	// HL7AEcg uses default extension: annotatedEcg
	// ClinicalTrial uses default extension: clinicalTrial
	// TrialSubject uses default extension: trialSubject
	// TrialSite uses default extension: trialSite
	// Series uses default extension: series
	// TrialInvestigator uses default extension: trialInvestigator
}
