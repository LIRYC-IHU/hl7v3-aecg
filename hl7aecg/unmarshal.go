package hl7aecg

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
)

// Unmarshal parses aECG XML data into the Hl7xml struct.
// The XML data must be a valid HL7 aECG document.
func (h *Hl7xml) Unmarshal(data []byte) error {
	// Reset the HL7AEcg field before unmarshalling
	h.HL7AEcg = types.HL7AEcg{}
	if err := xml.Unmarshal(data, &h.HL7AEcg); err != nil {
		return fmt.Errorf("unmarshal AnnotatedECG: %w", err)
	}
	return nil
}

// UnmarshalFromReader parses aECG XML from an io.Reader.
// Useful for streaming data from network or other sources.
func (h *Hl7xml) UnmarshalFromReader(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("read XML data: %w", err)
	}
	return h.Unmarshal(data)
}

// UnmarshalFromFile parses aECG XML from a file path.
// This is a convenience method that reads the file and calls Unmarshal.
func (h *Hl7xml) UnmarshalFromFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file %s: %w", filePath, err)
	}
	return h.Unmarshal(data)
}

// UnmarshalAndValidate parses aECG XML and validates the resulting document.
// This combines Unmarshal() and Validate() into a single operation.
// Returns parsing error if XML is invalid, validation error if document is invalid.
func (h *Hl7xml) UnmarshalAndValidate(data []byte) error {
	if err := h.Unmarshal(data); err != nil {
		return err
	}
	return h.Validate()
}
