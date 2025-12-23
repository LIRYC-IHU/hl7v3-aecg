package hl7aecg

import (
	"strconv"

	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
	"github.com/google/uuid"
)

// =============================================================================
// Builder Methods for ECG Components
// =============================================================================

// AddRhythmSeries creates and adds a rhythm series to the aECG.
//
// Parameters:
//   - startTime: Start timestamp (YYYYMMDDHHmmss.SSS format)
//   - endTime: End timestamp
//   - sampleRate: Sampling rate in Hz (e.g., 500 for 500 Hz)
//   - leads: Map of lead codes to their raw sample data
//   - origin: Baseline voltage value (typically 0)
//   - scale: Voltage resolution per digit (e.g., 5 for 5 µV)
//
// Example:
//
//	leads := map[string][]int{
//	    types.MDC_ECG_LEAD_I:  []int{1, 2, 3, 4, 5, ...},
//	    types.MDC_ECG_LEAD_II: []int{10, 11, 12, 13, 14, ...},
//	}
//	h.AddRhythmSeries("20021122091000.000", "20021122091010.000", 500, leads, 0, 5)
func (h *Hl7xml) AddRhythmSeries(
	startTime, endTime string,
	sampleRate float64,
	leads map[types.LeadCode][]int,
	origin, scale float64,
) *Hl7xml {
	series := h.buildSeries(
		types.RHYTHM_CODE,
		startTime,
		endTime,
		sampleRate,
		leads,
		origin,
		scale,
	)

	h.HL7AEcg.Component = append(h.HL7AEcg.Component, types.Component{Series: *series})
	return h
}

// AddRepresentativeBeatSeries adds a representative beat series.
//
// Similar to AddRhythmSeries but for derived representative beats.
func (h *Hl7xml) AddRepresentativeBeatSeries(
	startTime, endTime string,
	sampleRate float64,
	leads map[types.LeadCode][]int,
	origin, scale float64,
) *Hl7xml {
	series := h.buildSeries(
		types.REPRESENTATIVE_BEAT_CODE,
		startTime,
		endTime,
		sampleRate,
		leads,
		origin,
		scale,
	)

	h.HL7AEcg.Component = append(h.HL7AEcg.Component, types.Component{Series: *series})
	return h
}

// buildSeries constructs a Series with sequences for time and leads.
func (h *Hl7xml) buildSeries(
	seriesType types.SeriesTypeCode,
	startTime, endTime string,
	sampleRate float64,
	leads map[types.LeadCode][]int,
	origin, scale float64,
) *types.Series {
	series := &types.Series{
		ID:   &types.ID{Root: uuid.New().String()},
		Code: &types.Code[types.SeriesTypeCode, types.CodeSystemOID]{},
		EffectiveTime: types.EffectiveTime{
			Low:  types.Time{Value: startTime},
			High: types.Time{Value: endTime},
		},
	}
	series.Code.SetCode(seriesType, types.HL7_ActCode_OID, "")

	// Calculate increment from sample rate: increment = 1 / sampleRate seconds
	increment := 1.0 / sampleRate

	// Create sequence set with time sequence + lead sequences
	sequenceSet := types.SequenceSet{}

	// Add time sequence using the polymorphic SequenceValue (Typed + XsiType)
	timeSeq := types.SequenceComponent{
		Sequence: types.Sequence{
			Value: &types.SequenceValue{
				XsiType: "GLIST_TS",
				Typed: &types.GLIST_TS{
					Head: startTime,
					Increment: types.Increment{
						Value: formatFloat(increment),
						Unit:  "s",
					},
				},
			},
		},
	}
	// Set time code
	timeSeq.Sequence.Code.Time = &types.Code[types.TimeSequenceCode, types.CodeSystemOID]{}
	timeSeq.Sequence.Code.Time.SetCode(types.TIME_ABSOLUTE_CODE, types.HL7_ActCode_OID, "")
	sequenceSet.Component = append(sequenceSet.Component, timeSeq)

	// Add lead sequences
	for leadCode, samples := range leads {
		leadSeq := h.buildLeadSequence(leadCode, samples, origin, scale)
		sequenceSet.Component = append(sequenceSet.Component, leadSeq)
	}

	series.Component = []types.SeriesComponent{
		{SequenceSet: sequenceSet},
	}

	return series
}

func (h *Hl7xml) buildLeadSequence(
	leadCode types.LeadCode,
	samples []int,
	origin, scale float64,
) types.SequenceComponent {
	// Convert samples to space-separated string
	digitsStr := ""
	for i, sample := range samples {
		if i > 0 {
			digitsStr += " "
		}
		digitsStr += formatInt(sample)
	}

	seq := types.SequenceComponent{
		Sequence: types.Sequence{
			Value: &types.SequenceValue{
				XsiType: "SLIST_PQ",
				Typed: &types.SLIST_PQ{
					Origin: types.PhysicalQuantity{
						Value: formatFloat(origin),
						Unit:  "uV",
					},
					Scale: types.PhysicalQuantity{
						Value: formatFloat(scale),
						Unit:  "uV",
					},
					Digits: digitsStr,
				},
			},
		},
	}
	seq.Sequence.Code.Lead = &types.Code[types.LeadCode, types.CodeSystemOID]{}
	seq.Sequence.Code.Lead.SetCode(leadCode, types.MDC_OID, "")
	return seq
}

// SetSeriesAuthor sets the device that authored the series.
func (h *Hl7xml) SetSeriesAuthor(
	deviceID string,
	deviceType types.DeviceTypeCode,
	modelName string,
	softwareName string,
	manufacturerID string,
	manufacturerName string,
) *Hl7xml {
	if len(h.HL7AEcg.Component) == 0 {
		return h
	}

	// Apply to the most recently added series
	lastComponent := &h.HL7AEcg.Component[len(h.HL7AEcg.Component)-1]

	// buildLeadSequence creates a sequence for a single lead.
	lastComponent.Series.Author = &types.SeriesAuthor{
		ID: &types.ID{Root: deviceID},
		ManufacturedSeriesDevice: types.ManufacturedSeriesDevice{
			ID: &types.ID{
				Root:      manufacturerID,
				Extension: deviceID,
			},
			Code:                  types.NewCode(deviceType, types.CodeSystemOID(""), ""),
			ManufacturerModelName: &modelName,
			SoftwareName:          &softwareName,
		},
		ManufacturerOrganization: &types.ManufacturerOrganization{
			ID:   &types.ID{Root: manufacturerID},
			Name: &manufacturerName,
		},
	}

	return h
}

// =============================================================================
// Helper functions
// =============================================================================

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func formatInt(i int) string {
	return strconv.Itoa(i)
}
