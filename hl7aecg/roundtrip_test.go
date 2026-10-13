package hl7aecg

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
)

// A library that writes aECG and reads aECG has to be able to read what it
// writes. It could not: String() used a marshaller that does not invoke this
// package's custom MarshalXML methods, so it emitted lead codes in a nested
// form UnmarshalXML rejects and dropped the body of every <value> element. A
// document round-tripped to zero leads and no waveform data, silently.
//
// These tests pin the round trip itself rather than the marshaller, so they
// keep holding whatever it is replaced with.

func buildSample(t *testing.T) *Hl7xml {
	t.Helper()
	leads := map[types.LeadCode][]int{
		types.MDC_ECG_LEAD_I:   {1, 2, 3, 4, 5},
		types.MDC_ECG_LEAD_II:  {6, 7, 8, 9, 10},
		types.MDC_ECG_LEAD_AVR: {-1, -2, -3, -4, -5},
	}
	h := NewHl7xml("").Initialize("93000", types.CodeSystemOID("2.16.840.1.113883.6.12"), "CPT4", "ECG")
	h.SetEffectiveTime("20260101120000", "20260101120010", nil, nil)
	h.AddRhythmSeries("20260101120000", "20260101120010", nil, nil, 500, leads, 0, 5)
	return h
}

// writeAndReadBack renders the document with String() and parses it back.
func writeAndReadBack(t *testing.T, h *Hl7xml) *Hl7xml {
	t.Helper()
	out, err := h.String()
	if err != nil {
		t.Fatalf("String(): %v", err)
	}
	path := filepath.Join(t.TempDir(), "aecg.xml")
	if err := os.WriteFile(path, []byte(out), 0o600); err != nil {
		t.Fatalf("writing document: %v", err)
	}

	back := NewHl7xml("")
	if err := back.UnmarshalFromFile(path); err != nil {
		t.Fatalf("UnmarshalFromFile: %v", err)
	}
	return back
}

// leadSequences collects the voltage sequences of a parsed document.
func leadSequences(h *Hl7xml) map[types.LeadCode]*types.SLIST_PQ {
	out := map[types.LeadCode]*types.SLIST_PQ{}
	for _, c := range h.HL7AEcg.Component {
		for _, sc := range c.Series.Component {
			for _, comp := range sc.SequenceSet.Component {
				seq := comp.Sequence
				if seq.Code.Lead == nil || seq.Value == nil {
					continue
				}
				if v, ok := seq.Value.Typed.(*types.SLIST_PQ); ok {
					out[seq.Code.Lead.Code] = v
				}
			}
		}
	}
	return out
}

func TestStringRoundTripsLeadSequences(t *testing.T) {
	back := writeAndReadBack(t, buildSample(t))

	leads := leadSequences(back)
	if len(leads) != 3 {
		t.Fatalf("read back %d lead sequences, want 3 — the document did not survive String()", len(leads))
	}

	for _, code := range []types.LeadCode{types.MDC_ECG_LEAD_I, types.MDC_ECG_LEAD_II, types.MDC_ECG_LEAD_AVR} {
		v, ok := leads[code]
		if !ok {
			t.Errorf("lead %s missing after round trip", code)
			continue
		}
		digits, err := v.GetDigits()
		if err != nil {
			t.Errorf("lead %s digits: %v", code, err)
			continue
		}
		if len(digits) != 5 {
			t.Errorf("lead %s has %d samples after round trip, want 5", code, len(digits))
		}
	}
}

// The scale and origin are what turn digits into microvolts. Losing them is
// the difference between a waveform and a list of integers.
func TestStringRoundTripsScaleAndOrigin(t *testing.T) {
	back := writeAndReadBack(t, buildSample(t))

	v, ok := leadSequences(back)[types.MDC_ECG_LEAD_I]
	if !ok {
		t.Fatal("lead I missing after round trip")
	}
	if v.Scale.Value != "5" {
		t.Errorf("scale value = %q, want \"5\"", v.Scale.Value)
	}
	if v.Scale.Unit != "uV" {
		t.Errorf("scale unit = %q, want \"uV\"", v.Scale.Unit)
	}
	if v.Origin.Value != "0" {
		t.Errorf("origin value = %q, want \"0\"", v.Origin.Value)
	}
}

func TestStringRoundTripsDigitValues(t *testing.T) {
	back := writeAndReadBack(t, buildSample(t))

	v, ok := leadSequences(back)[types.MDC_ECG_LEAD_AVR]
	if !ok {
		t.Fatal("lead AVR missing after round trip")
	}
	digits, err := v.GetDigits()
	if err != nil {
		t.Fatalf("GetDigits: %v", err)
	}
	want := []int{-1, -2, -3, -4, -5}
	if len(digits) != len(want) {
		t.Fatalf("got %d digits, want %d", len(digits), len(want))
	}
	for i := range want {
		if digits[i] != want[i] {
			t.Errorf("digit %d = %d, want %d", i, digits[i], want[i])
		}
	}
}

// The time sequence carries the sampling interval; without it a consumer has
// no time base for the waveform.
func TestStringRoundTripsTimeSequence(t *testing.T) {
	back := writeAndReadBack(t, buildSample(t))

	var found bool
	for _, c := range back.HL7AEcg.Component {
		for _, sc := range c.Series.Component {
			for _, comp := range sc.SequenceSet.Component {
				seq := comp.Sequence
				if seq.Code.Time == nil || seq.Value == nil {
					continue
				}
				g, ok := seq.Value.Typed.(*types.GLIST_TS)
				if !ok {
					t.Fatalf("time sequence decoded as %T, want *GLIST_TS", seq.Value.Typed)
				}
				found = true
				if g.Increment.Value != "0.002" {
					t.Errorf("increment = %q, want \"0.002\" (500 Hz)", g.Increment.Value)
				}
				if g.Increment.Unit != "s" {
					t.Errorf("increment unit = %q, want \"s\"", g.Increment.Unit)
				}
			}
		}
	}
	if !found {
		t.Error("no time sequence survived the round trip")
	}
}

// The FDA aECG Implementation Guide writes sequence codes flat, as
// <code code="MDC_ECG_LEAD_I" codeSystem="…"/>. The nested form a previous
// marshaller produced is both non-conformant and unreadable by this package.
func TestStringWritesSequenceCodesFlat(t *testing.T) {
	out, err := buildSample(t).String()
	if err != nil {
		t.Fatalf("String(): %v", err)
	}

	if !contains(out, `<code code="MDC_ECG_LEAD_I"`) {
		t.Error("lead code is not written in the flat form the guide specifies")
	}
	if contains(out, "<Lead ") {
		t.Error("lead code is written nested; UnmarshalXML cannot read that form")
	}
	// The typed value body must be present, not an empty self-closing element.
	if !contains(out, "<digits>") {
		t.Error("the <value> body was dropped: no <digits> in the output")
	}
	if !contains(out, `<scale value="5"`) {
		t.Error("the <value> body was dropped: no <scale> in the output")
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && indexOf(haystack, needle) >= 0
}

func indexOf(haystack, needle string) int {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return i
		}
	}
	return -1
}

// The head of a GLIST_TS is an HL7 TS — an instant, not a measured quantity —
// and the FDA aECG Implementation Guide writes it as <head value="…"/> with no
// unit. This package used to emit unit="s" there, which a strict schema
// validator has every reason to reject.
func TestTimeSequenceHeadCarriesNoUnit(t *testing.T) {
	out, err := buildSample(t).String()
	if err != nil {
		t.Fatalf("String(): %v", err)
	}
	if contains(out, `<head value="20260101120000" unit=`) {
		t.Error(`the GLIST_TS head is written with a unit attribute; the guide writes <head value="…"/>`)
	}
	if !contains(out, `<head value="20260101120000">`) && !contains(out, `<head value="20260101120000"/>`) {
		t.Errorf("the head is not written as expected:\n%s", out)
	}
	// The increment is a duration and must keep its unit.
	if !contains(out, `<increment value="0.002" unit="s"`) {
		t.Error("the increment lost its unit; it is a duration and needs one")
	}
}
