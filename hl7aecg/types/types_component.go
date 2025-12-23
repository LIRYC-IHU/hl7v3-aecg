package types

// Component represents the component parts of the aECG.
//
// These are the waveforms and annotations. Even though the XML schema says
// this is optional, the message is not useful without at least one series.
//
// XML Structure:
//
//	<component>
//	  <series>...</series>
//	</component>
//
// Cardinality: Optional but strongly recommended
// Reference: HL7 aECG Implementation Guide, Page 26
type Component struct {
	// Series contains all sequences, regions of interest, and annotations
	// sharing a common frame of reference.
	//
	// XML Tag: <series>...</series>
	// Cardinality: Required (within Component)
	Series Series `xml:"series"`
}
