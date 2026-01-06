package hl7aecg

import (
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
	timeSeq.Sequence.Code.Time.SetCode(types.TIME_ABSOLUTE_CODE, types.HL7_ActCode_OID, "ActCode", "")
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

// AddControlVariable adds a control variable (observation) to the most recently added series.
//
// Control variables capture related information about the subject or ECG collection conditions,
// such as the subject's age, fasting status, or other clinical information.
//
// Parameters:
//   - observationCode: The observation code (e.g., "21612-7" for "Reported Age")
//   - codeSystem: The code system OID (e.g., types.LOINC_OID)
//   - displayName: Display name for the code (e.g., "Reported Age")
//   - value: The observed value (e.g., "34")
//   - unit: The unit of measurement (e.g., "a" for years)
//
// Example:
//
//	h.AddRhythmSeries(...).
//	  AddControlVariable("21612-7", types.LOINC_OID, "Reported Age", "34", "a")
//
// Returns the Hl7xml instance for method chaining.
func (h *Hl7xml) AddControlVariable(
	observationCode string,
	codeSystem types.CodeSystemOID,
	displayName string,
	value string,
	unit string,
) *Hl7xml {
	if len(h.HL7AEcg.Component) == 0 {
		return h
	}

	// Apply to the most recently added series
	lastComponent := &h.HL7AEcg.Component[len(h.HL7AEcg.Component)-1]

	// Create control variable
	controlVar := types.ControlVariable{
		RelatedObservation: types.RelatedObservation{},
	}

	// Set observation code
	controlVar.RelatedObservation.SetObservationCode(observationCode, codeSystem, displayName, "LOINC")

	// Set value
	if value != "" {
		controlVar.RelatedObservation.SetValue(value, unit)
	}

	// Add to series
	lastComponent.Series.ControlVariable = append(lastComponent.Series.ControlVariable, controlVar)

	return h
}

// AddControlVariableWithAuthor adds a control variable with author information to the most recently added series.
//
// This variant allows specifying who recorded the observation.
//
// Parameters:
//   - observationCode: The observation code
//   - codeSystem: The code system OID
//   - displayName: Display name for the code
//   - value: The observed value
//   - unit: The unit of measurement
//   - authorID: The author's ID root
//   - authorExtension: The author's ID extension
//   - authorName: The author's name
//
// Example:
//
//	h.AddRhythmSeries(...).
//	  AddControlVariableWithAuthor("21612-7", types.LOINC_OID, "Reported Age",
//	    "34", "a", "2.16.840.1.113883.3.5", "TECH_23", "JMK")
//
// Returns the Hl7xml instance for method chaining.
func (h *Hl7xml) AddControlVariableWithAuthor(
	observationCode string,
	codeSystem types.CodeSystemOID,
	displayName string,
	value string,
	unit string,
	authorID string,
	authorExtension string,
	authorName string,
) *Hl7xml {
	if len(h.HL7AEcg.Component) == 0 {
		return h
	}

	// Apply to the most recently added series
	lastComponent := &h.HL7AEcg.Component[len(h.HL7AEcg.Component)-1]

	// Create control variable
	controlVar := types.ControlVariable{
		RelatedObservation: types.RelatedObservation{},
	}

	// Set observation code and value
	controlVar.RelatedObservation.
		SetObservationCode(observationCode, codeSystem, displayName, "LOINC").
		SetValue(value, unit).
		SetAuthorPerson(authorID, authorExtension, authorName)

	// Add to series
	lastComponent.Series.ControlVariable = append(lastComponent.Series.ControlVariable, controlVar)

	return h
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
