package hl7aecg

import (
	"context"
	"fmt"
	"os"

	"github.com/ECUST-XX/xml"
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
func (h *Hl7xml) Initialize(code types.CPT_CODE, codeSystem types.CodeSystemOID, codeSystemName string) *Hl7xml {
	fmt.Println("Initialize HL7 aECG document")
	h.HL7AEcg.Code.SetCode(code, codeSystem, codeSystemName)
	return h
}

func (h *Hl7xml) String() (string, error) {
	data, err := xml.MarshalIndentShortForm(h.HL7AEcg, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (h *Hl7xml) Test() (*Hl7xml, error) {
	dir := "/tmp/hl7aecg_example.xml"
	data, err := xml.MarshalIndentShortForm(h.HL7AEcg, "", "  ")
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
