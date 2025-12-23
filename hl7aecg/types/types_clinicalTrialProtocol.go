package types

// ClinicalTrialProtocol represents the clinicalTrialProtocol element in an HL7 aECG document.
// It contains information about the clinical trial protocol used to define the trial.
//
// The protocol defines the procedures, objectives, design, and organization of a clinical trial.
// It includes codes used to name treatment groups, timepoint events, reference events, and relative timepoints.
//
// XML Structure:
//
//	<clinicalTrialProtocol>
//	  <id root="2.16.840.1.113883.3.2" extension="PUK-123-PROT-A"/>
//	  <title>Cardiac Safety Protocol for Compound PUK-123</title>
//	</clinicalTrialProtocol>
//
// Cardinality: Optional
// Reference: HL7 aECG Implementation Guide, Page 10
type ClinicalTrialProtocol struct {
	// ID is the unique identifier for the protocol used to define the trial.
	//
	// Structure:
	//   - Root: Required, must be a UID (OID or UUID)
	//   - Extension: Optional, traditional protocol identifier
	//
	// Best Practice:
	//   The sponsor should assign a unique OID to every new protocol.
	//   This OID goes into the root part. The traditional protocol identifier
	//   goes into the extension (e.g., "PUK-123-PROT-A").
	//
	// Important: If the protocol defines codes for treatment groups, timepoint events,
	// reference events, and relative timepoints, this OID will be used as the
	// coding system identifier for those codes.
	//
	// The combination of root and extension must be universally unique.
	//
	// XML Tag: <id root="..." extension="..."/>
	// Cardinality: Required
	ID ID `xml:"id"`

	// Title is the human-readable name of the protocol.
	//
	// Provides a descriptive name for the protocol that helps identify
	// its purpose and scope.
	//
	// Example: "Cardiac Safety Protocol for Compound PUK-123"
	//
	// XML Tag: <title>...</title>
	// Cardinality: Optional
	Title *string `xml:"title,omitempty"`
}
