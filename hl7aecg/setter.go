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
