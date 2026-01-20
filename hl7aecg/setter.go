package hl7aecg

import (
	"fmt"

	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
)

// SetText sets the Text field of the HL7AEcg instance.
func (h *Hl7xml) SetText(text string) *Hl7xml {
	h.HL7AEcg.Text = text
	return h
}

// SetEffectiveTime sets the EffectiveTime field of the HL7AEcg instance.
func (h *Hl7xml) SetEffectiveTime(low, high string, low_inclusive, high_inclusive *bool) *Hl7xml {
	h.HL7AEcg.EffectiveTime = types.NewEffectiveTime(low, high, low_inclusive, high_inclusive)
	return h
}

// SetSchemaLocation sets the SchemaLocation attribute of the HL7AEcg instance.
// Init by default sets SchemaLocation to "urn:hl7-org:v3".
func (h *Hl7xml) SetSchemaLocation(location string) *Hl7xml {
	h.HL7AEcg.SchemaLocation = location
	return h
}

// SetType sets the Type attribute of the HL7AEcg instance.
// Init by default sets Type to "Observation".
func (h *Hl7xml) SetType(typeValue string) *Hl7xml {
	h.HL7AEcg.Type = typeValue
	return h
}

// SetLocation configures the trial site location for the clinical trial.
//
// This method uses the ComponentOf structure to access ClinicalTrial.
//
// Parameters:
//   - siteID: Site identifier (e.g., "SITE_001")
//   - siteRoot: OID for the site (e.g., "2.16.840.1.113883.3.5")
//   - siteName: Name of the trial site (e.g., "1st Clinic of Milwaukee")
//   - city: City name (optional, use "" to skip)
//   - state: State/province (optional, use "" to skip)
//   - country: Country name (optional, use "" to skip)
//
// Example:
//
//	h.SetLocation("SITE_001", "2.16.840.1.113883.3.5", "Boston Medical Center", "Boston", "MA", "USA")
func (h *Hl7xml) SetLocation(siteID, siteRoot, siteName, city, state, country string) *Hl7xml {
	// Ensure ComponentOf structure exists (should be initialized by SetSubject)
	if h.HL7AEcg.ComponentOf == nil {
		h.SetSubject("", "", "") // Initialize structure
	}

	// Get reference to ClinicalTrial in ComponentOf structure
	clinicalTrial := &h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.SubjectAssignment.ComponentOf.ClinicalTrial

	// Initialize Location if nil
	if clinicalTrial.Location == nil {
		clinicalTrial.Location = &types.Location{
			TrialSite: types.TrialSite{
				ID:       types.ID{},
				Location: &types.SiteLocation{},
			},
		}
	}

	// Set site name if provided
	if siteName != "" {
		clinicalTrial.Location.TrialSite.Location.SetName(siteName)
	}

	// Set address if any component is provided
	if city != "" || state != "" || country != "" {
		clinicalTrial.Location.TrialSite.Location.SetFullAddress(city, state, country)
	}
	clinicalTrial.Location.TrialSite.ID.SetID(siteRoot, siteID)

	return h
}

// SetResponsibleParty configures the trial investigator responsible for ECG acquisition.
//
// This method uses the ComponentOf structure to access ClinicalTrial > Location > TrialSite.
//
// Parameters:
//   - investigatorRoot: OID for the investigator (use "" to use global root ID from singleton)
//   - investigatorID: Investigator identifier extension (e.g., "INV_001", "trialInvestigator")
//   - prefix: Title/honorific (e.g., "Dr.", "Prof.") - use "" to skip
//   - given: Given/first name - use "" to skip
//   - family: Family/last name - use "" to skip
//   - suffix: Suffix/credentials (e.g., "MD", "PhD") - use "" to skip
//
// Example:
//
//	h.SetResponsibleParty("", "INV_001", "Dr.", "John", "Smith", "MD")
//
// To set an empty name element (like <name/>), pass empty strings for all name components:
//
//	h.SetResponsibleParty("", "trialInvestigator", "", "", "", "")
func (h *Hl7xml) SetResponsibleParty(investigatorRoot, investigatorID, prefix, given, family, suffix string) *Hl7xml {
	// Ensure ComponentOf structure exists (should be initialized by SetSubject)
	if h.HL7AEcg.ComponentOf == nil {
		h.SetSubject("", "", "") // Initialize structure
	}

	// Get reference to ClinicalTrial in ComponentOf structure
	clinicalTrial := &h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.SubjectAssignment.ComponentOf.ClinicalTrial

	// Initialize Location if nil
	if clinicalTrial.Location == nil {
		clinicalTrial.Location = &types.Location{
			TrialSite: types.TrialSite{
				ID: types.ID{},
			},
		}
	}

	// Initialize ResponsibleParty if nil
	if clinicalTrial.Location.TrialSite.ResponsibleParty == nil {
		clinicalTrial.Location.TrialSite.ResponsibleParty = &types.ResponsibleParty{
			TrialInvestigator: types.TrialInvestigator{
				ID: types.ID{},
			},
		}
	}

	rp := clinicalTrial.Location.TrialSite.ResponsibleParty

	// Set investigator ID
	rp.SetInvestigatorID(investigatorRoot, investigatorID)

	// Set investigator name if any component is provided
	if prefix != "" || given != "" || family != "" || suffix != "" {
		rp.SetInvestigatorName(prefix, given, family, suffix)
	} else {
		// Set empty name element
		rp.SetEmptyInvestigatorName()
	}

	return h
}

// SetSeriesCode updates the code of the most recently added series with additional attributes.
//
// This method automatically finds the last series in the Component array and updates
// its code with the provided codeSystemName and displayName attributes.
//
// Parameters:
//   - code: The code value (e.g., RHYTHM_CODE, REPRESENTATIVE_BEAT_CODE)
//   - codeSystem: The OID of the code system (e.g., HL7_ActCode_OID)
//   - codeSystemName: Human-readable name of the code system (e.g., "ActCode")
//   - displayName: Human-readable name of the code (e.g., "Rhythm Waveforms")
//
// Returns an error if no series exists in the Component array.
//
// Example:
//
//	err := h.SetSeriesCode(
//	    types.RHYTHM_CODE,
//	    types.HL7_ActCode_OID,
//	    "ActCode",
//	    "Rhythm Waveforms",
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// This allows calling h.SetSeriesCode(...) instead of having to access
// h.HL7AEcg.Component[index].Series.SetSeriesCode(...) manually.
func (h *Hl7xml) SetSeriesCode(code types.SeriesTypeCode, codeSystem types.CodeSystemOID, codeSystemName, displayName string) error {
	// Check if there are any components
	if len(h.HL7AEcg.Component) == 0 {
		return fmt.Errorf("no series found: Component array is empty")
	}

	// Get the last component
	lastIdx := len(h.HL7AEcg.Component) - 1
	lastComponent := &h.HL7AEcg.Component[lastIdx]

	// Update the series code
	lastComponent.Series.SetSeriesCode(code, codeSystem, codeSystemName, displayName)

	return nil
}

// SetDerivedSeriesCode updates the series code for the most recently added derived series
// in the most recent parent series.
//
// Example:
//
//	h.AddDerivedSeries(...).
//	  SetDerivedSeriesCode(types.REPRESENTATIVE_BEAT_CODE, types.HL7_ActCode_OID, "ActCode", "Representative Beat")
func (h *Hl7xml) SetDerivedSeriesCode(
	code types.SeriesTypeCode,
	codeSystem types.CodeSystemOID,
	codeSystemName, displayName string,
) error {
	if len(h.HL7AEcg.Component) == 0 {
		return fmt.Errorf("no parent series found")
	}

	lastSeries := &h.HL7AEcg.Component[len(h.HL7AEcg.Component)-1].Series
	if len(lastSeries.Derivation) == 0 {
		return fmt.Errorf("no derived series found in parent series")
	}

	// Update most recent derived series
	lastDerived := &lastSeries.Derivation[len(lastSeries.Derivation)-1].DerivedSeries
	lastDerived.Code.SetCode(code, codeSystem, codeSystemName, displayName)
	return nil
}

// SetDerivedSeriesID sets the ID for the most recently added derived series.
//
// Parameters:
//   - root: OID root (use "" to inherit from parent document)
//   - extension: Extension identifier (e.g., "derivedSeries", "representativeBeat1")
//
// Example:
//
//	h.SetDerivedSeriesID("", "derivedSeries")
func (h *Hl7xml) SetDerivedSeriesID(root, extension string) error {
	if len(h.HL7AEcg.Component) == 0 {
		return fmt.Errorf("no parent series found")
	}

	lastSeries := &h.HL7AEcg.Component[len(h.HL7AEcg.Component)-1].Series
	if len(lastSeries.Derivation) == 0 {
		return fmt.Errorf("no derived series found in parent series")
	}

	// Update most recent derived series
	lastDerived := &lastSeries.Derivation[len(lastSeries.Derivation)-1].DerivedSeries
	if lastDerived.ID == nil {
		lastDerived.ID = &types.ID{}
	}
	lastDerived.ID.SetID(root, extension)
	return nil
}
