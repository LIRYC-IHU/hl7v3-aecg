package hl7aecg

import (
	"context"
	stdxml "encoding/xml"
	"os"

	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
)

type Hl7xml struct {
	HL7AEcg   types.HL7AEcg
	ctx       context.Context
	outputDir string
	vctx      *types.ValidationContext
}

func NewHl7xml(outputDir string) *Hl7xml {
	return &Hl7xml{
		HL7AEcg: types.HL7AEcg{
			XmlnsXsi:            "http://www.w3.org/2001/XMLSchema-instance",
			XmlnsVoc:            "urn:hl7-org:v3/voc",
			Xmlns:               "urn:hl7-org:v3",
			ID:                  &types.ID{},
			Code:                &types.Code[types.CPT_CODE, types.CodeSystemOID]{},
			ConfidentialityCode: &types.Code[types.ConfidentialityCode, string]{},
			ReasonCode:          &types.Code[types.ReasonCode, string]{},
		},
		ctx:       context.Background(),
		outputDir: outputDir,
		vctx:      types.NewValidationContext(true),
	}
}

// Initialize sets up the HL7 aECG document with the provided CPT code and code system OID.
func (h *Hl7xml) Initialize(code types.CPT_CODE, codeSystem types.CodeSystemOID, codeSystemName, displayName string) *Hl7xml {
	h.HL7AEcg.Code.SetCode(code, codeSystem, codeSystemName, displayName)

	// Set required HL7 v3 RIM attributes for schema validation
	h.HL7AEcg.Type = "Observation"
	h.HL7AEcg.SchemaLocation = "urn:hl7-org:v3"

	return h
}

// String renders the document as XML.
//
// It marshals with encoding/xml, and must keep doing so: several types in this
// package carry custom MarshalXML methods that the document depends on for
// correctness. SequenceCode flattens <code><Lead code="…"/></code> into the
// <code code="…"/> form the FDA aECG Implementation Guide specifies, and
// SequenceValue writes the typed <value> body — the origin, scale and digits of
// every waveform.
//
// This previously used a third-party marshaller for its short-form empty tags.
// That marshaller does not invoke those methods, so String() emitted lead codes
// in a nested form this package's own UnmarshalXML cannot read, and dropped the
// contents of every <value> element: a document that round-tripped to zero
// leads and no waveform data. Short-form empty tags are cosmetic; reading back
// what you wrote is not.
func (h *Hl7xml) String() (string, error) {
	data, err := stdxml.MarshalIndent(h.HL7AEcg, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (h *Hl7xml) Test() (*Hl7xml, error) {
	dir := "/tmp/hl7aecg_example.xml"

	data, err := stdxml.MarshalIndent(h.HL7AEcg, "", "  ")
	if err != nil {
		return h, err
	}
	file, err := os.Create(dir)
	if err != nil {
		return h, err
	}
	defer file.Close()
	file.Write(data)
	return h, err
}
