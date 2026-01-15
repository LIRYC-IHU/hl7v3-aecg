package types

import (
	"encoding/xml"
	"reflect"
	"testing"
)

// TestSLIST_PQ_GetDigits tests the GetDigits method for SLIST_PQ
func TestSLIST_PQ_GetDigits(t *testing.T) {
	tests := []struct {
		name    string
		slist   SLIST_PQ
		want    []int
		wantErr bool
	}{
		{
			name: "Valid digits with spaces",
			slist: SLIST_PQ{
				Digits: "1 2 3 4 5",
			},
			want:    []int{1, 2, 3, 4, 5},
			wantErr: false,
		},
		{
			name: "Valid digits with multiple spaces",
			slist: SLIST_PQ{
				Digits: "10  20   30",
			},
			want:    []int{10, 20, 30},
			wantErr: false,
		},
		{
			name: "Valid negative digits",
			slist: SLIST_PQ{
				Digits: "-5 -10 0 10 5",
			},
			want:    []int{-5, -10, 0, 10, 5},
			wantErr: false,
		},
		{
			name: "Empty digits",
			slist: SLIST_PQ{
				Digits: "",
			},
			want:    []int{},
			wantErr: false,
		},
		{
			name: "Invalid digits (non-numeric)",
			slist: SLIST_PQ{
				Digits: "1 abc 3",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.slist.GetDigits()

			if (err != nil) != tt.wantErr {
				t.Errorf("SLIST_PQ.GetDigits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SLIST_PQ.GetDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSLIST_PQ_GetLength tests the GetLength method for SLIST_PQ
func TestSLIST_PQ_GetLength(t *testing.T) {
	tests := []struct {
		name  string
		slist SLIST_PQ
		want  int
	}{
		{
			name: "Five digits",
			slist: SLIST_PQ{
				Digits: "1 2 3 4 5",
			},
			want: 5,
		},
		{
			name: "Three digits",
			slist: SLIST_PQ{
				Digits: "10 20 30",
			},
			want: 3,
		},
		{
			name: "Empty digits",
			slist: SLIST_PQ{
				Digits: "",
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.slist.GetLength(); got != tt.want {
				t.Errorf("SLIST_PQ.GetLength() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSLIST_PQ_GetActualValues tests the GetActualValues method for SLIST_PQ
func TestSLIST_PQ_GetActualValues(t *testing.T) {
	tests := []struct {
		name    string
		slist   SLIST_PQ
		want    []float64
		wantErr bool
	}{
		{
			name: "Basic calculation with origin 0 and scale 5",
			slist: SLIST_PQ{
				Origin: PhysicalQuantity{Value: "0", Unit: "uV"},
				Scale:  PhysicalQuantity{Value: "5", Unit: "uV"},
				Digits: "0 1 2 3 4",
			},
			want:    []float64{0, 5, 10, 15, 20},
			wantErr: false,
		},
		{
			name: "Calculation with non-zero origin",
			slist: SLIST_PQ{
				Origin: PhysicalQuantity{Value: "100", Unit: "uV"},
				Scale:  PhysicalQuantity{Value: "10", Unit: "uV"},
				Digits: "0 1 2 -1 -2",
			},
			want:    []float64{100, 110, 120, 90, 80},
			wantErr: false,
		},
		{
			name: "Calculation with fractional scale",
			slist: SLIST_PQ{
				Origin: PhysicalQuantity{Value: "0", Unit: "uV"},
				Scale:  PhysicalQuantity{Value: "2.5", Unit: "uV"},
				Digits: "0 2 4 -2 -4",
			},
			want:    []float64{0, 5, 10, -5, -10},
			wantErr: false,
		},
		{
			name: "Empty digits",
			slist: SLIST_PQ{
				Origin: PhysicalQuantity{Value: "0", Unit: "uV"},
				Scale:  PhysicalQuantity{Value: "5", Unit: "uV"},
				Digits: "",
			},
			want:    []float64{},
			wantErr: false,
		},
		{
			name: "Invalid origin value",
			slist: SLIST_PQ{
				Origin: PhysicalQuantity{Value: "abc", Unit: "uV"},
				Scale:  PhysicalQuantity{Value: "5", Unit: "uV"},
				Digits: "1 2 3",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid scale value",
			slist: SLIST_PQ{
				Origin: PhysicalQuantity{Value: "0", Unit: "uV"},
				Scale:  PhysicalQuantity{Value: "xyz", Unit: "uV"},
				Digits: "1 2 3",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.slist.GetActualValues()

			if (err != nil) != tt.wantErr {
				t.Errorf("SLIST_PQ.GetActualValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SLIST_PQ.GetActualValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSLIST_INT_GetDigits tests the GetDigits method for SLIST_INT
func TestSLIST_INT_GetDigits(t *testing.T) {
	tests := []struct {
		name    string
		slist   SLIST_INT
		want    []int
		wantErr bool
	}{
		{
			name: "Valid integer digits",
			slist: SLIST_INT{
				Digits: "1 2 3 4 5",
			},
			want:    []int{1, 2, 3, 4, 5},
			wantErr: false,
		},
		{
			name: "Valid negative integers",
			slist: SLIST_INT{
				Digits: "-10 -5 0 5 10",
			},
			want:    []int{-10, -5, 0, 5, 10},
			wantErr: false,
		},
		{
			name: "Empty digits",
			slist: SLIST_INT{
				Digits: "",
			},
			want:    []int{},
			wantErr: false,
		},
		{
			name: "Invalid digits (non-numeric)",
			slist: SLIST_INT{
				Digits: "1 abc 3",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.slist.GetDigits()

			if (err != nil) != tt.wantErr {
				t.Errorf("SLIST_INT.GetDigits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SLIST_INT.GetDigits() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSLIST_INT_GetActualValues tests the GetActualValues method for SLIST_INT
func TestSLIST_INT_GetActualValues(t *testing.T) {
	tests := []struct {
		name    string
		slist   SLIST_INT
		want    []int
		wantErr bool
	}{
		{
			name: "Basic calculation with origin 0 and scale 1",
			slist: SLIST_INT{
				Origin: 0,
				Scale:  1,
				Digits: "0 1 2 3 4",
			},
			want:    []int{0, 1, 2, 3, 4},
			wantErr: false,
		},
		{
			name: "Calculation with non-zero origin and scale",
			slist: SLIST_INT{
				Origin: 100,
				Scale:  10,
				Digits: "0 1 2 -1 -2",
			},
			want:    []int{100, 110, 120, 90, 80},
			wantErr: false,
		},
		{
			name: "Calculation with scale multiplier",
			slist: SLIST_INT{
				Origin: 0,
				Scale:  5,
				Digits: "0 2 4 -2 -4",
			},
			want:    []int{0, 10, 20, -10, -20},
			wantErr: false,
		},
		{
			name: "Empty digits",
			slist: SLIST_INT{
				Origin: 0,
				Scale:  1,
				Digits: "",
			},
			want:    []int{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.slist.GetActualValues()

			if (err != nil) != tt.wantErr {
				t.Errorf("SLIST_INT.GetActualValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SLIST_INT.GetActualValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestSLIST_PQ_XMLMarshal tests XML marshaling for SLIST_PQ
func TestSLIST_PQ_XMLMarshal(t *testing.T) {
	tests := []struct {
		name    string
		slist   SLIST_PQ
		wantXML string
	}{
		{
			name: "Complete SLIST_PQ",
			slist: SLIST_PQ{
				Origin: PhysicalQuantity{Value: "0", Unit: "uV"},
				Scale:  PhysicalQuantity{Value: "5", Unit: "uV"},
				Digits: "1 2 3 4 5",
			},
			wantXML: `<SLIST_PQ><origin value="0" unit="uV"></origin><scale value="5" unit="uV"></scale><digits>1 2 3 4 5</digits></SLIST_PQ>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := xml.Marshal(&tt.slist)
			if err != nil {
				t.Fatalf("xml.Marshal() error = %v", err)
			}

			got := string(data)
			if got != tt.wantXML {
				t.Errorf("xml.Marshal() = %v, want %v", got, tt.wantXML)
			}
		})
	}
}

// TestGLIST_TS_XMLMarshal tests XML marshaling for GLIST_TS
func TestGLIST_TS_XMLMarshal(t *testing.T) {
	tests := []struct {
		name    string
		glist   GLIST_TS
		wantXML string
	}{
		{
			name: "Complete GLIST_TS",
			glist: GLIST_TS{
				Head:      HeadTimestamp{Value: "20231223120000.000", Unit: "s"},
				Increment: Increment{Value: "0.002", Unit: "s"},
			},
			wantXML: `<GLIST_TS><head value="20231223120000.000" unit="s"></head><increment value="0.002" unit="s"></increment></GLIST_TS>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := xml.Marshal(&tt.glist)
			if err != nil {
				t.Fatalf("xml.Marshal() error = %v", err)
			}

			got := string(data)
			if got != tt.wantXML {
				t.Errorf("xml.Marshal() = %v, want %v", got, tt.wantXML)
			}
		})
	}
}

// TestSequenceValue_UnmarshalXML tests polymorphic unmarshaling of SequenceValue
func TestSequenceValue_UnmarshalXML(t *testing.T) {
	tests := []struct {
		name        string
		xmlData     string
		wantType    string
		wantErr     bool
		validateFn  func(*testing.T, *SequenceValue) // Custom validation function
	}{
		{
			name:     "Unmarshal GLIST_TS",
			xmlData:  `<value xsi:type="GLIST_TS"><head value="20231223120000.000" unit="s"/><increment value="0.002" unit="s"/></value>`,
			wantType: "GLIST_TS",
			wantErr:  false,
			validateFn: func(t *testing.T, sv *SequenceValue) {
				glist, ok := sv.Typed.(*GLIST_TS)
				if !ok {
					t.Errorf("Expected *GLIST_TS, got %T", sv.Typed)
					return
				}
				if glist.Head.Value != "20231223120000.000" {
					t.Errorf("Head.Value = %v, want %v", glist.Head.Value, "20231223120000.000")
				}
			},
		},
		{
			name:     "Unmarshal SLIST_PQ",
			xmlData:  `<value xsi:type="SLIST_PQ"><origin value="0" unit="uV"/><scale value="5" unit="uV"/><digits>1 2 3</digits></value>`,
			wantType: "SLIST_PQ",
			wantErr:  false,
			validateFn: func(t *testing.T, sv *SequenceValue) {
				slist, ok := sv.Typed.(*SLIST_PQ)
				if !ok {
					t.Errorf("Expected *SLIST_PQ, got %T", sv.Typed)
					return
				}
				if slist.Digits != "1 2 3" {
					t.Errorf("Digits = %v, want %v", slist.Digits, "1 2 3")
				}
			},
		},
		{
			name:     "Unmarshal SLIST_INT",
			xmlData:  `<value xsi:type="SLIST_INT"><origin>0</origin><scale>1</scale><digits>1 2 3</digits></value>`,
			wantType: "SLIST_INT",
			wantErr:  false,
			validateFn: func(t *testing.T, sv *SequenceValue) {
				slist, ok := sv.Typed.(*SLIST_INT)
				if !ok {
					t.Errorf("Expected *SLIST_INT, got %T", sv.Typed)
					return
				}
				if slist.Digits != "1 2 3" {
					t.Errorf("Digits = %v, want %v", slist.Digits, "1 2 3")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sv SequenceValue
			err := xml.Unmarshal([]byte(tt.xmlData), &sv)

			if (err != nil) != tt.wantErr {
				t.Errorf("SequenceValue.UnmarshalXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if sv.XsiType != tt.wantType {
					t.Errorf("XsiType = %v, want %v", sv.XsiType, tt.wantType)
				}

				if tt.validateFn != nil {
					tt.validateFn(t, &sv)
				}
			}
		})
	}
}
