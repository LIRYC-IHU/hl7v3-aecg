package types

import (
	"encoding/xml"
)

func (h *HL7AEcg) Unmarshal(data []byte) error {
	return xml.Unmarshal(data, h)
}
