package types

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// =============================================================================
// Sequence Set and Sequence Types
// =============================================================================

// SequenceSet is a set of sequences that all have the same length and contain
// related values.
//
// All sequences in a set must have the same length. The first value of each
// sequence is related to the first value of every other sequence. The second
// value of each sequence is related to the second value of every other sequence,
// and so on.
//
// A sequence set can be thought of as a table where every sequence is a column
// and the rows indicate which values are related.
//
// Examples:
//   - 12-lead ECG recorded simultaneously: 1 sequence set with time sequence + 12 lead sequences
//   - 12-lead ECG with 3 leads at a time: 4 sequence sets, each with time + 3 leads
//
// XML Structure:
//
//	<sequenceSet>
//	  <component>
//	    <sequence>
//	      <code code="TIME_ABSOLUTE" codeSystem="2.16.840.1.113883.5.4"/>
//	      <value xsi:type="GLIST_TS">...</value>
//	    </sequence>
//	  </component>
//	  <component>
//	    <sequence>
//	      <code code="MDC_ECG_LEAD_I" codeSystem="2.16.840.1.113883.6.24"/>
//	      <value xsi:type="SLIST_PQ">...</value>
//	    </sequence>
//	  </component>
//	</sequenceSet>
//
// Cardinality: Required (within SeriesComponent)
// Reference: HL7 aECG Implementation Guide, Page 34-36
type SequenceSet struct {
	// Component contains individual sequences (time and lead data).
	//
	// Must include at least one time sequence and one or more lead sequences.
	//
	// XML Tag: <component>...</component>
	// Cardinality: Required (1..*)
	Component []SequenceComponent `xml:"component"`
}

// SequenceComponent contains a single sequence within a sequence set.
//
// XML Structure:
//
//	<component>
//	  <sequence>...</sequence>
//	</component>
//
// Cardinality: Required (within SequenceSet, 1..*)
type SequenceComponent struct {
	// Sequence contains the ordered list of values.
	//
	// XML Tag: <sequence>...</sequence>
	// Cardinality: Required
	Sequence Sequence `xml:"sequence"`
}

// Sequence is an ordered list of values having a common code (or dimension).
//
// Sequence values are associated with other sequence values within a sequence set.
// For example, a 12-lead ECG series contains:
//   - 1 sequence for timestamps (when voltages were sampled)
//   - 12 sequences for voltages (measured at those times)
//
// XML Structure:
//
//	<sequence>
//	  <code code="TIME_ABSOLUTE" codeSystem="2.16.840.1.113883.5.4"/>
//	  <value xsi:type="GLIST_TS">
//	    <head value="20021122091000.000"/>
//	    <increment value="0.002" unit="s"/>
//	  </value>
//	</sequence>
//
// Or:
//
//	<sequence >
//	  <code code="MDC_ECG_LEAD_I" codeSystem="2.16.840.1.113883.6.24"/>
//	  <value xsi:type="SLIST_PQ">
//	    <origin value="0" unit="uV"/>
//	    <scale value="5" unit="uV"/>
//	    <digits>1 2 3 4 5</digits>
//	  </value>
//	</sequence>
//
// Cardinality: Required (within SequenceComponent)
// Reference: HL7 aECG Implementation Guide, Page 35-37
type Sequence struct {
	// Code names the dimension or type of values in the sequence.
	//
	// Time Sequence Codes (from HL7 ActCode, OID 2.16.840.1.113883.5.4):
	//   - "TIME_ABSOLUTE": Absolute time domain (Gregorian calendar)
	//   - "TIME_RELATIVE": Relative time domain (relative to series start)
	//
	// Voltage Sequence Codes (from MDC, OID 2.16.840.1.113883.6.24):
	//   - "MDC_ECG_LEAD_I"   through "MDC_ECG_LEAD_III"
	//   - "MDC_ECG_LEAD_AVR", "MDC_ECG_LEAD_AVL", "MDC_ECG_LEAD_AVF"
	//   - "MDC_ECG_LEAD_V1"  through "MDC_ECG_LEAD_V6"
	//
	// XML Tag: <code code="..." codeSystem="..."/>
	// Cardinality: Required
	Code SequenceCode `xml:"code"`

	// Value contains the list of values in the sequence.
	//
	// The value type depends on the sequence type:
	//   - Time sequences: GLIST_TS (Generated List of Timestamps)
	//   - Voltage sequences: SLIST_PQ (Scaled List of Physical Quantities)
	//   - Integer sequences: SLIST_INT (Scaled List of Integers)
	//
	// Note: In XML, this uses xsi:type to distinguish between types.
	// In Go, we'll need custom unmarshaling to handle different types.
	//
	// XML Tag: <value xsi:type="...">...</value>
	// Cardinality: Required
	Value *SequenceValue `xml:"value,omitempty"`
}

type SequenceCode struct {
	Time *Code[TimeSequenceCode, CodeSystemOID] `xml:",omitempty"`
	Lead *Code[LeadCode, CodeSystemOID]         `xml:",omitempty"`
}

// SequenceValue represents the polymorphic value field in a Sequence.
//
// This can contain:
//   - GLIST_TS: Generated list of timestamps
//   - SLIST_PQ: Scaled list of physical quantities
//   - SLIST_INT: Scaled list of integers
//
// In XML, the type is determined by the xsi:type attribute.
//
// Cardinality: Required (within Sequence)
type SequenceValue struct {
	XMLName xml.Name `xml:"value"`
	XsiType string   `xml:"xsi:type,attr,omitempty"`

	RawXML []byte `xml:",innerxml"`
	// Typed holds the decoded value as one of:
	//   - *GLIST_TS
	//   - *SLIST_PQ
	//   - *SLIST_INT
	Typed any `xml:"-"`
}

// =============================================================================
// GLIST_TS - Generated List of Timestamps
// =============================================================================

// GLIST_TS represents a generated list of timestamps.
//
// Used for time sequences where timestamps are generated using a simple algorithm.
// Most ECG devices use periodic sampling, making this ideal for generating the
// sequence of times at which voltages were sampled.
//
// The list is generated by:
//   - Starting at head value
//   - Adding increment repeatedly
//
// This efficiently represents thousands of timestamps with just 2 values!
//
// Example: 5000 samples at 500 Hz (0.002s increment)
//
//	<value xsi:type="GLIST_TS">
//	  <head value="20021122091000.000"/>
//	  <increment value="0.002" unit="s"/>
//	</value>
//
// This generates timestamps:
//   - 20021122091000.000 (head)
//   - 20021122091000.002 (head + 1*increment)
//   - 20021122091000.004 (head + 2*increment)
//   - ... (5000 total values)
//
// XML Attribute: xsi:type="GLIST_TS"
// Reference: HL7 aECG Implementation Guide, Page 36
type GLIST_TS struct {
	// Head is the first timestamp in the sequence.
	//
	// Format: HL7 timestamp (YYYYMMDDHHmmss.SSS)
	//
	// Example: "20021122091000.000" = Nov 22, 2002 at 09:10:00.000
	//
	// XML Tag: <head value="..."/>
	// Cardinality: Required
	Head string `xml:"head"`

	// Increment is the time interval between samples.
	//
	// Common Values:
	//   - "0.002" unit="s": 500 Hz sampling (standard 12-lead)
	//   - "0.001" unit="s": 1000 Hz sampling (high resolution)
	//   - "0.004" unit="s": 250 Hz sampling (Holter)
	//   - "0.008" unit="s": 125 Hz sampling (low resolution)
	//
	// The increment is added to the previous timestamp to generate the next.
	//
	// XML Tag: <increment value="..." unit="..."/>
	// Cardinality: Required
	Increment Increment `xml:"increment"`
}

// Increment represents a time increment with value and unit.
//
// Used in GLIST_TS to define the interval between samples.
//
// Cardinality: Required (within GLIST_TS)
type Increment struct {
	// Value is the numeric increment value.
	//
	// Examples: "0.002" (2 milliseconds), "0.001" (1 millisecond)
	//
	// XML Tag: value="..."
	// Cardinality: Required
	Value string `xml:"value,attr"`

	// Unit is the unit of time.
	//
	// Common Values:
	//   - "s": seconds (most common)
	//   - "ms": milliseconds
	//
	// XML Tag: unit="..."
	// Cardinality: Required
	Unit string `xml:"unit,attr"`
}

// =============================================================================
// SLIST_PQ - Scaled List of Physical Quantities
// =============================================================================

// SLIST_PQ represents a scaled list of physical quantities.
//
// Used to enumerate voltages in ECG waveforms efficiently. Instead of storing
// full floating-point values, this stores:
//   - Origin: baseline value
//   - Scale: multiplication factor
//   - Digits: raw integer values
//
// This allows efficient storage of raw integer values from the device.
// Actual values are calculated as: actualValue = origin + (digit * scale)
//
// Example: 5 samples with 5 µV resolution
//
//	<value xsi:type="SLIST_PQ">
//	  <origin value="0" unit="uV"/>
//	  <scale value="5" unit="uV"/>
//	  <digits>1 2 3 4 5</digits>
//	</value>
//
// Actual voltages:
//   - 0 + (1 * 5) = 5 µV
//   - 0 + (2 * 5) = 10 µV
//   - 0 + (3 * 5) = 15 µV
//   - 0 + (4 * 5) = 20 µV
//   - 0 + (5 * 5) = 25 µV
//
// XML Attribute: xsi:type="SLIST_PQ"
// Reference: HL7 aECG Implementation Guide, Page 36
type SLIST_PQ struct {
	// Origin is the baseline value (offset).
	//
	// This is added to the scaled digit values.
	// Often "0" for ECG data referenced to baseline.
	//
	// Example: value="0" unit="uV"
	//
	// XML Tag: <origin value="..." unit="..."/>
	// Cardinality: Required
	Origin PhysicalQuantity `xml:"origin"`

	// Scale is the multiplication factor applied to raw digits.
	//
	// This represents the amplitude resolution of the device.
	//
	// Common Values:
	//   - 2.5 µV (high resolution)
	//   - 5 µV (standard resolution)
	//   - 10 µV (lower resolution)
	//
	// Example: value="5" unit="uV"
	//
	// XML Tag: <scale value="..." unit="..."/>
	// Cardinality: Required
	Scale PhysicalQuantity `xml:"scale"`

	// Digits are the raw integer values from the device.
	//
	// These are space-separated integers that are scaled and offset
	// to get actual physical values.
	//
	// Calculation: actualValue = origin + (digit * scale)
	//
	// Example: "1 2 3 4 5" or "-2 -2 -2 -2 -3 -4 -3 -5 -5 -4 -6 -9 -9"
	//
	// The string may contain thousands of values separated by spaces.
	//
	// XML Tag: <digits>...</digits>
	// Cardinality: Required
	Digits string `xml:"digits"`
}

// PhysicalQuantity represents a physical measurement with value and unit.
//
// Used for origin and scale in SLIST_PQ.
//
// Cardinality: Required (in SLIST_PQ context)
type PhysicalQuantity struct {
	// Value is the numeric value.
	//
	// Example: "0", "5", "2.5"
	//
	// XML Tag: value="..."
	// Cardinality: Required
	Value string `xml:"value,attr"`

	// Unit is the unit of measurement.
	//
	// Common ECG Units:
	//   - "uV": microvolts (most common for ECG)
	//   - "mV": millivolts
	//   - "V": volts
	//
	// XML Tag: unit="..."
	// Cardinality: Required
	Unit string `xml:"unit,attr"`
}

// GetDigits parses the Digits string into a slice of integers.
//
// Returns a slice of integer values parsed from the space-separated string.
// Returns error if parsing fails.
func (s *SLIST_PQ) GetDigits() ([]int, error) {
	if s.Digits == "" {
		return []int{}, nil
	}

	digitStrs := strings.Fields(s.Digits)
	digits := make([]int, len(digitStrs))

	for i, str := range digitStrs {
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		digits[i] = val
	}

	return digits, nil
}

// GetLength returns the number of samples in this sequence.
func (s *SLIST_PQ) GetLength() int {
	if s.Digits == "" {
		return 0
	}
	return len(strings.Fields(s.Digits))
}

// GetActualValues calculates the actual physical values from digits.
//
// Calculation: actualValue = origin + (digit * scale)
//
// Returns a slice of floating-point values in the units specified by Origin/Scale.
// Returns error if parsing fails.
func (s *SLIST_PQ) GetActualValues() ([]float64, error) {
	digits, err := s.GetDigits()
	if err != nil {
		return nil, err
	}

	origin, err := strconv.ParseFloat(s.Origin.Value, 64)
	if err != nil {
		return nil, err
	}

	scale, err := strconv.ParseFloat(s.Scale.Value, 64)
	if err != nil {
		return nil, err
	}

	values := make([]float64, len(digits))
	for i, digit := range digits {
		values[i] = origin + (float64(digit) * scale)
	}

	return values, nil
}

// =============================================================================
// SLIST_INT - Scaled List of Integers
// =============================================================================

// SLIST_INT represents a scaled list of integer values.
//
// Similar to SLIST_PQ but for integer sequences rather than physical quantities.
// Less commonly used in ECG data but available for integer-valued measurements.
//
// Calculation: actualValue = origin + (digit * scale)
//
// XML Attribute: xsi:type="SLIST_INT"
// Reference: HL7 aECG Implementation Guide, Page 36
type SLIST_INT struct {
	// Origin is the baseline integer value (offset).
	//
	// XML Tag: <origin value="..."/>
	// Cardinality: Required
	Origin int `xml:"origin,attr"`

	// Scale is the multiplication factor applied to raw digits.
	//
	// XML Tag: <scale value="..."/>
	// Cardinality: Required
	Scale int `xml:"scale,attr"`

	// Digits are the raw integer values.
	//
	// Space-separated integers that are scaled and offset.
	//
	// XML Tag: <digits>...</digits>
	// Cardinality: Required
	Digits string `xml:"digits"`
}

// GetDigits parses the Digits string into a slice of integers.
//
// Returns a slice of integer values parsed from the space-separated string.
// Returns error if parsing fails.
func (s *SLIST_INT) GetDigits() ([]int, error) {
	if s.Digits == "" {
		return []int{}, nil
	}

	digitStrs := strings.Fields(s.Digits)
	digits := make([]int, len(digitStrs))

	for i, str := range digitStrs {
		val, err := strconv.Atoi(str)
		if err != nil {
			return nil, err
		}
		digits[i] = val
	}

	return digits, nil
}

// GetLength returns the number of values in this sequence.
func (s *SLIST_INT) GetLength() int {
	if s.Digits == "" {
		return 0
	}
	return len(strings.Fields(s.Digits))
}

// GetActualValues calculates the actual integer values from digits.
//
// Calculation: actualValue = origin + (digit * scale)
//
// Returns a slice of integer values.
// Returns error if parsing fails.
func (s *SLIST_INT) GetActualValues() ([]int, error) {
	digits, err := s.GetDigits()
	if err != nil {
		return nil, err
	}

	values := make([]int, len(digits))
	for i, digit := range digits {
		values[i] = s.Origin + (digit * s.Scale)
	}

	return values, nil
}

// UnmarshalXML handles decoding of <value xsi:type="...">...</value>
func (sv *SequenceValue) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type Alias SequenceValue
	aux := &Alias{}

	if err := d.DecodeElement(aux, &start); err != nil {
		return err
	}
	*sv = SequenceValue(*aux)

	switch sv.XsiType {
	case "GLIST_TS":
		var v GLIST_TS
		if err := xml.Unmarshal(sv.RawXML, &v); err != nil {
			return err
		}
		sv.Typed = &v
	case "SLIST_PQ":
		var v SLIST_PQ
		if err := xml.Unmarshal(sv.RawXML, &v); err != nil {
			return err
		}
		sv.Typed = &v
	case "SLIST_INT":
		var v SLIST_INT
		if err := xml.Unmarshal(sv.RawXML, &v); err != nil {
			return err
		}
		sv.Typed = &v
	default:
		// Unknown type — keep raw XML
		sv.Typed = nil
	}

	return nil
}

// MarshalXML handles encoding of the correct subtype based on XsiType.
func (sv *SequenceValue) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{Local: "value"}
	if sv.XsiType != "" {
		start.Attr = append(start.Attr, xml.Attr{
			Name:  xml.Name{Local: "xsi:type"},
			Value: sv.XsiType,
		})
	}

	var inner any

	switch sv.XsiType {
	case "GLIST_TS":
		inner, _ = sv.Typed.(*GLIST_TS)
	case "SLIST_PQ":
		inner, _ = sv.Typed.(*SLIST_PQ)
	case "SLIST_INT":
		inner, _ = sv.Typed.(*SLIST_INT)
	}

	// If type not recognized or Typed is nil, fallback to raw XML
	if inner == nil {
		if len(sv.RawXML) > 0 {
			e.EncodeToken(start)
			e.EncodeToken(xml.CharData(sv.RawXML))
			e.EncodeToken(xml.EndElement{Name: start.Name})
			return nil
		}
		return e.EncodeElement(nil, start)
	}

	// Encode with the correct struct
	return e.EncodeElement(inner, start)
}
