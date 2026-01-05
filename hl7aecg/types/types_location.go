package types

// Location represents the trial site location where ECG waveforms were acquired.
// It identifies the physical location and facility where the subject's data was collected.
//
// XML Structure:
//
//	<location>
//	  <trialSite>
//	    <id root="2.16.840.1.113883.3.5" extension="SITE_1"/>
//	    <location>
//	      <name>1st Clinic of Milwaukee</name>
//	      <addr>
//	        <city>Milwaukee</city>
//	        <state>WI</state>
//	        <country>USA</country>
//	      </addr>
//	    </location>
//	  </trialSite>
//	</location>
//
// Cardinality: Optional (in ClinicalTrial context)
// Reference: HL7 aECG Implementation Guide, Page 21
type Location struct {
	// TrialSite contains identification and details about the clinical trial site.
	//
	// XML Tag: <trialSite>...</trialSite>
	// Cardinality: Required (within Location)
	TrialSite TrialSite `xml:"trialSite"`
}

// TrialSite represents a specific location where clinical trial activities take place.
// This is the physical site where ECG waveforms were acquired from the subject.
//
// Each trial site should have a unique identifier to distinguish it from other
// sites participating in the same trial.
//
// XML Structure:
//
//	<trialSite>
//	  <id root="2.16.840.1.113883.3.5" extension="SITE_1"/>
//	  <location>...</location>
//	</trialSite>
//
// Cardinality: Required (within Location)
type TrialSite struct {
	// ID is the unique identifier for the trial site.
	//
	// Structure:
	//   - Root: Required, must be a UID (OID or UUID)
	//   - Extension: Optional, traditional site identifier
	//
	// Best Practice:
	//   The sponsor (or its vendor) should assign a unique OID to every trial site.
	//   This OID goes into the root part. The traditional site identifier
	//   (e.g., "SITE_1", "NYC_CARDIO") goes into the extension.
	//
	// The combination of root and extension must be universally unique
	// across all trial sites globally.
	//
	// Example: root="2.16.840.1.113883.3.5" extension="SITE_1"
	//
	// XML Tag: <id root="..." extension="..."/>
	// Cardinality: Required
	ID ID `xml:"id"`

	// Location contains detailed information about the physical location
	// of the trial site, including name and address.
	//
	// XML Tag: <location>...</location>
	// Cardinality: Optional
	Location *SiteLocation `xml:"location,omitempty"`

	// ResponsibleParty contains information about the trial investigator
	// responsible for ECG waveform acquisition at this site.
	//
	// XML Tag: <responsibleParty>...</responsibleParty>
	// Cardinality: Optional
	ResponsibleParty *ResponsibleParty `xml:"responsibleParty,omitempty"`
}

// SiteLocation represents the physical location details of a trial site.
// It includes the name and address of the facility where data was collected.
//
// XML Structure:
//
//	<location>
//	  <name>1st Clinic of Milwaukee</name>
//	  <addr>...</addr>
//	</location>
//
// Cardinality: Optional (within TrialSite)
type SiteLocation struct {
	// Name is the human-readable name of the trial site or facility.
	//
	// This is typically the official name of the clinic, hospital,
	// or research center conducting the trial.
	//
	// Example: "1st Clinic of Milwaukee", "Johns Hopkins Hospital"
	//
	// XML Tag: <name>...</name>
	// Cardinality: Optional
	Name *string `xml:"name,omitempty"`

	// Addr contains the physical address of the trial site.
	//
	// XML Tag: <addr>...</addr>
	// Cardinality: Optional
	Addr *Address `xml:"addr,omitempty"`
}

// Address represents a physical address.
// Contains standard address components for locating a physical site.
//
// Note: The HL7 aECG schema provides additional address fields beyond
// those listed here (street, postal code, etc.). Refer to the XML schema
// for a complete listing of available address fields.
//
// XML Structure:
//
//	<addr>
//	  <city>Milwaukee</city>
//	  <state>WI</state>
//	  <country>USA</country>
//	</addr>
//
// Cardinality: Optional (within SiteLocation)
type Address struct {
	// City is the city or municipality name.
	//
	// Example: "Milwaukee", "Boston", "London"
	//
	// XML Tag: <city>...</city>
	// Cardinality: Optional
	City *string `xml:"city,omitempty"`

	// State is the state, province, or region name.
	//
	// Use standard abbreviations where applicable (e.g., "WI" for Wisconsin).
	//
	// Example: "WI", "CA", "Ontario"
	//
	// XML Tag: <state>...</state>
	// Cardinality: Optional
	State *string `xml:"state,omitempty"`

	// Country is the country name or code.
	//
	// Example: "USA", "Canada", "United Kingdom"
	//
	// XML Tag: <country>...</country>
	// Cardinality: Optional
	Country *string `xml:"country,omitempty"`
}
