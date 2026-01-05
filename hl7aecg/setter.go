package hl7aecg

import (
	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
)

// SetText sets the Text field of the HL7AEcg instance.
func (h *Hl7xml) SetText(text string) *Hl7xml {
	h.HL7AEcg.Text = text
	return h
}

// SetEffectiveTime sets the EffectiveTime field of the HL7AEcg instance.
func (h *Hl7xml) SetEffectiveTime(low, high string) *Hl7xml {
	h.HL7AEcg.EffectiveTime = types.NewEffectiveTime(low, high)
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
//   h.SetLocation("SITE_001", "2.16.840.1.113883.3.5", "Boston Medical Center", "Boston", "MA", "USA")
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
				ID: types.ID{
					Root:      siteRoot,
					Extension: siteID,
				},
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

	return h
}
