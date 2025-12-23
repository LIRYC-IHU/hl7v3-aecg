package types

import (
	"encoding/xml"
	"testing"
)

// TestID_SetID tests the SetID method for ID type
func TestID_SetID(t *testing.T) {
	tests := []struct {
		name          string
		root          string
		extension     string
		wantRootEmpty bool // if root is empty, UUID should be generated
	}{
		{
			name:          "Set ID with both root and extension",
			root:          "2.16.840.1.113883.3.1234",
			extension:     "TEST-001",
			wantRootEmpty: false,
		},
		{
			name:          "Set ID with extension only (UUID should be generated)",
			root:          "",
			extension:     "TEST-002",
			wantRootEmpty: false, // UUID should be generated
		},
		{
			name:          "Set ID with root only",
			root:          "2.16.840.1.113883.3.1234",
			extension:     "",
			wantRootEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := &ID{}
			id.SetID(tt.root, tt.extension)

			if id.Root == "" && !tt.wantRootEmpty {
				t.Errorf("SetID() Root is empty, expected non-empty")
			}

			if tt.extension != "" && id.Extension != tt.extension {
				t.Errorf("SetID() Extension = %v, want %v", id.Extension, tt.extension)
			}

			// If root was empty, UUID should have been generated
			if tt.root == "" && id.Root == "" {
				t.Errorf("SetID() did not generate UUID when root was empty")
			}

			// If root was provided, it should be preserved
			if tt.root != "" && id.Root != tt.root {
				t.Errorf("SetID() Root = %v, want %v", id.Root, tt.root)
			}
		})
	}
}

// TestID_IsEmpty tests the IsEmpty method for ID type
func TestID_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		id   ID
		want bool
	}{
		{
			name: "Empty ID",
			id:   ID{},
			want: true,
		},
		{
			name: "ID with root only",
			id:   ID{Root: "2.16.840.1.113883.3.1234"},
			want: false,
		},
		{
			name: "ID with extension only",
			id:   ID{Extension: "TEST-001"},
			want: true, // Empty because Root is required
		},
		{
			name: "ID with both root and extension",
			id:   ID{Root: "2.16.840.1.113883.3.1234", Extension: "TEST-001"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.IsEmpty(); got != tt.want {
				t.Errorf("ID.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestID_String tests the String method for ID type
func TestID_String(t *testing.T) {
	tests := []struct {
		name string
		id   ID
		want string
	}{
		{
			name: "ID with both root and extension",
			id:   ID{Root: "2.16.840.1.113883.3.1234", Extension: "TEST-001"},
			want: "2.16.840.1.113883.3.1234^TEST-001",
		},
		{
			name: "ID with root only",
			id:   ID{Root: "2.16.840.1.113883.3.1234"},
			want: "2.16.840.1.113883.3.1234",
		},
		{
			name: "Empty ID",
			id:   ID{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.String(); got != tt.want {
				t.Errorf("ID.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCode_NewCode tests the NewCode constructor function
func TestCode_NewCode(t *testing.T) {
	tests := []struct {
		name            string
		code            string
		codeSystem      string
		displayName     string
		wantCode        string
		wantCodeSystem  string
		wantDisplayName string
	}{
		{
			name:            "Create code with all fields",
			code:            "93000",
			codeSystem:      "2.16.840.1.113883.6.12",
			displayName:     "Routine ECG",
			wantCode:        "93000",
			wantCodeSystem:  "2.16.840.1.113883.6.12",
			wantDisplayName: "Routine ECG",
		},
		{
			name:            "Create code without display name",
			code:            "RHYTHM",
			codeSystem:      "2.16.840.1.113883.5.4",
			displayName:     "",
			wantCode:        "RHYTHM",
			wantCodeSystem:  "2.16.840.1.113883.5.4",
			wantDisplayName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCode[string, string](tt.code, tt.codeSystem, tt.displayName)

			if got.Code != tt.wantCode {
				t.Errorf("NewCode() Code = %v, want %v", got.Code, tt.wantCode)
			}
			if got.CodeSystem != tt.wantCodeSystem {
				t.Errorf("NewCode() CodeSystem = %v, want %v", got.CodeSystem, tt.wantCodeSystem)
			}
			if got.DisplayName != tt.wantDisplayName {
				t.Errorf("NewCode() DisplayName = %v, want %v", got.DisplayName, tt.wantDisplayName)
			}
		})
	}
}

// TestCode_SetCode tests the SetCode method
func TestCode_SetCode(t *testing.T) {
	tests := []struct {
		name            string
		initialCode     *Code[string, string]
		code            string
		codeSystem      string
		displayName     string
		wantCode        string
		wantCodeSystem  string
		wantDisplayName string
	}{
		{
			name:            "Set code on empty code",
			initialCode:     &Code[string, string]{},
			code:            "M",
			codeSystem:      "2.16.840.1.113883.5.1",
			displayName:     "Male",
			wantCode:        "M",
			wantCodeSystem:  "2.16.840.1.113883.5.1",
			wantDisplayName: "Male",
		},
		{
			name: "Update existing code",
			initialCode: &Code[string, string]{
				Code:        "OLD",
				CodeSystem:  "OLD_SYSTEM",
				DisplayName: "Old Display",
			},
			code:            "NEW",
			codeSystem:      "NEW_SYSTEM",
			displayName:     "New Display",
			wantCode:        "NEW",
			wantCodeSystem:  "NEW_SYSTEM",
			wantDisplayName: "New Display",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initialCode.SetCode(tt.code, tt.codeSystem, tt.displayName)

			if tt.initialCode.Code != tt.wantCode {
				t.Errorf("SetCode() Code = %v, want %v", tt.initialCode.Code, tt.wantCode)
			}
			if tt.initialCode.CodeSystem != tt.wantCodeSystem {
				t.Errorf("SetCode() CodeSystem = %v, want %v", tt.initialCode.CodeSystem, tt.wantCodeSystem)
			}
			if tt.initialCode.DisplayName != tt.wantDisplayName {
				t.Errorf("SetCode() DisplayName = %v, want %v", tt.initialCode.DisplayName, tt.wantDisplayName)
			}
		})
	}
}

// TestCode_IsEmpty tests the IsEmpty method for Code type
func TestCode_IsEmpty(t *testing.T) {
	tests := []struct {
		name string
		code Code[string, string]
		want bool
	}{
		{
			name: "Empty code",
			code: Code[string, string]{},
			want: true,
		},
		{
			name: "Code with code value only",
			code: Code[string, string]{Code: "M"},
			want: false,
		},
		{
			name: "Code with all fields",
			code: Code[string, string]{
				Code:        "M",
				CodeSystem:  "2.16.840.1.113883.5.1",
				DisplayName: "Male",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.IsEmpty(); got != tt.want {
				t.Errorf("Code.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCode_String tests the String method for Code type
func TestCode_String(t *testing.T) {
	tests := []struct {
		name string
		code Code[string, string]
		want string
	}{
		{
			name: "Code with display name",
			code: Code[string, string]{
				Code:        "M",
				CodeSystem:  "2.16.840.1.113883.5.1",
				DisplayName: "Male",
			},
			want: "Male",
		},
		{
			name: "Code without display name",
			code: Code[string, string]{
				Code:       "M",
				CodeSystem: "2.16.840.1.113883.5.1",
			},
			want: "M",
		},
		{
			name: "Empty code",
			code: Code[string, string]{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.String(); got != tt.want {
				t.Errorf("Code.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestEffectiveTime_IsEmpty tests the IsEmpty method for EffectiveTime
func TestEffectiveTime_IsEmpty(t *testing.T) {
	tests := []struct {
		name          string
		effectiveTime EffectiveTime
		want          bool
	}{
		{
			name:          "Empty effective time",
			effectiveTime: EffectiveTime{},
			want:          true,
		},
		{
			name: "Effective time with low only",
			effectiveTime: EffectiveTime{
				Low: Time{Value: "20231223120000"},
			},
			want: false,
		},
		{
			name: "Effective time with high only",
			effectiveTime: EffectiveTime{
				High: Time{Value: "20231223120000"},
			},
			want: false,
		},
		{
			name: "Effective time with both low and high",
			effectiveTime: EffectiveTime{
				Low:  Time{Value: "20231223120000"},
				High: Time{Value: "20231223120010"},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.effectiveTime.IsEmpty(); got != tt.want {
				t.Errorf("EffectiveTime.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestID_XMLMarshal tests XML marshaling for ID type
func TestID_XMLMarshal(t *testing.T) {
	tests := []struct {
		name    string
		id      ID
		wantXML string
	}{
		{
			name: "ID with both root and extension",
			id: ID{
				Root:      "2.16.840.1.113883.3.1234",
				Extension: "TEST-001",
			},
			wantXML: `<ID root="2.16.840.1.113883.3.1234" extension="TEST-001"></ID>`,
		},
		{
			name: "ID with root only",
			id: ID{
				Root: "2.16.840.1.113883.3.1234",
			},
			wantXML: `<ID root="2.16.840.1.113883.3.1234"></ID>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := xml.Marshal(&tt.id)
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

// TestID_XMLUnmarshal tests XML unmarshaling for ID type
func TestID_XMLUnmarshal(t *testing.T) {
	tests := []struct {
		name      string
		xmlData   string
		wantRoot  string
		wantExt   string
		wantError bool
	}{
		{
			name:      "Unmarshal ID with both attributes",
			xmlData:   `<id root="2.16.840.1.113883.3.1234" extension="TEST-001"></id>`,
			wantRoot:  "2.16.840.1.113883.3.1234",
			wantExt:   "TEST-001",
			wantError: false,
		},
		{
			name:      "Unmarshal ID with root only",
			xmlData:   `<id root="2.16.840.1.113883.3.1234"></id>`,
			wantRoot:  "2.16.840.1.113883.3.1234",
			wantExt:   "",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id ID
			err := xml.Unmarshal([]byte(tt.xmlData), &id)

			if (err != nil) != tt.wantError {
				t.Errorf("xml.Unmarshal() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				if id.Root != tt.wantRoot {
					t.Errorf("xml.Unmarshal() Root = %v, want %v", id.Root, tt.wantRoot)
				}
				if id.Extension != tt.wantExt {
					t.Errorf("xml.Unmarshal() Extension = %v, want %v", id.Extension, tt.wantExt)
				}
			}
		})
	}
}

// TestCode_XMLMarshal tests XML marshaling for Code type
func TestCode_XMLMarshal(t *testing.T) {
	tests := []struct {
		name    string
		code    Code[string, string]
		wantXML string
	}{
		{
			name: "Code with all attributes",
			code: Code[string, string]{
				Code:        "93000",
				CodeSystem:  "2.16.840.1.113883.6.12",
				DisplayName: "Routine ECG",
			},
			wantXML: `<Code code="93000" codeSystem="2.16.840.1.113883.6.12" displayName="Routine ECG"></Code>`,
		},
		{
			name: "Code without display name",
			code: Code[string, string]{
				Code:       "RHYTHM",
				CodeSystem: "2.16.840.1.113883.5.4",
			},
			wantXML: `<Code code="RHYTHM" codeSystem="2.16.840.1.113883.5.4"></Code>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := xml.Marshal(&tt.code)
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

// TestTime_XMLMarshal tests XML marshaling for Time type
func TestTime_XMLMarshal(t *testing.T) {
	tests := []struct {
		name    string
		time    Time
		wantXML string
	}{
		{
			name:    "Time with HL7 timestamp",
			time:    Time{Value: "20231223120000"},
			wantXML: `<Time value="20231223120000"></Time>`,
		},
		{
			name:    "Time with milliseconds",
			time:    Time{Value: "20231223120000.123"},
			wantXML: `<Time value="20231223120000.123"></Time>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := xml.Marshal(&tt.time)
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
