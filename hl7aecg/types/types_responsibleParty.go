package types

// ResponsibleParty represents the trial investigator responsible for the acquisition
// of ECG waveforms at the trial site.
//
// The responsible party is typically the principal investigator or lead physician
// at the trial site who oversees the ECG data collection process.
//
// XML Structure:
//
//	<responsibleParty>
//	  <trialInvestigator>
//	    <id root="2.16.840.1.113883.3.6" extension="INV_001"/>
//	    <investigatorPerson>
//	      <name>
//	        <given>John</given>
//	        <family>Smith</family>
//	        <prefix>Dr.</prefix>
//	        <suffix>MD</suffix>
//	      </name>
//	    </investigatorPerson>
//	  </trialInvestigator>
//	</responsibleParty>
//
// Cardinality: Optional (in ClinicalTrial context)
// Reference: HL7 aECG Implementation Guide, Page 13
type ResponsibleParty struct {
	// TrialInvestigator contains identification and details about the
	// investigator responsible for ECG waveform acquisition.
	//
	// XML Tag: <trialInvestigator>...</trialInvestigator>
	// Cardinality: Required (within ResponsibleParty)
	TrialInvestigator TrialInvestigator `xml:"trialInvestigator"`
}

// TrialInvestigator represents the investigator responsible for acquiring
// ECG waveforms at a trial site.
//
// This is typically the principal investigator or lead physician overseeing
// the clinical trial activities at the specific site.
//
// XML Structure:
//
//	<trialInvestigator>
//	  <id root="2.16.840.1.113883.3.6" extension="INV_001"/>
//	  <investigatorPerson>...</investigatorPerson>
//	</trialInvestigator>
//
// Cardinality: Required (within ResponsibleParty)
type TrialInvestigator struct {
	// ID is the unique identifier for the investigator.
	//
	// Structure:
	//   - Root: Required, must be a UID (OID or UUID)
	//   - Extension: Optional, traditional investigator identifier
	//
	// Best Practice:
	//   The sponsor (or its vendor) should assign a unique OID to every investigator.
	//   This OID goes into the root part. The traditional investigator identifier
	//   (e.g., "INV_001", "PI_SMITH") goes into the extension.
	//
	// The combination of root and extension must be universally unique
	// across all investigators globally.
	//
	// Example: root="2.16.840.1.113883.3.6" extension="INV_001"
	//
	// XML Tag: <id root="..." extension="..."/>
	// Cardinality: Required
	ID ID `xml:"id"`

	// InvestigatorPerson contains personal information about the investigator,
	// including their name.
	//
	// XML Tag: <investigatorPerson>...</investigatorPerson>
	// Cardinality: Optional
	InvestigatorPerson *InvestigatorPerson `xml:"investigatorPerson,omitempty"`
}

// InvestigatorPerson represents personal information about a trial investigator.
//
// XML Structure:
//
//	<investigatorPerson>
//	  <name>
//	    <given>John</given>
//	    <family>Smith</family>
//	    <prefix>Dr.</prefix>
//	    <suffix>MD</suffix>
//	  </name>
//	</investigatorPerson>
//
// Cardinality: Optional (within TrialInvestigator)
type InvestigatorPerson struct {
	// Name contains the investigator's name components.
	//
	// Note: The HL7 aECG schema provides multiple fields for specifying
	// a person's name (prefix, given, family, suffix, etc.).
	// Refer to the XML schema for a complete listing of available name fields.
	//
	// XML Tag: <name>...</name>
	// Cardinality: Optional
	Name *PersonName `xml:"name,omitempty"`
}

// PersonName represents a person's name with its various components.
//
// This structure follows HL7 naming conventions, allowing for proper
// representation of names across different cultures and formats.
//
// XML Structure:
//
//	<name>
//	  <prefix>Dr.</prefix>
//	  <given>John</given>
//	  <family>Smith</family>
//	  <suffix>MD</suffix>
//	</name>
//
// Cardinality: Optional (within InvestigatorPerson)
type PersonName struct {
	// Prefix represents titles or honorifics that precede the name.
	//
	// Examples: "Dr.", "Prof.", "Mr.", "Ms."
	//
	// XML Tag: <prefix>...</prefix>
	// Cardinality: Optional
	Prefix *string `xml:"prefix,omitempty"`

	// Given represents the given name(s) or first name(s).
	//
	// Examples: "John", "Mary Jane"
	//
	// XML Tag: <given>...</given>
	// Cardinality: Optional
	Given *string `xml:"given,omitempty"`

	// Family represents the family name or surname.
	//
	// Examples: "Smith", "Johnson", "van der Berg"
	//
	// XML Tag: <family>...</family>
	// Cardinality: Optional
	Family *string `xml:"family,omitempty"`

	// Suffix represents suffixes or credentials that follow the name.
	//
	// Examples: "MD", "PhD", "Jr.", "III"
	//
	// XML Tag: <suffix>...</suffix>
	// Cardinality: Optional
	Suffix *string `xml:"suffix,omitempty"`
}
