package types

// ClinicalTrialSponsor represents the author of a clinical trial.
// It identifies the sponsoring organization of the clinical trial.
//
// XML Structure:
//
//	<author>
//	  <clinicalTrialSponsor>
//	    <sponsorOrganization>
//	      <id root="2.16.840.1.113883.3"/>
//	      <name>ABC Drug Company</name>
//	    </sponsorOrganization>
//	  </clinicalTrialSponsor>
//	</author>
//
// Cardinality: Optional (in ClinicalTrial context)
// Reference: HL7 aECG Implementation Guide, Page 11
type ClinicalTrialSponsor struct {
	// SponsorOrganization contains the identification and name of the
	// organization sponsoring the clinical trial.
	//
	// XML Tag: <sponsorOrganization>...</sponsorOrganization>
	// Cardinality: Required (within ClinicalTrialSponsor)
	SponsorOrganization SponsorOrganization `xml:"sponsorOrganization"`
}

// SponsorOrganization represents the organization that sponsors a clinical trial.
//
// The sponsor is responsible for initiating, managing, and financing the clinical trial.
// This organization is considered the "author" of the trial.
//
// XML Structure:
//
//	<sponsorOrganization>
//	  <id root="2.16.840.1.113883.3" extension="ABC-PHARMA"/>
//	  <name>ABC Drug Company</name>
//	</sponsorOrganization>
//
// Cardinality: Required (within ClinicalTrialSponsor)
type SponsorOrganization struct {
	// ID is the unique identifier for the sponsoring organization.
	//
	// Structure:
	//   - Root: Required, must be a UID (OID or UUID)
	//   - Extension: Optional, may be any identifier
	//
	// Best Practice:
	//   The sponsor should obtain an OID "root" for its organization.
	//   This OID "root" will be used as a prefix for all other OIDs the
	//   organization creates (trials, protocols, subjects, etc.).
	//
	//   The OID "root" goes into the root part of the ID.
	//   If the sponsor has another more common and well-known identifier
	//   (e.g., a stock ticker, registration number), that identifier
	//   should go into the extension.
	//
	// The combination of root and extension must be universally unique.
	//
	// Example: root="2.16.840.1.113883.3" extension="ABC-PHARMA"
	//
	// XML Tag: <id root="..." extension="..."/>
	// Cardinality: Optional
	ID ID `xml:"id,omitempty"`

	// Extension provides an additional identifier for the organization.
	//
	// Note: This field appears to be redundant with ID.Extension.
	// Consider removing this field or clarifying its purpose.
	//
	// XML Tag: <extension>...</extension>
	// Cardinality: Optional
	Extension *string `xml:"extension,omitempty"`

	// Name is the human-readable name of the sponsoring organization.
	//
	// This is the official or commonly used name of the organization
	// that is sponsoring the clinical trial.
	//
	// Example: "ABC Drug Company", "XYZ Pharmaceuticals Inc."
	//
	// XML Tag: <name>...</name>
	// Cardinality: Optional
	Name *string `xml:"name,omitempty"`
}
