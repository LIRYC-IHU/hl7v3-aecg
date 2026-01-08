package hl7aecg

import (
	"slices"
	"strconv"

	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
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
	inclusive_low, inclusive_high *bool,
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
	if inclusive_low != nil {
		series.EffectiveTime.Low.Inclusive = inclusive_low
	}
	if inclusive_high != nil {
		series.EffectiveTime.High.Inclusive = inclusive_high
	}

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
		ID:   &types.ID{},
		Code: &types.Code[types.SeriesTypeCode, types.CodeSystemOID]{},
		EffectiveTime: types.EffectiveTime{
			Low:  types.Time{Value: startTime},
			High: types.Time{Value: endTime},
		},
	}
	series.Code.SetCode(seriesType, types.HL7_ActCode_OID, "", "")

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
					Head: types.HeadTimestamp{
						Value: startTime,
						Unit:  "s",
					},
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
	timeSeq.Sequence.Code.Time.SetCode(types.TIME_ABSOLUTE_CODE, types.HL7_ActCode_OID, "ActCode", "")
	sequenceSet.Component = append(sequenceSet.Component, timeSeq)

	// Add lead sequences in the standard medical order:
	// Limb leads: I, II, III
	// Augmented leads: aVR, aVL, aVF
	// Precordial leads: V1, V2, V3, V4, V5, V6
	standardOrder := []types.LeadCode{
		types.MDC_ECG_LEAD_I,
		types.MDC_ECG_LEAD_II,
		types.MDC_ECG_LEAD_III,
		types.MDC_ECG_LEAD_AVR,
		types.MDC_ECG_LEAD_AVL,
		types.MDC_ECG_LEAD_AVF,
		types.MDC_ECG_LEAD_V1,
		types.MDC_ECG_LEAD_V2,
		types.MDC_ECG_LEAD_V3,
		types.MDC_ECG_LEAD_V4,
		types.MDC_ECG_LEAD_V5,
		types.MDC_ECG_LEAD_V6,
	}

	// Iterate in standard order, only adding leads that are present in the map
	for _, leadCode := range standardOrder {
		if samples, exists := leads[leadCode]; exists {
			leadSeq := h.buildLeadSequence(leadCode, samples, origin, scale)
			sequenceSet.Component = append(sequenceSet.Component, leadSeq)
		}
	}

	// Add any remaining leads that aren't in the standard 12-lead set
	for leadCode, samples := range leads {
		if !slices.Contains(standardOrder, leadCode) {
			leadSeq := h.buildLeadSequence(leadCode, samples, origin, scale)
			sequenceSet.Component = append(sequenceSet.Component, leadSeq)
		}
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
	seq.Sequence.Code.Lead.SetCode(leadCode, types.MDC_OID, "MDC", "")
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
	lastComponent.Series.Author = &types.Author{
		SeriesAuthor: types.SeriesAuthor{
			ID: &types.ID{Root: deviceID},
			ManufacturedSeriesDevice: types.ManufacturedSeriesDevice{
				ID: &types.ID{
					Root:      manufacturerID,
					Extension: deviceID,
				},
				Code:                  types.NewCode(deviceType, types.CodeSystemOID(""), "", ""),
				ManufacturerModelName: &modelName,
				SoftwareName:          &softwareName,
			},
			ManufacturerOrganization: &types.ManufacturerOrganization{
				ID:   &types.ID{Root: manufacturerID},
				Name: &manufacturerName,
			},
		},
	}

	return h
}

// AddSecondaryPerformer adds a secondary performer (technician) to the most recently added series.
//
// The secondary performer describes the technician who operated the device
// that captured the ECG waveforms.
//
// Parameters:
//   - functionCode: The function the technician was performing (e.g., PERFORMER_ECG_TECHNICIAN)
//   - performerID: Optional role-specific identifier root for the technician
//   - performerExtension: Optional role-specific identifier extension
//   - name: The technician's name (empty string for <name/> element, or actual name)
//
// Example:
//
//	h.AddRhythmSeries(...).
//	  AddSecondaryPerformer(types.PERFORMER_ECG_TECHNICIAN, "", "", "")
//
// Returns the Hl7xml instance for method chaining.
func (h *Hl7xml) AddSecondaryPerformer(
	functionCode types.PerformerFunctionCode,
	performerID string,
	performerExtension string,
	name string,
) *Hl7xml {
	if len(h.HL7AEcg.Component) == 0 {
		return h
	}

	// Apply to the most recently added series
	lastComponent := &h.HL7AEcg.Component[len(h.HL7AEcg.Component)-1]

	// Create secondary performer
	performer := types.SecondaryPerformer{
		SeriesPerformer: types.SeriesPerformer{},
	}

	// Set function code
	if functionCode != "" {
		performer.SetFunctionCode(functionCode, types.CodeSystemOID(""), "", "")
	}

	// Set performer details
	performer.SeriesPerformer.SetPerformer(performerID, performerExtension, name)

	// Add to series
	lastComponent.Series.SecondaryPerformer = append(lastComponent.Series.SecondaryPerformer, performer)

	return h
}

// AddControlVariable adds a control variable to the most recently added series.
//
// Control variables capture related information about ECG collection conditions,
// such as filter settings, subject observations (age, etc.), or other parameters.
//
// Parameters:
//   - cv: A pointer to the ControlVariable to add (use types.New* functions to create)
//
// Example:
//
//	// Add a low pass filter
//	h.AddControlVariable(types.NewLowPassFilter("35", "Hz"))
//
//	// Add a custom control variable
//	cv := types.NewControlVariable("CUSTOM_CODE", types.MDC_OID, "MDC", "Custom Parameter")
//	cv.ControlVariable.SetValue("100", "unit")
//	h.AddControlVariable(cv)
//
// Returns the Hl7xml instance for method chaining.
func (h *Hl7xml) AddControlVariable(cv *types.ControlVariable) *Hl7xml {
	if len(h.HL7AEcg.Component) == 0 {
		return h
	}

	// Apply to the most recently added series
	lastComponent := &h.HL7AEcg.Component[len(h.HL7AEcg.Component)-1]

	// Add to series
	lastComponent.Series.ControlVariable = append(lastComponent.Series.ControlVariable, *cv)

	return h
}

// AddLowPassFilter adds a low pass filter control variable to the most recently added series.
//
// This is a convenience method for the common case of adding a low pass filter specification.
//
// Parameters:
//   - cutoffFreq: The cutoff frequency value as a string (e.g., "35")
//   - unit: The frequency unit (typically "Hz")
//
// Example:
//
//	h.AddLowPassFilter("35", "Hz")
//
// Returns the Hl7xml instance for method chaining.
func (h *Hl7xml) AddLowPassFilter(cutoffFreq, unit string) *Hl7xml {
	return h.AddControlVariable(types.NewLowPassFilter(cutoffFreq, unit))
}

// AddHighPassFilter adds a high pass filter control variable to the most recently added series.
//
// This is a convenience method for the common case of adding a high pass filter specification.
//
// Parameters:
//   - cutoffFreq: The cutoff frequency value as a string (e.g., "0.56")
//   - unit: The frequency unit (typically "Hz")
//
// Example:
//
//	h.AddHighPassFilter("0.56", "Hz")
//
// Returns the Hl7xml instance for method chaining.
func (h *Hl7xml) AddHighPassFilter(cutoffFreq, unit string) *Hl7xml {
	return h.AddControlVariable(types.NewHighPassFilter(cutoffFreq, unit))
}

// AddNotchFilter adds a notch filter control variable to the most recently added series.
//
// This is a convenience method for the common case of adding a notch filter specification.
//
// Parameters:
//   - notchFreq: The notch frequency value as a string (e.g., "50" or "60")
//   - unit: The frequency unit (typically "Hz")
//
// Example:
//
//	h.AddNotchFilter("50", "Hz")  // 50 Hz for Europe
//	h.AddNotchFilter("60", "Hz")  // 60 Hz for North America
//
// Returns the Hl7xml instance for method chaining.
func (h *Hl7xml) AddNotchFilter(notchFreq, unit string) *Hl7xml {
	return h.AddControlVariable(types.NewNotchFilter(notchFreq, unit))
}

// AddAgeObservation adds an age observation control variable to the most recently added series.
//
// This is a convenience method for the common case of recording the subject's age.
//
// Parameters:
//   - age: The age value as a string (e.g., "34")
//   - unit: The unit (typically "a" for years)
//
// Example:
//
//	h.AddAgeObservation("34", "a")
//
// Returns the Hl7xml instance for method chaining.
func (h *Hl7xml) AddAgeObservation(age, unit string) *Hl7xml {
	return h.AddControlVariable(types.NewAgeObservation(age, unit))
}


// =============================================================================
// Add Confidentiality Code
// =============================================================================
func (h *Hl7xml) AddConfidentialityCode(code types.ConfidentialityCode) *Hl7xml {
	h.HL7AEcg.ConfidentialityCode = &types.Code[types.ConfidentialityCode, string]{}
	h.HL7AEcg.ConfidentialityCode.SetCode(code, "", "", "")
	return h
}

// ============================================================================
// Add Reason Code
// =============================================================================
func (h *Hl7xml) AddReasonCode(code types.ReasonCode) *Hl7xml {
	h.HL7AEcg.ReasonCode = &types.Code[types.ReasonCode, string]{}
	h.HL7AEcg.ReasonCode.SetCode(code, "", "", "")
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
